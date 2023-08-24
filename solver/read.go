package solver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ReadGraph reads graph configuration from file to Graph struct
//
// Returns an error if configuration is not valid
func ReadGraph(source string) (*Graph, error) {
	// remove all \r for strings created by MS DOS users
	cfgStr := strings.ReplaceAll(source, "\r", "")
	cfgStr += "\n"
	cfgStr = strings.ReplaceAll(cfgStr, "\n\n", "\n")

	m1 := regexp.MustCompile(` +`)
	cfgStr = m1.ReplaceAllLiteralString(cfgStr, " ")

	m2 := regexp.MustCompile(`^ +`)
	cfgStr = m2.ReplaceAllLiteralString(cfgStr, "")

	cfgStr = strings.ReplaceAll(cfgStr, " \n", "\n")

	// make ##start and ##end declarations inline
	// from:
	// ##start
	// start 0 0
	// to:
	// ##startstart 0 0
	cfgStr = strings.NewReplacer("##start\n", "##start", "##end\n", "##end").Replace(cfgStr)

	cfgLines := strings.Split(cfgStr, "\n")

	// antNum, start, end, edge
	if len(cfgLines) < 4 {
		return nil, fmt.Errorf("file is too short")
	}

	antNum, err := strconv.Atoi(cfgLines[0])
	if err != nil {
		return nil, fmt.Errorf("ant number '%v' is not valid: %w:wq:wq", cfgLines[0], err)
	}

	if antNum <= 0 {
		return nil, fmt.Errorf("ant number '%v' is not valid, must be positive integer", antNum)
	}

	graph := &Graph{
		AntNum:      antNum,
		PointByName: make(map[string]*Point),
	}

	edgeCfgStart := len(cfgLines)
	for i := 1; i < len(cfgLines); i++ {
		line := cfgLines[i]
		if strings.Contains(line, "-") {
			// the end of points' declaration, edges' declaration started
			edgeCfgStart = i
			break
		}
		err := graph.parsePointLine(line)
		if err != nil {
			return nil, fmt.Errorf("invalid point in line %v ('%v'): %v", i+1, line, err)
		}
	}

	if graph.Start == nil {
		return nil, fmt.Errorf("##start is not defined")
	}
	if graph.End == nil {
		return nil, fmt.Errorf("##end is not defined")
	}

	if edgeCfgStart == len(cfgLines) {
		return nil, fmt.Errorf("no edges defined")
	}

	for i := edgeCfgStart; i < len(cfgLines); i++ {
		line := cfgLines[i]
		if strings.HasPrefix(line, "#") || line == "" {
			continue // ignore comments and empty lines
		}
		edge, err := graph.ParseEdge(line)
		if err != nil {
			return nil, fmt.Errorf("invalid edge in line %v ('%v'): %v", i+1, line, err)
		}
		graph.Edges = append(graph.Edges, edge)
	}
	return graph, nil
}

// parsePointLine parses line from point configuration and adds it to Graph
// it takes care of comments, start and end points
func (g *Graph) parsePointLine(line string) error {
	var point *Point
	var err error

	switch {
	case strings.HasPrefix(line, "##start"):
		if g.Start != nil {
			return fmt.Errorf("only one ##start allowed")
		}
		line = strings.TrimPrefix(line, "##start")
		point, err = ParsePoint(line)
		if err != nil {
			return err
		}
		g.Start = point

	case strings.HasPrefix(line, "##end"):
		if g.End != nil {
			return fmt.Errorf("only one ##end allowed")
		}
		line = strings.TrimPrefix(line, "##end")
		point, err = ParsePoint(line)
		if err != nil {
			return err
		}
		g.End = point

	case strings.HasPrefix(line, "#") || line == "":
		// ignore comments and empty lines
		return nil

	default:
		point, err = ParsePoint(line)
		if err != nil {
			return err
		}
		if g.PointByName[point.Name] != nil {
			return fmt.Errorf("room name [%s] is duplicated", point.Name)
		}
		for _, p := range g.PointByName {
			if p.X == point.X && p.Y == point.Y {
				return fmt.Errorf("room coordinates [x %d, y %d] are duplicated", p.X, p.Y)
			}
		}

		g.Points = append(g.Points, point)
	}

	g.PointByName[point.Name] = point
	return nil
}
