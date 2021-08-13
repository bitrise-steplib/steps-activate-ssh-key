// Code generated by mockery v2.8.0. DO NOT EDIT.

package step

import (
	"io"
	"io/fs"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/mock"
)

// MockEnvValueClearer is an autogenerated mock type for the EnvValueClearer type
type MockEnvValueClearer struct {
	mock.Mock
}

// UnsetByValue provides a mock function with given fields: value
func (_m *MockEnvValueClearer) UnsetByValue(value string) error {
	ret := _m.Called(value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSSHKeyAgent is an autogenerated mock type for the SSHKeyAgent type
type MockSSHKeyAgent struct {
	mock.Mock
}

// AddKey provides a mock function with given fields: sshKeyPth
func (_m *MockSSHKeyAgent) AddKey(sshKeyPth string) error {
	ret := _m.Called(sshKeyPth)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(sshKeyPth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteKeys provides a mock function with given fields:
func (_m *MockSSHKeyAgent) DeleteKeys() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Kill provides a mock function with given fields:
func (_m *MockSSHKeyAgent) Kill() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListKeys provides a mock function with given fields:
func (_m *MockSSHKeyAgent) ListKeys() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields:
func (_m *MockSSHKeyAgent) Start() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOsRepository is an autogenerated mock type for the OsRepository type
type MockOsRepository struct {
	mock.Mock
}

// List provides a mock function with given fields:
func (_m *MockOsRepository) List() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Set provides a mock function with given fields: key, value
func (_m *MockOsRepository) Set(key string, value string) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unset provides a mock function with given fields: key
func (_m *MockOsRepository) Unset(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEnvmanRepository is an autogenerated mock type for the EnvmanRepository type
type MockEnvmanRepository struct {
	mock.Mock
}

// Set provides a mock function with given fields: key, value
func (_m *MockEnvmanRepository) Set(key string, value string) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unset provides a mock function with given fields: key
func (_m *MockEnvmanRepository) Unset(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFileWriter is an autogenerated mock type for the FileWriter type
type MockFileWriter struct {
	mock.Mock
}

// Write provides a mock function with given fields: path, value, mode
func (_m *MockFileWriter) Write(path string, value string, mode fs.FileMode) error {
	ret := _m.Called(path, value, mode)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, fs.FileMode) error); ok {
		r0 = rf(path, value, mode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

// Debugf provides a mock function with given fields: format, v
func (_m *MockLogger) Debugf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Donef provides a mock function with given fields: format, v
func (_m *MockLogger) Donef(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: format, v
func (_m *MockLogger) Errorf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Printf provides a mock function with given fields: format, v
func (_m *MockLogger) Printf(format string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// Println provides a mock function with given fields:
func (_m *MockLogger) Println() {
	_m.Called()
}

// MockTempDirProvider is an autogenerated mock type for the TempDirProvider type
type MockTempDirProvider struct {
	mock.Mock
}

// CreateTempDir provides a mock function with given fields: prefix
func (_m *MockTempDirProvider) CreateTempDir(prefix string) (string, error) {
	ret := _m.Called(prefix)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(prefix)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(prefix)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCommand is an autogenerated mock type for the Command type
type MockCommand struct {
	mock.Mock
}

// PrintableCommandArgs provides a mock function with given fields:
func (_m *MockCommand) PrintableCommandArgs() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Run provides a mock function with given fields:
func (_m *MockCommand) Run() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunAndReturnExitCode provides a mock function with given fields:
func (_m *MockCommand) RunAndReturnExitCode() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunAndReturnTrimmedOutput provides a mock function with given fields:
func (_m *MockCommand) RunAndReturnTrimmedOutput() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetStderr provides a mock function with given fields: stdout
func (_m *MockCommand) SetStderr(stdout io.Writer) *command.Model {
	ret := _m.Called(stdout)

	var r0 *command.Model
	if rf, ok := ret.Get(0).(func(io.Writer) *command.Model); ok {
		r0 = rf(stdout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*command.Model)
		}
	}

	return r0
}

// SetStdout provides a mock function with given fields: stdout
func (_m *MockCommand) SetStdout(stdout io.Writer) *command.Model {
	ret := _m.Called(stdout)

	var r0 *command.Model
	if rf, ok := ret.Get(0).(func(io.Writer) *command.Model); ok {
		r0 = rf(stdout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*command.Model)
		}
	}

	return r0
}
