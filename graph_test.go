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
	vord, eord := g.CanonicalPermutation()
	t.Log(vord)
	t.Log(eord)
	t.Fatal("fail")
}

