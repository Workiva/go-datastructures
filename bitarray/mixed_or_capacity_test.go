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
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/bitarray"
)

func TestMixedOrCapacityUnits(t *testing.T) {
	cases := []struct {
		name                    string
		left, right, want       []uint64
		denseSize, wantCapacity uint64
	}{
		{"one block", []uint64{0}, []uint64{1}, []uint64{0, 1}, 64, 64},
		{"sparse extends dense", []uint64{129}, []uint64{0}, []uint64{0, 129}, 64, 192},
		{"dense extends sparse", []uint64{0}, []uint64{511}, []uint64{0, 511}, 512, 512},
		{"overlap", []uint64{0, 64}, []uint64{64}, []uint64{0, 64}, 128, 128},
		{"empty sparse", nil, []uint64{64}, []uint64{64}, 128, 128},
		{"empty dense", []uint64{129}, nil, []uint64{129}, 64, 192},
	}
	for _, tc := range cases {
		for _, reverse := range []bool{false, true} {
			name := tc.name + "/sparse-first"
			if reverse {
				name = tc.name + "/dense-first"
			}
			t.Run(name, func(t *testing.T) {
				left := bitarray.NewSparseBitArray()
				right := bitarray.NewBitArray(tc.denseSize)
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
				a, b := left, right
				if reverse {
					a, b = b, a
				}
				got := a.Or(b)
				if got.Capacity() != tc.wantCapacity {
					t.Fatalf("Or capacity=%d, want %d", got.Capacity(), tc.wantCapacity)
				}
				if !reflect.DeepEqual(got.ToNums(), tc.want) {
					t.Fatalf("Or values=%v, want %v", got.ToNums(), tc.want)
				}
			})
		}
	}
}
