package user

import (
	"testing"
	"strings"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Valid username", "validuser", false},
		{"Empty username", "", true},
		{"Too long username", "thisusernameiswaytoolongandshouldfailvalidation_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true},
		{"Invalid characters", "invalid@user", true},
		{"Valid with numbers", "user123", false},
		{"Valid with underscore", "valid_user", false},
		{"Valid with hyphen", "valid-user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Username != strings.ToLower(tt.username) {
				t.Errorf("NewUser() got = %v, want %v", got.Username, strings.ToLower(tt.username))
			}
		})
	}
}