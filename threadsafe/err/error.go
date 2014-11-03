package err

import "sync"

// Error is a struct that holds an error and allows this error
// to be set and retrieved in a threadsafe manner.
type Error struct {
	lock sync.RWMutex
	err  error
}

// Set will set the error of this structure to the provided
// value.
func (e *Error) Set(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.err = err
}

// Get will return any error associated with this structure.
func (e *Error) Get() error {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.err
}

// New is a constructor to generate a new error object
// that can be set and retrieved in a threadsafe manner.
func New() *Error {
	return &Error{}
}
