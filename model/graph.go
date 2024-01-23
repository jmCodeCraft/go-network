package model

type Graph struct {
	Edges [][2]int
	Nodes map[int]bool
}

func (g *Graph) AddEdge(node1, node2 int) {
	g.Edges = append(g.Edges, [2]int{node1, node2})
}

func (g *Graph) AddNode(node int) {
	g.Nodes[node] = true
}

func (g *Graph) AddEdgesFromTupleList(edges [][2]int) {
	for _, nodes := range edges {
		g.AddEdge(nodes[0], nodes[1])
	}
}

func (g *Graph) AddEdgesFromEdgeList(sourceNode int, edges []int) {
	for _, node := range edges {
		g.AddEdge(sourceNode, node)
	}
}

func (g *Graph) AddNodes(nodes map[int]bool) {
	for node, _ := range nodes {
		g.AddNode(node)
	}
}
