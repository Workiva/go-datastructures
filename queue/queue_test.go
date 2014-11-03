package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	q := &Queue{}

	q.Put(`test`)
	if q.Len() != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, q.Len())
	}

	if q.items[0] != `test` {
		t.Errorf(`Expected: %s, received: %s`, `test`, q.items[0])
	}

	q.Put(`test2`)
	if q.Len() != 2 {
		t.Errorf(`Expected len: %d, received: %d`, 2, q.Len())
	}

	if q.items[1] != `test2` {
		t.Errorf(`Expected: %s, received: %s`, `test2`, q.items[1])
	}
}

func TestGet(t *testing.T) {
	q := &Queue{}

	q.Put(`test`)
	result, err := q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	if len(result) != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, len(result))
	}

	if result[0] != `test` {
		t.Errorf(`Expected: %s, received: %s`, `test`, result)
	}

	if len(q.items) != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, len(q.items))
	}

	q.Put(`1`)
	q.Put(`2`)

	result, err = q.Get(1)
	if !assert.Nil(t, err) {
		return
	}

	if len(result) != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, len(result))
	}

	if result[0] != `1` {
		t.Errorf(`Expected: %s, received: %s`, `1`, result[0])
	}

	if len(q.items) != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, len(q.items))
	}

	result, err = q.Get(2)
	if !assert.Nil(t, err) {
		return
	}

	if result[0] != `2` {
		t.Errorf(`Expected: %s, received: %s`, `2`, result[0])
	}
}

func TestAddEmptyPut(t *testing.T) {
	q := &Queue{}

	q.Put()

	if q.Len() != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, q.Len())
	}
}

func TestGetNonPositiveNumber(t *testing.T) {
	q := &Queue{}

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
	q := &Queue{}

	if !q.Empty() {
		t.Errorf(`Expected empty queue.`)
	}

	q.Put(`test`)
	if q.Empty() {
		t.Errorf(`Expected non-empty queue.`)
	}
}

func TestGetEmpty(t *testing.T) {
	q := &Queue{}

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
	q := &Queue{}
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
	q := &Queue{}
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

func TestGetPutDisposed(t *testing.T) {
	q := &Queue{}

	q.Dispose()

	_, err := q.Get(1)
	assert.IsType(t, DisposedError{}, err)

	err = q.Put(`a`)
	assert.IsType(t, DisposedError{}, err)
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
	q := &Queue{}
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
	q := &Queue{}
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
	q := &Queue{}
	q.Dispose()
	result, err := q.TakeUntil(func(item interface{}) bool {
		return true
	})

	assert.Nil(t, result)
	assert.IsType(t, DisposedError{}, err)
}
