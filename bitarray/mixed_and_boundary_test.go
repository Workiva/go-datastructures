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

func TestMixedAndDenseBoundary(t *testing.T) {
	cases := []struct {
		name              string
		left, right, want []uint64
		capacity          uint64
	}{
		{"first sparse block past dense", []uint64{64}, []uint64{0}, nil, 64},
		{"retain final common block", []uint64{0, 64}, []uint64{0}, []uint64{0}, 64},
		{"retain multiple common blocks", []uint64{0, 64, 128}, []uint64{0, 64}, []uint64{0, 64}, 128},
		{"remove zero result blocks", []uint64{0, 64}, []uint64{1}, nil, 64},
		{"empty dense", []uint64{0, 64}, nil, nil, 64},
		{"all blocks in range", []uint64{0, 64}, []uint64{0, 64}, []uint64{0, 64}, 128},
	}
	for _, tc := range cases {
		for _, reverse := range []bool{false, true} {
			name := tc.name + "/sparse-first"
			if reverse {
				name = tc.name + "/dense-first"
			}
			t.Run(name, func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("And panicked: %v", p)
					}
				}()
				left := bitarray.NewSparseBitArray()
				right := bitarray.NewBitArray(tc.capacity)
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
				got := a.And(b).ToNums()
				if !reflect.DeepEqual(append([]uint64{}, got...), append([]uint64{}, tc.want...)) {
					t.Fatalf("And = %v, want %v", got, tc.want)
				}
				if !reflect.DeepEqual(left.ToNums(), tc.left) {
					t.Fatal("left input changed")
				}
			})
		}
	}
}
