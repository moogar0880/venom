package venom

import (
	"fmt"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02T15:04:05.999Z"
	LOG_NAME = "[venom]"
)

// Logging is the interface which any logging mechanism must satisfy
// in order to be used as the logging mechanism for a LoggableConfigStore.
type Logging interface {
	Print(...interface{})
}

// DefaultLogger is the struct associated with the LoggableConfigStore. It has
// three configurable fields which ultimately compose structured log lines.
type DefaultLogger struct {
	log    Logging
	suffix string
}

// NewDefaultLogger returns a reference to an instance of a DefaultLogger
// struct with defaults for all fields.
func NewLogger(l Logging) DefaultLogger {
	return DefaultLogger {
		log:    l,
		suffix: LOG_NAME,
	}
}

// Write takes a slice of interfaces and prints them as a structured log line
// using the DefaultLogger's log which is a Logging.
func (lg *DefaultLogger) Write(a ...interface{}) {
	logLine := fmt.Sprintf("%s%s:", time.Now().UTC().Format(TIME_FORMAT), lg.suffix)
	for _, v := range a {
		logLine += fmt.Sprintf(" %v", v)
	}

	lg.log.Print(logLine)
}
