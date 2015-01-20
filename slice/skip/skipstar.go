package skip

import "log"

func init() {
	log.Printf(`I HATE THIS.`)
}

type SkipStarList struct {
	ary uint8
	num uint64
	sl  *SkipList
}

type entryBundle struct {
	key     uint64
	entries Entries
}

func (eb *entryBundle) Key() uint64 {
	return eb.key
}

func newEntryBundle(key uint64, size uint8) *entryBundle {
	return &entryBundle{
		key:     key,
		entries: make(Entries, 0, size),
	}
}

func (ssl *SkipStarList) init(ifc interface{}) {
	switch ifc.(type) {
	case uint8:
		ssl.ary = 8
	case uint16:
		ssl.ary = 16
	case uint32:
		ssl.ary = 32
	case uint64, uint:
		ssl.ary = 64
	}
	ssl.sl = New(ifc)
}

func (ssl *SkipStarList) getNormalizedKey(key uint64) uint64 {
	key = key/uint64(ssl.ary) + 1
	return key * uint64(ssl.ary)
}

func (ssl *SkipStarList) insert(entry Entry) Entry {
	key := ssl.getNormalizedKey(entry.Key())
	eb, ok := ssl.sl.Get(entry.Key())[0].(*entryBundle)
	if !ok {
		eb = newEntryBundle(key, ssl.ary)
		ssl.sl.Insert(eb)
	}

	e := eb.entries.insert(entry)
	if e == nil {
		ssl.num++
	}
	return e
}

func (ssl *SkipStarList) Insert(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, ssl.insert(e))
	}

	return overwritten
}

func (ssl *SkipStarList) get(key uint64) Entry {
	normalizedKey := ssl.getNormalizedKey(key)
	eb, ok := ssl.sl.Get(normalizedKey)[0].(*entryBundle)
	if ok {
		return eb.entries.get(key)
	}
	return nil
}

func (ssl *SkipStarList) Get(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, ssl.get(key))
	}

	return entries
}

func (ssl *SkipStarList) delete(key uint64) Entry {
	normalizedKey := ssl.getNormalizedKey(key)
	eb, ok := ssl.sl.Get(normalizedKey)[0].(*entryBundle)
	if !ok {
		return nil
	}

	deleted := eb.entries.delete(key)
	if deleted != nil {
		ssl.num--
		if len(eb.entries) == 0 {
			ssl.sl.Delete(eb.key)
		}
	}

	return deleted
}

func (ssl *SkipStarList) Delete(keys ...uint64) Entries {
	deleted := make(Entries, 0, len(keys))
	for _, key := range keys {
		deleted = append(deleted, ssl.delete(key))
	}

	return deleted
}

func (ssl *SkipStarList) iter(key uint64) *starIterator {
	normalizedKey := ssl.getNormalizedKey(key)
	iter := ssl.sl.Iter(normalizedKey)
	if !iter.Next() {
		return &starIterator{
			index: iteratorExhausted,
		}
	}

	eb := iter.Value().(*entryBundle)
	return &starIterator{
		index:   eb.entries.search(key) - 1,
		entries: eb.entries,
		iter:    iter,
	}
}

func (ssl *SkipStarList) Iter(key uint64) Iterator {
	return ssl.iter(key)
}

func (ssl *SkipStarList) Len() uint64 {
	return ssl.num
}

func NewStar(ifc interface{}) *SkipStarList {
	ssl := &SkipStarList{}
	ssl.init(ifc)
	return ssl
}
