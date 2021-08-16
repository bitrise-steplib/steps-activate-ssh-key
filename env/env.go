package env

import (
	"os"

	"github.com/bitrise-io/go-steputils/tools"
)

// TODO: Move to `go-utils`

// OsRepository ...
type OsRepository interface {
	Unset(key string) error
	Set(key string, value string) error
	List() ([]string, error)
}

type osRepository struct{}

// NewOsRepository ...
func NewOsRepository() OsRepository {
	return osRepository{}
}

// List ...
func (m osRepository) List() ([]string, error) {
	return os.Environ(), nil
}

// Unset ...
func (m osRepository) Unset(key string) error {
	return os.Unsetenv(key)
}

// Set ...
func (m osRepository) Set(key string, value string) error {
	return os.Setenv(key, value)
}

// TODO: Move to `go-steputils`

// EnvmanRepository ...
type EnvmanRepository interface {
	Unset(key string) error
	Set(key string, value string) error
}

type envmanRepository struct{}

// NewEnvmanRepository ...
func NewEnvmanRepository() EnvmanRepository {
	return envmanRepository{}
}

// Unset ...
func (e envmanRepository) Unset(key string) error {
	return tools.ExportEnvironmentWithEnvman(key, "")
}

// Set ...
func (e envmanRepository) Set(key string, value string) error {
	return tools.ExportEnvironmentWithEnvman(key, value)
}
