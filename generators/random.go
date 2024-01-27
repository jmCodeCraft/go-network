package generators

import (
	"math"
	"math/rand"

	"github.com/jmCodeCraft/go-network/model"
)

// Returns a $G_{n,p}$ random graph, also known as an Erdős-Rényi graph or a binomial graph.
// References: [1] Vladimir Batagelj and Ulrik Brandes, "Efficient generation of large random networks", Phys. Rev. E, 71, 036113, 2005.
func FastGNPRandomGraph(numberOfNodes int, probabilityForEdgeCreation float64) (g model.Graph) {
	g = model.Graph{
		Edges: nil,
		Nodes: nil,
	}
	lp := math.Log(1.0 - probabilityForEdgeCreation)
	// Nodes in graph are from 0,n-1 (start with v as the second node index).
	v := 1
	w := -1
	for v < numberOfNodes {
		lr := math.Log(1.0 - rand.Float64())
		w = w + 1 + int(lr/lp)
		for w >= v && v < numberOfNodes {
			w = w - v
			v = v + 1
			if v < numberOfNodes {
				g.AddEdge(v, w)
			}
		}
	}
	return g
}
