/*
Copyright 2015 Workiva, LLC

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

/*
Package list provides list implementations. Currently, this includes a
persistent, immutable linked list.
*/
package list

import "errors"

var (
	// Empty is an empty PersistentList.
	Empty PersistentList = &emptyList{}

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

	// Length returns the number of items in the list.
	Length() uint

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

	// Find applies the predicate function to the list and returns the first
	// item which matches.
	Find(func(interface{}) bool) (interface{}, bool)

	// FindIndex applies the predicate function to the list and returns the
	// index of the first item which matches or -1 if there is no match.
	FindIndex(func(interface{}) bool) int

	// Map applies the function to each entry in the list and returns the
	// resulting slice.
	Map(func(interface{}) interface{}) []interface{}
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

// Length returns the number of items in the list.
func (e *emptyList) Length() uint {
	return 0
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

// Find applies the predicate function to the list and returns the first item
// which matches.
func (e *emptyList) Find(func(interface{}) bool) (interface{}, bool) {
	return nil, false
}

// FindIndex applies the predicate function to the list and returns the index
// of the first item which matches or -1 if there is no match.
func (e *emptyList) FindIndex(func(interface{}) bool) int {
	return -1
}

// Map applies the function to each entry in the list and returns the resulting
// slice.
func (e *emptyList) Map(func(interface{}) interface{}) []interface{} {
	return nil
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

// Length returns the number of items in the list.
func (l *list) Length() uint {
	curr := l
	length := uint(0)
	for {
		length += 1
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return length
		}
		curr = tail.(*list)
	}
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

// Find applies the predicate function to the list and returns the first item
// which matches.
func (l *list) Find(pred func(interface{}) bool) (interface{}, bool) {
	if pred(l.head) {
		return l.head, true
	}
	return l.tail.Find(pred)
}

// FindIndex applies the predicate function to the list and returns the index
// of the first item which matches or -1 if there is no match.
func (l *list) FindIndex(pred func(interface{}) bool) int {
	curr := l
	idx := 0
	for {
		if pred(curr.head) {
			return idx
		}
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return -1
		}
		curr = tail.(*list)
		idx += 1
	}
}

// Map applies the function to each entry in the list and returns the resulting
// slice.
func (l *list) Map(f func(interface{}) interface{}) []interface{} {
	return append(l.tail.Map(f), f(l.head))
}
