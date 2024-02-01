package model

import (
	"fmt"
	"math/rand"

	"github.com/mroth/weightedrand"
)

type SamplingStrategy interface {
	Sample(graph *UndirectedGraph, samplingRate float32) UndirectedGraph
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

func (strategy *RandomNodeSampling) Sample(g *UndirectedGraph, samplingRate float32) (ng UndirectedGraph) {
	ng = UndirectedGraph{}
	sampleSize := int(float32(len(g.Nodes)) * samplingRate)
	for _, node := range rand.Perm(sampleSize) {
		ng.AddNode(Node{NodeId: node})
	}
	for node1 := range ng.Nodes {
		for node2 := range g.Edges[node1] {
			ng.AddEdge(Edge{
				Node1: Node{node1},
				Node2: Node{node2},
			})
		}
	}
	// todo verify if both edges exist in ng
	return ng
}

func (strategy *RandomDegreeNodeSampling) Sample(g *UndirectedGraph, samplingRate float32) (ng UndirectedGraph, err error) {
	ng = UndirectedGraph{}
	sampleSize := int(float32(len(g.Nodes)) * samplingRate)

	var choices []weightedrand.Choice
	for node := range g.Nodes {
		choices = append(choices, weightedrand.NewChoice(node, uint(len(g.Edges[node]))))
	}

	for i := 0; i < sampleSize; i++ {
		choice, err := weightedrand.NewChooser(choices...)
		if err != nil {
			return ng, fmt.Errorf("error gettint new chooser: %w", err)
		}
		pick := choice.Pick()
		// todo make sampling without replacement
		ng.AddNode(Node{pick.(int)})
	}
	return ng, nil
}

func (strategy *RandomEdgeSampling) Sample(g *UndirectedGraph, samplingRate float32) (ng UndirectedGraph) {
	ng, edges := UndirectedGraph{}, g.GetEdgeTuples()
	sampleSize := int(float32(len(edges)) * samplingRate)

	for _, edgeIndex := range rand.Perm(len(edges))[:sampleSize] {
		ng.AddEdge(edges[edgeIndex])
	}
	return ng
}
