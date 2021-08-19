package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/stepenv"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/bitrise-steplib/steps-activate-ssh-key/step"
)

// TODO debug log config? os.GetEnv? late init?
var logger = log.NewLogger(false)

func main() {
	if err := run(); err != nil {
		logger.Errorf("Step run failed: %s", err.Error())
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
	fileWriter := fileutil.NewFileWriter()
	tempDirProvider := pathutil.NewTempDirProvider()
	envRepository := stepenv.NewRepository(env.NewRepository())
	cmdFactory := command.NewFactory(envRepository)
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)
	stepInputParser := step.NewEnvInputParser()
	return step.NewActivateSSHKey(stepInputParser, envRepository, fileWriter, *agent, logger)
}
