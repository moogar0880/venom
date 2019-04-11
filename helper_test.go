package venom

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogger has a no-op Print() method
type TestLogger struct{}

func (tl *TestLogger) Print(a ...interface{}) {}

// TestLogWrapper has no-op ReadLog WriteLog methods
type TestLogWrapper struct{}

func (tl *TestLogWrapper) LogWrite(level ConfigLevel, key string, val interface{}) {}
func (tl *TestLogWrapper) LogRead(key string, val interface{}, bl bool)            {}

// kv is a test struct containing a (k)ey and a (v)alue
type kv struct {
	k string
	v interface{}
}

// lkv is a test struct containing a config (l)evel, a (k)ey, and a (v)alue
type lkv struct {
	l ConfigLevel
	k string
	v interface{}
}

func assertEqualErrors(t *testing.T, expect, actual error) {
	var msg string
	if actual != nil {
		msg = actual.Error()
	}
	assert.Equal(t, expect, actual, msg)
}

// redirectStdout is explicitly for TestNewLoggable test to pipe the contents
// sent to os.Stdout when implementing using the default logger.
func redirectStdout(test struct {
	tc  string
	f   func() *Venom
	log bool
	kv  kv
}) string {
	// make pipe to stdout and defer resetting
	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, _ := os.Pipe()
	os.Stdout = fakeStdout

	// run test case with fake stdout capture
	l := test.f()
	l.SetDefault(test.kv.k, test.kv.v)

	// close up pipe, return string, exec defer
	fakeStdout.Close()
	newOutBytes, _ := ioutil.ReadAll(r)
	r.Close()
	return string(newOutBytes)
}
