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
	logger := localLogger.NewDefaultLogger()

	fileWriter := filewriter.NewOsFileWriter()
	tempDirProvider := pathutil.NewOsTempDirProvider()
	osEnvRepository := env.NewOsRepository()
	cmdFactory := command.NewDefaultFactory(osEnvRepository)
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)

	stepInputParser := step.NewEnvInputParser()
	envmanEnvRepository := env.NewEnvmanRepository()
	envRepository := env.NewRepository(osEnvRepository, envmanEnvRepository)

	return step.NewActivateSSHKey(stepInputParser, envRepository, fileWriter, *agent, logger)
}
