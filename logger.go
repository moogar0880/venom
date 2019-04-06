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

// formatLogLine produces a default log line structure for the Logging interface
// print function to use.
func formatLogLine(a ...interface{}) string {
	logLine := fmt.Sprintf("%s%s:", time.Now().UTC().Format(TIME_FORMAT), LOG_NAME)
	for _, v := range a {
		logLine += fmt.Sprintf(" %v", v)
	}

	return logLine
}
