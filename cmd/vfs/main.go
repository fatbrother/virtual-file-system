package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Virtual File System")
	fmt.Println("Type 'exit' to quit the program")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "exit" {
			break
		}

		// TODO: Process the command
		fmt.Printf("You entered: %s\n", input)
	}

	fmt.Println("Goodbye!")
}