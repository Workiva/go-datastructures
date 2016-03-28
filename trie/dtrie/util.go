/*
Copyright (c) 2016, Theodore Butler
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package dtrie

import (
	"fmt"
	"hash/fnv"
)

func mask(hash, level uint32) uint32 {
	return (hash >> (5 * level)) & 0x01f
}

func setBit(bitmap uint32, pos uint32) uint32 {
	return bitmap | (1 << pos)
}

func clearBit(bitmap uint32, pos uint32) uint32 {
	return bitmap & ^(1 << pos)
}

func hasBit(bitmap uint32, pos uint32) bool {
	return (bitmap & (1 << pos)) != 0
}

func popCount(bitmap uint32) int {
	// bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	bitmap -= (bitmap >> 1) & 0x55555555
	bitmap = (bitmap>>2)&0x33333333 + bitmap&0x33333333
	bitmap += bitmap >> 4
	bitmap &= 0x0f0f0f0f
	bitmap *= 0x01010101
	return int(byte(bitmap >> 24))
}

func defaultHasher(value interface{}) uint32 {
	switch value.(type) {
	case uint8:
		return uint32(value.(uint8))
	case uint16:
		return uint32(value.(uint16))
	case uint32:
		return value.(uint32)
	case uint64:
		return uint32(value.(uint64))
	case int8:
		return uint32(value.(int8))
	case int16:
		return uint32(value.(int16))
	case int32:
		return uint32(value.(int32))
	case int64:
		return uint32(value.(int64))
	case uint:
		return uint32(value.(uint))
	case int:
		return uint32(value.(int))
	case uintptr:
		return uint32(value.(uintptr))
	case float32:
		return uint32(value.(float32))
	case float64:
		return uint32(value.(float64))
	}
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%#v", value)))
	return hasher.Sum32()
}
