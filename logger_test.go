package venom

import (
	"fmt"
	"os"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogableWith(t *testing.T) {
	testIO := []struct {
		tc   string
		lgr  LoggingInterface
	}{
		{
			tc:  "should be able to set a no-op logging interface",
			lgr: &TestLogger{},
		},
		{
			tc:  "should be able set to an io.Writer (stderr)",
			lgr: log.New(os.Stderr, "", 0),
		},
		{
			tc:  "should be able set to an io.Writer (stdout)",
			lgr: log.New(os.Stdout, "", 0),
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			l := NewLogableWith(test.lgr)
			assert.IsType(t, l, &Venom{})
		})
	}
}

func TestNewLogable(t *testing.T) {
	testIO := []struct {
		tc  string
		f   func() *Venom
		log bool
		kv  kv
	}{
		{
			tc: "should be receive a default to stdout if logable",
			f:  NewLogable,
			log: true,
			kv: kv{k: "foo", v: "bar"},
		},
		{
			tc: "should not log if new default store",
			f:  New,
			log: false,
			kv: kv{k: "foo", v: "bar"},
		},
		{
			tc: "should not log if new safe store",
			f: NewSafe,
			log: false,
			kv: kv{k: "foo", v: "bar"},
		},
		{
			tc: "should log to stdout with a value type int",
			f:  NewLogable,
			log: true,
			kv: kv{k: "foo", v: 100},
		},
		{
			tc: "should log to stdout with a value type float",
			f:  NewLogable,
			log: true,
			kv: kv{k: "foo", v: 10.0},
		},
		{
			tc: "should log to stdout with a value type boolean true",
			f:  NewLogable,
			log: true,
			kv: kv{k: "foo", v: true},
		},
		{
			tc: "should log to stdout with a value type boolean false",
			f:  NewLogable,
			log: true,
			kv: kv{k: "foo", v: false},
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			out := redirectStdout(test)
			if test.log {
				assert.Contains(t, out, fmt.Sprintf("[venom]:  [%s %s]", test.kv.k, test.kv.v))
			} else {
				assert.Empty(t, out)
			}
		})
	}
}
