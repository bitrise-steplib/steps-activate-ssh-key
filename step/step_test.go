package step

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const privateKey = "test-key"
const envKey = "env-key"

func TestStepRun_IfAgentRestarted_SSHKeyAdded(t *testing.T) {
	// Given Simple Activate SSH Key step
	config := createConfigWithDefaults()

	logger := log.NewDefaultLogger()

	osEnvRepository := new(MockOsRepository)
	osEnvRepository.On("Set", mock.Anything, mock.Anything).Return(nil)

	envValueClearer := new(MockEnvValueClearer)
	envValueClearer.On("UnsetByValue", mock.Anything).Return(nil)

	fileWriter := new(MockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tempDirProvider := new(MockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return("temp-dir", nil)

	// When SSH Agent list keys fails and agent gets restarted
	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything).Return(nil).Once()

	step := NewActivateSSHKey(nil, envValueClearer, nil, osEnvRepository, fileWriter, sshKeyAgent, logger)

	_, err := step.Run(config)

	// Then SSH Key gets activated
	assert.NoError(t, err)
	sshKeyAgent.AssertExpectations(t)
}

func Test_SSHPrivateKeyRemoved(t *testing.T) {
	osEnvRepository := createOsEnvRepositoryWithSSHKey()
	osEnvManager := createEnvmanEnvRepository()
	fileWriter := createFileWriter()
	logger := createLogger()
	tempDirProvider := createTempProvider()
	commandFactory := func(name string, args ...string) command.Command { return createCommand() }
	activateSSHKey := createActivateSSHKey(osEnvRepository, osEnvManager, fileWriter, logger, tempDirProvider, commandFactory)
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
	commandFactory := createCommandFactoryWithFailingSSHAdd()

	activateSSHKey := createActivateSSHKey(osEnvRepository, osEnvManager, fileWriter, logger, tempDirProvider, commandFactory)
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

func createCommandFactoryWithFailingSSHAdd() command.Factory {
	mockCommand := createCommand()
	failingMockCommand := new(MockCommand)
	failingMockCommand.On("RunAndReturnExitCode").Return(1, errors.New("mocked error"))
	failingMockCommand.On("SetStdout", mock.Anything).Return(nil)
	failingMockCommand.On("SetStderr", mock.Anything).Return(nil)
	failingMockCommand.On("PrintableCommandArgs").Return("")
	return func(name string, args ...string) command.Command {
		if name == "bash" && args[0] == "-c" {
			return failingMockCommand
		}
		return mockCommand
	}
}

func createActivateSSHKey(osEnvRepository *MockOsRepository, envmanEnvRepository *MockEnvmanRepository, fileWriter *MockFileWriter, logger *MockLogger, tempDirProvider *MockTempDirProvider, commandFactory command.Factory) *ActivateSSHKey {
	return &ActivateSSHKey{
		stepInputParser:     nil,
		envmanEnvRepository: envmanEnvRepository,
		osEnvRepository:     osEnvRepository,
		envValueClearer:     *NewCombinedEnvValueClearer(logger, osEnvRepository, envmanEnvRepository),
		fileWriter:          fileWriter,
		sshKeyAgent:         *sshkey.NewAgent(fileWriter, tempDirProvider, logger, commandFactory),
		logger:              logger,
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
