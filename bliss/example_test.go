package bliss_test

import (
	"fmt"
)

import (
	"github.com/timtadh/goiso/bliss"
)

type Vertex struct {
	Color int
}

type Edge struct {
	Src, Targ, Color int
}

type Vertices []Vertex
type Edges []Edge

func (V Vertices) Iterate() (vi bliss.VertexIterator) {
	i := 0
	vi = func() (color int, _ bliss.VertexIterator) {
		if i >= len(V) {
			return 0, nil
		}
		color = V[i].Color
		i++
		return color, vi
	}
	return vi
}

func (E Edges) Iterate() (ei bliss.EdgeIterator) {
	i := 0
	ei = func() (src, targ, color int, _ bliss.EdgeIterator) {
		if i >= len(E) {
			return 0, 0, 0, nil
		}
		src = E[i].Src
		targ = E[i].Targ
		color = E[i].Color
		i++
		return src, targ, color, ei
	}
	return ei
}

// A minimal example for implementing the necessary interfaces and iterators.
func ExampleMap() {
	nodes := Vertices{{1}, {1}, {0}, {0}, {0}, {0}}
	edges := Edges{{0, 2, 2}, {0, 3, 2}, {1, 4, 2}, {1, 5, 2}, {3, 5, 2}, {4, 2, 2}}
	m := bliss.NewMap(len(nodes), len(edges), nodes.Iterate(), edges.Iterate())
	vord, eord, _ := m.CanonicalPermutation()
	fmt.Println(vord, eord)
	// Output: [4 5 1 2 3 0] [3 0 1 2 4 5]
}

