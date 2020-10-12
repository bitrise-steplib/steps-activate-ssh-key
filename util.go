package main

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
)

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
