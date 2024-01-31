package generators

import (
	"github.com/jmCodeCraft/go-network/model"
)

func EmptyGraph() (g model.UndirectedGraph) {
	// TODO: review - how do we initialize, to avoid nil values?
	return model.UndirectedGraph{
		Edges: nil,
		Nodes: nil,
	}
}

func CompleteGraph(numberOfNodes int) (g model.UndirectedGraph) {
	g = EmptyGraph()
	for i := 0; i <= numberOfNodes; i++ {
		for j := i + 1; j <= numberOfNodes; i++ {
			g.AddEdge(model.Edge{model.Node{i}, model.Node{j}})
		}
	}
	return g
}

// Returns the Ladder graph of length n and 2n nodes
func LadderGraph(nodesInSinglePath int) (g model.UndirectedGraph) {
	g = EmptyGraph()
	edgesLadder1 := model.Pairwise(model.Range(nodesInSinglePath, 2*nodesInSinglePath))
	for i := 0; i <= len(edgesLadder1); i++ {
		g.AddEdge(edgesLadder1[i])
	}
	for i := 0; i <= nodesInSinglePath; i++ {
		g.AddEdge(model.Edge{model.Node{i}, model.Node{i + nodesInSinglePath}})
	}
	return g
}

// Returns the circular ladder graph $CL_n$ of length n
func CircularLadderGraph(nodesInSinglePath int) (g model.UndirectedGraph) {
	g := LadderGraph(nodesInSinglePath)
	g.AddEdge(model.Edge{model.Node{0}, model.Node{nodesInSinglePath - 1}})
	g.AddEdge(model.Edge{model.Node{nodesInSinglePath}, model.Node{2*nodesInSinglePath - 1}})
	return g
}
