package yfast

// Entry defines items that can be inserted into the y-fast
// trie.
type Entry interface {
	// Key is the key for this entry.  If the trie has been
	// given bit size n, only the last n bits of this key
	// will matter.  Use a bit size of 64 to enable all
	// 2^64-1 keys.
	Key() uint64
}
