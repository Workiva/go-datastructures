package datastructures

/*
Package datastructures exists solely to aid consumers of the
go-datastructures library when using dependency managers.  Depman,
for instance, will work correctly with any datastructure by simply
importing this package instead of each subpackage individually.
*/

import (
	_ "github.com/Workiva/go-datastructures/augmentedtree"
	_ "github.com/Workiva/go-datastructures/bitarray"
	_ "github.com/Workiva/go-datastructures/btree/palm"
	_ "github.com/Workiva/go-datastructures/btree/plus"
	_ "github.com/Workiva/go-datastructures/futures"
	_ "github.com/Workiva/go-datastructures/hashmap/fastinteger"
	_ "github.com/Workiva/go-datastructures/numerics/optimization"
	_ "github.com/Workiva/go-datastructures/queue"
	_ "github.com/Workiva/go-datastructures/rangetree"
	_ "github.com/Workiva/go-datastructures/rangetree/skiplist"
	_ "github.com/Workiva/go-datastructures/set"
	_ "github.com/Workiva/go-datastructures/slice"
	_ "github.com/Workiva/go-datastructures/slice/skip"
	_ "github.com/Workiva/go-datastructures/sort"
	_ "github.com/Workiva/go-datastructures/threadsafe/err"
	_ "github.com/Workiva/go-datastructures/tree/avl"
	_ "github.com/Workiva/go-datastructures/trie/xfast"
	_ "github.com/Workiva/go-datastructures/trie/yfast"
)
