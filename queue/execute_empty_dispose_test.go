package queue_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Workiva/go-datastructures/queue"
)

func TestExecuteInParallelDisposesEmptyQueue(t *testing.T) {
	for _, hint := range []int64{0, 8} {
		t.Run(map[int64]string{0: "zero capacity", 8: "preallocated"}[hint], func(t *testing.T) {
			q := queue.New(hint)
			queue.ExecuteInParallel(q, func(interface{}) { t.Error("callback called on empty queue") })
			done := make(chan bool, 1)
			go func() {
				disposed := q.Disposed()
				if q.Put("new") != queue.ErrDisposed {
					disposed = false
				}
				if _, err := q.Get(1); err != queue.ErrDisposed {
					disposed = false
				}
				done <- disposed
			}()
			select {
			case ok := <-done:
				if !ok {
					t.Error("empty queue was not disposed")
				}
			case <-time.After(time.Second):
				t.Fatal("queue remained locked after empty execution")
			}
		})
	}
}

func TestExecuteInParallelNormalDisposalControls(t *testing.T) {
	var calls int32
	queue.ExecuteInParallel(nil, func(interface{}) { t.Fatal("nil callback") })
	q := queue.New(3)
	q.Put(1, 2, 3)
	queue.ExecuteInParallel(q, func(interface{}) { atomic.AddInt32(&calls, 1) })
	if atomic.LoadInt32(&calls) != 3 || !q.Disposed() {
		t.Fatal("nonempty execution changed")
	}
}
