package step

import (
	"github.com/bitrise-steplib/steps-activate-ssh-key/log"
	"github.com/bitrise-steplib/steps-activate-ssh-key/sshkey"
	"reflect"
	"testing"
)

func Test_activateSSHKey_run(t *testing.T) {
	type fields struct {
		stepInputParse stepInputParser
		envManager     envManager
		fileWriter     fileWriter
		agent          sshkey.Agent
		logger         log.Logger
	}
	type args struct {
		cfg config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := activateSSHKey{
				stepInputParse: tt.fields.stepInputParse,
				envManager:     tt.fields.envManager,
				fileWriter:     tt.fields.fileWriter,
				agent:          tt.fields.agent,
				logger:         tt.fields.logger,
			}
			got, err := a.run(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
