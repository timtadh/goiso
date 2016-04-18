package bliss

import (
	"sort"
)

type Map struct {
	LenV int // number of vertices in the original graph
	LenE int // number of edges in the original graph
	V    []Vertex
	E    []Arc
}

type Vertex struct {
	IsEdge bool
	Idx    int
	Color  int
}

type Arc struct {
	Src, Targ int
}

type VertexIterator func() (color int, vi VertexIterator)
type EdgeIterator func() (src, targ, color int, ei EdgeIterator)

type perm struct{ idx, p int }
type perms []perm

func (self perms) Len() int           { return len(self) }
func (self perms) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self perms) Less(i, j int) bool { return self[i].p < self[j].p }

func NewMap(lenV, lenE int, vi VertexIterator, ei EdgeIterator) *Map {
	V := make([]Vertex, 0, lenV+lenE)
	E := make([]Arc, 0, lenE*2)
	vidx := 0
	for color, vi := vi(); vi != nil; color, vi = vi() {
		V = append(V, Vertex{
			IsEdge: false,
			Idx:    vidx,
			Color:  color,
		})
		vidx += 1
	}
	eidx := 0
	for src, targ, color, ei := ei(); ei != nil; src, targ, color, ei = ei() {
		eid := len(V)
		V = append(V, Vertex{
			IsEdge: true,
			Idx:    eidx,
			Color:  color,
		})
		E = append(E, Arc{
			Src:  src,
			Targ: eid,
		})
		E = append(E, Arc{
			Src:  eid,
			Targ: targ,
		})
		eidx += 1
	}
	return &Map{LenV: lenV, LenE: lenE, V: V, E: E}
}

// Construct the CanonicalPermutation from the Map. The map itself is
// unchanged the permutation is given in Vord and Eord
//
// Vord [original-index] -> new-index of vertices
// Eord [original-index] -> new-index of edges
// canonized is true if the graph was already in canonical order
// canonized is false otherwise
func (bm *Map) CanonicalPermutation() (Vord, Eord []int, canonized bool) {
	nodes := make([]uint32, 0, len(bm.V))
	edges := make([]BlissEdge, 0, len(bm.E))
	for i := range bm.V {
		nodes = append(nodes, uint32(bm.V[i].Color))
	}
	for i := range bm.E {
		edges = append(edges, BlissEdge{uint32(bm.E[i].Src), uint32(bm.E[i].Targ)})
	}
	P := Canonize(nodes, edges)
	VP := make(perms, 0, bm.LenV)
	EP := make(perms, 0, bm.LenE)
	canonized = true
	for i, p := range P {
		if uint(i) != p {
			canonized = false
		}
		v := bm.V[i]
		if v.IsEdge {
			EP = append(EP, perm{v.Idx, int(p)})
		} else {
			VP = append(VP, perm{v.Idx, int(p)})
		}
	}
	sort.Sort(VP)
	sort.Sort(EP)
	Vord = make([]int, bm.LenV)
	Eord = make([]int, bm.LenE)
	for p, vp := range VP {
		if canonized && p != vp.idx {
			panic("not canonized when it should have been!")
		}
		Vord[vp.idx] = p
	}
	for p, ep := range EP {
		if canonized && p != ep.idx {
			panic("not canonized when it should have been!")
		}
		Eord[ep.idx] = p
	}
	return Vord, Eord, canonized
}

// Construct a BlissGraph from the Map
func (bm *Map) Graph() *Graph {
	bg := NewGraph(0)
	for _, v := range bm.V {
		bg.AddVertex(uint(v.Color))
	}
	for _, e := range bm.E {
		bg.AddEdge(uint(e.Src), uint(e.Targ))
	}
	return bg
}
