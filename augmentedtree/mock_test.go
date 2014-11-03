package augmentedtree

import "fmt"

type dimension struct {
	low, high int64
}

type mockInterval struct {
	dimensions []*dimension
	id         uint64
}

func (mi *mockInterval) checkDimension(dimension uint64) {
	if dimension > uint64(len(mi.dimensions)) {
		panic(fmt.Sprintf(`Dimension: %d out of range.`, dimension))
	}
}

func (mi *mockInterval) LowAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].low
}

func (mi *mockInterval) HighAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].high
}

func (mi *mockInterval) OverlapsAtDimension(iv Interval, dimension uint64) bool {
	return mi.HighAtDimension(dimension) > iv.LowAtDimension(dimension) &&
		mi.LowAtDimension(dimension) < iv.HighAtDimension(dimension)
}

func (mi mockInterval) ID() uint64 {
	return mi.id
}

func constructSingleDimensionInterval(low, high int64, id uint64) *mockInterval {
	return &mockInterval{[]*dimension{&dimension{low: low, high: high}}, id}
}

func constructMultiDimensionInterval(id uint64, dimensions ...*dimension) *mockInterval {
	return &mockInterval{dimensions: dimensions, id: id}
}
