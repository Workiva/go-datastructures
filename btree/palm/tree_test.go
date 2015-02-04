package palm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInsert(t *testing.T) {
	tree := newTree(16)
	m1 := mockKey(1)

	result := tree.Insert(m1)
	assert.Equal(t, Keys{nil}, result)
	assert.Equal(t, Keys{m1}, tree.Get(m1))
	assert.Equal(t, uint64(1), tree.Len())
}
