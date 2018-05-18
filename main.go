package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	SSHRsaPrivateKey        string `env:"ssh_rsa_private_key,required"`
	SSHKeySavePath          string `env:"ssh_key_save_path"`
	IsRemoveOtherIdentities bool   `env:"is_remove_other_identities,required"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	log.Infof("# Configuration")
	log.Printf("* Path to save the RSA SSH private key: %s", cfg.SSHKeySavePath)
	log.Printf("* Should remove other identities from the ssh-agent? %t\n", cfg.IsRemoveOtherIdentities)

	if err := ensureSavePath(cfg.SSHKeySavePath); err != nil {
		failf("* Failed to create the provided path, %s", err)
	}

	if err := writheSSHKey(cfg.SSHKeySavePath, cfg.SSHRsaPrivateKey); err != nil {
		failf("* Failed to write the SSH key to the provided path, %s", err)
	}
}

func ensureSavePath(savePath string) error {
	dirpath := filepath.Dir(savePath)

	cmd := command.New("mkdir", "-p", dirpath)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	return cmd.Run()
}

func writheSSHKey(savePath string, SSHRsaPrivateKey string) error {
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(SSHRsaPrivateKey)
	return err
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
