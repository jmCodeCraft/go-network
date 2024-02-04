package model

import (
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
