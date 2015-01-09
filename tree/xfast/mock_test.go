package xfast

import "github.com/stretchr/testify/mock"

type mockEntry struct {
	mock.Mock
}

func (me *mockEntry) Key() uint64 {
	args := me.Called()
	return args.Get(0).(uint64)
}

func newMockEntry(key uint64) *mockEntry {
	me := new(mockEntry)
	me.On(`Key`).Return(key)
	return me
}
