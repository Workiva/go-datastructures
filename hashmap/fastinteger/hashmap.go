// Package fastinteger is designed to provide a very primitive
// implementation of a hash map for unsigned integer keys and
// values.  It is designed to have existence checks and insertions
// that are faster than Go's native implementation.  Like Go's
// native implementation, FastIntegerHashMap will dynamically
// grow in size.
//
// Current benchmarks on identical machine against native Go implementation:
// 		BenchmarkInsert-8	   10000	    131258 ns/op
//		BenchmarkGoMapInsert-8	   10000	    208787 ns/op
//		BenchmarkExists-8	  100000	     15820 ns/op
//		BenchmarkGoMapExists-8	  100000	     16394 ns/op
//		BenchmarkDelete-8	  100000	     17909 ns/op
//		BenchmarkGoDelete-8	   30000	     49376 ns/op
// 		BenchmarkInsertWithExpand-8	   20000	     90301 ns/op
//		BenchmarkGoInsertWithExpand-8	   10000	    142088 ns/op
//
//
// This performance could be further enhanced by using a
// better probing technique.
package fastinteger

const ratio = .75 // ratio sets the capacity the hashmap has to be at before it expands

// roundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type packet struct {
	key, value uint64
}

type packets []*packet

func (packets packets) find(key uint64) uint64 {
	h := hash(key)
	i := h & (uint64(len(packets)) - 1)
	for packets[i] != nil && packets[i].key != key {
		i = (i + 1) & (uint64(len(packets)) - 1)
	}

	return i
}

func (packets packets) set(packet *packet) {
	i := packets.find(packet.key)
	if packets[i] == nil {
		packets[i] = packet
		return
	}

	packets[i].value = packet.value
}

func (packets packets) get(key uint64) (uint64, bool) {
	i := packets.find(key)
	if packets[i] == nil {
		return 0, false
	}

	return packets[i].value, true
}

func (packets packets) delete(key uint64) bool {
	i := packets.find(key)
	if packets[i] == nil {
		return false
	}
	packets[i] = nil
	i = (i + 1) & (uint64(len(packets)) - 1)
	for packets[i] != nil {
		p := packets[i]
		packets[i] = nil
		packets.set(p)
		i = (i + 1) & (uint64(len(packets)) - 1)
	}
	return true
}

func (packets packets) exists(key uint64) bool {
	i := packets.find(key)
	return packets[i] != nil // technically, they can store nil
}

// FastIntegerHashMap is a simple hashmap to be used with
// integer only keys.  It supports few operations, and is designed
// primarily for cases where the consumer needs a very simple
// datastructure to set and check for existence of integer
// keys over a sparse range.
type FastIntegerHashMap struct {
	count   uint64
	packets packets
}

// rebuild is an expensive operation which requires us to iterate
// over the current bucket and rehash the keys for insertion into
// the new bucket.  The new bucket is twice as large as the old
// bucket by default.
func (fi *FastIntegerHashMap) rebuild() {
	packets := make(packets, roundUp(uint64(len(fi.packets))+1))
	for _, packet := range fi.packets {
		if packet == nil {
			continue
		}

		packets.set(packet)
	}
	fi.packets = packets
}

// Get returns an item from the map if it exists.  Otherwise,
// returns false for the second argument.
func (fi *FastIntegerHashMap) Get(key uint64) (uint64, bool) {
	return fi.packets.get(key)
}

// Set will set the provided key with the provided value.
func (fi *FastIntegerHashMap) Set(key, value uint64) {
	if float64(fi.count+1)/float64(len(fi.packets)) > ratio {
		fi.rebuild()
	}

	fi.packets.set(&packet{key: key, value: value})
	fi.count++
}

// Exists will return a bool indicating if the provided key
// exists in the map.
func (fi *FastIntegerHashMap) Exists(key uint64) bool {
	return fi.packets.exists(key)
}

// Delete will remove the provided key from the hashmap.  If
// the key cannot be found, this is a no-op.
func (fi *FastIntegerHashMap) Delete(key uint64) {
	if fi.packets.delete(key) {
		fi.count--
	}
}

// Len returns the number of items in the hashmap.
func (fi *FastIntegerHashMap) Len() uint64 {
	return fi.count
}

// Cap returns the capacity of the hashmap.
func (fi *FastIntegerHashMap) Cap() uint64 {
	return uint64(len(fi.packets))
}

// New returns a new FastIntegerHashMap with a bucket size specified
// by hint.
func New(hint uint64) *FastIntegerHashMap {
	if hint == 0 {
		hint = 16
	}

	hint = roundUp(hint)
	return &FastIntegerHashMap{
		count:   0,
		packets: make(packets, hint),
	}
}
