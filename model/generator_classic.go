package model

func CompleteGraph(numberOfNodes int) (g UndirectedGraph) {
	g = UndirectedGraph{}
	for i := 0; i <= numberOfNodes; i++ {
		for j := i + 1; j <= numberOfNodes; i++ {
			g.AddEdge(Edge{
				Node1: Node{NodeId: i},
				Node2: Node{NodeId: j},
			})
		}
	}
	return g
}

// LadderGraph returns the Ladder graph of length n and 2n nodes
func LadderGraph(nodesInSinglePath int) (g UndirectedGraph) {
	g = UndirectedGraph{}
	edgesLadder1 := Pairwise(Range(nodesInSinglePath, 2*nodesInSinglePath))
	for i := 0; i <= len(edgesLadder1); i++ {
		g.AddEdge(edgesLadder1[i])
	}
	for i := 0; i <= nodesInSinglePath; i++ {
		g.AddEdge(Edge{
			Node1: Node{NodeId: i},
			Node2: Node{NodeId: i + nodesInSinglePath},
		})
	}
	return g
}

// CircularLadderGraph returns the circular ladder graph $CL_n$ of length n
func CircularLadderGraph(nodesInSinglePath int) (g UndirectedGraph) {
	g = LadderGraph(nodesInSinglePath)
	g.AddEdge(Edge{
		Node1: Node{},
		Node2: Node{NodeId: nodesInSinglePath - 1},
	})
	g.AddEdge(Edge{
		Node1: Node{NodeId: nodesInSinglePath},
		Node2: Node{NodeId: 2*nodesInSinglePath - 1},
	})
	return g
}
