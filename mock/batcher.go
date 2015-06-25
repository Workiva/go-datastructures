package mock

import (
	"github.com/stretchr/testify/mock"
)

type MockBatcher struct {
	mock.Mock
	PutChan chan bool
}

func (m *MockBatcher) Put(items interface{}) error {
	args := m.Called(items)
	if m.PutChan != nil {
		m.PutChan <- true
	}
	return args.Error(0)
}

func (m *MockBatcher) Get() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockBatcher) Dispose() {
	m.Called()
}
