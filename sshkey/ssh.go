package sshkey

import (
	"fmt"
	"github.com/bitrise-io/go-utils/command"
	"os"
	"path/filepath"
)

type fileWriter interface {
	Write(path string, value string, mode os.FileMode) error
}

type logger interface {
	Printf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Println()
}

type commandRunner interface {
	runAndReturnExitCode(model *command.Model) (int, error)
	runAndReturnTrimmedOutput(model *command.Model) (string, error)
	run(model *command.Model) error
}

type tempDirProvider interface {
	createTempDir(prefix string) (string, error)
}

// Agent ...
type Agent struct {
	fileWriter      fileWriter
	tempDirProvider tempDirProvider
	commandRunner   commandRunner
	logger          logger
}

// NewAgent ...
func NewAgent(fileWriter fileWriter, tempDirProvider tempDirProvider, commandRunner commandRunner, logger logger) *Agent {
	return &Agent{fileWriter: fileWriter, tempDirProvider: tempDirProvider, commandRunner: commandRunner, logger: logger}
}

// Start ...
func (a Agent) Start() (string, error) {
	cmd := command.New("ssh-agent")

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return a.commandRunner.runAndReturnTrimmedOutput(cmd)
}

// Kill ...
func (a Agent) Kill() (int, error) {
	// try to kill the agent
	cmdKill := command.New("ssh-agent", "-k")
	cmdKill.SetStdout(os.Stdout)
	cmdKill.SetStderr(os.Stderr)

	a.logger.Println()
	a.logger.Printf("$ %s", cmdKill.PrintableCommandArgs())

	return a.commandRunner.runAndReturnExitCode(cmdKill)
}

// ListKeys ...
func (a Agent) ListKeys() (int, error) {
	cmd := command.New("ssh-add", "-l")
	cmd.SetStderr(os.Stderr)
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	return a.commandRunner.runAndReturnExitCode(cmd)
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

	pth, err := a.tempDirProvider.createTempDir("spawn")
	if err != nil {
		return err
	}

	filePth := filepath.Join(pth, "tmp_spawn.sh")
	if err := a.fileWriter.Write(filePth, spawnString, 0770); err != nil {
		return fmt.Errorf("failed to write the SSH key to the provided path, %s", err)
	}

	cmd := command.New("bash", "-c", filePth)
	cmd.SetStderr(os.Stderr)
	cmd.SetStdout(os.Stdout)

	a.logger.Println()
	a.logger.Printf("$ %s", cmd.PrintableCommandArgs())

	exitCode, err := a.commandRunner.runAndReturnExitCode(cmd)

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
	cmdRemove := command.New("ssh-add", "-D")
	cmdRemove.SetStdout(os.Stdout)
	cmdRemove.SetStderr(os.Stderr)

	a.logger.Println()
	a.logger.Println()
	a.logger.Printf("$ %s", cmdRemove.PrintableCommandArgs())

	return a.commandRunner.run(cmdRemove)
}
