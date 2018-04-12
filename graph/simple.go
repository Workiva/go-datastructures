/*
Copyright 2017 Julian Griggs

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

/*
Package graph provides graph implementations. Currently, this includes an
undirected simple graph.
*/
package graph

import (
	"errors"
	"sync"
)

var (
	// ErrVertexNotFound is returned when an operation is requested on a
	// non-existent vertex.
	ErrVertexNotFound = errors.New("vertex not found")

	// ErrSelfLoop is returned when an operation tries to create a disallowed
	// self loop.
	ErrSelfLoop = errors.New("self loops not permitted")

	// ErrParallelEdge is returned when an operation tries to create a
	// disallowed parallel edge.
	ErrParallelEdge = errors.New("parallel edges are not permitted")
)

// SimpleGraph is a mutable, non-persistent undirected graph.
// Parallel edges and self-loops are not permitted.
// Additional description: https://en.wikipedia.org/wiki/Graph_(discrete_mathematics)#Simple_graph
type SimpleGraph struct {
	mutex         sync.RWMutex
	adjacencyList map[interface{}]map[interface{}]struct{}
	v, e          int
}

// V returns the number of vertices in the SimpleGraph
func (g *SimpleGraph) V() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.v
}

// E returns the number of edges in the SimpleGraph
func (g *SimpleGraph) E() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.e
}

// AddEdge will create an edge between vertices v and w
func (g *SimpleGraph) AddEdge(v, w interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if v == w {
		return ErrSelfLoop
	}

	g.addVertex(v)
	g.addVertex(w)

	if _, ok := g.adjacencyList[v][w]; ok {
		return ErrParallelEdge
	}

	g.adjacencyList[v][w] = struct{}{}
	g.adjacencyList[w][v] = struct{}{}
	g.e++
	return nil
}

// Adj returns the list of all vertices connected to v
func (g *SimpleGraph) Adj(v interface{}) ([]interface{}, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	deg, err := g.Degree(v)
	if err != nil {
		return nil, ErrVertexNotFound
	}

	adj := make([]interface{}, deg)
	i := 0
	for key := range g.adjacencyList[v] {
		adj[i] = key
		i++
	}
	return adj, nil
}

// Degree returns the number of vertices connected to v
func (g *SimpleGraph) Degree(v interface{}) (int, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	val, ok := g.adjacencyList[v]
	if !ok {
		return 0, ErrVertexNotFound
	}
	return len(val), nil
}

func (g *SimpleGraph) addVertex(v interface{}) {
	mm, ok := g.adjacencyList[v]
	if !ok {
		mm = make(map[interface{}]struct{})
		g.adjacencyList[v] = mm
		g.v++
	}
}

// NewSimpleGraph creates and returns a SimpleGraph
func NewSimpleGraph() *SimpleGraph {
	return &SimpleGraph{
		adjacencyList: make(map[interface{}]map[interface{}]struct{}),
		v:             0,
		e:             0,
	}
}
