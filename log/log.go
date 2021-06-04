package log

import (
	"github.com/bitrise-io/go-utils/log"
	log2 "log"
)

// Logger ...
type Logger struct {
}

// NewLogger ...
func NewLogger() *Logger {
	return &Logger{}
}

// Printf ...
func (l Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Donef ...
func (l Logger) Donef(format string, v ...interface{}) {
	log.Donef(format, v...)
}

// Debugf ...
func (l Logger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Errorf ...
func (l Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Println ...
func (l Logger) Println() {
	log2.Print("\n")
}
