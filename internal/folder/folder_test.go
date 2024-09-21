package folder

import (
	"testing"
)

func TestNewFolder(t *testing.T) {
	tests := []struct {
		name        string
		folderName  string
		description string
		wantErr     bool
	}{
		{"Valid folder name", "validfolder", "A valid folder", false},
		{"Empty folder name", "", "Empty folder", true},
		{"Too long folder name", "thisfoldernameiswaytoolongandshouldfailvalidation_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "Too long folder name", true},
		{"Invalid characters", "invalid@folder", "Invalid folder name", true},
		{"Valid with numbers", "folder123", "Folder with numbers", false},
		{"Valid with underscore", "valid_folder", "Folder with underscore", false},
		{"Valid with hyphen", "valid-folder", "Folder with hyphen", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFolder(tt.folderName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
