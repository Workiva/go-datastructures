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

package btree

// Config defines all the parameters available to the UB-tree.
// Of most important are nodewidth and the persister to be used
// during commit phase.
type Config struct {
	// NodeWidth defines the branching factor of the tree.  Any node
	// wider than this value will get split and when the width of a node
	// falls to less than half this value the node gets merged.  This
	// ensures optimal performance while running to the key value store.
	NodeWidth int
	// Perister defines the key value store that the tree can use to
	// save and load nodes.
	Persister Persister
	// Comparator is the function used to determine ordering.
	Comparator Comparator `msg:"-"`
}

// DefaultConfig returns a configuration with the persister set.  All other
// fields are set to smart defaults for persistence.
func DefaultConfig(persister Persister, comparator Comparator) Config {
	return Config{
		NodeWidth:  10000,
		Persister:  persister,
		Comparator: comparator,
	}
}
