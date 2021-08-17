package step

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const privateKey = "test-key"
const envKey = "env-key"

func Test_GivenFailingSSHAgent_WhenStepRuns_ThenSSHAgentGetsRestartedAndSSHKeyGetsAdded(t *testing.T) {
	// Given
	logger := log.NewDefaultLogger()

	osEnvRepository := new(MockOsRepository)
	osEnvRepository.On("Set", mock.Anything, mock.Anything).Return(nil)

	envValueClearer := new(MockEnvValueClearer)
	envValueClearer.On("UnsetByValue", mock.Anything).Return(nil)

	fileWriter := new(MockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tempDirProvider := new(MockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return("temp-dir", nil)

	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything).Return(nil).Once()

	step := NewActivateSSHKey(nil, envValueClearer, nil, osEnvRepository, fileWriter, sshKeyAgent, logger)

	// When
	_, err := step.Run(config)

	// Then
	assert.NoError(t, err)
	sshKeyAgent.AssertCalled(t, "Start")
	sshKeyAgent.AssertCalled(t, "AddKey", mock.Anything)
}

func Test_WhenStepRuns_ThenPrivateKeyEnvGetsRemoved(t *testing.T) {
	// Given
	logger := createLogger()
	osEnvRepository := createOsEnvRepositoryWithSSHKey()
	envValueClearer := new(MockEnvValueClearer)
	envValueClearer.On("UnsetByValue", privateKey).Return(nil)
	fileWriter := createFileWriter()
	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything).Return(nil).Once()

	step := NewActivateSSHKey(nil, envValueClearer, nil, osEnvRepository, fileWriter, sshKeyAgent, logger)

	// When
	output, err := step.Run(config)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, output.sshAuthSock, "")
	envValueClearer.AssertCalled(t, "UnsetByValue", privateKey)
}

func Test_GivenSSHKeyAddFails_WhenStepRuns_ThenItFails(t *testing.T) {
	// Given
	logger := createLogger()
	osEnvRepository := createOsEnvRepositoryWithSSHKey()
	envValueClearer := new(MockEnvValueClearer)
	envValueClearer.On("UnsetByValue", mock.Anything).Return(nil)
	fileWriter := createFileWriter()
	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything).Return(errors.New("mocked error")).Once()

	step := NewActivateSSHKey(nil, envValueClearer, nil, osEnvRepository, fileWriter, sshKeyAgent, logger)

	// When
	output, err := step.Run(config)

	// Then
	wantOutput := Result{sshAuthSock: ""}
	wantErr := newStepError(
		"ssh_key_requires_passphrase",
		fmt.Errorf("SSH key requires passphrase: %v", errors.New("mocked error")),
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

func createConfigWithDefaults() Config {
	return Config{
		sshRsaPrivateKey:        privateKey,
		sshKeySavePath:          "test-path",
		isRemoveOtherIdentities: false,
		verbose:                 false,
	}
}
