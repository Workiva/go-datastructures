/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bitarray_test

import (
	"testing"

	"github.com/Workiva/go-datastructures/bitarray"
)

func TestSparseIntersectsSearchOffset(t *testing.T) {
	cases := []struct {
		name        string
		left, right []uint64
		dense       bool
		want        bool
	}{
		{"three matching blocks", []uint64{0, 64, 129}, []uint64{0, 64, 129}, false, true},
		{"skip receiver prefix", []uint64{0, 64, 129, 511}, []uint64{64, 129}, false, true},
		{"dense zero blocks", []uint64{0, 129}, []uint64{0, 129}, true, true},
		{"empty dense is contained", []uint64{129}, nil, true, true},
		{"missing other bit", []uint64{0, 64, 129}, []uint64{0, 65, 129}, false, false},
		{"other extends past receiver", []uint64{0, 64}, []uint64{511}, false, false},
		{"empty receiver", nil, []uint64{64}, false, false},
		{"empty sparse is contained", []uint64{129}, nil, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			left := bitarray.NewSparseBitArray()
			right := bitarray.NewSparseBitArray()
			if tc.dense {
				right = bitarray.NewBitArray(512)
			}
			for _, k := range tc.left {
				if err := left.SetBit(k); err != nil {
					t.Fatal(err)
				}
			}
			for _, k := range tc.right {
				if err := right.SetBit(k); err != nil {
					t.Fatal(err)
				}
			}
			if got := left.Intersects(right); got != tc.want {
				t.Fatalf("Intersects(%v,%v) = %t, want %t", tc.left, tc.right, got, tc.want)
			}
		})
	}
}
