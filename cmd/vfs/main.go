package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatbrother/virtual-file-system/internal/storage"
)

func main() {
	fmt.Println("Virtual File System")
	fmt.Println("Type 'exit' to quit the program")

	s := storage.NewStorage()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "exit" {
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		command := strings.ToLower(args[0])
		switch command {
		case "register":
			if len(args) != 2 {
				fmt.Println("Usage: register <username>")
				continue
			}
			username := args[1]
			err := s.AddUser(username)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("User '%s' registered successfully\n", username)
			}
		default:
			fmt.Println("Unknown command")
		}
	}

	fmt.Println("Goodbye!")
}