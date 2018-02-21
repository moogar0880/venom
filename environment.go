package venom

import (
	"os"
	"strings"
)

var defaultEnvResolver = &EnvironmentVariableResolver{}

func toUpperStringSlice(sSlice []string) []string {
	for i, s := range sSlice {
		sSlice[i] = strings.ToUpper(s)
	}
	return sSlice
}

// An EnvironmentVariableResolver is a resolver specifically capable of adding
// additional context in the form of a prefix to any loaded environment
// variables
type EnvironmentVariableResolver struct {
	Prefix string
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
	return os.LookupEnv(strings.Join(toUpperStringSlice(keysCopy), "_"))
}
