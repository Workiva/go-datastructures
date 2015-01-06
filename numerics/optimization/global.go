package optimization

import (
	"math"
	"sort"
)

type pbs []*vertexProbabilityBundle

type vertexProbabilityBundle struct {
	probability float64
	vertex      *nmVertex
}

// calculateVVP will calculate the variable variance probability
// of the provided vertex based on the previous best guess
// and the provided sigma.  The sigma changes with each run
// of the optimization algorithm and accounts for a changing
// number of guesses.
//
// VVP is defined as:
// 1/((2*pi)^(1/2)*sigma)*(1-e^(-dmin^2/2*sigma^2))
// where dmin = euclidean distance between this vertex and the best guess
// and sigma = (3*(m^(1/n)))^-1
//
func calculateVVP(guess, vertex *nmVertex, sigma float64) float64 {
	distance := -guess.euclideanDistance(vertex)
	lhs := 1 / (math.Sqrt(2*math.Pi) * sigma)
	rhs := 1 - math.Exp(math.Pow(distance, 2)/(2*math.Pow(sigma, 2)))
	return rhs * lhs
}

// calculateSigma will calculate sigma based on the provided information.
// Typically, sigma will decrease as the number of sampled points
// increases.
//
// sigma = (3*(m^(1/n)))^-1
//
func calculateSigma(dimensions, guesses int) float64 {
	return math.Pow(3*math.Pow(float64(guesses), 1/float64(dimensions)), -1)
}

func (pbs pbs) calculateProbabilities(bestGuess *nmVertex, sigma float64) {
	for _, v := range pbs {
		v.probability = calculateVVP(bestGuess, v.vertex, sigma)
	}
}

func (pbs pbs) sort() {
	sort.Sort(pbs)
}

func (pbs pbs) Less(i, j int) bool {
	return pbs[i].probability < pbs[j].probability
}

func (pbs pbs) Swap(i, j int) {
	pbs[i], pbs[j] = pbs[j], pbs[i]
}

func (pbs pbs) Len() int {
	return len(pbs)
}

// results stores the results of previous iterations of the
// nelder-mead algorithm
type results struct {
	vertices vertices
	config   NelderMeadConfiguration
}

// search will search this list of results based on order, order
// being defined in the NelderMeadConfiguration, that is a defined
// target will be treated
func (results *results) search(result *nmVertex) int {
	return sort.Search(len(results.vertices), func(i int) {
		return !results.vertices[i].less(results.config, result)
	})
}

func (results *results) exists(result *nmVertex, hint int) bool {
	if hint < 0 {
		hint = results.search(result)
	}

	// maximum hint here should be len(results.vertices)
	if hint > 0 && results.vertices[hint-1].approximatelyEqualToVertex(result) {
		return true
	}

	// -1 here because if hint == len(vertices) we would've already
	// checked the last value in the previous conditional
	if hint < len(results.vertices)-1 && results.vertices[hint].approximatelyEqualToVertex(result) {
		return true
	}

	return false
}

func (results *results) insert(vertex *nmVertex) {
	i := results.search(vertex)
	if results.exists(vertex, i) {
		return
	}

	if i == len(results.vertices) {
		results.vertices = append(results.vertices, vertex)
		return
	}

	results.vertices = append(results.vertices, nil)
	copy(results.vertices[i+1:], results.vertices[i:])
	results.vertices[i] = vertex
}

func (results *results) grab(num int) vertices {
	vs := make(vertices, num)
	// first, copy what you want to the list to return
	// not returning a sub-slice as we're about to mutate
	// the original slice
	copy(vs, results.vertices[:num])
	// now we overwrite the vertices that we are taking
	// from the beginning
	copy(results.vertices, results.vertices[num:])
	length := len(results.vertices) - num
}
