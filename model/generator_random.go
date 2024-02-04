package model

import (
	"math"
	"math/rand"
)

// FastGNPRandomGraph generates a random undirected graph using the G(n,p) model,
// where n is the number of nodes and p is the probability of edge creation between nodes.
// The algorithm efficiently constructs the graph by avoiding unnecessary calculations
// and improving overall performance.
//
// Parameters:
//   - numberOfNodes: The total number of nodes in the graph.
//   - probabilityForEdgeCreation: The probability of creating an edge between any two nodes,
//     ranging from 0.0 (no edges) to 1.0 (fully connected graph).
//
// Returns:
//
//	An UndirectedGraph generated using the G(n,p) model, with edges connecting nodes based on
//	the specified probability.
//
// Note:
//
//	The graph's nodes are labeled from 0 to n-1, and the edge creation follows a fast
//	non-uniform random process. This implementation is efficient for large graphs with
//	sparse connections.
//
// Returns a $G_{n,p}$ random graph, also known as an Erdős-Rényi graph or a binomial graph.
// References: [1] Vladimir Batagelj and Ulrik Brandes, "Efficient generation of large random networks", Phys. Rev. E, 71, 036113, 2005.
func FastGNPRandomGraph(numberOfNodes int, probabilityForEdgeCreation float64) (g UndirectedGraph) {
	g = UndirectedGraph{}
	g.Edges = make(map[Node][]Node)
	g.Nodes = make(map[Node]bool, numberOfNodes)
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
				g.AddEdge(Edge{Node(v), Node(w)})
			}
		}
	}
	return g
}

// In the $G_{n,m}$ model, a graph is chosen uniformly at random from the set
// of all graphs with $n$ nodes and $m$ edges.
// Algorithm by Keith M. Briggs Mar 31, 2006.
// Inspired by Knuth's Algorithm S (Selection sampling technique),
// in section 3.4.2 of [1]
// References: [1] Donald E. Knuth, The Art of Computer Programming,
// Volume 2/Seminumerical algorithms, Third Edition, Addison-Wesley, 1997.
func DenseGNMRandomGraph(numberOfNodes int, numberOfEdges int) (g *UndirectedGraph) {
	edgesMax := numberOfNodes * (numberOfNodes - 1) // 2
	if numberOfEdges >= edgesMax {
		return CompleteGraph(numberOfNodes)
	} else {
		g = &UndirectedGraph{}
	}
	if numberOfNodes == 1 || numberOfEdges >= edgesMax {
		return g
	}

	u, v, t, k := 0, 0, 0, 0
	for {
		if (t + rand.Int()*(edgesMax-t)) < (numberOfEdges - k) {
			g.AddEdge(Edge{Node(u), Node(v)})
			k = k + 1
			if k == numberOfEdges {
				return g
			}
		}
		t = t + 1
		v = v + 1
		if v == numberOfNodes {
			u = u + 1
			v = u + 1
		}
	}
}
