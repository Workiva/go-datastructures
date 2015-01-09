package fastinteger

// hash will convert the uint64 key into a hash based on Murmur3's 64-bit
// integer finalizer.
// Details here: https://code.google.com/p/smhasher/wiki/MurmurHash3
func hash(key uint64) uint64 {
	key ^= key >> 33
	key *= 0xff51afd7ed558ccd
	key ^= key >> 33
	key *= 0xc4ceb9fe1a85ec53
	key ^= key >> 33
	return key
}
