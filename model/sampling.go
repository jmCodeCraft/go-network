package model

import (
	"fmt"
	"math/rand"

	"github.com/mroth/weightedrand"
)

type SamplingStrategy interface {
	SamplingStage(graph *UndirectedGraph, howManyToDelete int) *UndirectedGraph
	Sample(graph *UndirectedGraph, ratioNodesToDelete float32) UndirectedGraph
}

// supported sampling methods so far
type RandomNodeSampling struct{}       // DONE
type RandomDegreeNodeSampling struct{} // WIP
type RandomPageRankNodeSampling struct{}
type NodeSamplingWithContraction struct{}
type RandomEdgeSampling struct{} // DONE
type RandomNodeEdgeSampling struct{}
type HybridSampling struct{}
type InducedRandomEdgeSampling struct{}
type EdgeSamplingWithContraction struct{}
type RandomWalkSampling struct{}
type RandomJumpSampling struct{}
type SnowballSampling struct{}
type ForestFireSampling struct{}
type FrontierSampling struct{}
type ExpansionSampling struct{}

func (strategy *RandomNodeSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) (g *UndirectedGraph) {
	for _, node := range rand.Perm(howManyToDelete) {
		g.RemoveNode(Node(node))
	}
	return g
}

func (strategy *RandomDegreeNodeSampling) SamplingStage(g *UndirectedGraph, howMany int) (g *UndirectedGraph, err error) {
	for i := 0; i < howMany; i++ {
		var choices []weightedrand.Choice
		for node := range g.Nodes {
			choices = append(choices, weightedrand.NewChoice(node, uint(len(g.Edges[node]))))
		}
		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return ng, fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		nodeToRemove := Node(pick.(int))
		g.RemoveNode(nodeToRemove)
	}
	return g, nil
}

func (strategy *RandomEdgeSampling) SamplingStage(g UndirectedGraph, howMany int) (ng UndirectedGraph) {
	ng, edges := UndirectedGraph{}, g.GetEdgeTuples()
	sampleSize := len(edges) - howMany

	for _, edgeIndex := range rand.Perm(len(edges))[:sampleSize] {
		ng.RemoveEdge(edges[edgeIndex])
	}
	return ng
}

func (strategy *SamplingStrategy) Sample(g *UndirectedGraph, ratioNodesToDelete float32) (ng UndirectedGraph) {
	ng = UndirectedGraph{} //TODO: deep copy g
	expectedFinalGraphSize := int(float32(len(ng.Nodes)) * (1 - ratioNodesToDelete))

	for {
		ng := strategy.SamplingStage(ng, int(0.03*float32(len(ng.Nodes))))

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
		if len(ng.Nodes) <= expectedFinalGraphSize {
			break
		}
	}
}
