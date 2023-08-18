package linkroom

import "fmt"

func QueueAnts(NumberOfants int, Paths [][]int) {
	var RoomsOfPath []int
	for i := range Paths {
		fmt.Printf("How many Rooms in path %v := %v\n", i, len(Paths[i]))
		RoomsOfPath = append(RoomsOfPath, len(Paths[i]))
	}

	j := 0
	for i := 1; i <= NumberOfants; i++ {
		j = (i - 1) % len(RoomsOfPath) // Ant 1 starts from path 0, Ant 2 from path 1, and so on...
		fmt.Printf("Ant No.%v goes through path %v\n", i, j)
	}
}
