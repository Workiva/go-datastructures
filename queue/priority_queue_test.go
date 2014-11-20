package queue

import (
	//"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityPut(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Put(mockItem(2))

	assert.Len(t, q.items, 1)
	assert.Equal(t, mockItem(2), q.items[0])

	q.Put(mockItem(1))

	if !assert.Len(t, q.items, 2) {
		return
	}
	assert.Equal(t, mockItem(1), q.items[0])
	assert.Equal(t, mockItem(2), q.items[1])
}

func TestPriorityGet(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Put(mockItem(2))
	result, err := q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, result, 1) {
		return
	}

	assert.Equal(t, mockItem(2), result[0])
	assert.Len(t, q.items, 0)

	q.Put(mockItem(2))
	q.Put(mockItem(1))

	result, err = q.Get(1)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, result, 1) {
		return
	}

	assert.Equal(t, mockItem(1), result[0])
	assert.Len(t, q.items, 1)

	result, err = q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, result, 1) {
		return
	}

	assert.Equal(t, mockItem(2), result[0])
}

func TestAddEmptyPriorityPut(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Put()

	assert.Len(t, q.items, 0)
}

func TestPriorityGetNonPositiveNumber(t *testing.T) {
	q := NewPriorityQueue(1)

	q.Put(mockItem(1))

	result, err := q.Get(0)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 0)

	result, err = q.Get(-1)
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, result, 0)
}

func TestPriorityEmpty(t *testing.T) {
	q := NewPriorityQueue(1)
	assert.True(t, q.Empty())

	q.Put(mockItem(1))

	assert.False(t, q.Empty())
}

func TestPriorityGetEmpty(t *testing.T) {
	q := NewPriorityQueue(1)

	go func() {
		q.Put(mockItem(1))
	}()

	result, err := q.Get(1)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, result, 1) {
		return
	}
	assert.Equal(t, mockItem(1), result[0])
}

func TestMultiplePriorityGetEmpty(t *testing.T) {
	q := NewPriorityQueue(1)
	var wg sync.WaitGroup
	wg.Add(2)
	results := make([][]Item, 2)

	go func() {
		wg.Done()
		local, _ := q.Get(1)
		results[0] = local
		wg.Done()
	}()

	go func() {
		wg.Done()
		local, _ := q.Get(1)
		results[1] = local
		wg.Done()
	}()

	wg.Wait()
	wg.Add(2)

	q.Put(mockItem(1), mockItem(3), mockItem(2))
	wg.Wait()

	if !assert.Len(t, results[0], 1) || !assert.Len(t, results[1], 1) {
		return
	}

	assert.True(
		t, (results[0][0] == mockItem(1) && results[1][0] == mockItem(2)) ||
			results[0][0] == mockItem(2) && results[1][0] == mockItem(1),
	)
}

func TestEmptyPriorityGetWithDispose(t *testing.T) {
	q := NewPriorityQueue(1)
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

	assert.IsType(t, DisposedError{}, err)
}

func TestPriorityGetPutDisposed(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Dispose()

	_, err := q.Get(1)
	assert.IsType(t, DisposedError{}, err)

	err = q.Put(mockItem(1))
	assert.IsType(t, DisposedError{}, err)
}

func BenchmarkPriorityQueue(b *testing.B) {
	q := NewPriorityQueue(b.N)
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
		q.Put(mockItem(i))
	}

	wg.Wait()
}

func TestPriorityPeek(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Put(mockItem(1))

	assert.Equal(t, mockItem(1), q.Peek())
}

func TestInsertDuplicate(t *testing.T) {
	q := NewPriorityQueue(1)
	q.Put(mockItem(1))
	q.Put(mockItem(1))

	assert.Equal(t, 1, q.Len())
}
