package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/step"
	"github.com/bitrise-io/go-steputils/stepconf"
	globallog "github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/env"
	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
)

// Input ...
type Input struct {
	SSHRsaPrivateKey        stepconf.Secret `env:"ssh_rsa_private_key,required"`
	SSHKeySavePath          string          `env:"ssh_key_save_path,required"`
	IsRemoveOtherIdentities bool            `env:"is_remove_other_identities,required"`
	Verbose                 bool            `env:"verbose"`
}

// Config ...
type Config struct {
	sshRsaPrivateKey        stepconf.Secret
	sshKeySavePath          string
	isRemoveOtherIdentities bool
	verbose                 bool
}

// Result ...
type Result struct {
	sshAuthSock string
}

// InputParser ...
type InputParser interface {
	Parse() (Input, error)
}

type envInputParser struct{}

// NewEnvInputParser ...
func NewEnvInputParser() InputParser {
	return envInputParser{}
}

// Parse ...
func (envInputParser) Parse() (Input, error) {
	var i Input
	if err := stepconf.Parse(&i); err != nil {
		return Input{}, err
	}
	return i, nil
}

// SSHKeyAgent ...
type SSHKeyAgent interface {
	Start() (string, error)
	Kill() (int, error)
	ListKeys() (int, error)
	AddKey(sshKeyPth, socket string) error
	DeleteKeys() error
}

// ActivateSSHKey ...
type ActivateSSHKey struct {
	stepInputParser InputParser
	envRepository   env.Repository
	fileWriter      filewriter.FileWriter
	sshKeyAgent     SSHKeyAgent
	logger          log.Logger
}

// NewActivateSSHKey ...
func NewActivateSSHKey(stepInputParse InputParser, envRepository env.Repository, fileWriter filewriter.FileWriter, agent SSHKeyAgent, logger log.Logger) *ActivateSSHKey {
	return &ActivateSSHKey{stepInputParser: stepInputParse, envRepository: envRepository, fileWriter: fileWriter, sshKeyAgent: agent, logger: logger}
}

// ProcessConfig ...
func (a ActivateSSHKey) ProcessConfig() (Config, error) {
	input, err := a.stepInputParser.Parse()
	if err != nil {
		return Config{}, err
	}
	stepconf.Print(input) // TODO: log.Infof(stepconf.toString(input))
	return Config{
		sshRsaPrivateKey:        input.SSHRsaPrivateKey,
		sshKeySavePath:          input.SSHKeySavePath,
		isRemoveOtherIdentities: input.IsRemoveOtherIdentities,
		verbose:                 input.Verbose,
	}, nil
}

// Run ...
func (a ActivateSSHKey) Run(cfg Config) (Result, error) {
	globallog.SetEnableDebugLog(cfg.verbose)
	if err := a.clearSSHKeys(string(cfg.sshRsaPrivateKey)); err != nil {
		return Result{}, err
	}
	output, err := a.activate(cfg.sshKeySavePath, string(cfg.sshRsaPrivateKey), cfg.isRemoveOtherIdentities)
	if err != nil {
		return Result{}, err
	}

	a.logger.Println()
	a.logger.Donef("Success")
	a.logger.Printf("The SSH key was saved to %s", cfg.sshKeySavePath)
	a.logger.Printf("and was successfully added to ssh-agent.")
	return Result{sshAuthSock: output}, nil
}

// Export ...
func (a ActivateSSHKey) Export(result Result) error {
	authSock := result.sshAuthSock
	if authSock == "" {
		return nil
	}
	if err := a.envRepository.Set("SSH_AUTH_SOCK", authSock); err != nil {
		return err
	}
	return nil
}

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

func (a ActivateSSHKey) clearSSHKeys(privateKey string) error {
	for _, env := range a.envRepository.List() {
		key, value := splitEnv(env)

		if value == privateKey {
			if err := a.envRepository.Unset(key); err != nil {
				return newStepError(
					"removing_private_key_data_failed",
					fmt.Errorf("failed to remove private key data from envs: %v", err),
					"Failed to remove private key data from envs",
				)
			}
		}
	}
	return nil
}

func (a ActivateSSHKey) activate(privateKeyPath, privateKey string, isRemoveOtherIdentities bool) (string, error) {
	// OpenSSH_8.1p1 on macOS requires a newline at the end of
	// private key using the new format (starting with -----BEGIN OPENSSH PRIVATE KEY-----).
	// See https://www.openssh.com/txt/release-7.8 for new format description.
	if err := a.fileWriter.Write(privateKeyPath, privateKey+"\n", 0600); err != nil {
		return "", newStepError(
			"writing_ssh_key_failed",
			fmt.Errorf("failed to write SSH key: %v", err),
			"Failed to write SSH key",
		)
	}

	socket, err := a.restartAgent(isRemoveOtherIdentities)
	if err != nil {
		return "", newStepError(
			"restarting_ssh_agent_failed",
			fmt.Errorf("failed to restart SSH Agent: %v", err),
			"Failed to restart SSH Agent",
		)
	}

	if err := a.sshKeyAgent.AddKey(privateKeyPath, socket); err != nil {
		return socket, newStepError(
			"ssh_key_requires_passphrase",
			fmt.Errorf("SSH key requires passphrase: %v", err),
			"SSH key requires passphrase",
		)
	}
	return socket, nil
}

func (a ActivateSSHKey) restartAgent(removeOtherIdentities bool) (string, error) {
	returnValue, err := a.sshKeyAgent.ListKeys()
	if err != nil {
		a.logger.Debugf("Exit code: %s", err)
	}

	shouldStartNewAgent := false
	//  as stated in the man page (https://developer.apple.com/library/mac/documentation/Darwin/Reference/ManPages/man1/ssh-add.1.html)
	//  ssh-add returns the exit code 2 if it could not connect to the ssh-agent
	if returnValue == 2 {
		a.logger.Printf("ssh_agent_check_result: %d", returnValue)
		a.logger.Printf("ssh-agent not started")
		shouldStartNewAgent = true
	} else {
		// ssh-agent loaded and accessible
		a.logger.Printf("ssh_agent_check_result: %d", returnValue)
		fmt.Printf("running / accessible ssh-agent detected")
		if removeOtherIdentities {
			if err := a.sshKeyAgent.DeleteKeys(); err != nil {
				return "", err
			}

			returnValue, err := a.sshKeyAgent.Kill()
			if err != nil {
				a.logger.Printf("Exit code: %s", err)
			}

			if returnValue == 0 {
				shouldStartNewAgent = true
			}
		}
	}

	if shouldStartNewAgent {
		returnValue, err := a.sshKeyAgent.Start()
		if err != nil {
			a.logger.Debugf("Exit code: %s", err)
		}

		a.logger.Printf("Expose SSH_AUTH_SOCK for the new ssh-agent, with envman")

		socket := strings.TrimPrefix(returnValue, "SSH_AUTH_SOCK=")
		socket = strings.Split(socket, ";")[0]

		return socket, nil
	}
	return "", nil
}

func newStepError(tag string, err error, shortMsg string) *step.Error {
	return step.NewError("activate-ssh-key", tag, err, shortMsg)
}
