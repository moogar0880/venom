package venom

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagsetResolver(t *testing.T) {
	testIO := []struct {
		tc     string
		flags  *flag.FlagSet
		args   []string
		keys   []string
		expect interface{}
		ok     bool
	}{
		{
			tc:     "should resolve nothing from os.Args",
			flags:  nil,
			args:   nil,
			keys:   []string{},
			expect: nil,
			ok:     false,
		},
		{
			tc:     "should resolve nothing from an empty FlagSet",
			flags:  flag.NewFlagSet("test", flag.ContinueOnError),
			args:   []string{},
			keys:   []string{},
			expect: nil,
			ok:     false,
		},
		{
			tc: "should resolve boolean flag from FlagSet",
			flags: func() *flag.FlagSet {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				fs.Bool("verbose", false, "enable verbose")
				return fs
			}(),
			args:   []string{"-verbose"},
			keys:   []string{"verbose"},
			expect: "true",
			ok:     true,
		},
		{
			tc: "should resolve hyphenated config from FlagSet",
			flags: func() *flag.FlagSet {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				fs.String("log-level", "WARNING", "set log level")
				return fs
			}(),
			args:   []string{"-log-level=INFO"},
			keys:   []string{"log", "level"},
			expect: "INFO",
			ok:     true,
		},
		{
			tc: "should resolve expected value from FlagSet",
			flags: func() *flag.FlagSet {
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				fs.Bool("verbose", false, "enable verbose")
				fs.String("log-level", "WARNING", "set log level")
				return fs
			}(),
			args:   []string{"-verbose", "-log-level=INFO"},
			keys:   []string{"log", "level"},
			expect: "INFO",
			ok:     true,
		},
		{
			tc:     "should fail to parse unknown flag",
			flags:  flag.NewFlagSet("test", flag.ContinueOnError),
			args:   []string{"-verbose"},
			keys:   []string{"verbose"},
			expect: nil,
			ok:     false,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			r := &FlagsetResolver{
				Flags:     test.flags,
				Arguments: test.args,
			}

			// iterate twice to ensure we only parse flags once
			for i := 0; i < 2; i++ {
				actual, ok := r.Resolve(test.keys, nil)
				assert.Equal(t, test.ok, ok)
				assert.Equal(t, test.expect, actual)
			}
		})
	}
}
