package envmanager

import (
	"github.com/bitrise-io/go-steputils/tools"
	"os"
)

// TODO: go-steputils

// OsEnvManager ...
type OsEnvManager struct{}

//NewOsEnvManager ...
func NewOsEnvManager() *OsEnvManager {
	return &OsEnvManager{}
}

//List ...
func (m OsEnvManager) List() []string {
	return os.Environ()
}

//Unset ...
func (m OsEnvManager) Unset(key string) error {
	return os.Unsetenv(key)
}

//Set ...
func (m OsEnvManager) Set(key string, value string) error {
	return os.Setenv(key, value)
}

//EnvmanEnvManager ...
type EnvmanEnvManager struct{}

//NewEnvmanEnvManager ...
func NewEnvmanEnvManager() *EnvmanEnvManager {
	return &EnvmanEnvManager{}
}

//Unset ...
func (e EnvmanEnvManager) Unset(key string) error {
	return tools.ExportEnvironmentWithEnvman(key, "")
}

//Set ...
func (e EnvmanEnvManager) Set(key string, value string) error {
	return tools.ExportEnvironmentWithEnvman(key, value)
}
