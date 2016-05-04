package goiso

/*
  Copyright (c) 2014 Tim Henderson
  Released under the GNU General Public License version 3.

  This file is part of goiso a wrapper around bliss.

  bliss is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License version 3
  as published by the Free Software Foundation.

  bliss is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with goiso.  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"strings"
)

import (
	"github.com/timtadh/goiso/bliss"
)

type Graph struct {
	V         Vertices
	E         Edges
	Kids      [][]*Edge
	Parents   [][]*Edge
	Colors    []string
	Labels    map[string]int
	colorFreq []int
	closed    bool
	canon     bool
	blissMap  *bliss.Map
}

type Vertices []Vertex
type Edges []Edge

type SubGraph struct {
	V           Vertices
	E           Edges
	Kids        [][]*Edge
	Parents     [][]*Edge
	G           *Graph
	vertexIndex map[int]*Vertex
	edgeIndex   map[ColoredArc]*Edge
}

type Vertex struct {
	Idx   int
	Id    int
	Color int
}

func (v *Vertex) Copy(idx int) Vertex {
	return Vertex{
		Idx:   idx,
		Id:    v.Id,
		Color: v.Color,
	}
}

type Arc struct {
	Src, Targ int
}

type Edge struct {
	Arc
	Idx   int
	Color int
}

func (e *Edge) Copy(idx, src, targ int) Edge {
	return Edge{
		Arc: Arc{
			Src:  src,
			Targ: targ,
		},
		Idx:   idx,
		Color: e.Color,
	}
}

func (V Vertices) Iterate() (vi bliss.VertexIterator) {
	i := 0
	vi = func() (color int, _ bliss.VertexIterator) {
		if i >= len(V) {
			return 0, nil
		}
		color = V[i].Color
		i++
		return color, vi
	}
	return vi
}

func (E Edges) Iterate() (ei bliss.EdgeIterator) {
	i := 0
	ei = func() (src, targ, color int, _ bliss.EdgeIterator) {
		if i >= len(E) {
			return 0, 0, 0, nil
		}
		src = E[i].Src
		targ = E[i].Targ
		color = E[i].Color
		i++
		return src, targ, color, ei
	}
	return ei
}

// Construct a new graph with V vertices and E edges.
func NewGraph(V, E int) Graph {
	return Graph{
		V:        make([]Vertex, 0, V),
		E:        make([]Edge, 0, E),
		Kids:     make([][]*Edge, 0, V),
		Parents:  make([][]*Edge, 0, V),
		Colors:   make([]string, 0, V),
		Labels: make(map[string]int, V),
	}
}

// Construct a subgraph. The vids are the vertices you are including.
// The filter_edges, are the edge labels you would like to ignore (can
// be nil). Note: these are the indexes into V not the vertex Ids. Also
// note: this subgraph will always be canonicalized! Finally: the Ids in
// the vertex will be not be the original Id on the graph but rather the
// Idx to the vertex on the original graph. This allows you to easily
// recover the embedding.
func (g *Graph) SubGraph(vids []int, filtered_edges map[string]bool) (sg *SubGraph, canonized bool) {
	V := g.find_vertices(vids)
	E := g.find_edges(vids, V, filtered_edges)
	return canonSubGraph(g, V, E)
}

func (g *Graph) VertexSubGraph(vid int) (sg *SubGraph, canonized bool) {
	V := g.find_vertices([]int{vid})
	return canonSubGraph(g, V, []Edge{})
}

func (g *Graph) EmptySubGraph() (sg *SubGraph, canonized bool) {
	return canonSubGraph(g, []Vertex{}, []Edge{})
}

func (g *Graph) find_vertices(vids []int) []Vertex {
	V := make([]Vertex, 0, len(vids))
	for i, vid := range vids {
		v := g.V[vid].Copy(i)
		v.Id = vid
		V = append(V, v)
	}
	return V
}

func (g *Graph) find_edges(vids []int, V []Vertex, filtered_edges map[string]bool) []Edge {
	vset := make(map[int]int)
	for i, vid := range vids {
		vset[vid] = i
	}
	edges := make([]Edge, 0, len(vids))
	for i, u := range vids {
		for _, e := range g.Kids[u] {
			if _, has := filtered_edges[g.Colors[e.Color]]; has {
				continue
			}
			if j, has := vset[e.Targ]; has {
				edges = append(edges, e.Copy(len(edges), (V)[i].Idx, (V)[j].Idx))
			}
		}
	}
	return edges
}

func safe_label(label string) string {
	label = strings.Replace(label, ":", "\\:", -1)
	label = strings.Replace(label, "(", "\\(", -1)
	label = strings.Replace(label, ")", "\\)", -1)
	label = strings.Replace(label, "[", "\\[", -1)
	label = strings.Replace(label, "]", "\\]", -1)
	label = strings.Replace(label, "->", "\\-\\>", -1)
	return label
}

// This is a short string useful as a unique (after canonicalization)
// label for the graph.
func (g *Graph) Label() string {
	V := make([]string, 0, len(g.V))
	E := make([]string, 0, len(g.E))
	for _, v := range g.V {
		V = append(V, fmt.Sprintf(
			"(%v:%v)",
			v.Idx,
			safe_label(g.Colors[v.Color]),
		))
	}
	for _, e := range g.E {
		E = append(E, fmt.Sprintf(
			"[%v->%v:%v]",
			e.Src,
			e.Targ,
			safe_label(g.Colors[e.Color]),
		))
	}
	return fmt.Sprintf("%d:%d%v%v", len(g.E), len(g.V), strings.Join(V, ""), strings.Join(E, ""))
}

// Stringifies the graph. This produces a String in the graphviz dot
// language.
func (g *Graph) String() string {
	V := make([]string, 0, len(g.V))
	E := make([]string, 0, len(g.E))
	for _, v := range g.V {
		V = append(V, fmt.Sprintf(
			"%v [label=\"%v\"];",
			v.Id,
			g.Colors[v.Color],
		))
	}
	for _, e := range g.E {
		E = append(E, fmt.Sprintf(
			"%v -> %v [label=\"%v\"];",
			g.V[e.Src].Id,
			g.V[e.Targ].Id,
			g.Colors[e.Color],
		))
	}
	return fmt.Sprintf(
		`digraph {
    %v
    %v
}
`, strings.Join(V, "\n    "), strings.Join(E, "\n    "))
}

// Finalize the graph. Once this method is called, edges and vertices
// can no longer be added. The reason is simple, the mapping between
// this graph and the graph is bliss has been constructed.
func (g *Graph) Finalize() {
	g.closed = true
	g.blissMap = bliss.NewMap(len(g.V), len(g.E), g.V.Iterate(), g.E.Iterate())
}

func (g *Graph) Canonized() bool {
	return g.canon
}

// Creates a new graph which is the canonical representation. This
// method does cause the graph to become finizalized as it makes use of
// CanonicalPermutation.
func (g *Graph) Canonical() (ng Graph, canonized bool) {
	ng = Graph{
		V:        make([]Vertex, len(g.V)),
		E:        make([]Edge, len(g.E)),
		Kids:     make([][]*Edge, len(g.Kids)),
		Parents:  make([][]*Edge, len(g.Parents)),
		Colors:   make([]string, len(g.Colors)),
		Labels: make(map[string]int),
		closed:   true,
		canon:    true,
	}
	copy(ng.Colors, g.Colors)
	for cid, color := range ng.Colors {
		ng.Labels[color] = cid
	}
	for i := range ng.Kids {
		ng.Kids[i] = make([]*Edge, 0, 5)
	}
	for i := range ng.Parents {
		ng.Parents[i] = make([]*Edge, 0, 5)
	}
	vord, eord, canonized := g.CanonicalPermutation()
	// i is the old vid, j is the new vid
	for i, j := range vord {
		ng.V[j] = g.V[i].Copy(j)
	}
	for i, j := range eord {
		ng.E[j] = g.E[i].Copy(j, vord[g.E[i].Src], vord[g.E[i].Targ])
		ng.Kids[vord[g.E[i].Src]] = append(ng.Kids[vord[g.E[i].Src]], &ng.E[j])
		ng.Parents[vord[g.E[i].Targ]] = append(ng.Parents[vord[g.E[i].Targ]], &ng.E[j])
	}
	return ng, canonized
}

// Computes the canonical (labeling) permutation of the graph. Vord is
// the mapping from vid->new-vid. Eord is eid->new-eid. Unless you need
// something special you probably just want to use Canonical(). canonized
// is true if the graph is already in canonical form.
// Note: this method does finalize the graph as it calls into bliss.
func (g *Graph) CanonicalPermutation() (Vord, Eord []int, canonized bool) {
	if !g.closed {
		g.Finalize()
	}
	return g.blissMap.CanonicalPermutation()
}

// Adds a vertex. The id is not used by this package but is preserved.
// The purpose is for you to track the identity of each vertex. The
// label is the label of the vertex.
func (g *Graph) AddVertex(id int, label string) *Vertex {
	if g.closed {
		return nil
	}
	v := Vertex{
		Idx:   len(g.V),
		Id:    id,
		Color: g.addColor(label),
	}
	g.V = append(g.V, v)
	g.Kids = append(g.Kids, make([]*Edge, 0, 5))
	g.Parents = append(g.Parents, make([]*Edge, 0, 5))
	return &v
}

// Adds and edge. The label is the label on the edge.
func (g *Graph) AddEdge(u, v *Vertex, label string) *Edge {
	if g.closed {
		return nil
	}
	e := Edge{
		Arc: Arc{
			Src:  u.Idx,
			Targ: v.Idx,
		},
		Idx:   len(g.E),
		Color: g.addColor(label),
	}
	g.E = append(g.E, e)
	g.Kids[e.Arc.Src] = append(g.Kids[e.Arc.Src], &e)
	g.Parents[e.Arc.Targ] = append(g.Parents[e.Arc.Targ], &e)
	return &e
}

// What is the frequency of a color?
func (g *Graph) ColorFrequency(color int) int {
	return g.colorFreq[color]
}

func (g *Graph) addColor(label string) int {
	if cid, has := g.Labels[label]; has {
		g.colorFreq[cid] += 1
		return cid
	}
	cid := len(g.Colors)
	g.Labels[label] = cid
	g.Colors = append(g.Colors, label)
	g.colorFreq = append(g.colorFreq, 1)
	return cid
}
