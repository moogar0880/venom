package venom

import "container/heap"

// NewConfigLevelHeap creates a pre-initialized ConfigLevelHeap instance
func NewConfigLevelHeap() *ConfigLevelHeap {
	h := new(ConfigLevelHeap)
	heap.Init(h)
	return h
}

// An ConfigLevelHeap is a max-heap of ConfigLevels
type ConfigLevelHeap []ConfigLevel

func (h ConfigLevelHeap) Len() int { return len(h) }
func (h ConfigLevelHeap) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater
	// than here.
	return h[i] > h[j]
}
func (h ConfigLevelHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push pushes a new value onto the heap, inserting in sort order
func (h *ConfigLevelHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(ConfigLevel))
}

// Pop removes the right-most (lowest) value from the heap
func (h *ConfigLevelHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
