package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/Workiva/go-datastructures/rangetree"
)

type RangeTree struct {
	mock.Mock
}

var _ rangetree.RangeTree = new(RangeTree)

func (m *RangeTree) Add(entries ...rangetree.Entry) rangetree.Entries {
	args := m.Called(entries)
	ifc := args.Get(0)
	if ifc == nil {
		return nil
	}

	return ifc.(rangetree.Entries)
}

func (m *RangeTree) Len() uint64 {
	return m.Called().Get(0).(uint64)
}

func (m *RangeTree) Delete(entries ...rangetree.Entry) rangetree.Entries {
	return m.Called(entries).Get(0).(rangetree.Entries)
}

func (m *RangeTree) Query(interval rangetree.Interval) rangetree.Entries {
	args := m.Called(interval)
	ifc := args.Get(0)
	if ifc == nil {
		return nil
	}

	return ifc.(rangetree.Entries)
}

func (m *RangeTree) InsertAtDimension(dimension uint64, index,
	number int64) (rangetree.Entries, rangetree.Entries) {

	args := m.Called(dimension, index, number)
	return args.Get(0).(rangetree.Entries), args.Get(1).(rangetree.Entries)
}

func (m *RangeTree) Apply(interval rangetree.Interval, fn func(rangetree.Entry) bool) {
	m.Called(interval, fn)
}
