package venom

import (
	"os"
	"strings"
	"unicode"
)

var defaultEnvResolver = &EnvironmentVariableResolver{}

// EnvSeparator is used as the delimiter for separating keys prior to looking
// them up in the current environment.
//
// ie, an EnvSeparator of "_" will result in a lookup for "log.level" searching
// for an environment variable named "LOG_LEVEL".
var EnvSeparator = "_"

// An EnvironmentVariableResolver is a resolver specifically capable of adding
// additional context in the form of a prefix to any loaded environment
// variables
type EnvironmentVariableResolver struct {
	Prefix     string
	Translator KeyTranslator
}

// Resolve is a Resolver implementation which attempts to load the requested
// configuration from an environment variable
func (r *EnvironmentVariableResolver) Resolve(keys []string, _ ConfigMap) (val interface{}, ok bool) {
	// copy the keys so we don't negatively impact subsequent lookups
	keysCopy := make([]string, len(keys))
	copy(keysCopy, keys)

	if len(r.Prefix) > 0 {
		keysCopy = append([]string{r.Prefix}, keysCopy...)
	}

	var translator KeyTranslator
	if r.Translator == nil {
		// If we weren't given a specific translator to use, then use the
		// default translator.
		translator = DefaultEnvironmentVariableKeyTranslator
	} else {
		// Otherwise, use the translator that was provided when this resolver
		// was created.
		translator = r.Translator
	}

	return os.LookupEnv(toEnvironmentVariable(keysCopy, translator))
}

// The DefaultEnvironmentVariableKeyTranslator is the default KeyTranslator
// used by the EnvironmentVariableResolver.
//
// It is responsible for mapping arbitrary input keys to an environment
// variable to lookup. This is done by converting all characters to upper case
// and replacing any hyphens with underscores.
func DefaultEnvironmentVariableKeyTranslator(b byte) byte {
	switch b {
	case '-':
		return '_'
	default:
		return byte(unicode.ToUpper(rune(b)))
	}
}

func toEnvironmentVariable(keys []string, translator KeyTranslator) string {
	// Convert the input keys into a single environment variable that we can
	// perform a lookup on.
	key := []byte(strings.Join(keys, EnvSeparator))

	// Next, perform any custom translations performed by the KeyTranslator.
	for index, char := range key {
		key[index] = translator(char)
	}

	// Finally, return the translated key to the caller.
	return string(key)
}
