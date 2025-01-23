package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
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
	index     int
	Edges     map[int]*Vertex
	community *Community
	CC        float32
}

func (graph Graph) getCommunity(node int) *Community {
	// TODO : Peut être optimiser en stockant les communautés pour chaque noeud (parce que la ça prends BEAUCOUP de temps)
	c := graph.Vertices[node].community
	return c
}

func (graph *Graph) AddEdge(srcKey, destKey int) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		graph.Vertices[srcKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), nil, -1}
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		graph.Vertices[destKey] = &Vertex{len(graph.Vertices) + 1, make(map[int]*Vertex), nil, -1}
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

func (graph *Graph) CountTrianglesVerticesCommunity(node int, c Community) int {
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
	for _, value := range neighborsFormTriangles {
		value = false
		if value {
			fmt.Println(value)
		}

	}
	neighbors := graph.Vertices[node].Edges
	for neighbor1 := range neighbors {
		for neighbor2 := range neighbors {
			if graph.Vertices[neighbor1].Edges[neighbor2] != nil && neighbor1 < neighbor2 {
				neighborsFormTriangles[neighbor1] = true
				neighborsFormTriangles[neighbor2] = true
			}
		}
	}

	/* if len(neighbors) != len(neighborsFormTriangles) {
		fmt.Println("Error in vt")
	} */
	/* fmt.Println(neighborsFormTriangles)
	fmt.Println(len(neighbors), len(neighborsFormTriangles)) */
	count := 0
	for _, formTriangle := range neighborsFormTriangles {
		if formTriangle {
			count++
		}
	}

	return count
}

func (graph *Graph) vtExcludingC(node int, c Community) int {
	/* Correspond to vt(x,V\C) in https://dl-acm-org.docelec.insa-lyon.fr/doi/pdf/10.1145/2566486.2568010 */
	if _, ok := graph.Vertices[node]; !ok {
		return 0
	}

	neighborsFormTriangles := make(map[int]bool)
	neighbors := graph.Vertices[node].Edges
	for neighbor := range neighbors {
		if c.Vertices[neighbor] == nil {
			for neighbor2 := range graph.Vertices[neighbor].Edges {
				if neighbor2 != node && neighbors[neighbor2] != nil && neighbor2 < neighbor {
					neighborsFormTriangles[neighbor] = true

				}
			}
		}
	}
	count := 0
	for _, formTriangle := range neighborsFormTriangles {
		if formTriangle {
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

func (graph *Graph) WccNode(node int, c Community) float64 {
	triangleInGraph := graph.CountTrianglesVertices(node)
	if triangleInGraph == 0 {
		return 0
	}

	triangleInC := graph.CountTrianglesVerticesCommunity(node, c)
	vtxV := graph.vt(node)
	vtxVexC := graph.vtExcludingC(node, c)
	if triangleInGraph == 0 || (float64(vtxVexC)+float64(len(c.Vertices)-1)) == 0 {
		return 0
	}
	res := float64(triangleInC) / float64(triangleInGraph) * float64(vtxV) / (float64(vtxVexC) + float64(len(c.Vertices)-1))
	/* 	if res > 1 || res < 0 {
		fmt.Println("Resultat incohérent : ", res)
		fmt.Println("Node : ", node)
		fmt.Println("triangleInC : ", triang	vtxV := graph.vt(node)
leInC)
		fmt.Println("triangleInGraph : ", triangleInGraph)
		fmt.Println("vtxV : ", vtxV)
		fmt.Println("vtxVexC : ", vtxVexC)
		fmt.Println("len(c.Vertices) : ", len(c.Vertices))
	} */

	return res
}

func (graph *Graph) WccCommunity(c Community) float64 {
	if len(c.Vertices) == 0 {
		return 0
	}

	avg := float64(0)
	for key := range c.Vertices {
		avg += graph.WccNode(key, c)
	}
	// if (avg / float64(len(c.Vertices))) > 1 {
	// 	fmt.Println("avg : ", avg/float64(len(c.Vertices)))
	// }
	return avg / float64(len(c.Vertices))
}

func (graph *Graph) Wcc() float64 {
	avg := float64(0)
	for _, community := range graph.Communities {
		avg += graph.WccCommunity(*community) * float64(len(community.Vertices))
	}
	return avg / float64(len(graph.Vertices))
}
func (graph *Graph) WccI(node int, c Community) float64 {
	// TODO : enlever les commentaires
	// https://preview.redd.it/9my9pzmf2s771.png?auto=webp&s=5ff7c47100f9ac3edb284f9612ec3fe934bf311a
	// maybe it works
	V := len(graph.Vertices)
	var newC Community
	newC.Vertices = c.Vertices
	newC.Vertices[node] = graph.Vertices[node]
	return 1 / float64(V) * (float64(len(newC.Vertices))*graph.WccCommunity(newC) - graph.WccCommunity(c)*float64(len(c.Vertices)))
}

func (graph *Graph) WccR(node int, c Community) float64 {
	// La meme chose que WccI mais en enlevant le noeud node de la communauté c
	V := len(graph.Vertices)
	var newC Community
	newC.Vertices = c.Vertices
	delete(newC.Vertices, node)
	return 1 / float64(V) * (float64(len(newC.Vertices))*graph.WccCommunity(newC) - graph.WccCommunity(c)*float64(len(c.Vertices)))

}
func (graph *Graph) WccT(node int, source Community, dest Community) float64 {
	// TODO : enlever les commentaires
	// https://preview.redd.it/9my9pzmf2s771.png?auto=webp&s=5ff7c47100f9ac3edb284f9612ec3fe934bf311a
	// maybe it works
	return graph.WccI(node, dest) + graph.WccR(node, source)
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
	for k, v := range graph.Vertices {
		v.index = k
	}
	return graph
}

func (graph *Graph) bestMovement(node int) (int, *Community) {
	movement := NO_ACTION
	sourceC := graph.getCommunity(node)
	var wccR float64
	if sourceC == nil {
		wccR = 0
	} else {
		wccR = graph.WccR(node, *sourceC)
	}	vtxV := graph.vt(node)

	wccT := 0.0
	var bestC *Community
	for dest := range graph.Vertices[node].Edges {
		destC := graph.getCommunity(dest)
		if destC != nil && destC != sourceC {
			if sourceC != nil {
				temp := graph.WccI(node, *destC) + graph.WccR(node, *sourceC)
				if temp > wccT {
					wccT = temp
					bestC = destC
				}
			} else {
				temp := graph.WccI(node, *destC)
				if temp > wccT {
					wccT = temp
					bestC = destC
				}
			}
		}
	}
	if wccT > wccR {
		movement = REMOVE
	} else if wccT < wccR {
		movement = MOVE
	}
	return movement, bestC
}

const (
	REMOVE    = -1
	NO_ACTION = 0
	MOVE      = 1
)

func main() {

	f, err := os.Create("myprogram.ezview")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	// Your application code here

	//filePath := "com-dblp.ungraph.txt"
	filePath := "com-amazon.ungraph.txt"
	//filePath := "test_graph.txt"
	//filePath := "test_graph copy.txt"
	graph := NewGraphFromFile(filePath)
	precision := -0.000000000001
	max_index := 0
	for key, vertex := range graph.Vertices {
		vertex.CC = graph.ClusteringCoeficient(key)
		vertex.index = key
		if key > max_index {
			max_index = key
		}

	}

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
				graph.Vertices[dest].community = community
				VisitedVertices[dest] = true
			}
		}

		graph.Communities = append(graph.Communities, community)
	}
	/* 	fmt.Println("Number of communities : ", len(graph.Communities))
	   	fmt.Println("Wcc : ", graph.Wcc())
	   	CommunityWithAllVertices := &Community{make(map[int]*Vertex)}
	   	for _, vertex := range graph.Vertices {
	   		CommunityWithAllVertices.Vertices[vertex.index] = vertex
	   	}
	   	fmt.Println("Wcc community 1 : ", graph.WccCommunity(*graph.Communities[1]))

	   	fmt.Println("Wcc with 1 community : ")
	   	fmt.Println("WCC community : ", graph.WccCommunity(*CommunityWithAllVertices)) */
	/* for _, community := range graph.Communities {
		fmt.Println("Community : ")
		for key := range community.Vertices {
			fmt.Print(key, " ")
		}
		fmt.Println()
	} */
	/* 	for k, v := range graph.Vertices {
		if v.index != k {
			panic("AAAAAAAAAAAA")
		}
		c := graph.getCommunity(v.index)
		if c != nil {
			if c.Vertices[v.index] == nil {
				fmt.Println("Error ")

			}
		} else {
			fmt.Println(v)
		}
	} */

	// for key, vertex := range graph.Vertices {
	// 	for dest := range vertex.Edges {
	// 		fmt.Println(key, "<-->", dest)
	// 	}
	// }
	/* for key, vertex := range graph.Vertices {
		for dest := range vertex.Edges {
			fmt.Println(key, "<-->", dest)
		}
	} */

	// ############################################################################################################
	// Main loop
	// ############################################################################################################
	WCC := graph.Wcc()
	var newWCC float64
	var newGraph Graph
	newGraph.Vertices = graph.Vertices
	for math.Abs(newWCC-WCC) > precision {
		panic("Testing")
		var listMovement = make([]int, max_index+1)
		fmt.Println("New iteration")
		fmt.Println("WCC : ", WCC)
		newWCC = WCC
		newGraph.Communities = graph.Communities
		i := 0
		pourcentage := 0
		listDest := make(map[int]*Community)
		for key, _ := range graph.Vertices {
			if i*100/len(graph.Vertices) > pourcentage {
				pourcentage += 1
				fmt.Println(pourcentage, "%")
			}
			var c *Community
			listMovement[key], c = graph.bestMovement(key)
			if listMovement[key] == MOVE {
				listDest[key] = c
			}
			i++
		}

		fmt.Println("Applying movements")
		for key, movement := range listMovement {

			if movement == REMOVE {
				// Remove node from community
				community := newGraph.getCommunity(key)
				if community != nil {
					delete(community.Vertices, key)
				}
			} else if movement != NO_ACTION {
				// Move node to another community
				sourceC := newGraph.getCommunity(key)
				destC := listDest[key]
				if sourceC != nil {
					sourceC.Vertices[key] = newGraph.Vertices[key]
				}
				if destC != nil {
					destC.Vertices[key] = newGraph.Vertices[key]
				}
				graph.Vertices[key].community = destC
			}
		}

	}

}
