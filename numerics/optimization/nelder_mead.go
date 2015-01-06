package optimization

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"
)

func init() {
	log.Println(`OH DO I HATE THIS.`)
}

const (
	alpha         = 1     // reflection, must be > 0
	beta          = 2     // expansion, must be > 1
	gamma         = .5    // contraction, 0 < gamma < 1
	sigma         = .5    // shrink, 0 < sigma < 1
	delta         = .0001 // going to use this to determine convergence
	maxRuns       = 130
	maxIterations = 5 // maxIterations defines the number of restarts that should
	// occur when attempting to find a global critical point
)

var (
	min = math.Inf(-1)
	max = math.Inf(1)
)

// generateRandomVerticesFromGuess will generate num number of vertices
// with random
func generateRandomVerticesFromGuess(guess *nmVertex, num int) vertices {
	// summed allows us to prevent duplicate guesses, checking
	// all previous guesses for every guess created would be too
	// time consuming so we take an indexed shortcut here.  summed
	// is a map of a sum of the vars to the vertices that have an
	// identical sum.  In this way, we can sum the vars of a new guess
	// and check only a small subset of previous guesses to determine
	// if this is an identical guess.
	summed := make(map[float64]vertices, num)
	dimensions := len(guess.vars)
	vs := make(vertices, 0, num)
	i := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

Guess:
	for i < num {
		sum := float64(0)
		vars := make([]float64, 0, dimensions)
		for j := 0; j < dimensions; j++ {
			v := r.Float64() * 1000
			// we do a separate random check here to determine
			// sign so we don't end up with all high v's one sign
			// and low v's another
			if r.Float64() > .5 {
				v = -v
			}
			sum += v
			vars = append(vars, v)
		}

		guess := &nmVertex{
			vars: vars,
		}

		if vs, ok := summed[sum]; !ok {
			vs = make(vertices, 0, dimensions) // dimensions is really just a guess, no real way of knowing what this is
			vs = append(vs, guess)
			summed[sum] = vs
		} else {
			for _, vertex := range vs {
				// if we've already guessed this, try the loop again
				if guess.equalToVertex(vertex) {
					continue Guess
				}
			}
			vs = append(vs, guess)
		}

		vs = append(vs, guess)
		i++
	}

	return vs
}

func isInf(num float64) bool {
	return math.IsInf(num, -1) || math.IsInf(num, 1)
}

func findMin(vertices ...*nmVertex) *nmVertex {
	min := vertices[0]
	for _, v := range vertices[1:] {
		if v.distance < min.distance {
			min = v
		}
	}

	return min
}

// findMidpoint will find the midpoint of the provided vertices
// and return a new vertex.
func findMidpoint(vertices ...*nmVertex) *nmVertex {
	num := len(vertices) // this is what we divide by
	vars := make([]float64, 0, num)

	for i := 0; i < num; i++ {
		sum := float64(0)
		for _, v := range vertices {
			sum += v.vars[i]
		}
		vars = append(vars, sum/float64(num))
	}

	return &nmVertex{
		vars: vars,
	}
}

// determineDistance will determine the distance between the value
// and the target.  If the target is positive or negative infinity,
// (ie find max or min), this is clamped to max or min float64.
func determineDistance(value, target float64) float64 {
	if math.IsInf(target, 1) { // positive infinity
		target = math.MaxFloat64
	} else if math.IsInf(target, -1) { // negative infinity
		target = -math.MaxFloat64
	}

	return math.Abs(target - value)
}

type vertices []*nmVertex

// evaluate will call evaluate on all the verticies in this list
// and order them by distance to target.
func (vertices vertices) evaluate(config NelderMeadConfiguration) {
	for _, v := range vertices {
		v.evaluate(config)
	}

	vertices.sort(config)
}

func (vertices vertices) sort(config NelderMeadConfiguration) {
	sorter := sorter{
		config:   config,
		vertices: vertices,
	}
	sorter.sort()
}

type sorter struct {
	config   NelderMeadConfiguration
	vertices vertices
}

func (sorter sorter) sort() {
	sort.Sort(sorter)
}

// the following methods are required for sort.Interface.  We
// use the standard libraries sort here as it uses an adaptive
// sort and we really don't expect there to be a ton of dimensions
// here so mulithreaded sort in this repo really isn't
// necessary.

func (sorter sorter) Less(i, j int) bool {
	return sorter.vertices[i].less(sorter.config, sorter.vertices[j])
}

func (sorter sorter) Len() int {
	return len(sorter.vertices)
}

func (sorter sorter) Swap(i, j int) {
	sorter.vertices[i], sorter.vertices[j] = sorter.vertices[j], sorter.vertices[i]
}

// String prints out a string representation of every vertex in this list.
// Useful for debugging :).
func (vertices vertices) String() string {
	result := ``
	for i, v := range vertices {
		result += fmt.Sprintf(`VERTEX INDEX: %+v, VERTEX: %+v`, i, v)
		result += fmt.Sprintln(``)
	}

	return result
}

// NelderMeadConfiguration is the struct that must be
// passed into the NelderMead function.  This defines
// the target value, the function to be run, and a guess
// of the variables.
type NelderMeadConfiguration struct {
	// Target is the target we are trying to converge
	// to.  Set this to positive or negative infinity
	// to find the min/max.
	Target float64
	// Fn defines the function that Nelder Mead is going
	// to call to determine if it is moving closer
	// to convergence.  In all likelihood, the execution
	// of this function is going to be the bottleneck.
	// The second value returns a bool indicating if the
	// calculated values are "good", that is, that no
	// constraint has been hit.
	Fn func([]float64) (float64, bool)
	// Vars is a guess and will determine what other
	// vertices will be used.  By convention, since
	// this guess will contain as many numbers as the
	// target function requires, the len of Vars determines
	// the dimension of this problem.
	Vars []float64
}

type nmVertex struct {
	// vars indicates the values used to calculate this vertex.
	vars []float64
	// distance is the distance between this vertex and the desired
	// value.  This metric has little meaning if the desired value
	// is +- inf.
	// result is the calculated result of this vertex.  This can
	// be used to measure distance or as a metrix to compare two
	// vertices if the desired result is a min/max.
	distance, result float64
	// good indicates if the calculated values here
	// are within all constraints, this should always
	// be true if this vertex is in a list of vertices.
	good bool
}

func (nm *nmVertex) evaluate(config NelderMeadConfiguration) {
	nm.result, nm.good = config.Fn(nm.vars)
	nm.distance = determineDistance(nm.result, config.Target)
}

func (nm *nmVertex) add(other *nmVertex) *nmVertex {
	vars := make([]float64, 0, len(nm.vars))
	for i := 0; i < len(nm.vars); i++ {
		vars = append(vars, nm.vars[i]+other.vars[i])
	}

	return &nmVertex{
		vars: vars,
	}
}

func (nm *nmVertex) multiply(scalar float64) *nmVertex {
	vars := make([]float64, 0, len(nm.vars))
	for i := 0; i < len(nm.vars); i++ {
		vars = append(vars, nm.vars[i]*scalar)
	}

	return &nmVertex{
		vars: vars,
	}
}

func (nm *nmVertex) subtract(other *nmVertex) *nmVertex {
	vars := make([]float64, 0, len(nm.vars))
	for i := 0; i < len(nm.vars); i++ {
		vars = append(vars, nm.vars[i]-other.vars[i])
	}

	return &nmVertex{
		vars: vars,
	}
}

// less defines a relationship between two points.  It is best not to
// think of less as returning a value indicating absolute relationship between
// two points, but instead think of less returning a bool indicating
// if this vertex is *closer* to the desired convergence, or a delta
// less than the other vertex.  For -inf, this returns a value indicating
// if this vertex has a less absolute value than the other vertex, if +inf
// less returns a bool indicating if this vertex has a *greater* absolute
// value than the other vertex.  Otherwise, this method returns a bool
// indicating if this vertex is closer to *converging* upon the desired
// value.
func (nm *nmVertex) less(config NelderMeadConfiguration, other *nmVertex) bool {
	if config.Target == min { // looking for a min
		return nm.result < other.result
	}
	if config.Target == max { // looking for a max
		return nm.result > other.result
	}

	return nm.distance < other.distance
}

func (nm *nmVertex) equal(config NelderMeadConfiguration, other *nmVertex) bool {
	if isInf(config.Target) {
		// if we are looking for a min or max, we compare result
		return nm.result == other.result
	}

	// otherwise, we compare distances
	return nm.distance == other.distance
}

// euclideanDistance determines the euclidean distance between two points.
func (nm *nmVertex) euclideanDistance(other *nmVertex) float64 {
	sum := float64(0)
	// first we want to sum all the distances between the points
	for i, otherPoint := range other.vars {
		// distance between points is defined by (qi-ri)^2
		sum += math.Pow(otherPoint-nm.vars[i], 2)
	}

	return math.Sqrt(sum)
}

// equalToVertex will compare this vertex to the provided vertex
// to determine if the two vertices are actually identical (that is,
// they fall on the same point).
func (nm *nmVertex) equalToVertex(other *nmVertex) bool {
	for i, n := range nm.vars {
		if n != other.vars[i] {
			return false
		}
	}

	return true
}

// approximatelyEqualToVertex returns a bool indicating if the
// *result* of this vertex is approximately equal to the vertex
// provided.  Approximately is 2 * delta as the algorithm may
// cease within a delta distance of the true value, so we may
// end up with a result that's 2*delta away if we came from
// the other direction.
func (nm *nmVertex) approximatelyEqualToVertex(other *nmVertex) bool {
	return math.Abs(nm.result-other.result) < 2*delta
}

type nelderMead struct {
	config  NelderMeadConfiguration
	results *results
}

// evaluateWithConstraints will safely evaluate the vertex while
// conforming to any imposed restraints.  If a constraint is found,
// this method will backtrack the vertex as described here:
// http://www.iccm-central.org/Proceedings/ICCM16proceedings/contents/pdf/MonK/MoKA1-04ge_ghiasimh224461p.pdf
// This should work with even non-linear constraints, but it is up to
// the consumer to check these constraints.
func (nm *nelderMead) evaluateWithConstraints(vertices vertices, vertex *nmVertex) *nmVertex {
	vertex.evaluate(nm.config)
	return vertex
	if vertex.good {
		return vertex
	}
	best := vertices[0]
	for i := 0; i < 5; i++ {
		vertex = best.add((vertex.subtract(best).multiply(alpha)))
		if vertex.good {
			return vertex
		}
	}

	return best
}

// reflect will find the reflection point between the two best guesses
// with the provided midpoint.
func (nm *nelderMead) reflect(vertices vertices, midpoint *nmVertex) *nmVertex {
	toScalar := midpoint.subtract(nm.lastVertex(vertices))
	toScalar = toScalar.multiply(alpha)
	toScalar = midpoint.add(toScalar)
	return nm.evaluateWithConstraints(vertices, toScalar)
}

func (nm *nelderMead) expand(vertices vertices, midpoint, reflection *nmVertex) *nmVertex {
	toScalar := reflection.subtract(midpoint)
	toScalar = toScalar.multiply(beta)
	toScalar = midpoint.add(toScalar)
	return nm.evaluateWithConstraints(vertices, toScalar)
}

// lastDimensionVertex returns the vertex that is represented by the
// last dimension, effectively, second to last in the list of
// vertices.
func (nm *nelderMead) lastDimensionVertex(vertices vertices) *nmVertex {
	return vertices[len(vertices)-2]
}

// lastVertex returns the last vertex in the list of vertices.
// It's important to remember that this vertex represents the
// number of dimensions + 1.
func (nm *nelderMead) lastVertex(vertices vertices) *nmVertex {
	return vertices[len(vertices)-1]
}

func (nm *nelderMead) outsideContract(vertices vertices, midpoint, reflection *nmVertex) *nmVertex {
	toScalar := reflection.subtract(midpoint)
	toScalar = toScalar.multiply(gamma)
	toScalar = midpoint.add(toScalar)
	return nm.evaluateWithConstraints(vertices, toScalar)
}

func (nm *nelderMead) insideContract(vertices vertices, midpoint, reflection *nmVertex) *nmVertex {
	toScalar := reflection.subtract(midpoint)
	toScalar = toScalar.multiply(gamma)
	toScalar = midpoint.subtract(toScalar)
	return nm.evaluateWithConstraints(vertices, toScalar)
}

func (nm *nelderMead) shrink(vertices vertices) {
	one := vertices[0]
	for i := 1; i < len(vertices); i++ {
		toScalar := vertices[i].subtract(one)
		toScalar = toScalar.multiply(sigma)
		vertices[i] = one.add(toScalar)
	}
}

// checkIteration checks some key values to determine if
// iteration should be complete.  Returns false if iteration
// should be terminated and true if iteration should continue.
func (nm *nelderMead) checkIteration(vertices vertices) bool {
	// this will never be true for += inf
	if math.Abs(vertices[0].result-nm.config.Target) < delta {
		return false
	}

	best := vertices[0]
	// here we are checking distance convergence.  If all vertices
	// are near convergence, that is they are all within some delta
	// from the expected value, we can go ahead and quit early.  This
	// can only be performed on convergence checks, not for finding
	// min/max.
	if !isInf(nm.config.Target) {
		for _, v := range vertices[1:] {
			if math.Abs(best.distance-v.distance) >= delta {
				return true
			}
		}
	}

	// next we want to check to see if the changes in our polytopes
	// dip below some threshold.  That is, we want to look at the
	// euclidean distances between the best guess and all the other
	// guesses to see if they are converged upon some point.  If
	// all of the vertices have converged close enough, it may be
	// worth it to cease iteration.
	for _, v := range vertices[1:] {
		if best.euclideanDistance(v) >= delta {
			return true
		}
	}

	return false
}

func (nm *nelderMead) evaluate() {
	vertices := nm.results.grab(len(nm.config.Vars) + 1)
	// if the initial guess provided is not good, then
	// we are going to die early, leave it up to the user
	// to create a good first guess.
	vertices[0].evaluate(nm.config)
	if !vertices[0].good {
		nm.results.insert(vertices[0])
		return
	}

	// the outer loop controls how hard we try to find
	// a global critical point
	for i := 0; i < maxIterations; i++ {
		// the inner loop controls the degenerate case where
		// we can't converge to a critical point
		for j := 0; j < maxRuns; j++ {
			// TODO: optimize this to prevent duplicate evaluations.
			vertices.evaluate(nm.config)
			best := vertices[0]
			if !nm.checkIteration(vertices) {
				break
			}

			midpoint := findMidpoint(vertices[:len(vertices)-1]...)
			// we are guaranteed to have two points here
			reflection := nm.reflect(vertices, midpoint)
			// we could not find a reflection that met constraints, the
			// best guess is the best guess.
			if reflection == best {
				break
			}
			// in this case, quality-wise, we are between the best
			// and second to best points
			if reflection.less(nm.config, nm.lastDimensionVertex(vertices)) &&
				!vertices[0].less(nm.config, reflection) {

				vertices[len(vertices)-1] = reflection
			}

			// midpoint is closer than our previous best guess
			if reflection.less(nm.config, vertices[0]) {
				expanded := nm.expand(vertices, midpoint, reflection)
				// we could not expand a valid guess, best is the best guess
				if expanded == best {
					break
				}

				// we only need to expand here
				if expanded.less(nm.config, reflection) {
					vertices[len(vertices)-1] = expanded
				} else {
					vertices[len(vertices)-1] = reflection
				}
				continue
			}

			// reflection is a bad guess, let's try to contract both
			// inside and outside and see if we can find a better value
			if reflection.less(nm.config, nm.lastVertex(vertices)) {
				oc := nm.outsideContract(vertices, midpoint, reflection)
				if oc == best {
					break
				}
				if oc.less(nm.config, reflection) || oc.equal(nm.config, reflection) {
					vertices[len(vertices)-1] = oc
					continue
				}
			} else if !reflection.less(nm.config, nm.lastVertex(vertices)) {
				ic := nm.insideContract(vertices, midpoint, reflection)
				if ic == best {
					break
				}
				if ic.less(nm.config, nm.lastVertex(vertices)) {
					vertices[len(vertices)-1] = ic
					continue
				}
			}

			// we could not guess a better value than nm.vertices[0], so
			// let's converge the other to guesses to our best guess.
			nm.shrink(vertices)
		}
		nm.results.reSort(vertices[0])
		vertices = nm.results.grab(len(nm.config.Vars) + 1)
	}
}

func newNelderMead(config NelderMeadConfiguration) *nelderMead {
	v := &nmVertex{vars: config.Vars}      // construct initial vertex with first guess
	results := newResults(v, config, 1000) // 1000 represents 1000 initial vertex guesses

	return &nelderMead{
		config:  config,
		results: results,
	}
}

// NelderMead takes a configuration and returns a list
// of floats that can be plugged into the provided function
// to converge at the target value.
func NelderMead(config NelderMeadConfiguration) []float64 {
	nm := newNelderMead(config)
	nm.evaluate()
	return nm.results.vertices[0].vars
}
