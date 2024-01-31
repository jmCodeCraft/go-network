package algorithm

import (
	"math/rand"
	"github.com/jmCodeCraft/go-network/generators"
	"github.com/jmCodeCraft/go-network/model"
),
	"github.com/mroth/weightedrand"
)

type SamplingStrategy interface {
	Sample(graph *model.Graph, samplingRate float32) model.Graph
}

// supported sampling methods so far
type RandomNodeSampling struct {} // DONE
type RandomDegreeNodeSampling struct {} //WIP
type RandomPageRankNodeSampling struct {}
type NodeSamplingWithContraction struct {}
type RandomEdgeSampling struct {} // DONE
type RandomNodeEdgeSampling struct {}
type HybridSampling struct {}
type InducedRandomEdgeSampling struct {}
type EdgeSamplingWithContraction struct {}
type RandomWalkSampling struct {}
type RandomJumpSampling struct {}
type SnowballSampling struct {}
type ForestFireSampling struct {}
type FrontierSampling struct {}
type ExpansionSampling struct {}

func (strategy *RandomNodeSampling) Sample(g *model.UndirectedGraph, samplingRate float32) (ng model.Graph) {
	ng := generators.EmptyGraph()
	sampleSize := int(float32(len(g.Nodes)) * samplingRate)
	for _, node := range rand.Perm(sampleSize) {
		ng.AddNode(model.Node{node})
	}
	for node, _ := range ng.Nodes {
		for edge := range g.Edges[node] {
			ng.AddEdge(model.Node{node}, edge)
		}
	}
	return ng
}

func (strategy *RandomDegreeNodeSampling) Sample(g *model.Graph, samplingRate float32) (ng model.Graph) {
	ng := model.Graph()
	sampleSize := int(float32(len(g.Nodes)) * samplingRate)

	choices := Choice[]
	for node, _ := range g.Nodes {
		choices = append(choices, NewChoice(node, len(g.Edges[node])))
	}
	for i := range [sampleSize]int{} {
		choice := NewChooser(choices).Pick()
		//https://pkg.go.dev/github.com/mroth/weightedrand#NewChooser
		choices = //TODO: delete the choice
	}
	return ng
}

// TODO: fix this function
func (strategy *RandomEdgeSampling) Sample(g *model.UndirectedGraph, samplingRate float32) (ng model.Graph) {
	ng, edges := generators.EmptyGraph(), g.GetEdgeTuples()
	sampleSize := int(float32(len(edges)) * samplingRate)
	for _, edge := range rand.Perm(edges) {
		ng.AddEdge(edge)
	}
	for node, _ := range ng.Nodes {
		for edge := range g.Edges[node] {
			ng.AddEdge(node, edge)
		}
	}
	return ng
}