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

package rangetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedAdd(t *testing.T) {
	nodes := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	overwritten := nodes.add(n1)
	assert.Nil(t, overwritten)

	overwritten = nodes.add(n2)
	assert.Nil(t, overwritten)

	assert.Equal(t, orderedNodes{n2, n1}, nodes)

	n3 := newNode(4, constructMockEntry(1, 4), false)

	overwritten = nodes.add(n3)

	assert.True(t, n1 == overwritten)
	assert.Equal(t, orderedNodes{n2, n3}, nodes)
}

func TestOrderedDelete(t *testing.T) {
	nodes := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	nodes.add(n1)
	nodes.add(n2)

	deleted := nodes.delete(n2.value)

	assert.Equal(t, orderedNodes{n1}, nodes)
	assert.Equal(t, n2, deleted)

	missingValue := int64(3)
	deleted = nodes.delete(missingValue)

	assert.Equal(t, orderedNodes{n1}, nodes)
	assert.Nil(t, deleted)

	deleted = nodes.delete(n1.value)

	assert.Empty(t, nodes)
	assert.Equal(t, n1, deleted)
}

func TestApply(t *testing.T) {
	ns := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	ns.add(n1)
	ns.add(n2)

	results := make(nodes, 0, 2)

	ns.apply(1, 2, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n2}, results)

	results = results[:0]

	ns.apply(0, 1, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(2, 4, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(4, 5, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n1}, results)
	results = results[:0]

	ns.apply(0, 5, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n2, n1}, results)
	results = results[:0]

	ns.apply(5, 10, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(0, 100, func(n *node) bool {
		results = append(results, n)
		return false
	})

	assert.Equal(t, nodes{n2}, results)
}

func TestInsertDelete(t *testing.T) {
	ns := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)
	n3 := newNode(2, constructMockEntry(3, 2), false)

	ns.add(n1)
	ns.add(n2)
	ns.add(n3)

	modified := make(Entries, 0, 1)
	deleted := make(Entries, 0, 1)

	ns.insert(2, 2, 2, 0, -5, &modified, &deleted)

	assert.Len(t, ns, 0)
	assert.Equal(t, Entries{n2.entry, n3.entry, n1.entry}, deleted)
}

func BenchmarkPrepend(b *testing.B) {
	numItems := 100000
	ns := make(orderedNodes, 0, numItems)

	for i := b.N; i < b.N+numItems; i++ {
		ns.add(newNode(int64(i), constructMockEntry(uint64(i), int64(i)), false))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ns.add(newNode(int64(i), constructMockEntry(uint64(i), int64(i)), false))
	}
}
