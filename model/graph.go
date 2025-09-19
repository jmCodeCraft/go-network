package model

import (
	"fmt"
	"strings"
)

type Graph interface {
	AddEdge(edge Edge)
	AddNode(node Node)
	GetEdgeTuples() []Edge
	Sample(sampler ISamplingStrategy, ratioNodesToDelete float32) (*UndirectedGraph, error)
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

type UndirectedGraph struct {
	Nodes map[Node]bool
	Edges map[Node][]Node
}

type Components struct {
	ComponentsArray     []*UndirectedGraph
	visitedNodes        map[Node]bool
	BiggestComponentIdx int
}

func (g *UndirectedGraph) Equals(other *UndirectedGraph) bool {
	if len(g.Nodes) != len(other.Nodes) {
		return false
	}

	for node := range g.Nodes {
		if !other.Nodes[node] {
			return false
		}
	}

	for node, edges := range g.Edges {
		otherEdges, ok := other.Edges[node]
		if !ok || len(edges) != len(otherEdges) {
			return false
		}

		for _, edge := range edges {
			found := false
			for _, otherEdge := range otherEdges {
				if edge == otherEdge {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

func (c *Components) AddComponent(component *UndirectedGraph) {
	c.ComponentsArray = append(c.ComponentsArray, component)

	if c.BiggestComponentIdx < 0 {
		c.BiggestComponentIdx = 0
	} else if len(c.ComponentsArray[c.BiggestComponentIdx].Nodes) < len(component.Nodes) {
		c.BiggestComponentIdx = len(c.ComponentsArray) - 1
	}
	for visitedNode := range component.Nodes {
		c.visitedNodes[visitedNode] = true
	}
}

func (c *Components) GetBiggestComponent() *UndirectedGraph {
	if c.BiggestComponentIdx >= 0 {
		return c.ComponentsArray[c.BiggestComponentIdx]
	}
	return nil
}

func (g *UndirectedGraph) String() string {
	var str strings.Builder

	str.WriteString("Nodes:\n")
	for node := range g.Nodes {
		str.WriteString(fmt.Sprintf("%d: true\t", node))
	}

	str.WriteString("\nEdges:\n")
	for node, edges := range g.Edges {
		str.WriteString(fmt.Sprintf("%d: %v\n", node, edges))
	}

	return str.String()
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

	// Only add if edge doesn’t already exist (undirected)
	if !g.HasEdge(edge.Node1, edge.Node2) {
		g.Edges[edge.Node1] = append(g.Edges[edge.Node1], edge.Node2)
		g.Edges[edge.Node2] = append(g.Edges[edge.Node2], edge.Node1)
	}
}

func (g *UndirectedGraph) DFS(startNode Node) *UndirectedGraph {
	if !g.Nodes[startNode] {
		return &UndirectedGraph{}
	}
	visited := make(map[Node]bool)
	visitedGraph := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	g.dfsUtil(startNode, visited, &visitedGraph)
	return &visitedGraph
}

func (g *UndirectedGraph) dfsUtil(node Node, visited map[Node]bool, visitedGraph *UndirectedGraph) {
	visited[node] = true
	visitedGraph.AddNode(node)

	for _, neighbor := range g.Edges[node] {
		if !visited[neighbor] {
			visitedGraph.AddEdge(Edge{Node1: node, Node2: neighbor})
			g.dfsUtil(neighbor, visited, visitedGraph)
		}
	}
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

func (g *UndirectedGraph) Sample(sampler ISamplingStrategy, ratioNodesToDelete float32) (*UndirectedGraph, error) {
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
HasEode checks if the UndirectedGraph contains a specific edge.

Parameters:
- edge: .

Returns:
- bool: True if the edge exists in the graph, otherwise false.

Description:
The function returns true if the specified edge is present in the UndirectedGraph, indicating its existence in the graph. Otherwise, it returns false.

Example:

	undirectedGraph := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
			4: true,
		},
		Edges: map[Node][]Node{
			1: {2, 3, 4},
			2: {1, 3},
			3: {1, 2, 4},
			4: {1, 3},
		},
	}

	result1 := undirectedGraph.HasEdge(2,3) // true
	result2 := undirectedGraph.HasNode(4,2) // false
*/
func (g *UndirectedGraph) HasEdge(u, v Node) bool {
	for _, nb := range g.Edges[u] {
		if nb == v {
			return true
		}
	}
	return false
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

func (g *UndirectedGraph) ContractNode(node Node) {
	neighbors := g.Edges[node]
	for i := 0; i < len(neighbors); i++ {
		for j := i + 1; j < len(neighbors); j++ {
			g.Edges[neighbors[i]] = append(g.Edges[neighbors[i]], neighbors[j])
			g.Edges[neighbors[j]] = append(g.Edges[neighbors[j]], neighbors[i])
		}
	}

	// Remove the node from the Nodes map
	delete(g.Nodes, node)

	// Remove the node from the Edges map and update neighbors' adjacency lists
	for neighbor, edges := range g.Edges {
		g.Edges[neighbor] = DeleteFromSlice(edges, node)
	}

	// Delete the entry for the removed node from the Edges map
	delete(g.Edges, node)
}

func (g *UndirectedGraph) ContractEdge(edge Edge) {
	node1 := edge.Node1
	node2 := edge.Node2
	//assign all of edges directed towards node1 to node2
	neighbors := g.Edges[node1]
	for i := 0; i < len(neighbors); i++ {
		g.AddEdge(Edge{
			Node1: neighbors[i],
			Node2: node2,
		})
	}

	// Remove the node from the Nodes map
	delete(g.Nodes, node1)

	// Remove the node from the Edges map and update neighbors' adjacency lists
	for neighbor, edges := range g.Edges {
		g.Edges[neighbor] = DeleteFromSlice(edges, node1)
	}

	// Delete the entry for the removed node from the Edges map
	delete(g.Edges, node1)
}

// ConnectedComponents finds the connected components in an undirected graph.
// It takes an undirected graph (g) as input and returns a Components struct.
// The Components struct contains an array of UndirectedGraphs, each representing
// a connected component in the input graph.
//
// Parameters:
//   - g: Pointer to an UndirectedGraph representing the input graph.
//
// Returns:
//   - components: Components struct containing an array of UndirectedGraphs,
//     each representing a connected component in the input graph. The biggest
//     connected component's index is stored in BiggestComponentIdx field of
//     the Components struct.
//
// The function iterates through each node in the input graph and performs DFS
// traversal from each unvisited node to identify connected components. It stores
// each connected component as an UndirectedGraph in the Components struct.
//
// Example usage:
//
//	graph := UndirectedGraph{
//	    Nodes: make(map[Node]bool),
//	    Edges: make(map[Node][]Node),
//	}
//	edges := []Edge{
//	    {Node1: 1, Node2: 2},
//	    {Node1: 3, Node2: 4},
//	    {Node1: 5, Node2: 6},
//	}
//	for _, edge := range edges {
//	    graph.AddEdge(edge)
//	}
//	components := ConnectedComponents(&graph)
//
//	// Access each connected component
//	for _, component := range components.ComponentsArray {
//	    fmt.Println("Connected component:")
//	    for node := range component.Nodes {
//	        fmt.Println(node)
//	    }
//	}
func ConnectedComponents(g *UndirectedGraph) (components Components) {
	components = Components{
		ComponentsArray:     make([]*UndirectedGraph, 0),
		visitedNodes:        make(map[Node]bool),
		BiggestComponentIdx: -1,
	}

	for node := range g.Nodes {
		if components.visitedNodes[node] {
			continue
		}
		components.visitedNodes[node] = true
		component := g.DFS(node)
		components.AddComponent(component)

	}
	return components
}
