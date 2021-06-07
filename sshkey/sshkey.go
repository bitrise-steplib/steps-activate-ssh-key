package sshkey

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"io"
	"os"
	"path/filepath"
)

//Command ...
type Command interface {
	PrintableCommandArgs() string
	RunAndReturnTrimmedOutput() (string, error)
	SetStdout(stdout io.Writer) *command.Model
	SetStderr(stdout io.Writer) *command.Model
	RunAndReturnExitCode() (int, error)
	Run() error
}

type fileWriter interface {
	Write(path string, value string, mode os.FileMode) error
}

type logger interface {
	Printf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Println()
}

type tempDirProvider interface {
	CreateTempDir(prefix string) (string, error)
}

// Agent ...
type Agent struct {
	fileWriter      fileWriter
	tempDirProvider tempDirProvider
	logger          logger
	commandFactory  func(name string, args ...string) *Command
}

// NewAgent ...
func NewAgent(fileWriter fileWriter, tempDirProvider tempDirProvider, logger logger, commandFactory func(name string, args ...string) *Command) *Agent {
	return &Agent{fileWriter: fileWriter, tempDirProvider: tempDirProvider, logger: logger, commandFactory: commandFactory}
}

// Start ...
func (a Agent) Start() (string, error) {
	cmd := *a.commandFactory("ssh-agent")

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnTrimmedOutput()
}

// Kill ...
func (a Agent) Kill() (int, error) {
	// try to kill the agent
	cmd := *a.commandFactory("ssh-agent", "-k")
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// ListKeys ...
func (a Agent) ListKeys() (int, error) {
	cmd := *a.commandFactory("ssh-add", "-l")
	cmd.SetStderr(os.Stderr)
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.RunAndReturnExitCode()
}

// AddKey ...
func (a Agent) AddKey(path string) error {
	spawnString := `expect <<EOD
spawn ssh-add ` + path + `
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
fi`

	pth, err := a.tempDirProvider.CreateTempDir("spawn")
	if err != nil {
		return err
	}

	filePth := filepath.Join(pth, "tmp_spawn.sh")
	if err := a.fileWriter.Write(filePth, spawnString, 0770); err != nil {
		return fmt.Errorf("failed to write the SSH key to the provided path, %s", err)
	}

	cmd := *a.commandFactory("bash", "-c", filePth)
	cmd.SetStderr(os.Stderr)
	cmd.SetStdout(os.Stdout)

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
	cmd := *a.commandFactory("ssh-add", "-D")
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	a.logger.Println()
	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return cmd.Run()
}
