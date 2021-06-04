package step

import (
	"fmt"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"os"
	"strings"
)

type input struct {
	SSHRsaPrivateKey        stepconf.Secret `env:"ssh_rsa_private_key,required"`
	SSHKeySavePath          string          `env:"ssh_key_save_path,required"`
	IsRemoveOtherIdentities bool            `env:"is_remove_other_identities,required"`
	Verbose                 bool            `env:"verbose"`
}

type config struct {
	sshRsaPrivateKey        stepconf.Secret
	sshKeySavePath          string
	isRemoveOtherIdentities bool
	verbose                 bool
}

type result struct {
	sshAuthSock string
}

type fileWriter interface {
	Write(path string, value string, mode os.FileMode) error
}

type envManager interface {
	unsetByValue(value string) error
	set(key string, value string) error
}

type stepInputParser interface {
	parse() (input, error)
}

type envStepInputParser struct{}

func newEnvStepInputParser() *envStepInputParser {
	return &envStepInputParser{}
}

func (envStepInputParser) parse() (input, error) {
	var i input
	if err := stepconf.Parse(&i); err != nil {
		return input{}, err
	}
	return i, nil
}

type activateSSHKey struct {
	stepInputParse      stepInputParser
	envManager          envManager
	fileWriter          fileWriter
	agent               sshkey.Agent
	logger              log.Logger
}

func newActivateSSHKey(stepInputParse stepInputParser, envManager envManager, fileWriter fileWriter, agent sshkey.Agent, logger log.Logger) *activateSSHKey {
	return &activateSSHKey{stepInputParse: stepInputParse, envManager: envManager, fileWriter: fileWriter, agent: agent, logger: logger}
}

func (a activateSSHKey) processConfig() (config, error) {
	input, err := a.stepInputParse.parse()
	if err != nil {
		return config{}, err
	}
	stepconf.Print(input) // TODO: log.Infof(stepconf.toString(input))
	return config{
		sshRsaPrivateKey:        input.SSHRsaPrivateKey,
		sshKeySavePath:          input.SSHKeySavePath,
		isRemoveOtherIdentities: input.IsRemoveOtherIdentities,
		verbose:                 input.Verbose,
	}, nil
}

func (a activateSSHKey) run(cfg config) (result, error) {
	if err := a.clearSSHKeys(string(cfg.sshRsaPrivateKey)); err != nil {
		return result{}, err
	}
	output, err := a.activate(cfg.sshKeySavePath, string(cfg.sshRsaPrivateKey), cfg.isRemoveOtherIdentities)
	if err != nil {
		return result{}, err
	}

	a.logger.Println()
	a.logger.Donef("Success")
	a.logger.Printf("The SSH key was saved to %s", cfg.sshKeySavePath)
	a.logger.Printf("and was successfully added to ssh-agent.")
	return result{sshAuthSock: output}, nil
}

func (a activateSSHKey) export(result result) error {
	authSock := result.sshAuthSock
	if len(authSock) < 1 {
		return nil
	}
	if err := a.envManager.set("SSH_AUTH_SOCK", authSock); err != nil {
		return err
	}
	return nil
}

func (a activateSSHKey) clearSSHKeys(privateKey string) error {
	if err := a.envManager.unsetByValue(privateKey); err != nil {
		return newStepError(
			"removing_private_key_data_failed",
			fmt.Errorf("failed to remove private key data from envs: %v", err),
			"Failed to remove private key data from envs",
		)
	}
	return nil
}

func (a activateSSHKey) activate(path string, privateKey string, isRemoveOtherIdentities bool) (string, error) {
	// OpenSSH_8.1p1 on macOS requires a newline at at the end of
	// private key using the new format (starting with -----BEGIN OPENSSH PRIVATE KEY-----).
	// See https://www.openssh.com/txt/release-7.8 for new format description.
	if err := a.fileWriter.Write(path, privateKey+"\n", 0600); err != nil {
		return "", newStepError(
			"writing_ssh_key_failed",
			fmt.Errorf("failed to write SSH key: %v", err),
			"Failed to write SSH key",
		)
	}

	result, err := a.restartAgent(isRemoveOtherIdentities)
	if err != nil {
		return "", newStepError(
			"restarting_ssh_agent_failed",
			fmt.Errorf("failed to restart SSH Agent: %v", err),
			"Failed to restart SSH Agent",
		)
	}

	if err := a.agent.AddKey(path); err != nil {
		return result, newStepError(
			"ssh_key_requries_passphrase",
			fmt.Errorf("SSH key requires passphrase: %v", err),
			"SSH key requires passphrase",
		)
	}
	return result, nil
}

func (a activateSSHKey) restartAgent(removeOtherIdentities bool) (string, error) {
	var shouldStartNewAgent bool
	returnValue, err := a.agent.ListKeys()
	if err != nil {
		a.logger.Debugf("Exit code: %s", err)
	}

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
			if err := a.agent.DeleteKeys(); err != nil {
				return "", err
			}

			returnValue, err := a.agent.Kill()
			if err != nil {
				a.logger.Printf("Exit code: %s", err)
			}

			if returnValue == 0 {
				shouldStartNewAgent = true
			}
		}
	}

	if shouldStartNewAgent {
		returnValue, err := a.agent.Start()
		if err != nil {
			a.logger.Debugf("Exit code: %s", err)
		}

		a.logger.Printf("Expose SSH_AUTH_SOCK for the new ssh-agent, with envman")

		returnValue = strings.TrimPrefix(returnValue, "SSH_AUTH_SOCK=")
		returnValue = strings.Split(returnValue, ";")[0]

		return returnValue, nil
	}
	return "", nil
}