package solver

// ConvertNodesToGraph takes the start Node (which is the entry to your entire graph due to its connections) and returns a Graph suitable for visualization.
func ConvertNodesToGraph(startNode *Node) *Graph {
	visited := make(map[string]bool)
	graph := &Graph{
		PointByName: make(map[string]*Point),
		Children:    make(map[*Point][]*Point),
		Edges:       []*Edge{},
	}

	// Recursive function to traverse nodes and convert them
	var traverseAndConvert func(node *Node)
	traverseAndConvert = func(node *Node) {
		// If the node is already visited, return
		if visited[node.Name] {
			return
		}
		visited[node.Name] = true

		// Convert the node to a point
		point := &Point{
			Name: node.Name,
			X:    node.X,
			Y:    node.Y,
		}
		graph.Points = append(graph.Points, point)
		graph.PointByName[point.Name] = point

		// Convert neighbors to points and edges
		for _, neighbor := range node.Neighbors {
			neighborPoint := &Point{
				Name: neighbor.Name,
				X:    neighbor.X,
				Y:    neighbor.Y,
			}
			edge := &Edge{
				From: point,
				To:   neighborPoint,
			}
			graph.Edges = append(graph.Edges, edge)
			graph.Children[point] = append(graph.Children[point], neighborPoint)

			// Recursively traverse neighbors
			traverseAndConvert(neighbor)
		}
	}

	// Start the traversal and conversion from the starting node
	traverseAndConvert(startNode)

	return graph
}
