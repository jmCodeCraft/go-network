package generators

import (
	"github.com/jmCodeCraft/go-network/model"
)

func EmptyGraph() (g model.Graph) {
	return model.Graph{
		Edges: nil,
		Nodes: nil,
	}
}

func CompleteGraph(numberOfNodes int) (g model.Graph) {
	g = EmptyGraph()
	for i := 0; i <= numberOfNodes; i++ {
		for j := i + 1; j <= numberOfNodes; i++ {
			g.AddEdge(i, j)
		}
	}
	return g
}
