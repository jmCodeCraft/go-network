package model

import (
	"fmt"
	"math/rand"

	"github.com/mroth/weightedrand"
)

type IDeletionSamplingStrategy interface {
	SamplingStage(graph *UndirectedGraph, howManyToDelete int) error
	Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error)
}

type ISamplingStrategy interface {
	Sample(graph *UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error)
}

type DeletionSamplingStrategy struct {
	IDeletionSamplingStrategy IDeletionSamplingStrategy
}

type SamplingStrategy struct {
	ISamplingStrategy ISamplingStrategy
}

// supported sampling methods so far
type DeletionRandomNodeSampling struct{ IDeletionSamplingStrategy }       // DONE
type DeletionRandomDegreeNodeSampling struct{ IDeletionSamplingStrategy } // DONE
type DeletionRandomEdgeSampling struct{ IDeletionSamplingStrategy }       // DONE
type DeletionRandomNodeEdgeSampling struct{ IDeletionSamplingStrategy }   //DONE
type DeletionHybridSampling struct{ IDeletionSamplingStrategy }

func (strategy *DeletionRandomNodeSampling) SamplingStage(g *UndirectedGraph, howMany int) error {
	nodes := GetDictKeys(g.Nodes)
	for _, node := range rand.Perm(len(nodes))[:howMany] {
		g.RemoveNode(nodes[node])
	}
	return nil
}

func (strategy *DeletionRandomDegreeNodeSampling) SamplingStage(g *UndirectedGraph, howMany int) error {
	for i := 0; i < howMany; i++ {
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

func (strategy *DeletionRandomEdgeSampling) SamplingStage(g *UndirectedGraph, howMany int) error {
	edges := g.GetEdgeTuples()
	sampleSize := len(edges) - howMany

	for _, edgeIndex := range rand.Perm(len(edges))[:sampleSize] {
		g.RemoveEdge(edges[edgeIndex])
	}
	return nil
}

func (strategy *DeletionRandomNodeEdgeSampling) SamplingStage(g *UndirectedGraph, howMany int) error {
	edges := g.GetEdgeTuples()
	sampleSize := len(edges) - howMany
	nodes := GetDictKeys(g.Nodes)

	for _, nodeIndex := range rand.Perm(len(nodes))[:sampleSize] {
		nodeEdges := g.Edges[nodes[nodeIndex]]
		for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
			g.RemoveEdge(edges[edgeIndex])
		}
	}
	return nil
}

func (strategy *DeletionHybridSampling) SamplingStage(g *UndirectedGraph, howMany int, w float32) error {
	edges := g.GetEdgeTuples()
	sampleSize := len(edges) - howMany
	nodes := GetDictKeys(g.Nodes)

	for i := 0; i < howMany; i++ {
		if rand.Float32() < w {
			for _, nodeIndex := range rand.Perm(len(nodes))[:1] {
				nodeEdges := g.Edges[nodes[nodeIndex]]
				for _, edgeIndex := range rand.Perm(len(nodeEdges))[:1] {
					g.RemoveEdge(edges[edgeIndex])
				}
			}
		} else {
			for _, edgeIndex := range rand.Perm(len(edges))[:sampleSize] {
				g.RemoveEdge(edges[edgeIndex])
			}
		}
	}
	return nil
}

func (strategy *DeletionSamplingStrategy) Sample(g *UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{} //TODO deep copy
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)

	for len(ng.Nodes) <= expectedFinalGraphSize {
		_ = *strategy.IDeletionSamplingStrategy.SamplingStage(&ng, int(0.03*float32(len(ng.Nodes))))

		// We retain the largest connected component and delete the rest
		components := ConnectedComponents(ng)
		biggestComponentArray := components.ComponentsArray[components.BiggestComponentIdx]
		biggestComponentDict := components.ComponentsDict[components.BiggestComponentIdx]
		for i := 0; i < len(g.Nodes); i++ {
			if len(biggestComponentArray) > 0 {
				for node := range ng.Nodes {
					if !biggestComponentDict[node] {
						ng.RemoveNode(node)
					}
				}
			}
		}
	}
	return ng, nil
}

type RandomNodeSampling struct{ ISamplingStrategy }       // DONE
type RandomDegreeNodeSampling struct{ ISamplingStrategy } // DONE
type NodeSamplingWithContraction struct{ ISamplingStrategy }

// type RandomPageRankNodeSampling struct{ ISamplingStrategy }

type RandomEdgeSampling struct{ ISamplingStrategy }     // DONE
type RandomNodeEdgeSampling struct{ ISamplingStrategy } //DONE
type HybridSampling struct{ ISamplingStrategy }         //DONE
// type InducedRandomEdgeSampling struct{ ISamplingStrategy }
type EdgeSamplingWithContraction struct{ ISamplingStrategy }

type RandomWalkSampling struct{ ISamplingStrategy }
type RandomJumpSampling struct{ ISamplingStrategy }
type SnowballSampling struct{ ISamplingStrategy }
type ForestFireSampling struct{ ISamplingStrategy }
type FrontierSampling struct{ ISamplingStrategy }
type ExpansionSampling struct{ ISamplingStrategy }

func (strategy *RandomNodeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	ng := UndirectedGraph{
		Nodes: make(map[Node]bool),
		Edges: make(map[Node][]Node),
	}
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(g.Nodes)
	selectedNodes := []Node{}

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

func (strategy *RandomDegreeNodeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
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

func (strategy *RandomEdgeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
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

func (strategy *RandomNodeEdgeSampling) Sample(g UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
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
				break
			}
		}
	}
	return ng, nil
}

func (strategy *HybridSampling) Sample(g *UndirectedGraph, sampledGraphSizeRatio float32) (UndirectedGraph, error) {
	w := float32(0.5)
	ng := UndirectedGraph{
		Nodes: map[Node]bool{},
		Edges: map[Node][]Node{},
	}

	edges := g.GetEdgeTuples()
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * sampledGraphSizeRatio)
	nodes := GetDictKeys(g.Nodes)
	newNodes := map[Node]bool{}

	for {
		if rand.Float32() < w {
			for _, nodeIndex := range rand.Perm(len(nodes))[:1] {
				nodeEdges := g.Edges[nodes[nodeIndex]]
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
