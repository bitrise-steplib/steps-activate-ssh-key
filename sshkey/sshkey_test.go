package sshkey

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-steplib/steps-activate-ssh-key/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_WhenSSHKeyIsAdded_ThenItCallsSSHAddScript(t *testing.T) {
	// Given
	logger := log.NewDefaultLogger()

	sshKeyPth := "ssh-key-path"
	tmpDir := "temp-dir"
	sshAddScriptPth := filepath.Join(tmpDir, addSSHKeyScriptFileName)

	tempDirProvider := new(MockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return(tmpDir, nil)

	fileWriter := new(MockFileWriter)
	fileWriter.On("Write", sshAddScriptPth, createAddSSHKeyScript(sshKeyPth), mock.Anything).Return(nil)

	cmd := new(MockCommand)
	cmd.On("RunAndReturnExitCode").Return(0, nil)
	cmd.On("PrintableCommandArgs").Return("")

	factory := new(MockFactory)
	factory.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(cmd)

	agent := NewAgent(fileWriter, tempDirProvider, logger, factory)

	// When
	err := agent.AddKey(sshKeyPth)

	// Then
	assert.NoError(t, err)
	cmd.AssertCalled(t, "RunAndReturnExitCode")
	factory.AssertCalled(t, "Create", "bash", []string{"-c", sshAddScriptPth}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}
