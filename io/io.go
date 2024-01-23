package io

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-network/model"
)

func FromAdjacencyList(filename string) {
	readFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	g := model.Graph{}

	lineCount := 0
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		g.AddEdgesFromEdgeList(lineCount, lineToList(fileScanner.Text()))
		lineCount++
	}
	readFile.Close()
}

func lineToList(line string) (integers []int) {
	integers = []int{}
	for _, i := range strings.Split(line, ",") {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		integers = append(integers, j)
	}
	return
}
