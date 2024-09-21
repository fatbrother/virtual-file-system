package user

import (
	"errors"
	"strings"
	"time"

	"github.com/fatbrother/virtual-file-system/pkg/trie"
	"github.com/fatbrother/virtual-file-system/pkg/validator"
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
	validators := []validator.Validator{
		validator.NewLengthValidator(1, 50),
		validator.NewPatternValidator("^[a-zA-Z0-9_-]+$"),
	}

	for _, v := range validators {
		if pass := v.Validate(username); !pass {
			return errors.New("The " + username + " is invalid.")
		}
	}

	return nil
}
