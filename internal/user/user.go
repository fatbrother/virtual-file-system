package user

import (
	"errors"
	"strings"
	"time"

	"github.com/fatbrother/virtual-file-system/pkg/trie"
)

// User represents a user in the virtual file system
type User struct {
	Username  string
	CreatedAt time.Time
	Folders   *trie.Trie
}

// NewUser creates a new User with the given username
func NewUser(username string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	return &User{
		Username:  strings.ToLower(username),
		CreatedAt: time.Now(),
		Folders:   trie.NewTrie(),
	}, nil
}

// validateUsername checks if the username is valid
func validateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("username cannot be empty")
	}
	if len(username) > 50 {
		return errors.New("username cannot be longer than 50 characters")
	}
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-') {
			return errors.New("username can only contain alphanumeric characters, underscores, and hyphens")
		}
	}
	return nil
}