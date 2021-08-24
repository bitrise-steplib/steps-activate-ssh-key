package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/stepenv"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/bitrise-steplib/steps-activate-ssh-key/step"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	fileWriter := fileutil.NewFileWriter()
	tempDirProvider := pathutil.NewTempDirProvider()
	envRepository := stepenv.NewRepository(env.NewRepository())
	cmdFactory := command.NewFactory(envRepository)
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)
	inputParser := stepconf.NewDefaultEnvParser()

	sshKeyActivator := step.NewActivateSSHKey(inputParser, envRepository, fileWriter, *agent, logger)

	config, err := sshKeyActivator.ProcessConfig()
	if err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	result, err := sshKeyActivator.Run(config)
	if err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	if err := sshKeyActivator.Export(result); err != nil {
		logger.Errorf(err.Error())
		return 1
	}

	return 0
}
