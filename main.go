package main

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/activatesshkey"
)

func main() {
	var cfg activatesshkey.Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()

	log.SetEnableDebugLog(cfg.Verbose)

	if err := activatesshkey.Execute(activatesshkey.Config{
		SSHRsaPrivateKey:        cfg.SSHRsaPrivateKey,
		SSHKeySavePath:          cfg.SSHKeySavePath,
		IsRemoveOtherIdentities: cfg.IsRemoveOtherIdentities,
	}); err != nil {
		failf(err.Error())
	}
}
