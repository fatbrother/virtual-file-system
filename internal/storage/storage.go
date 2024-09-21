package storage

import (
    "errors"
    "sort"
    "strings"
    "sync"

    "github.com/fatbrother/virtual-file-system/internal/user"
    "github.com/fatbrother/virtual-file-system/internal/folder"
    "github.com/fatbrother/virtual-file-system/internal/file"
    "github.com/fatbrother/virtual-file-system/pkg/trie"
)

// Storage represents the in-memory storage for the virtual file system
type Storage struct {
    users *trie.Trie
    mu    sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
    return &Storage{
        users: trie.NewTrie(),
    }
}

// AddUser adds a new user to the storage
func (s *Storage) AddUser(username string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    lowercaseUsername := strings.ToLower(username)
    if _, exists := s.users.Search(lowercaseUsername); exists {
        return errors.New("The " + username + " has already existed.")
    }

    newUser, err := user.NewUser(username)
    if err != nil {
        return err
    }

    s.users.Insert(lowercaseUsername, newUser)
    return nil
}

// GetUser retrieves a user from the storage
func (s *Storage) GetUser(username string) (*user.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    lowercaseUsername := strings.ToLower(username)
    if value, exists := s.users.Search(lowercaseUsername); exists {
        if user, ok := value.(*user.User); ok {
            return user, nil
        }
        return nil, errors.New("invalid user data")
    }
    return nil, errors.New("The " + username + " not found.")
}

// DeleteUser removes a user from the storage
func (s *Storage) DeleteUser(username string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    lowercaseUsername := strings.ToLower(username)
    if deleted := s.users.Delete(lowercaseUsername); !deleted {
        return errors.New("The " + username + " not found.")
    }
    return nil
}

// ListUsers returns a list of all usernames with the given prefix
func (s *Storage) ListUsers(prefix string) []string {
    s.mu.RLock()
    defer s.mu.RUnlock()

    results := s.users.PrefixSearch(strings.ToLower(prefix))
    usernames := make([]string, 0, len(results))
    for username := range results {
        usernames = append(usernames, username)
    }
    return usernames
}

// CreateFolder creates a new folder for a user
func (s *Storage) CreateFolder(username, folderName, description string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return err
    }

    lowercaseFolderName := strings.ToLower(folderName)
    if _, exists := user.Folders.Search(lowercaseFolderName); exists {
        return errors.New("The " + folderName + " has already existed.")
    }

    newFolder, err := folder.NewFolder(folderName, description)
    if err != nil {
        return err
    }

    user.Folders.Insert(lowercaseFolderName, newFolder)
    return nil
}

// DeleteFolder deletes a folder for a user
func (s *Storage) DeleteFolder(username, folderName string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return err
    }

    lowercaseFolderName := strings.ToLower(folderName)
    if deleted := user.Folders.Delete(lowercaseFolderName); !deleted {
        return errors.New("The " + folderName + " not found.")
    }

    return nil
}

// ListFolders returns a list of all folders for a user with sorting options
func (s *Storage) ListFolders(username, sortField, sortOrder string) ([]string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return nil, err
    }

    results := user.Folders.PrefixSearch("")
    folderNames := make([]string, 0, len(results))
    for folderName := range results {
        folderNames = append(folderNames, folderName)
    }

    sort.Slice(folderNames, func(i, j int) bool {
        if sortField == "created" {
            folderI, _ := user.Folders.Search(folderNames[i])
            folderJ, _ := user.Folders.Search(folderNames[j])
            folderICreatedAt := folderI.(*folder.Folder).CreatedAt
            folderJCreatedAt := folderJ.(*folder.Folder).CreatedAt
            if sortOrder == "asc" {
                return folderICreatedAt.Before(folderJCreatedAt)
            }
            return folderJCreatedAt.Before(folderICreatedAt)
        }
        if sortOrder == "asc" {
            return folderNames[i] < folderNames[j]
        }
        return folderNames[i] > folderNames[j]
    })

    return folderNames, nil
}

// CreateFile creates a new file in a folder for a user
func (s *Storage) CreateFile(username, folderName, fileName, description string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return err
    }

    folder, err := s.getFolderNoLock(user, folderName)
    if err != nil {
        return err
    }

    lowercaseFileName := strings.ToLower(fileName)
    if _, exists := folder.Files.Search(lowercaseFileName); exists {
        return errors.New("The " + fileName + " has already existed.")
    }

    newFile, err := file.NewFile(fileName, description)
    if err != nil {
        return err
    }

    folder.Files.Insert(lowercaseFileName, newFile)
    return nil
}

// DeleteFile deletes a file from a folder for a user
func (s *Storage) DeleteFile(username, folderName, fileName string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return err
    }

    folder, err := s.getFolderNoLock(user, folderName)
    if err != nil {
        return err
    }

    lowercaseFileName := strings.ToLower(fileName)
    if deleted := folder.Files.Delete(lowercaseFileName); !deleted {
        return errors.New("The " + fileName + " not found.")
    }

    return nil
}

// ListFiles returns a list of all files in a folder for a user with sorting options
func (s *Storage) ListFiles(username, folderName, sortField, sortOrder string) ([]string, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return nil, err
    }

    folder, err := s.getFolderNoLock(user, folderName)
    if err != nil {
        return nil, err
    }

    results := folder.Files.PrefixSearch("")
    fileNames := make([]string, 0, len(results))
    for fileName := range results {
        fileNames = append(fileNames, fileName)
    }

    sort.Slice(fileNames, func(i, j int) bool {
        if sortField == "created" {
            fileI, _ := folder.Files.Search(fileNames[i])
            fileJ, _ := folder.Files.Search(fileNames[j])
            fileICreatedAt := fileI.(*file.File).CreatedAt
            fileJCreatedAt := fileJ.(*file.File).CreatedAt
            if sortOrder == "asc" {
                return fileICreatedAt.Before(fileJCreatedAt)
            }
            return fileJCreatedAt.Before(fileICreatedAt)
        }
        if sortOrder == "asc" {
            return fileNames[i] < fileNames[j]
        }
        return fileNames[i] > fileNames[j]
    })

    return fileNames, nil
}

// getUserNoLock retrieves a user without locking (assumes caller holds the lock)
func (s *Storage) getUserNoLock(username string) (*user.User, error) {
    lowercaseUsername := strings.ToLower(username)
    if value, exists := s.users.Search(lowercaseUsername); exists {
        if user, ok := value.(*user.User); ok {
            return user, nil
        }
        return nil, errors.New("invalid user data")
    }
    return nil, errors.New("The " + username + " not found.")
}

// getFolderNoLock retrieves a folder
func (s *Storage) getFolderNoLock(user *user.User, folderName string) (*folder.Folder, error) {
    lowercaseFolderName := strings.ToLower(folderName)
    if value, exists := user.Folders.Search(lowercaseFolderName); exists {
        if folder, ok := value.(*folder.Folder); ok {
            return folder, nil
        }
        return nil, errors.New("invalid folder data")
    }
    return nil, errors.New("The " + folderName + " not found.")
}