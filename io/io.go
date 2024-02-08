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

type GraphFormatReader struct {
	IGraphFormatReader IGraphFormatReader
}

type AdjacencyListReader struct{ GraphFormatReader } // DONE
type EdgeListReader struct{ GraphFormatReader }      // DONE

func (strategy *GraphFormatReader) Read(reader io.Reader) (*model.UndirectedGraph, error) {
	ng := &model.UndirectedGraph{}

	csvReader := csv.NewReader(reader)
	lineCount := 0
	for {
		read, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error reading csv: %w", err)
		}
		slog.Info(fmt.Sprintf("read: %+v", read))
		nodes := lineToList(read)
		strategy.IGraphFormatReader.AddNodesToGraph(ng, nodes)
		lineCount++
	}
	return ng, nil
}

func (strategy *GraphFormatReader) ReadFromFile(filename string) (*model.UndirectedGraph, error) {
	readFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	ng, err := strategy.Read(readFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	err = readFile.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing file: %w", err)
	}
	return ng, nil
}

func (a *AdjacencyListReader) AddNodesToGraph(g *model.UndirectedGraph, nodes []model.Node) {
	g.AddEdgesFromIntEdgeList(nodes[0], nodes[1:])
}

func (a *EdgeListReader) AddNodesToGraph(g *model.UndirectedGraph, nodes []model.Node) {
	g.AddEdge(model.Edge{Node1: nodes[0], Node2: nodes[1]})
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
