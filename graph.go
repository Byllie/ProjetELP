package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Graph struct {
	Vertices    map[int]*Vertex
	Communities []*Community
}

type Community struct {
	Vertices map[int]*Vertex
}

type Vertex struct {
	index int
	Edges map[int]*Vertex
	CC    float32
}

func (graph *Graph) AddEdge(srcKey, destKey int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), -1}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), -1}
	}

	graph.Vertices[srcKey].Edges[destKey] = graph.Vertices[destKey]
	graph.Vertices[destKey].Edges[srcKey] = graph.Vertices[srcKey]
}

func (graph *Graph) CountTrianglesVertices(node int) int {
	if _, ok := graph.Vertices[node]; !ok {
		return 0
	}

	count := 0
	neighbors := graph.Vertices[node].Edges

	for neighbor1 := range neighbors {
		for neighbor2 := range neighbors {
			if neighbor1 < neighbor2 && graph.Vertices[neighbor1].Edges[neighbor2] != nil {
				count++
			}
		}
	}
	return count
}

func (graph *Graph) CountTrianglesEdge(src, dest int) int {
	if _, ok := graph.Vertices[src]; !ok {
		return 0
	}
	if _, ok := graph.Vertices[dest]; !ok {
		return 0
	}

	count := 0
	for neighbor := range graph.Vertices[src].Edges {
		if neighbor != dest && graph.Vertices[dest].Edges[neighbor] != nil {
			count++
		}
	}
	return count
}

func (graph *Graph) RemoveEdgesWithoutTriangles() {
	for src, vertex := range graph.Vertices {
		for dest := range vertex.Edges {
			if src < dest {
				// TODO : Optimiser ça (ne pas compter tout les triangles mais break des qu'il y en a 1)
				if graph.CountTrianglesEdge(src, dest) == 0 {
					delete(graph.Vertices[src].Edges, dest)
					delete(graph.Vertices[dest].Edges, src)
				}
				if graph.Vertices[src].Edges == nil {
					delete(graph.Vertices, src)
				}
				if graph.Vertices[dest].Edges == nil {
					delete(graph.Vertices, dest)
				}

			}
		}
	}
}

func (graph *Graph) ClusteringCoeficient(node int) float32 {
	degree := len(graph.Vertices[node].Edges)
	if degree < 2 {
		return 0
	}

	triangles := graph.CountTrianglesVertices(node)
	return 2 * float32(triangles) / float32(degree*(degree-1))
}

func (graph *Graph) SortVerticesByCC() []*Vertex {
	vertices := make([]*Vertex, len(graph.Vertices))
	i := 0
	for _, vertex := range graph.Vertices {
		vertices[i] = vertex
		i++
	}
	sort.Slice(vertices, func(i, j int) bool {
		if vertices[i].CC == vertices[j].CC {
			return len(vertices[i].Edges) > len(vertices[j].Edges)
		}
		return vertices[i].CC > vertices[j].CC
	})
	return vertices
}

func (graph *Graph) WccCommunity(node int, c Community) float32 {
	return 0.1
}

func NewGraphFromFile(filePath string) *Graph {
	file, _ := os.Open(filePath)
	defer file.Close()

	graph := &Graph{Vertices: make(map[int]*Vertex), Communities: []*Community{}}
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
	filePath := "com-amazon.ungraph.txt"
	graph := NewGraphFromFile(filePath)
	graph.RemoveEdgesWithoutTriangles()

	// ######################################################################################################################
	// Tests with http://snap.stanford.edu/data/index.html#communities datasets
	// ######################################################################################################################
	/* sumTriangles := 0
	for key := range graph.Vertices {
		sumTriangles += graph.CountTrianglesVertices(key)
	}
	fmt.Println("Triangles : ", sumTriangles/3)
	Supposed to be 667129 for com-amazon.ungraph.txt and is 667129
	*/
	// ######################################################################################################################
	/* for key, vertex := range graph.Vertices {
		vertex.CC = graph.ClusteringCoeficient(key)
	}
	avgCC := float32(0)
	for _, vertex := range graph.Vertices {
		avgCC += vertex.CC
	}
	avgCC /= float32(len(graph.Vertices))
	fmt.Println("Average Clustering Coeficient : ", avgCC)
	*/
	/* Average Clustering Coeficient :  0.3967 for com-amazon.ungraph.txt and is 0.3967 if you don't remove edges without triangles
	 */
	// ######################################################################################################################

	SortedVerticesByCC := graph.SortVerticesByCC()

	VisitedVertices := make(map[int]bool)

	for _, vertex := range SortedVerticesByCC {
		if VisitedVertices[vertex.index] {
			continue
		}
		community := &Community{make(map[int]*Vertex)}
		community.Vertices[vertex.index] = vertex
		VisitedVertices[vertex.index] = true

		for dest := range vertex.Edges {
			if !VisitedVertices[dest] {
				community.Vertices[dest] = graph.Vertices[dest]
				VisitedVertices[dest] = true
			}
		}

		graph.Communities = append(graph.Communities, community)
	}

	fmt.Println("Communities : ", len(graph.Communities))

	// for _, community := range graph.Communities {
	// 	fmt.Println("Community : ")
	// 	for key := range community.Vertices {
	// 		fmt.Print(key, " ")
	// 	}
	// 	fmt.Println()
	// }

	// for key, vertex := range graph.Vertices {
	// 	for dest := range vertex.Edges {
	// 		fmt.Println(key, "<-->", dest)
	// 	}
	// }
}
