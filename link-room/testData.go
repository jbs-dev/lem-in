package linkroom

// hardcoded data to demonstrate the structure
// in the real program, this will be read from file
func MakeTestData() {

	NumberOfants = 3

	ReadRoom("0 2 0", Normal)
	AddRoom("a", 12, 3, Normal)
	AddRoom("c", 5, 5, Start)
	AddRoom("b", 5, 2, Normal)
	AddRoom("d", 7, 3, End)
	AddRoom("k", 5, 2, Normal)

	//addRoom("a", 14, 2) // error test trying to add another room with the same name
	//addRoom("e", 14, 2, start) // error test trying to add another start room
	//AddRoom("f", 14, 2, End) // error test trying to add another end room

}
