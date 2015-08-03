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

/*
Package bitarray implements a bit array.  Useful for tracking bool type values in a space
efficient way.  This is *NOT* a threadsafe package.
*/
package bitarray

import (
	"encoding/binary"
	"io"
)

// bitArray is a struct that maintains state of a bit array.
type bitArray struct {
	blocks  []block
	lowest  uint64
	highest uint64
	anyset  bool
}

func getIndexAndRemainder(k uint64) (uint64, uint64) {
	return k / s, k % s
}

func (ba *bitArray) setLowest() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		if ba.blocks[i] == 0 {
			continue
		}

		pos := ba.blocks[i].findRightPosition()
		ba.lowest = (i * s) + pos
		ba.anyset = true
		return
	}

	ba.anyset = false
	ba.lowest = 0
	ba.highest = 0
}

func (ba *bitArray) setHighest() {
	for i := len(ba.blocks) - 1; i >= 0; i-- {
		if ba.blocks[i] == 0 {
			continue
		}

		pos := ba.blocks[i].findLeftPosition()
		ba.highest = (uint64(i) * s) + pos
		ba.anyset = true
		return
	}

	ba.anyset = false
	ba.highest = 0
	ba.lowest = 0
}

// capacity returns the total capacity of the bit array.
func (ba *bitArray) Capacity() uint64 {
	return uint64(len(ba.blocks)) * s
}

// ToNums converts this bitarray to a list of numbers contained within it.
func (ba *bitArray) ToNums() []uint64 {
	nums := make([]uint64, 0, ba.highest-ba.lowest/4)
	for i, block := range ba.blocks {
		block.toNums(uint64(i)*s, &nums)
	}

	return nums
}

// SetBit sets a bit at the given index to true.
func (ba *bitArray) SetBit(k uint64) error {
	if k >= ba.Capacity() {
		return OutOfRangeError(k)
	}

	if !ba.anyset {
		ba.lowest = k
		ba.highest = k
		ba.anyset = true
	} else {
		if k < ba.lowest {
			ba.lowest = k
		} else if k > ba.highest {
			ba.highest = k
		}
	}

	i, pos := getIndexAndRemainder(k)
	ba.blocks[i] = ba.blocks[i].insert(pos)
	return nil
}

// GetBit returns a bool indicating if the value at the given
// index has been set.
func (ba *bitArray) GetBit(k uint64) (bool, error) {
	if k >= ba.Capacity() {
		return false, OutOfRangeError(k)
	}

	i, pos := getIndexAndRemainder(k)
	result := ba.blocks[i]&block(1<<pos) != 0
	return result, nil
}

//ClearBit will unset a bit at the given index if it is set.
func (ba *bitArray) ClearBit(k uint64) error {
	if k >= ba.Capacity() {
		return OutOfRangeError(k)
	}

	if !ba.anyset { // nothing is set, might as well bail
		return nil
	}

	i, pos := getIndexAndRemainder(k)
	ba.blocks[i] &^= block(1 << pos)

	if k == ba.highest {
		ba.setHighest()
	} else if k == ba.lowest {
		ba.setLowest()
	}
	return nil
}

// Or will bitwise or two bit arrays and return a new bit array
// representing the result.
func (ba *bitArray) Or(other BitArray) BitArray {
	if dba, ok := other.(*bitArray); ok {
		return orDenseWithDenseBitArray(ba, dba)
	}

	return orSparseWithDenseBitArray(other.(*sparseBitArray), ba)
}

// And will bitwise and two bit arrays and return a new bit array
// representing the result.
func (ba *bitArray) And(other BitArray) BitArray {
	if dba, ok := other.(*bitArray); ok {
		return andDenseWithDenseBitArray(ba, dba)
	}

	return andSparseWithDenseBitArray(other.(*sparseBitArray), ba)
}

// Reset clears out the bit array.
func (ba *bitArray) Reset() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		ba.blocks[i] &= block(0)
	}
	ba.anyset = false
}

// Equals returns a bool indicating if these two bit arrays are equal.
func (ba *bitArray) Equals(other BitArray) bool {
	if other.Capacity() == 0 && ba.highest > 0 {
		return false
	}

	if other.Capacity() == 0 && !ba.anyset {
		return true
	}

	var selfIndex uint64
	for iter := other.Blocks(); iter.Next(); {
		toIndex, otherBlock := iter.Value()
		if toIndex > selfIndex {
			for i := selfIndex; i < toIndex; i++ {
				if ba.blocks[i] > 0 {
					return false
				}
			}
		}

		selfIndex = toIndex
		if !ba.blocks[selfIndex].equals(otherBlock) {
			return false
		}
		selfIndex++
	}

	lastIndex, _ := getIndexAndRemainder(ba.highest)
	if lastIndex >= selfIndex {
		return false
	}

	return true
}

// Intersects returns a bool indicating if the supplied bitarray intersects
// this bitarray.  This will check for intersection up to the length of the supplied
// bitarray.  If the supplied bitarray is longer than this bitarray, this
// function returns false.
func (ba *bitArray) Intersects(other BitArray) bool {
	if other.Capacity() > ba.Capacity() {
		return false
	}

	if sba, ok := other.(*sparseBitArray); ok {
		return ba.intersectsSparseBitArray(sba)
	}

	return ba.intersectsDenseBitArray(other.(*bitArray))
}

// Blocks will return an iterator over this bit array.
func (ba *bitArray) Blocks() Iterator {
	return newBitArrayIterator(ba)
}

// complement flips all bits in this array.
func (ba *bitArray) complement() {
	for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
		ba.blocks[i] = ^ba.blocks[i]
	}

	ba.setLowest()
	if ba.anyset {
		ba.setHighest()
	}
}

func (ba *bitArray) intersectsSparseBitArray(other *sparseBitArray) bool {
	for i, index := range other.indices {
		if !ba.blocks[index].intersects(other.blocks[i]) {
			return false
		}
	}

	return true
}

func (ba *bitArray) intersectsDenseBitArray(other *bitArray) bool {
	for i, block := range other.blocks {
		if !ba.blocks[i].intersects(block) {
			return false
		}
	}

	return true
}

func (ba *bitArray) copy() BitArray {
	blocks := make(blocks, len(ba.blocks))
	copy(blocks, ba.blocks)
	return &bitArray{
		blocks:  blocks,
		lowest:  ba.lowest,
		highest: ba.highest,
		anyset:  ba.anyset,
	}
}

// Write serializes the bitArray and its data and sends it to the writer.
func Write(w io.Writer, ba *bitArray) error {
	err := binary.Write(w, binary.LittleEndian, ba.lowest)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.LittleEndian, ba.highest)
	if err != nil {
		return err
	}

	var encodedanyset uint8
	if ba.anyset {
		encodedanyset = 1
	} else {
		encodedanyset = 0
	}
	err = binary.Write(w, binary.LittleEndian, encodedanyset)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, ba.blocks)
	return err
}

// Read takes a reader of a serialized bitArray created by the Write function,
// and returns a bitArray object.
func Read(r io.Reader) (*bitArray, error) {
	ret := &bitArray{}

	err := binary.Read(r, binary.LittleEndian, &ret.lowest)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.LittleEndian, &ret.highest)
	if err != nil {
		return nil, err
	}

	var encodedanyset uint8
	err = binary.Read(r, binary.LittleEndian, &encodedanyset)
	if err != nil {
		return nil, err
	}

	// anyset defaults to false so we don't need an else statement
	if encodedanyset == 1 {
		ret.anyset = true
	}

	var nextblock block
	err = binary.Read(r, binary.LittleEndian, &nextblock)
	for err == nil {
		ret.blocks = append(ret.blocks, nextblock)
		err = binary.Read(r, binary.LittleEndian, &nextblock)
	}
	if err != io.EOF {
		return nil, err
	}
	return ret, nil
}

// newBitArray returns a new dense BitArray at the specified size. This is a
// separate private constructor so unit tests don't have to constantly cast the
// BitArray interface to the concrete type.
func newBitArray(size uint64, args ...bool) *bitArray {
	i, r := getIndexAndRemainder(size)
	if r > 0 {
		i++
	}

	ba := &bitArray{
		blocks: make([]block, i),
		anyset: false,
	}

	if len(args) > 0 && args[0] == true {
		for i := uint64(0); i < uint64(len(ba.blocks)); i++ {
			ba.blocks[i] = maximumBlock
		}

		ba.lowest = 0
		ba.highest = i*s - 1
		ba.anyset = true
	}

	return ba
}

// NewBitArray returns a new BitArray at the specified size.  The
// optional arg denotes whether this bitarray should be set to the
// bitwise complement of the empty array, ie. sets all bits.
func NewBitArray(size uint64, args ...bool) BitArray {
	return newBitArray(size, args...)
}
