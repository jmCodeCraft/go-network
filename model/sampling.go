package model

import (
	"fmt"
	"math/rand"

	"github.com/mroth/weightedrand"
)

type ISamplingStrategy interface {
	SamplingStage(graph *UndirectedGraph, howManyToDelete int) *UndirectedGraph
	Sample(graph *UndirectedGraph, ratioNodesToDelete float32) UndirectedGraph
}

type SamplingStrategy struct {
	ISamplingStrategy ISamplingStrategy
}

// supported sampling methods so far
type RandomNodeSampling struct{ SamplingStrategy }       // DONE
type RandomDegreeNodeSampling struct{ SamplingStrategy } // WIP
type RandomPageRankNodeSampling struct{ SamplingStrategy }
type NodeSamplingWithContraction struct{ SamplingStrategy }
type RandomEdgeSampling struct{ SamplingStrategy } // DONE
type RandomNodeEdgeSampling struct{ SamplingStrategy }
type HybridSampling struct{ SamplingStrategy }
type InducedRandomEdgeSampling struct{ SamplingStrategy }
type EdgeSamplingWithContraction struct{ SamplingStrategy }
type RandomWalkSampling struct{ SamplingStrategy }
type RandomJumpSampling struct{ SamplingStrategy }
type SnowballSampling struct{ SamplingStrategy }
type ForestFireSampling struct{ SamplingStrategy }
type FrontierSampling struct{ SamplingStrategy }
type ExpansionSampling struct{ SamplingStrategy }

func (strategy *RandomNodeSampling) SamplingStage(g *UndirectedGraph, howManyToDelete int) *UndirectedGraph {
	for _, node := range rand.Perm(howManyToDelete) {
		g.RemoveNode(Node(node))
	}
	return g
}

func (strategy *RandomDegreeNodeSampling) SamplingStage(g *UndirectedGraph, howMany int) (*UndirectedGraph, error) {
	for i := 0; i < howMany; i++ {
		var choices []weightedrand.Choice
		for node := range g.Nodes {
			choices = append(choices, weightedrand.NewChoice(node, uint(len(g.Edges[node]))))
		}
		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return nil, fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		nodeToRemove := pick.(Node)
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

func (strategy *SamplingStrategy) Sample(g UndirectedGraph, ratioNodesToDelete float32) UndirectedGraph {
	expectedFinalGraphSize := int(float32(len(g.Nodes)) * (1 - ratioNodesToDelete))

	for len(g.Nodes) <= expectedFinalGraphSize {
		g = *strategy.ISamplingStrategy.SamplingStage(&g, int(0.03*float32(len(g.Nodes))))

		// We retain the largest connected component and delete the rest
		components := ConnectedComponents(g)
		biggestComponentArray := components.ComponentsArray[components.BiggestComponentIdx]
		biggestComponentDict := components.ComponentsDict[components.BiggestComponentIdx]
		for i := 0; i < len(g.Nodes); i++ {
			if len(biggestComponentArray) > 0 {
				for node := range g.Nodes {
					if !biggestComponentDict[node] {
						g.RemoveNode(node)
					}
				}
			}
		}
	}
	return g
}
