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

	cmd := new(MockCommand)
	cmd.On("RunAndReturnExitCode").Return(0, nil).Once()
	cmd.On("SetStdout", os.Stdout).Return(nil).Once()
	cmd.On("SetStderr", os.Stderr).Return(nil).Once()
	cmd.On("PrintableCommandArgs").Return("").Once()
	cmdFactory := func(name string, args ...string) command.Command {
		if len(args) == 2 && name == "bash" && args[0] == "-c" && args[1] == sshAddScriptPth {
			return cmd
		}
		t.Fatalf("Unknown command")
		return nil
	}

	agent := NewAgent(fileWriter, tempDirProvider, logger, cmdFactory)
	err := agent.AddKey(sshKeyPth)
	assert.NoError(t, err)

	tempDirProvider.AssertExpectations(t)
	fileWriter.AssertExpectations(t)
	cmd.AssertExpectations(t)
}
