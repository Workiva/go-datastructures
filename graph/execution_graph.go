package graph

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/queue"
)

type ExecutionGraph struct {
	size      int64
	toApply   []Nodes
	circulars Nodes
}

// Size returns the number of nodes in the graph.
func (graph *ExecutionGraph) Size() int64 {
	return graph.size
}

// RecursivelyApply calls the given function across the graph in dependency order.  If
// there are only circulars, a circular node is randomly chosen and all
// circulars will be called.  If the supplied function ever returns false,
// all further execution is halted.
func (graph *ExecutionGraph) RecursivelyApply(fn func(node INode) bool) {
	for _, nodes := range graph.combined() {
		for _, node := range nodes {
			result := fn(node)
			if !result {
				return
			}
		}
	}
}

func (graph *ExecutionGraph) worker(q *queue.Queue, fn func(node INode) bool,
	wg *sync.WaitGroup, done *uint64, todo uint64) {

	for {
		items, err := q.Get(1)
		if err != nil {
			break
		}

		node := items[0].(INode)

		fn(node)

		if atomic.AddUint64(done, 1) == todo {
			wg.Done()
			break
		}
	}
}

func (graph *ExecutionGraph) combined() []Nodes {
	toApply := graph.toApply
	toApply = append(toApply, graph.circulars)
	return toApply
}

// ParallelRecursivelyApply operates similarly to RecursivelyApply but does so
// in parallel if possible.
func (graph *ExecutionGraph) ParallelRecursivelyApply(fn func(node INode) bool) {
	if runtime.NumCPU() < 2 || graph.size < 20 { //20 is just some arbitrary number
		graph.RecursivelyApply(fn)
		return
	}

	var wg sync.WaitGroup

	for _, nodes := range graph.combined() {
		if int64(len(nodes)) < numNodesBeforeSplit {
			for _, node := range nodes {
				fn(node)
			}
		} else {
			q := queue.New(int64(len(nodes)))

			for _, node := range nodes {
				q.Put(node)
			}

			wg.Add(1)
			todo, done := uint64(q.Len()), uint64(0)
			for i := 0; i < runtime.NumCPU(); i++ {
				go graph.worker(q, fn, &wg, &done, todo)
			}

			wg.Wait()
			q.Dispose()
		}
	}
}

func newExecutionGraph() *ExecutionGraph {
	return &ExecutionGraph{}
}
