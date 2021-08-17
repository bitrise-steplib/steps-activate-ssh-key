package sshkey

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-steplib/steps-activate-ssh-key/command"

	"github.com/bitrise-steplib/steps-activate-ssh-key/filewriter"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/pathutil"
)

// Agent ...
type Agent struct {
	fileWriter      filewriter.FileWriter
	tempDirProvider pathutil.TempDirProvider
	logger          log.Logger
	cmdFactory      command.Factory
}

// NewAgent ...
func NewAgent(fileWriter filewriter.FileWriter, tempDirProvider pathutil.TempDirProvider, logger log.Logger, cmdFactory command.Factory) *Agent {
	return &Agent{fileWriter: fileWriter, tempDirProvider: tempDirProvider, logger: logger, cmdFactory: cmdFactory}
}

// Start ...
func (a Agent) Start() (string, error) {
	cmd := a.cmdFactory.Create("ssh-agent", nil, nil)

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnTrimmedOutput()
}

// Kill ...
func (a Agent) Kill() (int, error) {
	// try to kill the agent
	cmd := a.cmdFactory.Create("ssh-agent", []string{"-k"}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// ListKeys ...
func (a Agent) ListKeys() (int, error) {
	cmd := a.cmdFactory.Create("ssh-add", []string{"-l"}, &command.Opts{
		Stderr: os.Stderr,
	})
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

func createAddSSHKeyScript(sshKeyPth string) string {
	return fmt.Sprintf(`expect <<EOD
spawn ssh-add %s
expect {
	"Enter passphrase for" {
		exit 1
	}
	"Identity added" {
		exit 0
	}
}
send "nopass\n"
EOD
if [ $? -ne 0 ] ; then
exit 1
fi`, sshKeyPth)
}

const addSSHKeyScriptFileName = "tmp_spawn.sh"

// AddKey ...
func (a Agent) AddKey(sshKeyPth, socket string) error {
	pth, err := a.tempDirProvider.CreateTempDir("spawn")
	if err != nil {
		return err
	}

	filePth := filepath.Join(pth, addSSHKeyScriptFileName)
	if err := a.fileWriter.Write(filePth, createAddSSHKeyScript(sshKeyPth), 0770); err != nil {
		return fmt.Errorf("failed to write the SSH key to the provided path, %s", err)
	}

	cmd := a.cmdFactory.Create("bash", []string{"-c", filePth}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    []string{"SSH_AUTH_SOCK=" + socket},
	})

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	exitCode, err := cmd.RunAndReturnExitCode()

	if err != nil {
		a.logger.Debugf("Exit code: %s", err)
	}

	if exitCode != 0 {
		a.logger.Errorf("\nExit code: %d", exitCode)
		return fmt.Errorf("failed to add the SSH key to ssh-agent with an empty passphrase")
	}

	return nil
}

// DeleteKeys ...
func (a Agent) DeleteKeys() error {
	// remove all keys from the current agent
	cmd := a.cmdFactory.Create("ssh-add", []string{"-D"}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.Run()
}
