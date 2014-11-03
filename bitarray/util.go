package bitarray

// maxInt64 returns the highest integer in the provided list of int32s
func maxInt64(ints ...int64) int64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}

// maxUint64 returns the highest integer in the provided list of int32s
func maxUint64(ints ...uint64) uint64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}
