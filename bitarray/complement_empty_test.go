package bitarray_test

import (
	"testing"

	"github.com/Workiva/go-datastructures/bitarray"
)

func TestComplementedEmptyConstruction(t *testing.T) {
	cases := []struct {
		name string
		size uint64
		args []bool
	}{
		{name: "zero_complement", args: []bool{true}},
		{name: "zero_complement_first_argument", args: []bool{true, false}},
		{name: "zero_default"},
		{name: "zero_false", args: []bool{false}},
		{name: "zero_false_first_argument", args: []bool{false, true}},
		{name: "one_complement", size: 1, args: []bool{true}},
		{name: "one_block_complement", size: 64, args: []bool{true}},
		{name: "two_blocks_complement", size: 65, args: []bool{true}},
		{name: "nonzero_default", size: 65},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			array := bitarray.NewBitArray(tc.size, tc.args...)
			wantCount := 0
			if len(tc.args) > 0 && tc.args[0] {
				wantCount = int(array.Capacity())
			}
			if got := array.Count(); got != wantCount {
				t.Errorf("Count() = %d, want %d", got, wantCount)
			}
			if got := array.IsEmpty(); got != (wantCount == 0) {
				t.Errorf("IsEmpty() = %v for %d set bits", got, wantCount)
			}
			defer func() {
				if value := recover(); value != nil {
					t.Errorf("ToNums() panicked: %v", value)
				}
			}()
			nums := array.ToNums()
			if len(nums) != wantCount {
				t.Fatalf("ToNums() returned %d positions, want %d", len(nums), wantCount)
			}
			for i, position := range nums {
				if position != uint64(i) {
					t.Fatalf("ToNums()[%d] = %d, want %d", i, position, i)
				}
			}
		})
	}
}
