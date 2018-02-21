package venom

var defaultResolver = &DefaultResolver{}

// A Resolver is responsible for performing a lookup of a config returning the
// value stored for that config or an error.
type Resolver interface {
	Resolve([]string, ConfigMap) (val interface{}, ok bool)
}

// DefaultResolver is the default resolver function used to resolve
// configuration values for a level which does not specify a custom resolver.
type DefaultResolver struct{}

// Resolve will attempt to resolve the specified key using the configuration
// data stored in the provided ConfigMap
func (r *DefaultResolver) Resolve(keys []string, config ConfigMap) (val interface{}, ok bool) {
	for _, key := range keys {
		if val, ok = config[key]; ok {
			// if we're at the last key in the slice return the current value
			if len(keys) == 1 {
				return
			}

			switch actualValue := val.(type) {
			case ConfigMap:
				return r.Resolve(keys[1:], actualValue)
			}
		}
	}
	return nil, false
}
