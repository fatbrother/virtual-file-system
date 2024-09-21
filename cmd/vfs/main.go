package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/fatbrother/virtual-file-system/internal/storage"
)

func main() {
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
                fmt.Fprintln(os.Stderr, "Usage: register <username>")
                continue
            }
            username := args[1]
            err := s.AddUser(username)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("User '%s' registered successfully\n", username)
            }
        case "delete":
            if len(args) != 2 {
                fmt.Fprintln(os.Stderr, "Usage: delete <username>")
                continue
            }
            username := args[1]
            err := s.DeleteUser(username)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("User %s deleted successfully\n", username)
            }
        case "create-folder":
            if len(args) < 3 {
                fmt.Fprintln(os.Stderr, "Usage: create-folder <username> <foldername> [description]")
                continue
            }
            username, folderName := args[1], args[2]
            description := strings.Join(args[3:], " ")
            err := s.CreateFolder(username, folderName, description)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("Create %s successfully.\n", folderName)
            }
        case "delete-folder":
            if len(args) != 3 {
                fmt.Fprintln(os.Stderr, "Usage: delete-folder <username> <foldername>")
                continue
            }
            username, folderName := args[1], args[2]
            err := s.DeleteFolder(username, folderName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("Delete %s successfully for user %s\n", folderName, username)
            }
        case "list-folders":
            if len(args) < 2 || len(args) > 4 {
                fmt.Fprintln(os.Stderr, "Usage: list-folders <username> [--sort-name|--sort-created] [asc|desc]")
                continue
            }
            username := args[1]
            sortField, sortOrder := "name", "asc"
            if len(args) > 2 {
                if args[2] == "--sort-name" || args[2] == "--sort-created" {
                    sortField = strings.TrimPrefix(args[2], "--sort-")
                } else {
                    fmt.Fprintln(os.Stderr, "Usage: list-folders <username> [--sort-name|--sort-created] [asc|desc]")
                    continue
                }
            }
            if len(args) > 3 {
                if args[3] == "asc" || args[3] == "desc" {
                    sortOrder = args[3]
                } else {
                    fmt.Fprintln(os.Stderr, "Usage: list-folders <username> [--sort-name|--sort-created] [asc|desc]")
                    continue
                }
            }
            folders, err := s.ListFolders(username, sortField, sortOrder)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else if len(folders) == 0 {
                fmt.Printf("No folders found for user %s\n", username)
            } else {
                fmt.Printf("Folders for user %s:\n", username)
                for _, folder := range folders {
                    fmt.Printf("- %s\n", folder.Format())
                }
            }
        case "create-file":
            if len(args) < 4 {
                fmt.Fprintln(os.Stderr, "Usage: create-file <username> <foldername> <filename> [description]")
                continue
            }
            username, folderName, fileName := args[1], args[2], args[3]
            description := strings.Join(args[4:], " ")
            err := s.CreateFile(username, folderName, fileName, description)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("Create %s in %s/%s successfully.\n", fileName, username, folderName)
            }
        case "delete-file":
            if len(args) != 4 {
                fmt.Fprintln(os.Stderr, "Usage: delete-file <username> <foldername> <filename>")
                continue
            }
            username, folderName, fileName := args[1], args[2], args[3]
            err := s.DeleteFile(username, folderName, fileName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else {
                fmt.Printf("Delete %s in %s/%s successfully.\n", fileName, username, folderName)
            }
        case "list-files":
            if len(args) < 3 || len(args) > 5 {
                fmt.Fprintln(os.Stderr, "Usage: list-files <username> <foldername> [--sort-name|--sort-created] [asc|desc]")
                continue
            }
            username, folderName := args[1], args[2]
            sortField, sortOrder := "name", "asc"
            if len(args) > 3 {
                if args[3] == "--sort-name" || args[3] == "--sort-created" {
                    sortField = strings.TrimPrefix(args[3], "--sort-")
                } else {
                    fmt.Fprintln(os.Stderr, "Usage: list-files <username> <foldername> [--sort-name|--sort-created] [asc|desc]")
                    continue
                }
            }
            if len(args) > 4 {
                if args[4] == "asc" || args[4] == "desc" {
                    sortOrder = args[4]
                } else {
                    fmt.Fprintln(os.Stderr, "Usage: list-files <username> <foldername> [--sort-name|--sort-created] [asc|desc]")
                    continue
                }
            }
            files, err := s.ListFiles(username, folderName, sortField, sortOrder)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            } else if len(files) == 0 {
                fmt.Printf("No files found in folder %s for user %s\n", folderName, username)
            } else {
                fmt.Printf("Files in folder %s for user %s:\n", folderName, username)
                for _, file := range files {
                    fmt.Printf("- %s\n", file.Format())
                }
            }
        case "help":
            fmt.Println("Commands:")
            fmt.Println("  register <username>")
            fmt.Println("  delete <username>")
            fmt.Println("  list [prefix]")
            fmt.Println("  create-folder <username> <foldername> [description]")
            fmt.Println("  delete-folder <username> <foldername>")
            fmt.Println("  list-folders <username> [--sort-name|--sort-created] [asc|desc]")
            fmt.Println("  create-file <username> <foldername> <filename> [description]")
            fmt.Println("  delete-file <username> <foldername> <filename>")
            fmt.Println("  list-files <username> <foldername> [--sort-name|--sort-created] [asc|desc]")
            fmt.Println("  help")
            fmt.Println("  exit")
        default:
            fmt.Fprintln(os.Stderr, "Unknown command")
        }
    }

    fmt.Println("Goodbye!")
}