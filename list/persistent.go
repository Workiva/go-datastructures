package list

import "errors"

var (
	// Empty is an empty PersistentList.
	Empty = &emptyList{}

	// ErrEmptyList is returned when an invalid operation is performed on an
	// empty list.
	ErrEmptyList = errors.New("Empty list")
)

// PersistentList is an immutable, persistent linked list.
type PersistentList interface {
	// Head returns the head of the list. The bool will be false if the list is
	// empty.
	Head() (interface{}, bool)

	// Tail returns the tail of the list. The bool will be false if the list is
	// empty.
	Tail() (PersistentList, bool)

	// IsEmpty indicates if the list is empty.
	IsEmpty() bool

	// Add will add the item to the list, returning the new list.
	Add(head interface{}) PersistentList

	// Insert will insert the item at the given position, returning the new
	// list or an error if the position is invalid.
	Insert(val interface{}, pos uint) (PersistentList, error)

	// Get returns the item at the given position or an error if the position
	// is invalid.
	Get(pos uint) (interface{}, bool)

	// Remove will remove the item at the given position, returning the new
	// list or an error if the position is invalid.
	Remove(pos uint) (PersistentList, error)
}

type emptyList struct{}

// Head returns the head of the list. The bool will be false if the list is
// empty.
func (e *emptyList) Head() (interface{}, bool) {
	return nil, false
}

// Tail returns the tail of the list. The bool will be false if the list is
// empty.
func (e *emptyList) Tail() (PersistentList, bool) {
	return nil, false
}

// IsEmpty indicates if the list is empty.
func (e *emptyList) IsEmpty() bool {
	return true
}

// Add will add the item to the list, returning the new list.
func (e *emptyList) Add(head interface{}) PersistentList {
	return &list{head, e}
}

// Insert will insert the item at the given position, returning the new list or
// an error if the position is invalid.
func (e *emptyList) Insert(val interface{}, pos uint) (PersistentList, error) {
	if pos == 0 {
		return e.Add(val), nil
	}
	return nil, ErrEmptyList
}

// Get returns the item at the given position or an error if the position is
// invalid.
func (e *emptyList) Get(pos uint) (interface{}, bool) {
	return nil, false
}

// Remove will remove the item at the given position, returning the new list or
// an error if the position is invalid.
func (e *emptyList) Remove(pos uint) (PersistentList, error) {
	return nil, ErrEmptyList
}

type list struct {
	head interface{}
	tail PersistentList
}

// Head returns the head of the list. The bool will be false if the list is
// empty.
func (l *list) Head() (interface{}, bool) {
	return l.head, true
}

// Tail returns the tail of the list. The bool will be false if the list is
// empty.
func (l *list) Tail() (PersistentList, bool) {
	return l.tail, true
}

// IsEmpty indicates if the list is empty.
func (l *list) IsEmpty() bool {
	return false
}

// Add will add the item to the list, returning the new list.
func (l *list) Add(head interface{}) PersistentList {
	return &list{head, l}
}

// Insert will insert the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Insert(val interface{}, pos uint) (PersistentList, error) {
	if pos == 0 {
		return l.Add(val), nil
	}
	nl, err := l.tail.Insert(val, pos-1)
	if err != nil {
		return nil, err
	}
	return nl.Add(l.head), nil
}

// Get returns the item at the given position or an error if the position is
// invalid.
func (l *list) Get(pos uint) (interface{}, bool) {
	if pos == 0 {
		return l.head, true
	}
	return l.tail.Get(pos - 1)
}

// Remove will remove the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Remove(pos uint) (PersistentList, error) {
	if pos == 0 {
		nl, _ := l.Tail()
		return nl, nil
	}

	nl, err := l.tail.Remove(pos - 1)
	if err != nil {
		return nil, err
	}
	return &list{l.head, nl}, nil
}
