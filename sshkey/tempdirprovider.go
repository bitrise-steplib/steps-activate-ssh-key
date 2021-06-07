package sshkey

import "github.com/bitrise-io/go-utils/pathutil"

// TODO: pathutils / fs abstraction to go-utils

// OsTempDirProvider ...
type OsTempDirProvider struct{}

// NewOsTempDirProvider ...
func NewOsTempDirProvider() *OsTempDirProvider {
	return &OsTempDirProvider{}
}

// CreateTempDir ...
func (OsTempDirProvider) CreateTempDir(prefix string) (string, error) {
	return pathutil.NormalizedOSTempDirPath(prefix)
}
