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
	"sort"
	"strings"
)

import (
	"github.com/timtadh/goiso/bliss"
)

type Graph struct {
	V []Vertex
	E []Edge
	kids [][]*Edge
	Colors []string
	colorSet map[string]int
	closed bool
	blissMap *blissMap
}

type blissMap struct {
	V []blissVertex
	E []Arc
}

type blissVertex struct {
	edge bool
	idx int
	color int
}

type Vertex struct {
	Idx int
	Id int
	Color int
}

func (v *Vertex) Copy(idx int) Vertex {
	return Vertex{
		Idx: idx,
		Id: v.Id,
		Color: v.Color,
	}
}

type Arc struct {
	Src, Targ int
}

type Edge struct {
	Arc
	Idx int
	Color int
}

func (e *Edge) Copy(idx, src, targ int) Edge {
	return Edge{
		Arc: Arc{
			Src: src,
			Targ: targ,
		},
		Idx: idx,
		Color: e.Color,
	}
}

type perm struct {
	idx, p int
}

type perms []perm

func (self perms) Len() int           { return len(self) }
func (self perms) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self perms) Less(i, j int) bool { return self[i].p < self[j].p }

func NewGraph(V, E int) Graph {
	return Graph{
		V: make([]Vertex, 0, V),
		E: make([]Edge, 0, E),
		kids: make([][]*Edge, 0, V),
		Colors: make([]string, 0, V),
		colorSet: make(map[string]int),
	}
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
	return fmt.Sprintf("%v%v", strings.Join(V, ""), strings.Join(E, ""))
}

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

func (g *Graph) Finalize() {
	g.closed = true
	V := make([]blissVertex, 0, len(g.V) + len(g.E))
	E := make([]Arc, 0, len(g.E)*2)
	for _, v := range g.V {
		V = append(V, blissVertex{
			edge: false,
			idx: v.Idx,
			color: v.Color,
		})
	}
	for _, e := range g.E {
		eid := len(V)
		V = append(V, blissVertex{
			edge: true,
			idx: e.Idx,
			color: e.Color,
		})
		E = append(E, Arc{
			Src: e.Src,
			Targ: eid,
		})
		E = append(E, Arc{
			Src: eid,
			Targ: e.Targ,
		})
	}
	g.blissMap = &blissMap{V, E}
}

func (g *Graph) Canonical() Graph {
	ng := Graph{
		V: make([]Vertex, len(g.V)),
		E: make([]Edge, len(g.E)),
		kids: make([][]*Edge, len(g.kids)),
		Colors: make([]string, len(g.Colors)),
		colorSet: make(map[string]int),
	}
	copy(ng.Colors, g.Colors)
	for cid, color := range ng.Colors {
		ng.colorSet[color] = cid
	}
	for i := range ng.kids {
		ng.kids[i] = make([]*Edge, 0, 5)
	}
	vord, eord := g.CanonicalPermutation()
	// i is the old vid, j is the new vid
	for i, j := range vord {
		ng.V[j] = g.V[i].Copy(j)
	}
	for i, j := range eord {
		ng.E[j] = g.E[i].Copy(j, vord[g.E[i].Src], vord[g.E[i].Targ])
		ng.kids[vord[g.E[i].Src]] = append(ng.kids[vord[g.E[i].Src]], &ng.E[j])
	}
	return ng
}

func (g *Graph) CanonicalPermutation() (Vord, Eord []int) {
	if !g.closed {
		g.Finalize()
	}
	bg := g.blissGraph()
	defer bg.Release()
	P := bg.CanonicalPermutation()
	VP := make(perms, 0, len(g.V))
	EP := make(perms, 0, len(g.E))
	for i, p := range P {
		v := g.blissMap.V[i]
		if v.edge {
			EP = append(EP, perm{v.idx, int(p)})
		} else {
			VP = append(VP, perm{v.idx, int(p)})
		}
	}
	sort.Sort(VP)
	sort.Sort(EP)
	Vord = make([]int, len(g.V))
	Eord = make([]int, len(g.E))
	for p, vp := range VP {
		Vord[vp.idx] = p
	}
	for p, ep := range EP {
		Eord[ep.idx] = p
	}
	return Vord, Eord
}

func (g *Graph) AddVertex(id int, label string) *Vertex {
	if g.closed {
		return nil
	}
	v := Vertex{
		Idx: len(g.V),
		Id: id,
		Color: g.color(label),
	}
	g.V = append(g.V, v)
	g.kids = append(g.kids, make([]*Edge, 0, 5))
	return &v
}

func (g *Graph) AddEdge(u, v *Vertex, label string) *Edge {
	if g.closed {
		return nil
	}
	e := Edge{
		Arc: Arc{
			Src: u.Idx,
			Targ: v.Idx,
		},
		Idx: len(g.E),
		Color: g.color(label),
	}
	g.E = append(g.E, e)
	g.kids[e.Arc.Src] = append(g.kids[e.Arc.Src], &e)
	return &e
}

func (g *Graph) color(label string) int {
	if cid, has := g.colorSet[label]; has {
		return cid
	}
	cid := len(g.Colors)
	g.colorSet[label] = cid
	g.Colors = append(g.Colors, label)
	return cid
}

func (g *Graph) blissGraph() *bliss.BlissGraph {
	bg := bliss.NewGraph(0)
	for _, v := range g.blissMap.V {
		bg.AddVertex(uint(v.color))
	}
	for _, e := range g.blissMap.E {
		bg.AddEdge(uint(e.Src), uint(e.Targ))
	}
	return bg
}

