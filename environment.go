package venom

import (
	"os"
	"strings"
)

// EnvReplacer is the strings.Replacer used to map venom keys to their
// corresponding environment variables
var EnvReplacer = strings.NewReplacer(Delim, "_")

func findKeyInEnvironment(key string) (interface{}, bool) {
	return os.LookupEnv(EnvReplacer.Replace(key))
}

// RegisterEnv registers the key to be extracted from the environment when
// LoadEnvironment is run
func (v *Venom) RegisterEnv(key string) {
	v.envKeys = append(v.envKeys, key)
}

// LoadEnvironment reads all registered environment keys into the ConfigMap at
// the EnvironmentLevel
func (v *Venom) LoadEnvironment() {
	// iterate over registered keys and insert into map at EnvironmentLevel
	for _, key := range v.envKeys {
		if value, ok := findKeyInEnvironment(key); ok {
			v.setIfNotExists(EnvironmentLevel, key, value)
		}
	}
}
