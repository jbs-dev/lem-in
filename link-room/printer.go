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
}
