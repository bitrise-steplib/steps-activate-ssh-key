package step

import (
	"strings"
)

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

// EnvironmentRepositoryOrchestrator ...
type EnvironmentRepositoryOrchestrator interface {
	Set(key string, value string, setters ...Setter) error
	Get(key string, getter Getter) (string, error)
	GetKey(value string, lister Lister) (string, error)
	List(lister Lister) ([]string, error)
	Unset(key string, unSetters ...UnSetter) error
}

// DefaultEnvironmentRepositoryOrchestrator ...
type DefaultEnvironmentRepositoryOrchestrator struct{}

// Set ...
func (DefaultEnvironmentRepositoryOrchestrator) Set(key string, value string, setters ...Setter) error {
	for _, s := range setters {
		if err := s.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Get ...
func (DefaultEnvironmentRepositoryOrchestrator) Get(key string, getter Getter) (string, error) {
	return getter.Get(key)
}

// GetKey ...
func (DefaultEnvironmentRepositoryOrchestrator) GetKey(value string, lister Lister) (string, error) {
	environment, err := lister.List()
	if err != nil {
		return "", err
	}

	for _, e := range environment {
		key, val := splitEnv(e)
		if val == value {
			return key, nil
		}
	}
	return "", nil
}

// List ...
func (DefaultEnvironmentRepositoryOrchestrator) List(lister Lister) ([]string, error) {
	return lister.List()
}

// Unset ...
func (DefaultEnvironmentRepositoryOrchestrator) Unset(key string, unSetters ...UnSetter) error {
	for _, s := range unSetters {
		if err := s.Unset(key); err != nil {
			return err
		}
	}
	return nil
}

// Setter ...
type Setter interface {
	Set(key string, value string) error
}

// Getter ...
type Getter interface {
	Get(key string) (string, error)
}

// Lister ...
type Lister interface {
	List() ([]string, error)
}

// UnSetter ...
type UnSetter interface {
	Unset(value string) error
}
