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

package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	s := Int64Slice{3, 6, 1, 0, -1}
	s.Sort()

	assert.Equal(t, Int64Slice{-1, 0, 1, 3, 6}, s)
}

func TestSearch(t *testing.T) {
	s := Int64Slice{1, 3, 6}

	assert.Equal(t, 1, s.Search(3))
	assert.Equal(t, 1, s.Search(2))
	assert.Equal(t, 3, s.Search(7))
}

func TestExists(t *testing.T) {
	s := Int64Slice{1, 3, 6}

	assert.True(t, s.Exists(3))
	assert.False(t, s.Exists(4))
}

func TestInsert(t *testing.T) {
	s := Int64Slice{1, 3, 6}
	s = s.Insert(2)
	assert.Equal(t, Int64Slice{1, 2, 3, 6}, s)

	s = s.Insert(7)
	assert.Equal(t, Int64Slice{1, 2, 3, 6, 7}, s)
}
