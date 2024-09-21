package file

import (
	"errors"
	"strings"
	"time"
)

// File represents a file in the virtual file system
type File struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

// NewFile creates a new File instance
func NewFile(name, description string) (*File, error) {
	if err := ValidateFileName(name); err != nil {
		return nil, err
	}

	return &File{
		Name:        strings.ToLower(name),
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

// ValidateFileName checks if the file name is valid
func ValidateFileName(name string) error {
	if len(name) == 0 {
		return errors.New("file name cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("file name cannot be longer than 50 characters")
	}
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' || char == '.') {
			return errors.New("file name can only contain alphanumeric characters, underscores, hyphens, and dots")
		}
	}
	return nil
}