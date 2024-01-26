package algorithm

import (
	"fmt"
	"math/rand"

	"github.com/jmCodeCraft/go-network/model",
	"github.com/mroth/weightedrand"
)

type SamplingStrategy struct {
	Graph model.Graph
	Strategy string
}

func (strategy *SamplingStrategy) Sample(g *model.Graph, samplingRate float32) (model.Graph) {
	switch strategy.Strategy {
		case "RN":
			return strategy.randomNodeSampling(g, samplingRate)
		case "RDN":
			return strategy.randomDegreeNodeSampling(g, samplingRate)
		default:
			return strategy.randomNodeSampling(g, samplingRate)
	}
}

func (strategy *SamplingStrategy) randomNodeSampling(g *model.Graph, samplingRate float32) (ng model.Graph) {
	ng := model.Graph()
	sampleSize := int(float32(len(g.Nodes)) * samplingRate)
	for _, node := range rand.Perm(sampleSize) {
		ng.AddNode(node)
	}
	for node, _ := range ng.Nodes {
		for edge := range g.Edges[node] {
			ng.AddEdge(node, edge)
		}
	}
	return ng
}

func (strategy *SamplingStrategy) randomDegreeNodeSampling(g *model.Graph, samplingRate float32) (ng model.Graph) {
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