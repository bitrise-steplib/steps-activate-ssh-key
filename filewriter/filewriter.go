package filewriter

import (
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/fileutil"
)

// TODO: Move to `go-utils`

// FileWriter ...
type FileWriter interface {
	Write(path string, value string, mode os.FileMode) error
}

type osFileWriter struct{}

// NewOsFileWriter ...
func NewOsFileWriter() FileWriter {
	return osFileWriter{}
}

// Write ...
func (osFileWriter) Write(path string, value string, mode os.FileMode) error {
	if err := ensureSavePath(path); err != nil {
		return err
	}

	if err := fileutil.WriteStringToFile(path, value); err != nil {
		return err
	}

	if err := os.Chmod(path, mode); err != nil {
		return err
	}
	return nil
}

func ensureSavePath(savePath string) error {
	dirPath := filepath.Dir(savePath)
	return os.MkdirAll(dirPath, 0700)
}
