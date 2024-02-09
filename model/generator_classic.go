package model

func CompleteGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
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

// LadderGraph returns the Ladder graph of length n and 2n nodes
func LadderGraph(nodesInSinglePath int) *UndirectedGraph {
	g := &UndirectedGraph{}

	// Generate and add edges for the ladder structure
	for _, edge := range Pairwise(Range(nodesInSinglePath, 2*nodesInSinglePath)) {
		g.AddEdge(edge)
	}

	// Add rung edges between the two paths of the ladder
	for i := 0; i < nodesInSinglePath; i++ { // nodesInSinglePath = 3
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node(i + nodesInSinglePath),
		})

		if i != nodesInSinglePath-1 { // i != 2
			g.AddEdge(Edge{
				Node1: Node(i),
				Node2: Node(i + 1),
			})
		}
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
