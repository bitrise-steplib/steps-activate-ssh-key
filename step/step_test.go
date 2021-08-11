package step

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const privateKey = "test-key"
const envKey = "env-key"

func Test_SSHPrivateKeyRemoved(t *testing.T) {
	osEnvRepository := createOsEnvRepositoryWithSSHKey()
	osEnvManager := createEnvmanEnvRepository()
	fileWriter := createFileWriter()
	logger := createLogger()
	tempDirProvider := createTempProvider()
	mockCommand := createCommand()
	activateSSHKey := createActivateSSHKey(osEnvRepository, osEnvManager, fileWriter, logger, tempDirProvider, mockCommand)
	config := createConfigWithDefaults()

	output, err := activateSSHKey.Run(config)

	assert.NoError(t, err)
	assert.Equal(t, output.sshAuthSock, "")
	osEnvRepository.AssertNumberOfCalls(t, "Unset", 1)
	osEnvRepository.AssertCalled(t, "Unset", envKey)
	osEnvManager.AssertNumberOfCalls(t, "Unset", 1)
	osEnvManager.AssertCalled(t, "Unset", envKey)
}

func Test_ErrorRaisedIfSSHAddFails(t *testing.T) {
	osEnvRepository := createOsEnvRepositoryWithSSHKey()
	osEnvManager := createEnvmanEnvRepository()
	fileWriter := createFileWriter()
	logger := createLogger()
	tempDirProvider := createTempProvider()

	mockCommand := new(MockCommand)
	mockCommand.On("RunAndReturnExitCode", mock.Anything).Return(1, errors.New("mocked error"))
	mockCommand.On("RunAndReturnTrimmedOutput", mock.Anything).Return("", nil)
	mockCommand.On("Run", mock.Anything).Return(nil)
	mockCommand.On("PrintableCommandArgs").Return("")
	mockCommand.On("SetStdout", mock.Anything).Return(nil)
	mockCommand.On("SetStderr", mock.Anything).Return(nil)

	activateSSHKey := createActivateSSHKey(osEnvRepository, osEnvManager, fileWriter, logger, tempDirProvider, mockCommand)
	config := createConfigWithDefaults()

	output, err := activateSSHKey.Run(config)

	wantOutput := Result{sshAuthSock: ""}
	wantErr := newStepError(
		"ssh_key_requires_passphrase",
		fmt.Errorf("SSH key requires passphrase: %v", errors.New("failed to add the SSH key to ssh-agent with an empty passphrase")),
		"SSH key requires passphrase",
	)
	assert.Equal(t, wantOutput, output)
	assert.Error(t, err)
	assert.Equal(t, wantErr, err)
}

func createOsEnvRepositoryWithSSHKey() (osEnvRepository *MockOsRepository) {
	osEnvRepository = new(MockOsRepository)
	osEnvRepository.On("Set", mock.Anything, mock.Anything).Return(nil)
	osEnvRepository.On("Unset", mock.Anything).Return(nil)
	osEnvRepository.On("List").Return([]string{envKey + "=" + privateKey})
	return
}

func createEnvmanEnvRepository() (envmanEnvRepository *MockEnvmanRepository) {
	envmanEnvRepository = new(MockEnvmanRepository)
	envmanEnvRepository.On("Set", mock.Anything, mock.Anything).Return(nil)
	envmanEnvRepository.On("Unset", mock.Anything).Return(nil)
	return
}

func createFileWriter() (fileWriter *MockFileWriter) {
	fileWriter = new(MockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return
}

func createLogger() (logger *MockLogger) {
	logger = new(MockLogger)
	logger.On("Debugf", mock.Anything, mock.Anything).Return()
	logger.On("Donef", mock.Anything, mock.Anything).Return()
	logger.On("Printf", mock.Anything, mock.Anything).Return()
	logger.On("Errorf", mock.Anything, mock.Anything).Return()
	logger.On("Println").Return()
	return
}

func createTempProvider() (tempDirProvider *MockTempDirProvider) {
	tempDirProvider = new(MockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return("temp-dir", nil)
	return
}

func createCommand() (command *MockCommand) {
	command = new(MockCommand)
	command.On("RunAndReturnExitCode", mock.Anything).Return(0, nil)
	command.On("RunAndReturnTrimmedOutput", mock.Anything).Return("", nil)
	command.On("Run", mock.Anything).Return(nil)
	command.On("PrintableCommandArgs").Return("")
	command.On("SetStdout", mock.Anything).Return(nil)
	command.On("SetStderr", mock.Anything).Return(nil)
	return
}

func createActivateSSHKey(osEnvRepository *MockOsRepository, envmanEnvRepository *MockEnvmanRepository, fileWriter *MockFileWriter, logger *MockLogger, tempDirProvider *MockTempDirProvider, mockCommand *MockCommand) *ActivateSSHKey {
	return &ActivateSSHKey{
		stepInputParser:     nil,
		envmanEnvRepository: envmanEnvRepository,
		osEnvRepository:     osEnvRepository,
		envValueClearer:     *NewCombinedEnvValueClearer(logger, osEnvRepository, envmanEnvRepository),
		fileWriter:          fileWriter,
		sshKeyAgent: *sshkey.NewAgent(fileWriter, tempDirProvider, logger, func(name string, args ...string) command.Command {
			return mockCommand
		}),
		logger: logger,
	}
}

func createConfigWithDefaults() Config {
	return Config{
		sshRsaPrivateKey:        privateKey,
		sshKeySavePath:          "test-path",
		isRemoveOtherIdentities: false,
		verbose:                 false,
	}
}
