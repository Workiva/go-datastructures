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

/*
Package err implements a threadsafe error interface.  In my places,
I found myself needing a lock to protect writing to a common error interface
from multiple go routines (channels are great but slow).  This just makes
that process more convenient.
*/

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
