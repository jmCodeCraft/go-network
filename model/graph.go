package model

import (
	"fmt"
	"log/slog"
	"strings"
)

type Graph interface {
	AddEdge(edge Edge)
	AddNode(node Node)
	GetEdgeTuples() []Edge
	Sample(sampler ISamplingStrategy, samplingRate float32) UndirectedGraph
	NodeDegree(node Node) int
	NumberOfEdges() int
	HasNode(node Node) bool
	RemoveEdge(edge Edge)
	RemoveNode(node Node)
}

type Node int

type Edge struct {
	Node1 Node
	Node2 Node
}

type Components struct {
	ComponentsArray     [][]Node
	ComponentsDict      []map[Node]bool
	BiggestComponentIdx int
}

type UndirectedGraph struct {
	Nodes map[Node]bool
	Edges map[Node][]Node
}

func (c *Components) AddComponent(component []Node) {
	c.ComponentsArray = append(c.ComponentsArray, component)
	componentDict := map[Node]bool{}
	for i := 0; i < len(component); i++ {
		componentDict[component[i]] = true
	}
	c.ComponentsDict = append(c.ComponentsDict, componentDict)

	if c.BiggestComponentIdx > -1 && len(c.ComponentsArray[c.BiggestComponentIdx]) < len(component) {
		c.BiggestComponentIdx = len(c.ComponentsArray) - 1
	}
}

// String returns a string representation of the UndirectedGraph.
func (g UndirectedGraph) String() string {
	var result strings.Builder

	result.WriteString("Nodes: {")
	for node := range g.Nodes {
		result.WriteString(fmt.Sprintf("%v, ", node))
	}
	if len(g.Nodes) > 0 {
		result.WriteString(result.String()[:result.Len()-2]) // Remove trailing comma and space
	}
	result.WriteString("}\n")

	result.WriteString("Edges: {\n")
	for node, neighbors := range g.Edges {
		result.WriteString(fmt.Sprintf("  %v: {%v}\n", node, neighborsString(neighbors)))
	}
	result.WriteString("}")

	return result.String()
}

// neighborsString returns a formatted string representation of a slice of neighbors.
func neighborsString(neighbors []Node) string {
	var result strings.Builder

	result.WriteString("{")
	for _, neighbor := range neighbors {
		result.WriteString(fmt.Sprintf("%v, ", neighbor))
	}
	if len(neighbors) > 0 {
		result.WriteString(result.String()[:result.Len()-2]) // Remove trailing comma and space
	}
	result.WriteString("}")

	return result.String()
}

/*
AddEdge adds an undirected edge to the UndirectedGraph.

Parameters:
- edge: An Edge struct representing the edge to be added, with Node1 and Node2 as the connected nodes.

Description:
The function ensures the existence of the Edges map in the UndirectedGraph and adds the specified edge. It also adds both connected nodes to the graph if they do not already exist.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	edge := Edge{Node1: 1, Node2: 2}
	undirectedGraph.AddEdge(edge)

	fmt.Println(undirectedGraph.Edges) // Output: map[1:[2] 2:[1]]
*/
func (g *UndirectedGraph) AddEdge(edge Edge) {
	// Ensure the existence of the Edges map
	if g.Edges == nil {
		g.Edges = make(map[Node][]Node)
	}

	// Add both connected nodes to the graph
	g.AddNode(edge.Node1)
	g.AddNode(edge.Node2)

	// Add the edge to the Edges map
	g.Edges[edge.Node1] = append(g.Edges[edge.Node1], edge.Node2)
	g.Edges[edge.Node2] = append(g.Edges[edge.Node2], edge.Node1)
}

/*
AddNode adds a node to the UndirectedGraph.

Parameters:
- node: A Node representing the node to be added to the graph.

Description:
The function ensures the existence of the Nodes map in the UndirectedGraph and adds the specified node if it does not already exist.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	undirectedGraph.AddNode(1)
	undirectedGraph.AddNode(2)

	fmt.Println(undirectedGraph.Nodes) // Output: map[1:true 2:true]
*/
func (g *UndirectedGraph) AddNode(node Node) {
	// Ensure the existence of the Nodes map
	if g.Nodes == nil {
		g.Nodes = make(map[Node]bool)
	}

	// Add the node to the Nodes map
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

/*
AddNodes adds multiple nodes to the UndirectedGraph.

Parameters:
- nodes: A slice of Node representing the nodes to be added to the graph.

Description:
The function iterates over the provided slice of nodes and adds each node to the graph using the AddNode method.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	nodeSlice := []Node{1, 2, 3}
	undirectedGraph.AddNodes(nodeSlice)

	fmt.Println(undirectedGraph.Nodes) // Output: map[1:true 2:true 3:true]
*/
func (g *UndirectedGraph) AddNodes(nodes []Node) {
	// Iterate over the provided slice of nodes and add each node to the graph
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

/*
GetEdgeTuples returns a slice of Edge representing all the edges in the UndirectedGraph.

Description:
The function iterates through the graph's Edges map, creating an Edge tuple for each pair of connected nodes. The resulting slice provides a comprehensive list of all edges in the undirected graph.

Returns:
- edges: A slice of Edge, where each Edge is a tuple representing a connection between two nodes in the graph.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: map[Node][]Node{
			1: {2, 3},
			2: {1, 3},
			3: {1, 2},
		},
	}

	edgeTuples := undirectedGraph.GetEdgeTuples()
	fmt.Println(edgeTuples) // Output: [{1 2} {1 3} {2 1} {2 3} {3 1} {3 2}]
*/
func (g *UndirectedGraph) GetEdgeTuples() []Edge {
	var edges []Edge
	for node1, array := range g.Edges {
		for _, node2 := range array {
			edges = append(edges, Edge{node1, node2})
		}
	}
	return edges
}

func (g *UndirectedGraph) Sample(sampler ISamplingStrategy, ratioNodesToDelete float32) UndirectedGraph {
	return sampler.Sample(g, ratioNodesToDelete)
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

/*
HasNode checks if the UndirectedGraph contains a specific node.

Parameters:
- node: A Node representing the node to check for existence in the graph.

Returns:
- bool: True if the node exists in the graph, otherwise false.

Description:
The function returns true if the specified node is present in the Nodes map of the UndirectedGraph, indicating its existence in the graph. Otherwise, it returns false.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: map[Node][]Node{
			1: {2, 3},
			2: {1, 3},
			3: {1, 2},
		},
	}

	result1 := undirectedGraph.HasNode(2) // true
	result2 := undirectedGraph.HasNode(4) // false
*/
func (g *UndirectedGraph) HasNode(node Node) bool {
	return g.Nodes[node]
}

/*
RemoveEdge removes an undirected edge from the UndirectedGraph.

Parameters:
- edge: An Edge struct representing the edge to be removed, with Node1 and Node2 as the connected nodes.

Description:
The function removes the specified edge from the Edges map of the UndirectedGraph, effectively disconnecting the two nodes associated with the edge.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: map[Node][]Node{
			1: {2, 3},
			2: {1, 3},
			3: {1, 2},
		},
	}

	edgeToRemove := Edge{Node1: 1, Node2: 2}
	undirectedGraph.RemoveEdge(edgeToRemove)

	// After removal, Edges map becomes: map[1:[3] 2:[3] 3:[1 2]]
*/
func (g *UndirectedGraph) RemoveEdge(edge Edge) {
	if len(g.Edges[edge.Node1]) > 0 {
		g.Edges[edge.Node1] = DeleteFromSlice(g.Edges[edge.Node1], edge.Node2)
	}

	if len(g.Edges[edge.Node2]) > 0 {
		g.Edges[edge.Node2] = DeleteFromSlice(g.Edges[edge.Node2], edge.Node1)
	}
}

/*
RemoveNode removes a node from the UndirectedGraph and all associated edges.

Parameters:
- node: A Node representing the node to be removed from the graph.

Description:
The function removes the specified node from the Nodes map of the UndirectedGraph, as well as all edges connected to the node in the Edges map. This operation effectively disconnects the node and eliminates all edges involving that node.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: map[Node][]Node{
			1: {2, 3},
			2: {1, 3},
			3: {1, 2},
		},
	}

	nodeToRemove := Node(2)
	undirectedGraph.RemoveNode(nodeToRemove)

	// After removal, Nodes map becomes: map[1:true 3:true]
	// After removal, Edges map becomes: map[1:[3] 3:[1]]
*/
func (g *UndirectedGraph) RemoveNode(node Node) {
	// Remove the node from the Nodes map
	delete(g.Nodes, node)

	// Remove the node from the Edges map and update neighbors' adjacency lists
	for neighbor, edges := range g.Edges {
		g.Edges[neighbor] = DeleteFromSlice(edges, node)
	}

	// Delete the entry for the removed node from the Edges map
	delete(g.Edges, node)
}

func ConnectedComponents(g UndirectedGraph) (components Components) {
	visited := map[Node]bool{}
	components = Components{
		ComponentsDict:      make([]map[Node]bool, 0),
		ComponentsArray:     make([][]Node, 0),
		BiggestComponentIdx: 0,
	}

	for node := range g.Nodes {
		if !visited[node] {
			component := []Node{}
			stack := []Node{node}
			for len(stack) > 0 {
				// todo remove this
				slog.Info(fmt.Sprintf("%v", len(stack)))
				current := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				component = append(component, current)
				visited[current] = true
				stack = append(stack, g.Edges[current]...)
			}
			components.AddComponent(component)
		}
	}
	return components
}
