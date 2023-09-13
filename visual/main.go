package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"lem-in/solver"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//go:embed solution.gohtml
var tmplFile string

const setBold = "\033[1m"
const reset = "\033[0m"

type AntAnimations struct {
	XAnimations string
	YAnimations string
}

type GraphPage struct {
	Graph          *solver.Graph
	Ants           []AntAnimations
	ViewBox        string
	RoomHeight     string
	RoomWidth      string
	RoomHeightHalf string
	RoomWidthHalf  string
	EdgeWidth      string
	AllPaths       [][]*solver.Point
	Paths          [][]*solver.Point
}

const HELP = `Usage: ` + setBold + `./visualizer [file]` + reset + ` to read from input and visualize solution to [file]
Example: ` + setBold + `lem-in ./examples/example00.txt | ./visualizer` + reset + ` to solve example00 and visualize it`

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Println(HELP)
	}
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)

	outputFile := "./visual/solution.html"
	if flag.Arg(0) != "" {
		outputFile = flag.Arg(0)
	}

	var cfg string
	var solution string
	var allPathsStr string
	var pathsStr string
	readingCfg := true
	for {
		text, err := reader.ReadString('\n')
		// fmt.Printf("Read line: %v\n", text)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if strings.HasPrefix(text, "All paths found:") {
			allPathsStr = text
		} else if strings.HasPrefix(text, "Allocated paths:") {
			pathsStr = text
		} else if strings.HasPrefix(text, "L") {
			solution += text
			readingCfg = false
		} else if readingCfg {
			cfg += text
		} else {
			solution += text
		}
	}
	// fmt.Printf("AllPathsStr: %v\n", allPathsStr) // New log statement
	// fmt.Printf("PathsStr: %v\n", pathsStr)       // New log statement

	graph, err := solver.ReadGraph(cfg, allPathsStr, pathsStr)
	if err != nil {
		log.Fatal(err)
	}

	allPathsStr = strings.TrimSpace(allPathsStr) // New line
	pathsStr = strings.TrimSpace(pathsStr)       // New line

	// fmt.Printf("All paths found: %v\n", allPathsStr) // New log statement
	// fmt.Printf("Allocated paths: %v\n", pathsStr)    // New log statement

	allPaths, err := parsePaths(allPathsStr, graph)
	if err != nil {
		log.Fatal(err)
	}
	paths, err := parsePaths(pathsStr, graph)
	if err != nil {
		log.Fatal(err)
	}

	solution = strings.ReplaceAll(solution, "\r", "")

	lines := strings.Split(solution, "\n")

	ants := make([][]AntAnimations, graph.AntNum)
	for i := range ants {
		ants[i] = make([]AntAnimations, len(lines))
	}

	for i, line := range strings.Split(solution, "\n") {
		if line == "" {
			continue
		}
		commands := strings.Split(line, " ")

		parseCommand := func(command string) error {
			if !strings.HasPrefix(command, "L") {
				return fmt.Errorf("should begin with the letter L")
			}
			command = strings.TrimPrefix(command, "L")
			antNumPointName := strings.Split(command, "-")
			if len(antNumPointName) != 2 {
				return fmt.Errorf("wrong format, expected 'LantNum-pointName'(e.g. 'L0-0'))")
			}
			antNum, err := strconv.Atoi(antNumPointName[0])
			if err != nil || antNum <= 0 {
				return fmt.Errorf("ant number %v is not valid, expected positive integer", antNum)
			}
			pointName := antNumPointName[1]
			point := graph.PointByName[pointName]
			if point == nil {
				return fmt.Errorf("no such point: %v", pointName)
			}

			ants[antNum-1][i].XAnimations = strconv.Itoa(point.X)
			ants[antNum-1][i].YAnimations = strconv.Itoa(point.Y)

			return nil
		}
		for _, command := range commands {
			err := parseCommand(command)
			if err != nil {
				log.Fatalf("invalid command '%v' in line %v ('%v'): %v", command, i, line, err)
			}
		}
	}

	tmpl, err := template.New("solution").Parse(tmplFile)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	minX, minY, maxX, maxY := graph.Start.X, graph.Start.Y, graph.Start.X, graph.Start.Y

	longestNameSize := len(graph.Start.Name)

	if graph.End == nil {
		log.Fatal("graph.End is nil")
	}
	for i, point := range graph.Points {
		if point == nil {
			log.Fatalf("graph.Points[%d] is nil", i)
		}
	}

	for _, point := range append(graph.Points, graph.End) {
		if len(point.Name) > longestNameSize {
			longestNameSize = len(point.Name)
		}
		if point.X < minX {
			minX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	roomHeight := float32(maxY-minY+2) / 15
	roomWidth := roomHeight * 1.5

	if roomWidth < roomHeight/3*float32(longestNameSize) {
		roomWidth = roomHeight / 3 * float32(longestNameSize)
	}

	antAnimations := make([]AntAnimations, len(ants))
	for ant := range ants {
		animation := 0
		for ; animation < len(ants[ant]) && ants[ant][animation].XAnimations == ""; animation++ {
			antAnimations[ant].XAnimations += strconv.Itoa(graph.Start.X) + ";"
			antAnimations[ant].YAnimations += strconv.Itoa(graph.Start.Y) + ";"
		}
		for ; animation < len(ants[ant]) && ants[ant][animation].XAnimations != ""; animation++ {
			antAnimations[ant].XAnimations += ants[ant][animation].XAnimations + ";"
			antAnimations[ant].YAnimations += ants[ant][animation].YAnimations + ";"
		}
		for ; animation < len(ants[ant]); animation++ {
			antAnimations[ant].XAnimations += strconv.Itoa(graph.End.X) + ";"
			antAnimations[ant].YAnimations += strconv.Itoa(graph.End.Y) + ";"
		}
	}

	err = tmpl.Execute(file,
		GraphPage{
			Graph:          graph,
			Ants:           antAnimations,
			AllPaths:       allPaths,
			Paths:          paths,
			ViewBox:        fmt.Sprintf("%v %v %v %v", minX-1, minY-1, maxX-minX+2, maxY-minY+2),
			RoomHeight:     fmt.Sprintf("%.2f", roomHeight),
			RoomWidth:      fmt.Sprintf("%.2f", roomWidth),
			RoomHeightHalf: fmt.Sprintf("%.2f", roomHeight/2),
			RoomWidthHalf:  fmt.Sprintf("%.2f", roomWidth/2),
			EdgeWidth:      fmt.Sprintf("%.2f", roomHeight/5),
		})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Visualization saved to %v\n", file.Name())
}

func parsePaths(pathsStr string, graph *solver.Graph) ([][]*solver.Point, error) {
	// fmt.Printf("parsePaths: pathsStr: %v\n", pathsStr) // New log statement

	if pathsStr == "" {
		return nil, fmt.Errorf("pathsStr is empty")
	}

	// Remove the prefix "All paths found:" or "Allocated paths:" from the pathsStr
	pathsStr = strings.TrimPrefix(pathsStr, "All paths found:")
	pathsStr = strings.TrimPrefix(pathsStr, "Allocated paths:")

	// fmt.Printf("pathsStr: %v\n", pathsStr)

	// Remove the leading and trailing white space, and then the brackets
	pathsStr = strings.TrimSpace(pathsStr)
	pathsStr = strings.Trim(pathsStr, "[]")

	// fmt.Printf("pathsStr: %v\n", pathsStr)

	// Split the pathsStr into individual paths
	pathsStrs := strings.Split(pathsStr, "] [")
	// fmt.Printf("parsePaths: pathsStrs: %v\n", pathsStrs)

	// Initialize the result slice
	var paths [][]*solver.Point

	for _, pathStr := range pathsStrs {
		// Split the pathStr into individual points
		pointStrs := strings.Split(pathStr, " ")
		// fmt.Printf("parsePaths: pointStrs: %v\n", pointStrs)

		// Initialize the path slice
		var path []*solver.Point

		for _, pointName := range pointStrs {
			if pointName == "" {
				fmt.Printf("parsePaths: ERROR empty pointName in pointStrs: %v\n", pointStrs)
				continue
			}
			point := graph.PointByName[pointName]
			if point == nil {
				fmt.Printf("parsePaths: ERROR no such point in pointStrs: %v\n", pointStrs)
				return nil, fmt.Errorf("no such point: %v", pointName)
			}
			path = append(path, point)
		}

		paths = append(paths, path)
	}

	return paths, nil
}
