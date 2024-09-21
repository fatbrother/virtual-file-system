package folder

import (
	"errors"
	"strings"
	"time"

	"github.com/fatbrother/virtual-file-system/pkg/trie"
	"github.com/fatbrother/virtual-file-system/pkg/validator"
)

// Folder represents a folder in the virtual file system
type Folder struct {
	Name        string
	Description string
	CreatedAt   time.Time
	Files       *trie.Trie
}

// NewFolder creates a new Folder with the given name and description
func NewFolder(name, description string) (*Folder, error) {
	if err := validateFolderName(name); err != nil {
		return nil, err
	}
	return &Folder{
		Name:        strings.ToLower(name),
		Description: description,
		CreatedAt:   time.Now(),
		Files:       trie.NewTrie(),
	}, nil
}

// validateFolderName checks if the folder name is valid
func validateFolderName(name string) error {
	validators := []validator.Validator{
		validator.NewLengthValidator(1, 50),
		validator.NewPatternValidator("^[a-zA-Z0-9_-]+$"),
	}

	for _, v := range validators {
		if pass := v.Validate(name); !pass {
			return errors.New("The " + name + " is invalid.")
		}
	}

	return nil
}