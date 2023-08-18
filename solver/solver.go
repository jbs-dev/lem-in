package solver

import (
	"container/list"
)

// Node represents a structure for a node in a graph.
// It contains data related to the node's name, position, connectivity, and traversal status.
type Node struct {
	Name      string  // Unique identifier for the node
	X         int     // X-coordinate of the node's position
	Y         int     // Y-coordinate of the node's position
	Capacity  int     // Maximum capacity of the node (if applicable)
	Used      bool    // Indicates if the node has been used (for some operation or traversal)
	Visited   bool    // Indicates if the node has been visited during a traversal
	Neighbors []*Node // List of neighboring nodes to this node
}

// BFS performs Breadth-First Search traversal from a start node to an end node in a graph.
// It returns the path from the start node to the end node if such a path exists.
func BFS(start, end *Node, n map[string]*Node) []string {
	q := list.New()                                            // Create a new queue to maintain nodes for BFS traversal
	q.PushBack(start)                                          // Enqueue the start node
	start.Visited = true                                       // Mark the start node as visited
	var mapOfLinks map[string]string = make(map[string]string) // Map to store links (child -> parent) during traversal

	// Traverse all nodes in BFS manner
	for q.Len() > 0 {
		e := q.Front() // Get the front node from the queue
		node := e.Value.(*Node)
		for _, v := range node.Neighbors {
			if !v.Visited {
				v.Visited = true               // Mark the neighbor as visited
				mapOfLinks[v.Name] = node.Name // Store the link (neighbor -> current node)
				q.PushBack(v)                  // Enqueue the neighbor for BFS traversal
			}
		}
		q.Remove(e) // Dequeue the current node
	}

	// Extract the path from start to end using the constructed links
	path := extractPath(mapOfLinks, n, start, end)
	return path
}

// ClearVisited resets the Visited flag for nodes in the graph which haven't been used.
func ClearVisited(n map[string]*Node) {
	for _, v := range n {
		if !v.Used {
			v.Visited = false
		}
	}
}

// extractPath backtracks from the end node to the start node using the given map of links.
// It constructs the path taken during the BFS traversal.
func extractPath(l map[string]string, n map[string]*Node, start, end *Node) []string {
	name := end.Name // Start backtracking from the end node
	var result []string
	result = append(result, name) // Add the end node to the result
	ok := false

	// Backtrack from end to start using the links
	for name != start.Name {
		name, ok = l[name] // Get the parent node for the current node
		if !ok {
			// If a link doesn't exist for a node, then return nil (path doesn't exist)
			return nil
		}
		result = append(result, name) // Add the node to the result
	}
	return result
}
