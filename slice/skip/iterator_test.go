/*
Copyright 2014 Workiva, LLC

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

package skip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	e1 := newMockEntry(1)
	n1 := newNode(e1, 8)
	iter := &iterator{
		n:     n1,
		first: true,
	}

	assert.True(t, iter.Next())
	assert.Equal(t, e1, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e2 := newMockEntry(2)
	n2 := newNode(e2, 8)
	n1.forward[0] = n2

	iter = &iterator{
		n:     n1,
		first: true,
	}

	assert.True(t, iter.Next())
	assert.Equal(t, e1, iter.Value())
	assert.True(t, iter.Next())
	assert.Equal(t, e2, iter.Value())
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	iter = nilIterator()
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())
}
