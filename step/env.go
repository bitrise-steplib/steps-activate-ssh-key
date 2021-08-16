package step

import (
	"strings"
)

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

type EnvironmentManagerOrchestrator interface {
	Set(key string, value string, setters ...Setter) error
	Get(key string, getter Getter) (string, error)
	GetKey(value string, lister Lister) (string, error)
	List(lister Lister) ([]string, error)
	Unset(key string, unSetters ...UnSetter) error
	UnsetByValue(value string, byValueUnSetters ...ByValueUnSetter) error
}

type DefaultEnvironmentManagerOrchestrator struct{}

func (DefaultEnvironmentManagerOrchestrator) Set(key string, value string, setters ...Setter) error {
	for _, s := range setters {
		if err := s.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (DefaultEnvironmentManagerOrchestrator) Get(key string, getter Getter) (string, error) {
	return getter.Get(key)
}

func (DefaultEnvironmentManagerOrchestrator) GetKey(value string, lister Lister) (string, error) {
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

func (DefaultEnvironmentManagerOrchestrator) List(lister Lister) ([]string, error) {
	return lister.List()
}

func (DefaultEnvironmentManagerOrchestrator) Unset(key string, unSetters ...UnSetter) error {
	for _, s := range unSetters {
		if err := s.Unset(key); err != nil {
			return err
		}
	}
	return nil
}

func (DefaultEnvironmentManagerOrchestrator) UnsetByValue(value string, byValueUnSetters ...ByValueUnSetter) error {
	for _, s := range byValueUnSetters {
		if err := s.UnsetByValue(value); err != nil {
			return err
		}
	}
	return nil
}

type Setter interface {
	Set(key string, value string) error
}

type Getter interface {
	Get(key string) (string, error)
}

type Lister interface {
	List() ([]string, error)
}

type UnSetter interface {
	Unset(value string) error
}

type ByValueUnSetter interface {
	UnsetByValue(value string) error
}
