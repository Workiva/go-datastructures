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

import "github.com/Workiva/go-datastructures/common"

// Iterator defines an interface that allows a consumer to iterate
// all results of a query.  All values will be visited in-order.
type Iterator interface {
	// Next returns a bool indicating if there is future value
	// in the iterator and moves the iterator to that value.
	Next() bool
	// Value returns a Comparator representing the iterator's current
	// position.  If there is no value, this returns nil.
	Value() common.Comparator
	// exhaust is a helper method that will iterate this iterator
	// to completion and return a list of resulting Entries
	// in order.
	exhaust() common.Comparators
}
