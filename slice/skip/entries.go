/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package skip

import "sort"

// search will look for the provided key and return the index
// where that key would be inserted in this list.  This could
// be equal to the length of the list which means no suitable entry
// point was found.
func (entries Entries) search(e Entry) int {
	return sort.Search(len(entries), func(i int) bool {
		return entries[i].Compare(e) > -1
	})
}

// insert will insert the provided entry into this list.
func (entries *Entries) insert(entry Entry) Entry {
	i := entries.search(entry)
	if i >= len(*entries) {
		*entries = append(*entries, entry)
		return nil
	}

	if (*entries)[i].Compare(entry) == 0 {
		oldEntry := (*entries)[i]
		(*entries)[i] = entry
		return oldEntry
	}

	*entries = append(*entries, nil)
	copy((*entries)[i+1:], (*entries)[i:])
	(*entries)[i] = entry
	return nil
}

// delete will delete the provided key from this list.
func (entries *Entries) delete(e Entry) Entry {
	i := entries.search(e)
	if i >= len(*entries) {
		return nil
	}

	if (*entries)[i].Compare(e) != 0 {
		return nil
	}

	oldEntry := (*entries)[i]
	copy((*entries)[i:], (*entries)[i+1:])
	(*entries)[len(*entries)-1] = nil
	*entries = (*entries)[:len(*entries)-1]
	return oldEntry
}

// get will return the entry associated with the provided key.
// If no such Entry exists, this returns nil.
func (entries Entries) get(e Entry) Entry {
	i := entries.search(e)
	if i >= len(entries) {
		return nil
	}

	if entries[i].Compare(e) == 0 {
		return entries[i]
	}

	return nil
}
