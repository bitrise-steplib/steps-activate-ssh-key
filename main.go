package main

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

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
	deployDir := os.Getenv("BITRISE_DEPLOY_DIR")

	cpuProfilePth := filepath.Join(deployDir, "cpu.prof")
	cpuProfile, err := os.Create(cpuProfilePth)
	if err != nil {
		panic("could not create CPU profile: " + err.Error())
	}
	defer cpuProfile.Close()
	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		panic("could not start CPU profile: " + err.Error())
	}
	defer pprof.StopCPUProfile()

	memProfilePth := filepath.Join(deployDir, "mem.prof")
	memProfile, err := os.Create(memProfilePth)
	if err != nil {
		panic("could not create memory profile: " + err.Error())
	}
	defer memProfile.Close()

	runtime.GC()

	exitCode := run()

	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		panic("could not write memory profile: " + err.Error())
	}

	os.Exit(exitCode)
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
