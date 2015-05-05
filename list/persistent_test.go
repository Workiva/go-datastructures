package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyList(t *testing.T) {
	assert := assert.New(t)
	head, ok := Empty.Head()
	assert.Nil(head)
	assert.False(ok)

	tail, ok := Empty.Tail()
	assert.Nil(tail)
	assert.False(ok)

	assert.True(Empty.IsEmpty())
}

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	l1 := Empty.Add(1)

	// l1: [1]
	assert.False(l1.IsEmpty())
	head, ok := l1.Head()
	assert.True(ok)
	assert.Equal(1, head)
	tail, ok := l1.Tail()
	assert.True(ok)
	assert.Equal(Empty, tail)

	l1 = l1.Add(2)

	// l1: [2, 1]
	head, ok = l1.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = l1.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	l2, err := l1.Insert("a", 1)
	assert.Nil(err)

	// l1: [2, 1]
	// l2: [2, "a", 1]
	head, ok = l1.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = l1.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	head, ok = l2.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = l2.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal("a", head)
	tail, ok = tail.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)
}

func TestInsertAndGet(t *testing.T) {
	assert := assert.New(t)
	_, err := Empty.Insert(1, 5)
	assert.Error(err)

	l, err := Empty.Insert(1, 0)
	assert.Nil(err)

	// [1]
	item, ok := l.Get(0)
	assert.True(ok)
	assert.Equal(1, item)

	l, err = l.Insert(2, 0)
	assert.Nil(err)

	// [2, 1]
	item, ok = l.Get(0)
	assert.True(ok)
	assert.Equal(2, item)
	item, ok = l.Get(1)
	assert.True(ok)
	assert.Equal(1, item)

	_, ok = l.Get(2)
	assert.False(ok)

	l, err = l.Insert("a", 3)
	assert.Nil(l)
	assert.Error(err)
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	l, err := Empty.Remove(0)
	assert.Nil(l)
	assert.Error(err)

	l = Empty.Add(1)
	l = l.Add(2)
	l = l.Add(3)

	// [3, 2, 1]
	l1, err := l.Remove(3)
	assert.Nil(l1)
	assert.Error(err)

	l2, err := l.Remove(0)

	// l: [3, 2, 1]
	// l2: [2, 1]
	assert.Nil(err)
	head, ok := l.Head()
	assert.True(ok)
	assert.Equal(3, head)
	tail, ok := l.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = tail.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	assert.Nil(err)
	head, ok = l2.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = l2.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	l2, err = l.Remove(1)

	// l: [3, 2, 1]
	// l2: [3, 1]
	assert.Nil(err)
	head, ok = l.Head()
	assert.True(ok)
	assert.Equal(3, head)
	tail, ok = l.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = tail.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	assert.Nil(err)
	head, ok = l2.Head()
	assert.True(ok)
	assert.Equal(3, head)
	tail, ok = l2.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	l2, err = l.Remove(2)

	// l: [3, 2, 1]
	// l2: [3, 2]
	assert.Nil(err)
	head, ok = l.Head()
	assert.True(ok)
	assert.Equal(3, head)
	tail, ok = l.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(2, head)
	tail, ok = tail.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(1, head)

	assert.Nil(err)
	head, ok = l2.Head()
	assert.True(ok)
	assert.Equal(3, head)
	tail, ok = l2.Tail()
	assert.True(ok)
	head, ok = tail.Head()
	assert.True(ok)
	assert.Equal(2, head)
}
