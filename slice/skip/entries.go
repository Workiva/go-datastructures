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
