package activatesshkey

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
)

// Config is the activate SSH key step configuration
type Config struct {
	SSHRsaPrivateKey        stepconf.Secret `env:"ssh_rsa_private_key,required"`
	SSHKeySavePath          string          `env:"ssh_key_save_path,required"`
	IsRemoveOtherIdentities bool            `env:"is_remove_other_identities,required"`
	Verbose                 bool            `env:"verbose"`
}

// Execute activates a given SSH key
func Execute(cfg Config) error {
	// Remove SSHRsaPrivateKey from envs
	if err := unsetEnvsBy(string(cfg.SSHRsaPrivateKey)); err != nil {
		return &StepError{
			Tag:      "remove-private-key-data",
			ShortMsg: "Failed to remove private key data from envs",
			Err:      err,
		}
	}

	if err := ensureSavePath(cfg.SSHKeySavePath); err != nil {
		return &StepError{
			Tag:      "create-ssh-save-path",
			ShortMsg: "Failed to create the provided path",
			Err:      err,
		}
	}

	if err := fileutil.WriteStringToFile(cfg.SSHKeySavePath, string(cfg.SSHRsaPrivateKey)); err != nil {
		return &StepError{
			Tag:      "write-ssh-key",
			ShortMsg: "Failed to write the SSH key to the provided path",
			Err:      err,
		}
	}

	if err := os.Chmod(cfg.SSHKeySavePath, 0600); err != nil {
		return &StepError{
			Tag:      "change-ssh-key-permission",
			ShortMsg: "Failed to change file's access permission",
			Err:      err,
		}
	}

	if err := restartAgent(cfg.IsRemoveOtherIdentities); err != nil {
		return &StepError{
			Tag:      "restart-ssh-agent",
			ShortMsg: "Failed to restart SSH Agent",
			Err:      err,
		}
	}

	if err := checkPassphrase(cfg.SSHKeySavePath); err != nil {
		return &StepError{
			Tag:      "check-passphrase",
			ShortMsg: "SSH key requires passphrase",
			Err:      err,
		}
	}

	fmt.Println()
	log.Donef("Success")
	log.Printf("The SSH key was saved to %s", cfg.SSHKeySavePath)
	log.Printf("and was successfully added to ssh-agent.")

	return nil
}
