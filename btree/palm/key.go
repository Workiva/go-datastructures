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

package palm

import "sort"

func (keys Keys) search(key Key) int {
	return sort.Search(len(keys), func(i int) bool {
		return keys[i].Compare(key) > -1
	})
}

func (keys *Keys) insert(key Key) Key {
	i := keys.search(key)
	return keys.insertAt(key, i)
}

func (keys *Keys) insertAt(key Key, i int) Key {
	if i == len(*keys) {
		*keys = append(*keys, key)
		return nil
	}

	if (*keys)[i].Compare(key) == 0 { //overwrite case
		oldKey := (*keys)[i]
		(*keys)[i] = key
		return oldKey
	}

	*keys = append(*keys, nil)
	copy((*keys)[i+1:], (*keys)[i:])
	(*keys)[i] = key
	return nil
}

func (keys *Keys) splitAt(i int) (Keys, Keys) {
	right := make(Keys, len(*keys)-i-1, cap(*keys))
	copy(right, (*keys)[i+1:])
	for j := i + 1; j < len(*keys); j++ {
		(*keys)[j] = nil
	}
	*keys = (*keys)[:i+1]

	return *keys, right
}

func (keys Keys) reverse() Keys {
	reversed := make(Keys, len(keys))
	for i := len(keys) - 1; i >= 0; i-- {
		reversed[len(keys)-1-i] = keys[i]
	}

	return reversed
}

func chunkKeys(keys Keys, numParts int64) []Keys {
	parts := make([]Keys, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = keys[i*int64(len(keys))/numParts : (i+1)*int64(len(keys))/numParts]
	}
	return parts
}
