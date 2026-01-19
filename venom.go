package venom

// Default the default set of available config levels.
const (
	DefaultLevel     ConfigLevel = iota
	FileLevel                    = DefaultLevel + 10
	EnvironmentLevel             = FileLevel + 10
	FlagLevel                    = EnvironmentLevel + 10
	OverrideLevel    ConfigLevel = 99
)

const defaultDelim = "."

// Delim is the delimiter used for separating nested key spaces. The default is
// to separate on "." characters.
var Delim = defaultDelim

// A KeyTranslator is used for translating portions of config keys in a
// level-agnostic manner.
//
// Functions like these may be useful for performing operations such as
// normalizing hyphens ('-') to underscores ('_') when performing environment
// variable lookups, or perhaps the inverse when performing command line flag
// lookups.
type KeyTranslator func(b byte) byte

// The NoOpKeyTranslator is a KeyTranslator which returns all input bytes
// unmodified.
func NoOpKeyTranslator(b byte) byte {
	return b
}

// ConfigLevel is a type alias used to identify various configuration levels.
type ConfigLevel int

// ConfigMap defines the inner map type which holds actual config data. These
// are nested under a ConfigLevel which determines their priority.
type ConfigMap map[string]interface{}

func (c ConfigMap) merge(d ConfigMap) ConfigMap {
	for key, val := range d {
		switch actual := val.(type) {
		case map[string]interface{}:
			var existing ConfigMap
			if _, ok := c[key]; !ok {
				existing = make(ConfigMap)
			} else {
				existing = c[key].(ConfigMap)
			}

			c[key] = existing.merge(actual)
		case map[interface{}]interface{}:
			var existing ConfigMap
			if _, ok := c[key]; !ok {
				existing = make(ConfigMap)
			} else {
				existing = c[key].(ConfigMap)
			}
			c[key] = existing.merge(mapInterfaceInterfaceToStrInterface(actual))
		default:
			c[key] = val
		}
	}
	return c
}

func mapInterfaceInterfaceToStrInterface(src map[interface{}]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for key, value := range src {
		if actualKey, ok := key.(string); ok {
			data[actualKey] = value
		}
	}
	return data
}

// ConfigLevelMap is a mapping of config levels to the maps which contain
// various configuration values at those levels.
type ConfigLevelMap map[ConfigLevel]ConfigMap

// Venom is the configuration registry responsible for storing and managing
// arbitrary configuration keys and values.
type Venom struct {
	Store ConfigStore
}

// New returns a newly initialized Venom instance.
//
// The internal config map is created empty, only allocating space for a given
// config level once a value is set to that level.
func New() *Venom {
	return NewWithStore(NewDefaultConfigStore())
}

// NewSafe returns a newly initialized Venom instance that is safe to read and
// write from multiple goroutines.
//
// The internal config map is created empty, only allocating space for a given
// config level once a value is set to that level.
func NewSafe() *Venom {
	return NewWithStore(NewSafeConfigStore())
}

// NewLoggableWith takes a Logger and returns a newly initialized Venom
// instance that will log to a Logger interface upon reads and writes.
func NewLoggableWith(l Logger) *Venom {
	lcs := NewLoggableConfigStoreWith(l)
	return NewWithStore(lcs)
}

// NewLoggable returns a Venom instance with a default log set to standard out.
func NewLoggable() *Venom {
	return NewWithStore(NewLoggableConfigStore())
}

// NewWithStore returns a newly initialized Venom instance that wraps the
// provided ConfigStore.
func NewWithStore(s ConfigStore) *Venom {
	return &Venom{
		Store: s,
	}
}

// Default returns a new venom instance with some default resolver
// configuration applied to it.
func Default() *Venom {
	ven := New()
	ven.RegisterResolver(EnvironmentLevel, defaultEnvResolver)
	return ven
}

// DefaultSafe returns a new goroutine-safe venom instance with some default
// resolver configuration applied to it.
func DefaultSafe() *Venom {
	ven := NewSafe()
	ven.RegisterResolver(EnvironmentLevel, defaultEnvResolver)
	return ven
}

// RegisterResolver registers a custom config resolver for the specified
// ConfigLevel.
//
// Additionally, if the provided level is not already in the current collection
// of active config levels, it will be added automatically
func (v *Venom) RegisterResolver(level ConfigLevel, r Resolver) {
	v.Store.RegisterResolver(level, r)
}

// Alias registers an alias for a given key. This allows consumers to access
// the same config via a different key, increasing the backwards
// compatibility of an application.
func (v *Venom) Alias(from, to string) {
	v.Store.Alias(from, to)
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the map, conditionally creating a new ConfigMap if
// one didn't previously exist.
func (v *Venom) SetLevel(level ConfigLevel, key string, value interface{}) {
	v.Store.SetLevel(level, key, value)
}

// SetDefault sets the provided key and value into the DefaultLevel of the
// config collection.
func (v *Venom) SetDefault(key string, value interface{}) {
	v.Store.SetLevel(DefaultLevel, key, value)
}

// SetOverride sets the provided key and value into the OverrideLevel of the
// config collection.
func (v *Venom) SetOverride(key string, value interface{}) {
	v.Store.SetLevel(OverrideLevel, key, value)
}

// Get performs a fetch on a given key from the inner config collection.
func (v *Venom) Get(key string) interface{} {
	val, _ := v.Store.Find(key)
	return val
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether the key was found.
func (v *Venom) Find(key string) (interface{}, bool) {
	return v.Store.Find(key)
}

// Merge merges the provided config map into the ConfigLevel l, allocating
// space for ConfigLevel l if the level hasn't already been allocated.
func (v *Venom) Merge(l ConfigLevel, data ConfigMap) {
	v.Store.Merge(l, data)
}

// Clear removes all data from the ConfigLevelMap and resets the heap of config
// levels.
func (v *Venom) Clear() {
	v.Store.Clear()
}

// Debug returns the current venom ConfigLevelMap as a pretty-printed JSON
// string.
func (v *Venom) Debug() string {
	return v.Store.Debug()
}

// Size returns the number of config levels stored in the underlying
// ConfigStore.
func (v *Venom) Size() int {
	return v.Store.Size()
}
