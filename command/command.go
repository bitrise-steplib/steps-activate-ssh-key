package command

import (
	"io"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
)

// TODO: Move to `go-utils`

// Opts ...
type Opts struct {
	Stdout io.Writer
	Stderr io.Writer
}

// Builder ...
type Builder interface {
	Command(name string, args []string, opts ...Opts) Runnable
}

// DefaultBuilder ...
type DefaultBuilder struct{}

// Command ...
func (b *DefaultBuilder) Command(name string, args []string, opts ...Opts) Runnable {
	cmd := exec.Command(name, args...)
	if len(opts) > 0 {
		opt := opts[0]
		cmd.Stdout = opt.Stdout
		cmd.Stderr = opt.Stderr
	}
	return Command{cmd}
}

// Runnable ...
type Runnable interface {
	PrintableCommandArgs() string
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnExitCode() (int, error)
	Run() error
}

// Command ...
type Command struct {
	cmd *exec.Cmd
}

// RunAndReturnTrimmedOutput ...
func (r Command) RunAndReturnTrimmedOutput() (string, error) {
	return command.RunCmdAndReturnTrimmedOutput(r.cmd)
}

// RunAndReturnExitCode ...
func (r Command) RunAndReturnExitCode() (int, error) {
	return command.RunCmdAndReturnExitCode(r.cmd)
}

// Run ...
func (r Command) Run() error {
	return r.cmd.Run()
}

// PrintableCommandArgs ...
func (r Command) PrintableCommandArgs() string {
	return command.PrintableCommandArgs(false, r.cmd.Args)
}
