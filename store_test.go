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
	t.Run("LogableConfigStore", func(t *testing.T) {
		testVenom(t, NewLogableWith(&TestLogger{}))
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
}

func TestConfigStoreDebug(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testDebug(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testDebug(t, NewSafeConfigStore())
	})
	t.Run("LogableConfigStore", func(t *testing.T) {
		testDebug(t, NewLogableConfigStore())
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
}

func TestConfigStoreAlias(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testAlias(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testAlias(t, NewSafeConfigStore())
	})
	t.Run("LogableConfigStore", func(t *testing.T) {
		testAlias(t, NewLogableWith(&TestLogger{}))
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
}

func TestConfigStoreEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("DefaultConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewDefaultConfigStore())
	})
	t.Run("SafeConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewSafeConfigStore())
	})
	t.Run("LogableConfigStore", func(t *testing.T) {
		testEdgeCases(t, NewLogableWith(&TestLogger{}))
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
}
