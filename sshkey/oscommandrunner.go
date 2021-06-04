package sshkey

import "github.com/bitrise-io/go-utils/command"

// OsCommandRunner ...
type OsCommandRunner struct{}

// NewOsCommandRunner ...
func NewOsCommandRunner() *OsCommandRunner {
	return &OsCommandRunner{}
}

func (OsCommandRunner) runAndReturnExitCode(model *command.Model) (int, error) {
	return model.RunAndReturnExitCode()
}

func (OsCommandRunner) runAndReturnTrimmedOutput(model *command.Model) (string, error) {
	return model.RunAndReturnTrimmedOutput()
}

func (OsCommandRunner) run(model *command.Model) error {
	return model.Run()
}