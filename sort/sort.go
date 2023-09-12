package merge

import (
	"runtime"
	"sort"
	"sync"
)

func sortBucket(comparators Comparators) {
	sort.Sort(comparators)
}

func copyChunk(chunk []Comparators) []Comparators {
	cp := make([]Comparators, len(chunk))
	copy(cp, chunk)
	return cp
}

// MultithreadedSortComparators will take a list of comparators
// and sort it using as many threads as are available.  The list
// is split into buckets for a bucket sort and then recursively
// merged using SymMerge.
func MultithreadedSortComparators(comparators Comparators) Comparators {
	toBeSorted := make(Comparators, len(comparators))
	copy(toBeSorted, comparators)

	var wg sync.WaitGroup

	numCPU := int64(runtime.NumCPU())
	if numCPU == 1 { // single core machine
		numCPU = 2
	} else {
		// otherwise this algo only works with a power of two
		numCPU = int64(prevPowerOfTwo(uint64(numCPU)))
	}

	chunks := chunk(toBeSorted, numCPU)
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			sortBucket(chunks[i])
			wg.Done()
		}(i)
	}

	wg.Wait()
	todo := make([]Comparators, len(chunks)/2)
	for {
		todo = todo[:len(chunks)/2]
		wg.Add(len(chunks) / 2)
		for i := 0; i < len(chunks); i += 2 {
			go func(i int) {
				todo[i/2] = SymMerge(chunks[i], chunks[i+1])
				wg.Done()
			}(i)
		}

		wg.Wait()

		chunks = copyChunk(todo)
		if len(chunks) == 1 {
			break
		}
	}

	return chunks[0]
}

func chunk(comparators Comparators, numParts int64) []Comparators {
	parts := make([]Comparators, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = comparators[i*int64(len(comparators))/numParts : (i+1)*int64(len(comparators))/numParts]
	}
	return parts
}

func prevPowerOfTwo(x uint64) uint64 {
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >> 16)
	x = x | (x >> 32)
	return x - (x >> 1)
}
