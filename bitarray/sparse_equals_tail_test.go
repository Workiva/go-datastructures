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

func TestSparseEqualsRequiresAllReceiverBlocks(t *testing.T) {
	cases := []struct {
		name        string
		left, right []uint64
		dense       bool
		want        bool
	}{
		{name: "sparse prefix", left: []uint64{5, 129}, right: []uint64{5}},
		{name: "multiple trailing blocks", left: []uint64{5, 129, 256}, right: []uint64{5}},
		{name: "dense prefix", left: []uint64{5, 129}, right: []uint64{5}, dense: true},
		{name: "empty dense", left: []uint64{129}, dense: true},
		{name: "equal sparse", left: []uint64{5, 129}, right: []uint64{5, 129}, want: true},
		{name: "equal dense", left: []uint64{5, 129}, right: []uint64{5, 129}, dense: true, want: true},
		{name: "empty arrays", dense: true, want: true},
		{name: "other has trailing block", left: []uint64{5}, right: []uint64{5, 129}},
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
			if got := left.Equals(right); got != tc.want {
				t.Fatalf("Equals(%v, %v) = %t, want %t", tc.left, tc.right, got, tc.want)
			}
		})
	}
}
