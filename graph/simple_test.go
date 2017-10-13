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

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV(t *testing.T) {
	assert := assert.New(t)
	sgraph := NewSimpleGraph()
	assert.Equal(0, sgraph.V())

	sgraph.AddEdge("A", "B")
	assert.Equal(2, sgraph.V())

	sgraph.AddEdge("B", "C")
	assert.Equal(3, sgraph.V())

	sgraph.AddEdge("A", "C")
	assert.Equal(3, sgraph.V())

	// Parallel edges not allowed
	sgraph.AddEdge("A", "C")
	assert.Equal(3, sgraph.V())
	sgraph.AddEdge("C", "A")
	assert.Equal(3, sgraph.V())

	// Self loops not allowed
	sgraph.AddEdge("C", "C")
	assert.Equal(3, sgraph.V())
	sgraph.AddEdge("D", "D")
	assert.Equal(3, sgraph.V())
}

func TestE(t *testing.T) {
	assert := assert.New(t)
	sgraph := NewSimpleGraph()

	assert.Equal(0, sgraph.E())

	sgraph.AddEdge("A", "B")
	assert.Equal(1, sgraph.E())

	sgraph.AddEdge("B", "C")
	assert.Equal(2, sgraph.E())

	sgraph.AddEdge("A", "C")
	assert.Equal(3, sgraph.E())

	// Parallel edges not allowed
	sgraph.AddEdge("A", "C")
	assert.Equal(3, sgraph.E())
	sgraph.AddEdge("C", "A")
	assert.Equal(3, sgraph.E())

	// Self loops not allowed so no edges added
	sgraph.AddEdge("C", "C")
	assert.Equal(3, sgraph.E())
	sgraph.AddEdge("D", "D")
	assert.Equal(3, sgraph.E())
}

func TestDegree(t *testing.T) {
	assert := assert.New(t)
	sgraph := NewSimpleGraph()

	// No edges added so degree is 0
	v, err := sgraph.Degree("A")
	assert.Zero(v)
	assert.Error(err)

	// One edge added
	sgraph.AddEdge("A", "B")
	v, err = sgraph.Degree("A")
	assert.Equal(1, v)
	assert.Nil(err)

	// Self loops are not allowed
	sgraph.AddEdge("A", "A")
	v, err = sgraph.Degree("A")
	assert.Equal(1, v)
	assert.Nil(err)

	// Parallel edges are not allowed
	sgraph.AddEdge("A", "B")
	v, err = sgraph.Degree("A")
	assert.Equal(1, v)
	assert.Nil(err)
	sgraph.AddEdge("B", "A")
	v, err = sgraph.Degree("A")
	assert.Equal(1, v)
	assert.Nil(err)

	v, err = sgraph.Degree("B")
	assert.Equal(1, v)
	assert.Nil(err)

	sgraph.AddEdge("C", "D")
	sgraph.AddEdge("A", "C")
	sgraph.AddEdge("E", "F")
	sgraph.AddEdge("E", "G")
	sgraph.AddEdge("H", "G")

	v, err = sgraph.Degree("A")
	assert.Equal(2, v)
	assert.Nil(err)

	v, err = sgraph.Degree("B")
	assert.Equal(1, v)
	assert.Nil(err)

	v, err = sgraph.Degree("C")
	assert.Equal(2, v)
	assert.Nil(err)

	v, err = sgraph.Degree("D")
	assert.Equal(1, v)
	assert.Nil(err)

	v, err = sgraph.Degree("E")
	assert.Equal(2, v)
	assert.Nil(err)

	v, err = sgraph.Degree("G")
	assert.Equal(2, v)
	assert.Nil(err)
}

func TestAddEdge(t *testing.T) {
	assert := assert.New(t)
	sgraph := NewSimpleGraph()

	err := sgraph.AddEdge("A", "B")
	assert.Nil(err)

	err = sgraph.AddEdge("A", "B")
	assert.Error(err)

	err = sgraph.AddEdge("B", "A")
	assert.Error(err)

	err = sgraph.AddEdge("A", "A")
	assert.Error(err)

	err = sgraph.AddEdge("C", "C")
	assert.Error(err)

	err = sgraph.AddEdge("B", "C")
	assert.Nil(err)

}

func TestAdj(t *testing.T) {
	assert := assert.New(t)
	sgraph := NewSimpleGraph()

	v, err := sgraph.Adj("A")
	assert.Zero(v)
	assert.Error(err)

	// Self loops not allowed
	sgraph.AddEdge("A", "A")
	v, err = sgraph.Adj("A")
	assert.Zero(v)
	assert.Error(err)

	sgraph.AddEdge("A", "B")
	v, err = sgraph.Adj("A")
	assert.Equal(1, len(v))
	assert.Nil(err)
	assert.Equal("B", v[0])

	v, err = sgraph.Adj("B")
	assert.Equal(1, len(v))
	assert.Nil(err)
	assert.Equal("A", v[0])

	// Parallel Edges not allowed
	sgraph.AddEdge("A", "B")
	sgraph.AddEdge("B", "A")
	v, err = sgraph.Adj("B")
	assert.Equal(1, len(v))
	assert.Nil(err)
	assert.Equal("A", v[0])

	sgraph.AddEdge("C", "D")
	sgraph.AddEdge("A", "C")
	sgraph.AddEdge("E", "F")
	sgraph.AddEdge("E", "G")
	sgraph.AddEdge("H", "G")

	v, err = sgraph.Adj("A")
	assert.Equal(2, len(v))
	assert.Nil(err)
	assert.Contains(v, "B")
	assert.Contains(v, "C")
	assert.NotContains(v, "A")
	assert.NotContains(v, "D")

	v, err = sgraph.Adj("B")
	assert.Equal(1, len(v))
	assert.Nil(err)
	assert.Contains(v, "A")
	assert.NotContains(v, "B")
	assert.NotContains(v, "C")
	assert.NotContains(v, "D")

	v, err = sgraph.Adj("C")
	assert.Equal(2, len(v))
	assert.Nil(err)
	assert.Contains(v, "A")
	assert.Contains(v, "D")
	assert.NotContains(v, "B")
	assert.NotContains(v, "C")

	v, err = sgraph.Adj("E")
	assert.Equal(2, len(v))
	assert.Nil(err)
	assert.Contains(v, "F")
	assert.Contains(v, "G")
	assert.NotContains(v, "A")

	v, err = sgraph.Adj("G")
	assert.Equal(2, len(v))
	assert.Nil(err)
	assert.Contains(v, "E")
	assert.Contains(v, "H")
	assert.NotContains(v, "A")
}
