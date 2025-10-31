package heap

// This file contains fuzz/property tests. We keep them separate from deterministic
// unit tests and benchmarks to:
//  (1) make it explicit they run under the fuzzing engine (go test -run Fuzz -fuzz=...)
//  (2) avoid mixing fuzz-specific helpers and seeds with regular unit tests
//  (3) simplify CI configuration where fuzzing may be opt-in or longer-running

import (
	"math/rand"
	"testing"
)

// helper to convert bytes to ints with negatives
func bytesToInts(data []byte) []int {
	res := make([]int, len(data))
	for i, b := range data {
		res[i] = int(int8(b))
	}
	return res
}

// FuzzHeapProperties validates ordering for binary heap across random inputs.
func FuzzHeapProperties(f *testing.F) {
	seeds := [][]byte{
		{},
		{1},
		{5, 4, 3, 2, 1},
		{0, 0, 0},
		{251, 10, 254, 7, 7, 3}, // negative via int8
	}
	for _, s := range seeds {
		f.Add(s)
	}

	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		arr := bytesToInts(data)
		rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		h := NewHeap[int](cmp)
		for _, v := range arr {
			h.Push(v)
		}

		if len(arr) == 0 {
			if _, ok := h.Peek(); ok {
				t.Fatalf("expected empty heap to have no peek")
			}
		} else {
			min := arr[0]
			for _, v := range arr[1:] {
				if v < min {
					min = v
				}
			}
			if top, ok := h.Peek(); !ok || top != min {
				t.Fatalf("peek mismatch: got %v %v, want %v true", top, ok, min)
			}
		}

		prevSet := false
		var prev int
		for {
			v, ok := h.Pop()
			if !ok {
				break
			}
			if prevSet && v < prev {
				t.Fatalf("heap order violated: %v < %v", v, prev)
			}
			prev = v
			prevSet = true
		}
	})
}

// FuzzDaryHeapProperties validates ordering for several d across random inputs.
func FuzzDaryHeapProperties(f *testing.F) {
	seeds := [][]byte{
		{}, {1}, {2, 1}, {3, 1, 2}, {255, 255, 0, 5},
	}
	for _, s := range seeds {
		f.Add(s)
	}

	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	dVals := []int{2, 3, 4, 6, 8}
	f.Fuzz(func(t *testing.T, data []byte) {
		arr := bytesToInts(data)
		rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		for _, d := range dVals {
			h := NewDaryHeapFromSlice[int](d, arr, cmp)
			prevSet := false
			var prev int
			for {
				v, ok := h.Pop()
				if !ok {
					break
				}
				if prevSet && v < prev {
					t.Fatalf("d-ary(%d) heap order violated: %v < %v", d, v, prev)
				}
				prev = v
				prevSet = true
			}
		}
	})
}
