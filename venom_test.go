package venom

import (
	"fmt"
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
				lkv{DefaultLevel, "foo", "bar"},
			},
			expected: []kv{
				kv{"foo", "bar"},
			},
		},
		// slightly more complex test case with use of both the default and
		// override config levels, asserting that the higher priority wins
		{
			inp: []lkv{
				lkv{DefaultLevel, "foo", "bar"},
				lkv{OverrideLevel, "foo", "baz"},
			},
			expected: []kv{
				kv{"foo", "baz"},
			},
		},
		// simple test case using nested keys
		{
			inp: []lkv{
				lkv{DefaultLevel, "foo.bar", "baz"},
			},
			expected: []kv{
				kv{"foo.bar", "baz"},
				kv{"foo", configMap{"bar": "baz"}},
			},
		},
		// more complex test case using multiple config levels and nested key space
		{
			inp: []lkv{
				lkv{DefaultLevel, "foo.bar", "baz"},
				lkv{OverrideLevel, "foo.bar", 12},
			},
			expected: []kv{
				kv{"foo.bar", 12},
				kv{"foo", configMap{"bar": 12}},
			},
		},
		// simple case of an absent key
		{
			inp: []lkv{
				lkv{DefaultLevel, "foo", "bar"},
			},
			expected: []kv{
				kv{"bar", nil},
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
