package filewriter

import (
	"github.com/bitrise-io/go-utils/fileutil"
	"os"
	"path/filepath"
)

// OsFileWriter ...
type OsFileWriter struct{}

// NewOsFileWriter ...
func NewOsFileWriter() *OsFileWriter {
	return &OsFileWriter{}
}

// Write ...
func (OsFileWriter) Write(path string, value string, mode os.FileMode) error {
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