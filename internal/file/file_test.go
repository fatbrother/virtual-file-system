package file

import (
	"testing"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		name		string
		fileName	string
		description	string
		wantErr		bool
	}{
		{"Valid file name", "validfile.txt", "A valid file", false},
		{"Empty file name", "", "Empty file", true},
		{"Too long file name", "thisfilenameiswaytoolongandshouldfailvalidation_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.txt", "Too long file name", true},
		{"Invalid characters", "**&&^%$#@!", "Invalid file name", true},
		{"Valid with numbers", "file123", "File with numbers", false},
		{"Valid with underscore", "valid_file", "File with underscore", false},
		{"Valid with hyphen", "valid-file", "File with hyphen", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFile(tt.fileName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
