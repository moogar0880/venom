package venom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testVenom(t *testing.T, v ConfigStore) {
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
			for _, inp := range test.inp {
				v.SetLevel(inp.l, inp.k, inp.v)
			}

			for _, expect := range test.expected {
				val, _ := v.Find(expect.k)
				if !assert.Equal(t, expect.v, val) {
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
			assert.Equal(t, 0, v.Size())
		})
	}
}

func testDebug(t *testing.T, v ConfigStore) {
	testIO := []struct {
		tc      string
		configs ConfigMap
		expect  string
	}{
		{
			tc:      "should debug empty config map",
			configs: ConfigMap{},
			expect: `{
  "2": {}
}`,
		},
		{
			tc: "should debug populated ConfigMap",
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
			v.Merge(EnvironmentLevel, test.configs)

			actual := v.Debug()
			assert.Equal(t, test.expect, actual, actual)
		})
	}
}

func testAlias(t *testing.T, v ConfigStore) {
	testIO := []struct {
		tc       string
		from     string
		to       string
		key      string
		expect   interface{}
		expectOK bool
	}{
		{
			tc:       "should alias foo to bar",
			from:     "bar",
			to:       "foo",
			key:      "bar",
			expect:   12,
			expectOK: true,
		},
		{
			tc:       "should not alias foo to itself",
			from:     "foo",
			to:       "foo",
			key:      "foo",
			expect:   12,
			expectOK: true,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			v.Alias(test.from, test.to)
			v.SetLevel(OverrideLevel, "foo", 12)

			actual, ok := v.Find(test.key)
			assert.Equal(t, test.expectOK, ok)
			assert.Equal(t, test.expect, actual)
		})
	}
}

func testEdgeCases(t *testing.T, v ConfigStore) {
	testIO := []struct {
		tc       string
		setup    func(ConfigStore)
		key      string
		expect   interface{}
		expectOK bool
	}{
		{
			tc: "should load ConfigMap if retrieving first key of nested keys",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo.bar", "baz")
			},
			key:      "foo",
			expect:   ConfigMap{"bar": "baz"},
			expectOK: true,
		},
		{
			tc: "should load map if retrieving first key of nested keys",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo", map[string]interface{}{"bar": "baz"})
			},
			key:      "foo",
			expect:   map[string]interface{}{"bar": "baz"},
			expectOK: true,
		},
		{
			tc: "should load ConfigMap when retrieving arbitrarily nested keys",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo.bar", "baz")
				v.SetLevel(DefaultLevel, "foo.baz", "bar")
			},
			key: "foo",
			expect: ConfigMap{
				"bar": "baz",
				"baz": "bar",
			},
			expectOK: true,
		},
		{
			tc: "should load ConfigMap when retrieving value with nested sub-keys",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo.bar", "baz")
				v.SetLevel(DefaultLevel, "foo.baz.bat", "bar")
			},
			key: "foo",
			expect: ConfigMap{
				"bar": "baz",
				"baz": ConfigMap{
					"bat": "bar",
				},
			},
			expectOK: true,
		},
		{
			tc: "should return nil when joining string key/value during lookup",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo", "bar.baz")
			},
			key:      "foo.bar",
			expect:   nil,
			expectOK: false,
		},
		{
			tc: "should return nil when joining string/ConfigMap key/value during lookup",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo", ConfigMap{"bar.baz": "bat"})
			},
			key:      "foo.bar",
			expect:   nil,
			expectOK: false,
		},
		{
			tc: "should return nil when joining string/map key/value during lookup",
			setup: func(v ConfigStore) {
				v.SetLevel(DefaultLevel, "foo", map[string]interface{}{"bar.baz": "bat"})
			},
			key:      "foo.bar",
			expect:   nil,
			expectOK: false,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			test.setup(v)

			actual, ok := v.Find(test.key)
			assert.Equal(t, test.expectOK, ok)
			assert.Equal(t, test.expect, actual)

			v.Clear()
		})
	}
}
