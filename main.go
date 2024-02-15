package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-steputils/v2/stepenv"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/fileutil"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/bitrise-steplib/steps-activate-ssh-key/step"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	fileWriter := fileutil.NewFileManager()
	tempDirProvider := pathutil.NewPathProvider()
	envRepository := env.NewRepository()
	stepEnvRepository := stepenv.NewRepository(envRepository)
	cmdFactory := command.NewFactory(envRepository)
	agent := sshkey.NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)
	inputParser := stepconf.NewInputParser(envRepository)

	sshKeyActivator := step.NewActivateSSHKey(inputParser, stepEnvRepository, fileWriter, agent, logger)

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
