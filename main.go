package main

import (
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/activatesshkey"
	"os"
)

func main() {
	if err := activatesshkey.Run(); err != nil {
		log.Errorf("Step run failed: %s", err.Error())
		os.Exit(1)
	}
}
