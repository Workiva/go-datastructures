package augmentedtree

import "runtime"

type trees []*tree

func (t trees) split() []trees {
	numParts := runtime.NumCPU()
	parts := make([]trees, numParts)
	for i := 0; i < numParts; i++ {
		parts[i] = t[i*len(t)/numParts : (i+1)*len(t)/numParts]
	}
	return parts
}
