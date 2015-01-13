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

func TestHello(t *testing.T) {
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

