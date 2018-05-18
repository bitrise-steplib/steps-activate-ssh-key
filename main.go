package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
)

// Config ...
type Config struct {
	SSHRsaPrivateKey        stepconf.Secret `env:"ssh_rsa_private_key,required"`
	SSHKeySavePath          string          `env:"ssh_key_save_path"`
	IsRemoveOtherIdentities bool            `env:"is_remove_other_identities,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	// Remove SSHRsaPrivateKey from envs
	if err := os.Setenv("ssh_rsa_private_key", ""); err != nil {
		failf("* Failed to remove ssh_rsa_private_key")
	}

	if err := ensureSavePath(cfg.SSHKeySavePath); err != nil {
		failf("* Failed to create the provided path, %s", err)
	}

	if err := writeSSHKey(cfg.SSHKeySavePath, string(cfg.SSHRsaPrivateKey)); err != nil {
		failf("* Failed to write the SSH key to the provided path, %s", err)
	}

	if err := os.Chmod(cfg.SSHKeySavePath, 0600); err != nil {
		failf("* Failed to change file's access permission, %s", err)
	}

	if err := restartAgent(cfg.IsRemoveOtherIdentities); err != nil {
		failf("* Failed to restart SSH Agent, %s", err)
	}

	if err := checkPassphrase(cfg.SSHKeySavePath); err != nil {
		fmt.Println()
		failf("Error, %s", err)
	}

	fmt.Println()
	log.Printf("# Success")
	log.Printf("The SSH key was saved to *%s*", cfg.SSHKeySavePath)
	log.Printf("and was successfully added to ssh-agent.")
}

func ensureSavePath(savePath string) error {
	dirpath := filepath.Dir(savePath)
	return os.MkdirAll(dirpath, 0600)
}

func writeSSHKey(savePath string, SSHRsaPrivateKey string) error {
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}

	defer func() {
		cerr := f.Close()
		if cerr != nil {
			log.Errorf("Failed to close the file, err: %s", cerr)
		}
	}()

	_, err = f.WriteString(SSHRsaPrivateKey)
	return err
}

func checkPassphrase(savePath string) error {
	cmd := command.New("spawn", "ssh-add", savePath)
	fmt.Println()
	log.Printf("-> %s", cmd.PrintableCommandArgs())

	if returnValue, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return err
	} else if strings.HasPrefix(returnValue, "Identity added") {
		log.Printf(returnValue)
		return nil
	} else {
		log.Printf(returnValue)
		log.Warnf("returnValue: %s", returnValue)

	}

	return fmt.Errorf("* Failed to add the SSH key to ssh-agent with an empty passphrase")
}

func restartAgent(removeOtherIdentities bool) error {
	var shouldStartNewAgent bool
	cmd := command.New("ssh-add", "-l")
	cmd.SetStderr(os.Stderr)
	log.Printf("-> %s", cmd.PrintableCommandArgs())

	returnValue, err := cmd.RunAndReturnExitCode()
	if err != nil {
		log.Debugf("Exit code: %s", err)
	}

	if returnValue == 2 {
		log.Printf(" (i) ssh_agent_check_result: %d", returnValue)
		log.Printf(" (i) ssh-agent not started")
		shouldStartNewAgent = true
	} else {
		// ssh-agent loaded and accessible
		log.Printf(" (i) ssh_agent_check_result: %d", returnValue)
		fmt.Printf(" (i) running / accessible ssh-agent detected")
		if removeOtherIdentities {
			// remove all keys from the current agent
			cmdRemove := command.New("ssh-add", "-D")
			cmdRemove.SetStdout(os.Stdout)
			cmdRemove.SetStderr(os.Stderr)

			fmt.Println()
			fmt.Println()
			log.Printf("-> %s", cmdRemove.PrintableCommandArgs())

			if err := cmdRemove.Run(); err != nil {
				return err
			}

			// try to kill the agent
			cmdKill := command.New("ssh-agent", "-k")
			cmdKill.SetStdout(os.Stdout)
			cmdKill.SetStderr(os.Stderr)

			fmt.Println()
			log.Printf("-> %s", cmdKill.PrintableCommandArgs())

			returnValue, err := cmdKill.RunAndReturnExitCode()
			if err != nil {
				log.Debugf("Exit code: %s", err)
			}

			if returnValue == 0 {
				shouldStartNewAgent = true
			}
		}
	}

	if shouldStartNewAgent {
		fmt.Printf(" (i) starting a new ssh-agent and exporting connection information with envman")
		cmd := command.New("ssh-agent")
		cmd.SetStderr(os.Stderr)
		log.Printf("-> %s", cmd.PrintableCommandArgs())

		returnValue, err := cmd.RunAndReturnExitCode()
		if err != nil {
			log.Debugf("Exit code: %s", err)
		}

		if returnValue != 0 {
			return fmt.Errorf("[!] Failed to load SSH agent")
		}

		SSHAuthSock := os.Getenv("SSH_AUTH_SOCK")
		fmt.Printf(" (i) Expose SSH_AUTH_SOCK for the new ssh-agent, with envman")
		fmt.Println()
		log.Printf("-> export SSH_AUTH_SOCK: %s", SSHAuthSock)

		if err := tools.ExportEnvironmentWithEnvman("SSH_AUTH_SOCK", SSHAuthSock); err != nil {
			failf("Failed to generate output")
		}
	}
	return nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
