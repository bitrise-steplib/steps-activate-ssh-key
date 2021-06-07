package step

import (
	"strings"
)

type envManager interface {
	Unset(key string) error
	Set(key string, value string) error
}

type extendedEnvManager interface {
	Unset(key string) error
	Set(key string, value string) error
	List() []string
}

//CombinedEnvValueClearer ...
type CombinedEnvValueClearer struct {
	logger           logger
	osEnvManager     extendedEnvManager
	envmanEnvManager envManager
}

//NewCombinedEnvValueClearer ...
func NewCombinedEnvValueClearer(logger logger, osEnvManager extendedEnvManager, envmanEnvManager envManager) *CombinedEnvValueClearer {
	return &CombinedEnvValueClearer{logger: logger, osEnvManager: osEnvManager, envmanEnvManager: envmanEnvManager}
}

//UnsetByValue ...
func (o CombinedEnvValueClearer) UnsetByValue(value string) error {
	for _, env := range o.osEnvManager.List() {
		key, val := splitEnv(env)

		if val == value {
			if err := o.osEnvManager.Unset(key); err != nil {
				return err
			}

			if err := o.envmanEnvManager.Unset(key); err != nil {
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
