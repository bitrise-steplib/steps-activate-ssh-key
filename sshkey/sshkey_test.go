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

func TestAgent_AddKey(t *testing.T) {
	logger := log.NewDefaultLogger()

	sshKeyPth := "ssh-key-path"
	tmpDir := "temp-dir"
	sshAddScriptPth := filepath.Join(tmpDir, addSSHKeyScriptFileName)

	tempDirProvider := new(MockTempDirProvider)
	tempDirProvider.On("CreateTempDir", mock.Anything).Return(tmpDir, nil).Once()

	fileWriter := new(MockFileWriter)
	fileWriter.On("Write", sshAddScriptPth, createAddSSHKeyScript(sshKeyPth), mock.Anything).Return(nil).Once()

	cmd := new(MockRunnable)
	cmd.On("RunAndReturnExitCode").Return(0, nil).Once()
	cmd.On("PrintableCommandArgs").Return("").Once()

	cmdBuilder := new(MockCommandBuilder)
	cmdBuilder.On("Command", "bash", []string{"-c", sshAddScriptPth}, command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}).Return(cmd)

	agent := NewAgent(fileWriter, tempDirProvider, logger, cmdBuilder)
	err := agent.AddKey(sshKeyPth)
	assert.NoError(t, err)

	tempDirProvider.AssertExpectations(t)
	fileWriter.AssertExpectations(t)
	cmd.AssertExpectations(t)
}
