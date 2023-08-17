package linkroom

import (
	"fmt"
	"log"
)

// declare a new type "roomtype", so that we can conviently refer to it
type roomtype byte

const (
	Start  roomtype = 0
	End    roomtype = 1
	Normal roomtype = 2
)

type Room struct {
	Name  string
	X, Y  int      // coordinates
	Rtype roomtype // whether the room is a start, finsish or nomal room
	Links []*Room  // links to other rooms
}

var Rooms []Room
var NumberOfants int

// function to add a room, create a []Rooms. start from "start room", to "end room"
func AddRoom(name string, x_cord, y_cord int, t roomtype) {
	// check if room name already exists
	for _, r1 := range Rooms {
		if r1.Name == name {
			err := fmt.Errorf("ERROR: Room name \"%s\" used more than once", name)
			log.Fatal(err)
		}
	}

	newRoom := Room{Name: name, X: x_cord, Y: y_cord, Rtype: t}
	newRooms := []Room{newRoom}

	hasstart, hasend := CheckStartEnd() // check if start or end rooms already exist

	last := len(Rooms) - 1 //index of the last element of the Room
	// Now insert the newRoom at the appropriate position in the Rooms slice,
	// depending on whether it is a start, finsish or normal room
	switch t {
	case Start:
		if hasstart { // check if there already is a Start room
			err := fmt.Errorf("ERROR: Room name \"%s\" cannot be the Start, because Start is already at room \"%s\"", name, Rooms[0].Name)
			log.Fatal(err)
		}

		// insert start room at beginning of the slice
		/* Example usage
		a := []int{1, 2}
		b := []int{11, 22}
		a = append(a, b...) // a == [1 2 11 22]
		The ... unpacks b. Without the dots, the code would attempt to append the slice as a whole, which is invalid.*/
		Rooms = append(newRooms, Rooms...)

	case End:
		if hasend { // check if there already is a end room
			err := fmt.Errorf("ERROR: Room name \"%s\" cannot be the end, because end is already at room \"%s\"", name, Rooms[last].Name)
			log.Fatal(err)
		}
		// append end room at end of the slice
		Rooms = append(Rooms, newRoom)
	case Normal:
		if hasend {
			Rooms = append(Rooms[:last+1], Rooms[last:]...) // append just before the end Room
			Rooms[last] = newRoom
		} else { // otherwise just append at the end of the slice
			Rooms = append(Rooms, newRoom)
		}
	}
}

// function to add a link between rooms matching room names r1 and r2
func AddLink(n1, n2 string) {
	if n1 == n2 { // cannot create a link to itself
		err := fmt.Errorf("ERROR: Link \"%v-%v\" links to itself", n1, n2)
		log.Fatal(err)
	}

	for i := range Rooms {
		r1 := &Rooms[i] // grab pointer to the room in outer loop i
		if r1.Name == n1 {
			for j := range Rooms {
				r2 := &Rooms[j]
				if r2.Name == n2 {
					for _, link := range r1.Links { // check if rooms alredy linked
						if link == r2 {
							return
						}
					} // if rooms not yet linked -> create link data
					r1.Links = append(r1.Links, r2)
					r2.Links = append(r2.Links, r1)
					return
				}
			}
		}
	}
	err := fmt.Errorf("ERROR: Link \"%v-%v\" cannot be created. One or both of the rooms do not exist", n1, n2)
	log.Fatal(err)
}
