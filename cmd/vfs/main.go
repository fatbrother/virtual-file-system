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
		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <username>")
				continue
			}
			username := args[1]
			err := s.DeleteUser(username)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("User '%s' deleted successfully\n", username)
			}
		case "list":
			prefix := ""
			if len(args) > 1 {
				prefix = args[1]
			}
			users := s.ListUsers(prefix)
			if len(users) == 0 {
				fmt.Println("No users found")
			} else {
				fmt.Println("Users:")
				for _, user := range users {
					fmt.Printf("- %s\n", user)
				}
			}
		case "create-folder":
			if len(args) < 3 {
				fmt.Println("Usage: create-folder <username> <foldername> [description]")
				continue
			}
			username, folderName := args[1], args[2]
			description := strings.Join(args[3:], " ")
			err := s.CreateFolder(username, folderName, description)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Folder '%s' created successfully for user '%s'\n", folderName, username)
			}
		case "delete-folder":
			if len(args) != 3 {
				fmt.Println("Usage: delete-folder <username> <foldername>")
				continue
			}
			username, folderName := args[1], args[2]
			err := s.DeleteFolder(username, folderName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Folder '%s' deleted successfully for user '%s'\n", folderName, username)
			}
		case "list-folders":
			if len(args) != 2 {
				fmt.Println("Usage: list-folders <username>")
				continue
			}
			username := args[1]
			folders, err := s.ListFolders(username)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if len(folders) == 0 {
				fmt.Printf("No folders found for user '%s'\n", username)
			} else {
				fmt.Printf("Folders for user '%s':\n", username)
				for _, folder := range folders {
					fmt.Printf("- %s\n", folder)
				}
			}
		case "create-file":
			if len(args) < 4 {
				fmt.Println("Usage: create-file <username> <foldername> <filename> <description>")
				continue
			}
			username, folderName, fileName := args[1], args[2], args[3]
			description := strings.Join(args[4:], " ")
			err := s.CreateFile(username, folderName, fileName, description)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("File '%s' created successfully in folder '%s' for user '%s'\n", fileName, folderName, username)
			}
		case "delete-file":
			if len(args) != 4 {
				fmt.Println("Usage: delete-file <username> <foldername> <filename>")
				continue
			}
			username, folderName, fileName := args[1], args[2], args[3]
			err := s.DeleteFile(username, folderName, fileName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("File '%s' deleted successfully from folder '%s' for user '%s'\n", fileName, folderName, username)
			}
		case "list-files":
			if len(args) != 3 {
				fmt.Println("Usage: list-files <username> <foldername>")
				continue
			}
			username, folderName := args[1], args[2]
			files, err := s.ListFiles(username, folderName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if len(files) == 0 {
				fmt.Printf("No files found in folder '%s' for user '%s'\n", folderName, username)
			} else {
				fmt.Printf("Files in folder '%s' for user '%s':\n", folderName, username)
				for _, file := range files {
					fmt.Printf("- %s\n", file)
				}
			}
		case "help":
			fmt.Println("Commands:")
			fmt.Println("  register <username>")
			fmt.Println("  delete <username>")
			fmt.Println("  list [prefix]")
			fmt.Println("  create-folder <username> <foldername> [description]")
			fmt.Println("  delete-folder <username> <foldername>")
			fmt.Println("  list-folders <username>")
			fmt.Println("  create-file <username> <foldername> <filename> <description>")
			fmt.Println("  delete-file <username> <foldername> <filename>")
			fmt.Println("  list-files <username> <foldername>")
			fmt.Println("  help")
			fmt.Println("  exit")
		default:
			fmt.Println("Unknown command")
		}
	}

	fmt.Println("Goodbye!")
}