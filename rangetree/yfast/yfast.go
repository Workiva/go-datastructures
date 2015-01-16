package yfast

import (
	"log"

	"github.com/Workiva/go-datastructures/trie/yfast"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

// isLastDimension returns a bool indicating if the provided
// value is the last dimension.  Because the user provides
// a 1-index number, this returns dimension == lastDimension-1.
func isLastDimension(dimension, lastDimension uint64) bool {
	return dimension == lastDimension-1
}

type dimensionalWrapper struct {
	key   uint64
	entry Entry
	trie  *yfast.YFastTrie
}

// Key returns the key of the wrapper which is the value
// of the entry in this wrapper's dimension.
func (dw *dimensionalWrapper) Key() uint64 {
	return dw.key
}

func newWrapper(key uint64, entry Entry, trie *yfast.YFastTrie) *dimensionalWrapper {
	return &dimensionalWrapper{
		key:   key,
		trie:  trie,
		entry: entry,
	}
}

type rangeTree struct {
	top        *yfast.YFastTrie // the first dimension, required
	dimensions uint64           // number of dimensions in this RT
	bitsize    interface{}      // will be uint of some sort
}

func (rt *rangeTree) add(entry Entry) Entry {
	trie := rt.top
	var wrapper *dimensionalWrapper
	var key uint64 // so we don't make too many calls to entry.ValueAtDimension
	var overwritten Entry
	var yfastEntry yfast.Entry
	for i := uint64(0); i < rt.dimensions; i++ {
		key = entry.ValueAtDimension(i)
		yfastEntry = trie.Get(key)
		if yfastEntry == nil {
			if isLastDimension(i, rt.dimensions) {
				wrapper = newWrapper(key, entry, nil)
			} else {
				wrapper = newWrapper(key, nil, yfast.New(rt.bitsize))
			}
			trie.Insert(wrapper)
			trie = wrapper.trie
			continue
		}

		wrapper = yfastEntry.(*dimensionalWrapper)
		trie = wrapper.trie
		if isLastDimension(i, rt.dimensions) {
			overwritten = wrapper.entry
			wrapper.entry = entry
		}
	}

	return overwritten
}

func (rt *rangeTree) Add(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, rt.add(e))
	}

	return overwritten
}

func (rt *rangeTree) get(entry Entry) Entry {
	var yfastEntry yfast.Entry
	trie := rt.top
	println(`STARTING LOOP`)
	for i := uint64(0); i < rt.dimensions; i++ {
		yfastEntry = trie.Get(entry.ValueAtDimension(i))
		if yfastEntry == nil {
			return nil
		}
		trie = yfastEntry.(*dimensionalWrapper).trie
	}

	println(`LOOP DONE`)
	return yfastEntry.(*dimensionalWrapper).entry
}

func (rt *rangeTree) Get(entries ...Entry) Entries {
	result := make(Entries, 0, len(entries))
	for _, e := range entries {
		result = append(result, rt.get(e))
	}

	return result
}

func (rt *rangeTree) delete(entry Entry) Entry {
	return nil
}

func new(dimensions uint64, bitsize interface{}) *rangeTree {
	return &rangeTree{
		bitsize:    bitsize,
		dimensions: dimensions,
		top:        yfast.New(bitsize),
	}
}

func New(dimensions uint64, bitsize interface{}) RangeTree {
	return new(dimensions, bitsize)
}
