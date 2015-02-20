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
Package Hilbert is designed to allow consumers to find the Hilbert
distance on the Hilbert curve if given a 2 dimensional coordinate.
This could be useful for hashing or constructing a Hilbert R-Tree.
Algorithm taken from here:

http://en.wikipedia.org/wiki/Hilbert_curve

This expects coordinates in the range [0, 0] to [MaxInt32, MaxInt32].
Using negative values for x and y will have undefinied behavior.

Benchmarks:
BenchmarkEncode-8	10000000	       181 ns/op
BenchmarkDecode-8	10000000	       191 ns/op
*/
package hilbert

// n defines the maximum power of 2 that can define a bound,
// this is the value for 2-d space if you want to support
// all hilbert ids with a single integer variable
const n = 1 << 31

func boolToInt(value bool) int32 {
	if value {
		return int32(1)
	}

	return int32(0)
}

func rotate(n, rx, ry int32, x, y *int32) {
	if ry == 0 {
		if rx == 1 {
			*x = n - 1 - *x
			*y = n - 1 - *y
		}

		t := *x
		*x = *y
		*y = t
	}
}

// Encode will encode the provided x and y coordinates into a Hilbert
// distance.
func Encode(x, y int32) int64 {
	var rx, ry int32
	var d int64
	for s := int32(n / 2); s > 0; s /= 2 {
		rx = boolToInt(x&s > 0)
		ry = boolToInt(y&s > 0)
		d += int64(int64(s) * int64(s) * int64(((3 * rx) ^ ry)))
		rotate(s, rx, ry, &x, &y)
	}

	return d
}

// Decode will decode the provided Hilbert distance into a corresponding
// x and y value, respectively.
func Decode(h int64) (int32, int32) {
	var ry, rx int64
	var x, y int32
	t := h

	for s := int64(1); s < int64(n); s *= 2 {
		rx = 1 & (t / 2)
		ry = 1 & (t ^ rx)
		rotate(int32(s), int32(rx), int32(ry), &x, &y)
		x += int32(s * rx)
		y += int32(s * ry)
		t /= 4
	}

	return x, y
}
