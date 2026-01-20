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

	// A map containing only the names and values of flags that were actually
	// specified on the command line, as determined by a call to flag.Visit or
	// flag.FlagSet.Visit.
	cachedValueMap map[string]string
}

// parse will parse the specified flagset if it has not already been parsed.
func (r *FlagsetResolver) parse() error {
	// Ensure we cache the set of flags that were specified after parsing the
	// values on the commandline.
	//
	// Note that this function will immediately return if the flag values have
	// already been cached.
	defer r.cacheFlagValues()

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

func (r *FlagsetResolver) cacheFlagValues() {
	// Ensure we don't re-cache the parsed flags every time we check for a key
	// in the wrapped flagset.
	if len(r.cachedValueMap) > 0 {
		return
	}

	// Allocate our map of cached flag values.
	r.cachedValueMap = make(map[string]string)

	// For all flags in the flagset that were actually specified, inject their
	// name and value into our cache of values.
	if r.Flags == nil {
		flag.Visit(func(fl *flag.Flag) {
			r.cachedValueMap[fl.Name] = fl.Value.String()
		})
	} else {
		r.Flags.Visit(func(fl *flag.Flag) {
			r.cachedValueMap[fl.Name] = fl.Value.String()
		})
	}
}

// Resolve is a Resolver function and will lookup the requested config value
// from a FlagSet.
func (r *FlagsetResolver) Resolve(keys []string, _ ConfigMap) (val interface{}, ok bool) {
	// bail early if we fail to parse any of the provided flags
	if err := r.parse(); err != nil {
		return nil, false
	}

	// Leverage our cached map of flags and their values (generated as a part
	// of the call to r.parse() above) rather than iterating over all provided
	// flags every time Resolve is called.
	if value, ok := r.cachedValueMap[strings.Join(keys, FlagSeparator)]; ok {
		return value, ok
	}

	return nil, false
}
