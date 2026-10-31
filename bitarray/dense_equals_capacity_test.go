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

func TestDenseEqualsUnequalCapacities(t *testing.T) {
	cases := []struct {
		name                         string
		leftSize, rightSize          uint64
		left, right                  []uint64
		rightSparse, resetLeft, want bool
	}{
		{name: "zero storage and empty dense", want: true},
		{name: "zero storage and populated dense", rightSize: 128, right: []uint64{64}},
		{name: "short empty receiver", leftSize: 64, rightSize: 512, right: []uint64{129}},
		{name: "extra dense blocks", leftSize: 64, rightSize: 512, left: []uint64{0}, right: []uint64{0, 129}},
		{name: "extra sparse blocks", leftSize: 64, left: []uint64{0}, right: []uint64{129}, rightSparse: true},
		{name: "zero storage and sparse", right: []uint64{129}, rightSparse: true},
		{name: "reset left and empty sparse", leftSize: 512, left: []uint64{129}, rightSparse: true, resetLeft: true, want: true},
		{name: "reset left and zero dense", leftSize: 512, left: []uint64{129}, resetLeft: true, want: true},
		{name: "same values different capacities", leftSize: 64, rightSize: 512, left: []uint64{0}, right: []uint64{0}, want: true},
		{name: "empty dense and sparse", leftSize: 64, rightSparse: true, want: true},
		{name: "different values", leftSize: 64, rightSize: 64, left: []uint64{0}, right: []uint64{1}},
		{name: "unmatched left blocks", leftSize: 512, rightSize: 64, left: []uint64{0, 129}, right: []uint64{0}},
		{name: "sparse gap equality", leftSize: 512, left: []uint64{129}, right: []uint64{129}, rightSparse: true, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("Equals panicked: %v", p)
				}
			}()
			left := bitarray.NewBitArray(tc.leftSize)
			right := bitarray.NewBitArray(tc.rightSize)
			if tc.rightSparse {
				right = bitarray.NewSparseBitArray()
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
			if tc.resetLeft {
				left.Reset()
			}
			if got := left.Equals(right); got != tc.want {
				t.Fatalf("Equals = %t, want %t", got, tc.want)
			}
		})
	}
}
