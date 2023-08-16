package linkroom

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// check the room if correct format, then call Addroom to add the room to slice of rooms
func CheckValidRoom(data string, t roomtype) {
	if !strings.HasPrefix(data, "#") {
		d := strings.Fields(data)
		if len(d) != 3 {
			log.Fatal("ERROR: invalid data format, excepted 'name x y'")
		}
		if strings.HasPrefix(d[0], "L") || strings.HasPrefix(d[0], "#") {
			log.Fatal("ERROR: invalid data format, the name of room can't start with the letter 'L' or with '#")
		}
		x, err := strconv.Atoi(d[1])
		if err != nil {
			err := fmt.Errorf("ERROR: invalid data format, coordinates 'x' should be an integer number, not \"%v\"", x)
			log.Fatal(err)
		}
		y, err := strconv.Atoi(d[2])
		if err != nil {
			err := fmt.Errorf("ERROR: invalid data format, coordinates 'x' should be an integer number, not \"%v\"", y)
			log.Fatal(err)
		}
		for _, r1 := range Rooms {
			if x == r1.X && y == r1.Y {
				err := fmt.Errorf("ERROR: invalid data format, these coordinates are already taken by room \"%v\"", r1.Name)
				log.Fatal(err)
			}
		}
		AddRoom(d[0], x, y, t)
	}
}

// check if Rooms slice has a start or end room
func CheckStartEnd() (bool, bool) {
	last := len(Rooms) - 1
	hasend, hasstart := false, false
	if last > 0 { // this check avoids index out of range errors if the slice is still empty
		hasend = Rooms[last].Rtype == End
		hasstart = Rooms[0].Rtype == Start
	}
	return hasstart, hasend
}
