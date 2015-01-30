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

package plus

// Keys is a typed list of Key interfaces.
type Keys []Key

type Key interface {
	// Compare should return an int indicating how this key relates
	// to the provided key.  -1 will indicate less than, 0 will indicate
	// equality, and 1 will indicate greater than.  Duplicate keys
	// are allowed, but duplicate IDs are not.
	Compare(Key) int
}

// Iterator will be called with matching keys until either false is
// returned or we run out of keys to iterate.
type Iterator interface {
	// Next will move the iterator to the next position and return
	// a bool indicating if there is a value.
	Next() bool
	// Value returns a Key at the associated iterator position.  Returns
	// nil if the iterator is exhausted or has never been nexted.
	Value() Key
	// exhaust is an internal helper method to iterate this iterator
	// until exhausted and returns the resulting list of keys.
	exhaust() keys
}
