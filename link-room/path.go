package linkroom

import "fmt"

var Path []int
var Paths [][]int

// This function will keep call itself until get all valid paths.
// A valid path needs to start with start room and finish with end room.
func FindValidPaths(r int) [][]int {
	// to store all index of rooms from []Rooms for a valid path ex:Room[r].Links[i]
	for i := 0; i < len(Rooms[r].Links); i++ {
		/*if i == len(Rooms[r].Links)-1 {
			return Paths
		}*/
		// counter count how many times reach to End
		counter := 1
		if Rooms[Rooms[r].Links[i]].Rtype != Start && !Rooms[Rooms[r].Links[i]].Visited {
			fmt.Printf("Visit Room Index:= %v\n", Rooms[r].Links[i])
			fmt.Println(counter)
			if Rooms[Rooms[r].Links[i]].Rtype == End { // Found a whole valid path, append to paths slice
				Path = append(Path, Rooms[r].Links[i])
				Paths = append(Paths, Path)
				//fmt.Println("All valid paths:=", Paths)
				Path = nil
				Path = append(Path, 0)
				counter++
				if counter == len(Rooms[0].Links) {
					Path = append(Path, Rooms[r].Links[i])
					Paths = append(Paths, Path)
					return Paths
				} else {
					return FindValidPaths(0) // recursive function itself
				}
			} else {
				Rooms[Rooms[r].Links[i]].Visited = true
				fmt.Printf("Append index := %v\n", Rooms[r].Links[i])
				Path = append(Path, Rooms[r].Links[i])
				fmt.Printf("path is now :=%v\n", Path)
				return FindValidPaths(Rooms[r].Links[i]) // recursive function itself
			}
		}
	}

	return Paths
}
