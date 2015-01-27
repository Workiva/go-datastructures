package skiplist

import (
	"github.com/Workiva/go-datastructures/rangetree"
	"github.com/Workiva/go-datastructures/slice/skip"
)

func isLastDimension(dimension, lastDimension uint64) bool {
	if dimension >= lastDimension { // useful in testing and denotes a serious problem
		panic(`Dimension is greater than possible dimensions.`)
	}

	return dimension == lastDimension-1
}

type dimensionalBundle struct {
	key uint64
	sl  *skip.SkipList
}

func (db *dimensionalBundle) Key() uint64 {
	return db.key
}

type lastBundle struct {
	key   uint64
	entry rangetree.Entry
}

func (lb *lastBundle) Key() uint64 {
	return lb.key
}

type skipListRT struct {
	top                *skip.SkipList
	dimensions, number uint64
}

func (rt *skipListRT) init(dimensions uint64) {
	rt.dimensions = dimensions
	rt.top = skip.New(uint64(0))
}

func (rt *skipListRT) add(entry rangetree.Entry) rangetree.Entry {
	var (
		value int64
		e     skip.Entry
		sl    = rt.top
		db    *dimensionalBundle
		lb    *lastBundle
	)

	for i := uint64(0); i < rt.dimensions; i++ {
		value = entry.ValueAtDimension(i)
		e = sl.Get(uint64(value))[0]
		if isLastDimension(i, rt.dimensions) {
			if e != nil { // this is an overwrite
				lb = e.(*lastBundle)
				oldEntry := lb.entry
				lb.entry = entry
				return oldEntry
			}

			// need to add new sl entry
			lb = &lastBundle{key: uint64(value), entry: entry}
			rt.number++
			sl.Insert(lb)
			return nil
		}

		if e == nil { // we need the intermediate dimension
			db = &dimensionalBundle{key: uint64(value), sl: skip.New(uint64(0))}
			sl.Insert(db)
		} else {
			db = e.(*dimensionalBundle)
		}

		sl = db.sl
	}

	panic(`Ran out of dimensions before for loop completed.`)
}

func (rt *skipListRT) Add(entries ...rangetree.Entry) rangetree.Entries {
	overwritten := make(rangetree.Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, rt.add(e))
	}

	return overwritten
}

func (rt *skipListRT) get(entry rangetree.Entry) rangetree.Entry {
	var (
		sl    = rt.top
		e     skip.Entry
		value uint64
	)
	for i := uint64(0); i < rt.dimensions; i++ {
		value = uint64(entry.ValueAtDimension(i))
		e = sl.Get(value)[0]
		if e == nil {
			return nil
		}

		if isLastDimension(i, rt.dimensions) {
			return e.(*lastBundle).entry
		}

		sl = e.(*dimensionalBundle).sl
	}

	panic(`Reached past for loop without finding last dimension.`)
}

func (rt *skipListRT) Get(entries ...rangetree.Entry) rangetree.Entries {
	results := make(rangetree.Entries, 0, len(entries))
	for _, e := range entries {
		results = append(results, rt.get(e))
	}

	return results
}

func (rt *skipListRT) Len() uint64 {
	return rt.number
}

func new(dimensions uint64) *skipListRT {
	sl := &skipListRT{}
	sl.init(dimensions)
	return sl
}

func New() rangetree.RangeTree {
	return nil
}
