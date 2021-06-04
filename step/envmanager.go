package step

import (
	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"os"
	"strings"
)

// TODO: go-steputils
type osEnvManager struct {
	logger log.Logger
}

func newOsEnvManager(logger log.Logger) *osEnvManager {
	return &osEnvManager{logger: logger}
}

func (o osEnvManager) set(key string, value string) error {
	if err := os.Setenv(key, value); err != nil {
		return err
	}
	if err := tools.ExportEnvironmentWithEnvman(key, value); err != nil {
		return err
	}
	return nil
}

func (o osEnvManager) unsetByValue(value string) error {
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