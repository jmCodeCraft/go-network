package model

import (
	"reflect"
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
	// Test case 1: Ladder graph with 2 nodes in each path
	result := LadderGraph(2)
	expected := &UndirectedGraph{
		Nodes: map[Node]bool{Node(0): true, Node(1): true, Node(2): true, Node(3): true},
		Edges: map[Node][]Node{
			Node(0): {Node(1)},
			Node(1): {Node(0), Node(2)},
			Node(2): {Node(1), Node(3)},
			Node(3): {Node(2)},
		},
	}
	if result.NumberOfEdges() != 2*len(result.Nodes) {
		t.Errorf("graph has no 2n vertices. Num of edges: %v; num of nodes: %v", result.NumberOfEdges(), len(result.Nodes))
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected\n%v\n, but got\n%v", expected, result)
	}

	// Test case 2: Ladder graph with 3 nodes in each path
	result = LadderGraph(3)
	expected = &UndirectedGraph{
		Nodes: map[Node]bool{Node(0): true, Node(1): true, Node(2): true, Node(3): true, Node(4): true, Node(5): true},
		Edges: map[Node][]Node{
			Node(0): {Node(1)},
			Node(1): {Node(0), Node(2)},
			Node(2): {Node(1), Node(3)},
			Node(3): {Node(2), Node(4)},
			Node(4): {Node(3), Node(5)},
			Node(5): {Node(4)},
		},
	}
	if result.NumberOfEdges() != 2*len(result.Nodes) {
		t.Errorf("graph has no 2n vertices. Num of edges: %v; num of nodes: %v", result.NumberOfEdges(), len(result.Nodes))
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected\n%v\n, but got\n%v", expected, result)
	}

	// Test case 3: Ladder graph with 0 nodes in each path (empty graph)
	result = LadderGraph(0)
	expected = &UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected\n%v\n, but got\n%v", expected, result)
	}

	// Add more test cases as needed...
}
