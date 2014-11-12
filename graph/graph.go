/*
Graph is used to store a pre-indxed (flattened) version of the
graph.  This is useful when attempting to flatten a large number
of nodes.  We can use some heuristics to determine if we should
reuse the pre-indexed positions or if it is beneficial to flatten
the nodes to find the optimal subgraph for execution.
*/
package graph

import (
	"log"
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/bitarray"
	"github.com/Workiva/go-datastructures/queue"
)

const (
	numNodesBeforeSplit int64 = 100
	goldenRatio               = .5
)

// bundle is a helper construct so that we bundle a node
// with its position in the flattened graph.
type bundle struct {
	INode
	position int64
}

// layer helps us order lists of nodes in a dense list of layers
type layer struct {
	position uint64
	nodes    Nodes
}

// layers is a helper struct that implements sort.Interface to help
// us sort a dense list of layers.
type layers []*layer

// Len returns the length of layers.
func (layers layers) Len() int {
	return len(layers)
}

// Less returns a bool indicating the relationship between
// items at the supplied positions.  IE, is layers[i] less than
// layers[j].
func (layers layers) Less(i, j int) bool {
	return layers[i].position < layers[j].position
}

// Swap will swap values at the provided positions.
func (layers layers) Swap(i, j int) {
	layers[i], layers[j] = layers[j], layers[i]
}

// positions stores a list of int64s which represent the layer
// of the subgraph where nodes can live.  The node's id represents
// the placement within the nodes slice (ie, position corresponds
// to node id).  A layer value of -1 indicates the formula is circular.
type positions []*bundle

// extractLayer returns a list of nodes that exist in the map of
// layer indices provided.
func (positions positions) extractLayers(layers map[int64]bool) Nodes {
	nodes := NewNodes()
	for _, node := range positions {
		if node == nil {
			continue
		}
		if _, ok := layers[node.position]; ok {
			nodes = append(nodes, node.INode)
		}
	}

	return nodes
}

// lowest returns the lowest index found in the group of nodes provided
func (positions positions) lowest(nodes Nodes) int64 {
	lowest := int64(math.MaxInt64)
	for _, node := range nodes {
		bundle := positions[node.ID()]
		if bundle == nil {
			continue
		}
		if bundle.position < lowest {
			lowest = bundle.position
		}
	}

	if lowest == int64(math.MaxInt64) {
		lowest = 0
	}

	return lowest
}

// returns the highest position seen in the index
func (positions positions) highestSeen() int64 {
	highest := int64(-1)
	for _, node := range positions {
		if node == nil {
			continue
		}

		if node.position > highest {
			highest = node.position
		}
	}

	return highest
}

// flatten will take nodes and return a flattened list that
// is safe to execute in parallel.  Because this method relies
// on pre-indexed positions, this subgraph may not be optimal.
func (positions positions) flatten(nodes Nodes) (layers, Nodes) {
	var (
		circulars    = make(Nodes, 0, 5)
		layerMap     = make(map[uint64]Nodes, 5)
		position     int64
		nodePosition uint64
	)

	for _, node := range nodes {
		if node == nil {
			continue
		}
		id := node.ID() // reducing number of non-inlined method calls
		if id >= uint64(len(positions)) {
			log.Printf(
				`Trying to flatten a node that hasn't 
						been added to the graph.  Node: %+v`, node,
			)
			continue
		}

		bundle := positions[id]
		if bundle == nil {
			continue
		}

		position = bundle.position
		if position == -1 {
			circulars = append(circulars, node)
			continue
		}

		// we've already confirmed position isn't negative so
		// this should be safe
		nodePosition = uint64(position)
		layerMap[nodePosition] = append(layerMap[nodePosition], node)
	}

	layers := make(layers, 0, len(layerMap))
	for position, nodes := range layerMap {
		layers = append(layers, &layer{
			position: position,
			nodes:    nodes,
		})
	}

	sort.Sort(layers)
	return layers, circulars
}

type Graph struct {
	positions positions
	maxLayer  int64
}

// GetSubgraph will return an execution graph that needs to be calculated
// based on the supplied nodes.
func (g *Graph) GetSubgraph(nodes Nodes) *ExecutionGraph {
	// a large proportion of the graph needs to be recalculated
	layers, circulars := g.positions.flatten(nodes)
	flattened := make([]Nodes, 0, len(layers))
	size := int64(len(circulars))
	for _, layer := range layers {
		flattened = append(flattened, layer.nodes)
		size += int64(len(layer.nodes))
	}

	return &ExecutionGraph{
		size:      size,
		toApply:   flattened,
		circulars: circulars,
	}
}

// GetLowestNodes will return a list of nodes that match
// the lowest layer of the nodes provided.
func (g *Graph) GetLowestNodes(nodes Nodes) Nodes {
	if len(nodes) == 0 {
		return nil
	}

	toReturn := make(Nodes, 0, len(nodes))
	isSet := false
	lowest := int64(0)

	for _, node := range nodes {
		nb := g.positions[node.ID()]
		if nb == nil {
			log.Printf(`Node flattened that hasn't been added: %+v`, node)
			continue
		}

		if !isSet {
			isSet = true
			toReturn = append(toReturn, node)
			lowest = nb.position
			continue
		}

		if nb.position == lowest {
			toReturn = append(toReturn, node)
		} else if nb.position < lowest {
			toReturn = toReturn[:0]
			toReturn = append(toReturn, node)
			lowest = nb.position
		}
	}

	return toReturn
}

// AddNodes will add the provided nodes to the flattened index
// of the graph and return an execution graph that is ready to
// be calculated.
func (g *Graph) AddNodes(dp IDependencyProvider, nodes Nodes) *ExecutionGraph {
	dependentNodes := dp.GetDependents(nodes)
	highest := nodes.Highest()

	// want to make sure we don't overflow here
	if highest >= uint64(len(g.positions)) {
		diff := highest - uint64(len(g.positions))
		g.positions = append(g.positions, make(positions, diff+1)...)
	}

	lowest := g.positions.lowest(dependentNodes)

	dependentNodes = append(dependentNodes, nodes...)
	flattened, circulars, size := flatten(dp, dependentNodes)
	g.insert(flattened, circulars, lowest)

	return &ExecutionGraph{
		size:      size,
		toApply:   flattened,
		circulars: circulars,
	}
}

func (g *Graph) findUniqueDependencies(dp IDependencyProvider, nodes Nodes) Nodes {
	if len(nodes) == 0 {
		return nil
	}

	ba := bitarray.NewSparseBitArray()
	for _, node := range nodes {
		ba.SetBit(node.ID())
	}
	ids := make([]uint64, 0, len(nodes))
	chunks := nodes.Split()
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for _, chunk := range chunks {
		go func(nodes Nodes) {
			for _, node := range nodes {
				deps := dp.GetDependencies(node)
				if deps == nil {
					continue
				}

				uints := deps.ToNums()
				lock.Lock()
				for _, id := range uints {
					if ok, _ := ba.GetBit(id); ok {
						continue
					}

					ids = append(ids, id)
					ba.SetBit(id)
				}
				lock.Unlock()
			}
			wg.Done()
		}(chunk)
	}

	wg.Wait()

	deps := make(Nodes, 0, len(ids))
	for _, id := range ids {
		nb := g.positions[id]
		if nb == nil {
			continue
		}

		deps = append(deps, nb.INode)
	}

	return deps
}

// RemoveNodes will remove the provided nodes from the graph.
func (g *Graph) RemoveNodes(dp IDependencyProvider, nodes Nodes) *ExecutionGraph {
	if len(nodes) == 0 {
		return &ExecutionGraph{}
	}

	lowest := g.positions.lowest(nodes)
	for _, node := range nodes {
		g.positions[node.ID()] = nil // kill these first
	}

	maxLayer := g.maxLayer
	if maxLayer == -1 {
		maxLayer = 1
	}

	toExtract := make(map[int64]bool, maxLayer)
	if lowest == -1 { // we end up reordering the whole damn thing
		toExtract[-1] = true
		lowest = 0
	}

	for i := lowest; i <= g.maxLayer; i++ {
		toExtract[i] = true
	}

	nodesToFlatten := g.positions.extractLayers(toExtract)
	flattened, circulars, size := flatten(dp, nodesToFlatten)

	g.insert(flattened, circulars, lowest)
	// have to set new max layer, only in the delete case
	g.maxLayer = g.positions.highestSeen()

	return &ExecutionGraph{
		toApply:   flattened,
		circulars: circulars,
		size:      size,
	}
}

func (g *Graph) insert(flattened []Nodes, circulars Nodes, offset int64) {
	maxLayer := int64(-1)
	for i, nodes := range flattened {
		for _, node := range nodes {
			g.positions[node.ID()] = &bundle{
				INode:    node,
				position: int64(i) + offset,
			}
		}
		maxLayer = int64(i) + offset
	}

	for _, node := range circulars {
		g.positions[node.ID()] = &bundle{
			INode: node, position: -1,
		}
	}

	if maxLayer > g.maxLayer {
		g.maxLayer = maxLayer
	}
}

// FromNodes will create a new graph from the given nodes.
func FromNodes(dp IDependencyProvider, nodes Nodes) *Graph {
	if len(nodes) == 0 {
		return &Graph{}
	}

	highest := nodes.Highest()
	positions := make(positions, highest+1)
	g := &Graph{
		positions: positions,
		maxLayer:  -1, // it's possible they are all circular
	}
	flattened, circulars, _ := flatten(dp, nodes)

	g.insert(flattened, circulars, 0)

	return g
}

func flatten(dp IDependencyProvider, nodes Nodes) ([]Nodes, Nodes, int64) {
	if len(nodes) == 0 {
		return nil, nil, 0
	}

	ba := bitarray.NewBitArray(dp.MaxNode(), true)
	remaining := make(Nodes, 0, len(nodes))
	helper := make(Nodes, len(nodes))
	for _, node := range nodes {
		if node != nil {
			ba.ClearBit(node.ID())
			remaining = append(remaining, node)
		}
	}

	var circulars Nodes
	results := make([]Nodes, 0, 5) // 5 is a guess
	layer := make([]INode, 0, len(nodes))
	size := int64(0)

	var wg sync.WaitGroup
	var lock sync.Mutex

	for {
		if len(remaining) == 0 && len(layer) == 0 {
			break
		}

		q := queue.New(int64(len(remaining)))
		for _, node := range remaining {
			q.Put(node)
		}

		completed := uint64(0)
		required := uint64(len(remaining))

		wg.Add(1)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					items, err := q.Get(1)
					if err != nil { // queue was disposed
						break
					}

					node := items[0].(INode)
					okToAdd := true
					deps := dp.GetDependencies(node)
					if deps != nil && !ba.Intersects(deps) {
						okToAdd = false
					}

					if okToAdd {
						if node.IsCircular() {
							node.SetCircular(false)
						}
						lock.Lock()
						layer = append(layer, node)
						lock.Unlock()
					}

					if atomic.AddUint64(&completed, 1) == required {
						wg.Done()
						break
					}
				}
			}()
		}

		wg.Wait()
		q.Dispose()

		if len(layer) == 0 { // this is a problem, nothing new found
			circulars = make(Nodes, 0, len(remaining))
			for _, n := range remaining {
				n.SetCircular(true)
				circulars = append(circulars, n)
				size++
			}
			break
		}

		for _, node := range layer {
			size++
			ba.SetBit(node.ID())
		}

		results = append(results, layer)
		layer = make([]INode, 0, len(layer))
		copy(helper, remaining)
		remaining = remaining[:0]
		for i := 0; i < len(helper); i++ {
			node := helper[i]
			if node == nil {
				break
			}
			if ok, _ := ba.GetBit(node.ID()); ok {
				continue
			}

			remaining = append(remaining, node)
			helper[i] = nil
		}

		helper = helper[:len(remaining)]
	}

	return results, circulars, size
}

func New() *Graph {
	return &Graph{
		maxLayer: -1,
	}
}
