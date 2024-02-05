package model

func Pairwise(nodeIds []int) (e []Edge) {
	e = []Edge{}
	for i := 0; i < len(nodeIds)-1; i++ {
		e = append(e, Edge{Node(nodeIds[i]), Node(nodeIds[i+1])})
	}
	return e
}

func Range(start int, end int) []int {
	values := []int{}
	for i := 0; i < end-start; i++ {
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
