package model

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/jinzhu/copier"
	"github.com/mroth/weightedrand"
)

type IDeletionSamplingStrategy interface {
	SamplingStage(graph *UndirectedGraph, howManyToDelete int) error
	Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (*UndirectedGraph, error)
}

type ISamplingStrategy interface {
	Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (*UndirectedGraph, error)
}

type DeletionSamplingStrategy struct {
	IDeletionSamplingStrategy IDeletionSamplingStrategy
}

type SamplingStrategy struct {
	ISamplingStrategy ISamplingStrategy
}

/*
DELETION GRAPH SAMPLING METHODS
*/
type DeletionRandomNodeSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomNodeNeighbourSampling struct{ IDeletionSamplingStrategy }
type DeletionInclusiveRandomNodeNeighbourSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomDegreeNodeSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomEdgeSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomNodeEdgeSampling struct{ IDeletionSamplingStrategy }
type DeletionHybridSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomWalkSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomWalkWithJumpSampling struct{ IDeletionSamplingStrategy }
type DeletionRandomWalkWithRestartSampling struct{ IDeletionSamplingStrategy }

func (strategy *DeletionRandomNodeSampling) SamplingStage(g *UndirectedGraph, howMany int) error {
	nodes := GetDictKeys(g.Nodes)
	for _, node := range rand.Perm(len(nodes))[:howMany] {
		g.RemoveNode(nodes[node])
	}
	return nil
}

func (strategy *DeletionRandomNodeNeighbourSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	for i := 0; i < howManyToDelete; i++ {
		nodes := GetDictKeys(g.Nodes)
		nodeIndex := rand.Perm(len(nodes))[0]
		neighbours := g.Edges[nodes[nodeIndex]]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		g.RemoveNode(neighbours[neighbourIndex])
	}
	return nil
}

func (strategy *DeletionInclusiveRandomNodeNeighbourSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	for i := 0; i < howManyToDelete; i++ {
		nodes := GetDictKeys(g.Nodes)
		nodeIndex := rand.Perm(len(nodes))[0]
		neighbours := g.Edges[nodes[nodeIndex]]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		g.RemoveNode(nodes[nodeIndex])
		g.RemoveNode(neighbours[neighbourIndex])
	}
	return nil
}

func (strategy *DeletionRandomDegreeNodeSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	for i := 0; i < howManyToDelete; i++ {
		var choices []weightedrand.Choice
		for node := range g.Nodes {
			choices = append(choices, weightedrand.NewChoice(node, uint(len(g.Edges[node]))))
		}
		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		nodeToRemove := pick.(Node)
		g.RemoveNode(nodeToRemove)
	}
	return nil
}

func (strategy *DeletionRandomEdgeSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	edges := g.GetEdgeTuples()

	for _, edgeIndex := range rand.Perm(len(edges))[:howManyToDelete] {
		g.RemoveEdge(edges[edgeIndex])
	}
	return nil
}

func (strategy *DeletionRandomNodeEdgeSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	edges := g.GetEdgeTuples()
	nodes := GetDictKeys(g.Nodes)

	for _, nodeIndex := range rand.Perm(len(nodes))[:howManyToDelete] {
		nodeEdges := g.Edges[nodes[nodeIndex]]
		edgeIndex := rand.Perm(len(nodeEdges))[0]
		g.RemoveEdge(edges[edgeIndex])
	}
	return nil
}

func (strategy *DeletionHybridSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	w := float32(0.42) // TODO update
	edges := g.GetEdgeTuples()

	nodes := GetDictKeys(g.Nodes)

	for i := 0; i < howManyToDelete; i++ {
		if rand.Float32() < w {
			for _, nodeIndex := range rand.Perm(len(nodes))[:1] {
				nodeEdges := g.Edges[nodes[nodeIndex]]
				edgeIndex := rand.Perm(len(nodeEdges))[0]
				g.RemoveEdge(edges[edgeIndex])
			}
		} else {
			edgeIndex := rand.Perm(len(edges))[0]
			g.RemoveEdge(edges[edgeIndex])
		}
	}
	return nil
}

func (strategy *DeletionRandomWalkSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	startNode := g.pickRandomNode()
	currentNode := startNode

	for i := 0; i < howManyToDelete; i++ {
		neighbors := g.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			g.RemoveNode(currentNode)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			currentNode = nextNode
		} else {
			// If the current node has no neighbors, break the walk
			break
		}
	}
	return nil
}

func (strategy *DeletionRandomWalkWithRestartSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	startNode := g.pickRandomNode()
	neighbors := g.Edges[startNode]
	nodeToInclude := neighbors[rand.Intn(len(neighbors))]

	for i := 0; i < howManyToDelete; i++ {
		neighbors := g.Edges[nodeToInclude]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			g.RemoveNode(nodeToInclude)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				neighbors = g.Edges[startNode]
				nodeToInclude = neighbors[rand.Intn(len(neighbors))]
			} else {
				nodeToInclude = nextNode
			}
		} else {
			// If the current node has no neighbors, go to first node
			neighbors = g.Edges[startNode]
			nodeToInclude = neighbors[rand.Intn(len(neighbors))]
		}
	}
	return nil
}

func (strategy *DeletionRandomWalkWithJumpSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) error {
	startNode := g.pickRandomNode()
	currentNode := startNode

	for i := 0; i < howManyToDelete; i++ {
		neighbors := g.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			g.RemoveNode(currentNode)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				currentNode = g.pickRandomNode()
			} else {
				currentNode = nextNode
			}
		} else {
			// If the current node has no neighbors, jump to random node
			currentNode = g.pickRandomNode()
		}
	}
	return nil
}

func (strategy *DeletionSamplingStrategy) Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (*UndirectedGraph, error) {
	ng := &UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return nil, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		err = strategy.IDeletionSamplingStrategy.SamplingStage(ng, int(0.03*float32(len(ng.Nodes))))
		if err != nil {
			return nil, fmt.Errorf("error performing sampling stage: %w", err)
		}

		// We retain the largest connected component and delete the rest
		connectedComponents := ConnectedComponents(ng)
		ng = connectedComponents.GetBiggestComponent()
	}
	return ng, nil
}

/*
	PRESERVATION GRAPH SAMPLING METHODS
*/

type PreservationRandomNodeSampling struct{ ISamplingStrategy }
type PreservationRandomNodeNeighbourSampling struct{ ISamplingStrategy }
type PreservationInclusiveRandomNodeNeighbourSampling struct{ ISamplingStrategy }
type PreservationRandomDegreeNodeSampling struct{ ISamplingStrategy }
type PreservationNodeSamplingWithContraction struct{ ISamplingStrategy }
type RandomPageRankNodeSampling struct{ ISamplingStrategy }
type PreservationRandomEdgeSampling struct{ ISamplingStrategy }
type PreservationRandomNodeEdgeSampling struct{ ISamplingStrategy }
type PreservationHybridSampling struct{ ISamplingStrategy }
type PreservationRandomWalkSampling struct{ ISamplingStrategy }
type PreservationRandomWalkWithJumpSampling struct{ ISamplingStrategy }
type PreservationRandomWalkWithRestartSampling struct{ ISamplingStrategy }
type PreservationTopKEdgeSampling struct{ ISamplingStrategy }
type PreservationTopRatioEdgeSampling struct{ ISamplingStrategy }

func (strategy *PreservationRandomNodeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(g.Nodes)
	var selectedNodes []Node

	for _, node := range rand.Perm(len(nodes))[:expectedFinalGraphSize] {
		ng.AddNode(nodes[node])
		selectedNodes = append(selectedNodes, nodes[node])
	}

	for i := 0; i < len(selectedNodes); i++ {
		neighbors := g.Edges[selectedNodes[i]]
		for _, neighbor := range neighbors {
			if ng.HasNode(neighbor) {
				ng.AddEdge(Edge{Node1: selectedNodes[i], Node2: neighbor})
			}
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomNodeNeighbourSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(graph.Nodes)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		nodeIndex := rand.Perm(len(nodes))[0]
		node := nodes[nodeIndex]
		ng.AddNode(node)
		neighbours := graph.Edges[node]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		neighbour := neighbours[neighbourIndex]
		ng.AddNode(neighbour)
		ng.AddEdge(Edge{Node1: node, Node2: neighbour})
	}
	return ng, nil
}

func (strategy *PreservationInclusiveRandomNodeNeighbourSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		nodes := GetDictKeys(graph.Nodes)
		nodeIndex := rand.Perm(len(nodes))[0]
		node := nodes[nodeIndex]
		neighbours := graph.Edges[node]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		neighbour := neighbours[neighbourIndex]
		ng_nodes := GetDictKeys(ng.Nodes)
		ng.AddNode(node)
		for _, ng_node := range ng_nodes {
			for _, sample_neighbour := range neighbours {
				if ng_node == sample_neighbour {
					ng.AddEdge(Edge{Node1: node, Node2: ng_node})
				}
			}
		}
		ng.AddNode(neighbour)
		ng.AddEdge(Edge{Node1: node, Node2: neighbour})
		neighbours_of_neighbour := graph.Edges[neighbour]
		for _, ng_node := range ng_nodes {
			for _, sample_neighbour := range neighbours_of_neighbour {
				if ng_node == sample_neighbour {
					ng.AddEdge(Edge{Node1: neighbour, Node2: ng_node})
				}
			}
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomDegreeNodeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)
	selectedNodes := map[Node]bool{}
	selectedNodesArray := []Node{}

	nodesCounter := 0
	for {
		var choices []weightedrand.Choice
		for node := range g.Nodes {
			choices = append(choices, weightedrand.NewChoice(node, uint(len(g.Edges[node]))))
		}

		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return UndirectedGraph{
				Nodes: nil,
				Edges: nil,
			}, fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		if !selectedNodes[pick.(Node)] {
			ng.AddNode(pick.(Node))
			selectedNodes[pick.(Node)] = true
			selectedNodesArray = append(selectedNodesArray, pick.(Node))
			nodesCounter = nodesCounter + 1
		}
		if expectedFinalGraphSize <= nodesCounter {
			break
		}
	}

	for i := 0; i < len(selectedNodesArray); i++ {
		neighbors := g.Edges[selectedNodesArray[i]]
		for _, neighbor := range neighbors {
			if ng.HasNode(neighbor) {
				ng.AddEdge(Edge{Node1: selectedNodesArray[i], Node2: neighbor})
			}
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomEdgeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)

	edges := g.GetEdgeTuples()
	nodes := map[Node]bool{}

	for _, edgeIndex := range rand.Perm(len(edges)) {
		ng.AddEdge(edges[edgeIndex])
		nodes[edges[edgeIndex].Node1] = true
		nodes[edges[edgeIndex].Node2] = true
		if expectedFinalGraphSize <= len(nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomNodeEdgeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)
	edges := g.GetEdgeTuples()

	nodes := GetDictKeys(g.Nodes)
	newNodes := map[Node]bool{}

	for _, nodeIndex := range rand.Perm(len(nodes)) {
		nodeEdges := g.Edges[nodes[nodeIndex]]
		for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
			ng.AddEdge(edges[edgeIndex])
			newNodes[edges[edgeIndex].Node1] = true
			newNodes[edges[edgeIndex].Node2] = true
			if expectedFinalGraphSize <= len(newNodes) {
				return ng, nil
			}
		}
	}
	return ng, nil
}

func (strategy *PreservationHybridSampling) Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (*UndirectedGraph, error) {
	w := float32(0.5)
	ng := &UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}

	edges := graph.GetEdgeTuples()
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(graph.Nodes)
	newNodes := map[Node]bool{}

	for {
		if rand.Float32() < w {
			for _, nodeIndex := range rand.Perm(len(nodes))[:1] {
				nodeEdges := graph.Edges[nodes[nodeIndex]]
				for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
					ng.AddEdge(edges[edgeIndex])
					newNodes[edges[edgeIndex].Node1] = true
					newNodes[edges[edgeIndex].Node2] = true
				}
			}
		} else {
			for _, edgeIndex := range rand.Perm(len(edges)) {
				ng.AddEdge(edges[edgeIndex])
				newNodes[edges[edgeIndex].Node1] = true
				newNodes[edges[edgeIndex].Node2] = true
			}
		}
		if expectedFinalGraphSize <= len(newNodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomWalkSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := graph.pickRandomNode()
	currentNode := startNode

	for {
		neighbors := graph.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.AddEdge(Edge{Node1: currentNode, Node2: nextNode})
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			currentNode = nextNode
		} else {
			// If the current node has no neighbors, break the walk
			break
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomWalkWithRestartSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := graph.pickRandomNode()
	neighbors := graph.Edges[startNode]
	for {
		if len(neighbors) > 0 {
			break
		} else {
			startNode = graph.pickRandomNode()
			neighbors = graph.Edges[startNode]
		}
	}
	//nodeToInclude := neighbors[rand.Intn(len(neighbors))]
	currentNode := startNode

	for {
		//neighbors := graph.Edges[nodeToInclude]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.AddEdge(Edge{Node1: currentNode, Node2: nextNode})
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				currentNode = startNode
			} else {
				currentNode = nextNode
			}
			neighbors = graph.Edges[currentNode]
		} else {
			// If the current node has no neighbors, go to first node
			currentNode = startNode
			neighbors = graph.Edges[currentNode]
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *PreservationRandomWalkWithJumpSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := graph.pickRandomNode()
	currentNode := startNode

	for {
		neighbors := graph.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.AddEdge(Edge{Node1: currentNode, Node2: nextNode})
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				currentNode = ng.pickRandomNode()
			} else {
				currentNode = nextNode
			}
		} else {
			// If the current node has no neighbors, jump to random node
			currentNode = graph.pickRandomNode()
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *PreservationTopRatioEdgeSampling) Sample(
	g UndirectedGraph,
	sampleNeighborRatio float32, // e.g., 0.3 keeps top 30% per node
) (UndirectedGraph, error) {

	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	for u := range g.Nodes {
		nbrs := g.Edges[u]
		if len(nbrs) == 0 {
			continue
		}

		// Build (neighbor, weight) where weight = neighbor's degree
		type we struct {
			v      Node
			weight float32
		}
		ws := make([]we, 0, len(nbrs))
		for _, v := range nbrs {
			if v == u { // skip self-loops if any
				continue
			}
			ws = append(ws, we{
				v:      v,
				weight: float32(g.NodeDegree(v)),
			})
		}
		if len(ws) == 0 {
			continue
		}

		// Sort by weight DESC (keep highest-degree neighbors)
		sort.SliceStable(ws, func(i, j int) bool { return ws[i].weight > ws[j].weight })

		// Top-K per node, where K = ceil(ratio * deg(u))
		k := int(math.Ceil(float64(sampleNeighborRatio) * float64(len(ws))))
		if k < 0 {
			k = 0
		}
		if k > len(ws) {
			k = len(ws)
		}
		if k == 0 {
			continue
		}

		// Add the top-K edges (once for undirected graphs)
		for i := 0; i < k; i++ {
			v := ws[i].v
			if !ng.HasEdge(u, v) {
				ng.AddEdge(Edge{Node1: u, Node2: v})
			}
		}
	}

	return ng, nil
}

func (strategy *PreservationTopKEdgeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	// TODO: where do we get the K parameter?
	// shall we change the interface to passing sampling parameters object?
	// and some Factory, to create the appropriate map of parameters
	topK := 5
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}

	for nodeId := range g.Nodes {
		neighbours := g.Edges[nodeId]
		if len(neighbours) == 0 {
			continue
		}
		weightedNeighbours := make([]WeightedElement, 0)
		for k := 0; k < len(neighbours); k++ {
			weightedNeighbours = append(weightedNeighbours, WeightedElement{
				Payload: neighbours[k],
				Weight:  float32(g.NodeDegree(neighbours[k])),
			})
		}
		sort.SliceStable(weightedNeighbours, func(i, j int) bool {
			return weightedNeighbours[i].Weight > weightedNeighbours[j].Weight
		})

		topK_adjusted := min(topK, len(neighbours))

		for idx := 0; idx < topK_adjusted; idx++ {
			nodeId_1 := nodeId
			nodeId_2 := weightedNeighbours[idx].Payload.(Node)
			if !ng.HasEdge(nodeId_1, nodeId_2) {
				ng.AddEdge(Edge{
					Node1: nodeId_1,
					Node2: nodeId_2,
				})
			}
		}
	}
	return ng, nil
}

/*
	CONTRACTION GRAPH SAMPLING METHODS
*/

type ContractionRandomNodeSampling struct{ ISamplingStrategy }
type ContractionRandomNodeNeighbourSampling struct{ ISamplingStrategy }
type ContractionInclusiveRandomNodeNeighbourSampling struct{ ISamplingStrategy }
type ContractionRandomDegreeNodeSampling struct{ ISamplingStrategy }
type ContractionRandomEdgeSampling struct{ ISamplingStrategy }
type ContractionRandomNodeEdgeSampling struct{ ISamplingStrategy }
type ContractionHybridSampling struct{ ISamplingStrategy }
type ContractionRandomWalkSampling struct{ ISamplingStrategy }
type ContractionRandomWalkWithRestartSampling struct{ ISamplingStrategy }
type ContractionRandomWalkWithJumpSampling struct{ ISamplingStrategy }

type ContractionPageRankNodeSampling struct{ ISamplingStrategy }

func (strategy *ContractionRandomNodeSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		nodes := GetDictKeys(graph.Nodes)
		for _, node := range rand.Perm(len(nodes))[:1] {
			graph.ContractNode(nodes[node])
		}
	}
	return ng, nil
}

func (strategy *ContractionRandomNodeNeighbourSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		nodes := GetDictKeys(graph.Nodes)
		nodeIndex := rand.Perm(len(nodes))[0]
		neighbours := ng.Edges[Node(nodes[nodeIndex])]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		ng.ContractNode(neighbours[neighbourIndex])
	}
	return ng, nil
}

func (strategy *ContractionInclusiveRandomNodeNeighbourSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		nodes := GetDictKeys(graph.Nodes)
		nodeIndex := rand.Perm(len(nodes))[0]
		neighbours := ng.Edges[Node(nodes[nodeIndex])]
		neighbourIndex := rand.Perm(len(neighbours))[0]
		ng.ContractNode(nodes[nodeIndex])
		ng.ContractNode(neighbours[neighbourIndex])
	}
	return ng, nil
}

func (strategy *ContractionRandomDegreeNodeSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		var choices []weightedrand.Choice
		for node := range ng.Nodes {
			choices = append(choices, weightedrand.NewChoice(node, uint(len(ng.Edges[node]))))
		}
		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return ng, fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		nodeToContract := pick.(Node)
		graph.ContractNode(nodeToContract)
	}
	return ng, nil
}

func (strategy *ContractionRandomEdgeSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	edges := graph.GetEdgeTuples()
	nodes := map[Node]bool{}

	for _, edgeIndex := range rand.Perm(len(edges)) {
		ng.ContractEdge(edges[edgeIndex])
		if expectedFinalGraphSize <= len(nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *ContractionRandomNodeEdgeSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}

	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)
	edges := graph.GetEdgeTuples()

	nodes := GetDictKeys(graph.Nodes)
	newNodes := map[Node]bool{}

	for _, nodeIndex := range rand.Perm(len(nodes)) {
		nodeEdges := graph.Edges[nodes[nodeIndex]]
		for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
			ng.ContractEdge(edges[edgeIndex])
			if expectedFinalGraphSize <= len(newNodes) {
				break
			}
		}
	}
	return ng, nil
}

func (strategy *ContractionHybridSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}
	w := float32(0.5)

	edges := graph.GetEdgeTuples()
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(graph.Nodes)
	newNodes := map[Node]bool{}

	for {
		if rand.Float32() < w {
			for _, nodeIndex := range rand.Perm(len(nodes))[:1] {
				nodeEdges := graph.Edges[nodes[nodeIndex]]
				for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
					ng.ContractEdge(edges[edgeIndex])
				}
			}
		} else {
			for _, edgeIndex := range rand.Perm(len(edges)) {
				ng.ContractEdge(edges[edgeIndex])
			}
		}
		if expectedFinalGraphSize <= len(newNodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *ContractionRandomWalkSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := ng.pickRandomNode()
	currentNode := startNode

	for {
		neighbors := graph.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.ContractNode(currentNode)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			currentNode = nextNode
		} else {
			// If the current node has no neighbors, break the walk
			break
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *ContractionRandomWalkWithRestartSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := ng.pickRandomNode()
	neighbors := graph.Edges[startNode]
	contractNode := neighbors[rand.Intn(len(neighbors))]

	for {
		neighbors := graph.Edges[contractNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.ContractNode(contractNode)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				neighbors = graph.Edges[startNode]
				contractNode = neighbors[rand.Intn(len(neighbors))]
			} else {
				contractNode = nextNode
			}
		} else {
			// If the current node has no neighbors, go to first node
			neighbors = graph.Edges[startNode]
			contractNode = neighbors[rand.Intn(len(neighbors))]
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

func (strategy *ContractionRandomWalkWithJumpSampling) Sample(graph UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{}
	// deep copy
	err := copier.Copy(ng, graph)
	if err != nil {
		return ng, fmt.Errorf("error performing deep copy: %w", err)
	}
	expectedFinalGraphSize := int(float32(len(graph.Nodes)) * sampledGraphSizeRatio)

	startNode := ng.pickRandomNode()
	currentNode := startNode

	for {
		neighbors := graph.Edges[currentNode]
		if len(neighbors) > 0 {
			nextNode := neighbors[rand.Intn(len(neighbors))]
			ng.ContractNode(currentNode)
			// c value taken from Leskovec, Jure, and Christos Faloutsos. "Sampling from large graphs." Proceedings of the 12th ACM SIGKDD international conference on Knowledge discovery and data mining. 2006.
			if rand.Float32() < 0.15 {
				currentNode = ng.pickRandomNode()
			} else {
				currentNode = nextNode
			}

		} else {
			// If the current node has no neighbors, jump to random node
			currentNode = ng.pickRandomNode()
		}
		if expectedFinalGraphSize <= len(ng.Nodes) {
			break
		}
	}
	return ng, nil
}

// Helper method to pick a random node from the graph
func (g *UndirectedGraph) pickRandomNode() Node {
	var nodes []Node
	for node := range g.Nodes {
		nodes = append(nodes, node)
	}
	return nodes[rand.Intn(len(nodes))]
}
