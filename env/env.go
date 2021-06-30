package env

import (
	"os"

	"github.com/bitrise-io/go-steputils/tools"
)

// TODO: go-steputils

type OsRepository interface {
	Unset(key string) error
	Set(key string, value string) error
	List() []string
}

type osRepository struct{}

//NewOsEnvManager ...
func NewOsEnvManager() OsRepository {
	return &osRepository{}
}

//List ...
func (m osRepository) List() []string {
	return os.Environ()
}

//Unset ...
func (m osRepository) Unset(key string) error {
	return os.Unsetenv(key)
}

//Set ...
func (m osRepository) Set(key string, value string) error {
	return os.Setenv(key, value)
}

type EnvmanRepository interface {
	Unset(key string) error
	Set(key string, value string) error
}

type envmanRepository struct{}

//NewEnvmanEnvManager ...
func NewEnvmanEnvManager() EnvmanRepository {
	return &envmanRepository{}
}

//Unset ...
func (e envmanRepository) Unset(key string) error {
	return tools.ExportEnvironmentWithEnvman(key, "")
}

//Set ...
func (e envmanRepository) Set(key string, value string) error {
	return tools.ExportEnvironmentWithEnvman(key, value)
}
