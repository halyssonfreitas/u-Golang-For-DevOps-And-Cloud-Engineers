// execute 1 :
//   - go run main.go
//   - echo $?
//
// execute 2 :
//   - go run main.go argument1 argument2
//   - echo $?
package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./main <argument>\n")
		os.Exit(1)
	}

	fmt.Printf("hello world\nos.Args: %v\nArguments: %v\n", args, args[1:])
}
