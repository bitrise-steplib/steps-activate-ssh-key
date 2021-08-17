package main

import (
	"github.com/bitrise-steplib/steps-activate-ssh-key/pathutil"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/env"
	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	localLogger "github.com/bitrise-steplib/steps-activate-ssh-key/log"
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
	osEnvRepository := env.NewOsRepository()
	envmanEnvRepository := env.NewEnvmanRepository()
	envValueClearer := step.NewCombinedEnvValueClearer(logger, osEnvRepository, envmanEnvRepository)

	fileWriter := filewriter.NewOsFileWriter()
	tempDirProvider := pathutil.NewOsTempDirProvider()
	cmdFactory := command.NewDefaultFactory()
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)

	return step.NewActivateSSHKey(stepInputParser, *envValueClearer, envmanEnvRepository, osEnvRepository, fileWriter, *agent, logger)
}
