package model

// Pairwise creates a slice of edges representing pairwise connections between consecutive nodes in the given sequence.
//
// The Pairwise function takes a slice of node identifiers and generates edges connecting each node
// to its consecutive neighbor in the sequence. It returns a slice of edges capturing the pairwise connections.
//
// Parameters:
//   - nodeIds: A slice of integers representing node identifiers.
//
// Returns:
//
//	A slice of Edges representing pairwise connections between consecutive nodes in the sequence.
func Pairwise(nodeIds []int) []Edge {
	if len(nodeIds) == 0 {
		return []Edge{}
	}

	edges := make([]Edge, 0, len(nodeIds)-1)
	for i := 0; i < len(nodeIds)-1; i++ {
		edges = append(edges, Edge{Node(nodeIds[i]), Node(nodeIds[i+1])})
	}
	return edges
}

// Range generates a sequence of integers starting from 'start' (inclusive) to 'end' (exclusive).
//
// The Range function creates a slice containing integers in the range [start, end).
//
// Parameters:
//   - start: The starting value of the range (inclusive).
//   - end: The ending value of the range (exclusive).
//
// Returns:
//
//	A slice of integers representing the sequence from 'start' to 'end-1'.
func Range(start, end int) []int {
	values := make([]int, end-start)
	for i := 0; i < len(values); i++ {
		values[i] = start + i
	}
	return values
}
