package venom

import (
	"container/heap"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

// The ConfigStore interface defines a type capable of performing basic
// operations required for configuration management within a venom
// instance.
type ConfigStore interface {
	RegisterResolver(level ConfigLevel, r Resolver)
	SetLevel(level ConfigLevel, key string, value interface{})
	Merge(l ConfigLevel, data ConfigMap)
	Alias(from, to string)
	Find(key string) (interface{}, bool)
	Clear()
	Debug() string
	Size() int
}

// DefaultConfigStore is the minimum implementation of a ConfigStore. It is
// capable of storing and managing arbitrary configuration keys and values.
type DefaultConfigStore struct {
	// config is the global config space. Values are stored at a specified
	// ConfigLevel for prioritized retrieval
	config ConfigLevelMap

	// usedLevels is a sorted slice of all ConfigLevels currently stored in the
	// config map
	usedLevels *ConfigLevelHeap

	// resolvers is the definitive list of any customer ConfigLevel resolvers
	// provided to this Venom instance
	resolvers map[ConfigLevel]Resolver

	// aliases contains the collection of any aliased config values
	aliases map[string]string
}

// NewDefaultConfigStore returns a newly allocated DefaultConfigStore.
func NewDefaultConfigStore() *DefaultConfigStore {
	return &DefaultConfigStore{
		config:     make(ConfigLevelMap),
		usedLevels: NewConfigLevelHeap(),
		resolvers:  make(map[ConfigLevel]Resolver),
		aliases:    make(map[string]string),
	}
}

// RegisterResolver registers a custom config resolver for the specified
// ConfigLevel.
//
// Additionally, if the provided level is not already in the current collection
// of active config levels, it will be added automatically
func (s *DefaultConfigStore) RegisterResolver(level ConfigLevel, r Resolver) {
	s.resolvers[level] = r
	heap.Push(s.usedLevels, level)
}

// Alias registers an alias for a given key. This allows consumers to access
// the same config via a different key, increasing the backwards
// compatibility of an application.
func (s *DefaultConfigStore) Alias(from, to string) {
	if from == to {
		return
	}
	s.aliases[from] = to
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the map, conditionally creating a new ConfigMap if
// one didn't previously exist.
func (s *DefaultConfigStore) SetLevel(level ConfigLevel, key string, value interface{}) {
	s.setIfNotExists(level, key, value)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found
func (s *DefaultConfigStore) Find(key string) (interface{}, bool) {
	return s.find(key)
}

// Merge merges the provided config map into the ConfigLevel l, allocating
// space for ConfigLevel l if the level hasn't already been allocated.
func (s *DefaultConfigStore) Merge(l ConfigLevel, data ConfigMap) {
	if _, ok := s.config[l]; !ok {
		s.config[l] = make(ConfigMap)
		heap.Push(s.usedLevels, l)
	}
	s.config[l].merge(data)
}

// Size returns the number of config levels stored in this ConfigStore.
func (s *DefaultConfigStore) Size() int {
	return len(s.config)
}

func (s *DefaultConfigStore) setIfNotExists(l ConfigLevel, key string, value interface{}) {
	if _, ok := s.config[l]; !ok {
		s.config[l] = make(ConfigMap)
		heap.Push(s.usedLevels, l)
	}
	setNested(s.config[l], strings.Split(key, Delim), value)
}

// setNested inserts the provided value into the nested keyspace as defined by
// the delim separated keys
func setNested(config ConfigMap, keys []string, value interface{}) {
	for index, key := range keys {
		// if we're at the end of our slice of keys, set the value and return
		if index == len(keys)-1 {
			config[key] = value
			return
		}

		// make sure we won't overwrite existing keys, before creating a new
		// ConfigMap at the current node and continuing
		if _, ok := config[key]; !ok {
			config[key] = make(ConfigMap)
		}
		config = config[key].(ConfigMap)
	}
}

func (s *DefaultConfigStore) find(key string) (val interface{}, ok bool) {
	// check for aliases before beginning search
	if actual, isAliased := s.aliases[key]; isAliased {
		key = actual
	}

	keys := strings.Split(key, Delim)
	for _, level := range *s.usedLevels {
		resolver, resolverExists := s.resolvers[level]
		if !resolverExists {
			resolver = defaultResolver
		}

		if val, ok = resolver.Resolve(keys, s.config[level]); ok {
			return
		}
	}
	return nil, false
}

// Clear removes all data from the ConfigLevelMap and resets the heap of config
// levels.
func (s *DefaultConfigStore) Clear() {
	s.config = make(ConfigLevelMap)
	s.usedLevels = NewConfigLevelHeap()
}

// Debug returns the current venom ConfigLevelMap as a pretty-printed JSON
// string.
func (s *DefaultConfigStore) Debug() string {
	b, _ := json.MarshalIndent(s.config, "", "  ")
	return string(b)
}

// SafeConfigStore implements the ConfigStore interface and provides a store
// that is safe to read and write from multiple go routines.
type SafeConfigStore struct {
	c  *DefaultConfigStore
	mu sync.Mutex
}

// NewSafeConfigStore returns a new SafeConfigStore.
func NewSafeConfigStore() ConfigStore {
	return &SafeConfigStore{
		c: NewDefaultConfigStore(),
	}
}

// RegisterResolver registers a custom config resolver for the specified
// ConfigLevel.
//
// Additionally, if the provided level is not already in the current collection
// of active config levels, it will be added automatically
func (s *SafeConfigStore) RegisterResolver(level ConfigLevel, r Resolver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.RegisterResolver(level, r)
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the map, conditionally creating a new ConfigMap if
// one didn't previously exist.
func (s *SafeConfigStore) SetLevel(level ConfigLevel, key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.SetLevel(level, key, value)
}

// Merge merges the provided config map into the ConfigLevel l, allocating
// space for ConfigLevel l if the level hasn't already been allocated.
func (s *SafeConfigStore) Merge(l ConfigLevel, data ConfigMap) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.Merge(l, data)
}

// Alias registers an alias for a given key. This allows consumers to access
// the same config via a different key, increasing the backwards
// compatibility of an application.
func (s *SafeConfigStore) Alias(from, to string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.Alias(from, to)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found
func (s *SafeConfigStore) Find(key string) (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.Find(key)
}

// Clear removes all data from the ConfigLevelMap and resets the heap of config
// levels.
func (s *SafeConfigStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.c.Clear()
}

// Debug returns the current venom ConfigLevelMap as a pretty-printed JSON
// string.
func (s *SafeConfigStore) Debug() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.Debug()
}

// Size returns the number of config levels stored in this ConfigStore.
func (s *SafeConfigStore) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.Size()
}

// LoggableConfigStore implements the ConfigStore interface and provides a store
// a field for reference to a logging mechanism to log on reads/writes
type LoggableConfigStore struct {
	c   ConfigStore
	log Logging
}

// NewLoggableConfigStore takes a Logging and returns a new ConfigStore
// with said interface as the logging mechanism used for read and writes.
func NewLoggableConfigStoreWith(l Logging) ConfigStore {
	return &LoggableConfigStore{
		c: NewDefaultConfigStore(),
		log: l,
	}
}

// NewLoggableConfigStore returns a ConfigStore with a default logging mechanism
// set to write to os.Stdout.
func NewLoggableConfigStore() ConfigStore {
	l := log.New(os.Stdout, "", 0)
	return NewLoggableConfigStoreWith(l)
}

// RegisterResolver registers a custom config resolver for the specified
// ConfigLevel.
//
// Additionally, if the provided level is not already in the current collection
// of active config levels, it will be added automatically
func (l *LoggableConfigStore) RegisterResolver(level ConfigLevel, r Resolver) {
	l.c.RegisterResolver(level, r)
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the map, conditionally creating a new ConfigMap if
// one didn't previously exist.
func (l *LoggableConfigStore) SetLevel(level ConfigLevel, key string, value interface{}) {
	l.c.SetLevel(level, key, value)
	l.log.Print(formatLogLine("SET", level, key, value))
}

// Merge merges the provided config map into the ConfigLevel l, allocating
// space for ConfigLevel l if the level hasn't already been allocated.
func (l *LoggableConfigStore) Merge(cl ConfigLevel, data ConfigMap) {
	l.c.Merge(cl, data)
}

// Alias registers an alias for a given key. This allows consumers to access
// the same config via a different key, increasing the backwards
// compatibility of an application.
func (l *LoggableConfigStore) Alias(from, to string) {
	l.c.Alias(from, to)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found
func (l *LoggableConfigStore) Find(key string) (interface{}, bool) {
	a, b := l.c.Find(key)
	l.log.Print(formatLogLine("GET", a, b))
	return a, b
}

// Clear removes all data from the ConfigLevelMap and resets the heap of config
// levels.
func (l *LoggableConfigStore) Clear() {
	l.c.Clear()
}

// Debug returns the current venom ConfigLevelMap as a pretty-printed JSON
// string.
func (l *LoggableConfigStore) Debug() string {
	return l.c.Debug()
}

// Size returns the number of config levels stored in this ConfigStore.
func (l *LoggableConfigStore) Size() int {
	return l.c.Size()
}
