package linkroom

import "fmt"

func ProgramOutput() {

	fmt.Println(NumberOfants)
	for i, curr := range Rooms {
		if i == 0 {
			fmt.Println("##start")
		}
		if i == len(Rooms)-1 {
			fmt.Println("##end")
		}
		fmt.Printf("%s ", curr.Name)
		fmt.Printf("%d %d", curr.X, curr.Y)
		fmt.Println()
	}
	for i := range Links {
		fmt.Println(Links[i])
	}
}

func PrintRoomInfo() {
	/*fmt.Println("Room Slice:")
	fmt.Println(Rooms)
	fmt.Println()*/
	for _, r1 := range Rooms {
		fmt.Println("\nRoom Info:=", r1)
		fmt.Printf("Room %v:\n\tCoordinates:(%v,%v) \n", r1.Name, r1.X, r1.Y)
		fmt.Printf("\tType: %d\n", r1.Rtype)
		fmt.Printf("\tLinks: ")
		for i, link := range r1.Links {
			fmt.Printf("%v", link.Name)
			if i < len(r1.Links)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}
}
