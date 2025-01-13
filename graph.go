package main

type Graph struct {
	Vertices map[int]*Vertex
}

type Vertex struct {
	Edges map[int]*Vertex
}

func (graph *Graph) AddEdge(srcKey, destKey int, weight int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{map[int]*Vertex{}}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{map[int]*Vertex{}}
	}

	graph.Vertices[srcKey].Edges[destKey] = graph.Vertices[destKey]
	graph.Vertices[destKey].Edges[srcKey] = graph.Vertices[srcKey]
}
