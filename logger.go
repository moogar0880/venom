package venom

import (
	"fmt"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02T15:04:05.999Z"
	LOG_NAME    = "[venom]"
)

// Logger is the interface which any logging mechanism must satisfy in order to
// be used as the Logger mechanism wrapped by LogWrapper.
type Logger interface {
	Print(...interface{})
}

// LogWrapper is the wrapping interface by which a user must adhere in order to
// be used by the LogableConfigStore.
type LogWrapper interface {
	LogWrite(level ConfigLevel, key string, val interface{})
	LogRead(key string, val interface{}, bl bool)
}

// Entry is Venom's default LogWrapper and borrows the wording from logrus.
type Entry struct {
	Log Logger
}

// NewEntry returns a LogWrapper suitable for default use case.
func NewEntry(l Logger) LogWrapper {
	return &Entry{
		Log: l,
	}
}

// LogWrite is the default logging behavior of a LoggableConfigStore on an
// action to set a value in the ConfigStore.
func (e *Entry) LogWrite(level ConfigLevel, key string, val interface{}) {
	logLine := fmt.Sprintf("%s%s: writing level=%v key=%s val=%s", time.Now().UTC().Format(TIME_FORMAT), LOG_NAME, level, key, val)
	e.Log.Print(logLine)
}

// LogWrite is the default logging behavior of a LoggableConfigStore on an
// action to read a value in the ConfigStore.
func (e *Entry) LogRead(key string, val interface{}, bl bool) {
	logLine := fmt.Sprintf("%s%s: reading key=%s val=%s exist=%v", time.Now().UTC().Format(TIME_FORMAT), LOG_NAME, key, val, bl)
	e.Log.Print(logLine)
}
