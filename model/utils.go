package model

type WeightedElement struct {
	Payload any
	Weight  float32
}

// Pairwise generates a list of edges connecting consecutive nodes in the given list of node IDs.
// It takes a slice of node IDs and returns a slice of edges, where each edge connects a node
// with the next node in the input slice.
// For example, if nodeIds is [1, 2, 3, 4], Pairwise will return [{1, 2}, {2, 3}, {3, 4}].
func Pairwise(nodeIds []int) []Edge {
	if len(nodeIds) == 0 {
		return []Edge{}
	}

	edges := make([]Edge, len(nodeIds)-1)
	for i := 0; i < len(nodeIds)-1; i++ {
		edges[i] = Edge{Node(nodeIds[i]), Node(nodeIds[i+1])}
	}
	return edges
}

// Range returns a slice containing integers from 'start' (inclusive) to 'end' (exclusive).
// It generates a sequence of integers starting from 'start' and ending at 'end-1'.
// If 'start' is greater than or equal to 'end', an empty slice is returned.
func Range(start int, end int) []int {
	values := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		values = append(values, i)
	}
	return values
}

func DeleteFromSlice(slice []Node, objectToRemove Node) []Node {
	newSlice := []Node{}
	for i := 0; i < len(slice); i++ {
		if slice[i] != objectToRemove {
			newSlice = append(newSlice, slice[i])
		}
	}
	return newSlice
}

func GetDictKeys(dict map[Node]bool) []Node {
	keys := make([]Node, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}
