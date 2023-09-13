package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	s "lem-in/solver"
)

func main() {
	// Collect command line arguments
	args := os.Args[1:]
	if len(args) != 1 {
		return
	}
	filename := args[0]

	// Read number of ants and map lines from file
	numberOfAnts, lines := readLines(filename)
	if numberOfAnts < 1 {
		var err error = fmt.Errorf("number of ants is invalid")
		abortOnError(err)
	}
	// Note the start time for performance measurement
	start := time.Now()

	// Parse the map lines to extract node details
	mapOfNodes, startNode, endNode := parseLines(lines)

	// Print simple aesthetic header using printf
	// fmt.Printf("----------------------- %s.txt -----------------------\n\n", filename)

	// Print file content for debugging
	bytes, err := os.ReadFile(fmt.Sprintf("./examples/%s.txt", filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(string(bytes))
	fmt.Println()
	fmt.Println() // Add this line to ensure a newline after printing the map content.

	// Find all paths from start node to end node
	var takenPath []string

	allPaths := Solver(takenPath, mapOfNodes, startNode, endNode)
	fmt.Println("All paths found:", allPaths)

	paths := allocatePaths(allPaths, startNode, mapOfNodes)
	fmt.Println("Allocated paths:", paths)

	for i, path := range paths {
		paths[i] = path[1:]
	}

	lemin(numberOfAnts, paths, mapOfNodes)

	// Find optimal paths for ants to travel
	/* 	antQueues := lemin(numberOfAnts, paths, mapOfNodes)
	   	usedPaths := make([][]string, 0)
	   	for i, queue := range antQueues {
	   		if len(queue) > 0 {
	   			usedPaths = append(usedPaths, paths[i])
	   		}
	   	}
	   	fmt.Println("Used paths:", usedPaths) */

	// Print how long the program took to run
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
}

// allocatePaths determines the optimal set of paths among all found paths.
// It aims to allocate as many non-overlapping paths as possible, considering the overall length.
func allocatePaths(paths [][]string, startNode *s.Node, mapOfNodes map[string]*s.Node) [][]string {
	// Initialize a 3D slice to hold potential path combinations.
	var maxNPaths [][][]string = make([][][]string, len(paths))

	// Generate combinations of paths that don't overlap.
	for i := 0; i < len(paths); i++ {
		maxNPaths[i] = append(maxNPaths[i], paths[i])
		for j := i + 1; j < len(paths); j++ {
			// Check if paths do not overlap and then append them.
			if norm(maxNPaths[i], paths[j]) {
				maxNPaths[i] = append(maxNPaths[i], paths[j])
			}
		}
	}

	// Filter the combinations to get the ones with maximum length (most paths).
	var result [][][]string
	max := 0
	for _, v := range maxNPaths {
		if max < len(v) {
			max = len(v)
		}
	}
	for _, v := range maxNPaths {
		if max == len(v) {
			result = append(result, v)
		}
	}

	// Among the selected combinations, find the one with the smallest total path length.
	var res [][]string
	min := int(^uint(0) >> 1)
	for _, v := range result {
		tempmin := 0
		for _, vv := range v {
			tempmin += len(vv)
		}
		if tempmin < min {
			min = tempmin
			res = v
		}
	}

	return res
}

// norm checks if path n2 intersects with any paths in n1.
// It does so by checking if nodes (excluding start and end nodes) overlap.
func norm(n1 [][]string, n2 []string) bool {
	// Iterate through paths in n1.
	for _, v := range n1 {
		// For each path in n1, check its nodes against nodes in n2.
		for _, k := range v[1 : len(v)-1] {
			for _, kk := range n2[1 : len(n2)-1] {
				// If there's an overlap, return false.
				if k == kk {
					return false
				}
			}
		}
	}
	// If there's no overlap found, return true.
	return true
}

// Solver finds all paths from startNode to endNode.
// It uses a recursive approach, marking nodes as used to prevent cycles.
// It combines BFS and backtracking to find all possible paths.
func Solver(takenPath []string, mapOfNodes map[string]*s.Node, startNode, endNode *s.Node) [][]string {
	// If the start node and end node are the same, return nil.
	if startNode == endNode {
		return nil
	}

	// Use BFS to find a path from startNode to endNode.
	extractedPath := reverse(s.BFS(startNode, endNode, mapOfNodes))
	// fmt.Println("Extracted path:", extractedPath)

	// Construct the full path by appending the taken path and the extracted path.
	var fullpath []string
	if extractedPath != nil {
		fullpath = append(fullpath, takenPath...)
		fullpath = append(fullpath, extractedPath...)
	}
	// Mark the current start node as used and add its name to the taken path.
	takenPath = append(takenPath, startNode.Name)
	startNode.Used = true

	// Reset visited markers for all nodes.
	s.ResetVisited(mapOfNodes)

	// Recurse into neighbors to find additional paths.
	var result [][]string
	if fullpath != nil {
		result = append(result, fullpath)
	}
	for _, v := range startNode.Neighbors {
		if !v.Used {
			paths := Solver(takenPath, mapOfNodes, v, endNode)
			for _, path := range paths {
				if !isEqual(fullpath, path) {
					result = append(result, path)
				}
			}
		}
	}

	// Reset the used and visited markers for the current start node.
	startNode.Used = false
	startNode.Visited = false

	return result
}

// isEqual compares two slices of strings and determines if they are the same.
func isEqual(fullpath, path []string) bool {
	if len(fullpath) != len(path) {
		return false
	}
	for i := range fullpath {
		if fullpath[i] != path[i] {
			return false
		}
	}
	return true
}

// minLen finds the minimum length after summing the lengths of the paths and their respective ant queues.
func minLen(p [][]string, ants [][]int) int {
	min := int(^uint(0) >> 1)

	for i := range p {
		if min > len(p[i])+len(ants[i]) {
			min = len(p[i]) + len(ants[i])
		}
	}

	return min
}

// MaxLen finds the maximum length after summing the lengths of the paths and their respective ant queues.
func MaxLen(p [][]string, ants [][]int) int {
	min := 0
	for i := range p {
		if min < len(p[i])+len(ants[i]) {
			min = len(p[i]) + len(ants[i])
		}
	}
	return min
}

type AntMovement struct {
	AntNumber int
	Movement  string
}

func lemin(n int, p [][]string, m map[string]*s.Node) [][]int {
	// Initialize a list of queues to assign ants to paths.
	var antQueues [][]int = make([][]int, len(p))
	i := 1
	min := minLen(p, antQueues)

	// Distribute ants across the paths based on minimum length criteria.
	for i <= n {
		for k := 0; k < len(p); k++ {
			if len(p[k])+len(antQueues[k]) <= min {
				antQueues[k] = append(antQueues[k], i)
				min = minLen(p, antQueues)
				break
			}
		}
		i++
	}

	// Construct a movement solution for the ants.
	var solution [][]AntMovement = make([][]AntMovement, MaxLen(p, antQueues)-1)
	for i := 0; i < len(p); i++ {
		for j, v := range antQueues[i] {
			for k, w := range p[i] {
				movement := AntMovement{
					AntNumber: v,
					Movement:  w,
				}
				solution[k+j] = append(solution[k+j], movement)
			}
		}
	}

	for _, v := range solution {
		sort.Slice(v, func(i, j int) bool {
			if v[i].AntNumber == v[j].AntNumber {
				return v[i].Movement < v[j].Movement
			}
			return v[i].AntNumber < v[j].AntNumber
		})
		movements := make([]string, len(v))
		for idx, movement := range v {
			movements[idx] = fmt.Sprintf("L%d-%s", movement.AntNumber, movement.Movement)
		}
		fmt.Println(strings.Join(movements, " "))
	}

	return antQueues
}

// reverse reverses the order of a slice of strings.
func reverse(nodes []string) []string {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}

// abortOnError checks for an error and if found, it prints an error message and exits the program.
func abortOnError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(0)
	}
}

// buildNodes constructs a map of nodes from a slice of node data.
func buildNodes(nodes [][]string) map[string]*s.Node {
	var n map[string]*s.Node = make(map[string]*s.Node)
	for _, v := range nodes {
		x, err := strconv.Atoi(v[1])
		abortOnError(err)
		y, err := strconv.Atoi(v[2])
		abortOnError(err)
		if _, ok := n[v[0]]; ok {
			var err error = fmt.Errorf("duplicate rooms")
			abortOnError(err)
		}
		n[v[0]] = &s.Node{Name: v[0], X: x, Y: y}
	}
	return n
}

// parseLines processes lines of input data and returns node information, including start and end nodes.
func parseLines(lines []string) (map[string]*s.Node, *s.Node, *s.Node) {
	// Variable initializations
	var start bool
	var end bool
	var startNode string
	var endNode string
	var nodes [][]string
	var links [][]string

	// Parse each line to extract node or link information
	for _, v := range lines {
		if !strings.HasPrefix(v, "#") {
			splitted := strings.Split(v, " ")
			if len(splitted) == 3 {
				nodes = append(nodes, splitted)
				if start {
					// Handle start flag and node information
					if startNode != "" {
						fmt.Printf("ERROR: invalid start data format %s", v)
						os.Exit(0)
					}
					startNode = splitted[0]
					start = false
				}
				if end {
					// Handle end flag and node information
					if endNode != "" {
						fmt.Println("ERROR: invalid finish data format")
						os.Exit(0)
					}
					endNode = splitted[0]
					end = false
				}
			} else {
				splitted = strings.Split(v, "-")
				links = append(links, splitted)
			}
		}

		// Mark start or end when they are found.
		if v == "##start" {
			start = true
		}

		if v == "##end" {
			end = true
		}
	}

	// Handle missing start or end node case.
	if startNode == "" || endNode == "" {
		var err error = fmt.Errorf("no start or end")
		abortOnError(err)
	}

	n := buildNodes(nodes)
	createLinks(n, links)

	return n, n[startNode], n[endNode]
}

// createLinks updates node information by establishing links between nodes.
func createLinks(n map[string]*s.Node, links [][]string) {
	for _, link := range links {
		node1, ok := n[link[0]]
		if !ok {
			var err error = fmt.Errorf("unknown room")
			abortOnError(err)
		}
		node2, ok := n[link[1]]
		if !ok {
			var err error = fmt.Errorf("unknown room")
			abortOnError(err)
		}
		if node1 == node2 {
			var err error = fmt.Errorf("self linkage")
			abortOnError(err)
		}
		node1.Neighbors = append(node1.Neighbors, node2)
		node2.Neighbors = append(node2.Neighbors, node1)
	}
}

// readLines reads the content of a given file and returns the number of nodes and a slice of lines.
func readLines(file string) (int, []string) {
	bytes, err := os.ReadFile(fmt.Sprintf("./examples/%s.txt", file))
	if err != nil {
		fmt.Printf("error reading file: %s", err)
	}
	lines := strings.Split(string(bytes), "\n")
	n, err := strconv.Atoi(lines[0])
	abortOnError(err)
	return n, lines[1:]
}

func printNode(n map[string]*s.Node) {
	for _, v := range n {
		fmt.Printf("Name: %s, X: %d, Y: %d, Used: %t, Visited: %t, Neighbors: ", v.Name, v.X, v.Y, v.Used, v.Visited)
		for _, vv := range v.Neighbors {
			fmt.Printf("%s ", vv.Name)
		}
		fmt.Println()
	}
}
