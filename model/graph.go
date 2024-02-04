package model

import (
	"fmt"
)

type Graph interface {
	AddEdge(edge Edge)
	AddNode(node Node)
	GetEdgeTuples() []Edge
	Sample(sampler SamplingStrategy, samplingRate float32) UndirectedGraph
	NodeDegree(node Node) int
	NumberOfEdges() int
}

type Node int

type Edge struct {
	Node1 Node
	Node2 Node
}

type UndirectedGraph struct {
	Nodes map[Node]bool
	Edges map[Node][]Node
}

func (g UndirectedGraph) String() string {
	return fmt.Sprintf("graph has %d nodes and %d edges. Nodes: %v; Edges: %v", len(g.Nodes), len(g.Edges), g.Nodes, g.Edges)
}

func (g *UndirectedGraph) AddEdge(edge Edge) {
	if g.Edges == nil {
		g.Edges = make(map[Node][]Node)
	}
	g.AddNode(edge.Node1)
	g.AddNode(edge.Node2)
	g.Edges[edge.Node1] = append(g.Edges[edge.Node1], edge.Node2)
	g.Edges[edge.Node2] = append(g.Edges[edge.Node2], edge.Node1)
}

func (g *UndirectedGraph) AddNode(node Node) {
	if g.Nodes == nil {
		g.Nodes = make(map[Node]bool)
	}
	g.Nodes[node] = true
}

func (g *UndirectedGraph) AddEdgesFromIntTupleList(edges [][2]int) {
	for _, nodes := range edges {
		g.AddEdge(Edge{Node(nodes[0]), Node(nodes[1])})
	}
}

func (g *UndirectedGraph) AddEdgesFromIntEdgeList(sourceNode Node, edges []Node) {
	for _, node := range edges {
		g.AddEdge(Edge{sourceNode, node})
	}
}

func (g *UndirectedGraph) AddNodes(nodes []Node) {
	for _, node := range nodes {
		g.AddNode(node)
	}
}

// NodeDegree returns the degree (number of incident edges) of the specified node in the graph.
func (g *UndirectedGraph) NodeDegree(node Node) int {
	// Check if the node exists in the graph
	if _, exists := g.Nodes[node]; !exists {
		return 0 // Node not found, degree is 0
	}

	// Retrieve the neighbors of the node and return the degree
	return len(g.Edges[node])
}

// todo suggest rename to GetEdges
func (g *UndirectedGraph) GetEdgeTuples() []Edge {
	var edges []Edge
	for node1, array := range g.Edges {
		for _, node2 := range array {
			edges = append(edges, Edge{node1, node2})
		}
	}
	return edges
}

func (g *UndirectedGraph) Sample(sampler SamplingStrategy, samplingRate float32) UndirectedGraph {
	return sampler.Sample(g, samplingRate)
}

// NumberOfEdges returns the total number of edges in the undirected graph.
func (g *UndirectedGraph) NumberOfEdges() int {
	totalEdges := 0

	// Iterate over each node's neighbors and count the unique edges
	for _, neighbors := range g.Edges {
		totalEdges += len(neighbors)
	}

	// Divide by 2 to account for the fact that each edge is counted twice (undirected graph)
	return totalEdges / 2
}
