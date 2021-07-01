package step

import (
	"strings"

	"github.com/bitrise-steplib/steps-activate-ssh-key/env"
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
)

// CombinedEnvValueClearer ...
type CombinedEnvValueClearer struct {
	logger              log.Logger
	osEnvRepository     env.OsRepository
	envmanEnvRepository env.EnvmanRepository
}

// NewCombinedEnvValueClearer ...
func NewCombinedEnvValueClearer(logger log.Logger, osEnvRepository env.OsRepository, envmanEnvRepository env.EnvmanRepository) *CombinedEnvValueClearer {
	return &CombinedEnvValueClearer{logger: logger, osEnvRepository: osEnvRepository, envmanEnvRepository: envmanEnvRepository}
}

// UnsetByValue ...
func (o CombinedEnvValueClearer) UnsetByValue(value string) error {
	for _, e := range o.osEnvRepository.List() {
		key, val := splitEnv(e)

		if val == value {
			if err := o.osEnvRepository.Unset(key); err != nil {
				return err
			}

			if err := o.envmanEnvRepository.Unset(key); err != nil {
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
