package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	Vertices map[int]*Vertex
}

type Vertex struct {
	Edges map[int]*Vertex
}

func (graph *Graph) AddEdge(srcKey, destKey int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{make(map[int]*Vertex)}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{make(map[int]*Vertex)}
	}

	graph.Vertices[srcKey].Edges[destKey] = graph.Vertices[destKey]
	graph.Vertices[destKey].Edges[srcKey] = graph.Vertices[srcKey]
}

func NewGraphFromFile(filePath string) *Graph {
	file, _ := os.Open(filePath)
	defer file.Close()

	graph := &Graph{Vertices: make(map[int]*Vertex)}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		srcKey, _ := strconv.Atoi(parts[0])

		destKey, _ := strconv.Atoi(parts[1])

		graph.AddEdge(srcKey, destKey)
	}
	return graph
}

func main() {
	filePath := "graph.txt"
	graph := NewGraphFromFile(filePath)
}
