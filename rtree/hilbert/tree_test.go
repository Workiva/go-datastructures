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

package hilbert

import (
	"log"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/rtree"
)

func getConsoleLogger() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}

func (n *node) print(log *log.Logger) {
	log.Printf(`NODE: %+v, MBR: %+v, %p`, n, n.mbr, n)
	if n.isLeaf {
		for i, wh := range n.nodes.list {
			xlow, ylow := wh.LowerLeft()
			xhigh, yhigh := wh.UpperRight()
			log.Printf(`KEY: %+v, XLOW: %+v, YLOW: %+v, XHIGH: %+v, YHIGH: %+v`, n.keys.list[i], xlow, ylow, xhigh, yhigh)
		}
	} else {
		for _, wh := range n.nodes.list {
			wh.(*node).print(log)
		}
	}
}

func (t *tree) print(log *log.Logger) {
	log.Println(`PRINTING TREE`)
	if t.root == nil {
		log.Println(`EMPTY TREE.`)
		return
	}

	t.root.print(log)
}

func constructMockPoints(num int) rtree.Rectangles {
	rects := make(rtree.Rectangles, 0, num)
	for i := int32(0); i < int32(num); i++ {
		rects = append(rects, newMockRectangle(i, i, i, i))
	}
	return rects
}

func constructRandomMockPoints(num int) rtree.Rectangles {
	rects := make(rtree.Rectangles, 0, num)
	for i := 0; i < num; i++ {
		r := rand.Int31()
		rects = append(rects, newMockRectangle(r, r, r, r))
	}

	return rects
}

func TestSimpleInsert(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	tree := newTree(3, 3)

	tree.Insert(r1)
	assert.Equal(t, uint64(1), tree.Len())

	q := newMockRectangle(5, 5, 15, 15)
	result := tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1}, result)
}

func TestTwoInsert(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	tree := newTree(3, 3)

	tree.Insert(r1, r2)
	assert.Equal(t, uint64(2), tree.Len())

	q := newMockRectangle(0, 0, 20, 20)
	result := tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1, r2}, result)

	q = newMockRectangle(0, 0, 4, 4)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1}, result)

	q = newMockRectangle(11, 11, 20, 20)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r2}, result)
}

func TestInsertCausesRootSplitOddAry(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	r3 := newMockRectangle(10, 10, 20, 20)
	tree := newTree(3, 3)

	tree.Insert(r1, r2, r3)
	assert.Equal(t, uint64(3), tree.Len())

	q := newMockRectangle(0, 0, 20, 20)
	result := tree.Search(q)
	assert.Contains(t, result, r1)
	assert.Contains(t, result, r2)
	assert.Contains(t, result, r3)
}

func TestInsertCausesRootSplitEvenAry(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	r3 := newMockRectangle(10, 10, 20, 20)
	r4 := newMockRectangle(15, 15, 25, 25)
	tree := newTree(4, 4)

	tree.Insert(r1, r2, r3, r4)
	assert.Equal(t, uint64(4), tree.Len())

	q := newMockRectangle(0, 0, 25, 25)
	result := tree.Search(q)
	assert.Contains(t, result, r1)
	assert.Contains(t, result, r2)
	assert.Contains(t, result, r3)
	assert.Contains(t, result, r4)
}

func TestQueryWithLine(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	tree := newTree(3, 3)
	tree.Insert(r1, r2)

	// vertical line at x=5
	q := newMockRectangle(5, 0, 5, 10)
	result := tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1, r2}, result)

	// horizontal line at y=5
	q = newMockRectangle(0, 5, 10, 5)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1, r2}, result)

	// vertical line at x=15
	q = newMockRectangle(15, 0, 15, 20)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r2}, result)

	// horizontal line at y=15
	q = newMockRectangle(0, 15, 20, 15)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r2}, result)

	// vertical line on the y-axis
	q = newMockRectangle(0, 0, 0, 10)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1}, result)

	// horizontal line on the x-axis
	q = newMockRectangle(0, 0, 10, 0)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1}, result)

	// vertical line at x=20
	q = newMockRectangle(20, 0, 20, 20)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{}, result)

	// horizontal line at y=20
	q = newMockRectangle(0, 20, 20, 20)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{}, result)
}

func TestQueryForPoint(t *testing.T) {
	r1 := newMockRectangle(5, 5, 5, 5)     // (5, 5)
	r2 := newMockRectangle(10, 10, 10, 10) // (10, 10)
	tree := newTree(3, 3)
	tree.Insert(r1, r2)

	q := newMockRectangle(0, 0, 5, 5)
	result := tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r1}, result)

	q = newMockRectangle(0, 0, 20, 20)
	result = tree.Search(q)
	assert.Contains(t, result, r1)
	assert.Contains(t, result, r2)

	q = newMockRectangle(6, 6, 20, 20)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{r2}, result)

	q = newMockRectangle(20, 20, 30, 30)
	result = tree.Search(q)
	assert.Equal(t, rtree.Rectangles{}, result)
}

func TestMultipleInsertsCauseInternalSplitOddAry(t *testing.T) {
	points := constructMockPoints(100)
	tree := newTree(3, 3)

	tree.Insert(points...)

	assert.Equal(t, uint64(len(points)), tree.Len())

	q := newMockRectangle(0, 0, int32(len(points)), int32(len(points)))
	result := tree.Search(q)
	succeeded := true
	for _, p := range points {
		if !assert.Contains(t, result, p) {
			succeeded = false
		}
	}

	if !succeeded {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertsCauseInternalSplitOddAryRandomPoints(t *testing.T) {
	points := constructRandomMockPoints(100)
	tree := newTree(3, 3)

	tree.Insert(points...)

	assert.Equal(t, uint64(len(points)), tree.Len())

	q := newMockRectangle(0, 0, math.MaxInt32, math.MaxInt32)
	result := tree.Search(q)
	succeeded := true
	for _, p := range points {
		if !assert.Contains(t, result, p) {
			succeeded = false
		}
	}

	if !succeeded {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertsCauseInternalSplitEvenAry(t *testing.T) {
	points := constructMockPoints(100)
	tree := newTree(4, 4)

	tree.Insert(points...)

	assert.Equal(t, uint64(len(points)), tree.Len())

	q := newMockRectangle(0, 0, math.MaxInt32, math.MaxInt32)
	result := tree.Search(q)
	succeeded := true
	for _, p := range points {
		if !assert.Contains(t, result, p) {
			succeeded = false
		}
	}

	if !succeeded {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertsCauseInternalSplitEvenAryRandomOrder(t *testing.T) {
	points := constructRandomMockPoints(100)
	tree := newTree(4, 4)

	tree.Insert(points...)

	assert.Equal(t, uint64(len(points)), tree.Len())

	q := newMockRectangle(0, 0, math.MaxInt32, math.MaxInt32)
	result := tree.Search(q)
	succeeded := true
	for _, p := range points {
		if !assert.Contains(t, result, p) {
			succeeded = false
		}
	}

	if !succeeded {
		tree.print(getConsoleLogger())
	}
}

func BenchmarkBulkAddPoints(b *testing.B) {
	numItems := 1000
	points := constructMockPoints(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newTree(8, 8)
		tree.Insert(points...)
	}
}

func BenchmarkBulkUpdatePoints(b *testing.B) {
	numItems := 1000
	points := constructMockPoints(numItems)
	tree := newTree(8, 8)
	tree.Insert(points...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(points...)
	}
}

func BenchmarkPointInsertion(b *testing.B) {
	numItems := b.N
	points := constructMockPoints(numItems)
	tree := newTree(8, 8)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(points[i%numItems])
	}
}

func BenchmarkQueryPoints(b *testing.B) {
	numItems := b.N
	points := constructMockPoints(numItems)
	tree := newTree(8, 8)
	tree.Insert(points...)

	b.ResetTimer()

	for i := int32(0); i < int32(b.N); i++ {
		tree.Search(newMockRectangle(i, i, i+10, i+10))
	}
}
