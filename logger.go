package venom

import (
	"fmt"
	"log"
	"time"
)

const (
	TimeFormat = "2006-01-02T15:04:05.999Z"
	LogName    = "[venom]"
)

// Logger is the interface a user must implement in order to be used by
// the LogableConfigStore.
type Logger interface {
	LogWrite(level ConfigLevel, key string, val interface{})
	LogRead(key string, val interface{}, bl bool)
}

// StoreLogger is Venom's default Logger.
type StoreLogger struct {
	Log *log.Logger
}

// NewStoreLogger returns a Logger suitable for default use case.
func NewStoreLogger(l *log.Logger) Logger {
	return &StoreLogger{
		Log: l,
	}
}

// LogWrite is the default logging behavior of a LoggableConfigStore on an
// action to set a value in the ConfigStore.
func (sl *StoreLogger) LogWrite(level ConfigLevel, key string, val interface{}) {
	logLine := fmt.Sprintf("%s%s: writing level=%v key=%s val=%s",
		time.Now().UTC().Format(TimeFormat),
		LogName,
		level,
		key,
		val,
	)
	sl.Log.Print(logLine)
}

// LogRead is the default logging behavior of a LoggableConfigStore on an
// action to read a value in the ConfigStore.
func (sl *StoreLogger) LogRead(key string, val interface{}, bl bool) {
	logLine := fmt.Sprintf("%s%s: reading key=%s val=%s exist=%v",
		time.Now().UTC().Format(TimeFormat),
		LogName,
		key,
		val,
		bl,
	)
	sl.Log.Print(logLine)
}
