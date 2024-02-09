package model

import (
	"fmt"
	"testing"
)

func TestCompleteGraph(t *testing.T) {
	tests := []struct {
		name           string
		numberOfNodes  int
		expectedEdges  int
		expectedDegree int
	}{
		{
			name:           "CompleteGraph with 0 nodes",
			numberOfNodes:  0,
			expectedEdges:  0,
			expectedDegree: 0,
		},
		{
			name:           "CompleteGraph with 2 nodes",
			numberOfNodes:  2,
			expectedEdges:  1, // 2C1 (combination)
			expectedDegree: 1,
		},
		{
			name:           "CompleteGraph with 5 nodes",
			numberOfNodes:  5,
			expectedEdges:  10, // 5C2 (combination)
			expectedDegree: 4,
		},
		{
			name:           "CompleteGraph with 10 nodes",
			numberOfNodes:  10,
			expectedEdges:  45, // 10C2 (combination)
			expectedDegree: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := CompleteGraph(tt.numberOfNodes)

			// Check the number of edges
			actualEdges := g.NumberOfEdges()
			if actualEdges != tt.expectedEdges {
				t.Errorf("Expected %d edges, but got %d", tt.expectedEdges, actualEdges)
			}

			// Check the degree of each node
			for i := 0; i < tt.numberOfNodes; i++ {
				actualDegree := g.NodeDegree(Node(i))
				if actualDegree != tt.expectedDegree {
					t.Errorf("Expected degree of node %d to be %d, but got %d", i, tt.expectedDegree, actualDegree)
				}
			}
		})
	}
}

func TestLadderGraph(t *testing.T) {
	testCases := []struct {
		nodes         int
		expectedNodes map[Node]bool
		expectedEdges map[Node][]Node
	}{
		{
			nodes:         0,
			expectedNodes: map[Node]bool{},
			expectedEdges: map[Node][]Node{},
		},
		{
			nodes: 2,
			expectedNodes: map[Node]bool{
				0: true, 1: true, 2: true, 3: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1, 2},
				1: {0, 3},
				2: {0, 3},
				3: {1, 2},
			},
		},
		{
			nodes: 3,
			expectedNodes: map[Node]bool{
				0: true, 1: true, 2: true, 3: true, 4: true, 5: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1, 3},
				1: {0, 2, 4},
				2: {1, 5},
				3: {0, 4},
				4: {1, 3, 5},
				5: {2, 4},
			},
		},
		{
			nodes: 5,
			expectedNodes: map[Node]bool{
				0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1, 5},
				1: {0, 2, 6},
				2: {1, 3, 7},
				3: {2, 4, 8},
				4: {3, 9},
				5: {0, 6},
				6: {5, 1, 7},
				7: {6, 2, 8},
				8: {7, 3, 9},
				9: {8}, // Adjusted to connect with node 8 only
			},
		},
		{
			nodes: 1,
			expectedNodes: map[Node]bool{
				0: true, 1: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1},
				1: {0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Nodes: %d", tc.nodes), func(t *testing.T) {
			g := LadderGraph(tc.nodes)
			validateGraph(t, g, tc.expectedNodes, tc.expectedEdges)
		})
	}
}

func TestCircularLadderGraph(t *testing.T) {
	// Test case 1: Basic Test
	nodesInSinglePath := 4
	g := CircularLadderGraph(nodesInSinglePath)
	expectedNodes := map[Node]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true}
	expectedEdges := map[Node][]Node{
		1: {2, 6},
		2: {1, 3},
		3: {2, 4},
		4: {3, 5},
		5: {4, 6},
		6: {1, 5},
	}
	validateGraph(t, g, expectedNodes, expectedEdges)

	// Test case 2: Edge Case Test
	nodesInSinglePath = 2
	g = CircularLadderGraph(nodesInSinglePath)
	expectedNodes = map[Node]bool{1: true, 2: true, 3: true}
	expectedEdges = map[Node][]Node{
		1: {2, 3},
		2: {1, 3},
		3: {1, 2},
	}
	validateGraph(t, g, expectedNodes, expectedEdges)

	// Test case 3: Larger Graph Test
	nodesInSinglePath = 6
	g = CircularLadderGraph(nodesInSinglePath)
	expectedNodes = map[Node]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 11: true, 12: true}
	expectedEdges = map[Node][]Node{
		1:  {2, 12},
		2:  {1, 3},
		3:  {2, 4},
		4:  {3, 5},
		5:  {4, 6},
		6:  {5, 7},
		7:  {6, 8},
		8:  {7, 9},
		9:  {8, 10},
		10: {9, 11},
		11: {10, 12},
		12: {1, 11},
	}
	validateGraph(t, g, expectedNodes, expectedEdges)

	// Test case 4: Connectivity Test
	// Validate if all nodes are properly connected
	for node, neighbors := range g.Edges {
		for _, neighbor := range neighbors {
			if !contains(g.Edges[neighbor], node) {
				t.Errorf("Test case 4 failed: Nodes %d and %d are not properly connected", node, neighbor)
			}
		}
	}

	// Test case 5: Node Existence Test
	// Validate if all expected nodes exist
	for node := range expectedNodes {
		if !g.Nodes[node] {
			t.Errorf("Test case 5 failed: Expected node %d does not exist in the generated graph", node)
		}
	}
}

// Helper function to validate the generated graph
func validateGraph(t *testing.T, g *UndirectedGraph, expectedNodes map[Node]bool, expectedEdges map[Node][]Node) {
	// Validate nodes
	if len(g.Nodes) != len(expectedNodes) {
		t.Errorf("Nodes mismatch, expected: %v, got: %v", expectedNodes, g.Nodes)
	} else {
		for node := range expectedNodes {
			if !g.Nodes[node] {
				t.Errorf("Nodes mismatch, expected: %v, got: %v", expectedNodes, g.Nodes)
			}
		}
	}
	// Validate edges
	for node, edges := range expectedEdges {
		for _, edge := range edges {
			if !contains(g.Edges[node], edge) {
				t.Errorf("Edges mismatch for node %d, expected: %v, got: %v. Full graph: %v", node, edges, g.Edges[node], g)
			}
		}
	}
}

// Helper function to check if a node exists in a slice
func contains(slice []Node, value Node) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
