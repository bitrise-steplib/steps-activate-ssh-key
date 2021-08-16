package main

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	localLogger "github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/pathutil"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/bitrise-steplib/steps-activate-ssh-key/step"
)

func main() {
	if err := run(); err != nil {
		log.Errorf("Step run failed: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	activateSSHKey := createActivateSSHKey()
	config, err := activateSSHKey.ProcessConfig()
	if err != nil {
		return err
	}
	result, err := activateSSHKey.Run(config)
	if err != nil {
		return err
	}
	if err := activateSSHKey.Export(result); err != nil {
		return err
	}
	return nil
}

func createActivateSSHKey() *step.ActivateSSHKey {
	stepInputParser := step.NewEnvInputParser()

	logger := localLogger.NewDefaultLogger()
	envManagerOrchestrator := step.DefaultEnvironmentManagerOrchestrator{}

	fileWriter := filewriter.NewOsFileWriter()
	tempDirProvider := pathutil.NewOsTempDirProvider()
	commandFactory := command.NewCommand
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, commandFactory)

	return step.NewActivateSSHKey(stepInputParser, envManagerOrchestrator, fileWriter, *agent, logger)
}
