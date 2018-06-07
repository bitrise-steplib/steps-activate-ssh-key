package main

import (
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/tools"
)

func unsetEnvsBy(value string) error {
	for _, env := range os.Environ() {
		key, val := splitEnv(env)

		if val == value {
			if err := os.Unsetenv(key); err != nil {
				return err
			}

			if err := tools.ExportEnvironmentWithEnvman(key, ""); err != nil {
				return err
			}
			log.Debugf("%s has been unset", key)
		}
	}
	return nil
}

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
