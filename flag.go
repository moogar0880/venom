package venom

import (
	"flag"
	"os"
	"strings"
)

// FlagSeparator is used as the delimiter for separating keys prior to looking
// up their value from a FlagSet.
//
// ie, a FlagSeparator of "-" will result in a lookup for "log.level" searching
// for a flag named "log-level".
var FlagSeparator = "-"

// The FlagsetResolver is used to resolve configuration from a FlagSet using
// the standard flag library.
//
// If no FlagSet is provided, then flag.Parse and flag.Lookup are used to pull
// values directly from os.Args. If a FlagSet is provided, then any arguments
// required when Parsing the FlagSet may also be specified, otherwise the
// arguments default to os.Args[1:].
type FlagsetResolver struct {
	Flags     *flag.FlagSet
	Arguments []string
}

// parse will parse the specified flagset if it has not already been parsed
func (r *FlagsetResolver) parse() error {
	// if no flagset was provided, use the default parser which parses
	// os.Args[1:] by default
	if r.Flags == nil {
		if flag.Parsed() {
			return nil
		}
		flag.Parse()
		return nil
	}

	// if our flagset is already parsed, we can return
	if r.Flags.Parsed() {
		return nil
	}

	// our flagset needs to be parsed. if no arguments were specified, fall
	// back on os.Args
	if r.Arguments == nil {
		return r.Flags.Parse(os.Args[1:])
	}
	return r.Flags.Parse(r.Arguments)
}

// Resolve is a Resolver function and will lookup the requested config value
// from a FlagSet
func (r *FlagsetResolver) Resolve(keys []string, _ ConfigMap) (val interface{}, ok bool) {
	// bail early if we fail to parse any of the provided flags
	if err := r.parse(); err != nil {
		return nil, false
	}

	var f *flag.Flag
	if r.Flags == nil {
		f = flag.Lookup(strings.Join(keys, FlagSeparator))
	} else {
		f = r.Flags.Lookup(strings.Join(keys, FlagSeparator))
	}

	if f != nil {
		return f.Value.String(), true
	}
	return nil, false
}
