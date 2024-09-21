package storage

import (
	"errors"
	"strings"
	"sync"

	"github.com/fatbrother/virtual-file-system/internal/user"
	"github.com/fatbrother/virtual-file-system/internal/folder"
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
		return errors.New("user already exists")
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
	return nil, errors.New("user not found")
}


// DeleteUser removes a user from the storage
func (s *Storage) DeleteUser(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lowercaseUsername := strings.ToLower(username)
	if deleted := s.users.Delete(lowercaseUsername); !deleted {
		return errors.New("user not found")
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
		return errors.New("folder already exists")
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
		return errors.New("folder not found")
	}

	return nil
}

// ListFolders returns a list of all folders for a user
func (s *Storage) ListFolders(username string) ([]string, error) {
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
	return folderNames, nil
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
	return nil, errors.New("user not found")
}