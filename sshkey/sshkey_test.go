package sshkey

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	mockcommand "github.com/bitrise-io/go-utils/command/mocks"
	mockfileutil "github.com/bitrise-io/go-utils/fileutil/mocks"
	mocklog "github.com/bitrise-io/go-utils/log/mocks"
	mockpathutil "github.com/bitrise-io/go-utils/pathutil/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_WhenSSHKeyIsAdded_ThenItCallsSSHAddScript(t *testing.T) {
	// Given
	logger := new(mocklog.Logger)
	logger.On("Printf", mock.Anything, mock.Anything).Return()
	logger.On("Println").Return()

	sshKeyPth := "ssh-key-path"
	tmpDir := "temp-dir"
	sshAddScriptPth := filepath.Join(tmpDir, addSSHKeyScriptFileName)

	tempDirProvider := new(mockpathutil.TempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return(tmpDir, nil)

	fileWriter := new(mockfileutil.FileWriter)
	fileWriter.On("Write", sshAddScriptPth, createAddSSHKeyScript(sshKeyPth), mock.Anything).Return(nil)

	cmd := new(mockcommand.Command)
	cmd.On("RunAndReturnExitCode").Return(0, nil)
	cmd.On("PrintableCommandArgs").Return("")

	factory := new(mockcommand.Factory)
	factory.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(cmd)

	agent := NewAgent(fileWriter, tempDirProvider, logger, factory)

	// When
	err := agent.AddKey(sshKeyPth, "socket")

	// Then
	assert.NoError(t, err)
	cmd.AssertCalled(t, "RunAndReturnExitCode")
	factory.AssertCalled(t, "Create", "bash", []string{"-c", sshAddScriptPth}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    []string{"SSH_AUTH_SOCK=socket"},
	})
}
