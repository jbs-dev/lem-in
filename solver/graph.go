package solver

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	Name string
	X    int
	Y    int
}

// ParsePoint reads point in format `name x y` to Point
func ParsePoint(str string) (*Point, error) {
	NameXY := strings.Split(str, " ")
	if len(NameXY) != 3 {
		return nil, fmt.Errorf("wrong format, expected 'name x y'")
	}
	if strings.HasPrefix(NameXY[0], "L") {
		return nil, fmt.Errorf("wrong format, room name cannot start with 'L'")
	}
	X, err := parseCoordinate(NameXY[1])
	if err != nil {
		return nil, fmt.Errorf("wrong format, X should be integer number, not '%v'", NameXY[1])
	}
	Y, err := parseCoordinate(NameXY[2])
	if err != nil {
		return nil, fmt.Errorf("wrong format, Y should be integer number, not '%v'", NameXY[2])
	}
	return &Point{Name: NameXY[0], X: X, Y: Y}, nil
}

func parseCoordinate(coordStr string) (int, error) {
	coord, err := strconv.Atoi(coordStr)
	if err != nil {
		return 0, err
	}
	return coord, nil
}

type Edge struct {
	From *Point
	To   *Point
}

// ParseEdge reads edge in format `from-to` to Edge
func (g *Graph) ParseEdge(str string) (*Edge, error) {
	fromTo := strings.Split(str, "-")
	if len(fromTo) != 2 {
		return nil, fmt.Errorf("wrong format, expected 'from-to'")
	}
	from, to := fromTo[0], fromTo[1]
	fromPoint := g.PointByName[from]
	if fromPoint == nil {
		return nil, fmt.Errorf("point %v does not exist", from)
	}
	toPoint := g.PointByName[to]
	if toPoint == nil {
		return nil, fmt.Errorf("point %v does not exist", to)
	}
	return &Edge{From: fromPoint, To: toPoint}, nil
}

// ParseChildren creates a map Children to store all the points, connected with Point
func (g *Graph) ParseChildren() {
	g.Children = make(map[*Point][]*Point)
	for _, edge := range g.Edges {
		g.Children[edge.From] = append(g.Children[edge.From], edge.To)
		g.Children[edge.To] = append(g.Children[edge.To], edge.From)
	}
}

type Path []*Point

type Graph struct {
	AntNum      int
	Start       *Point
	End         *Point
	Points      []*Point
	PointByName map[string]*Point
	Children    map[*Point][]*Point
	Edges       []*Edge
	AllPaths    [][]*Point
	Paths       [][]*Point
}

func (g *Graph) ParsePaths(pathsStr string) error {
	if pathsStr == "" {
		return fmt.Errorf("AllPaths line is empty")
	}

	// Remove the prefix "All paths found:" or "Allocated paths:" from the pathsStr
	pathsStr = strings.TrimPrefix(pathsStr, "All paths found:")
	pathsStr = strings.TrimPrefix(pathsStr, "Allocated paths:")

	// Remove the leading and trailing white space, and then the brackets
	pathsStr = strings.TrimSpace(pathsStr)
	pathsStr = strings.Trim(pathsStr, "[]")

	// fmt.Printf("pathsStr: %v\n", pathsStr)

	// Split the pathsStr into individual path strings
	pathStrs := strings.SplitAfter(pathsStr, "]")
	// fmt.Printf("pathStrs: %v\n", pathStrs)

	// Initialize the result slice
	var paths [][]*Point

	for _, pathStr := range pathStrs {
		// Remove the "]" and "[" characters from pathStr
		pathStr = strings.Trim(pathStr, "[] ")

		// Skip empty pathStr
		if pathStr == "" {
			continue
		}

		// Split the pathStr into individual point names
		pointNames := strings.Split(pathStr, " ")

		// Initialize the path slice
		var path []*Point

		for _, pointName := range pointNames {
			point := g.PointByName[pointName]
			if point == nil {
				return fmt.Errorf("no such point: %v", pointName)
			}
			path = append(path, point)
		}

		paths = append(paths, path)
	}

	g.AllPaths = paths

	return nil
}
