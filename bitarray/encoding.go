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
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Marshal takes a dense or sparse bit array and serializes it to a
// byte slice.
func Marshal(ba BitArray) ([]byte, error) {
	if eba, ok := ba.(*bitArray); ok {
		return eba.Serialize()
	} else if sba, ok := ba.(*sparseBitArray); ok {
		return sba.Serialize()
	} else {
		return nil, errors.New("not a valid BitArray")
	}
}

// Unmarshal takes a byte slice, of the same format produced by Marshal,
// and returns a BitArray.
func Unmarshal(input []byte) (BitArray, error) {
	if len(input) == 0 {
		return nil, errors.New("no data in input")
	}
	if input[0] == 'B' {
		ret := newBitArray(0)
		err := ret.Deserialize(input)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else if input[0] == 'S' {
		ret := newSparseBitArray()
		err := ret.Deserialize(input)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		return nil, errors.New("unrecognized encoding")
	}
}

// Serialize converts the sparseBitArray to a byte slice
func (ba *sparseBitArray) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)

	var identifier uint8 = 'S'
	err := binary.Write(w, binary.LittleEndian, identifier)
	if err != nil {
		return nil, err
	}

	blocksLen := uint64(len(ba.blocks))
	indexLen := uint64(len(ba.indices))

	err = binary.Write(w, binary.LittleEndian, blocksLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.blocks)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, indexLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.indices)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Deserialize takes the incoming byte slice, and populates the sparseBitArray
// with data in the bytes. Note that this will overwrite any capacity
// specified when creating the sparseBitArray. Also note that if an error
// is returned, the sparseBitArray this is called on might be populated
// with partial data.
func (ret *sparseBitArray) Deserialize(incoming []byte) error {
	r := bytes.NewReader(incoming[1:]) // Discard identifier

	var intsToRead uint64
	err := binary.Read(r, binary.LittleEndian, &intsToRead)
	if err != nil {
		return err
	}

	var nextblock block
	for i := intsToRead; i > uint64(0); i-- {
		err = binary.Read(r, binary.LittleEndian, &nextblock)
		if err != nil {
			return err
		}
		ret.blocks = append(ret.blocks, nextblock)
	}

	err = binary.Read(r, binary.LittleEndian, &intsToRead)
	if err != nil {
		return err
	}

	var nextuint uint64
	for i := intsToRead; i > uint64(0); i-- {
		err = binary.Read(r, binary.LittleEndian, &nextuint)
		if err != nil {
			return err
		}
		ret.indices = append(ret.indices, nextuint)
	}

	return nil
}

// Serialize converts the bitArray to a byte slice.
func (ba *bitArray) Serialize() ([]byte, error) {
	w := new(bytes.Buffer)

	var identifier uint8 = 'B'
	err := binary.Write(w, binary.LittleEndian, identifier)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.lowest)
	if err != nil {
		return nil, err
	}
	err = binary.Write(w, binary.LittleEndian, ba.highest)
	if err != nil {
		return nil, err
	}

	var encodedanyset uint8
	if ba.anyset {
		encodedanyset = 1
	} else {
		encodedanyset = 0
	}
	err = binary.Write(w, binary.LittleEndian, encodedanyset)
	if err != nil {
		return nil, err
	}

	err = binary.Write(w, binary.LittleEndian, ba.blocks)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Deserialize takes the incoming byte slice, and populates the bitArray
// with data in the bytes. Note that this will overwrite any capacity
// specified when creating the bitArray. Also note that if an error is returned,
// the bitArray this is called on might be populated with partial data.
func (ret *bitArray) Deserialize(incoming []byte) error {
	r := bytes.NewReader(incoming[1:]) // Discard identifier

	err := binary.Read(r, binary.LittleEndian, &ret.lowest)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &ret.highest)
	if err != nil {
		return err
	}

	var encodedanyset uint8
	err = binary.Read(r, binary.LittleEndian, &encodedanyset)
	if err != nil {
		return err
	}

	// anyset defaults to false so we don't need an else statement
	if encodedanyset == 1 {
		ret.anyset = true
	}

	var nextblock block
	err = binary.Read(r, binary.LittleEndian, &nextblock)
	for err == nil {
		ret.blocks = append(ret.blocks, nextblock)
		err = binary.Read(r, binary.LittleEndian, &nextblock)
	}
	if err != io.EOF {
		return err
	}
	return nil
}
