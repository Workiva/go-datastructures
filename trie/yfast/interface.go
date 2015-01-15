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

package yfast

// Entry defines items that can be inserted into the y-fast
// trie.
type Entry interface {
	// Key is the key for this entry.  If the trie has been
	// given bit size n, only the last n bits of this key
	// will matter.  Use a bit size of 64 to enable all
	// 2^64-1 keys.
	Key() uint64
}
