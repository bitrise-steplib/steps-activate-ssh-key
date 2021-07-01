package pathutil

import "github.com/bitrise-io/go-utils/pathutil"

// TODO: Move to `go-utils`

// TempDirProvider ...
type TempDirProvider interface {
	CreateTempDir(prefix string) (string, error)
}

type osTempDirProvider struct{}

// NewOsTempDirProvider ...
func NewOsTempDirProvider() TempDirProvider {
	return osTempDirProvider{}
}

// CreateTempDir ...
func (osTempDirProvider) CreateTempDir(prefix string) (string, error) {
	return pathutil.NormalizedOSTempDirPath(prefix)
}
