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

func nandSparseWithSparseBitArray(sba, other *sparseBitArray) BitArray {
	// nand is an operation on the incoming array only, so the size will never
	// be more than the incoming array, regardless of the size of the other
	max := len(sba.indices)
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	var resultBlock block

	// move through the array and compare the blocks if they happen to
	// intersect
	for {
		if selfIndex == len(sba.indices) {
			// The bitarray being operated on is exhausted, so just return
			break
		} else if otherIndex == len(other.indices) {
			// The other array is exhausted. In this case, we assume that we
			// are calling nand on empty bit arrays, which is the same as just
			// copying the value in the sba array
			indices = append(indices, sba.indices[selfIndex])
			blocks = append(blocks, sba.blocks[selfIndex])
			selfIndex++
			continue
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
			// Here, the sba array has blocks that the other array doesn't
			// have. In this case, we just copy exactly the sba array values
			indices = append(indices, selfValue)
			blocks = append(blocks, sba.blocks[selfIndex])

			// This is the exact logical inverse of the above case.
			selfIndex++

		default:
			// Here, our indices match for both `sba` and `other`.
			// Time to do the bitwise AND operation and add a block
			// to our result list if the block has values in it.
			resultBlock = sba.blocks[selfIndex].nand(other.blocks[otherIndex])
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

func nandSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	// Since nand is non-commutative, the resulting array should be sparse,
	// and the same length or less than the sparse array
	indices := make(uintSlice, 0, len(sba.indices))
	blocks := make(blocks, 0, len(sba.indices))

	var resultBlock block

	// Loop through the sparse array and match it with the dense array.
	for selfIndex, selfValue := range sba.indices {
		if selfValue >= uint64(len(other.blocks)) {
			// Since the dense array is exhausted, just copy over the data
			// from the sparse array
			resultBlock = sba.blocks[selfIndex]
			indices = append(indices, selfValue)
			blocks = append(blocks, resultBlock)
			continue
		}

		resultBlock = sba.blocks[selfIndex].nand(other.blocks[selfValue])
		if resultBlock > 0 {
			indices = append(indices, selfValue)
			blocks = append(blocks, resultBlock)
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func nandDenseWithSparseBitArray(sba *bitArray, other *sparseBitArray) BitArray {
	// Since nand is non-commutative, the resulting array should be dense,
	// and the same length or less than the dense array
	tmp := sba.copy()
	ret := tmp.(*bitArray)

	// Loop through the other array and match it with the sba array.
	for otherIndex, otherValue := range other.indices {
		if otherValue >= uint64(len(ret.blocks)) {
			break
		}

		ret.blocks[otherValue] = sba.blocks[otherValue].nand(other.blocks[otherIndex])
	}

	ret.setLowest()
	ret.setHighest()

	return ret
}

func nandDenseWithDenseBitArray(dba, other *bitArray) BitArray {
	min := uint64(len(dba.blocks))

	ba := newBitArray(min * s)

	for i := uint64(0); i < min; i++ {
		ba.blocks[i] = dba.blocks[i].nand(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}
