
go-datastructures
=================

Go-datastructes is a collection of useful and performant Go datastructures.

### NOTE: only tested with Go 1.3+.

#### Augmented Tree

Interval tree for collision in n-dimensional ranges.  Implemented via a
red-black augmented tree.  Extra dimensions are handled in simultaneous
inserts/queries to save space although this may result in suboptimal time
complexity.  Intersection determined using bit arrays.  In a single dimension,
inserts, deletes, and queries should be in O(log n) time.

#### Bitarray

Bitarray used to detect existence without having to resort to hashing with
hashmaps.  Requires entities have a uint64 unique identifier.  Two
implementations exist, regular and sparse.  Sparse saves a great deal of space
but insertions are O(log n).  There are some useful functions on the BitArray
interface to detect intersection between two bitarrays.

#### Futures

A helpful tool to send a "broadcast" message to listeners.  Channels have the
issue that once one listener takes a message from a channel the other listeners
aren't notified.  There were many cases when I wanted to notify many listeners
of a single event and this package helps.

#### Queue

Package contains both a normal and priority queue.  Both implementations never
block on send and grow as much as necessary.  Both also only return errors if
you attempt to push to a disposed queue and will not panic like sending a
message on a closed channel.  The priority queue also allows you to place items
in priority order inside the queue.  If you give a useful hint to the regular
queue, it is actually faster than a channel.  The priority queue is somewhat
slow currently and targeted for an update to a Fibonacci heap.

Also included in the queue package is a MPMC threadsafe ring buffer. This is a
block full/empty queue, but will return a blocked thread if the queue is
disposed while a thread is blocked.  This can be used to synchronize goroutines
and ensure goroutines quit so objects can be GC'd.  Threadsafety is acheived
using only CAS operations making this queue quite fast.  Benchmarks can be found
in that package.

#### Range Tree

Useful to determine if n-dimensional points fall within an n-dimensional range.
Not a typical range tree however, as we are actually using an n-dimensional
sorted list of points as this proved to be simpler and faster than attempting a
traditional range tree while saving space on any dimension greater than one.
Inserts are typical BBST times at O(log n^d) where d is the number of
dimensions.

#### Set
Self explanatory.  Could be further optimized by getting the uintptr of the
generic interface{} used and using that as the key as Golang maps handle that
much better than the generic struct type.

#### Threadsafe
A package that is meant to contain some commonly used items but in a threadsafe
way.  Example: there's a threadsafe error in there as I commonly found myself
wanting to set an error in many threads at the same time (yes, I know, but
channels are slow).

#### AVL Tree

This is an example of a branch copy immutable AVL BBST.  Any operation on a node
makes a copy of that node's branch.  Because of this, this tree is inherently
threadsafe although the writes will likely still need to be serialized.  This
structure is good if your use case is a large number of reads and infrequent
writes as reads will be highly available but writes somewhat slow due to the
copying.  This structure serves as a basis for a large number of functional data
structures.

#### X-Fast Trie

An interesting design that treats integers as words and uses a trie structure to
reduce time complexities by matching prefixes.  This structure is really fast
for finding values or making predecessor/successor types of queries, but also
results in greater than linear space consumption.  The exact time complexities
can be found in that package.

#### Y-Fast Trie

An extension of the X-Fast trie in which an X-Fast trie is combined with some
other ordered data structure to reduce space consumption and improve CRUD types
of operations.  These secondary structures are often BSTs, but our implemention
uses a simple ordered list as I believe this improves cache locality.  We also
use fixed size buckets to aid in parallelization of operations.  Exact time
complexities are in that package.

#### Fast integer hashmap

A datastructure used for checking existence but without knowing the bounds of
your data.  If you have a limited small bounds, the bitarray package might be a
better choice.  This implementation uses a fairly simple hashing alogrithm
combined with linear probing and a flat datastructure to provide optimal
performance up to a few million integers (faster than the native Golang
implementation).  Beyond that, the native implementation is faster (I believe
they are using a large -ary B-tree).  In the future, this will be implemented
with a B-tree for scale.

#### Skiplist

An ordered structure that provides amoritized logarithmic operations but without
the complication of rotations that are required by BSTs.  In testing, however,
the performance of the skip list is often far worse than the guaranteed log n
time of a BBST.  Tall nodes tend to "cast shadows", especially when large
bitsizes are required as the optimum maximum height for a node is often based on
this.  More detailed performance characteristics are provided in that package.

#### Sort

The sort package implements a multithreaded bucket sort that can be up to 3x
faster than the native Golang sort package.  These buckets are then merged using
a symmetrical merge, similar to the stable sort in the Golang package.  However,
our algorithm is modified so that two sorted lists can be merged by using
symmetrical decomposition.

#### Numerics

Early work on some nonlinear optimization problems.  The initial implementation
allows a simple use case with either linear or nonlinear constraints.  You can
find min/max or target an optimal value.  The package currently employs a
probablistic global restart system in an attempt to avoid local critical points.
More details can be found in that package.

#### B+ Tree

Initial implementation of a B+ tree.  Delete method still needs added as well as
some performance optimization.  Specific performance characteristics can be
found in that package.  Despite the theoretical superiority of BSTs, the B-tree
often has better all around performance due to cache locality.  The current
implementation is mutable, but the immutable AVL tree can be used to build an
immutable version.  Unfortunately, to make the B-tree generic we require an
interface and the most expensive operation in CPU profiling is the interface
method which in turn calls into runtime.assertI2T.  We need generics.

### Installation

 1. Install Go 1.3 or higher.
 2. Run `go get github.com/Workiva/go-datastructures/...`

### Updating

When new code is merged to master, you can use

	go get -u github.com/Workiva/go-datastructures/...

To retrieve the latest version of go-datastructures.

### Testing

To run all the unit tests use these commands:

	cd $GOPATH/src/github.com/Workiva/go-datastructures
	go get -t -u ./...
	go test ./...

Once you've done this once, you can simply use this command to run all unit tests:

	go test ./...


### Contributing

Requirements to commit here:

 - Branch off master, PR back to master.
 - `gofmt`'d code.
 - Compliance with [these guidelines](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
 - Unit test coverage
 - [Good commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)


### Maintainers

 - Dustin Hiatt <[dustin.hiatt@workiva.com](mailto:dustin.hiatt@workiva.com)>
 - Alexander Campbell <[alexander.campbell@workiva.com](mailto:alexander.campbell@workiva.com)>

