package daggo

import (
	"fmt"
	"sort"
)

// Attrs is vertice' attributes.
type Attrs interface{}

// Vertice is a vertice formed DAG.
type Vertice interface {
	ID() string
	Attrs() Attrs
}

// Vertices is a slice of vertices.
type Vertices []Vertice

func (s Vertices) Len() int           { return len(s) }
func (s Vertices) Less(i, j int) bool { return s[i].ID() < s[j].ID() }
func (s Vertices) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Sort returns a slice of vertices in increasing order by vertice ID.
func (s Vertices) Sort() Vertices {
	sort.Stable(s)
	return s
}

// IDs returns a slice of vertices' IDs
func (s Vertices) IDs() []string {
	res := make([]string, len(s))
	for i, v := range s {
		res[i] = v.ID()
	}
	return res
}

// Attrs returns a slice of vertices' Attrs
func (s Vertices) Attrs() []Attrs {
	res := make([]Attrs, len(s))
	for i, v := range s {
		res[i] = v.Attrs()
	}
	return res
}

// Filter ...
func (s Vertices) Filter(fn func(v Vertice) bool) Vertices {
	res := make([]Vertice, 0, len(s)/2)
	for _, v := range s {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}

// DAG is a directed acyclic graph.
type DAG struct {
	blocks map[string]*block
}

// New returns a new DAG.
func New() *DAG {
	return &DAG{
		blocks: make(map[string]*block),
	}
}

type block struct {
	vertice Vertice
	prev    map[string]int
	next    map[string]int
}

// Len returns vertices count in the DAG.
func (d *DAG) Len() int {
	return len(d.blocks)
}

// GetVertice returns a vertice in the DAG by ID, returns nil if not found.
func (d *DAG) GetVertice(id string) Vertice {
	block, ok := d.blocks[id]
	if ok {
		return block.vertice
	}
	return nil
}

// Vertices returns all vertices in the DAG.
func (d *DAG) Vertices() Vertices {
	res := make([]Vertice, 0, len(d.blocks))
	for _, b := range d.blocks {
		res = append(res, b.vertice)
	}
	return res
}

// StartingVertices returns starting vertices in the DAG that have no other vertices connected to them.
func (d *DAG) StartingVertices() Vertices {
	res := make([]Vertice, 0)
	for _, b := range d.blocks {
		if len(b.prev) == 0 {
			res = append(res, b.vertice)
		}
	}
	return res
}

// EndingVertices returns ending vertices in the DAG that don't connected to any other vertices.
func (d *DAG) EndingVertices() Vertices {
	res := make([]Vertice, 0)
	for _, b := range d.blocks {
		if len(b.next) == 0 && len(b.prev) != 0 {
			res = append(res, b.vertice)
		}
	}
	return res
}

// ToVertices returns vertices in the DAG that the vertice v connected to them.
func (d *DAG) ToVertices(v Vertice) Vertices {
	res := make([]Vertice, 0)
	if v == nil {
		return res
	}

	b, ok := d.blocks[v.ID()]
	if !ok {
		return res
	}

	for k := range b.next {
		res = append(res, d.blocks[k].vertice)
	}
	return res
}

// FromVertices returns vertices in the DAG that connected to the vertice v.
func (d *DAG) FromVertices(v Vertice) Vertices {
	res := make([]Vertice, 0)
	if v == nil {
		return res
	}

	b, ok := d.blocks[v.ID()]
	if !ok {
		return res
	}

	for k := range b.prev {
		res = append(res, d.blocks[k].vertice)
	}
	return res
}

// Equal asserts that two DAG are equal.
func (d *DAG) Equal(a *DAG) bool {
	if len(d.blocks) != len(a.blocks) {
		return false
	}
	for k, b := range d.blocks {
		x, ok := a.blocks[k]
		if !ok || x.vertice.ID() != b.vertice.ID() || len(x.prev) != len(b.prev) || len(x.next) != len(b.next) {
			return false
		}
		for id, w := range b.prev {
			xw, ok := x.prev[id]
			if !ok || xw != w {
				return false
			}
		}
		for id, w := range b.next {
			xw, ok := x.next[id]
			if !ok || xw != w {
				return false
			}
		}
	}
	return true
}

// AddEdge adds a connecting pairs of vertices into the DAG.
// the vertices should not be nil, not be equal, and not form a cyclic graph.
// the method can be called multiple times.
func (d *DAG) AddEdge(start, end Vertice, weight int) error {
	if start == nil || start.ID() == "" {
		return fmt.Errorf("invalid starting vertice: %#v", start)
	}
	if end == nil || end.ID() == "" {
		return fmt.Errorf("invalid ending vertice: %#v", end)
	}

	startID := start.ID()
	endID := end.ID()
	if startID == endID {
		return fmt.Errorf("starting vertice is ending vertice: %s", startID)
	}

	startBlock, ok1 := d.blocks[startID]
	if !ok1 {
		startBlock = &block{
			vertice: start,
			prev:    make(map[string]int),
			next:    make(map[string]int),
		}
		d.blocks[startID] = startBlock
	}

	endBlock, ok2 := d.blocks[endID]
	if !ok2 {
		endBlock = &block{
			vertice: end,
			prev:    make(map[string]int),
			next:    make(map[string]int),
		}
		d.blocks[endID] = endBlock
	}

	if ok1 && ok2 {
		if d.isReachable(endBlock, startBlock.vertice.ID()) {
			return fmt.Errorf("cyclic graph will come into being")
		}
	}
	startBlock.next[endBlock.vertice.ID()] = weight
	endBlock.prev[startBlock.vertice.ID()] = weight
	return nil
}

// RemoveEdge remove the direct connecting in the vertices pair.
func (d *DAG) RemoveEdge(start, end Vertice) {
	if start == nil || end == nil {
		return
	}

	startBlock, ok := d.blocks[start.ID()]
	if !ok {
		return
	}

	endBlock, ok := d.blocks[end.ID()]
	if !ok {
		return
	}

	delete(startBlock.next, endBlock.vertice.ID())
	delete(endBlock.prev, startBlock.vertice.ID())
}

// ReachDAG returns a new sub DAG with the most edges that starting vertice may reach to.
func (d *DAG) ReachDAG(start Vertice) *DAG {
	nd := New()
	if start == nil {
		return nd
	}

	startBlock, ok := d.blocks[start.ID()]
	if !ok {
		return nd
	}
	var iterator func(n *block)
	iterator = func(n *block) {
		for k, w := range n.next {
			b := d.blocks[k]
			nd.AddEdge(n.vertice, b.vertice, w)
			iterator(b)
		}
	}
	iterator(startBlock)
	return nd
}

// CloseDAG returns a new transitive closure DAG with the most edges that represents the same reachability relation.
func (d *DAG) CloseDAG(start, end Vertice) *DAG {
	nd := New()
	if start == nil {
		return nd
	}

	startBlock, ok := d.blocks[start.ID()]
	if !ok {
		return nd
	}
	if end == nil {
		return nd
	}

	endBlock, ok := d.blocks[end.ID()]
	if !ok {
		return nd
	}
	if startBlock == endBlock {
		return nd
	}

	var iterator func(n *block) bool
	iterator = func(n *block) bool {
		ok := false
		for k, w := range n.next {
			b := d.blocks[k]
			if b == endBlock || iterator(b) {
				nd.AddEdge(n.vertice, b.vertice, w)
				ok = true
			}
		}
		return ok
	}

	iterator(startBlock)
	return nd
}

// ReduceDAG returns a new transitive reduction DAG with the fewest edges that represents the same reachability relation.
func (d *DAG) ReduceDAG(start, end Vertice) *DAG {
	nd := d.CloseDAG(start, end)
	if nd.Len() == 0 {
		return nd
	}

	var iterator func(n *block)
	iterator = func(n *block) {
		target := n.vertice.ID()
		for k := range n.prev {
			b := nd.blocks[k]
			// try remove relation and check other relations
			delete(b.next, target)
			if nd.isReachable(b, target) {
				// clear relation
				delete(n.prev, k)
			} else {
				// fix relation
				b.next[target] = n.prev[k]
			}
			iterator(b)
		}
	}

	iterator(nd.blocks[end.ID()])
	return nd
}

// Reverse returns a new DAG that all edges relation reversed.
func (d *DAG) Reverse() *DAG {
	nd := New()
	for k, b := range d.blocks {
		nb := &block{
			vertice: b.vertice,
			prev:    make(map[string]int),
			next:    make(map[string]int),
		}
		nd.blocks[k] = nb
	}
	for _, b := range d.blocks {
		id := b.vertice.ID()
		for kk, w := range b.next {
			nd.blocks[kk].next[id] = w
		}
		for kk, w := range b.prev {
			nd.blocks[kk].prev[id] = w
		}
	}
	return nd
}

// Iterate iterate the DAG' vertices with the most reachability relation paths.
func (d *DAG) Iterate(start Vertice, init []Attrs, fn func(cur Vertice, weight int, acc []Attrs) []Attrs) []Attrs {
	res := make([]Attrs, 0)
	b := d.blocks[start.ID()]
	if b == nil {
		return res
	}

	var iterator func(b *block, weight int, acc []Attrs)
	iterator = func(b *block, weight int, acc []Attrs) {
		r := fn(b.vertice, weight, acc)
		if len(b.next) == 0 {
			res = append(res, r...)
			return
		}
		keys := make([]string, 0, len(b.next))
		for k := range b.next {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			iterator(d.blocks[k], b.next[k], r[:])
		}
	}
	if init == nil {
		init = make([]Attrs, 0)
	}
	iterator(b, 0, init)
	return res
}

type pathAcc struct {
	weight int
	paths  []*block
}

func (d *DAG) findPaths(start, end Vertice) []*pathAcc {
	res := make([]*pathAcc, 0)
	if start == nil || end == nil {
		return res
	}

	startBlock, ok := d.blocks[start.ID()]
	if !ok {
		return res
	}

	endBlock, ok := d.blocks[end.ID()]
	if !ok {
		return res
	}
	if startBlock == endBlock {
		return res
	}
	if !d.isReachable(startBlock, end.ID()) {
		return res
	}

	var iterator func(n *block) []*pathAcc
	iterator = func(n *block) []*pathAcc {
		res := make([]*pathAcc, 0)
		for k, w := range n.next {
			b := d.blocks[k]
			if b == endBlock {
				res = append(res, &pathAcc{weight: w, paths: []*block{b}})
				continue
			}

			for _, acc := range iterator(b) {
				acc.weight += w
				acc.paths = append(acc.paths, b)
				res = append(res, acc)
			}
		}
		return res
	}

	return iterator(startBlock)
}

// Shortest find a shortest paths.
func (d *DAG) Shortest(start, end Vertice, withWeight bool) Vertices {
	accs := d.findPaths(start, end)
	if len(accs) == 0 {
		return make([]Vertice, 0)
	}

	if withWeight {
		sort.SliceStable(accs, func(i, j int) bool { return accs[i].weight < accs[j].weight })
	} else {
		sort.SliceStable(accs, func(i, j int) bool { return len(accs[i].paths) < len(accs[j].paths) })
	}
	res := make([]Vertice, 0, len(accs[0].paths)+1)
	res = append(res, d.blocks[start.ID()].vertice)
	for i := len(accs[0].paths) - 1; i >= 0; i-- {
		res = append(res, accs[0].paths[i].vertice)
	}
	return res
}

// Longest find a longest paths.
func (d *DAG) Longest(start, end Vertice, withWeight bool) Vertices {
	accs := d.findPaths(start, end)
	if len(accs) == 0 {
		return make([]Vertice, 0)
	}

	if withWeight {
		sort.SliceStable(accs, func(i, j int) bool { return accs[i].weight > accs[j].weight })
	} else {
		sort.SliceStable(accs, func(i, j int) bool { return len(accs[i].paths) > len(accs[j].paths) })
	}
	res := make([]Vertice, 0, len(accs[0].paths)+1)
	res = append(res, d.blocks[start.ID()].vertice)
	for i := len(accs[0].paths) - 1; i >= 0; i-- {
		res = append(res, accs[0].paths[i].vertice)
	}
	return res
}

func (d *DAG) isReachable(x *block, target string) bool {
	for k := range x.next {
		if k == target {
			return true
		}
		if d.isReachable(d.blocks[k], target) {
			return true
		}
	}
	return false
}
