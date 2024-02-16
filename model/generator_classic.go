package model

import "fmt"

// CompleteGraph generates a complete graph with the specified number of nodes.
// A complete graph is a simple undirected graph where each pair of distinct nodes is connected by a unique edge.
// The graph is represented by an UndirectedGraph object.
//
// Parameters:
//
//	numberOfNodes: The number of nodes in the complete graph.
//
// Returns:
//
//	An UndirectedGraph representing the complete graph with the specified number of nodes.
//
// Example:
//
//	// Generate a complete graph with 4 nodes
//	graph := CompleteGraph(4)
func CompleteGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
	for i := 0; i < numberOfNodes; i++ {
		for j := i + 1; j < numberOfNodes; j++ {
			g.AddEdge(Edge{
				Node1: Node(i),
				Node2: Node(j),
			})
		}
	}
	return g
}

// LadderGraph returns the Ladder graph of length n and 2n nodes
func LadderGraph(nodesInSinglePath int) *UndirectedGraph {
	g := &UndirectedGraph{}

	// Generate and add edges for the ladder structure
	for _, edge := range Pairwise(Range(nodesInSinglePath, 2*nodesInSinglePath)) {
		g.AddEdge(edge)
	}

	// Add rung edges between the two paths of the ladder
	for i := 0; i < nodesInSinglePath; i++ { // nodesInSinglePath = 3
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node(i + nodesInSinglePath),
		})

		if i != nodesInSinglePath-1 { // i != 2
			g.AddEdge(Edge{
				Node1: Node(i),
				Node2: Node(i + 1),
			})
		}
	}

	return g
}

// CircularLadderGraph returns the circular ladder graph CL_n of length n
func CircularLadderGraph(nodesInSinglePath int) (*UndirectedGraph, error) {
	if nodesInSinglePath < 3 {
		return nil, fmt.Errorf("nodesInSinglePath must be at least 3")
	}

	g := LadderGraph(nodesInSinglePath)
	lastNode := Node(nodesInSinglePath - 1)
	g.AddEdge(Edge{
		Node1: 0,
		Node2: lastNode,
	})
	g.AddEdge(Edge{
		Node1: Node(nodesInSinglePath),
		Node2: 2*Node(nodesInSinglePath) - 1,
	})
	return g, nil
}

// WheelGraph returns the wheel graph
func WheelGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
	g.AddNode(0)
	for i := 1; i < numberOfNodes; i++ {
		g.AddEdge(Edge{
			Node1: Node(i - 1),
			Node2: Node(i),
		})
		g.AddEdge(Edge{
			Node1: Node(0),
			Node2: Node(i),
		})
	}
	return g
}

// TuranGraph returns the Turán graph
func TuranGraph(numberOfNodes int, numberOfPartitions int) *UndirectedGraph {
	g := &UndirectedGraph{}

	numberOfPartitionsA := numberOfPartitions - (numberOfNodes % numberOfPartitions)
	sizeOfPartitionsA := numberOfNodes / numberOfPartitions
	numberOfPartitionsB := numberOfNodes % numberOfPartitions
	sizeOfPartitionsB := numberOfNodes / (numberOfPartitions + 1)

	partitionsA := make(map[int]map[Node]bool, 0)
	partitionsB := make(map[int]map[Node]bool, 0)

	nodeId := 0
	for p := 0; p < numberOfPartitionsA; p++ {
		for n := 0; n < sizeOfPartitionsA; n++ {
			g.AddNode(Node(nodeId))
			partitionsA[p][Node(nodeId)] = true
			nodeId = nodeId + 1
		}
	}

	for p := 0; p < numberOfPartitionsB; p++ {
		for n := 0; n < sizeOfPartitionsB; n++ {
			g.AddNode(Node(nodeId))
			partitionsB[p][Node(nodeId)] = true
			nodeId = nodeId + 1
		}
	}

	// for nodes in partitions
	// generate connections to nodes outside the partition
	for p := 0; p < numberOfPartitionsA; p++ {
		for node := range partitionsA[p] {
			for i := 0; i < numberOfNodes; i++ {
				if !partitionsA[p][Node(i)] {
					g.AddEdge(Edge{
						Node1: node,
						Node2: Node(i),
					})
				}
			}
		}
	}

	for p := 0; p < numberOfPartitionsB; p++ {
		for node := range partitionsB[p] {
			for i := 0; i < numberOfNodes; i++ {
				if !partitionsB[p][Node(i)] {
					g.AddEdge(Edge{
						Node1: node,
						Node2: Node(i),
					})
				}
			}
		}
	}

	return g
}

// TrivialGraph returns a graph with one node (with label 0) and no edges
func TrivialGraph() *UndirectedGraph {
	g := &UndirectedGraph{}
	g.AddNode(0)
	return g
}

// NullGraph returns a graph without nodes and edges
func NullGraph() *UndirectedGraph {
	g := &UndirectedGraph{}
	return g
}

// TadpoleGraph returns a Tadpole graph consisting of a cycle graph on cycleSize (at least 3) vertices and a path graph of pathSize vertices, connected with a bridge.
func TadpoleGraph(cycleSize int, pathSize int) (*UndirectedGraph, error) {
	if cycleSize < 3 {
		return nil, fmt.Errorf("cycle size can't be < 3")
	}
	g := &UndirectedGraph{}
	// generate cycle graph
	for i := 0; i < cycleSize; i++ {
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node((i + 1) % cycleSize),
		})
	}
	for i := cycleSize; i < cycleSize+pathSize; i++ {
		g.AddEdge(Edge{
			Node1: Node(i - 1),
			Node2: Node(i),
		})
	}
	return g, nil
}

// StarGraph returns a star graph.
func StarGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
	// generate a star graph
	for i := 1; i < numberOfNodes; i++ {
		g.AddEdge(Edge{
			Node1: Node(0),
			Node2: Node(i),
		})
	}
	return g
}

// PathGraph returns a path graph.
func PathGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
	// generate a path graph
	for i := 1; i < numberOfNodes; i++ {
		g.AddEdge(Edge{
			Node1: Node(i - 1),
			Node2: Node(i),
		})
	}
	return g
}

// LollipopGraph returns a path graph.
func LollipopGraph(completeGraphSize int, pathGraphSize int) *UndirectedGraph {
	g := &UndirectedGraph{}
	// generate a Lollipop graph
	for i := 0; i < completeGraphSize; i++ {
		for j := i + 1; j < completeGraphSize; j++ {
			g.AddEdge(Edge{
				Node1: Node(i),
				Node2: Node(j),
			})
		}
	}
	for i := completeGraphSize; i < completeGraphSize+pathGraphSize; i++ {
		g.AddEdge(Edge{
			Node1: Node(i - 1),
			Node2: Node(i),
		})
	}
	return g
}

// CycleGraph returns a cyrcle graph.
func CycleGraph(numberOfNodes int) *UndirectedGraph {
	g := &UndirectedGraph{}
	// generate a Cycle graph
	for i := 0; i < numberOfNodes; i++ {
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node((i + 1) % numberOfNodes),
		})
	}
	return g
}

// CirculantGraph returns a circulant graph of n nodes and .
func CirculantGraph(numberOfNodes int, offset int) *UndirectedGraph {
	g := &UndirectedGraph{}
	// generate a Circulant graph
	for i := 0; i < numberOfNodes; i++ {
		g.AddEdge(Edge{
			Node1: Node(i),
			Node2: Node((i + offset) % numberOfNodes),
		})
	}

	return g
}

// TODO: balanced tree, binomial tree, barbell graph, complete multipartite graph, dorogovtsev goltsev mendes graph, full rary tree
