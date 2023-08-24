// TODO: refactor and add comments
// TODO: add edges
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

	outputFile := "./solution.html"
	if flag.Arg(0) != "" {
		outputFile = flag.Arg(0)
	}

	var cfg string
	var solution string
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Fatal("unexpected EOF, solution should be separated from configuration with an empty line")
			}
			log.Fatal(err)
		}
		if strings.HasPrefix(text, "L") {
			solution += text
			break
		}
		cfg += text
	}

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		solution += text
	}

	graph, err := solver.ReadGraph(cfg)
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
