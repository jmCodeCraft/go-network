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
				9: {8, 4},
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
	testCases := []struct {
		nodesInSinglePath int
		expectedNodes     map[Node]bool
		expectedEdges     map[Node][]Node
	}{
		{
			nodesInSinglePath: 3,
			expectedNodes: map[Node]bool{
				0: true, 1: true, 2: true,
				3: true, 4: true, 5: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1, 2, 3},
				1: {0, 2, 4},
				2: {0, 1, 5},
				3: {0, 4, 5},
				4: {1, 3, 5},
				5: {4, 3, 2},
			},
		},
		{
			nodesInSinglePath: 4,
			expectedNodes: map[Node]bool{
				0: true, 1: true, 2: true, 3: true,
				4: true, 5: true, 6: true, 7: true,
			},
			expectedEdges: map[Node][]Node{
				0: {1, 4, 3},
				1: {0, 2, 5},
				2: {1, 3, 6},
				3: {0, 2, 7},
				4: {0, 5, 7},
				5: {1, 4, 6},
				6: {2, 5, 7},
				7: {3, 4, 6},
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("NodesInSinglePath=%d", tc.nodesInSinglePath), func(t *testing.T) {
			graph, _ := CircularLadderGraph(tc.nodesInSinglePath)
			validateGraph(t, graph, tc.expectedNodes, tc.expectedEdges)
		})
	}
}

func TestCircularLadderGraph_Error(t *testing.T) {
	testCases := []struct {
		nodesInSinglePath int
		expectedError     string
	}{
		{nodesInSinglePath: 2, expectedError: "nodesInSinglePath must be at least 3"},
		{nodesInSinglePath: 0, expectedError: "nodesInSinglePath must be at least 3"},
		{nodesInSinglePath: -5, expectedError: "nodesInSinglePath must be at least 3"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("NodesInSinglePath=%d", tc.nodesInSinglePath), func(t *testing.T) {
			graph, err := CircularLadderGraph(tc.nodesInSinglePath)
			if err == nil {
				t.Errorf("Expected an error, but got nil")
			} else if err.Error() != tc.expectedError {
				t.Errorf("Unexpected error message, expected: %s, got: %s", tc.expectedError, err.Error())
			}
			if graph != nil {
				t.Errorf("Expected nil graph, but got %+v", graph)
			}
		})
	}
}

// Helper function to validate the generated graph
func validateGraph(t *testing.T, g *UndirectedGraph, expectedNodes map[Node]bool, expectedEdges map[Node][]Node) {
	expectedGraph := &UndirectedGraph{Nodes: expectedNodes, Edges: expectedEdges}
	if !g.Equals(expectedGraph) {
		t.Errorf("Graph mismatch, expected: %v, got: %v", expectedGraph, g)
	}
}
