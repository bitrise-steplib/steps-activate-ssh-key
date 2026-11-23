package step

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/stretchr/testify/assert"
)

func TestActivateSSHKey_Integration(t *testing.T) {
	for _, tool := range []string{"ssh-agent", "ssh-add", "ssh-keygen", "expect"} {
		if _, err := exec.LookPath(tool); err != nil {
			// A skip would look like a pass in CI and hide that this test never ran.
			if os.Getenv("CI") != "" {
				t.Fatalf("%s not found in PATH", tool)
			}
			t.Skipf("%s not found in PATH", tool)
		}
	}

	tests := []struct {
		name    string
		key     func(t *testing.T) (privateKey, fingerprint string)
		wantErr bool
	}{
		{
			name: "RSA key in PEM format",
			key:  func(t *testing.T) (string, string) { return generateKey(t, "-t", "rsa", "-b", "2048", "-m", "PEM") },
		},
		{
			name: "RSA key in OpenSSH format",
			key:  func(t *testing.T) (string, string) { return generateKey(t, "-t", "rsa", "-b", "2048") },
		},
		{
			name: "Ed25519 key",
			key:  func(t *testing.T) (string, string) { return generateKey(t, "-t", "ed25519") },
		},
		{
			name: "OpenSSH format key without trailing newline",
			key: func(t *testing.T) (string, string) {
				key, fingerprint := generateKey(t, "-t", "ed25519")
				return strings.TrimRight(key, "\n"), fingerprint
			},
		},
		{
			name: "passphrase protected key",
			key: func(t *testing.T) (string, string) {
				return generateKey(t, "-t", "ed25519", "-N", "secret")
			},
			wantErr: true,
		},
		{
			name:    "invalid key",
			key:     func(t *testing.T) (string, string) { return "not a private key", "" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startSSHAgent(t)
			privateKey, fingerprint := tt.key(t)
			keyPath := filepath.Join(t.TempDir(), "id")

			step := newIntegrationStep()
			_, err := step.Run(Config{
				sshRsaPrivateKey:        stepconf.Secret(privateKey),
				sshKeySavePath:          keyPath,
				isRemoveOtherIdentities: false,
			})

			if tt.wantErr {
				assert.Error(t, err)
				if fingerprint != "" {
					assert.NotContains(t, listAgentKeys(t), fingerprint)
				}
				return
			}
			assert.NoError(t, err)
			assert.FileExists(t, keyPath)
			assert.Contains(t, listAgentKeys(t), fingerprint)
		})
	}
}

func newIntegrationStep() *ActivateSSHKey {
	logger := log.NewLogger()
	fileManager := fileutil.NewFileManager()
	envRepository := env.NewRepository()
	agent := sshkey.NewAgent(fileManager, pathutil.NewPathProvider(), logger, command.NewFactory(envRepository))
	return NewActivateSSHKey(nil, envRepository, fileManager, agent, logger)
}

// startSSHAgent runs a dedicated agent for the test so it never touches the developer's own agent and keys.
func startSSHAgent(t *testing.T) {
	// Unix socket paths are limited to ~104 bytes on macOS, which t.TempDir() can exceed.
	dir, err := os.MkdirTemp("", "ssh")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	socket := filepath.Join(dir, "agent.sock")

	out, err := exec.Command("ssh-agent", "-s", "-a", socket).Output()
	if err != nil {
		t.Fatalf("ssh-agent: %v", err)
	}

	match := regexp.MustCompile(`SSH_AGENT_PID=(\d+)`).FindStringSubmatch(string(out))
	if match == nil {
		t.Fatalf("unexpected ssh-agent output: %s", out)
	}
	t.Cleanup(func() {
		_ = exec.Command("kill", match[1]).Run()
	})

	t.Setenv("SSH_AUTH_SOCK", socket)
	t.Setenv("SSH_AGENT_PID", match[1])
}

func generateKey(t *testing.T, args ...string) (privateKey, fingerprint string) {
	path := filepath.Join(t.TempDir(), "key")
	args = append([]string{"-q", "-f", path, "-C", "test"}, args...)
	if !contains(args, "-N") {
		args = append(args, "-N", "")
	}
	if out, err := exec.Command("ssh-keygen", args...).CombinedOutput(); err != nil {
		t.Fatalf("ssh-keygen: %v: %s", err, out)
	}

	key, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("ssh-keygen", "-l", "-f", path+".pub").Output()
	if err != nil {
		t.Fatalf("ssh-keygen -l: %v", err)
	}
	return string(key), strings.Fields(string(out))[1]
}

func listAgentKeys(t *testing.T) string {
	out, err := exec.Command("ssh-add", "-l").CombinedOutput()
	var exitErr *exec.ExitError
	// Exit code 1 means the agent has no identities.
	if err != nil && !(errors.As(err, &exitErr) && exitErr.ExitCode() == 1) {
		t.Fatalf("ssh-add -l: %v: %s", err, out)
	}
	return string(out)
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
