package main

import (
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

func unsetSSHEnvs(sshRsaPrivateKey string) error {
	for _, env := range os.Environ() {
		key, value := splitEnv(env)

		if value == sshRsaPrivateKey {
			log.Debugf("%s has been unset", key)
			if err := os.Unsetenv(key); err != nil {
				return err
			}
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
