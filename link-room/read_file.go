package linkroom

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var Links []string

// check the room, then call Addroom to add the room to slice of rooms
/*func ReadRoom(data string, t roomtype) {
	d := strings.Fields(data)
	x, _ := strconv.Atoi(d[1])
	y, _ := strconv.Atoi(d[2])
	AddRoom(d[0], x, y, t)
}
*/

// read file and seperate the data and restore in corespond slice or valiable
func Readfile(filename string) {
	data := ReadByLine(filename)
	//fmt.Printf("File data is:=%#v", data)

	if len(data) < 6 {
		log.Fatal("ERROR: invalid data format, file is too short")
	}
	int, err := strconv.Atoi(data[0])
	if err == nil {
		NumberOfants = int
		if NumberOfants == 0 {
			log.Fatal("ERROR: invalid data format, invalid number of Ants")
		}
	}
	for i := 1; i < len(data); i++ {
		//fmt.Printf("File data %v := %v\n", i, data[i])
		if strings.Contains(data[i], "-") {
			linkR := strings.Split(data[i], "-")
			if len(linkR) != 2 {
				log.Fatal("ERROR: invalid data format, links are defined by \"name1:name2\", and will usually look like \"1-2\", \"2-5\"")
			} else {
				Links = append(Links, data[i])
			}
			AddLink(linkR[0], linkR[1])
		} else {
			t := Normal
			if data[i] == "##start" {
				i++
				t = Start
			}
			if data[i] == "##end" {
				i++
				t = End
			}
			CheckValidRoom(data[i], t)
		}
	}
	//fmt.Println(Links)
}

// readByline will retrun a slice of string
func ReadByLine(filename string) []string {
	file, err := os.Open("./examples/" + filename) //open the file first
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var newstring []string
	for scanner.Scan() {
		readstring := scanner.Text() //read the whole line, then next line .....
		newstring = append(newstring, readstring)
	}
	return newstring
}
