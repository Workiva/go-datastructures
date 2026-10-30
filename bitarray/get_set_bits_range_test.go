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
	"fmt"
	"math"
	"testing"

	"github.com/Workiva/go-datastructures/bitarray"
)

func TestDenseGetSetBitsPastCapacity(t *testing.T) {
	for _, size := range []uint64{0, 1, 64, 65, 129} {
		for _, from := range []uint64{size + 128, math.MaxUint64} {
			t.Run(fmt.Sprintf("size_%d/from_%d", size, from), func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Errorf("GetSetBits panicked: %v", p)
					}
				}()
				array := bitarray.NewBitArray(size)
				if size > 0 {
					if err := array.SetBit(0); err != nil {
						t.Fatal(err)
					}
				}
				for _, buffer := range [][]uint64{nil, {}, make([]uint64, 0, 2), {99, 100}} {
					result := array.GetSetBits(from, buffer)
					if len(result) != 0 {
						t.Fatalf("GetSetBits(%d) = %v, want no bits", from, result)
					}
					if cap(result) != cap(buffer) {
						t.Fatalf("buffer capacity changed: %d != %d", cap(result), cap(buffer))
					}
				}
			})
		}
	}
}

func TestDenseGetSetBitsBoundaryAndBuffer(t *testing.T) {
	array := bitarray.NewBitArray(65)
	for _, k := range []uint64{0, 63, 64, 127} {
		if err := array.SetBit(k); err != nil {
			t.Fatal(err)
		}
	}
	buffer := []uint64{99, 100}
	got := array.GetSetBits(64, buffer[:0])
	if len(got) != 2 || got[0] != 64 || got[1] != 127 {
		t.Fatalf("in-range results = %v", got)
	}
	got = array.GetSetBits(array.Capacity(), buffer)
	if len(got) != 0 || cap(got) != 2 {
		t.Fatalf("boundary result = %v (capacity %d)", got, cap(got))
	}
	if buffer[0] != 64 || buffer[1] != 127 {
		t.Fatalf("out-of-range call changed buffer: %v", buffer)
	}
}
