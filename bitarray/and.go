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

func andSparseWithSparseBitArray(sba, other *sparseBitArray) BitArray {
	max := maxInt64(int64(len(sba.indices)), int64(len(other.indices)))
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	var resultBlock block

	// move through the array and compare the blocks if they happen to
	// intersect
	for {
		if selfIndex == len(sba.indices) || otherIndex == len(other.indices) {
			// One of the arrays has been exhausted. We don't need
			// to compare anything else for a bitwise and; the
			// operation is complete.
			break
		}

		selfValue := sba.indices[selfIndex]
		otherValue := other.indices[otherIndex]

		switch {
		case otherValue < selfValue:
			// The `sba` bitarray has a block with a index position
			// greater than us. We want to compare with that block
			// if possible, so move our `other` index closer to that
			// block's index.
			otherIndex++

		case otherValue > selfValue:
			// This is the exact logical inverse of the above case.
			selfIndex++

		default:
			// Here, our indices match for both `sba` and `other`.
			// Time to do the bitwise AND operation and add a block
			// to our result list if the block has values in it.
			resultBlock = sba.blocks[selfIndex].and(other.blocks[otherIndex])
			if resultBlock > 0 {
				indices = append(indices, selfValue)
				blocks = append(blocks, resultBlock)
			}
			selfIndex++
			otherIndex++
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func andSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	// Use a duplicate of the sparse array to store the results of the
	// bitwise and. More memory-efficient than allocating a new dense bit
	// array.
	//
	// NOTE: this could be faster if we didn't copy the values as well
	// (since they are overwritten), but I don't want this method to know
	// too much about the internals of sparseBitArray. The performance hit
	// should be minor anyway.
	ba := sba.copy()

	// Run through the sparse array and attempt comparisons wherever
	// possible against the dense bit array.
	for selfIndex, selfValue := range ba.indices {

		if selfValue >= uint64(len(other.blocks)) {
			// The dense bit array has been exhausted. This is the
			// annoying case because we have to trim the sparse
			// array to the size of the dense array.
			ba.blocks = ba.blocks[:selfIndex]
			ba.indices = ba.indices[:selfIndex]

			// once this is done, there are no more comparisons.
			// We're ready to return
			break
		}
		ba.blocks[selfIndex] = ba.blocks[selfIndex].and(
			other.blocks[selfValue])
	}

	return ba
}

func andDenseWithDenseBitArray(dba, other *bitArray) BitArray {
	min := minUint64(uint64(len(dba.blocks)), uint64(len(other.blocks)))

	ba := newBitArray(min * s)

	for i := uint64(0); i < min; i++ {
		ba.blocks[i] = dba.blocks[i].and(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}
