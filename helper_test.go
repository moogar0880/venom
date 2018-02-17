package venom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
