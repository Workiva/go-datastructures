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

package queue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	q := New(10)

	q.Put(`test`)
	assert.Equal(t, int64(1), q.Len())

	results, err := q.Get(1)
	assert.Nil(t, err)

	result := results[0]
	assert.Equal(t, `test`, result)
	assert.True(t, q.Empty())

	q.Put(`test2`)
	assert.Equal(t, int64(1), q.Len())

	results, err = q.Get(1)
	assert.Nil(t, err)

	result = results[0]
	assert.Equal(t, `test2`, result)
	assert.True(t, q.Empty())
}

func TestGet(t *testing.T) {
	q := New(10)

	q.Put(`test`)
	result, err := q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `test`, result[0])
	assert.Equal(t, int64(0), q.Len())

	q.Put(`1`)
	q.Put(`2`)

	result, err = q.Get(1)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `1`, result[0])
	assert.Equal(t, int64(1), q.Len())

	result, err = q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, `2`, result[0])
}

func TestPoll(t *testing.T) {
	q := New(10)

	q.Put(`test`)
	result, err := q.Poll(2, 0)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `test`, result[0])
	assert.Equal(t, int64(0), q.Len())

	q.Put(`1`)
	q.Put(`2`)

	result, err = q.Poll(1, time.Millisecond)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `1`, result[0])
	assert.Equal(t, int64(1), q.Len())

	result, err = q.Poll(2, time.Millisecond)
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, `2`, result[0])

	before := time.Now()
	_, err = q.Poll(1, 5*time.Millisecond)
	assert.InDelta(t, 5, time.Since(before).Seconds()*1000, 2)
	assert.Equal(t, ErrTimeout, err)
}

func TestAddEmptyPut(t *testing.T) {
	q := New(10)

	q.Put()

	if q.Len() != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, q.Len())
	}
}

func TestGetNonPositiveNumber(t *testing.T) {
	q := New(10)

	q.Put(`test`)
	result, err := q.Get(0)
	if !assert.Nil(t, err) {
		return
	}

	if len(result) != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, len(result))
	}
}

func TestEmpty(t *testing.T) {
	q := New(10)

	if !q.Empty() {
		t.Errorf(`Expected empty queue.`)
	}

	q.Put(`test`)
	if q.Empty() {
		t.Errorf(`Expected non-empty queue.`)
	}
}

func TestGetEmpty(t *testing.T) {
	q := New(10)

	go func() {
		q.Put(`a`)
	}()

	result, err := q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 1)
	assert.Equal(t, `a`, result[0])
}

func TestMultipleGetEmpty(t *testing.T) {
	q := New(10)
	var wg sync.WaitGroup
	wg.Add(2)
	results := make([][]interface{}, 2)

	go func() {
		wg.Done()
		local, err := q.Get(1)
		assert.Nil(t, err)
		results[0] = local
		wg.Done()
	}()

	go func() {
		wg.Done()
		local, err := q.Get(1)
		assert.Nil(t, err)
		results[1] = local
		wg.Done()
	}()

	wg.Wait()
	wg.Add(2)

	q.Put(`a`, `b`, `c`)
	wg.Wait()

	if assert.Len(t, results[0], 1) && assert.Len(t, results[1], 1) {
		assert.True(t, (results[0][0] == `a` && results[1][0] == `b`) ||
			(results[0][0] == `b` && results[1][0] == `a`),
			`The array should be a, b or b, a`)
	}
}

func TestEmptyGetWithDispose(t *testing.T) {
	q := New(10)
	var wg sync.WaitGroup
	wg.Add(1)

	var err error

	go func() {
		wg.Done()
		_, err = q.Get(1)
		wg.Done()
	}()

	wg.Wait()
	wg.Add(1)

	q.Dispose()

	wg.Wait()

	assert.IsType(t, ErrDisposed, err)
}

func TestGetPutDisposed(t *testing.T) {
	q := New(10)

	q.Dispose()

	_, err := q.Get(1)
	assert.IsType(t, ErrDisposed, err)

	err = q.Put(`a`)
	assert.IsType(t, ErrDisposed, err)
}

func BenchmarkQueue(b *testing.B) {
	q := New(int64(b.N))
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0

	go func() {
		for {
			q.Get(1)
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		q.Put(`a`)
	}

	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan interface{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0

	go func() {
		for {
			<-ch
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		ch <- `a`
	}

	wg.Wait()
}

func TestTakeUntil(t *testing.T) {
	q := New(10)
	q.Put(`a`, `b`, `c`)
	result, err := q.TakeUntil(func(item interface{}) bool {
		return item != `c`
	})

	if !assert.Nil(t, err) {
		return
	}

	expected := []interface{}{`a`, `b`}
	assert.Equal(t, expected, result)
}

func TestTakeUntilEmptyQueue(t *testing.T) {
	q := New(10)
	result, err := q.TakeUntil(func(item interface{}) bool {
		return item != `c`
	})

	if !assert.Nil(t, err) {
		return
	}

	expected := []interface{}{}
	assert.Equal(t, expected, result)
}

func TestTakeUntilOnDisposedQueue(t *testing.T) {
	q := New(10)
	q.Dispose()
	result, err := q.TakeUntil(func(item interface{}) bool {
		return true
	})

	assert.Nil(t, result)
	assert.IsType(t, ErrDisposed, err)
}

func TestExecuteInParallel(t *testing.T) {
	q := New(10)
	for i := 0; i < 10; i++ {
		q.Put(i)
	}

	numCalls := uint64(0)

	ExecuteInParallel(q, func(item interface{}) {
		t.Logf("ExecuteInParallel called us with %+v", item)
		atomic.AddUint64(&numCalls, 1)
	})

	assert.Equal(t, uint64(10), numCalls)
	assert.True(t, q.Disposed())
}

func TestExecuteInParallelEmptyQueue(t *testing.T) {
	q := New(1)

	// basically just ensuring we don't deadlock here
	ExecuteInParallel(q, func(interface{}) {
		t.Fail()
	})
}

func BenchmarkQueuePut(b *testing.B) {
	numItems := int64(1000)

	qs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		q := New(10)
		qs = append(qs, q)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := qs[i]
		for j := int64(0); j < numItems; j++ {
			q.Put(j)
		}
	}
}

func BenchmarkQueueGet(b *testing.B) {
	numItems := int64(1000)

	qs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		q := New(numItems)
		for j := int64(0); j < numItems; j++ {
			q.Put(j)
		}
		qs = append(qs, q)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := qs[i]
		for j := int64(0); j < numItems; j++ {
			q.Get(1)
		}
	}
}

func BenchmarkQueuePoll(b *testing.B) {
	numItems := int64(1000)

	qs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		q := New(numItems)
		for j := int64(0); j < numItems; j++ {
			q.Put(j)
		}
		qs = append(qs, q)
	}

	b.ResetTimer()

	for _, q := range qs {
		for j := int64(0); j < numItems; j++ {
			q.Poll(1, time.Millisecond)
		}
	}
}

func BenchmarkExecuteInParallel(b *testing.B) {
	numItems := int64(1000)

	qs := make([]*Queue, 0, b.N)

	for i := 0; i < b.N; i++ {
		q := New(numItems)
		for j := int64(0); j < numItems; j++ {
			q.Put(j)
		}
		qs = append(qs, q)
	}

	var counter int64
	fn := func(ifc interface{}) {
		c := ifc.(int64)
		atomic.AddInt64(&counter, c)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := qs[i]
		ExecuteInParallel(q, fn)
	}
}
