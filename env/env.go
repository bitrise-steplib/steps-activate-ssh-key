package env

import (
	"os"

	"github.com/bitrise-io/go-steputils/tools"
)

// TODO: Move to `go-utils`

// Repository ...
type Repository interface {
	Set(key, value string) error
	Unset(key string) error
	List() []string
}

// NewRepository ...
func NewRepository(osRepository OsRepository, envmanRepository EnvmanRepository) Repository {
	return defaultRepository{
		osRepository:     osRepository,
		envmanRepository: envmanRepository,
	}
}

type defaultRepository struct {
	osRepository     OsRepository
	envmanRepository EnvmanRepository
}

// Set ...
func (r defaultRepository) Set(key, value string) error {
	if err := r.osRepository.Set(key, value); err != nil {
		return err
	}
	if err := r.envmanRepository.Set(key, value); err != nil {
		return err
	}
	return nil
}

func (r defaultRepository) Unset(key string) error {
	if err := r.osRepository.Unset(key); err != nil {
		return err
	}
	if err := r.envmanRepository.Unset(key); err != nil {
		return err
	}
	return nil
}

func (r defaultRepository) List() []string {
	return r.osRepository.List()
}

// OsRepository ...
type OsRepository interface {
	Unset(key string) error
	Set(key string, value string) error
	List() []string
}

type osRepository struct{}

// NewOsRepository ...
func NewOsRepository() OsRepository {
	return osRepository{}
}

// List ...
func (m osRepository) List() []string {
	return os.Environ()
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
