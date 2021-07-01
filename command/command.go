package command

import (
	"io"

	"github.com/bitrise-io/go-utils/command"
)

// TODO: Move to `go-utils`

// Command ...
type Command interface {
	PrintableCommandArgs() string
	RunAndReturnTrimmedOutput() (string, error)
	SetStdout(stdout io.Writer) *command.Model
	SetStderr(stdout io.Writer) *command.Model
	RunAndReturnExitCode() (int, error)
	Run() error
}

// Factory ...
type Factory func(name string, args ...string) Command

// NewCommand ...
func NewCommand(name string, args ...string) Command {
	return command.New(name, args...)
}
