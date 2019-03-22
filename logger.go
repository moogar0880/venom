package venom

import (
    "fmt"
    "log"
    "os"
    "time"
)

const (
    TIME_FORMAT = "2006-01-02T15:04:05.999Z"
)

// LoggingInterface is the interface which any logging mechanism must satisfy
// in order to be used as the logging mechanism for a LoggableConfigStore.
type LoggingInterface interface {
    Print(...interface{})
    // Printf(string, ...interface{})
}

// DefaultLogger is the struct associated with the LoggableConfigStore. It has
// three configurable fields which ultimately compose structured log lines.
type DefaultLogger struct {
    log    LoggingInterface
    Suffix string
    Prefix string
}

// NewDefaultLogger returns a reference to an instance of a DefaultLogger
// struct with defaults for all fields.
func NewDefaultLogger() *DefaultLogger {
    return &DefaultLogger {
        log:    log.New(os.Stdout, "", 0),
        Suffix: "venom",
        Prefix: "INFO",
    }
}

// Write takes a slice of interfaces and prints them as a structured log line
// using the DefaultLogger's log which is a LoggingInterface.
func (lg *DefaultLogger) Write(a ...interface{}) {
    logLine := fmt.Sprintf("%s %s - %s: ", lg.Prefix, time.Now().UTC().Format(TIME_FORMAT), lg.Suffix)
    for _, v := range a {
        logLine += fmt.Sprintf(" %s", v)
    }

    lg.log.Print(logLine)
}

// SetLogger takes an parameter of type LoggingInterface and sets it as the log
// field for a DefaultLogger.
func (lg *DefaultLogger) SetLogger(logger LoggingInterface) {
    lg.log = logger
}

// SetPrefix takes a string and sets the DefaultLogger's Prefix field.
func (lg *DefaultLogger) SetPrefix(s string) {
    lg.Prefix = s
}

// SetSuffix takes a string and sets the DefaultLogger's Suffix field.
func (lg *DefaultLogger) SetSuffix(s string) {
    lg.Suffix = s
}
