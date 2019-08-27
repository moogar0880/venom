package venom

import "testing"

func TestConfigStoreSetAndFind(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testVenom(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testVenom(t, NewSafeConfigStore())
	})
	t.Run("LoggableConfigStore", func(t *testing.T) {
		testVenom(t, NewLoggableWith(&TestLogger{}))
	})
	t.Run("Venom", func(t *testing.T) {
		testVenom(t, New())
	})
	t.Run("DefaultVenom", func(t *testing.T) {
		testVenom(t, Default())
	})
	t.Run("SafeVenom", func(t *testing.T) {
		testVenom(t, NewSafe())
	})
	t.Run("SafeDefaultVenom", func(t *testing.T) {
		testVenom(t, DefaultSafe())
	})
	t.Run("SubscriptionStore", func(t *testing.T) {
		store, clear := NewSubscriptionStore(NewDefaultConfigStore())
		defer clear()
		testVenom(t, store)
	})
}

func TestConfigStoreDebug(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testDebug(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testDebug(t, NewSafeConfigStore())
	})
	t.Run("LoggableConfigStore", func(t *testing.T) {
		testDebug(t, NewLoggableConfigStore())
	})
	t.Run("Venom", func(t *testing.T) {
		testDebug(t, New())
	})
	t.Run("DefaultVenom", func(t *testing.T) {
		testDebug(t, Default())
	})
	t.Run("SafeVenom", func(t *testing.T) {
		testDebug(t, NewSafe())
	})
	t.Run("SafeDefaultVenom", func(t *testing.T) {
		testDebug(t, DefaultSafe())
	})
	t.Run("SubscriptionStore", func(t *testing.T) {
		store, clear := NewSubscriptionStore(NewDefaultConfigStore())
		defer clear()
		testDebug(t, store)
	})
}

func TestConfigStoreAlias(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testAlias(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testAlias(t, NewSafeConfigStore())
	})
	t.Run("LoggableConfigStore", func(t *testing.T) {
		testAlias(t, NewLoggableWith(&TestLogger{}))
	})
	t.Run("Venom", func(t *testing.T) {
		testAlias(t, New())
	})
	t.Run("DefaultVenom", func(t *testing.T) {
		testAlias(t, Default())
	})
	t.Run("SafeVenom", func(t *testing.T) {
		testAlias(t, NewSafe())
	})
	t.Run("SafeDefaultVenom", func(t *testing.T) {
		testAlias(t, DefaultSafe())
	})
	t.Run("SubscriptionStore", func(t *testing.T) {
		store, clear := NewSubscriptionStore(NewDefaultConfigStore())
		defer clear()
		testAlias(t, store)
	})
}

func TestConfigStoreEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewSafeConfigStore())
	})
	t.Run("LoggableConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewLoggableWith(&TestLogger{}))
	})
	t.Run("Venom", func(t *testing.T) {
		testEdgeCases(t, New())
	})
	t.Run("DefaultVenom", func(t *testing.T) {
		testEdgeCases(t, Default())
	})
	t.Run("SafeVenom", func(t *testing.T) {
		testEdgeCases(t, NewSafe())
	})
	t.Run("SafeDefaultVenom", func(t *testing.T) {
		testEdgeCases(t, DefaultSafe())
	})
	t.Run("SubscriptionStore", func(t *testing.T) {
		store, clear := NewSubscriptionStore(NewDefaultConfigStore())
		defer clear()
		testEdgeCases(t, store)
	})
}
