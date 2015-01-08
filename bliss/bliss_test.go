package bliss
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
  along with Bliss.  If not, see <http://www.gnu.org/licenses/>.
*/

import "testing"

func TestNew(t *testing.T) {
	Graph(0, func(g *BlissGraph) {
		a := g.AddVertex(1)
		b := g.AddVertex(1)
		c := g.AddVertex(2)
		d := g.AddVertex(2)
		g.AddEdge(a, b)
		g.AddEdge(c, d)
		g.AddEdge(a, c)
		g.AddEdge(b, d)
	})
}

func TestCompare(t *testing.T) {
	Graph(0, func(g1 *BlissGraph) {
		a := g1.AddVertex(1)
		b := g1.AddVertex(1)
		c := g1.AddVertex(2)
		d := g1.AddVertex(2)
		g1.AddEdge(a, b)
		g1.AddEdge(c, d)
		g1.AddEdge(a, c)
		g1.AddEdge(b, d)
		Graph(0, func(g2 *BlissGraph) {
			a := g2.AddVertex(1)
			b := g2.AddVertex(1)
			c := g2.AddVertex(2)
			d := g2.AddVertex(2)
			g2.AddEdge(b, a)
			g2.AddEdge(d, c)
			g2.AddEdge(a, c)
			g2.AddEdge(b, d)
			if g1.Cmp(g2) != -1 {
				t.Error("unexpected compare result")
			}
			if !g1.Iso(g2) {
				t.Error("should have been isomorphic")
			}
		})
		Graph(0, func(g2 *BlissGraph) {
			a := g2.AddVertex(1)
			b := g2.AddVertex(1)
			c := g2.AddVertex(2)
			d := g2.AddVertex(2)
			g2.AddEdge(b, a)
			g2.AddEdge(c, d)
			g2.AddEdge(a, c)
			g2.AddEdge(b, d)
			if g1.Cmp(g2) != -1 {
				t.Error("unexpected compare result")
			}
			if g1.Iso(g2) {
				t.Error("should not have been isomorphic")
			}
		})
	})
}


