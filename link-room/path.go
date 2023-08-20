package linkroom

import (
	"fmt"
)

var Path []int

// to store all index of rooms from []Rooms for a valid path ex:Room[r].Links[i]
var Paths [][]int

// This function will keep call itself until get all valid paths.
// A valid path needs to start with start room and finish with end room.
func FindValidPaths(currentRoomIndex int) bool {

	// colour definitions just for debug messages
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorCyan := "\033[36m"
	colorYellow := "\033[33m"

	currentRoom := &Rooms[currentRoomIndex]
	currentRoom.Visited = true
	fmt.Printf("%sNow in Room Index:= %v%s\n", colorCyan, currentRoomIndex, colorReset)

	for i, linkedIndex := range currentRoom.Links {
		fmt.Printf("Loop turn: %v, Checking Room Index:= %v\n", i, linkedIndex)
		linkedRoom := &Rooms[linkedIndex]
		// if linkedRoom is already visited, try next one
		if linkedRoom.Visited {
			fmt.Printf("Room Index %v already visted, trying next one\n", linkedIndex)
			continue
		}
		// if LinkedRoom is End, complete and save the Path and then backtrack to start
		if linkedRoom.Rtype == End {
			fmt.Printf("%sRoom Index %v is End, completing Path%s\n", colorYellow, linkedIndex, colorReset)
			Path = append(Path, linkedIndex)
			// last element of Path
			fmt.Printf("path is now :=%v\n", Path)
			// save Path t Paths
			Paths = append(Paths, Path)
			fmt.Println("All valid paths:=", Paths)
			// reset Path
			Path = nil
			Path = append(Path, 0)
			return true
		}

		// linkedRoom is neither start,visited or end
		// so, we append to path and recurse into it
		fmt.Printf("Append index := %v\n", linkedIndex)
		//fmt.Printf("Append Room Info:= %v\n", linkedRoom)
		Path = append(Path, linkedIndex)
		fmt.Printf("path is now :=%v\n", Path)
		fmt.Printf("recursing into %v\n", linkedIndex)
		// recursive function itself
		result := FindValidPaths(linkedIndex)
		fmt.Printf("%sbacktracking into Room %v%s\n", colorCyan, currentRoomIndex, colorReset)
		if result && currentRoomIndex != 0 {
			return true
		}
	} // end of the loop
	// Backtrack: Mark the current room as unvisited and remove it from the path
	fmt.Println(colorRed, "No path found.Start backtracking", colorReset)
	currentRoom.Visited = false
	Path = Path[:len(Path)-1]
	return false
}
