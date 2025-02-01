package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func worker(id int, jobs <-chan rune, results chan<- bool) {
	for char := range jobs {
		fmt.Printf("Le worker %d travaille sur le charactère: %c\n", id, char)
		time.Sleep(6 * time.Second)
		results <- true
	}
}

func createJobs(input string, jobs chan<- rune) {
	for _, char := range input {
		jobs <- char
	}
	close(jobs)
}

func processString(input string, numWorkers int) {
	jobs := make(chan rune, 2*numWorkers)
	results := make(chan bool, 2*numWorkers)

	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	go createJobs(input, jobs)

	for a := 1; a <= len(input); a++ {
		<-results
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur de lecture:", err)
		return
	}
	input = strings.TrimSpace(input)
	fmt.Printf("Recu: %s\n", input)
	const numWorkers = 3
	processString(input, numWorkers)
	fmt.Println("Fini:", input)
}

func main() {
	listener, err := net.Listen("tcp", ":5827")
	if err != nil {
		fmt.Println("Erreur dans le démarage du serveur:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 5827...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
