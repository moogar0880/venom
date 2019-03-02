package venom

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVenom(t *testing.T) {
	testIO := []struct {
		inp      []lkv
		expected []kv
	}{
		// simple test case with only the default level
		{
			inp: []lkv{
				{DefaultLevel, "foo", "bar"},
			},
			expected: []kv{
				{"foo", "bar"},
			},
		},
		// slightly more complex test case with use of both the default and
		// override config levels, asserting that the higher priority wins
		{
			inp: []lkv{
				{DefaultLevel, "foo", "bar"},
				{OverrideLevel, "foo", "baz"},
			},
			expected: []kv{
				{"foo", "baz"},
			},
		},
		// simple test case using nested keys
		{
			inp: []lkv{
				{DefaultLevel, "foo.bar", "baz"},
			},
			expected: []kv{
				{"foo.bar", "baz"},
				{"foo", ConfigMap{"bar": "baz"}},
			},
		},
		// more complex test case using multiple config levels and nested key space
		{
			inp: []lkv{
				{DefaultLevel, "foo.bar", "baz"},
				{OverrideLevel, "foo.bar", 12},
			},
			expected: []kv{
				{"foo.bar", 12},
				{"foo", ConfigMap{"bar": 12}},
			},
		},
		// simple case of an absent key
		{
			inp: []lkv{
				{DefaultLevel, "foo", "bar"},
			},
			expected: []kv{
				{"bar", nil},
			},
		},
	}

	for i, test := range testIO {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			v := New()
			for _, inp := range test.inp {
				switch inp.l {
				case DefaultLevel:
					v.SetDefault(inp.k, inp.v)
				case OverrideLevel:
					v.SetOverride(inp.k, inp.v)
				}
			}

			for _, expect := range test.expected {
				if !assert.Equal(t, expect.v, v.Get(expect.k)) {
					fmt.Println(v.Debug())
				}

				_, exists := v.Find(expect.k)
				if expect.v == nil {
					assert.False(t, exists)
				} else {
					assert.True(t, exists)
				}
			}

			// clear the instance and assert that it's empty
			v.Clear()
			assert.Empty(t, v.config)
		})
	}
}

func TestDebug(t *testing.T) {
	testIO := []struct {
		tc      string
		venom   *Venom
		configs ConfigMap
		expect  string
	}{
		{
			tc:    "should debug",
			venom: New(),
			configs: ConfigMap{
				"foo": "bar",
				"baz": map[string]interface{}{
					"bar": "foo",
				},
			},
			expect: `{
  "2": {
    "baz": {
      "bar": "foo"
    },
    "foo": "bar"
  }
}`,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			test.venom.config[EnvironmentLevel] = test.configs

			actual := test.venom.Debug()
			assert.Equal(t, test.expect, actual, actual)
		})
	}
}

func TestVenomEdgeCases(t *testing.T) {
	testIO := []struct {
		tc     string
		v      *Venom
		setup  func(*Venom)
		key    string
		expect interface{}
	}{
		{
			tc: "should load ConfigMap if retrieving first key of nested keys",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo.bar", "baz")
			},
			key:    "foo",
			expect: ConfigMap{"bar": "baz"},
		},
		{
			tc: "should load map if retrieving first key of nested keys",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo", map[string]interface{}{"bar": "baz"})
			},
			key:    "foo",
			expect: map[string]interface{}{"bar": "baz"},
		},
		{
			tc: "should load ConfigMap when retrieving arbitrarily nested keys",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo.bar", "baz")
				v.SetDefault("foo.baz", "bar")
			},
			key: "foo",
			expect: ConfigMap{
				"bar": "baz",
				"baz": "bar",
			},
		},
		{
			tc: "should load ConfigMap when retrieving value with nested sub-keys",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo.bar", "baz")
				v.SetDefault("foo.baz.bat", "bar")
			},
			key: "foo",
			expect: ConfigMap{
				"bar": "baz",
				"baz": ConfigMap{
					"bat": "bar",
				},
			},
		},
		{
			tc: "should return nil when joining string key/value during lookup",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo", "bar.baz")
			},
			key:    "foo.bar",
			expect: nil,
		},
		{
			tc: "should return nil when joining string/ConfigMap key/value during lookup",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo", ConfigMap{"bar.baz": "bat"})
			},
			key:    "foo.bar",
			expect: nil,
		},
		{
			tc: "should return nil when joining string/map key/value during lookup",
			v:  New(),
			setup: func(v *Venom) {
				v.SetDefault("foo", map[string]interface{}{"bar.baz": "bat"})
			},
			key:    "foo.bar",
			expect: nil,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			test.setup(test.v)

			actual := test.v.Get(test.key)
			assert.Equal(t, test.expect, actual)
		})
	}
}

func TestVenomAliases(t *testing.T) {
	testIO := []struct {
		tc    string
		key   string
		alias string
		value interface{}
	}{
		{
			tc:    "should resolve aliased field",
			key:   "foo",
			alias: "bar",
			value: 867.5309,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			ven := New()
			ven.SetDefault(test.key, test.value)
			ven.Alias(test.alias, test.key)

			unAliased := ven.Get(test.key)
			assert.Equal(t, test.value, unAliased)

			aliased := ven.Get(test.alias)
			assert.Equal(t, test.value, aliased)
		})
	}
}

var readResult interface{}

func benchmarkVenomRead(ven *Venom, key string, b *testing.B) {
	var v interface{}
	for n := 0; n < b.N; n++ {
		v = ven.Get(key)
	}
	readResult = v
}

type keyer func(i int) string
type leveler func(i int) ConfigLevel

func benchmarkVenomWrite(ven *Venom, l leveler, key keyer, val interface{}, b *testing.B) {
	for n := 0; n < b.N; n++ {
		ven.SetLevel(l(n), key(n), val)
	}
}

func BenchmarkVenomGet(b *testing.B) {
	envVarsToUnset := make([]string, 0)

	defer func() {
		// clean up any env vars
		for _, v := range envVarsToUnset {
			os.Unsetenv(v)
		}
		envVarsToUnset = []string{}
	}()

	testIO := []struct {
		tc  string
		ven *Venom
		key string
	}{
		{
			tc: "single ConfigLevel with one key/value pair",
			ven: func() *Venom {
				ven := New()
				ven.SetDefault("test", 100)
				return ven
			}(),
			key: "test",
		},
		{
			tc: "many key/value pairs in a single ConfigLevel",
			ven: func() *Venom {
				ven := New()
				for i := 0; i < 10000; i++ {
					ven.SetDefault(fmt.Sprintf("test_%d", i), i)
				}
				return ven
			}(),
			key: "test_5000",
		},
		{
			tc: "many key/value pairs spread across multiple ConfigLevels",
			ven: func() *Venom {
				ven := New()
				var key string
				for i := 0; i < 10000; i++ {
					key = fmt.Sprintf("test_%d", i)
					if i%int(FileLevel) == 0 {
						ven.SetLevel(FileLevel, key, i)
					} else if i%int(EnvironmentLevel) == 0 {
						os.Setenv(strings.ToUpper(key), fmt.Sprintf("%d", i))
						envVarsToUnset = append(envVarsToUnset, strings.ToUpper(key))
					} else if i%int(OverrideLevel) == 0 {
						ven.SetOverride(key, i)
					} else {
						ven.SetDefault(key, i)
					}
				}
				return ven
			}(),
			key: "test_5000",
		},
	}

	for _, test := range testIO {
		b.Run(test.tc, func(b *testing.B) {
			benchmarkVenomRead(test.ven, test.key, b)
		})
	}
}

func BenchmarkVenomWrite(b *testing.B) {
	envVarsToUnset := make([]string, 0)

	defer func() {
		// clean up any env vars
		for _, v := range envVarsToUnset {
			os.Unsetenv(v)
		}
		envVarsToUnset = []string{}
	}()

	testIO := []struct {
		tc    string
		ven   *Venom
		level leveler
		key   keyer
		value interface{}
	}{
		{
			tc:    "single key/value pair in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(int) string { return "test" },
			value: 1000,
		},
		{
			tc:    "many key/value pairs in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(i int) string { return fmt.Sprintf("test_%d", i) },
			value: 1000,
		},
		{
			tc:    "many nested key/value pairs in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(i int) string { return fmt.Sprintf("test_%d.test.%d", i, i) },
			value: 1000,
		},
		{
			tc:  "many key/value pairs in many ConfigLevels",
			ven: New(),
			level: func(i int) ConfigLevel {
				if i%int(FileLevel) == 0 {
					return FileLevel
				} else if i%int(EnvironmentLevel) == 0 {
					envVarsToUnset = append(envVarsToUnset, strings.ToUpper(fmt.Sprintf("test_%d", i)))
					return EnvironmentLevel
				} else if i%int(OverrideLevel) == 0 {
					return OverrideLevel
				}
				return DefaultLevel
			},
			key:   func(i int) string { return fmt.Sprintf("test_%d", i) },
			value: 1000,
		},
		{
			tc:  "many nested key/value pairs in many ConfigLevels",
			ven: New(),
			level: func(i int) ConfigLevel {
				if i%int(FileLevel) == 0 {
					return FileLevel
				} else if i%int(EnvironmentLevel) == 0 {
					envVarsToUnset = append(envVarsToUnset, strings.ToUpper(fmt.Sprintf("test_%d", i)))
					return EnvironmentLevel
				} else if i%int(OverrideLevel) == 0 {
					return OverrideLevel
				}
				return DefaultLevel
			},
			key:   func(i int) string { return fmt.Sprintf("test_%d.test.%d", i, i) },
			value: 1000,
		},
	}

	for _, test := range testIO {
		b.Run(test.tc, func(b *testing.B) {
			benchmarkVenomWrite(test.ven, test.level, test.key, test.value, b)
		})
	}
}
