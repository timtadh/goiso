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

import "testing"

func TestCanon(t *testing.T) {
	g := NewGraph(4, 4)
	a := g.AddVertex(12, "blue")
	b := g.AddVertex(7, "blue")
	c := g.AddVertex(57, "green")
	d := g.AddVertex(9, "green")
	g.AddEdge(a, b, "purple")
	g.AddEdge(c, d, "purple")
	g.AddEdge(a, c, "purple")
	g.AddEdge(b, d, "purple")

	g2 := NewGraph(4, 4)
	w := g2.AddVertex(12, "blue")
	x := g2.AddVertex(7, "blue")
	y := g2.AddVertex(57, "green")
	z := g2.AddVertex(9, "green")
	g2.AddEdge(x, w, "purple")
	g2.AddEdge(z, y, "purple")
	g2.AddEdge(w, y, "purple")
	g2.AddEdge(x, z, "purple")

	can := g.Canonical()
	can2 := g2.Canonical()
	e := []int{7, 12, 9, 57}
	for i, vid := range e {
		if can.V[i].Id != vid {
			t.Error("canonical was incorrect")
		}
	}
	if can.Label() != can2.Label() {
		t.Error("graphs should have been equal")
	}
}

func TestSubgraph(t *testing.T) {
	g := NewGraph(8, 8)
	{
		a := g.AddVertex(12, "blue")
		b := g.AddVertex(7, "blue")
		c := g.AddVertex(57, "green")
		d := g.AddVertex(9, "green")
		g.AddEdge(a, b, "purple")
		g.AddEdge(c, d, "purple")
		g.AddEdge(a, c, "purple")
		g.AddEdge(b, d, "purple")

		w := g.AddVertex(12, "blue")
		x := g.AddVertex(7, "blue")
		y := g.AddVertex(57, "green")
		z := g.AddVertex(9, "green")
		g.AddEdge(x, w, "purple")
		g.AddEdge(z, y, "purple")
		g.AddEdge(w, y, "purple")
		g.AddEdge(x, z, "purple")
	}

	g2 := NewGraph(4, 4)
	{
		w := g2.AddVertex(12, "blue")
		x := g2.AddVertex(7, "blue")
		y := g2.AddVertex(57, "green")
		z := g2.AddVertex(9, "green")
		g2.AddEdge(x, w, "purple")
		g2.AddEdge(z, y, "purple")
		g2.AddEdge(w, y, "purple")
		g2.AddEdge(x, z, "purple")
	}
	can2 := g2.Canonical()

	g3 := NewGraph(4, 4)
	{
		a := g3.AddVertex(12, "blue")
		b := g3.AddVertex(7, "blue")
		c := g3.AddVertex(57, "green")
		d := g3.AddVertex(9, "green")
		g3.AddEdge(a, b, "purple")
		g3.AddEdge(c, d, "purple")
		g3.AddEdge(a, c, "purple")
		g3.AddEdge(b, d, "purple")
	}
	can3 := g2.Canonical()

	t.Log(can2.Label())
	t.Log(can3.Label())

	sg1 := g.SubGraph([]int{0, 1, 2, 3}, nil)
	sg2 := g.SubGraph([]int{4, 5, 6, 7}, nil)
	t.Log(sg1.Label())
	t.Log(sg2.Label())
	if can2.Label() != can3.Label() {
		t.Error("g2 != g3")
	}
	if can2.Label() != sg1.Label() {
		t.Error("g2 != sg1")
	}
	if sg1.Label() != sg2.Label() {
		t.Error("sg1 != sg2")
	}
}

