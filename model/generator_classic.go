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
func TuranGraph(n, r int) *UndirectedGraph {
	if r <= 0 || n < 0 {
		return &UndirectedGraph{Nodes: map[Node]bool{}, Edges: map[Node][]Node{}}
	}

	g := &UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	// Partition sizes: first 'rem' partitions have size q+1, others have size q.
	q := n / r
	rem := n % r

	partitions := make([][]Node, r)
	id := 0
	for p := 0; p < r; p++ {
		size := q
		if p < rem {
			size = q + 1
		}
		part := make([]Node, 0, size)
		for i := 0; i < size; i++ {
			u := Node(id)
			g.AddNode(u)
			part = append(part, u)
			id++
		}
		partitions[p] = part
	}

	// Add edges between every pair of *different* partitions (complete r-partite)
	for a := 0; a < r; a++ {
		for b := a + 1; b < r; b++ {
			for _, u := range partitions[a] {
				for _, v := range partitions[b] {
					g.AddEdge(Edge{Node1: u, Node2: v})
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
