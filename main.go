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

	// print the program output
	fmt.Println("\nProgram output:")
	link.ProgramOutput()
}
