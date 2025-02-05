package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Graph struct {
	Vertices    map[int]*Vertex
	Communities []*Community
}

type Community struct {
	Vertices map[int]*Vertex
}

type Vertex struct {
	index              int
	Edges              map[int]*Vertex
	community          *Community
	CC                 float32
	Triangles          int
	NbVerticesTriangle int
}
type ResultBestMouvement struct {
	node      int
	movement  int
	community *Community
}

func (graph *Graph) removeCommunity(key int) {
	community := graph.Communities[key]
	if len(community.Vertices) != 0 {
		panic("Trying to remove a non empty community")
	}
	graph.Communities[key] = graph.Communities[len(graph.Communities)-1]
	graph.Communities = graph.Communities[:len(graph.Communities)-1]

}

func (graph *Graph) GetCommunity(node int) *Community {
	// TODO : Peut être optimiser en stockant les communautés pour chaque noeud (parce que la ça prends BEAUCOUP de temps)
	c := graph.Vertices[node].community
	if c == nil {
		WriteLog("No community for node "+strconv.Itoa(node), graph)
		panic("No community for node " + strconv.Itoa(node) + "see log.txt for more informations")

	}
	return c
}

func (graph *Graph) AddEdge(srcKey, destKey int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), nil, -1, 0, 0}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), nil, -1, 0, 0}
	}

	graph.Vertices[srcKey].Edges[destKey] = graph.Vertices[destKey]
	graph.Vertices[destKey].Edges[srcKey] = graph.Vertices[srcKey]
}

func (graph *Graph) CountTrianglesVertices(node int) int {
	if _, ok := graph.Vertices[node]; !ok {
		return 0
	}

	countForGraph := 0
	neighbors := graph.Vertices[node].Edges

	for neighbor1 := range neighbors {
		for neighbor2 := range neighbors {
			if neighbor1 < neighbor2 && graph.Vertices[neighbor1].Edges[neighbor2] != nil {
				countForGraph++
			}
		}
	}
	return countForGraph
}

func (graph *Graph) CountTrianglesVerticesCommunity(node int, c *Community) int {
	/* Renvoit le nombre de triangles dans la communauté c */
	if _, ok := c.Vertices[node]; !ok {
		return 0
	}

	countForC := 0

	neighbors := graph.Vertices[node].Edges

	for neighbor1 := range neighbors {
		for neighbor2 := range neighbors {
			if neighbor1 < neighbor2 && graph.Vertices[neighbor1].Edges[neighbor2] != nil {
				if c.Vertices[neighbor1] != nil && c.Vertices[neighbor2] != nil {
					countForC++
				}
			}
		}
	}
	return countForC

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

func (graph *Graph) vt(node int) int {
	/* vt(x,V) in https://dl-acm-org.docelec.insa-lyon.fr/doi/pdf/10.1145/2566486.2568010 */

	if _, ok := graph.Vertices[node]; !ok {
		return 0
	}

	neighborsFormTriangles := make(map[int]bool)
	neighbors := graph.Vertices[node].Edges
	n := len(neighbors)
	i := 0 // Utiliser pour sortir de la boucle for sans passer par l'implementation du range qui prends du temps
	for neighbor1 := range neighbors {
		j := 0 // Utiliser pour sortir de la boucle for sans passer par l'implementation du range qui prends du temps
		for neighbor2 := range neighbors {
			if graph.Vertices[neighbor1].Edges[neighbor2] != nil && neighbor1 < neighbor2 {
				neighborsFormTriangles[neighbor1] = true
				neighborsFormTriangles[neighbor2] = true
			}
			j++
			if j == n {
				break
			}
		}
		i++
		if i == n {
			break
		}
	}

	//fmt.Println(neighborsFormTriangles)
	//fmt.Println(len(neighbors), len(neighborsFormTriangles))
	count := 0
	for _, formTriangle := range neighborsFormTriangles {
		if formTriangle {
			count++
		}
	}

	return count
}

func (graph *Graph) vtExcludingC(node int, c *Community) int {
	/* Correspond to vt(x,V\C) in https://dl-acm-org.docelec.insa-lyon.fr/doi/pdf/10.1145/2566486.2568010 */
	if _, ok := graph.Vertices[node]; !ok {
		return 0
	}

	neighborsFormTriangles := make(map[int]bool)
	neighbors := graph.Vertices[node].Edges

	for neighbor1 := range neighbors {
		if c.Vertices[neighbor1] == nil {
			for neighbor2 := range neighbors {
				if graph.Vertices[neighbor2].Edges[neighbor1] != nil {
					neighborsFormTriangles[neighbor1] = true
				}
			}
		}
	}
	count := 0

	for k := range neighborsFormTriangles {
		if k >= 0 {

			count++
		}

	}
	neighborsFormTriangles = nil

	return count
}

func (graph *Graph) RemoveEdgesWithoutTriangles() {
	for src, vertex := range graph.Vertices {
		for dest := range vertex.Edges {
			if src < dest {
				triangle := false
				for neighbor := range graph.Vertices[src].Edges {
					if neighbor != dest && graph.Vertices[dest].Edges[neighbor] != nil {
						triangle = true
					}
				}
				if !triangle {
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
	for j := 0; j < len(graph.Vertices); {
		if graph.Vertices[i] != nil {
			vertices[j] = graph.Vertices[i]
			j++
		}
		i++
	}
	sort.SliceStable(vertices, func(i, j int) bool {
		if vertices[i].CC == vertices[j].CC {
			return len(vertices[i].Edges) > len(vertices[j].Edges)
		}
		return vertices[i].CC > vertices[j].CC
	})
	return vertices
}

func (graph *Graph) WccNode(node int, c *Community) float64 {
	triangleInGraph := graph.Vertices[node].Triangles
	if triangleInGraph == 0 {
		return 0
	}
	triangleInC := graph.CountTrianglesVerticesCommunity(node, c)
	vtxV := graph.Vertices[node].NbVerticesTriangle
	vtxVexC := graph.vtExcludingC(node, c)
	if triangleInGraph == 0 || (float64(vtxVexC)+float64(len(c.Vertices)-1)) == 0 {
		return 0
	}
	res := float64(triangleInC) / float64(triangleInGraph) * float64(vtxV) / (float64(vtxVexC) + float64(len(c.Vertices)-1))
	if res > 1 || res < 0 {
		str := "Resultat incohérent de WccNode : " + strconv.FormatFloat(res, 'f', 6, 64) + "\n"
		str += "\nNode : "
		str += strconv.Itoa(node)
		str += "\ntriangleInC : "
		str += strconv.Itoa(triangleInC)
		str += "\ntriangleInGraph : "
		str += strconv.Itoa(triangleInGraph)
		str += "\nvtxV : "
		str += strconv.Itoa(vtxV)
		str += "\nvtxVexC : "
		str += strconv.Itoa(vtxVexC)
		str += "\nlen(c.Vertices) : "
		str += strconv.Itoa(len(c.Vertices))
		str += "\n\n"
		fmt.Println(str)
		//WriteLog(str, graph)
		//panic("Resultat incohérent de WccNode : " + strconv.FormatFloat(res, 'f', 6, 64))
	}

	return res
}

func (graph *Graph) WccCommunity(c *Community) float64 {
	if len(c.Vertices) == 0 {
		return 0
	}

	avg := float64(0)
	for key := range c.Vertices {
		avg += graph.WccNode(key, c)
	}
	if (avg / float64(len(c.Vertices))) > 1 {
		WriteLog("avg : "+strconv.FormatFloat(avg/float64(len(c.Vertices)), 'f', 6, 64), graph)
		panic("avg : " + strconv.FormatFloat(avg/float64(len(c.Vertices)), 'f', 6, 64))
	}
	return avg / float64(len(c.Vertices))
}

func (graph *Graph) Wcc() float64 {
	avg := float64(0)
	for _, v := range graph.Vertices {
		//v.ActualWCC = graph.WccNode(v.index, v.community)
		avg += graph.WccNode(v.index, v.community)
	}
	return avg / float64(len(graph.Vertices))
}
func (graph *Graph) WccI(node int, c *Community) float64 {
	// TODO : enlever les commentaires
	// https://preview.redd.it/9my9pzmf2s771.png?auto=webp&s=5ff7c47100f9ac3edb284f9612ec3fe934bf311a
	// maybe it works

	if len(c.Vertices) == 0 {
		return 0
	}
	newC := &Community{make(map[int]*Vertex)}
	for key, value := range c.Vertices {
		newC.Vertices[key] = value
	}
	newC.Vertices[node] = graph.Vertices[node]
	avgC := float64(0)
	avgNewC := float64(0)
	for key, _ := range c.Vertices {
		//avgC += v.ActualWCC
		avgC += graph.WccNode(key, c)
		avgNewC += graph.WccNode(key, newC)
	}

	if newC.Vertices[node] == nil {
		WriteLog("Error of community pointer in WccI. newC.Vertices[node]", graph)
		panic("Error of community pointer in WccI. See log.txt for more informations")
	}
	if c.Vertices[node] != nil {
		WriteLog("Error of community pointer in WccI. c.Vertices[node] != nil", graph)
		panic("Error of community pointer in WccI. See log.txt for more informations")
	}

	return (avgNewC - avgC + graph.WccNode(node, newC)) / float64(len(c.Vertices))

	/* V := len(graph.Vertices)
	newC := &Community{make(map[int]*Vertex)}
	for key, value := range c.Vertices {
		newC.Vertices[key] = value
	}

	newC.Vertices[node] = graph.Vertices[node]
	if newC.Vertices[node] == nil {
		WriteLog("Error of community pointer in WccI. ", graph)
		panic("Error of community pointer in WccI. See log.txt for more informations")
	}
	if c.Vertices[node] != nil {
		WriteLog("Error of community pointer in WccI. ", graph)
		panic("Error of community pointer in WccI. See log.txt for more informations")
	}
	return 1 / float64(V) * (float64(len(newC.Vertices))*graph.WccCommunity(newC) - graph.WccCommunity(c)*float64(len(c.Vertices))) */
}

func (graph *Graph) WccR(node int, c *Community) float64 {
	// La meme chose que WccI mais en enlevant le noeud node de la communauté c

	if len(c.Vertices) == 0 {
		return 0
	}
	newC := &Community{make(map[int]*Vertex)}
	for key, value := range c.Vertices {
		newC.Vertices[key] = value
	}
	delete(newC.Vertices, node)

	if c.Vertices[node] == nil {
		WriteLog("Error of community pointer in WccR. c.Vertices[node] == nil ", graph)
		panic("Error of community pointer in WccR. See log.txt for more informations")
	}
	if newC.Vertices[node] != nil {
		WriteLog("Error of community pointer in WccR. newC.Vertices[node] != nil ", graph)
		panic("Error of community pointer in WccR. See log.txt for more informations")
	}

	return -graph.WccI(node, newC)

	/* V := len(graph.Vertices)
	newC := &Community{make(map[int]*Vertex)}
	for key, value := range c.Vertices {
		newC.Vertices[key] = value
	}

	delete(newC.Vertices, node)
	if c.Vertices[node] == nil {
		WriteLog("Error of community pointer in WccR. c.Vertices[node] == nil ", graph)
		panic("Error of community pointer in WccR. See log.txt for more informations")
	}
	if newC.Vertices[node] != nil {
		WriteLog("Error of community pointer in WccR. newC.Vertices[node] != nil ", graph)
		panic("Error of community pointer in WccR. See log.txt for more informations")
	}
	return 1 / float64(V) * (float64(len(newC.Vertices))*graph.WccCommunity(newC) - graph.WccCommunity(c)*float64(len(c.Vertices))) */

}
func (graph *Graph) WccT(node int, source Community, dest Community) float64 {
	// TODO : enlever les commentaires
	// https://preview.redd.it/9my9pzmf2s771.png?auto=webp&s=5ff7c47100f9ac3edb284f9612ec3fe934bf311a
	// maybe it works
	return graph.WccI(node, &dest) + graph.WccR(node, &source)
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
		fmt.Println("Edge added")
	}
	fmt.Println("Graph created")
	for k, v := range graph.Vertices {
		v.index = k
	}
	return graph
}

func NewGraphFromFile(filePath string) *Graph {
	file, ok := os.Open(filePath)
	if ok != nil {
		WriteLog("Error opening file :"+filePath, nil)
		panic("Error opening file. See log.txt for more informations")
	}
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
	for k, v := range graph.Vertices {
		v.index = k
	}
	return graph
}

func WriteLog(err string, graph *Graph) {
	file, ok := os.Create("log.txt")
	if ok != nil {
		panic("Error opening file")
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	file2, ok := os.Create("memprof.ezview")
	if ok != nil {

		fmt.Println(err)
		return

	}
	pprof.WriteHeapProfile(file2)
	file.WriteString(err)
	if graph != nil {
		for key, c := range graph.Communities {
			if c != nil {
				file.WriteString("Community : " + strconv.Itoa(key) + " Vertices : " + strconv.Itoa(len(c.Vertices)) + "\n")
			}
		}
	} else {
		file.WriteString("Graph is nil")
	}

}

func WriteCommunityInFile(graph *Graph, filePath string) {
	file, ok := os.Create(filePath)
	if ok != nil {
		panic("Error opening file")
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	for key, c := range graph.Communities {
		if c != nil {
			file.WriteString("Community\t" + strconv.Itoa(key) + "\tVertices\t" + strconv.Itoa(len(c.Vertices)) + "\n")
			for key2 := range c.Vertices {
				file.WriteString("Node\t" + strconv.Itoa(key2) + "\n")
			}
		}
	}
}

func ParseCommunityFile(filePath string) []*Community {
	file, ok := os.Open(filePath)
	if ok != nil {
		panic("Error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var communities []*Community
	var c *Community
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if parts[0] == "Community" {
			if c != nil {
				communities = append(communities, c)
			}
			c = &Community{make(map[int]*Vertex)}
		} else if parts[0] == "Node" {
			node, _ := strconv.Atoi(parts[1])
			c.Vertices[node] = nil
		}
	}
	if c != nil {
		communities = append(communities, c)
	}
	return communities
}

func (graph *Graph) bestMovement(node int) (int, int, *Community) {
	movement := NO_ACTION
	sourceC := graph.GetCommunity(node)
	var wccR float64
	if sourceC == nil {
		wccR = 0
	} else {
		wccR = graph.WccR(node, sourceC)
	}
	wccT := 0.0
	var bestC *Community
	visitedCommunities := make(map[*Community]bool)
	for dest := range graph.Vertices[node].Edges {
		destC := graph.GetCommunity(dest)

		if destC != nil && destC != sourceC && !visitedCommunities[destC] {
			visitedCommunities[destC] = true
			if sourceC != nil {
				temp := graph.WccI(node, destC) + graph.WccR(node, sourceC)
				if temp > wccT {
					wccT = temp
					bestC = destC
				}
			} else {
				temp := graph.WccI(node, destC)
				if temp > wccT {
					wccT = temp
					bestC = destC
				}
			}
		}
	}
	if wccT > wccR && wccT > 1e-10 {
		movement = MOVE
	} else if wccT < wccR && wccR > 1e-10 {
		movement = REMOVE
	}
	/* fmt.Println("Movement : ", movement)
	if movement == MOVE || movement == REMOVE {
		if bestC != nil {
			fmt.Println("Node : ", node, " WccT : ", wccT, " WccR : ", wccR, " WccI : ", graph.WccI(node, bestC))
		} else {
			fmt.Println("Node : ", node, " WccT : ", wccT, " WccR : ", wccR)
		}
	}
	fmt.Println("Node : ", node, " movement : ", movement, " BestC : ", bestC) */
	return node, movement, bestC
}

func (graph *Graph) UpdateCommunities() {
	VisitedVertices := make(map[int]bool)

	for _, vertex := range graph.Vertices {
		if vertex.community == nil {
			community := &Community{make(map[int]*Vertex)}
			community.Vertices[vertex.index] = vertex
			vertex.community = community
			VisitedVertices[vertex.index] = true

			graph.Communities = append(graph.Communities, community)
		}

	}
	n := len(graph.Communities)
	for i := 0; i < n; i++ {
		community := graph.Communities[i]
		if len(community.Vertices) == 0 {
			graph.removeCommunity(i)
			n--
		}
	}
}

const (
	REMOVE    = -1
	NO_ACTION = 0
	MOVE      = 1
)

func worker(id int, jobs <-chan int, results chan<- ResultBestMouvement, graph *Graph) {
	for key := range jobs {
		//fmt.Printf("Le worker %d travaille sur le noeud: %d \n", id, key)
		n, mouvement, c := graph.bestMovement(key)
		results <- ResultBestMouvement{node: n, movement: mouvement, community: c}
	}
}
func createJobs(graph *Graph, jobs chan<- int) {
	for key := range graph.Vertices {
		jobs <- key
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection Started")
	f, err := os.Create("myprogram.ezview")
	if err != nil {

		fmt.Println(err)
		return

	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//filePath := "com-dblp.ungraph.txt"
	filePath := "com-amazon.ungraph.txt"
	//filePath := "com-youtube.ungraph.txt"
	//filePath := "test_graph.txt"
	//filePath := "test_graph copy.txt"
	//filePath := "test_graph5.txt"
	//graph := NewGraphFromTCP(conn)
	graph := NewGraphFromFile(filePath)
	precision := 0.0001
	max_index := 0
	for key, vertex := range graph.Vertices {
		vertex.CC = graph.ClusteringCoeficient(key)
		vertex.Triangles = graph.CountTrianglesVertices(key)
		vertex.NbVerticesTriangle = graph.vt(key)
		vertex.index = key
		if key > max_index {
			max_index = key
		}

	}

	graph.RemoveEdgesWithoutTriangles()
	fmt.Println("Preprocessing done")
	for key, vertex := range graph.Vertices {
		vertex.CC = graph.ClusteringCoeficient(key)
	}

	SortedVerticesByCC := graph.SortVerticesByCC()

	VisitedVertices := make(map[int]bool)

	for _, vertex := range SortedVerticesByCC {
		if VisitedVertices[vertex.index] {
			continue
		}
		community := &Community{make(map[int]*Vertex)}
		community.Vertices[vertex.index] = vertex
		vertex.community = community
		VisitedVertices[vertex.index] = true

		for dest := range vertex.Edges {
			if !VisitedVertices[dest] {
				community.Vertices[dest] = graph.Vertices[dest]
				graph.Vertices[dest].community = community
				VisitedVertices[dest] = true
			}
		}

		graph.Communities = append(graph.Communities, community)
	}

	// ############################################################################################################
	// Main loop
	// ############################################################################################################
	WCC := graph.Wcc()
	var newWCC float64
	listWCC := make([]float64, 0)
	listWCC = append(listWCC, WCC)
	newWCC = 0
	const numWorkers = 8
	jobs := make(chan int, 2*numWorkers)
	results := make(chan ResultBestMouvement, 2*numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, graph)
	}
	startTotalTime := time.Now()
	for math.Abs((WCC-newWCC)/newWCC) > precision {
		startTime := time.Now()
		var listMovement = make([]int, max_index+1)
		oldCommunities := make([]*Community, len(graph.Communities))

		for i, c := range graph.Communities {
			oldCommunities[i] = &Community{make(map[int]*Vertex)}
			for key, value := range c.Vertices {
				oldCommunities[i].Vertices[key] = value
			}
		}
		fmt.Println("New iteration")
		fmt.Println("WCC : ", WCC)
		fmt.Println("Nombre de communautés : ", len(graph.Communities))
		newWCC = WCC

		//pourcentage := 0
		listDest := make(map[int]*Community)
		go createJobs(graph, jobs)
		for i := 0; i < len(graph.Vertices); i++ {
			r := <-results
			c := r.community
			listMovement[r.node] = r.movement
			if listMovement[r.node] == MOVE {
				listDest[r.node] = c
			}
		}
		fmt.Println("Applying movements")

		for key, movement := range listMovement {

			if movement == REMOVE {
				// Remove node from communityenzor@enzo-83dd
				community := graph.GetCommunity(key)
				if community != nil {
					delete(community.Vertices, key)
				}
				community = &Community{make(map[int]*Vertex)}
				community.Vertices[key] = graph.Vertices[key]
				graph.Vertices[key].community = community
				graph.Communities = append(graph.Communities, community)

			} else if movement != NO_ACTION {
				// Move node to another community
				sourceC := graph.GetCommunity(key)
				destC := listDest[key]

				if sourceC != nil {
					delete(sourceC.Vertices, key)
				}
				if destC != nil {
					destC.Vertices[key] = graph.Vertices[key]
				}
				graph.Vertices[key].community = destC

			}
		}
		graph.UpdateCommunities()

		WCC = graph.Wcc()

		fmt.Println("WCC : ", WCC, " old WCC : ", newWCC)
		fmt.Println("Time to calculate WCC : ", time.Since(startTime))
		fmt.Println("Conditions : ", (WCC-newWCC)/newWCC, " Precision : ", precision)
		startTime = time.Now()
		var stabilise = false
		for _, indivWcc := range listWCC {
			if indivWcc == WCC || (WCC-newWCC) < 0 {

				stabilise = true
				break
			}
		}

		if stabilise {
			graph.Communities = oldCommunities
			for _, c := range graph.Communities {
				for key2 := range c.Vertices {
					graph.Vertices[key2].community = c
				}

			}
			fmt.Println("Stabilisation")
			for key := range graph.Vertices {
				n, mouvement, c := graph.bestMovement(key)
				if mouvement == MOVE {
					sourceC := graph.GetCommunity(n)
					destC := c
					if sourceC != nil {
						delete(sourceC.Vertices, n)
					}
					if destC != nil {
						destC.Vertices[n] = graph.Vertices[n]
					}
					graph.Vertices[n].community = destC
				} else if mouvement == REMOVE {
					community := graph.GetCommunity(n)
					if community != nil {
						delete(community.Vertices, n)
					}
					community = &Community{make(map[int]*Vertex)}
					community.Vertices[n] = graph.Vertices[n]
					graph.Vertices[n].community = community
					graph.Communities = append(graph.Communities, community)
				}
			}
			graph.UpdateCommunities()
			listWCC = append(listWCC, WCC)
			if len(listWCC) > 10 {
				listWCC = listWCC[1:]
			}
			WCC = graph.Wcc()
			fmt.Println("WCC : ", WCC)
			fmt.Println("Time to stabilise : ", time.Since(startTime))
		}

		listWCC = append(listWCC, WCC)
		if len(listWCC) > 10 {
			listWCC = listWCC[1:]
		}
		WriteCommunityInFile(graph, "communities.txt")
	}

	fmt.Println("Done in : ", time.Since(startTotalTime), " with ", numWorkers, " workers")
}

func main() {
	listener, err := net.Listen("tcp", ":5828")
	if err != nil {
		fmt.Println("Erreur dans le démarrage du serveur:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 5828...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur de connection tcp:", err)
			continue
		}
		go handleConnection(conn)
	}

}
