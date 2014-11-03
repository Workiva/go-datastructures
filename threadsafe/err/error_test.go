package err

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSetError(t *testing.T) {
	e := New()
	assert.Nil(t, e.Get())

	err := fmt.Errorf(`test`)
	e.Set(err)

	assert.Equal(t, err, e.Get())
}
