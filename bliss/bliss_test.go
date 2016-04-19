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

import "testing"

import "reflect"

func TestCanonize(t *testing.T) {
	expected := []uint{5, 4, 1, 2, 3, 0}
	nodes := []uint32{1, 1, 0, 0, 0, 0}
	edges := []BlissEdge{{0, 2}, {0, 3}, {1, 4}, {1, 5}, {3, 5}, {4, 2}}
	actual := Canonize(nodes, edges)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestNew(t *testing.T) {
	Do(0, func(g *Graph) {
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

func TestPermutation(t *testing.T) {
	Do(0, func(g *Graph) {
		a := g.AddVertex(1)
		b := g.AddVertex(1)
		c := g.AddVertex(2)
		d := g.AddVertex(2)
		g.AddEdge(a, b)
		g.AddEdge(c, d)
		g.AddEdge(a, c)
		g.AddEdge(b, d)
		p := g.CanonicalPermutation()
		e := []uint{1, 0, 3, 2}
		if !reflect.DeepEqual(p, e) {
			t.Errorf("Expected %v got %v", e, p)
		}
	})
}

func TestCompare(t *testing.T) {
	Do(0, func(g1 *Graph) {
		a := g1.AddVertex(1)
		b := g1.AddVertex(1)
		c := g1.AddVertex(2)
		d := g1.AddVertex(2)
		g1.AddEdge(a, b)
		g1.AddEdge(c, d)
		g1.AddEdge(a, c)
		g1.AddEdge(b, d)
		Do(0, func(g2 *Graph) {
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
		Do(0, func(g2 *Graph) {
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
