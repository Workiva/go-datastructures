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

func TestEntriesInsert(t *testing.T) {
	e1 := newMockEntry(10)
	e2 := newMockEntry(5)

	entries := Entries{e1}
	result := entries.insert(e2)
	assert.Equal(t, Entries{e2, e1}, entries)
	assert.Nil(t, result)

	e3 := newMockEntry(15)
	result = entries.insert(e3)
	assert.Equal(t, Entries{e2, e1, e3}, entries)
	assert.Nil(t, result)
}

func TestInsertOverwrite(t *testing.T) {
	e1 := newMockEntry(10)
	e2 := newMockEntry(10)

	entries := Entries{e1}
	result := entries.insert(e2)
	assert.Equal(t, e1, result)
}

func TestEntriesDelete(t *testing.T) {
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)

	entries := Entries{e1, e2}

	result := entries.delete(e2)
	assert.Equal(t, Entries{e1}, entries)
	assert.Equal(t, e2, result)

	result = entries.delete(e1)
	assert.Equal(t, Entries{}, entries)
	assert.Equal(t, e1, result)
}
