package venom

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigHeap(t *testing.T) {
	testIO := []struct {
		inp      []ConfigLevel
		expected *ConfigLevelHeap
	}{
		// empty slice
		{
			inp:      []ConfigLevel{},
			expected: new(ConfigLevelHeap),
		},
		// boring test with slice of size 1
		{
			inp:      []ConfigLevel{OverrideLevel},
			expected: &ConfigLevelHeap{OverrideLevel},
		},
		// already sorted list
		{
			inp:      []ConfigLevel{OverrideLevel, EnvironmentLevel, FileLevel},
			expected: &ConfigLevelHeap{OverrideLevel, EnvironmentLevel, FileLevel},
		},
		// unsorted list
		{
			inp:      []ConfigLevel{FileLevel, DefaultLevel, OverrideLevel, EnvironmentLevel},
			expected: &ConfigLevelHeap{OverrideLevel, EnvironmentLevel, FileLevel, DefaultLevel},
		},
	}

	for i, test := range testIO {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := NewConfigLevelHeap()
			for _, v := range test.inp {
				heap.Push(h, v)
			}
			assert.Equal(t, test.expected, h)

			slc := *test.expected
			assert.Equal(t, *test.expected, *h)
			for i := len(slc) - 1; i >= 0; i-- {
				if h.Len() == 0 {
					break
				}
				assert.Equal(t, slc[i], h.Pop())
			}
		})
	}
}
