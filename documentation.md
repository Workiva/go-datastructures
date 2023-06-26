# Introducing go-datastructures

The goal of the go-datastructures library is to port implementations of some common datastructures to Go or to improve on some existing datastructures.  These datastructures are designed to be re-used for anyone that needs them throughout the community. (and hopefully improved upon).

Given the commonality and popularity of these datastructures in other languages, it is hoped that by open sourcing this library we leverage a great deal of institutional knowledge to improve upon the Go-specific implementations.

# Datastructures

## Augmented Tree

Designed for determining intersections between ranges. For example, we can query the augmentedtree for any ranges that intersect with a cell, which can be represented as a range of size one (ie, Cell A1 can be represented as range A1:B2 where B2 is exclusive).  In this way, we can walk through the graph looking for exact or approximate intersections.

The current implementation exists in n-dimensions, but is quickest when an n-dimensional query can be reduced in its first dimension.  That is, queries only filtered in anything but the first dimension will be slowest.

The actual implementation is a top-down red-black binary search tree.

### Future

Implement a bottom-up version as well.

## Bit Array

Also known as a bitmap, a bitarray is useful for comparing two sets of data that can be represented as an integer.  It's useful because bitwise operations can compare a number of these integers at once instead of independently.  For instance, the sets {1, 3, 5} and {3, 5, 7} can be intersected in a single clock cycle if these sets were represented in their associated bit array.  Included in this package is the ability to convert a bitarray back to integers.

There are two implementations of bit arrays in this package, one is dense and the other borrows concepts from linear algebra's compressed row sparse matrix to represent bitarrays in much smaller spaces.  Unfortunately, the sparse version has logarithmic insertions and existence checks but retains some speed advantages when checking for intersections.

Incidentally, this is one of two things needed to build a native Go database.

### Future

Implement a dense but expandable bit array.  Optimize the current package to utilize larger amounts of mechanical sympathy.

## Futures

We ran into some cases where we wanted to indicate to a goroutine that an operation had started in another goroutine and to pause go routines until the initial routine had completed.  You can do this with buffered channels, but it seems somewhat redundant to send the same result to a channel to ensure all waiting threads were alerted.  Futures operate similarly to how ndb futures work in GAE and might be thought of as a "broadcast" channel.

## Queue

Pretty self-explanatory, this package includes both a queue and a priority queue.  Currently, waitgroups are used to orchestrate threads but with a proper constructor hint, this does end up being faster than channels when attempting to send data to a go routine.  The other advantage over a channel is that the queue will return an error if you attempt to put to a queue that has had Dispose called on it instead of panicking like what would happen if you attempted to send to a closed channel.  I believe this is closer to the Golang's stated design goals.

Speaking of Dispose, calling dispose on a queue will immediately return any waiting threads with an error.

### Future

When I get time, I'd like to implement a lockless ring buffer for further performance enhancements.

## Range Tree

The range tree is a way to store n-dimensional points of data in a manner that permits logarithmic-complexity queries.  These points are usually represented as points on a Cartesian graph represented by integers.

There are two implementations of a range tree in this package, one that is mutable and one that is immutable.  The mutable version can be faster, but involves lock contention if the consumer needs to ensure threadsafety.  The immutable version is a copy-on-write range tree that is optimized by only copying portions of the rangetree on write and is best written to in batches.  Operations on the immutable version are slower, but it is safe to read and write from this version at the same time from different threads.

Although rangetrees are often represented as BBSTs as described above, the n-dimensional nature of this rangetree actually made the design easier to implement as a sparse n-dimensional array.

### Future

Unite both implementations of the rangetree under the same interface.  The implementations (especially the immutable one) could use some further performance optimizations.

## Fibonacci Heap

The usual Fibonacci Heap with a floating-point priority key. Does a good job as a priority queue, especially for large n. Should be useful in writing an optimal solution for Dijkstra's and Prim's algorithms. (because of it's efficient decrease-key)

### Future

I'd like to add a value interface{} pointer that will be able to hold any user data attached to each node in the heap. Another thing would be writing a fast implementation of Dijkstra and Prim using this structure. And a third would be analysing thread-safety and coming up with a thread-safe variant.

## Set

Not much to say here.  This is an unordered set back by a Go map.  This particular version is threadsafe which does hinder performance a bit, although reads can happen simultaneously.  

### Future

I'd like to experiment with a ground-up implementation of a hash map using the standard library's hash/fnv hashing function, which is a non-cryptographic hash that's proven to be very fast.  I'd also like to experiment with a lockless hashmap.

## Slice

Golang's standard library "sort" includes a slice of ints that contain some sorting and searching functions.  This is like that standard library package but with Int64s, which requires a new package as Go doesn't want us to have generics.  I also added a method for inserting to the slice.

## Threadsafe

This package just wraps some common interfaces with a lock to make them threadsafe.  Golang would tell us to forget about locks and use channels (even though channels themselves are just fancy abstractions around locks as evidenced in their source code) but I found some situations where I wanted to protect some memory that was accessible from multiple goroutines where channels would be ugly, slow, and unnecessary.  The only interface with an implemntation thusfar is error, which is useful if you need to indicate that an error was returned from logic running in any number of goroutines.

# Going Forward

There is a PR into the datastructures repo that contains some pieces required for implementing a B+ tree.  With a B+ tree and bitmap, the pieces are in place to write a native Go database.  Going forward, I'd like to take these pieces, expand upon them, and implement a fast database in Go.  

As always, any optimizations or bug fixes in any of this code would be greatly appreciated and encouraged :).  These datastructures can and are the foundations of many programs and algorithms, even if they are abstracted away in different libraries which makes working with them a lot of fun and very informative.
