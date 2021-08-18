package command

import (
	"io"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-steplib/steps-activate-ssh-key/env"
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

type defaultFactory struct {
	envLister env.Lister
}

// NewDefaultFactory ...
func NewDefaultFactory(envLister env.Lister) Factory {
	return &defaultFactory{
		envLister: envLister,
	}
}

// Create ...
func (b *defaultFactory) Create(name string, args []string, opts *Opts) Command {
	cmd := exec.Command(name, args...)
	if opts != nil {
		cmd.Stdout = opts.Stdout
		cmd.Stderr = opts.Stderr

		// If Env is nil, the new process uses the current process's
		// environment.
		// If we pass env vars we want to append them to the
		// current process's environment.
		cmd.Env = append(b.envLister.List(), opts.Env...)
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
