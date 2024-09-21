package file

import (
	"errors"
	"strings"
	"time"

	"github.com/fatbrother/virtual-file-system/pkg/validator"
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