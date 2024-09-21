package storage

import (
	"testing"
)

func TestStorage_AddUser(t *testing.T) {
	s := NewStorage()

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Valid user", "testuser", false},
		{"Duplicate user", "testuser", true},
		{"Invalid username", "invalid@user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.AddUser(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_GetUser(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Existing user", "testuser", false},
		{"Non-existent user", "nonexistent", true},
		{"Case-insensitive", "TestUser", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetUser(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}