package fibheap_test

import (
	"math"
	"testing"

	"github.com/Workiva/go-datastructures/fibheap"
)

func TestDeleteRemovesRequestedEntryWithInfinitePriorities(t *testing.T) {
	for _, tc := range []struct {
		name       string
		priorities []float64
		target     int
	}{
		{"finite target beside negative infinity", []float64{math.Inf(-1), 0, 1}, 1},
		{"positive infinity target", []float64{math.Inf(-1), 1, math.Inf(1)}, 2},
		{"one of two negative infinities", []float64{math.Inf(-1), math.Inf(-1), 1}, 0},
		{"smallest finite target", []float64{math.Inf(-1), -math.MaxFloat64, math.MaxFloat64}, 1},
		{"ordinary finite target", []float64{-5, 0, 1}, 1},
		{"single negative infinity", []float64{math.Inf(-1)}, 0},
		{"current minimum", []float64{-5, 0, 1}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := fibheap.NewFloatFibHeap()
			nodes := make([]*fibheap.Entry, len(tc.priorities))
			for i, p := range tc.priorities {
				nodes[i] = h.Enqueue(p)
			}
			if err := h.Delete(nodes[tc.target]); err != nil {
				t.Fatal(err)
			}
			if h.Size() != uint(len(nodes)-1) {
				t.Fatal("Delete changed the wrong number of entries")
			}
			seen := map[*fibheap.Entry]bool{}
			previous := math.Inf(-1)
			for !h.IsEmpty() {
				n, err := h.DequeueMin()
				if err != nil {
					t.Fatal(err)
				}
				if n == nodes[tc.target] {
					t.Error("deleted entry remains in heap")
				}
				if seen[n] {
					t.Error("entry returned twice")
				}
				seen[n] = true
				if n.Priority < previous {
					t.Error("remaining priorities not ordered")
				}
				previous = n.Priority
			}
			for i, n := range nodes {
				if i != tc.target && !seen[n] {
					t.Errorf("unrequested entry%d was removed", i)
				}
			}
		})
	}
}
