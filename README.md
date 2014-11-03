go-datastructure
================

Package containing some useful datastructures when writing code in Go.

### NOTE: only tested with Go 1.3+.

#### Augmented Tree: 
Interval tree for collision in n-dimensional ranges.  Implemented via a red-black augmented tree.  Extra dimensions are handled in simultaneous inserts/queries to save space although this may result in suboptimal time complexity.  Intersection determined using bit arrays.  In a single dimension, inserts, deletes, and queries should be in O(log n) time.

#### Bitarray: 
Bitarray used to detect existence without having to resort to hashing with hashmaps.  Requires entities have a uint64 unique identifier.  Two implementations exist, regular and sparse.  Sparse saves a great deal of space but insertions are O(log n).  There are some useful functions on the BitArray interface to detect intersection between two bitarrays.

#### Futures: 
A helpful tool to send a "broadcast" message to listeners.  Channels have the issue that once one listener takes a message from a channel the other listeners aren't notified.  There were many cases when I wanted to notify many listeners of a single event and this package helps.

#### Graph: 
Still pretty specific to gotable, but contains logic required to maintain graph state.  Also has logic to flatten graph into executable chunks.

#### Queue: 
Package contains both a normal and priority queue.  Both implementations never block on send and grow as much as necessary.  Both also only return errors if you attempt to push to a disposed queue and will not panic like sending a message on a closed channel.  The priority queue also allows you to place items in priority order inside the queue.  If you give a useful hint to the regular queue, it is actually faster than a channel.

#### Range Tree: 
Useful to determine if n-dimensional points fall within an n-dimensional range.  Not a typical range tree however, as we are actually using an n-dimensional sorted list of points as this proved to be simpler and faster than attempting a traditional range tree while saving space on any dimension greater than one.  Inserts are typical BBST times at O(log n^d) where d is the number of dimensions.

#### Set: 
Self explanatory.  Could be further optimized by getting the uintptr of the generic interface{} used and using that as the key as Golang maps handle that much better than the generic struct type.

#### Threadsafe: 
A package that is meant to contain some commonly used items but in a threadsafe way.  Example: there's a threadsafe error in there as I commonly found myself wanting to set an error in many threads at the same time (yes, I know, but channels are slow).

### Installation

1) Install Go 1.3 or higher.

2) Configure git to use SSH instead of HTTPS for github repositories. This
allows `go get` to use private repositories.

	# ~/.gitconfig
	[url "git@github.com:"]
		insteadOf = https://github.com

3) go get github.com/Workiva/go-datastructures ...

### Updating

When new code is merged to master, you can use 

	go get -u github.com/Workiva/go-datastructures/...

To retrieve the latest version of go-datastructures.

### Testing

To run all the unit tests use these commands:

	cd $GOPATH/src/github.com/Workiva/go-datastructures
	go get -t -u ./...
	go test ./...

Once you've done this once, you can simply use

	go test ./...

### Notice

Requirements to commit here:

 - `gofmt`'d code.
 - Compliance with [these guidelines](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
 - Unit test coverage
 - [Good commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)