package step

import (
	"testing"

	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const privateKey = "test-key"
const envKey = "env-key"

func Test_SSHPrivateKeyRemoved(t *testing.T) {
	osManager := preparePrePopulatedOsEnvManager()
	envmanManager := prepareDefaultEnvmanEnvManager()
	fileWriter := prepareDefaultFileWriter()
	logger := prepareDefaultLogger()
	tempDirProvider := prepareDefaultTempDirProvider()
	activateSSHKey := prepareActivateSSHKey(osManager, envmanManager, fileWriter, logger, tempDirProvider)
	config := getDefaultConfig()

	output, err := activateSSHKey.Run(config)

	assert.NoError(t, err)
	assert.Equal(t, output.sshAuthSock, "")
	osManager.AssertNumberOfCalls(t, "Unset", 1)
	osManager.AssertCalled(t, "Unset", envKey)
	envmanManager.AssertNumberOfCalls(t, "Unset", 1)
	envmanManager.AssertCalled(t, "Unset", envKey)
}

func prepareDefaultCommand() (command *MockCommand) {
	command = new(MockCommand)
	command.On("RunAndReturnExitCode", mock.Anything).Return(0, nil)
	command.On("RunAndReturnTrimmedOutput", mock.Anything).Return("", nil)
	command.On("Run", mock.Anything).Return(nil)
	command.On("PrintableCommandArgs").Return("")
	command.On("SetStdout", mock.Anything).Return(nil)
	command.On("SetStderr", mock.Anything).Return(nil)
	return
}

func prepareDefaultFileWriter() (fileWriter *mockFileWriter) {
	fileWriter = new(mockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return
}

func prepareDefaultLogger() (logger *mockLogger) {
	logger = new(mockLogger)
	logger.On("Debugf", mock.Anything, mock.Anything).Return()
	logger.On("Donef", mock.Anything, mock.Anything).Return()
	logger.On("Printf", mock.Anything, mock.Anything).Return()
	logger.On("Errorf", mock.Anything, mock.Anything).Return()
	logger.On("Println").Return()
	return
}

func prepareDefaultEnvmanEnvManager() (envManager *mockEnvManager) {
	envManager = new(mockEnvManager)
	envManager.On("Set", mock.Anything, mock.Anything).Return(nil)
	envManager.On("Unset", mock.Anything).Return(nil)
	return
}

func preparePrePopulatedOsEnvManager() (envManager *mockExtendedEnvManager) {
	envManager = new(mockExtendedEnvManager)
	envManager.On("Set", mock.Anything, mock.Anything).Return(nil)
	envManager.On("Unset", mock.Anything).Return(nil)
	envManager.On("List").Return([]string{envKey + "=" + privateKey})
	return
}

func prepareDefaultTempDirProvider() (tempDirProvider *mockTempDirProvider) {
	tempDirProvider = new(mockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return("temp-dir", nil)
	return
}

func getDefaultConfig() Config {

	return Config{
		sshRsaPrivateKey:        privateKey,
		sshKeySavePath:          "test-path",
		isRemoveOtherIdentities: false,
		verbose:                 false,
	}
}

func prepareActivateSSHKey(osManager *mockExtendedEnvManager, envmanManager *mockEnvManager, writer *mockFileWriter, logger *mockLogger, provider *mockTempDirProvider) *ActivateSSHKey {
	return &ActivateSSHKey{
		stepInputParse:      nil,
		envmanEnvRepository: envmanManager,
		osEnvRepository:     osManager,
		envValueClearer:     *NewCombinedEnvValueClearer(logger, osManager, envmanManager),
		fileWriter:          writer,
		agent: *sshkey.NewAgent(writer, provider, logger, func(name string, args ...string) command.Command {
			var c command.Command = prepareDefaultCommand()
			return c
		}),
		logger: logger,
	}
}
