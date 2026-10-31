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

func TestDenseNandUnequalCapacities(t *testing.T) {
	cases := []struct {
		name                string
		leftSize, rightSize uint64
		left, right, want   []uint64
	}{
		{"empty right storage", 128, 0, []uint64{0, 64}, nil, []uint64{0, 64}},
		{"preserve left tail", 128, 64, []uint64{0, 64}, []uint64{0}, []uint64{64}},
		{"remove only common bits", 512, 128, []uint64{0, 64, 511}, []uint64{1, 64}, []uint64{0, 511}},
		{"empty left values", 128, 64, nil, []uint64{0}, nil},
		{"equal capacity", 128, 128, []uint64{0, 64}, []uint64{0}, []uint64{64}},
		{"longer right", 64, 512, []uint64{0, 1}, []uint64{0, 511}, []uint64{1}},
		{"empty left storage", 0, 128, nil, []uint64{64}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("Nand panicked: %v", p)
				}
			}()
			left := bitarray.NewBitArray(tc.leftSize)
			right := bitarray.NewBitArray(tc.rightSize)
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
			got := left.Nand(right)
			if !reflect.DeepEqual(append([]uint64{}, got.ToNums()...), append([]uint64{}, tc.want...)) {
				t.Fatalf("Nand = %v, want %v", got.ToNums(), tc.want)
			}
			if got.Capacity() != left.Capacity() {
				t.Fatalf("capacity changed: %d != %d", got.Capacity(), left.Capacity())
			}
			if !reflect.DeepEqual(append([]uint64{}, left.ToNums()...), append([]uint64{}, tc.left...)) {
				t.Fatal("left input changed")
			}
		})
	}
}
