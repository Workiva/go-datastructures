package rangetree

type mockEntry struct {
	id         uint64
	dimensions []int64
}

func (me *mockEntry) ID() uint64 {
	return me.id
}

func (me *mockEntry) ValueAtDimension(dimension uint64) int64 {
	return me.dimensions[dimension-1]
}

func constructMockEntry(id uint64, values ...int64) *mockEntry {
	return &mockEntry{
		id:         id,
		dimensions: values,
	}
}

type dimension struct {
	low, high int64
}

type mockInterval struct {
	dimensions []dimension
}

func (mi *mockInterval) LowAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].low
}

func (mi *mockInterval) HighAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].high
}

func constructMockInterval(dimensions ...dimension) *mockInterval {
	return &mockInterval{
		dimensions: dimensions,
	}
}
