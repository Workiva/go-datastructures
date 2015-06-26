package batcher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoCalculateBytes(t *testing.T) {
	_, err := New(0, 0, 100, 5, nil)
	assert.Error(t, err)
}

func TestMaxItems(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 100, 100000, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)

	go func() {
		for i := 0; i < 1000; i++ {
			assert.Nil(b.Put("foo bar baz"))
		}
	}()

	batch, err := b.Get()
	assert.Len(batch, 100)
	assert.Nil(err)
}

func TestMaxBytes(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 10000, 100, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)

	go func() {
		for i := 0; i < 1000; i++ {
			b.Put("a")
		}
	}()

	batch, err := b.Get()
	assert.Len(batch, 100)
	assert.Nil(err)
}

func TestMaxTime(t *testing.T) {
	assert := assert.New(t)
	b, err := New(time.Millisecond*200, 100000, 100000, 10,
		func(str interface{}) uint {
			return uint(len(str.(string)))
		},
	)
	assert.Nil(err)

	go func() {
		for i := 0; i < 10000; i++ {
			b.Put("a")
			time.Sleep(time.Millisecond)
		}
	}()

	before := time.Now()
	batch, err := b.Get()
	assert.InDelta(200, time.Since(before).Seconds()*1000, 5)
	assert.True(len(batch) > 0)
	assert.Nil(err)
}

func TestMultiConsumer(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 100, 100000, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			batch, err := b.Get()
			assert.Len(batch, 100)
			assert.Nil(err)
			wg.Done()
		}()
	}

	go func() {
		for i := 0; i < 500; i++ {
			b.Put("a")
		}
	}()

	wg.Wait()
}

func TestDispose(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 100000, 100000, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)
	b.Put("a")
	wait := make(chan bool)
	go func() {
		_, err := b.Get()
		assert.Equal(ErrDisposed, err)
		wait <- true
	}()

	b.Dispose()

	assert.Equal(ErrDisposed, b.Put("a"))

	<-wait
}
