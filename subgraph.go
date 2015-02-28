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
	"encoding/binary"
	"fmt"
	"strings"
)

func canonSubGraph(g *Graph, V *[]Vertex, E *[]Edge) *SubGraph {
	if len(*V) == 1 && len(*E) == 0 {
		sg := &SubGraph{
			V: *V,
			E: *E,
			Kids: make([][]*Edge, len(*V)),
			G: g,
			idIndex: make(map[int]*Vertex),
		}
		sg.Kids[0] = make([]*Edge, 0)
		sg.idIndex[sg.V[0].Id] = &sg.V[0]
		return sg
	}
	bMap := makeBlissMap(V, E)
	sg := &SubGraph{
		V:       make([]Vertex, len(*V)),
		E:       make([]Edge, len(*E)),
		Kids:    make([][]*Edge, len(*V)),
		G:       g,
		idIndex: make(map[int]*Vertex),
	}
	for i := range sg.Kids {
		sg.Kids[i] = make([]*Edge, 0, 5)
	}
	vord, eord := bMap.canonicalPermutation(len(*V), len(*E))
	// i is the old vid, j is the new vid
	for i, j := range vord {
		sg.V[j] = (*V)[i].Copy(j)
		sg.idIndex[sg.V[j].Id] = &sg.V[j]
	}
	for i, j := range eord {
		sg.E[j] = (*E)[i].Copy(j, vord[(*E)[i].Src], vord[(*E)[i].Targ])
		sg.Kids[vord[(*E)[i].Src]] = append(sg.Kids[vord[(*E)[i].Src]], &sg.E[j])
	}
	return sg
}

// This is a useful method for finding out if the subgraph has a vertex from the
// parent graph
func (sg *SubGraph) Has(id int) bool {
	_, has := sg.idIndex[id]
	return has
}

// Checks to see if these two subgraphs are isomorphic. It relies on the fact that
// subgraph are always stored in the cannonical ordering.
func (sg *SubGraph) Equals(o *SubGraph) bool {
	if len(sg.V) != len(o.V) {
		return false
	}
	if len(sg.E) != len(o.E) {
		return false
	}
	for i := range sg.V {
		if sg.V[i].Color != o.V[i].Color {
			return false
		}
	}
	for i := range sg.E {
		if sg.E[i].Color != o.E[i].Color {
			return false
		}
		if sg.V[sg.E[i].Src].Color != o.V[o.E[i].Src].Color {
			return false
		}
		if sg.V[sg.E[i].Targ].Color != o.V[o.E[i].Targ].Color {
			return false
		}
	}
	return true
}

// This will extend the current subgraph and return a new larger subgraph. Note:
// this will not modify the current subgraph in any way.
func (sg *SubGraph) Extend(vids ...int) *SubGraph {
	avids := make([]int, 0, len(sg.V)+len(vids))
	for _, v := range sg.V {
		avids = append(avids, v.Id)
	}
	for _, vid := range vids {
		avids = append(avids, vid)
	}
	return sg.G.SubGraph(avids, nil)
}

// See SubGraph.Serialize for the format
func DeserializeSubGraph(g *Graph, bytes []byte) *SubGraph {
	lenV := binary.LittleEndian.Uint32(bytes[0:4])
	lenE := binary.LittleEndian.Uint32(bytes[4:8])
	off := 8
	V := make([]Vertex, lenV)
	E := make([]Edge, lenE)
	kids := make([][]*Edge, len(V))
	idIndex := make(map[int]*Vertex)
	for i := range kids {
		kids[i] = make([]*Edge, 0, 5)
	}
	for i := 0; i < int(lenV); i++ {
		s := off + i*4
		e := s + 4
		id := int(binary.LittleEndian.Uint32(bytes[s:e]))
		v := Vertex{
			Idx: i,
			Id: id,
			Color: g.V[id].Color,
		}
		V[i] = v
		idIndex[v.Id] = &V[i]
	}
	off += len(V)*4
	for i := 0; i < int(lenE); i++ {
		s := off + i*12
		e := s + 4
		src := int(binary.LittleEndian.Uint32(bytes[s:e]))
		s += 4
		e += 4
		targ := int(binary.LittleEndian.Uint32(bytes[s:e]))
		s += 4
		e += 4
		color := int(binary.LittleEndian.Uint32(bytes[s:e]))
		edge := Edge{
			Arc: Arc{
				Src: src,
				Targ: targ,
			},
			Idx: i,
			Color: color,
		}
		E[i] = edge
		kids[E[i].Src] = append(kids[E[i].Src], &E[i])
	}
	return &SubGraph{
		V:       V,
		E:       E,
		Kids:    kids,
		G:       g,
		idIndex: idIndex,
	}
}

// format: (vertex count : 4)(edge count : 4)(vertex id : 4)+[edge (src idx : 4)(targ idx : 4)(label color : 4)]+
//
// vertices are in idx order.
// edges are in idx order.
// the order is the canonical order.
func (sg *SubGraph) Serialize() []byte {
	bytes := make([]byte, 8 + len(sg.V)*4 + len(sg.E)*12)
	binary.LittleEndian.PutUint32(bytes[0:4], uint32(len(sg.V)))
	binary.LittleEndian.PutUint32(bytes[4:8], uint32(len(sg.E)))
	off := 8
	for i, v := range sg.V {
		s := off + i*4
		e := s + 4
		binary.LittleEndian.PutUint32(bytes[s:e], uint32(v.Id)) // Idx in *Graph
	}
	off += len(sg.V)*4
	for i, edge := range sg.E {
		s := off + i*12
		e := s + 4
		binary.LittleEndian.PutUint32(bytes[s:e], uint32(edge.Src))
		s += 4
		e += 4
		binary.LittleEndian.PutUint32(bytes[s:e], uint32(edge.Targ))
		s += 4
		e += 4
		binary.LittleEndian.PutUint32(bytes[s:e], uint32(edge.Color))
	}
	return bytes
}

func (sg *SubGraph) ShortLabel() []byte {
	label := make([]byte, len(sg.V)*4 + len(sg.E)*12)
	for i, v := range sg.V {
		s := i*4
		e := s + 4
		binary.LittleEndian.PutUint32(label[s:e], uint32(v.Color))
	}
	off := len(sg.V)*4
	for i, edge := range sg.E {
		s := off + i*12
		e := s + 4
		binary.LittleEndian.PutUint32(label[s:e], uint32(edge.Src))
		s += 4
		e += 4
		binary.LittleEndian.PutUint32(label[s:e], uint32(edge.Targ))
		s += 4
		e += 4
		binary.LittleEndian.PutUint32(label[s:e], uint32(edge.Color))
	}
	return label
}

// This is a short string useful as a unique (after canonicalization)
// label for the graph.
func (sg *SubGraph) Label() string {
	V := make([]string, 0, len(sg.V))
	E := make([]string, 0, len(sg.E))
	for _, v := range sg.V {
		V = append(V, fmt.Sprintf(
			"(%v:%v)",
			v.Idx,
			safe_label(sg.G.Colors[v.Color]),
		))
	}
	for _, e := range sg.E {
		E = append(E, fmt.Sprintf(
			"[%v->%v:%v]",
			e.Src,
			e.Targ,
			safe_label(sg.G.Colors[e.Color]),
		))
	}
	return fmt.Sprintf("%v%v", strings.Join(V, ""), strings.Join(E, ""))
}

// Stringifies the graph. This produces a String in the graphviz dot
// language.
func (sg *SubGraph) String() string {
	V := make([]string, 0, len(sg.V))
	E := make([]string, 0, len(sg.E))
	for _, v := range sg.V {
		V = append(V, fmt.Sprintf(
			"%v [label=\"%v\"];",
			sg.G.V[v.Id].Id,
			sg.G.Colors[v.Color],
		))
	}
	for _, e := range sg.E {
		E = append(E, fmt.Sprintf(
			"%v -> %v [label=\"%v\"];",
			sg.G.V[sg.V[e.Src].Id].Id,
			sg.G.V[sg.V[e.Targ].Id].Id,
			sg.G.Colors[e.Color],
		))
	}
	return fmt.Sprintf(
		`digraph {
    %v
    %v
}
`, strings.Join(V, "\n    "), strings.Join(E, "\n    "))
}
