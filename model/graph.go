package model

import (
	"github.com/jmCodeCraft/go-network/algorithm"
)

type Graph struct {
	Edges map[int][]int
	Nodes map[int]bool
}

func (g *Graph) AddEdge(node1, node2 int) {
	g.Edges[node1] = append(g.Edges[node1], node2)
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
	for node := range nodes {
		g.AddNode(node)
	}
}

func (g *Graph) Sample(sampler *algorithm.SamplingStrategy) (g Graph) {
	return sampler.Sample(g)
}
