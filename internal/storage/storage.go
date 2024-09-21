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
func (s *Storage) ListUsers(prefix string) ([]user.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    results := s.users.PrefixSearch(strings.ToLower(prefix))
    users := make([]user.User, 0, len(results))

    for _, value := range results {
        if u, ok := value.(*user.User); ok {
            users = append(users, *u)
        } else {
            return nil, errors.New("invalid user data")
        }
    }
    return users, nil
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
func (s *Storage) ListFolders(username, sortField, sortOrder string) ([]folder.Folder, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    user, err := s.getUserNoLock(username)
    if err != nil {
        return nil, err
    }

    results := user.Folders.PrefixSearch("")
    folders := make([]folder.Folder, 0, len(results))
    for _, value := range results {
        if f, ok := value.(*folder.Folder); ok {
            folders = append(folders, *f)
        } else {
            return nil, errors.New("invalid folder data")
        }
    }

    sort.Slice(folders, func(i, j int) bool {
        if sortField == "created" {
            if sortOrder == "asc" {
                return folders[i].CreatedAt.Before(folders[j].CreatedAt)
            }
            return folders[j].CreatedAt.Before(folders[i].CreatedAt)
        }
        if sortOrder == "asc" {
            return folders[i].Name < folders[j].Name
        }
        return folders[i].Name > folders[j].Name
    })

    return folders, nil
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
func (s *Storage) ListFiles(username, folderName, sortField, sortOrder string) ([]file.File, error) {
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
    files := make([]file.File, 0, len(results))
    for _, value := range results {
        if f, ok := value.(*file.File); ok {
            files = append(files, *f)
        }
    }

    sort.Slice(files, func(i, j int) bool {
        if sortField == "created" {
            if sortOrder == "asc" {
                return files[i].CreatedAt.Before(files[j].CreatedAt)
            }
            return files[j].CreatedAt.Before(files[i].CreatedAt)
        }
        if sortOrder == "asc" {
            return files[i].Name < files[j].Name
        }
        return files[i].Name > files[j].Name
    })

    return files, nil
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