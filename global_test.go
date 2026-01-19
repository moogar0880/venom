package venom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalVenom(t *testing.T) {
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
		// // simple test case using nested keys
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
			// reset the global venom instance for each test run
			v = Default()

			for _, inp := range test.inp {
				switch inp.l {
				case DefaultLevel:
					SetDefault(inp.k, inp.v)
				case OverrideLevel:
					SetOverride(inp.k, inp.v)
				}
			}

			for _, expect := range test.expected {
				if !assert.Equal(t, expect.v, Get(expect.k)) {
					fmt.Println(Debug())
				}

				_, exists := Find(expect.k)
				if expect.v == nil {
					assert.False(t, exists)
				} else {
					assert.True(t, exists)
				}
			}

			// clear the instance and assert that it's empty
			Clear()
			st := v.Store.(*DefaultConfigStore)
			assert.Empty(t, st.config)
		})
	}
}

func TestGlobalLoadFile(t *testing.T) {
	v = New()
	err := LoadFile("testdata/config.json")
	assertEqualErrors(t, nil, err)

	assert.Equal(t, v.Get("foo"), "bar")
	assert.Equal(t, v.Get("level"), 5.0)
}

func TestGlobalLoadDirectory(t *testing.T) {
	v = New()
	err := LoadDirectory("testdata/sub", false)
	assertEqualErrors(t, nil, err)

	assert.Equal(t, v.Get("foo"), "bar")
	assert.Equal(t, v.Get("level"), 5.0)
}

func TestGlobalDebug(t *testing.T) {
	testIO := []struct {
		tc      string
		configs ConfigMap
		expect  string
	}{
		{
			tc: "should debug",
			configs: ConfigMap{
				"foo": "bar",
				"baz": map[string]interface{}{
					"bar": "foo",
				},
			},
			expect: `{
  "20": {
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
			v = New()
			st := v.Store.(*DefaultConfigStore)
			st.config[EnvironmentLevel] = test.configs

			actual := Debug()
			assert.Equal(t, test.expect, actual, actual)
		})
	}
}

func TestGlobalAlias(t *testing.T) {
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
			defer v.Clear()
			SetLevel(DefaultLevel, test.key, test.value)
			Alias(test.alias, test.key)

			unAliased := Get(test.key)
			assert.Equal(t, test.value, unAliased)

			aliased := Get(test.alias)
			assert.Equal(t, test.value, aliased)
		})
	}
}

func TestGlobalRegisterResolver(t *testing.T) {
	v.RegisterResolver(EnvironmentLevel, defaultEnvResolver)
	st := v.Store.(*DefaultConfigStore)
	assert.Contains(t, st.resolvers, EnvironmentLevel)
}
