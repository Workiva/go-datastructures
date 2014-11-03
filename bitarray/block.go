package bitarray

// s denotes the size of any element in the block array.  Cannot use
// unsafe.SizeOf here as you can't take the size of a type.
const s = uint64(64)

// block defines how we split apart the bit array.  This also
// determines the size of s.
type block uint64

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

func (b block) get(position uint64) bool {
	return b&block(1<<position) != 0
}

func (b block) equals(other block) bool {
	return b == other
}

func (b block) intersects(other block) bool {
	return b&other == other
}
