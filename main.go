package main

import (
	"os"

	utilcommand "github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/envmanager"
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

func createStep() *step.ActivateSSHKey {
	logger := *localLogger.NewLogger()
	writer := *filewriter.NewOsFileWriter()
	osEnvManager := *envmanager.NewOsEnvManager()
	envmanEnvManager := *envmanager.NewEnvmanEnvManager()
	return step.NewActivateSSHKey(step.NewEnvStepInputParser(), *step.NewCombinedEnvValueClearer(logger, osEnvManager, envmanEnvManager), envmanEnvManager, osEnvManager, writer, *sshkey.NewAgent(writer, sshkey.NewOsTempDirProvider(), logger, func(name string, args ...string) *sshkey.Command {
		var c sshkey.Command
		c = utilcommand.New(name, args...)
		return &c
	}), logger)
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
