package main

import "fmt"

func worker(id int, graph *Graph, jobs <-chan int, results chan<- int, done chan struct{}) {
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				close(done)
				return
			}

			results <- j * 2
		case <-done:
			return
		}
	}
}

func create_worker() {
	const numWorkers = 5
	jobs := make(chan int, numWorkers*2)    // Numéro du node dans le graph
	results := make(chan int, numWorkers*2) //-1 remove; 0 no action; x transfer to comunity x
	done := make(chan struct{})
	graph := &Graph{
		Vertices: map[int]*Vertex{
			1: {Edges: map[int]*Vertex{2: {}, 3: {}}},
			2: {Edges: map[int]*Vertex{1: {}, 3: {}}},
			3: {Edges: map[int]*Vertex{1: {}, 2: {}}},
		},
	}

	for i := 1; i <= numWorkers; i++ {
		go worker(i, graph, jobs, results, done)
	}

	for node := range graph.Vertices {
		jobs <- node
	}
	close(jobs)

	for a := 1; a <= 9; a++ {
		result := <-results
		fmt.Println("Result:", result)
	}

	<-done
}
