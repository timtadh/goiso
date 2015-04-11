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
			Parents: make([][]*Edge, len(*V)),
			G: g,
			vertexIndex: make(map[int]*Vertex),
			edgeIndex: make(map[Arc]*Edge),
		}
		sg.Kids[0] = make([]*Edge, 0)
		sg.Parents[0] = make([]*Edge, 0)
		sg.vertexIndex[sg.V[0].Id] = &sg.V[0]
		return sg
	}
	bMap := makeBlissMap(V, E)
	sg := &SubGraph{
		V:       make([]Vertex, len(*V)),
		E:       make([]Edge, len(*E)),
		Kids:    make([][]*Edge, len(*V)),
		Parents: make([][]*Edge, len(*V)),
		G:       g,
		vertexIndex: make(map[int]*Vertex),
		edgeIndex: make(map[Arc]*Edge),
	}
	for i := range sg.Kids {
		sg.Kids[i] = make([]*Edge, 0, 5)
	}
	for i := range sg.Parents {
		sg.Parents[i] = make([]*Edge, 0, 5)
	}
	vord, eord := bMap.canonicalPermutation(len(*V), len(*E))
	// i is the old vid, j is the new vid
	for i, j := range vord {
		sg.V[j] = (*V)[i].Copy(j)
		sg.vertexIndex[sg.V[j].Id] = &sg.V[j]
	}
	for i, j := range eord {
		sg.E[j] = (*E)[i].Copy(j, vord[(*E)[i].Src], vord[(*E)[i].Targ])
		sg.Kids[vord[(*E)[i].Src]] = append(sg.Kids[vord[(*E)[i].Src]], &sg.E[j])
		sg.Parents[vord[(*E)[i].Targ]] = append(sg.Parents[vord[(*E)[i].Targ]], &sg.E[j])
		idArc := Arc{sg.V[sg.E[j].Src].Id, sg.V[sg.E[j].Targ].Id}
		sg.edgeIndex[idArc] = &sg.E[j]
	}
	return sg
}

// This is a useful method for finding out if the subgraph has a
// vertex from the parent graph
func (sg *SubGraph) HasVertex(id int) bool {
	_, has := sg.vertexIndex[id]
	return has
}

func (sg *SubGraph) HasEdge(a Arc, color int) bool {
	e, has := sg.edgeIndex[a]
	if !has {
		return false
	}
	return e.Color == color
}

// Checks to see if these two subgraphs are isomorphic. It relies on
// the fact that subgraph are always stored in the cannonical
// ordering.
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

// This will extend the current subgraph and return a new larger
// subgraph. This does "vertex" extension. It adds a vertices listed by
// G.Idx in vids to the extension and all edges contained in the parent
// graph. If you want to add an edge at a time use EdgeExtend.
// Note: this will not modify the current subgraph in any way.
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

// This will extend the current subgraph with the given edge. Only the
// Arc and Color attributes of the edge are used. The Idx attribute is
// ignored. The edge.Arc.Src must be in the SubGraph, edge.Arg.Targ
// does not have to be in the subgraph. If it is not already there it
// will be added. The Src and Targ should contain the Idx of the
// vertices in the original graph. (This becomes the Id field in the
// SubGraph).
func (sg *SubGraph) EdgeExtend(edge *Edge) *SubGraph {
	avids := make([]int, 0, len(sg.V) + 1)
	hasTarg := false
	hasSrc := false
	var src int = -1
	var targ int = -1
	for _, v := range sg.V {
		avids = append(avids, v.Id)
		if edge.Src == v.Id {
			hasSrc = true
			src = v.Idx
		}
		if edge.Targ == v.Id {
			hasTarg = true
			targ = v.Idx
		}
	}
	if !hasTarg && !hasSrc {
		panic(fmt.Errorf("both Src or Targ not in graph"))
	}
	if !hasTarg {
		targ = len(avids)
		avids = append(avids, edge.Targ)
	}
	if !hasSrc {
		src = len(avids)
		avids = append(avids, edge.Src)
	}
	if src == -1 || targ == -1 {
		panic(fmt.Errorf("Src or Targ not in graph"))
	}
	V := sg.G.find_vertices(avids)
	E := make([]Edge, 0, len(sg.E) + 1)
	for _, e := range sg.E {
		E = append(E, e.Copy(len(E), e.Src, e.Targ))
	}
	E = append(E, Edge{
		Arc: Arc{
			Src: src,
			Targ: targ,
		},
		Idx: len(E),
		Color: edge.Color,
	})
	return canonSubGraph(sg.G, &V, &E)
}

// Removes the edge at the given idx and if necessary an attached
// vertex. It returns a new subgraph which has been canonicalized. If
// the graph only has two vertices and one edge it will return a graph
// with only the Src of the edge. The target will be dropped.
func (sg *SubGraph) RemoveEdge(edgeIdx int) *SubGraph {
	rmSrc := true
	rmTarg := true
	edge := &sg.E[edgeIdx]
	if len(sg.E) == 1 && len(sg.V) == 2 {
		return sg.G.SubGraph([]int{sg.V[edge.Src].Id}, nil)
	}
	for _, e := range sg.Kids[edge.Src] {
		if e == edge {
			continue
		}
		rmSrc = false
	}
	for _, e := range sg.Parents[edge.Src] {
		if e == edge {
			continue
		}
		rmSrc = false
	}
	for _, e := range sg.Kids[edge.Targ] {
		if e == edge {
			continue
		}
		rmTarg = false
	}
	for _, e := range sg.Parents[edge.Targ] {
		if e == edge {
			continue
		}
		rmTarg = false
	}
	if rmSrc && rmTarg {
		panic("would have removed both source and target")
	}
	rmV := rmSrc || rmTarg
	var rmVidx int
	if rmSrc {
		rmVidx = edge.Src
	}
	if rmTarg {
		rmVidx = edge.Targ
	}
	adjustIdx := func(idx int) int {
		if rmV && idx > rmVidx {
			return idx - 1
		}
		return idx
	}
	avids := make([]int, 0, len(sg.V))
	for idx, v := range sg.V {
		if rmV && rmVidx == idx {
			continue
		}
		avids = append(avids, v.Id)
	}
	V := sg.G.find_vertices(avids)
	E := make([]Edge, 0, len(sg.E) + 1)
	for idx, e := range sg.E {
		if idx == edgeIdx {
			continue
		}
		E = append(E, e.Copy(len(E), adjustIdx(e.Src), adjustIdx(e.Targ)))
	}
	return canonSubGraph(sg.G, &V, &E)
}

// See SubGraph.Serialize for the format
func DeserializeSubGraph(g *Graph, bytes []byte) *SubGraph {
	lenV := binary.LittleEndian.Uint32(bytes[0:4])
	lenE := binary.LittleEndian.Uint32(bytes[4:8])
	off := 8
	V := make([]Vertex, lenV)
	E := make([]Edge, lenE)
	kids := make([][]*Edge, len(V))
	parents := make([][]*Edge, len(V))
	vertexIndex := make(map[int]*Vertex)
	edgeIndex := make(map[Arc]*Edge)
	for i := range kids {
		kids[i] = make([]*Edge, 0, 5)
	}
	for i := range parents {
		parents[i] = make([]*Edge, 0, 5)
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
		vertexIndex[v.Id] = &V[i]
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
		parents[E[i].Targ] = append(parents[E[i].Targ], &E[i])
		idArc := Arc{V[E[i].Src].Id, V[E[i].Targ].Id}
		edgeIndex[idArc] = &E[i]
	}
	return &SubGraph{
		V:       V,
		E:       E,
		Kids:    kids,
		Parents: parents,
		G:       g,
		vertexIndex: vertexIndex,
		edgeIndex: edgeIndex,
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
	size := 8 + len(sg.V)*4 + len(sg.E)*12
	label := make([]byte, size)
	binary.BigEndian.PutUint32(label[0:4], uint32(len(sg.E)))
	binary.BigEndian.PutUint32(label[4:8], uint32(len(sg.V)))
	off := 8
	for i, v := range sg.V {
		s := off + i*4
		e := s + 4
		binary.BigEndian.PutUint32(label[s:e], uint32(v.Color))
	}
	off += len(sg.V)*4
	for i, edge := range sg.E {
		s := off + i*12
		e := s + 4
		binary.BigEndian.PutUint32(label[s:e], uint32(edge.Src))
		s += 4
		e += 4
		binary.BigEndian.PutUint32(label[s:e], uint32(edge.Targ))
		s += 4
		e += 4
		binary.BigEndian.PutUint32(label[s:e], uint32(edge.Color))
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
			sg.G.Colors[v.Color],
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
	return fmt.Sprintf("%v%v%v%v", len(sg.E), len(sg.V), strings.Join(V, ""), strings.Join(E, ""))
}

// Stringifies the graph. This produces a String in the graphviz dot
// language.
func (sg *SubGraph) String() string {
	return sg.StringWithAttrs(nil)
}

func (sg *SubGraph) StringWithAttrs(attrs map[int]map[string]interface{}) string {
	V := make([]string, 0, len(sg.V))
	E := make([]string, 0, len(sg.E))
	safeStr := func(i interface{}) string{
		s := fmt.Sprint(i)
		s = strings.Replace(s, "\n", "\\n", -1)
		s = strings.Replace(s, "\"", "\\\"", -1)
		return s
	}
	renderAttrs := func(v *Vertex) string {
		a := attrs[v.Id]
		label := sg.G.Colors[v.Color]
		strs := make([]string, 0, len(a)+1)
		strs = append(strs, fmt.Sprintf("label=\"%v\"", label))
		for name, value := range a {
			if name == "label" || name == "id" {
				continue
			}
			strs = append(strs, fmt.Sprintf("%v=\"%v\"", name, safeStr(value)))
		}
		return strings.Join(strs, ",")
	}
	for _, v := range sg.V {
		V = append(V, fmt.Sprintf(
			"%v [%v];",
			sg.G.V[v.Id].Id,
			renderAttrs(&v),
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


