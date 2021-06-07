package step

import (
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"os"
	"strings"
)

// TODO: go-steputils

// OsEnvManager ...
type OsEnvManager struct {
	logger log.Logger
}

// NewOsEnvManager ...
func NewOsEnvManager(logger log.Logger) *OsEnvManager {
	return &OsEnvManager{logger: logger}
}

//Set ...
func (o OsEnvManager) Set(key string, value string) error {
	if err := os.Setenv(key, value); err != nil {
		return err
	}
	if err := tools.ExportEnvironmentWithEnvman(key, value); err != nil {
		return err
	}
	return nil
}

//UnsetByValue ...
func (o OsEnvManager) UnsetByValue(value string) error {
	for _, env := range os.Environ() {
		key, val := splitEnv(env)

		if val == value {
			if err := os.Unsetenv(key); err != nil {
				return err
			}

			if err := tools.ExportEnvironmentWithEnvman(key, ""); err != nil {
				return err
			}
			o.logger.Debugf("%s has been unset", key)
		}
	}
	return nil
}

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}