package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Println("try pwd docs/")
		return
	} else if len(os.Args) == 2 {
		log.Printf("try docs/ from '%s'", os.Args[1])
		return
	}

	switch os.Args[1] {
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}
