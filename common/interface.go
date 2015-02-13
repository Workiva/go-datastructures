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

package common

// Comparator is a generic interface that represents items that can
// be compared.
type Comparator interface {
	// Compare compares this interface with another.  Returns a positive
	// number if this interface is greater, 0 if equal, negative number
	// if less.
	Compare(Comparator) int
}

// Comparators is a typed list of type Comparator.
type Comparators []Comparator
