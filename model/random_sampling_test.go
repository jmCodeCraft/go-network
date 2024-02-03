package model

import (
	"fmt"
	"math/rand"
	"testing"
)

// BenchmarkFastGNPRandomGraph benchmarks the performance of FastGNPRandomGraph function
// by generating a graph with 10,000 nodes and a specified edge creation probability.
func BenchmarkFastGNPRandomGraph(b *testing.B) {
	// Set up the benchmark parameters
	numberOfNodes := 10000
	probabilityForEdgeCreation := 0.1 // You can adjust the probability as needed

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Reset the random seed for each iteration to ensure consistent results
		rand.NewSource(int64(i))

		// Run the FastGNPRandomGraph function
		g := FastGNPRandomGraph(numberOfNodes, probabilityForEdgeCreation)
		fmt.Println(g)
	}
}
