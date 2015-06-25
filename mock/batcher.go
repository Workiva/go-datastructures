package mock

import (
	"github.com/stretchr/testify/mock"
)

type Batcher struct {
	mock.Mock
	PutChan chan bool
}

func (m *Batcher) Put(items interface{}) error {
	args := m.Called(items)
	if m.PutChan != nil {
		m.PutChan <- true
	}
	return args.Error(0)
}

func (m *Batcher) Get() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *Batcher) Dispose() {
	m.Called()
}
