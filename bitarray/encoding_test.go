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

package bitarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparseBitArraySerialization(t *testing.T) {
	numItems := uint64(1280)
	input := newSparseBitArray()

	for i := uint64(0); i < numItems; i++ {
		if i%3 == 0 {
			input.SetBit(i)
		}
	}

	outBytes, err := input.Serialize()
	assert.Equal(t, err, nil)

	assert.Equal(t, len(outBytes), 337)
	assert.True(t, outBytes[0] == 'S')
	expected := []byte{83, 20, 0, 0, 0, 0, 0, 0, 0, 73}
	assert.Equal(t, expected, outBytes[:10])

	output := newSparseBitArray()
	err = output.Deserialize(outBytes)
	assert.Equal(t, err, nil)
	assert.True(t, input.Equals(output))
}

func TestBitArraySerialization(t *testing.T) {
	numItems := uint64(1280)
	input := newBitArray(numItems)

	for i := uint64(0); i < numItems; i++ {
		if i%3 == 0 {
			input.SetBit(i)
		}
	}

	outBytes, err := input.Serialize()
	assert.Equal(t, err, nil)

	// 1280 bits = 20 blocks = 160 bytes, plus lowest and highest at
	// 128 bits = 16 bytes plus 1 byte for the anyset param and the identifer
	assert.Equal(t, len(outBytes), 178)

	expected := []byte{66, 0, 0, 0, 0, 0, 0, 0, 0, 254}
	assert.Equal(t, expected, outBytes[:10])

	output := newBitArray(0)
	err = output.Deserialize(outBytes)
	assert.Equal(t, err, nil)
	assert.True(t, input.Equals(output))
}

func TestBitArrayMarshalUnmarshal(t *testing.T) {
	numItems := uint64(1280)
	input := newBitArray(numItems)

	for i := uint64(0); i < numItems; i++ {
		if i%3 == 0 {
			input.SetBit(i)
		}
	}

	outputBytes, err := Marshal(input)
	assert.Equal(t, err, nil)
	assert.Equal(t, outputBytes[0], byte('B'))
	assert.Equal(t, len(outputBytes), 178)

	output, err := Unmarshal(outputBytes)
	assert.Equal(t, err, nil)

	assert.True(t, input.Equals(output))
}

func TestSparseBitArrayMarshalUnmarshal(t *testing.T) {
	numItems := uint64(1280)
	input := newSparseBitArray()

	for i := uint64(0); i < numItems; i++ {
		if i%3 == 0 {
			input.SetBit(i)
		}
	}

	outputBytes, err := Marshal(input)
	assert.Equal(t, err, nil)
	assert.Equal(t, outputBytes[0], byte('S'))
	assert.Equal(t, len(outputBytes), 337)

	output, err := Unmarshal(outputBytes)
	assert.Equal(t, err, nil)

	assert.True(t, input.Equals(output))
}

func TestUnmarshalErrors(t *testing.T) {
	numItems := uint64(1280)
	input := newBitArray(numItems)

	for i := uint64(0); i < numItems; i++ {
		if i%3 == 0 {
			input.SetBit(i)
		}
	}

	outputBytes, err := Marshal(input)

	outputBytes[0] = 'C'

	output, err := Unmarshal(outputBytes)
	assert.Error(t, err)
	assert.Equal(t, output, nil)

	output, err = Unmarshal(nil)
	assert.Error(t, err)
	assert.Equal(t, output, nil)
}
