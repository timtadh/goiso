package bliss

/*
  Copyright (c) 2016 Tim Henderson
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
  along with .  If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"sort"
)

// Maps a labeled digraph with edge and vertex labels to a  digraph with just
// vertex labels. Bliss does not support edge labels so using a mapping is
// necessary in order to canonically order edge labeled graphs.
type Map struct {
	LenV      int         // number of vertices in the original graph
	LenE      int         // number of edges in the original graph
	FirstEdge int         // index of the first vertex representing an edge
	Nodes     []uint32    // The colors of each mapped vertex/edge
	Edges     []BlissEdge // Mapped edges
}

type VertexIterator func() (color int, vi VertexIterator)
type EdgeIterator func() (src, targ, color int, ei EdgeIterator)

type perm struct{ idx, p int }
type perms []perm

func (self perms) Len() int           { return len(self) }
func (self perms) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self perms) Less(i, j int) bool { return self[i].p < self[j].p }

// Construct the Mapping from the number of vertices and edges plus iterators
// over the vertices and edges.
func NewMap(lenV, lenE int, vi VertexIterator, ei EdgeIterator) *Map {
	nodes := make([]uint32, 0, lenV+lenE)
	edges := make([]BlissEdge, 0, lenE*2)
	for color, vi := vi(); vi != nil; color, vi = vi() {
		nodes = append(nodes, uint32(color))
	}
	firstEdge := len(nodes)
	for src, targ, color, ei := ei(); ei != nil; src, targ, color, ei = ei() {
		eid := uint32(len(nodes))
		nodes = append(nodes, uint32(color))
		edges = append(
			edges,
			BlissEdge{Src: uint32(src), Targ: eid},
			BlissEdge{Src: eid, Targ: uint32(targ)},
		)
	}
	return &Map{
		LenV:      lenV,
		LenE:      lenE,
		FirstEdge: firstEdge,
		Nodes:     nodes,
		Edges:     edges,
	}
}

// Construct the CanonicalPermutation from the Map. The map itself is
// unchanged the permutation is given in Vord and Eord
//
// Vord [original-index] -> new-index of vertices
// Eord [original-index] -> new-index of edges
// canonized is true if the graph was already in canonical order
// canonized is false otherwise
func (m *Map) CanonicalPermutation() (Vord, Eord []int, canonized bool) {
	P := Canonize(m.Nodes, m.Edges)
	VP := make(perms, 0, m.LenV)
	EP := make(perms, 0, m.LenE)
	canonized = true
	for i, p := range P {
		if uint(i) != p {
			canonized = false
		}
		if i < m.FirstEdge {
			// it is a vertex
			VP = append(VP, perm{i, int(p)})
		} else {
			// it is an edge
			EP = append(EP, perm{i - m.FirstEdge, int(p)})
		}
	}
	sort.Sort(VP)
	sort.Sort(EP)
	Vord = make([]int, m.LenV)
	Eord = make([]int, m.LenE)
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

// Construct a BlissDigraph from the Map
func (m *Map) Digraph() *Digraph {
	bg := NewDigraph(0)
	for _, color := range m.Nodes {
		bg.AddVertex(uint(color))
	}
	for _, e := range m.Edges {
		bg.AddEdge(uint(e.Src), uint(e.Targ))
	}
	return bg
}
