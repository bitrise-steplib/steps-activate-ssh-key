package command

import (
	"io"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
)

// TODO: Move to `go-utils`

// Opts ...
type Opts struct {
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
}

// Factory ...
type Factory interface {
	Create(name string, args []string, opts *Opts) Command
}

type defaultFactory struct{}

// NewDefaultFactory ...
func NewDefaultFactory() Factory {
	return &defaultFactory{}
}

// Create ...
func (b *defaultFactory) Create(name string, args []string, opts *Opts) Command {
	cmd := exec.Command(name, args...)
	if opts != nil {
		cmd.Stdout = opts.Stdout
		cmd.Stderr = opts.Stderr
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	return defaultCommand{cmd}
}

// Command ...
type Command interface {
	PrintableCommandArgs() string
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnExitCode() (int, error)
	Run() error
}

type defaultCommand struct {
	cmd *exec.Cmd
}

// RunAndReturnTrimmedOutput ...
func (r defaultCommand) RunAndReturnTrimmedOutput() (string, error) {
	return command.RunCmdAndReturnTrimmedOutput(r.cmd)
}

// RunAndReturnExitCode ...
func (r defaultCommand) RunAndReturnExitCode() (int, error) {
	return command.RunCmdAndReturnExitCode(r.cmd)
}

// Run ...
func (r defaultCommand) Run() error {
	return r.cmd.Run()
}

// PrintableCommandArgs ...
func (r defaultCommand) PrintableCommandArgs() string {
	return command.PrintableCommandArgs(false, r.cmd.Args)
}
