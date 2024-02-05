package io

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/jmCodeCraft/go-network/model"
)

type IGraphFormatReader interface {
	Read(reader io.Reader) model.UndirectedGraph
	ReadFromFile(filename string) model.UndirectedGraph
	AddNodesToGraph(g *model.UndirectedGraph, nodes []model.Node)
}

type AdjacencyListReader struct{ IGraphFormatReader } // DONE
type EdgeListReader struct{ IGraphFormatReader }      // DONE

func (strategy *IGraphFormatReader) Read(reader io.Reader) model.UndirectedGraph {
	ng := model.UndirectedGraph{}

	csvReader := csv.NewReader(reader)
	lineCount := 0
	for {
		read, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil
		}
		slog.Info(fmt.Sprintf("read: %+v", read))
		nodes := lineToList(read)
		IGraphFormatReader.AddNodesToGraph(ng, nodes)
		lineCount++
	}
	return ng
}

func (strategy *IGraphFormatReader) ReadFromFile(filename string) model.UndirectedGraph {
	readFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	ng := strategy.Read(readFile)

	err = readFile.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	return ng
}

func (a *AdjacencyListReader) AddNodesToGraph(g *model.UndirectedGraph, nodes []model.Node) {
	g.AddEdgesFromIntEdgeList(model.Node(nodes[0]), nodes[1:])
}

func (a *EdgeListReader) AddNodesToGraph(g *model.UndirectedGraph, nodes []model.Node) {
	g.AddEdge(model.Edge{nodes[0], nodes[1]})
}

/*
lineToList converts a slice of strings representing numerical values into a slice of model.Node integers.

Parameters:
- values: A slice of strings containing numerical values to be converted.

Returns:
- integers: A slice of model.Node integers representing the converted values.

Panics:
- If any string in the input slice cannot be converted to an integer, a panic with the corresponding error is triggered.
*/
func lineToList(values []string) (integers []model.Node) {
	integers = make([]model.Node, len(values))
	for index, value := range values {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		integers[index] = model.Node(valueInt)
	}
	return
}
