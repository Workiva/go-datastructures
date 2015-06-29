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

	"github.com/stretchr/testify/assert"
)

func TestRingInsert(t *testing.T) {
	rb := NewRingBuffer(5)
	assert.Equal(t, uint64(8), rb.Cap())

	err := rb.Put(5)
	if !assert.Nil(t, err) {
		return
	}

	result, err := rb.Get()
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, 5, result)
}

func TestRingMultipleInserts(t *testing.T) {
	rb := NewRingBuffer(5)

	err := rb.Put(1)
	if !assert.Nil(t, err) {
		return
	}

	err = rb.Put(2)
	if !assert.Nil(t, err) {
		return
	}

	result, err := rb.Get()
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, 1, result)

	result, err = rb.Get()
	if assert.Nil(t, err) {
		return
	}

	assert.Equal(t, 2, result)
}

func TestIntertwinedGetAndPut(t *testing.T) {
	rb := NewRingBuffer(5)
	err := rb.Put(1)
	if !assert.Nil(t, err) {
		return
	}

	result, err := rb.Get()
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, 1, result)

	err = rb.Put(2)
	if !assert.Nil(t, err) {
		return
	}

	result, err = rb.Get()
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, 2, result)
}

func TestPutToFull(t *testing.T) {
	rb := NewRingBuffer(3)

	for i := 0; i < 4; i++ {
		err := rb.Put(i)
		if !assert.Nil(t, err) {
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := rb.Put(4)
		assert.Nil(t, err)
		wg.Done()
	}()

	go func() {
		defer wg.Done()
		result, err := rb.Get()
		if !assert.Nil(t, err) {
			return
		}

		assert.Equal(t, 0, result)
	}()

	wg.Wait()
}

func TestOffer(t *testing.T) {
	rb := NewRingBuffer(2)

	ok, err := rb.Offer("foo")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = rb.Offer("bar")
	assert.True(t, ok)
	assert.Nil(t, err)
	ok, err = rb.Offer("baz")
	assert.False(t, ok)
	assert.Nil(t, err)

	item, err := rb.Get()
	assert.Nil(t, err)
	assert.Equal(t, "foo", item)
	item, err = rb.Get()
	assert.Nil(t, err)
	assert.Equal(t, "bar", item)
}

func TestRingGetEmpty(t *testing.T) {
	rb := NewRingBuffer(3)

	var wg sync.WaitGroup
	wg.Add(1)

	// want to kick off this consumer to ensure it blocks
	go func() {
		wg.Done()
		result, err := rb.Get()
		assert.Nil(t, err)
		assert.Equal(t, 0, result)
		wg.Done()
	}()

	wg.Wait()
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := rb.Put(0)
		assert.Nil(t, err)
	}()

	wg.Wait()
}

func TestRingLen(t *testing.T) {
	rb := NewRingBuffer(4)
	assert.Equal(t, uint64(0), rb.Len())

	rb.Put(1)
	assert.Equal(t, uint64(1), rb.Len())

	rb.Get()
	assert.Equal(t, uint64(0), rb.Len())

	for i := 0; i < 4; i++ {
		rb.Put(1)
	}
	assert.Equal(t, uint64(4), rb.Len())

	rb.Get()
	assert.Equal(t, uint64(3), rb.Len())
}

func TestDisposeOnGet(t *testing.T) {
	numThreads := 8
	var wg sync.WaitGroup
	wg.Add(numThreads)
	rb := NewRingBuffer(4)
	var spunUp sync.WaitGroup
	spunUp.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		go func() {
			spunUp.Done()
			defer wg.Done()
			_, err := rb.Get()
			assert.NotNil(t, err)
		}()
	}

	spunUp.Wait()
	rb.Dispose()

	wg.Wait()
	assert.True(t, rb.IsDisposed())
}

func TestDisposeOnPut(t *testing.T) {
	numThreads := 8
	var wg sync.WaitGroup
	wg.Add(numThreads)
	rb := NewRingBuffer(4)
	var spunUp sync.WaitGroup
	spunUp.Add(numThreads)

	// fill up the queue
	for i := 0; i < 4; i++ {
		rb.Put(i)
	}

	// it's now full
	for i := 0; i < numThreads; i++ {
		go func(i int) {
			spunUp.Done()
			defer wg.Done()
			err := rb.Put(i)
			assert.NotNil(t, err)
		}(i)
	}

	spunUp.Wait()

	rb.Dispose()

	wg.Wait()

	assert.True(t, rb.IsDisposed())
}

func BenchmarkRBLifeCycle(b *testing.B) {
	rb := NewRingBuffer(64)

	counter := uint64(0)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			_, err := rb.Get()
			assert.Nil(b, err)

			if atomic.AddUint64(&counter, 1) == uint64(b.N) {
				return
			}
		}
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rb.Put(i)
	}

	wg.Wait()
}

func BenchmarkRBPut(b *testing.B) {
	rbs := make([]*RingBuffer, 0, b.N)

	for i := 0; i < b.N; i++ {
		rbs = append(rbs, NewRingBuffer(2))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbs[i].Put(i)
	}
}

func BenchmarkRBGet(b *testing.B) {
	rbs := make([]*RingBuffer, 0, b.N)

	for i := 0; i < b.N; i++ {
		rbs = append(rbs, NewRingBuffer(2))
		rbs[i].Put(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rbs[i].Get()
	}
}
