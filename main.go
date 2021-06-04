package main

import (
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/bitrise-steplib/steps-activate-ssh-key/step"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Errorf("Step run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	logger := *log.NewLogger()
	writer := *filewriter.NewOsFileWriter()
	activateSSHKey := step.newActivateSSHKey(step.newEnvStepInputParser(), step.newOsEnvManager(logger), step.newOsFileSSHKeyActivator(writer, *sshkey.NewAgent(writer, sshkey.NewOsTempDirProvider(), sshkey.NewOsCommandRunner(), logger), logger), logger)
	processConfig, err := activateSSHKey.processConfig()
	if err != nil {
		return err
	}
	result, err := activateSSHKey.run(processConfig)
	if err != nil {
		return err
	}
	if err := activateSSHKey.export(result); err != nil {
		return err
	}
	return nil
}
