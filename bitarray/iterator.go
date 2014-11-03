package bitarray

type sparseBitArrayIterator struct {
	index int64
	sba   *sparseBitArray
}

// Next increments the index and returns a bool indicating
// if any further items exist.
func (iter *sparseBitArrayIterator) Next() bool {
	iter.index++
	return iter.index < int64(len(iter.sba.indices))
}

// Value returns the index and block at the given index.
func (iter *sparseBitArrayIterator) Value() (uint64, block) {
	return iter.sba.indices[iter.index], iter.sba.blocks[iter.index]
}

func newCompressedBitArrayIterator(sba *sparseBitArray) *sparseBitArrayIterator {
	return &sparseBitArrayIterator{
		sba:   sba,
		index: -1,
	}
}

type bitArrayIterator struct {
	index     int64
	stopIndex uint64
	ba        *bitArray
}

// Next increments the index and returns a bool indicating if any further
// items exist.
func (iter *bitArrayIterator) Next() bool {
	iter.index++
	return uint64(iter.index) <= iter.stopIndex
}

// Value returns an index and the block at this index.
func (iter *bitArrayIterator) Value() (uint64, block) {
	return uint64(iter.index), iter.ba.blocks[iter.index]
}

func newBitArrayIterator(ba *bitArray) *bitArrayIterator {
	stop, _ := getIndexAndRemainder(ba.highest)
	start, _ := getIndexAndRemainder(ba.lowest)
	return &bitArrayIterator{
		ba:        ba,
		index:     int64(start) - 1,
		stopIndex: stop,
	}
}
