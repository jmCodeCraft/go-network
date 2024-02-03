package model

import "log"

// CompleteGraph generates a complete undirected graph with the specified number of nodes.
//
// The CompleteGraph function creates a graph where every pair of distinct nodes is connected by an edge,
// resulting in a fully connected graph. It is commonly used in graph theory to model complete relationships.
//
// Parameters:
//   - numberOfNodes: The total number of nodes in the complete graph.
//
// Returns:
//
//	A pointer to an UndirectedGraph representing the generated complete graph.
func CompleteGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}

	// Create edges between all pairs of distinct nodes
	for i := 0; i < numberOfNodes; i++ {
		for j := i + 1; j < numberOfNodes; j++ {
			g.AddEdge(Edge{
				Node1: Node(i),
				Node2: Node(j),
			})
		}
	}

	return g
}

// LadderGraph generates a Ladder graph with 2n nodes, formed by two parallel paths of length n.
//
// The Ladder graph consists of two parallel paths, each of length n, connected by edges
// between corresponding nodes in the two paths. It is often used as a model in graph theory.
//
// Parameters:
//   - nodesInSinglePath: The length of each parallel path in the Ladder graph.
//
// Returns:
//
//	A pointer to an UndirectedGraph representing the generated Ladder graph.
func LadderGraph(nodesInSinglePath int) *UndirectedGraph {
	g := &UndirectedGraph{}
	edgesLadder1 := Pairwise(Range(nodesInSinglePath, 2*nodesInSinglePath))
	log.Println("HAEEEEY")
	log.Println(edgesLadder1)
	for i := 0; i < len(edgesLadder1); i++ {
		g.AddEdge(edgesLadder1[i])
	}
	for i := 0; i < nodesInSinglePath; i++ {
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node(i + nodesInSinglePath),
		})
	}
	return g
}

// CircularLadderGraph returns the circular ladder graph $CL_n$ of length n
func CircularLadderGraph(nodesInSinglePath int) *UndirectedGraph {
	g := LadderGraph(nodesInSinglePath)
	g.AddEdge(Edge{
		Node2: Node(nodesInSinglePath - 1),
	})
	g.AddEdge(Edge{
		Node1: Node(nodesInSinglePath),
		Node2: Node(2*nodesInSinglePath - 1),
	})
	return g
}
