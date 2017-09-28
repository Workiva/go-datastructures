package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyStack(t *testing.T) {
	assert := assert.New(t)
	s := Empty()
	top, ok := s.Top()
	assert.Nil(top)
	assert.False(ok)

	top, err := s.Pop()
	assert.Nil(top)
	assert.Equal(err, ErrEmptyStack)

	assert.True(s.IsEmpty())
}

func TestPushPop(t *testing.T) {
	assert := assert.New(t)
	s := Empty()

	// stack [10]
	s.Push(10)
	assert.False(s.IsEmpty())
	top, ok := s.Top()
	assert.True(ok)
	assert.Equal(top, 10)
	assert.Equal(uint(1), s.Size())

	s.Push(3)
	// stack [3 10]
	assert.Equal(uint(2), s.Size())
	top, err := s.Pop()
	assert.Nil(err)
	assert.Equal(top, 3)
	assert.Equal(uint(1), s.Size())
}

func TestPopDropWhile(t *testing.T) {
	assert := assert.New(t)
	s := Empty()
	for i := 0; i < 11; i++ {
		s.Push(i * i)
	}
	assert.Equal(uint(11), s.Size())

	pred := func(it interface{}) bool {
		return it.(int) >= 64
	}

	its := s.PopWhile(pred)

	for _, it := range its {
		assert.True(pred(it))
	}

	assert.Equal(uint(8), s.Size())

	pred = func(it interface{}) bool {
		return s.Size() > 3
	}

	s.DropWhile(pred)

	assert.Equal(uint(3), s.Size())
}

func TestClearStack(t *testing.T) {
	assert := assert.New(t)

	s := Empty()
	s.Push("a")
	s.Push("b")
	s.Push("c")
	s.Push("d")

	assert.Equal(uint(4), s.Size())
	top, ok := s.Top()
	assert.True(ok)
	assert.Equal(top, "d")

	s.Clear()
	assert.True(s.IsEmpty())
}
