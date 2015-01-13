package xfast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	iter := &Iterator{
		first: true,
	}

	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e1 := newMockEntry(5)
	n1 := newNode(nil, e1)
	iter = &Iterator{
		first: true,
		n:     n1,
	}

	assert.True(t, iter.Next())
	assert.Equal(t, e1, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e2 := newMockEntry(10)
	n2 := newNode(nil, e2)
	n1.children[1] = n2

	iter = &Iterator{
		first: true,
		n:     n1,
	}

	assert.True(t, iter.Next())
	assert.True(t, iter.Next())
	assert.Equal(t, e2, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())
}
