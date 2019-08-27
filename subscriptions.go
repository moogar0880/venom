package venom

import (
	"fmt"
	"strings"
)

// An Event represents an update published to a ConfigStore. The Key represents
// the config value that was updated. The Value contains the value that the key
// was updated to.
type Event struct {
	Key   string
	Value interface{}
}

// A SubscriptionStore allows callers to subscribe to specific key namespaces,
// allowing them to be notified when a config under that space is modified.
//
// Key-spaces can be subscribed to via the Subscribe method, which return a
// channel over which events for that key-space are emitted. This means that an
// update to the config `db.host` would trigger an event to be emitted if a
// subscription existed on either `db` or `db.host`.
type SubscriptionStore struct {
	channels map[string]chan Event
	store    ConfigStore
	bufSize  int
}

// NewSubscriptionStore returns a newly allocated SubscriptionStore which wraps
// the specified ConfigStore. The channels created by this method will be
// un-buffered by default. To customize the buffer size of these channels, see
// NewSubscriptionStoreWithSize.
func NewSubscriptionStore(s ConfigStore) (*SubscriptionStore, func()) {
	return NewSubscriptionStoreWithSize(s, 0)
}

// NewSubscriptionStoreWithSize returns a newly allocated SubscriptionStore
// which wraps the provided ConfigStore as well as a closure which, when
// called, will close all channels managed by this SubscriptionStore. The size
// parameter of this function allows for the caller to toggle the buffer size
// of the channels created by this ConfigStore.
func NewSubscriptionStoreWithSize(s ConfigStore, size int) (*SubscriptionStore, func()) {
	store := &SubscriptionStore{
		channels: make(map[string]chan Event),
		store:    s,
		bufSize:  size,
	}
	return store, store.Close
}

// RegisterResolver registers a custom config resolver for the specified
// ConfigLevel.
//
// Additionally, if the provided level is not already in the current collection
// of active config levels, it will be added automatically
func (s *SubscriptionStore) RegisterResolver(level ConfigLevel, r Resolver) {
	s.store.RegisterResolver(level, r)
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the map, conditionally creating a new ConfigMap if
// one didn't previously exist.
//
// Once written, a new event will be emitted by this Subscription store, if any
// matching key-spaces have subscription channels.
func (s *SubscriptionStore) SetLevel(level ConfigLevel, key string, value interface{}) {
	s.store.SetLevel(level, key, value)
	s.emit(key, value)
}

// Merge merges the provided config map into the ConfigLevel l, allocating
// space for ConfigLevel l if the level hasn't already been allocated.
func (s *SubscriptionStore) Merge(l ConfigLevel, data ConfigMap) {
	s.store.Merge(l, data)
}

// Alias registers an alias for a given key. This allows consumers to access
// the same config via a different key, increasing the backwards
// compatibility of an application.
func (s *SubscriptionStore) Alias(from, to string) {
	s.store.Alias(from, to)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found.
func (s *SubscriptionStore) Find(key string) (interface{}, bool) {
	return s.store.Find(key)
}

// Clear removes all data from the ConfigLevelMap and resets the heap of config
// levels.
func (s *SubscriptionStore) Clear() {
	s.store.Clear()
}

// Debug returns the current venom ConfigLevelMap as a pretty-printed JSON
// string.
func (s *SubscriptionStore) Debug() string {
	return s.store.Debug()
}

// Size returns the number of config levels stored in this ConfigStore.
func (s *SubscriptionStore) Size() int {
	return s.store.Size()
}

// Close iterates over all allocated channels, closes them, and removes the
// subscriptions from the map of event subscriptions.
func (s *SubscriptionStore) Close() {
	for subscriptionKey, channel := range s.channels {
		close(channel)
		delete(s.channels, subscriptionKey)
	}
}

// Subscribe returns a channel over which updates to any value located at, or
// under, the specified key will be emitted.
//
// Note that if no buffer size was specified when the SubscriptionStore was
// created, the channels will default to being un-buffered.
//
// Also note that multiple subscriptions to the same key will result in the
// same channel being returned. To clear an existing subscription and to close
// it's corresponding channel, use the Unsubscribe method.
//
// Subscriptions to the root space (empty string) will result in events being
// emitted for any and all config updates.
func (s *SubscriptionStore) Subscribe(key string) <-chan Event {
	if channel, ok := s.channels[key]; ok {
		return channel
	}

	var newChan chan Event
	if s.bufSize == 0 {
		newChan = make(chan Event)
	} else {
		newChan = make(chan Event, s.bufSize)
	}
	s.channels[key] = newChan
	return s.channels[key]
}

// Unsubscribe removes an existing subscription. The removal of this
// subscription results in the subscription channel being closed, and the
// subscription being completely removed from this Store.
//
// To remove all existing subscriptions, use Close.
func (s *SubscriptionStore) Unsubscribe(key string) error {
	// if the channel exists in the map, close it and remove the subscription
	// from the map
	if channel, ok := s.channels[key]; ok {
		close(channel)
		delete(s.channels, key)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("venom: no such subscription: %s", key)
}

// emit emits an update event, if a subscription was made to the updated key or
// to any of it's parent key-spaces, for every key-space that matches.
//
// For example, if a subscription is made to the `db` space and an update is
// made to `db.host`, first a subscription to `db.host` is checked, since there
// is no such subscription no event is emitted. Then `db` is checked, since
// there is a subscription, an event is emitted.
//
// In the above example, if there were subscriptions on both `db` and `db.host`
// then both channels would have unique events emitted over them.
func (s *SubscriptionStore) emit(key string, value interface{}) {
	keys := strings.Split(key, Delim)
	for i := len(keys); i >= 0; i-- {
		if channel, ok := s.channels[strings.Join(keys[:i], Delim)]; ok {
			channel <- Event{
				Key:   key,
				Value: value,
			}
		}
	}
}
