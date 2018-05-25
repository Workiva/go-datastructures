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
	"log"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ephemeral struct {
	mp   map[string]*Payload
	lock sync.RWMutex
}

func (e *ephemeral) Save(items ...*Payload) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		e.mp[string(item.Key)] = item
	}

	return nil
}

func (e *ephemeral) Load(keys ...[]byte) ([]*Payload, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	if len(keys) == 0 {
		return nil, nil
	}

	items := make([]*Payload, 0, len(keys))
	for _, k := range keys {
		items = append(items, e.mp[string(k)])
	}

	return items, nil
}

const (
	maxValue = int64(100000)
)

func init() {
	rand.Seed(time.Now().Unix())
}

type valueSortWrapper struct {
	comparator Comparator
	values     []interface{}
}

func (v *valueSortWrapper) Len() int {
	return len(v.values)
}

func (v *valueSortWrapper) Swap(i, j int) {
	v.values[i], v.values[j] = v.values[j], v.values[i]
}

func (v *valueSortWrapper) Less(i, j int) bool {
	return v.comparator(v.values[i], v.values[j]) < 0
}

func (v *valueSortWrapper) sort() {
	sort.Sort(v)
}

func reverse(items items) items {
	for i := 0; i < len(items)/2; i++ {
		items[i], items[len(items)-1-i] = items[len(items)-1-i], items[i]
	}

	return items
}

var comparator = func(item1, item2 interface{}) int {
	int1, int2 := item1.(int64), item2.(int64)
	if int1 < int2 {
		return -1
	}

	if int1 > int2 {
		return 1
	}

	return 0
}

// orderedItems is going to contain our "master" copy of items in
// sorted order.  Because the operations on a flat list are well
// understood, we can use this type to do generative type testing and
// confirm the results.
type orderedItems []*Item

func (o orderedItems) Len() int {
	return len(o)
}

func (o orderedItems) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o orderedItems) Less(i, j int) bool {
	return comparator(o[i].Value, o[j].Value) < 0
}

func (o orderedItems) equal(item1, item2 *Item) bool {
	return comparator(item1.Value, item2.Value) == 0
}

func (o orderedItems) copy() orderedItems {
	cp := make(orderedItems, len(o))
	copy(cp, o)
	return cp
}

func (o orderedItems) search(value interface{}) int {
	return sort.Search(len(o), func(i int) bool {
		return comparator(o[i].Value, value) >= 0
	})
}

func (o orderedItems) add(item *Item) orderedItems {
	cp := make(orderedItems, len(o))
	copy(cp, o)
	i := cp.search(item.Value)
	if i < len(o) && o.equal(o[i], item) {
		cp[i] = item
		return cp
	}

	if i == len(cp) {
		cp = append(cp, item)
		return cp
	}

	cp = append(cp, nil)
	copy(cp[i+1:], cp[i:])
	cp[i] = item
	return cp
}

func (o orderedItems) delete(item *Item) orderedItems {
	i := o.search(item.Value)
	if i == len(o) {
		return o
	}

	if !o.equal(o[i], item) {
		return o
	}

	cp := make(orderedItems, len(o))
	copy(cp, o)

	copy(cp[i:], cp[i+1:])
	cp[len(cp)-1] = nil // or the zero value of T
	cp = cp[:len(cp)-1]
	return cp
}

func (o orderedItems) toItems() items {
	cp := make(items, 0, len(o))
	for _, item := range o {
		cp = append(cp, item)
	}

	return cp
}

func (o orderedItems) query(start, stop interface{}) items {
	items := make(items, 0, len(o))

	for i := o.search(start); i < len(o); i++ {
		if comparator(o[i], stop) > 0 {
			break
		}

		items = append(items, o[i])
	}

	return items
}

func generateRandomQuery() (interface{}, interface{}) {
	start := int64(rand.Intn(int(maxValue)))
	offset := int64(rand.Intn(100))
	return start, start + offset
}

func newItem(value interface{}) *Item {
	return &Item{
		Value:   value,
		Payload: newID(),
	}
}

func newEphemeral() Persister {
	return &ephemeral{
		mp: make(map[string]*Payload),
	}
}

type delayedPersister struct {
	Persister
}

func (d *delayedPersister) Load(keys ...[]byte) ([]*Payload, error) {
	time.Sleep(5 * time.Millisecond)
	return d.Persister.Load(keys...)
}

func newDelayed() Persister {
	return &delayedPersister{newEphemeral()}
}

func defaultConfig() Config {
	return Config{
		NodeWidth:  10, // easy number to test with
		Persister:  newEphemeral(),
		Comparator: comparator,
	}
}

func generateRandomItem() *Item {
	return newItem(int64(rand.Intn(int(maxValue))))
}

// generateRandomItems will generate a list of random items with
// no duplicates.
func generateRandomItems(num int) items {
	items := make(items, 0, num)
	mp := make(map[interface{}]struct{}, num)
	for len(items) < num {
		c := generateRandomItem()
		if _, ok := mp[c.Value]; ok {
			continue
		}
		mp[c.Value] = struct{}{}
		items = append(items, c)
	}

	return items
}

// generateLinearItems is similar to random item generation except that
// items are returned in sorted order.
func generateLinearItems(num int) items {
	items := make(items, 0, num)
	for i := 0; i < num; i++ {
		c := newItem(int64(i))
		items = append(items, c)
	}

	return items
}

func toOrdered(items items) orderedItems {
	oc := make(orderedItems, 0, len(items))
	for _, item := range items {
		oc = oc.add(item)
	}

	return oc
}

// the following 3 methods are in the _test file as they are only used
// in a testing environment.
func (t *Tr) toList(values ...interface{}) (items, error) {
	items := make(items, 0, t.Count)
	err := t.Apply(func(item *Item) {
		items = append(items, item)
	}, values...)

	return items, err
}

func (t *Tr) pprint(id ID) {
	n, _ := t.contextOrCachedNode(id, true)
	if n == nil {
		log.Printf(`NODE: %+v`, n)
		return
	}
	log.Printf(`NODE: %+v, LEN(ids): %+v, LEN(values): %+v`, n, n.lenKeys(), n.lenValues())
	for i, key := range n.ChildKeys {
		child, _ := t.contextOrCachedNode(key.ID(), true)
		if child == nil {
			continue
		}
		log.Printf(`CHILD %d: %+v`, i, child)
	}

	for _, key := range n.ChildKeys {
		child, _ := t.contextOrCachedNode(key.ID(), true)
		if child == nil {
			continue
		}
		t.pprint(key.ID())
	}
}

func (t *Tr) verify(id ID, tb testing.TB) (interface{}, interface{}) {
	n, err := t.contextOrCachedNode(id, true)
	require.NoError(tb, err)

	cp := n.copy() // copy the values and sort them, ensure node values are sorted
	cpValues := cp.ChildValues

	(&valueSortWrapper{comparator: comparator, values: cpValues}).sort()
	assert.Equal(tb, cpValues, n.ChildValues)

	if !assert.False(tb, n.needsSplit(t.config.NodeWidth)) {
		tb.Logf(`NODE NEEDS SPLIT: NODE: %+v`, n)
	}
	if string(t.Root) != string(n.ID) {
		assert.True(tb, n.lenValues() >= t.config.NodeWidth/2)
	}

	if n.IsLeaf {
		assert.Equal(tb, n.lenValues(), n.lenKeys()) // assert lens are equal
		return n.firstValue(), n.lastValue()         // return last value
	} else {
		for _, key := range n.ChildKeys {
			assert.Empty(tb, key.Payload)
		}
	}

	for i, key := range n.ChildKeys {
		min, max := t.verify(key.ID(), tb)
		if i == 0 {
			assert.True(tb, comparator(max, n.valueAt(i)) <= 0)
		} else if i == n.lenValues() {
			assert.True(tb, comparator(min, n.lastValue()) > 0)
		} else {
			assert.True(tb, comparator(max, n.valueAt(i)) <= 0)
			assert.True(tb, comparator(min, n.valueAt(i-1)) > 0)
		}
	}

	return n.firstValue(), n.lastValue()
}

func itemsToValues(items ...*Item) []interface{} {
	values := make([]interface{}, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}

	return values
}

func TestNodeSplit(t *testing.T) {
	number := 100
	items := generateLinearItems(number)
	cfg := defaultConfig()

	rt := New(cfg)
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.NoError(t, err)
	assert.Equal(t, number, mutable.Len())
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)

	result, err := mutable.(*Tr).toList(itemsToValues(items[5:10]...)...)
	require.NoError(t, err)
	if !assert.Equal(t, items[5:10], result) {
		mutable.(*Tr).pprint(mutable.(*Tr).Root)
		for i, c := range items[5:10] {
			t.Logf(`EXPECTED: %+v, RESULT: %+v`, c, result[i])
		}
		t.FailNow()
	}

	mutable = rt.AsMutable()
	for _, c := range items {
		_, err := mutable.AddItems(c)
		require.NoError(t, err)
	}

	result, err = mutable.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)

	rt, err = mutable.Commit()
	require.NoError(t, err)
	rt, err = Load(cfg.Persister, rt.ID(), comparator)

	result, err = mutable.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	rt.(*Tr).verify(rt.(*Tr).Root, t)
}

func TestReverseNodeSplit(t *testing.T) {
	number := 400
	items := generateLinearItems(number)

	reversed := make([]*Item, len(items))
	copy(reversed, items)
	reversed = reverse(reversed)
	rt := New(defaultConfig())
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(reversed...)
	require.NoError(t, err)

	result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	if !assert.Equal(t, items, result) {
		for _, c := range result {
			t.Logf(`RESULT: %+v`, c)
		}
	}

	mutable = rt.AsMutable()
	for _, c := range reversed {
		_, err := mutable.AddItems(c)
		require.NoError(t, err)
	}

	result, err = mutable.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
}

func TestDuplicate(t *testing.T) {
	item1 := newItem(int64(1))
	item2 := newItem(int64(1))

	rt := New(defaultConfig())
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(item1)
	require.NoError(t, err)
	_, err = mutable.AddItems(item2)
	require.NoError(t, err)

	assert.Equal(t, 1, mutable.Len())
	result, err := mutable.(*Tr).toList(int64(1))
	require.NoError(t, err)

	assert.Equal(t, items{item2}, result)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
}

func TestCommit(t *testing.T) {
	items := generateRandomItems(5)
	rt := New(defaultConfig())
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.Nil(t, err)

	rt, err = mutable.Commit()
	require.NoError(t, err)
	expected := toOrdered(items).toItems()
	result, err := rt.(*Tr).toList(itemsToValues(expected...)...)
	require.NoError(t, err)
	if !assert.Equal(t, expected, result) {
		require.Equal(t, len(expected), len(result))
		for i, c := range expected {
			if !assert.Equal(t, c, result[i]) {
				t.Logf(`EXPECTED: %+v, RESULT: %+v`, c, result[i])
			}
		}
	}

	rt.(*Tr).verify(rt.(*Tr).Root, t)
}

func TestRandom(t *testing.T) {
	items := generateRandomItems(1000)
	rt := New(defaultConfig())
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.Nil(t, err)

	require.NoError(t, err)
	expected := toOrdered(items).toItems()
	result, err := mutable.(*Tr).toList(itemsToValues(expected...)...)
	if !assert.Equal(t, expected, result) {
		assert.Equal(t, len(expected), len(result))
		for i, c := range expected {
			assert.Equal(t, c, result[i])
		}
	}
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
}

func TestLoad(t *testing.T) {
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(1000)
	_, err := mutable.AddItems(items...)
	require.NoError(t, err)

	id := mutable.ID()
	_, err = mutable.Commit()
	require.NoError(t, err)

	rt, err = Load(cfg.Persister, id, comparator)
	require.NoError(t, err)
	sort.Sort(orderedItems(items))
	result, err := rt.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	rt.(*Tr).verify(rt.(*Tr).Root, t)
}

func TestDeleteFromRoot(t *testing.T) {
	number := 5
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateLinearItems(number)

	mutable.AddItems(items...)
	mutable.DeleteItems(items[0].Value, items[1].Value, items[2].Value)

	result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)
	assert.Equal(t, items[3:], result)
	assert.Equal(t, 2, mutable.Len())

	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
}

func TestDeleteAllFromRoot(t *testing.T) {
	num := 5
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateLinearItems(num)

	mutable.AddItems(items...)
	mutable.DeleteItems(itemsToValues(items...)...)

	result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, mutable.Len())
}

func TestDeleteAfterSplitIncreasing(t *testing.T) {
	num := 11
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateLinearItems(num)

	mutable.AddItems(items...)
	for i := 0; i < num-1; i++ {
		mutable.DeleteItems(itemsToValues(items[i])...)
		result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
		require.Nil(t, err)
		assert.Equal(t, items[i+1:], result)
		mutable.(*Tr).verify(mutable.(*Tr).Root, t)
	}
}

func TestDeleteMultipleLevelsRandomlyBulk(t *testing.T) {
	num := 200
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(num)
	mutable.AddItems(items...)
	mutable.DeleteItems(itemsToValues(items[:100]...)...)
	result, _ := mutable.(*Tr).toList(itemsToValues(items...)...)
	assert.Len(t, result, 100)
}

func TestDeleteAfterSplitDecreasing(t *testing.T) {
	num := 11
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateLinearItems(num)

	mutable.AddItems(items...)
	for i := num - 1; i >= 0; i-- {
		mutable.DeleteItems(itemsToValues(items[i])...)
		result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
		require.Nil(t, err)
		assert.Equal(t, items[:i], result)
		if i > 0 {
			mutable.(*Tr).verify(mutable.(*Tr).Root, t)
		}
	}
}

func TestDeleteMultipleLevels(t *testing.T) {
	num := 20
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(num)
	mutable.AddItems(items...)
	ordered := toOrdered(items)

	for i, c := range ordered {
		_, err := mutable.DeleteItems(c.Value)
		require.NoError(t, err)
		result, err := mutable.(*Tr).toList(itemsToValues(ordered...)...)
		require.NoError(t, err)
		if !assert.Equal(t, ordered[i+1:].toItems(), result) {
			log.Printf(`LEN EXPECTED: %+v, RESULT: %+v`, len(ordered[i+1:]), len(result))
			mutable.(*Tr).pprint(mutable.(*Tr).Root)
			assert.Equal(t, len(ordered[i+1:]), len(result))
			for i, c := range ordered[i+1:] {
				log.Printf(`EXPECTED: %+v`, c)
				if i < len(result) {
					log.Printf(`RECEIVED: %+v`, result[i])
				}
			}
			break
		}
		if len(ordered[i+1:]) > 0 {
			mutable.(*Tr).verify(mutable.(*Tr).Root, t)
		}
	}

	assert.Nil(t, mutable.(*Tr).Root)
}

func TestDeleteMultipleLevelsRandomly(t *testing.T) {
	num := 200
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(num)
	mutable.AddItems(items...)
	ordered := toOrdered(items)

	for _, c := range items {
		_, err := mutable.DeleteItems(c.Value)
		require.NoError(t, err)
		ordered = ordered.delete(c)

		result, err := mutable.(*Tr).toList(itemsToValues(ordered...)...)
		require.NoError(t, err)
		assert.Equal(t, ordered.toItems(), result)
		if len(ordered) > 0 {
			mutable.(*Tr).verify(mutable.(*Tr).Root, t)
		}
	}

	assert.Nil(t, mutable.(*Tr).Root)
}

func TestDeleteMultipleLevelsWithCommit(t *testing.T) {
	num := 20
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(num)
	mutable.AddItems(items...)
	rt, _ = mutable.Commit()

	rt, _ = Load(cfg.Persister, rt.ID(), comparator)
	result, err := rt.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	mutable = rt.AsMutable()

	for _, c := range items[:10] {
		_, err := mutable.DeleteItems(c.Value)
		require.Nil(t, err)
	}

	result, err = mutable.(*Tr).toList(itemsToValues(items[10:]...)...)
	require.Nil(t, err)
	assert.Equal(t, items[10:], result)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)

	result, err = rt.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
	rt.(*Tr).verify(rt.(*Tr).Root, t)
}

func TestCommitAfterDelete(t *testing.T) {
	num := 15
	cfg := defaultConfig()
	rt := New(cfg)
	mutable := rt.AsMutable()
	items := generateRandomItems(num)
	mutable.AddItems(items...)
	for _, c := range items[:5] {
		mutable.DeleteItems(c.Value)
		mutable.(*Tr).verify(mutable.(*Tr).Root, t)
	}

	rt, err := mutable.Commit()
	require.Nil(t, err)
	result, err := rt.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)
	assert.Equal(t, items[5:], result)
	rt.(*Tr).verify(rt.(*Tr).Root, t)

}

func TestSecondCommitSplitsRoot(t *testing.T) {
	number := 15
	cfg := defaultConfig()
	rt := New(cfg)
	items := generateLinearItems(number)

	mutable := rt.AsMutable()
	mutable.AddItems(items[:10]...)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
	rt, _ = mutable.Commit()
	rt.(*Tr).verify(rt.(*Tr).Root, t)
	mutable = rt.AsMutable()
	mutable.AddItems(items[10:]...)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)

	result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)

	if !assert.Equal(t, items, result) {
		for i, c := range items {
			log.Printf(`EXPECTED: %+v, RECEIVED: %+v`, c, result[i])
		}
	}
}

func TestSecondCommitMultipleSplits(t *testing.T) {
	num := 50
	cfg := defaultConfig()
	rt := New(cfg)
	items := generateRandomItems(num)

	mutable := rt.AsMutable()
	mutable.AddItems(items[:25]...)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
	rt, err := mutable.Commit()
	rt.(*Tr).verify(rt.(*Tr).Root, t)

	result, err := rt.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)
	assert.Equal(t, items[:25], result)

	mutable = rt.AsMutable()
	mutable.AddItems(items[25:]...)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)

	sort.Sort(orderedItems(items))
	result, err = mutable.(*Tr).toList(itemsToValues(items...)...)
	require.Nil(t, err)
	if !assert.Equal(t, items, result) {
		mutable.(*Tr).pprint(mutable.(*Tr).Root)
	}
}

func TestLargeAdd(t *testing.T) {
	cfg := defaultConfig()
	number := cfg.NodeWidth * 5
	rt := New(cfg)
	items := generateLinearItems(number)

	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.NoError(t, err)

	id := mutable.ID()
	result, err := mutable.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)

	_, err = mutable.Commit()
	require.NoError(t, err)

	rt, err = Load(cfg.Persister, id, comparator)
	require.NoError(t, err)
	result, err = rt.(*Tr).toList(itemsToValues(items...)...)
	require.NoError(t, err)
	assert.Equal(t, items, result)
}

func TestNodeInfiniteLoop(t *testing.T) {
	cfg := defaultConfig()
	rt := New(cfg)
	items := generateLinearItems(3)

	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.NoError(t, err)

	result, err := mutable.DeleteItems(items[1].Value, items[2].Value)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// all remaining tests are generative in nature to catch things
// I can't think of.

func TestGenerativeAdds(t *testing.T) {
	if testing.Short() {
		t.Skipf(`skipping generative add`)
		return
	}

	number := 100
	cfg := defaultConfig()
	rt := New(cfg)
	oc := make(orderedItems, 0)
	for i := 0; i < number; i++ {
		num := int(rand.Int31n(100))
		if num == 0 {
			num++
		}

		items := generateRandomItems(num)
		mutated := oc.copy()
		for _, c := range items {
			mutated = mutated.add(c)
		}

		mutable := rt.AsMutable()
		_, err := mutable.AddItems(items...)
		require.Nil(t, err)
		mutable.(*Tr).verify(mutable.(*Tr).Root, t)

		rtMutated, err := mutable.Commit()
		require.Nil(t, err)
		rtMutated.(*Tr).verify(rtMutated.(*Tr).Root, t)

		result, err := rtMutated.(*Tr).toList(itemsToValues(mutated.toItems()...)...)
		require.Nil(t, err)
		if !assert.Equal(t, mutated.toItems(), result) {
			rtMutated.(*Tr).pprint(rtMutated.(*Tr).Root)
			if len(mutated) == len(result) {
				for i, c := range mutated.toItems() {
					log.Printf(`EXPECTED: %+v, RECEIVED: %+v`, c, result[i])
				}
			}
		}
		assert.Equal(t, len(mutated), rtMutated.Len())

		result, err = rt.(*Tr).toList(itemsToValues(oc.toItems()...)...)
		require.Nil(t, err)
		assert.Equal(t, oc.toItems(), result)

		oc = mutated
		rt = rtMutated
	}
}

func TestGenerativeDeletes(t *testing.T) {
	if testing.Short() {
		t.Skipf(`skipping generative delete`)
		return
	}

	number := 100
	var err error
	cfg := defaultConfig()
	rt := New(cfg)
	oc := toOrdered(generateRandomItems(1000))
	mutable := rt.AsMutable()
	mutable.AddItems(oc.toItems()...)
	mutable.(*Tr).verify(mutable.(*Tr).Root, t)
	rt, err = mutable.Commit()
	require.NoError(t, err)
	for i := 0; i < number; i++ {
		mutable = rt.AsMutable()
		index := rand.Intn(len(oc))
		c := oc[index]
		mutated := oc.delete(c)

		result, err := rt.(*Tr).toList(itemsToValues(oc.toItems()...)...)
		require.NoError(t, err)
		assert.Equal(t, oc.toItems(), result)
		assert.Equal(t, len(oc), rt.Len())

		_, err = mutable.DeleteItems(c.Value)
		require.NoError(t, err)
		mutable.(*Tr).verify(mutable.(*Tr).Root, t)
		result, err = mutable.(*Tr).toList(itemsToValues(mutated.toItems()...)...)
		require.NoError(t, err)
		assert.Equal(t, len(mutated), len(result))
		require.Equal(t, mutated.toItems(), result)
		oc = mutated
		rt, err = mutable.Commit()
		require.NoError(t, err)
	}
}

func TestGenerativeOperations(t *testing.T) {
	if testing.Short() {
		t.Skipf(`skipping generative operations`)
		return
	}

	number := 100
	cfg := defaultConfig()
	rt := New(cfg)

	// seed the tree
	items := generateRandomItems(1000)
	oc := toOrdered(items)

	mutable := rt.AsMutable()
	mutable.AddItems(items...)

	result, err := mutable.(*Tr).toList(itemsToValues(oc.toItems()...)...)
	require.NoError(t, err)
	require.Equal(t, oc.toItems(), result)

	rt, err = mutable.Commit()
	require.NoError(t, err)

	for i := 0; i < number; i++ {
		mutable = rt.AsMutable()
		if rand.Float64() < .5 && len(oc) > 0 {
			c := oc[rand.Intn(len(oc))]
			oc = oc.delete(c)
			_, err = mutable.DeleteItems(c.Value)
			require.NoError(t, err)
			mutable.(*Tr).verify(mutable.(*Tr).Root, t)
			result, err := mutable.(*Tr).toList(itemsToValues(oc.toItems()...)...)
			require.NoError(t, err)
			require.Equal(t, oc.toItems(), result)
			assert.Equal(t, len(oc), mutable.Len())
		} else {
			c := generateRandomItem()
			oc = oc.add(c)
			_, err = mutable.AddItems(c)
			require.NoError(t, err)
			mutable.(*Tr).verify(mutable.(*Tr).Root, t)
			result, err = mutable.(*Tr).toList(itemsToValues(oc.toItems()...)...)
			require.NoError(t, err)
			require.Equal(t, oc.toItems(), result)
			assert.Equal(t, len(oc), mutable.Len())
		}

		rt, err = mutable.Commit()
		require.NoError(t, err)
	}
}

func BenchmarkGetitems(b *testing.B) {
	number := 100
	cfg := defaultConfig()
	cfg.Persister = newDelayed()
	rt := New(cfg)

	items := generateRandomItems(number)
	mutable := rt.AsMutable()
	_, err := mutable.AddItems(items...)
	require.NoError(b, err)

	rt, err = mutable.Commit()
	require.NoError(b, err)
	id := rt.ID()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt, err = Load(cfg.Persister, id, comparator)
		require.NoError(b, err)
		_, err = rt.(*Tr).toList(itemsToValues(items...)...)
		require.NoError(b, err)
	}
}

func BenchmarkBulkAdd(b *testing.B) {
	number := 1000000
	items := generateLinearItems(number)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New(defaultConfig())
		mutable := tr.AsMutable()
		mutable.AddItems(items...)
	}
}
