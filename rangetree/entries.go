package rangetree

import "sync"

var entriesPool = sync.Pool{
	New: func() interface{} {
		return make(Entries, 0, 10)
	},
}

// Entries is a typed list of Entry that can be reused if Dispose
// is called.
type Entries []Entry

// Dispose will free the resources consumed by this list and
// allow the list to be reused.
func (entries *Entries) Dispose() {
	for i := 0; i < len(*entries); i++ {
		(*entries)[i] = nil
	}

	*entries = (*entries)[:0]
	entriesPool.Put(*entries)
}

// NewEntries will return a reused list of entries.
func NewEntries() Entries {
	return entriesPool.Get().(Entries)
}
