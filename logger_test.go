package venom

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggableWith(t *testing.T) {
	testIO := []struct {
		tc  string
		lgr *log.Logger
	}{
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
			lw := NewStoreLogger(test.lgr)
			l := NewLoggableWith(lw)
			assert.IsType(t, l, &Venom{})
		})
	}
}

func TestNewLoggable(t *testing.T) {
	testIO := []struct {
		tc  string
		f   func() *Venom
		log bool
		kv  kv
	}{
		{
			tc:  "should be receive a default to stdout if Loggable",
			f:   NewLoggable,
			log: true,
			kv:  kv{k: "foo", v: "bar"},
		},
		{
			tc:  "should not log if new default store",
			f:   New,
			log: false,
			kv:  kv{k: "foo", v: "bar"},
		},
		{
			tc:  "should not log if new safe store",
			f:   NewSafe,
			log: false,
			kv:  kv{k: "foo", v: "bar"},
		},
		{
			tc:  "should log to stdout with a value type int",
			f:   NewLoggable,
			log: true,
			kv:  kv{k: "foo", v: 100},
		},
		{
			tc:  "should log to stdout with a value type float",
			f:   NewLoggable,
			log: true,
			kv:  kv{k: "foo", v: 10.0},
		},
		{
			tc:  "should log to stdout with a value type boolean true",
			f:   NewLoggable,
			log: true,
			kv:  kv{k: "foo", v: true},
		},
		{
			tc:  "should log to stdout with a value type boolean false",
			f:   NewLoggable,
			log: true,
			kv:  kv{k: "foo", v: false},
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			out := redirectStdout(test)
			if test.log {
				assert.Contains(t, out, fmt.Sprintf("writing level=0 key=%s val=%s", test.kv.k, test.kv.v))
			} else {
				assert.Empty(t, out)
			}
		})
	}
}
