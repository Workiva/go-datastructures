/*
Copyright 2015 Workiva, LLC

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

	for i := 0; i < 1000; i++ {
		assert.Nil(b.Put("foo bar baz"))
	}

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

	// This delta is normally 1-3 ms but running tests in CI with -race causes
	// this to run much slower. For now, just bump up the threshold.
	assert.InDelta(200, time.Since(before).Seconds()*1000, 50)
	assert.True(len(batch) > 0)
	assert.Nil(err)
}

func TestFlush(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 10, 10, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)
	b.Put("a")
	wait := make(chan bool)
	go func() {
		batch, err := b.Get()
		assert.Equal([]interface{}{"a"}, batch)
		assert.Nil(err)
		wait <- true
	}()

	b.Flush()
	<-wait
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
	b, err := New(1, 2, 100000, 2, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)
	b.Put("a")
	b.Put("b")
	b.Put("c")

	batch1, err := b.Get()
	assert.Equal([]interface{}{"a", "b"}, batch1)
	assert.Nil(err)

	batch2, err := b.Get()
	assert.Equal([]interface{}{"c"}, batch2)
	assert.Nil(err)

	b.Put("d")
	b.Put("e")
	b.Put("f")

	b.Dispose()

	_, err = b.Get()
	assert.Equal(ErrDisposed, err)

	assert.Equal(ErrDisposed, b.Put("j"))
	assert.Equal(ErrDisposed, b.Flush())

}

func TestIsDisposed(t *testing.T) {
	assert := assert.New(t)
	b, err := New(0, 10, 10, 10, func(str interface{}) uint {
		return uint(len(str.(string)))
	})
	assert.Nil(err)
	assert.False(b.IsDisposed())
	b.Dispose()
	assert.True(b.IsDisposed())
}
