package futures

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitOnGetResult(t *testing.T) {
	completer := make(chan interface{})
	f := New(completer, time.Duration(30*time.Minute))
	var result interface{}
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		result, err = f.GetResult()
		wg.Done()
	}()

	completer <- `test`
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, `test`, result)

	// ensure we don't get paused on the next iteration.
	result, err = f.GetResult()

	assert.Equal(t, `test`, result)
	assert.Nil(t, err)
}

func TestTimeout(t *testing.T) {
	completer := make(chan interface{})
	f := New(completer, time.Duration(0))

	result, err := f.GetResult()

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func BenchmarkFuture(b *testing.B) {
	completer := make(chan interface{})
	timeout := time.Duration(30 * time.Minute)
	var wg sync.WaitGroup

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		f := New(completer, timeout)
		go func() {
			f.GetResult()
			wg.Done()
		}()

		completer <- `test`
		wg.Wait()
	}
}
