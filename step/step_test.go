package step

import (
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_SSHPrivateKeyRemoved(t *testing.T) {
	envManager := prepareDefaultEnvManager()
	fileWriter := prepareDefaultFileWriter()
	logger := prepareDefaultLogger()
	commandRunner := prepareDefaultCommandRunner()
	tempDirProvider := prepareDefaultTempDirProvider()
	activateSSHKey := prepareActivateSSHKey(envManager, fileWriter, logger, commandRunner, tempDirProvider)
	config := getDefaultConfig()

	output, err := activateSSHKey.Run(config)

	assert.NoError(t, err)
	assert.Equal(t, output.sshAuthSock, "")
	envManager.AssertNumberOfCalls(t, "UnsetByValue", 1)
	envManager.AssertCalled(t, "UnsetByValue", "test-key")
}

func prepareDefaultCommandRunner() (commandRunner *mockCommandRunner) {
	commandRunner = new(mockCommandRunner)
	commandRunner.On("RunAndReturnExitCode", mock.Anything).Return(0, nil)
	commandRunner.On("RunAndReturnTrimmedOutput", mock.Anything).Return("", nil)
	commandRunner.On("Run", mock.Anything).Return(nil)
	return commandRunner
}

func prepareDefaultFileWriter() (fileWriter *mockFileWriter) {
	fileWriter = new(mockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return fileWriter
}

func prepareDefaultLogger() (logger *mockLogger) {
	logger = new(mockLogger)
	logger.On("Debugf", mock.Anything, mock.Anything).Return()
	logger.On("Donef", mock.Anything, mock.Anything).Return()
	logger.On("Printf", mock.Anything, mock.Anything).Return()
	logger.On("Errorf", mock.Anything, mock.Anything).Return()
	logger.On("Println").Return()
	return logger
}

func prepareDefaultEnvManager() (envManager *mockEnvManager) {
	envManager = new(mockEnvManager)
	envManager.On("UnsetByValue", mock.Anything).Return(nil)
	return envManager
}

func getDefaultConfig() Config {
	return Config{
		sshRsaPrivateKey:        "test-key",
		sshKeySavePath:          "test-path",
		isRemoveOtherIdentities: false,
		verbose:                 false,
	}
}

func prepareDefaultTempDirProvider() (tempDirProvider *mockTempDirProvider) {
	tempDirProvider = new(mockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return("temp-dir", nil)
	return tempDirProvider
}

func prepareActivateSSHKey(manager *mockEnvManager, writer *mockFileWriter, logger *mockLogger, runner *mockCommandRunner, provider *mockTempDirProvider) *ActivateSSHKey {
	return &ActivateSSHKey{
		stepInputParse: nil,
		envManager:     manager,
		fileWriter:     writer,
		agent:          *sshkey.NewAgent(writer, provider, runner, logger),
		logger:         logger,
	}
}
