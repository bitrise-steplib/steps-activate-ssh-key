package command

import (
	"io"

	"github.com/bitrise-io/go-utils/command"
)

// TODO: move to a util lib

//Command ...
type Command interface {
	PrintableCommandArgs() string
	RunAndReturnTrimmedOutput() (string, error)
	SetStdout(stdout io.Writer) *command.Model
	SetStderr(stdout io.Writer) *command.Model
	RunAndReturnExitCode() (int, error)
	Run() error
}

// CommandFactory ...
type CommandFactory func(name string, args ...string) Command

// NewCommand ...
func NewCommand(name string, args ...string) Command {
	return command.New(name, args...)
}
