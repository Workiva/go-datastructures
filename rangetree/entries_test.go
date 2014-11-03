package rangetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisposeEntries(t *testing.T) {
	entries := NewEntries()
	entries = append(entries, constructMockEntry(0, 0))

	entries.Dispose()

	assert.Len(t, entries, 0)
}
