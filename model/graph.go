package model

import "fmt"

type Graph interface {
	AddEdge(edge Edge)
	AddNode(node Node)
	GetEdgeTuples() []Edge
	Sample(sampler *SamplingStrategy, samplingRate float32) Graph
}

type UndirectedGraph struct {
	Nodes map[int]bool
	Edges map[int][]int
}

func (g UndirectedGraph) String() string {
	return fmt.Sprintf("graph has %d nodes and %d edges", len(g.Nodes), len(g.Edges))
}

type Edge struct {
	Node1 Node
	Node2 Node
}

type Node struct {
	NodeId int
}

func (g *UndirectedGraph) AddEdge(edge Edge) {
	g.AddNode(edge.Node1)
	g.AddNode(edge.Node2)
	g.Edges[edge.Node1.NodeId] = append(g.Edges[edge.Node1.NodeId], edge.Node2.NodeId)
}

func (g *UndirectedGraph) AddNode(node Node) {
	g.Nodes[node.NodeId] = true
}

func (g *UndirectedGraph) AddEdgesFromIntTupleList(edges [][2]int) {
	for _, nodes := range edges {
		g.AddEdge(Edge{Node{nodes[0]}, Node{nodes[1]}})
	}
}

func (g *UndirectedGraph) AddEdgesFromIntEdgeList(sourceNode int, edges []int) {
	for _, node := range edges {
		g.AddEdge(Edge{Node{sourceNode}, Node{node}})
	}
}

func (g *UndirectedGraph) AddNodes(nodes map[int]bool) {
	for node := range nodes {
		g.AddNode(Node{node})
	}
}

// todo suggest rename to GetEdges
func (g *UndirectedGraph) GetEdgeTuples() []Edge {
	var edges []Edge
	for node1, array := range g.Edges {
		for node2 := range array {
			edges = append(edges, Edge{Node{node1}, Node{node2}})
		}
	}
	return edges
}

func (g *UndirectedGraph) Sample(sampler SamplingStrategy, samplingRate float32) UndirectedGraph {
	return sampler.Sample(g, samplingRate)
}
