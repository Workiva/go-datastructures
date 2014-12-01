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
	_ "github.com/Workiva/go-datastructures/futures"
	_ "github.com/Workiva/go-datastructures/graph"
	_ "github.com/Workiva/go-datastructures/rangetree"
	_ "github.com/Workiva/go-datastructures/set"
	_ "github.com/Workiva/go-datastructures/threadsafe/err"
)
