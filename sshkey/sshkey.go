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
	cmdBuilder      command.Builder
}

// NewAgent ...
func NewAgent(fileWriter filewriter.FileWriter, tempDirProvider pathutil.TempDirProvider, logger log.Logger, cmdBuilder command.Builder) *Agent {
	return &Agent{fileWriter: fileWriter, tempDirProvider: tempDirProvider, logger: logger, cmdBuilder: cmdBuilder}
}

// Start ...
func (a Agent) Start() (string, error) {
	cmd := a.cmdBuilder.Command("ssh-agent", nil)

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnTrimmedOutput()
}

// Kill ...
func (a Agent) Kill() (int, error) {
	// try to kill the agent
	cmd := a.cmdBuilder.Command("ssh-agent", []string{"-k"}, command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// ListKeys ...
func (a Agent) ListKeys() (int, error) {
	cmd := a.cmdBuilder.Command("ssh-add", []string{"-l"}, command.Opts{
		Stdout: os.Stdout,
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
func (a Agent) AddKey(sshKeyPth string) error {
	pth, err := a.tempDirProvider.CreateTempDir("spawn")
	if err != nil {
		return err
	}

	filePth := filepath.Join(pth, addSSHKeyScriptFileName)
	if err := a.fileWriter.Write(filePth, createAddSSHKeyScript(sshKeyPth), 0770); err != nil {
		return fmt.Errorf("failed to write the SSH key to the provided path, %s", err)
	}

	cmd := a.cmdBuilder.Command("bash", []string{"-c", filePth}, command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
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
	cmd := a.cmdBuilder.Command("ssh-add", []string{"-D"}, command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})

	a.logger.Println()
	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.Run()
}
