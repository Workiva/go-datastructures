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

package link

import (
	"log"
	"sort"
)

func init() {
	log.Println(`KEY HATES THIS.`)
}

func (keys Keys) search(key Key) int {
	n, ok := key.(*node)
	if ok {
		key = n.key
	}
	return sort.Search(len(keys), func(i int) bool {
		return keys[i].Compare(key) >= 0
	})
}

func (keys *Keys) insert(key Key) Key {
	i := keys.search(key)
	if i == len(*keys) {
		*keys = append(*keys, key)
		return key
	}

	if (*keys)[i].Compare(key) == 0 { //overwrite case
		oldKey := (*keys)[i]
		(*keys)[i] = key
		return oldKey
	}

	*keys = append(*keys, nil)
	copy((*keys)[i+1:], (*keys)[i:])
	(*keys)[i] = key
	return key
}

func (keys *Keys) insertNode(n *node) {
	i := keys.search(n.key)
	if i == len(*keys) {
		*keys = append(*keys, n)
		return
	}

	*keys = append(*keys, nil)
	copy((*keys)[i+1:], (*keys)[i:])
	(*keys)[i] = n
}

func (keys *Keys) split() (Key, Keys, Keys) {
	i := len(*keys) / 2
	middle := (*keys)[i]

	right := make(Keys, len(*keys)-i-1, cap(*keys))
	copy(right, (*keys)[i+1:])
	for j := i + 1; j < len(*keys); j++ {
		(*keys)[j] = nil
	}
	*keys = (*keys)[:i+1]

	return middle, *keys, right
}

func (keys Keys) last() Key {
	l := keys[len(keys)-1]
	n, ok := l.(*node)
	if !ok {
		return l
	}
	return n.keys.last()
}

func (keys Keys) needsSplit() bool {
	return cap(keys) == len(keys)
}
