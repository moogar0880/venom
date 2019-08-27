package venom

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionStore_Subscribe(t *testing.T) {
	testIO := []struct {
		name         string
		init         func(ven *Venom)
		subscribeKey string
		expect       []Event
		updates      func(ven *Venom)
	}{
		{
			name: "should track updates when subscribed to parent space",
			init: func(ven *Venom) {
				ven.SetDefault("db.host", "localhost")
				ven.SetDefault("db.port", "1234")
			},
			subscribeKey: "db",
			expect: []Event{
				{
					Key:   "db.host",
					Value: "example.com",
				},
			},
			updates: func(ven *Venom) {
				ven.SetOverride("db.host", "example.com")
			},
		},
		{
			name: "should track updates when subscribed to deeply nested parent space",
			init: func(ven *Venom) {
				ven.SetDefault("db.connection.details.host", "localhost")
				ven.SetDefault("db.connection.details.port", "1234")
			},
			subscribeKey: "db",
			expect: []Event{
				{
					Key:   "db.connection.details.host",
					Value: "example.com",
				},
			},
			updates: func(ven *Venom) {
				ven.SetOverride("db.connection.details.host", "example.com")
			},
		},
		{
			name: "should track all updates when subscribed to root space",
			init: func(ven *Venom) {
				ven.SetDefault("db.host", "localhost")
				ven.SetDefault("db.port", "1234")
			},
			subscribeKey: "",
			expect: []Event{
				{
					Key:   "db.host",
					Value: "example.com",
				},
			},
			updates: func(ven *Venom) {
				ven.SetOverride("db.host", "example.com")
			},
		},
	}

	for _, test := range testIO {
		t.Run(test.name, func(t *testing.T) {
			// create our subscription store and ensure we clean it up when the
			// test case has completed
			store, clear := NewSubscriptionStore(NewDefaultConfigStore())
			defer clear()

			// create a venom instance that wraps the new store, then
			// initialize it specifically for this test case and subscribe to
			// the test key space
			ven := NewWithStore(store)
			test.init(ven)
			events := store.Subscribe(test.subscribeKey)

			// subscribe to changes on our unbuffered channel before attempting
			// to write updates to it
			done := make(chan bool, 1)
			go func() {
				for i, expect := range test.expect {
					assert.Equal(t, expect, <-events)
					if i == len(test.expect)-1 {
						done <- true
					}
				}
			}()

			// apply the updates for the current test case
			test.updates(ven)

			// wait for up to 2 seconds for the test to complete. if the test
			// does not complete in the allotted time (the expected events were
			// not emitted over the channel) then the test is immediately
			// marked as failed
			tick := time.NewTicker(2 * time.Second)
			select {
			case <-done:
				tick.Stop()
			case <-tick.C:
				tick.Stop()
				t.Error("test timed out waiting for events")
			}

			// ensure that we can properly unsubscribe from our key-space once
			// we're done testing it
			assert.Nil(t, store.Unsubscribe(test.subscribeKey))
		})
	}
}
