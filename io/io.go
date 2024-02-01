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

func FromAdjacencyListFile(filename string) error {
	readFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	FromAdjacencyList(readFile)

	err = readFile.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	return nil
}

// FromAdjacencyList reads an adjacency list from a reader.
// The reader should be a CSV file with the following format:
//
// 0,1,2,3
// 1,2
// 2,3
// 3
//
// Each line specifies...
func FromAdjacencyList(reader io.Reader) {
	g := model.UndirectedGraph{}

	csvReader := csv.NewReader(reader)
	lineCount := 0
	for {
		read, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}
		slog.Info(fmt.Sprintf("read: %+v", read))
		g.AddEdgesFromIntEdgeList(lineCount, lineToList(read))
		lineCount++
	}

	// todo remove this log
	slog.Info(fmt.Sprintf("graph: %+v", g))
}

func lineToList(values []string) (integers []int) {
	integers = make([]int, len(values))
	for index, value := range values {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		integers[index] = valueInt
	}
	return
}
