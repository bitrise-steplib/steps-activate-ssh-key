package main

import (
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	localLogger "github.com/bitrise-steplib/steps-activate-ssh-key/log"
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
	logger := *localLogger.NewLogger()
	writer := *filewriter.NewOsFileWriter()
	activateSSHKey := step.NewActivateSSHKey(step.NewEnvStepInputParser(), step.NewOsEnvManager(logger), writer, *sshkey.NewAgent(writer, sshkey.NewOsTempDirProvider(), sshkey.NewOsCommandRunner(), logger), logger)
	processConfig, err := activateSSHKey.ProcessConfig()
	if err != nil {
		return err
	}
	result, err := activateSSHKey.Run(processConfig)
	if err != nil {
		return err
	}
	if err := activateSSHKey.Export(result); err != nil {
		return err
	}
	return nil
}
