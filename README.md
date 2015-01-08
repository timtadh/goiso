# goiso - Graph Isomorphism Testing

A wrapper around [bliss](http://www.tcs.hut.fi/Software/bliss/) for graph
isomorphism testing and canonical labeling. Bliss is the work Tommi Junttila and
Petteri Kaski. You should cite their papers:

- Tommi Junttila and Petteri Kaski. Engineering an efficient canonical labeling
  tool for large and sparse graphs. In: Proceedings of the Ninth Workshop on
  Algorithm Engineering and Experiments (ALENEX07), pages 135-149, SIAM, 2007.

- Tommi Junttila and Petteri Kaski. Conflict Propagation and Component Recursion
  for Canonical Labeling. In: Proceedings of the 1st International ICST
  Conference on Theory and Practice of Algorithms (TAPAS 2011), Springer, 2011

I have made a few modifications to their library mainly in the C interface side
so I can interface more easily with Go.

### Licensing

All work is licensed under the GPL Version 3. Please respect the license.

Copyright Holders:

      Copyright (c) 2006-2011 Tommi Junttila
      Copyright (c) 2014 Tim Henderson

### Usage

So far there is just a low level interface on top the C interface to bliss.
Eventually, I am hoping to build a more goish interface in the goiso package
that lets you extract the canonical labeling. This interface will likely remain
the most efficient.

You can read the docs over at godoc:
[documentation](http://godoc.org/github.com/timtadh/goiso/bliss)

import

    import "github.com/timtadh/goiso/bliss"

api example

    bliss.Graph(0, func(g1 *BlissGraph) {
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
    })

Or with manual memory management:

    g1 := bliss.NewGraph(0)
    defer g1.Release()
    a := g1.AddVertex(1)
    b := g1.AddVertex(1)
    c := g1.AddVertex(2)
    d := g1.AddVertex(2)
    g1.AddEdge(a, b)
    g1.AddEdge(c, d)
    g1.AddEdge(a, c)
    g1.AddEdge(b, d)

    // do stuff with g1

