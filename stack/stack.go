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

/* Package stack implements a simple linked-list based stack */
package stack

import "errors"

var (
	ErrEmptyStack = errors.New("stack is empty")
)

// Stack is an immutable stack
type Stack interface {
	// Top returns the top-most item of the stack.
	// If the stack is empty, the bool is set to false.
	Top() (interface{}, bool)

	// Pop returns the top-most item of the stack and removes it.
	// The error is set to ErrEmptyStack should the stack be empty.
	Pop() (interface{}, error)

	// Drop drops the top-most item of the stack.
	// Error is set to ErrEmptyStack should the stack be empty.
	Drop() error

	// Push pushes an item onto the stack.
	Push(interface{})

	// PopWhile creates a channel of interfaces and pops items from the stack
	// as long as the predicate passed holds or the stack is emptied.
	PopWhile(func(interface{}) bool) []interface{}

	// DropWhile drops items from the stack as long as the predicate passed holds or the stack is emptied.
	DropWhile(func(interface{}) bool)

	// IsEmpty returns whether the stack is empty.
	IsEmpty() bool

	// Size returns the amount of items in the stack.
	Size() uint

	// Clear empties the stack
	Clear()
}

type stack struct {
	size uint
	top  *item
}

type item struct {
	item interface{}
	next *item
}

// Top returns the top-most item of the stack.
// If the stack is empty, the bool is set to false.
func (s *stack) Top() (interface{}, bool) {
	if s.top == nil {
		return nil, false
	}
	return s.top.item, true
}

// Pop returns the top-most item of the stack and removes it.
// The error is set to ErrEmptyStack should the stack be empty.
func (s *stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, ErrEmptyStack
	}

	s.size--
	top := s.top
	s.top = s.top.next
	return top.item, nil
}

// Drop drops the top-most item of the stack.
// Error is set to ErrEmptyStack should the stack be empty.
func (s *stack) Drop() error {
	if s.IsEmpty() {
		return ErrEmptyStack
	}

	s.size--
	top := s.top
	ntop := top.next
	s.top = ntop
	return nil
}

// Push pushes an item onto the stack.
func (s *stack) Push(it interface{}) {
	s.size++
	s.top = &item{it, s.top}
}

// PopWhile creates a channel of interfaces and pops items from the stack
// as long as the predicate passed holds or the stack is emptied.
func (s *stack) PopWhile(pred func(interface{}) bool) []interface{} {
	its := make([]interface{}, 0)
	for !s.IsEmpty() {
		// We are sure this cannot return an error
		it, _ := s.Top()
		if pred(it) {
			s.Pop()
			its = append(its, it)
			continue
		}
		break
	}
	return its
}

// DropWhile drops items from the stack as long as the predicate passed holds or the stack is emptied.
func (s *stack) DropWhile(pred func(interface{}) bool) {
	for !s.IsEmpty() {
		// We are sure this cannot return an error
		it, _ := s.Top()
		if pred(it) {
			s.Pop()
			continue
		}
		return
	}
}

// IsEmpty returns whether the stack is empty.
func (s *stack) IsEmpty() bool {
	return s.size == 0
}

// Size returns the amount of items in the stack.
func (s *stack) Size() uint {
	return s.size
}

// Clear empties the stack
func (s *stack) Clear() {
	s.size = 0
	s.top = nil
}

// Empty returns a new empty stack
func Empty() Stack {
	return &stack{0, nil}
}
