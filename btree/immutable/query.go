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

package btree

import (
	"runtime"
	"sync"

	terr "github.com/Workiva/go-datastructures/threadsafe/err"
)

func (t *Tr) Apply(fn func(item *Item), keys ...interface{}) error {
	if t.Root == nil || len(keys) == 0 {
		return nil
	}

	positions := make(map[interface{}]int, len(keys))
	for i, key := range keys {
		positions[key] = i
	}

	chunks := splitValues(keys, runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	lerr := terr.New()
	result := make(Keys, len(keys))

	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			defer wg.Done()

			chunk := chunks[i]
			if len(chunk) == 0 {
				return
			}

			for _, value := range chunk {
				n, _, err := t.iterativeFindWithoutPath(value, t.Root)
				if err != nil {
					lerr.Set(err)
					return
				}

				if n == nil {
					continue
				}

				k, _ := n.searchKey(t.config.Comparator, value)
				if k != nil && t.config.Comparator(k.Value, value) == 0 {
					result[positions[value]] = k
				}
			}
		}(i)
	}

	wg.Wait()

	if lerr.Get() != nil {
		return lerr.Get()
	}

	for _, k := range result {
		if k == nil {
			continue
		}

		item := k.ToItem()
		fn(item)
	}

	return nil
}

// filter performs an after fetch filtering of the values in the provided node.
// Due to the nature of the UB-Tree, we may get results in the node that
// aren't in the provided range.  The returned list of keys is not necessarily
// in the correct row-major order.
func (t *Tr) filter(start, stop interface{}, n *Node, fn func(key *Key) bool) bool {
	for iter := n.iter(t.config.Comparator, start, stop); iter.next(); {
		id, _ := iter.value()
		if !fn(id) {
			return false
		}
	}

	return true
}

func (t *Tr) iter(start, stop interface{}, fn func(*Key) bool) error {
	if len(t.Root) == 0 {
		return nil
	}

	cur := start
	seen := make(map[string]struct{}, 10)

	for t.config.Comparator(stop, cur) > 0 {
		n, highestValue, err := t.iterativeFindWithoutPath(cur, t.Root)
		if err != nil {
			return err
		}

		if n == nil && highestValue == nil {
			break
		} else if n != nil {
			if _, ok := seen[string(n.ID)]; ok {
				break
			}
			if !t.filter(cur, stop, n, fn) {
				break
			}
		}

		cur = n.lastValue()
		seen[string(n.ID)] = struct{}{}
	}

	return nil
}

// iterativeFind searches for the node with the provided value.  This
// is an iterative function and returns an error if there was a problem
// with persistence.
func (t *Tr) iterativeFind(value interface{}, id ID) (*path, error) {
	if len(id) == 0 { // can't find a matching node
		return nil, nil
	}

	path := &path{}
	var n *Node
	var err error
	var i int
	var key *Key

	for {
		n, err = t.contextOrCachedNode(id, t.mutable)
		if err != nil {
			return nil, err
		}

		key, i = n.searchKey(t.config.Comparator, value)

		pb := &pathBundle{i: i, n: n}
		path.append(pb)
		if n.IsLeaf {
			return path, nil
		}
		id = key.ID()
	}

	return path, nil
}

func (t *Tr) iterativeFindWithoutPath(value interface{}, id ID) (*Node, interface{}, error) {
	var n *Node
	var err error
	var i int
	var key *Key
	var highestValue interface{}

	for {
		n, err = t.contextOrCachedNode(id, t.mutable)
		if err != nil {
			return nil, highestValue, err
		}

		if n.IsLeaf {
			if t.config.Comparator(n.lastValue(), value) < 0 {
				return nil, highestValue, nil
			}
			highestValue = n.lastValue()
			return n, highestValue, nil
		}

		key, i = n.searchKey(t.config.Comparator, value)
		if i < n.lenValues() {
			highestValue = n.valueAt(i)
		}
		id = key.ID()
	}

	return n, highestValue, nil
}
