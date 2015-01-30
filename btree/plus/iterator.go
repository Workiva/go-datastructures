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

const iteratorExhausted = -2

type iterator struct {
	node  *lnode
	index int
}

func (iter *iterator) Next() bool {
	if iter.index == iteratorExhausted {
		return false
	}

	iter.index++
	if iter.index >= len(iter.node.keys) {
		iter.node = iter.node.pointer
		if iter.node == nil {
			iter.index = iteratorExhausted
			return false
		}
		iter.index = 0
	}

	return true
}

func (iter *iterator) Value() Key {
	if iter.index == iteratorExhausted ||
		iter.index < 0 || iter.index >= len(iter.node.keys) {

		return nil
	}

	return iter.node.keys[iter.index]
}

// exhaust is a test function that's not exported
func (iter *iterator) exhaust() keys {
	keys := make(keys, 0, 10)
	for iter := iter; iter.Next(); {
		keys = append(keys, iter.Value())
	}

	return keys
}

func nilIterator() *iterator {
	return &iterator{
		index: iteratorExhausted,
	}
}
