package main

import (
	"fmt"
	link "lem-in/link-room"
	"log"
	"os"
)

func main() {
	filename := os.Args[1:]
	if len(filename) < 1 {
		err := fmt.Errorf("ERROR: Please provide a filename")
		log.Fatal(err)
	}
	link.Readfile(filename[0])

	//temp function, hard coding. later make a function read from file
	//link.MakeTestData()

	hasstart, hasend := link.CheckStartEnd()
	if !(hasstart && hasend) {
		err := fmt.Errorf("ERROR: No start room or no end room defined")
		log.Fatal(err)
	}

	link.PrintRoomInfo()
	fmt.Printf("\n----------------------- %s-----------------------", filename)
	fmt.Println("\nProgram output:")
	link.ProgramOutput()
	start := 0 //parse first element of []Rooms witch is start room
	link.Path = append(link.Path, 0)
	link.FindValidPaths(start)
	fmt.Println("All valid paths:=\n", link.Paths)
	link.QueueAnts(link.NumberOfants, link.Paths)
	//link.SendingAnts()
}
