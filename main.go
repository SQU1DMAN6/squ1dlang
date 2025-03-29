package main

import (
	"fmt"
	// "os"
	"os/user"
	// "squ1d/repl"
)

// ~ a class, abstract thing, blueprint
type Programming struct {
	name     string
	yearborn int
}

//first method ==> it belongs to specific struct (Programming)
func (quanwhatever *Programming) manipulate() string {
	return quanwhatever.name  // 
}


func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the SQU1D programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	// repl.Start(os.Stdin, os.Stdout)

	
	//~ instance, concrete thing from the abstract one above
	var first = Programming{
		name:     "Q lang",
		yearborn: 2024,
	}


	things := first.manipulate()
	fmt.Println(things)

	//https://go.dev/play/p/a00CNWNxsla
}
