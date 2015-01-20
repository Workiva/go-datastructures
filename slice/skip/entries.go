package skip

import "sort"

func (entries Entries) search(key uint64) int {
	return sort.Search(len(entries), func(i int) bool {
		return entries[i].Key() >= key
	})
}

func (entries *Entries) insert(entry Entry) Entry {
	i := entries.search(entry.Key())
	if i >= len(*entries) {
		*entries = append(*entries, entry)
		return nil
	}

	if (*entries)[i].Key() == entry.Key() {
		oldEntry := (*entries)[i]
		(*entries)[i] = entry
		return oldEntry
	}

	*entries = append(*entries, nil)
	copy((*entries)[i+1:], (*entries)[i:])
	(*entries)[i] = entry
	return nil
}

func (entries *Entries) delete(key uint64) Entry {
	i := entries.search(key)
	if i >= len(*entries) {
		return nil
	}

	if (*entries)[i].Key() != key {
		return nil
	}

	oldEntry := (*entries)[i]
	copy((*entries)[i:], (*entries)[i+1:])
	(*entries)[len(*entries)-1] = nil
	*entries = (*entries)[:len(*entries)-1]
	return oldEntry
}

func (entries Entries) get(key uint64) Entry {
	i := entries.search(key)
	if i >= len(entries) {
		return nil
	}

	if entries[i].Key() == key {
		return entries[i]
	}

	return nil
}
