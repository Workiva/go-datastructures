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

func orSparseWithSparseBitArray(sba *sparseBitArray,
	other *sparseBitArray) BitArray {

	if len(other.indices) == 0 {
		return sba.copy()
	}

	if len(sba.indices) == 0 {
		return other.copy()
	}

	max := maxInt64(int64(len(sba.indices)), int64(len(other.indices)))
	indices := make(uintSlice, 0, max)
	blocks := make(blocks, 0, max)

	selfIndex := 0
	otherIndex := 0
	for {
		// last comparison was a real or, we are both exhausted now
		if selfIndex == len(sba.indices) && otherIndex == len(other.indices) {
			break
		} else if selfIndex == len(sba.indices) {
			indices = append(indices, other.indices[otherIndex:]...)
			blocks = append(blocks, other.blocks[otherIndex:]...)
			break
		} else if otherIndex == len(other.indices) {
			indices = append(indices, sba.indices[selfIndex:]...)
			blocks = append(blocks, sba.blocks[selfIndex:]...)
			break
		}

		selfValue := sba.indices[selfIndex]
		otherValue := other.indices[otherIndex]

		switch diff := int(otherValue) - int(selfValue); {
		case diff > 0:
			indices = append(indices, selfValue)
			blocks = append(blocks, sba.blocks[selfIndex])
			selfIndex++
		case diff < 0:
			indices = append(indices, otherValue)
			blocks = append(blocks, other.blocks[otherIndex])
			otherIndex++
		default:
			indices = append(indices, otherValue)
			blocks = append(blocks, sba.blocks[selfIndex].or(other.blocks[otherIndex]))
			selfIndex++
			otherIndex++
		}
	}

	return &sparseBitArray{
		indices: indices,
		blocks:  blocks,
	}
}

func orSparseWithDenseBitArray(sba *sparseBitArray, other *bitArray) BitArray {
	if other.Capacity() == 0 || !other.anyset {
		return sba.copy()
	}

	if sba.Capacity() == 0 {
		return other.copy()
	}

	max := maxUint64(uint64(sba.Capacity()), uint64(other.Capacity()))

	ba := newBitArray(max * s)
	selfIndex := 0
	otherIndex := 0
	for {
		if selfIndex == len(sba.indices) && otherIndex == len(other.blocks) {
			break
		} else if selfIndex == len(sba.indices) {
			copy(ba.blocks[otherIndex:], other.blocks[otherIndex:])
			break
		} else if otherIndex == len(other.blocks) {
			for i, value := range sba.indices[selfIndex:] {
				ba.blocks[value] = sba.blocks[i+selfIndex]
			}
			break
		}

		selfValue := sba.indices[selfIndex]
		if selfValue == uint64(otherIndex) {
			ba.blocks[otherIndex] = sba.blocks[selfIndex].or(other.blocks[otherIndex])
			selfIndex++
			otherIndex++
			continue
		}

		ba.blocks[otherIndex] = other.blocks[otherIndex]
		otherIndex++
	}

	ba.setHighest()
	ba.setLowest()

	return ba
}

func orDenseWithDenseBitArray(dba *bitArray, other *bitArray) BitArray {
	if dba.Capacity() == 0 || !dba.anyset {
		return other.copy()
	}

	if other.Capacity() == 0 || !other.anyset {
		return dba.copy()
	}

	max := maxUint64(uint64(len(dba.blocks)), uint64(len(other.blocks)))

	ba := newBitArray(max * s)

	for i := uint64(0); i < max; i++ {
		if i == uint64(len(dba.blocks)) {
			copy(ba.blocks[i:], other.blocks[i:])
			break
		}

		if i == uint64(len(other.blocks)) {
			copy(ba.blocks[i:], dba.blocks[i:])
			break
		}

		ba.blocks[i] = dba.blocks[i].or(other.blocks[i])
	}

	ba.setLowest()
	ba.setHighest()

	return ba
}
