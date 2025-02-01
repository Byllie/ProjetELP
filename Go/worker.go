package main

import (
	"bufio"
	"fmt"
	"net"
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
	Edges     map[int]*Vertex
	community *Community
	CC        float32
}

func (graph *Graph) AddEdge(srcKey, destKey int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{
			Edges:     make(map[int]*Vertex),
			community: nil,
			CC:        -1,
		}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{
			Edges:     make(map[int]*Vertex),
			community: nil,
			CC:        -1,
		}
	}

	graph.Vertices[srcKey].Edges[destKey] = graph.Vertices[destKey]
	graph.Vertices[destKey].Edges[srcKey] = graph.Vertices[srcKey]
}

func NewGraphFromTCP(conn net.Conn) *Graph {
	graph := &Graph{Vertices: make(map[int]*Vertex), Communities: []*Community{}}
	scanner := bufio.NewScanner(conn)
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

	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur de lecture du graph:", err)
	}

	return graph
}

func worker(id int, jobs <-chan int, results chan<- bool, graph *Graph) {
	for key := range jobs {
		neighbors := []int{}
		for neighborKey := range graph.Vertices[key].Edges {
			neighbors = append(neighbors, neighborKey)
		}
		fmt.Printf("Le worker %d travaille sur le noeud: %d avec les voisins %v\n", id, key, neighbors)
		results <- true
	}
}

func createJobs(graph *Graph, jobs chan<- int) {
	for key := range graph.Vertices {
		jobs <- key
	}
	close(jobs)
}

func processGraph(graph *Graph, numWorkers int) {
	jobs := make(chan int, 2*numWorkers)
	results := make(chan bool, 2*numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, graph)
	}

	go createJobs(graph, jobs)

	for a := 1; a <= len(graph.Vertices); a++ {
		<-results
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	graph := NewGraphFromTCP(conn)
	const numWorkers = 8
	processGraph(graph, numWorkers)
	fmt.Println("Fin Graph")
}

func main() {
	listener, err := net.Listen("tcp", ":5827")
	if err != nil {
		fmt.Println("Erreur dans le démarrage du serveur:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 5827...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur de connection tcp:", err)
			continue
		}
		go handleConnection(conn)
	}
}
