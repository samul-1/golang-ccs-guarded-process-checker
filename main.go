package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	/* Example tests

	  not guarded
		rec x. (a|x)+b
		rec x. x
		rec x. a.rec y. x+y
		rec x. a.rec y. x|y

		guarded
		rec x. a.x+b
		rec x. a|b.x
		rec x. a.rec y. x
		rec x. a.rec y. x|b.y
	*/
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter process: ")
	input, _ := reader.ReadString('\n')

	// lex & parse the input string into a CCS process
	l := NewLexer(input)
	p := NewParser(l)

	astProcess := p.ParseAstProcess()

	// convert the ast version into a Process object
	process := NewProcessFromInterface(astProcess)

	// create a channel for receiving the result of the check
	resultChanP := make(chan bool)

	// test G(P, emptyset)
	go process.IsGuarded(map[string]bool{}, resultChanP)

	if <-resultChanP {
		fmt.Println("Entered process is guarded")
	} else {
		fmt.Println("Entered process is not guarded")
	}
}
