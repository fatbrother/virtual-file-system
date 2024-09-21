package storage

import (
	"errors"
	"strings"
	"sync"

	"github.com/fatbrother/virtual-file-system/internal/user"
)

// Storage represents the in-memory storage for the virtual file system
type Storage struct {
	users map[string]*user.User
	mu    sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		users: make(map[string]*user.User),
	}
}

// AddUser adds a new user to the storage
func (s *Storage) AddUser(username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lowercaseUsername := strings.ToLower(username)
	if _, exists := s.users[lowercaseUsername]; exists {
		return errors.New("user already exists")
	}

	newUser, err := user.NewUser(username)
	if err != nil {
		return err
	}

	s.users[lowercaseUsername] = newUser
	return nil
}

// GetUser retrieves a user from the storage
func (s *Storage) GetUser(username string) (*user.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lowercaseUsername := strings.ToLower(username)
	user, exists := s.users[lowercaseUsername]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}