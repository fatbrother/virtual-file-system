package folder

import (
	"errors"
	"strings"
	"time"

	"github.com/fatbrother/virtual-file-system/pkg/trie"
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
	if len(name) == 0 {
		return errors.New("folder name cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("folder name cannot be longer than 50 characters")
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' || char == '.') {
			return errors.New("folder name can only contain alphanumeric characters, underscores, hyphens, and dots")
		}
	}
	return nil
}