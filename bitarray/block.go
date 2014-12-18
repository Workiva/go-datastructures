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

package bitarray

import (
	"fmt"
	"unsafe"
)

// block defines how we split apart the bit array. This also determines the size
// of s. This can be changed to any unsigned integer type: uint8, uint16,
// uint32, and so on.
type block uint64

// s denotes the size of any element in the block array.
// For a block of uint64, s will be equal to 64
// For a block of uint32, s will be equal to 32
// and so on...
const s = uint64(unsafe.Sizeof(block(0)) * 8)

// maximumBlock represents a block of all 1s and is used in the constructors.
const maximumBlock = block(0) | ^block(0)

func (b block) toNums(offset uint64, nums *[]uint64) {
	for i := uint64(0); i < s; i++ {
		if b&block(1<<i) > 0 {
			*nums = append(*nums, i+offset)
		}
	}
}

func (b block) findLeftPosition() uint64 {
	for i := s - 1; i < s; i-- {
		test := block(1 << i)
		if b&test == test {
			return i
		}
	}

	return s
}

func (b block) findRightPosition() uint64 {
	for i := uint64(0); i < s; i++ {
		test := block(1 << i)
		if b&test == test {
			return i
		}
	}

	return s
}

func (b block) insert(position uint64) block {
	return b | block(1<<position)
}

func (b block) remove(position uint64) block {
	return b & ^block(1<<position)
}

func (b block) or(other block) block {
	return b | other
}

func (b block) and(other block) block {
	return b & other
}

func (b block) get(position uint64) bool {
	return b&block(1<<position) != 0
}

func (b block) equals(other block) bool {
	return b == other
}

func (b block) intersects(other block) bool {
	return b&other == other
}

func (b block) String() string {
	return fmt.Sprintf(fmt.Sprintf("%%0%db", s), uint64(b))
}
