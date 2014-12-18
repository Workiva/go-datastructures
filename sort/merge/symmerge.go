package merge

import (
	"math"
	"runtime"
	"sort"
	"sync"
)

func symSearch(u, w Comparators) int {
	start, stop, p := 0, len(u), len(w)-1
	for start < stop {
		mid := (start + stop) / 2
		if u[mid].Compare(w[p-mid]) <= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}

	return start
}

func swap(u, w Comparators, index int) {
	for i := index; i < len(u); i++ {
		u[i], w[i-index] = w[i-index], u[i]
	}
}

func decomposeForSymMerge(length int,
	comparators Comparators) (v1 Comparators,
	w Comparators, v2 Comparators) {

	if length >= len(comparators) {
		panic(`INCORRECT PARAMS FOR SYM MERGE.`)
	}

	overhang := (len(comparators) - length) / 2
	v1 = comparators[:overhang]
	w = comparators[overhang : overhang+length]
	v2 = comparators[overhang+length:]
	return
}

func symBinarySearch(u Comparators, start, stop, total int) int {
	for start < stop {
		mid := (start + stop) / 2
		if u[mid].Compare(u[total-mid]) <= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}

	return start
}

func symSwap(u Comparators, start1, start2, end int) {
	for i := 0; i < end; i++ {
		u[start1+i], u[start2+i] = u[start2+i], u[start1+i]
	}
}

func symRotate(u Comparators, start1, start2, end int) {
	i := start2 - start1
	if i == 0 {
		return
	}

	j := end - start2
	if j == 0 {
		return
	}

	if i == j {
		symSwap(u, start1, start2, i)
		return
	}

	p := start1 + i
	for i != j {
		if i > j {
			symSwap(u, p-i, p, j)
			i -= j
		} else {
			symSwap(u, p-i, p+j-i, i)
			j -= i
		}
	}
	symSwap(u, p-i, p, i)
}

func symMerge(u Comparators, start1, start2, last int) {
	if start1 < start2 && start2 < last {
		mid := (start1 + last) / 2
		n := mid + start2
		var start int
		if start2 > mid {
			start = symBinarySearch(u, n-last, mid, n-1)
		} else {
			start = symBinarySearch(u, start1, start2, n-1)
		}
		end := n - start

		symRotate(u, start, start2, end)
		symMerge(u, start1, start, mid)
		symMerge(u, mid, end, last)
	}
}

func SymMerge(u, w Comparators) Comparators {
	lenU, lenW := len(u), len(w)
	if lenU == 0 {
		return w
	}

	if lenW == 0 {
		return u
	}

	diff := lenU - lenW
	if math.Abs(float64(diff)) > 1 {
		u1, w1, u2, w2 := prepareForSymMerge(u, w)

		lenU1 := len(u1)
		lenU2 := len(u2)
		u = append(u1, w1...)
		w = append(u2, w2...)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			symMerge(u, 0, lenU1, len(u))
			wg.Done()
		}()
		go func() {
			symMerge(w, 0, lenU2, len(w))
			wg.Done()
		}()

		wg.Wait()
		u = append(u, w...)
		return u
	}

	u = append(u, w...)
	symMerge(u, 0, lenU, len(u))
	return u
}

func prepareForSymMerge(u, w Comparators) (u1, w1, u2, w2 Comparators) {
	if u.Len() > w.Len() {
		u, w = w, u
	}
	v1, w, v2 := decomposeForSymMerge(len(u), w)

	i := symSearch(u, w)

	u1 = make(Comparators, i)
	copy(u1, u[:i])
	w1 = append(v1, w[:len(w)-i]...)

	u2 = make(Comparators, len(u)-i)
	copy(u2, u[i:])

	w2 = append(w[len(w)-i:], v2...)
	return
}

func sortBucket(comparators Comparators) {
	sort.Sort(comparators)
}

func copyChunk(chunk []Comparators) []Comparators {
	cp := make([]Comparators, len(chunk))
	copy(cp, chunk)
	return cp
}

func multithreadedSortComparators(comparators Comparators) Comparators {
	toBeSorted := make(Comparators, len(comparators))
	copy(toBeSorted, comparators)

	var wg sync.WaitGroup
	chunks := chunk(toBeSorted, int64(runtime.NumCPU()))
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
