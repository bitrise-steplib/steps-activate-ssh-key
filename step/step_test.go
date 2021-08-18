package step

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const privateKey = "test-key"
const envKey = "env-key"

func Test_GivenFailingSSHAgent_WhenStepRuns_ThenSSHAgentGetsRestartedAndSSHKeyGetsAdded(t *testing.T) {
	// Given
	logger := createLogger()

	envRepository := new(MockRepository)
	envRepository.On("List").Return(nil)

	fileWriter := new(MockFileWriter)
	fileWriter.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything, mock.Anything).Return(nil).Once()

	step := NewActivateSSHKey(nil, envRepository, fileWriter, sshKeyAgent, logger)

	// When
	_, err := step.Run(config)

	// Then
	assert.NoError(t, err)
	sshKeyAgent.AssertCalled(t, "Start")
	sshKeyAgent.AssertCalled(t, "AddKey", mock.Anything, mock.Anything)
}

func Test_WhenStepRuns_ThenPrivateKeyEnvGetsRemoved(t *testing.T) {
	// Given
	logger := createLogger()

	envRepository := new(MockRepository)
	envRepository.On("List").Return([]string{envKey + "=" + privateKey})
	envRepository.On("Unset", mock.Anything).Return(nil)

	fileWriter := createFileWriter()
	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything, mock.Anything).Return(nil).Once()

	step := NewActivateSSHKey(nil, envRepository, fileWriter, sshKeyAgent, logger)

	// When
	output, err := step.Run(config)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, output.sshAuthSock, "")
	envRepository.AssertCalled(t, "Unset", envKey)
}

func Test_GivenSSHKeyAddFails_WhenStepRuns_ThenItFails(t *testing.T) {
	// Given
	logger := createLogger()
	envRepository := new(MockRepository)
	envRepository.On("List").Return(nil)
	fileWriter := createFileWriter()
	config := createConfigWithDefaults()

	sshKeyAgent := new(MockSSHKeyAgent)
	sshKeyAgent.On("ListKeys").Return(2, errors.New("exit status 2")).Once()
	sshKeyAgent.On("Start").Return("", nil)
	sshKeyAgent.On("AddKey", mock.Anything, mock.Anything).Return(errors.New("mocked error")).Once()

	step := NewActivateSSHKey(nil, envRepository, fileWriter, sshKeyAgent, logger)

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
