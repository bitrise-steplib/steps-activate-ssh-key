package sshkey

import "github.com/bitrise-io/go-utils/command"

// OsCommandRunner ...
type OsCommandRunner struct{}

// NewOsCommandRunner ...
func NewOsCommandRunner() *OsCommandRunner {
	return &OsCommandRunner{}
}

// RunAndReturnExitCode ...
func (OsCommandRunner) RunAndReturnExitCode(model *command.Model) (int, error) {
	return model.RunAndReturnExitCode()
}

// RunAndReturnTrimmedOutput ...
func (OsCommandRunner) RunAndReturnTrimmedOutput(model *command.Model) (string, error) {
	return model.RunAndReturnTrimmedOutput()
}

// Run ...
func (OsCommandRunner) Run(model *command.Model) error {
	return model.Run()
}