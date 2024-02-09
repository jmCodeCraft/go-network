package model

import (
	"reflect"
	"sort"
	"testing"
)

func TestUndirectedGraph_AddEdge(t *testing.T) {
	// Test case 1: Adding an edge to an empty graph
	graph1 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	edge1 := Edge{Node1: 1, Node2: 2}
	graph1.AddEdge(edge1)

	expectedEdges1 := map[Node][]Node{
		1: {2},
		2: {1},
	}

	if !reflect.DeepEqual(graph1.Edges, expectedEdges1) {
		t.Errorf("Expected %v, but got %v", expectedEdges1, graph1.Edges)
	}

	// Test case 2: Adding an edge to a non-empty graph
	graph2 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
		},
		Edges: map[Node][]Node{
			1: {2},
			2: {1},
		},
	}

	edge2 := Edge{Node1: 2, Node2: 3}
	graph2.AddEdge(edge2)

	expectedEdges2 := map[Node][]Node{
		1: {2},
		2: {1, 3},
		3: {2},
	}

	if !reflect.DeepEqual(graph2.Edges, expectedEdges2) {
		t.Errorf("Expected %v, but got %v", expectedEdges2, graph2.Edges)
	}

	// Test case 3: Adding an edge with existing nodes
	graph3 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: map[Node][]Node{
			1: {2},
			2: {1, 3},
			3: {2},
		},
	}

	edge3 := Edge{Node1: 3, Node2: 1}
	graph3.AddEdge(edge3)

	expectedEdges3 := map[Node][]Node{
		1: {2, 3},
		2: {1, 3},
		3: {2, 1},
	}

	if !reflect.DeepEqual(graph3.Edges, expectedEdges3) {
		t.Errorf("Expected %v, but got %v", expectedEdges3, graph3.Edges)
	}
}

func TestUndirectedGraph_AddNode(t *testing.T) {
	// Test case 1: Adding a node to an empty graph
	graph1 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	node1 := Node(1)
	graph1.AddNode(node1)

	expectedNodes1 := map[Node]bool{
		1: true,
	}

	if !reflect.DeepEqual(graph1.Nodes, expectedNodes1) {
		t.Errorf("Expected %v, but got %v", expectedNodes1, graph1.Nodes)
	}

	// Test case 2: Adding a node to a non-empty graph
	graph2 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
		},
		Edges: make(map[Node][]Node),
	}

	node2 := Node(3)
	graph2.AddNode(node2)

	expectedNodes2 := map[Node]bool{
		1: true,
		2: true,
		3: true,
	}

	if !reflect.DeepEqual(graph2.Nodes, expectedNodes2) {
		t.Errorf("Expected %v, but got %v", expectedNodes2, graph2.Nodes)
	}

	// Test case 3: Adding an existing node
	graph3 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: make(map[Node][]Node),
	}

	node3 := Node(2)
	graph3.AddNode(node3)

	expectedNodes3 := map[Node]bool{
		1: true,
		2: true,
		3: true,
	}

	if !reflect.DeepEqual(graph3.Nodes, expectedNodes3) {
		t.Errorf("Expected %v, but got %v", expectedNodes3, graph3.Nodes)
	}
}

func TestUndirectedGraph_AddNodes(t *testing.T) {
	// Test case 1: Adding nodes to an empty graph
	graph1 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	nodes1 := []Node{1, 2, 3}
	graph1.AddNodes(nodes1)

	expectedNodes1 := map[Node]bool{
		1: true,
		2: true,
		3: true,
	}

	if !reflect.DeepEqual(graph1.Nodes, expectedNodes1) {
		t.Errorf("Expected %v, but got %v", expectedNodes1, graph1.Nodes)
	}

	// Test case 2: Adding nodes to a non-empty graph
	graph2 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
		},
		Edges: make(map[Node][]Node),
	}

	nodes2 := []Node{3, 4}
	graph2.AddNodes(nodes2)

	expectedNodes2 := map[Node]bool{
		1: true,
		2: true,
		3: true,
		4: true,
	}

	if !reflect.DeepEqual(graph2.Nodes, expectedNodes2) {
		t.Errorf("Expected %v, but got %v", expectedNodes2, graph2.Nodes)
	}

	// Test case 3: Adding an empty slice of nodes
	graph3 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	nodes3 := []Node{}
	graph3.AddNodes(nodes3)

	expectedNodes3 := map[Node]bool{}

	if !reflect.DeepEqual(graph3.Nodes, expectedNodes3) {
		t.Errorf("Expected %v, but got %v", expectedNodes3, graph3.Nodes)
	}
}

func TestUndirectedGraph_NodeDegree(t *testing.T) {
	// Test case 1: Node with incident edges
	graph1 := UndirectedGraph{
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

	degree1 := graph1.NodeDegree(1)
	expectedDegree1 := 2

	if degree1 != expectedDegree1 {
		t.Errorf("Expected degree %v for node 1, but got %v", expectedDegree1, degree1)
	}

	// Test case 2: Node with no incident edges
	graph2 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
			3: true,
		},
		Edges: map[Node][]Node{},
	}

	degree2 := graph2.NodeDegree(2)
	expectedDegree2 := 0

	if degree2 != expectedDegree2 {
		t.Errorf("Expected degree %v for node 2, but got %v", expectedDegree2, degree2)
	}

	// Test case 3: Node not present in the graph
	graph3 := UndirectedGraph{
		Nodes: map[Node]bool{
			1: true,
			2: true,
		},
		Edges: map[Node][]Node{
			1: {2},
			2: {1},
		},
	}

	degree3 := graph3.NodeDegree(3)
	expectedDegree3 := 0

	if degree3 != expectedDegree3 {
		t.Errorf("Expected degree %v for node 3, but got %v", expectedDegree3, degree3)
	}
}

func TestUndirectedGraph_HasNode(t *testing.T) {
	// Test case 1: Node exists in the graph
	graph1 := UndirectedGraph{
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

	result1 := graph1.HasNode(2)
	expectedResult1 := true

	if result1 != expectedResult1 {
		t.Errorf("Expected %v, but got %v", expectedResult1, result1)
	}

	// Test case 2: Node does not exist in the graph
	graph2 := UndirectedGraph{
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

	result2 := graph2.HasNode(4)
	expectedResult2 := false

	if result2 != expectedResult2 {
		t.Errorf("Expected %v, but got %v", expectedResult2, result2)
	}

	// Test case 3: Empty graph
	graph3 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	result3 := graph3.HasNode(1)
	expectedResult3 := false

	if result3 != expectedResult3 {
		t.Errorf("Expected %v, but got %v", expectedResult3, result3)
	}
}

func TestUndirectedGraph_RemoveEdge(t *testing.T) {
	// Test case 1: Removing an existing edge
	graph1 := UndirectedGraph{
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

	edgeToRemove1 := Edge{Node1: 1, Node2: 2}
	graph1.RemoveEdge(edgeToRemove1)

	expectedEdges1 := map[Node][]Node{
		1: {3},
		2: {3},
		3: {1, 2},
	}

	if !reflect.DeepEqual(graph1.Edges, expectedEdges1) {
		t.Errorf("Expected %v, but got %v", expectedEdges1, graph1.Edges)
	}

	// Test case 2: Removing a non-existing edge
	graph2 := UndirectedGraph{
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

	edgeToRemove2 := Edge{Node1: 1, Node2: 4}
	graph2.RemoveEdge(edgeToRemove2)

	expectedEdges2 := map[Node][]Node{
		1: {2, 3},
		2: {1, 3},
		3: {1, 2},
	}

	if !reflect.DeepEqual(graph2.Edges, expectedEdges2) {
		t.Errorf("Expected %v, but got %v", expectedEdges2, graph2.Edges)
	}

	// Test case 3: Removing an edge from an empty graph
	graph3 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	edgeToRemove3 := Edge{Node1: 1, Node2: 2}
	graph3.RemoveEdge(edgeToRemove3)

	expectedEdges3 := map[Node][]Node{}

	if !reflect.DeepEqual(graph3.Edges, expectedEdges3) {
		t.Errorf("Expected %v, but got %v", expectedEdges3, graph3.Edges)
	}
}

func TestUndirectedGraph_RemoveNode(t *testing.T) {
	// Test case 1: Removing an existing node with associated edges
	graph1 := UndirectedGraph{
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

	nodeToRemove1 := Node(2)
	graph1.RemoveNode(nodeToRemove1)

	expectedNodes1 := map[Node]bool{
		1: true,
		3: true,
	}

	expectedEdges1 := map[Node][]Node{
		1: {3},
		3: {1},
	}

	if !reflect.DeepEqual(graph1.Nodes, expectedNodes1) {
		t.Errorf("Expected %v, but got %v", expectedNodes1, graph1.Nodes)
	}

	if !reflect.DeepEqual(graph1.Edges, expectedEdges1) {
		t.Errorf("Expected %v, but got %v", expectedEdges1, graph1.Edges)
	}

	// Test case 2: Removing a non-existing node
	graph2 := UndirectedGraph{
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

	nodeToRemove2 := Node(4)
	graph2.RemoveNode(nodeToRemove2)

	expectedNodes2 := map[Node]bool{
		1: true,
		2: true,
		3: true,
	}

	expectedEdges2 := map[Node][]Node{
		1: {2, 3},
		2: {1, 3},
		3: {1, 2},
	}

	if !reflect.DeepEqual(graph2.Nodes, expectedNodes2) {
		t.Errorf("Expected %v, but got %v", expectedNodes2, graph2.Nodes)
	}

	if !reflect.DeepEqual(graph2.Edges, expectedEdges2) {
		t.Errorf("Expected %v, but got %v", expectedEdges2, graph2.Edges)
	}

	// Test case 3: Removing a node from an empty graph
	graph3 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	nodeToRemove3 := Node(1)
	graph3.RemoveNode(nodeToRemove3)

	expectedNodes3 := map[Node]bool{}
	expectedEdges3 := map[Node][]Node{}

	if !reflect.DeepEqual(graph3.Nodes, expectedNodes3) {
		t.Errorf("Expected %v, but got %v", expectedNodes3, graph3.Nodes)
	}

	if !reflect.DeepEqual(graph3.Edges, expectedEdges3) {
		t.Errorf("Expected %v, but got %v", expectedEdges3, graph3.Edges)
	}
}
func TestDFS(t *testing.T) {
	// Test case 1: Regular graph traversal
	graph1 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	edges1 := []Edge{
		{Node1: 1, Node2: 2},
		{Node1: 1, Node2: 3},
		{Node1: 2, Node2: 4},
		{Node1: 2, Node2: 5},
		{Node1: 3, Node2: 6},
	}
	for _, edge := range edges1 {
		graph1.AddEdge(edge)
	}
	visitedGraph1 := graph1.DFS(1)
	visitedNodes1 := getSortedNodes(visitedGraph1)
	expectedVisitedNodes1 := []Node{1, 2, 3, 4, 5, 6}
	if !sliceEqual(visitedNodes1, expectedVisitedNodes1) {
		t.Errorf("Test case 1 failed: Visited nodes mismatch, expected: %v, got: %v", expectedVisitedNodes1, visitedNodes1)
	}

	// Test case 2: Graph with isolated nodes
	graph2 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	edges2 := []Edge{
		{Node1: 1, Node2: 2},
		{Node1: 2, Node2: 3},
		{Node1: 4, Node2: 5},
	}
	for _, edge := range edges2 {
		graph2.AddEdge(edge)
	}
	visitedGraph2 := graph2.DFS(1)
	visitedNodes2 := getSortedNodes(visitedGraph2)
	expectedVisitedNodes2 := []Node{1, 2, 3}
	if !sliceEqual(visitedNodes2, expectedVisitedNodes2) {
		t.Errorf("Test case 2 failed: Visited nodes mismatch, expected: %v, got: %v", expectedVisitedNodes2, visitedNodes2)
	}

	// Test case 3: Empty graph
	graph3 := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	visitedGraph3 := graph3.DFS(1)
	visitedNodes3 := getSortedNodes(visitedGraph3)
	if len(visitedNodes3) != 0 {
		t.Errorf("Test case 3 failed: Visited nodes should be empty, got: %v", visitedNodes3)
	}
}

func getSortedNodes(graph *UndirectedGraph) []Node {
	visitedNodes := make([]Node, 0, len(graph.Nodes))
	for node := range graph.Nodes {
		visitedNodes = append(visitedNodes, node)
	}
	sort.Slice(visitedNodes, func(i, j int) bool { return visitedNodes[i] < visitedNodes[j] })
	return visitedNodes
}

func sliceEqual(a, b []Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestConnectedComponents(t *testing.T) {
	testCases := []struct {
		name               string
		graph              UndirectedGraph
		expectedComponents []*UndirectedGraph
	}{
		{
			name: "Regular graph with multiple connected components",
			graph: UndirectedGraph{
				Nodes: make(map[Node]bool),
				Edges: make(map[Node][]Node),
			},
			expectedComponents: []*UndirectedGraph{
				{
					Nodes: map[Node]bool{1: true, 2: true, 3: true},
					Edges: map[Node][]Node{1: {2}, 2: {1, 3}, 3: {2}},
				},
				{
					Nodes: map[Node]bool{4: true, 5: true},
					Edges: map[Node][]Node{4: {5}, 5: {4}},
				},
				{
					Nodes: map[Node]bool{6: true, 7: true},
					Edges: map[Node][]Node{6: {7}, 7: {6}},
				},
			},
		},
		{
			name: "Graph with one connected component",
			graph: UndirectedGraph{
				Nodes: make(map[Node]bool),
				Edges: make(map[Node][]Node),
			},
			expectedComponents: []*UndirectedGraph{
				{
					Nodes: map[Node]bool{1: true, 2: true, 3: true, 4: true},
					Edges: map[Node][]Node{
						1: {2},
						2: {1, 3},
						3: {2, 4},
						4: {3},
					},
				},
			},
		},
		{
			name: "Empty graph",
			graph: UndirectedGraph{
				Nodes: make(map[Node]bool),
				Edges: make(map[Node][]Node),
			},
			expectedComponents: []*UndirectedGraph{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			components := ConnectedComponents(&tc.graph)
			for _, computedComponent := range components.ComponentsArray {
				found := false
				for _, expectedComponent := range tc.expectedComponents {
					if computedComponent.Equals(expectedComponent) {
						found = true
					}
				}
				if !found {
					t.Errorf("Connected component %v not found in expected components", computedComponent)
				}
			}
		})
	}
}
