package main

import (
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
	step := createStep()
	config, err := step.ProcessConfig()
	if err != nil {
		return err
	}
	result, err := step.Run(config)
	if err != nil {
		return err
	}
	if err := step.Export(result); err != nil {
		return err
	}
	return nil
}

func createStep() *step.ActivateSSHKey {
	logger := localLogger.NewLogger()
	writer := filewriter.NewOsFileWriter()
	osEnvRepository := env.NewOsEnvManager()
	envmanEnvRepository := env.NewEnvmanEnvManager()
	stepInputParser := step.NewEnvStepInputParser()
	combinedEnvValueClearer := step.NewCombinedEnvValueClearer(logger, osEnvRepository, envmanEnvRepository)
	tempDirProvider := sshkey.NewOsTempDirProvider()
	cmdFactory := command.NewCommand
	agent := sshkey.NewAgent(writer, tempDirProvider, logger, cmdFactory)

	return step.NewActivateSSHKey(stepInputParser, *combinedEnvValueClearer, envmanEnvRepository, osEnvRepository, writer, *agent, logger)
}
