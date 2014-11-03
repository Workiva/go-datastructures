package augmentedtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalsDispose(t *testing.T) {
	intervals := intervalsPool.Get().(Intervals)
	intervals = append(intervals, constructSingleDimensionInterval(0, 1, 0))

	intervals.Dispose()

	assert.Len(t, intervals, 0)
}
