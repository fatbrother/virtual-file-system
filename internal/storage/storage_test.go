package storage

import (
	"reflect"
	"testing"
	"sort"
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

func TestStorage_DeleteUser(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Existing user", "testuser", false},
		{"Case-insensitive", "TestUser", false},
		{"Non-existent user", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteUser(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		s.AddUser("testuser")
	}
}

func TestStorage_ListUsers(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("user1")
	_ = s.AddUser("user2")
	_ = s.AddUser("admin1")

	tests := []struct {
		name   string
		prefix string
		want   []string
	}{
		{"All users", "", []string{"admin1", "user1", "user2"}},
		{"User prefix", "user", []string{"user1", "user2"}},
		{"Admin prefix", "admin", []string{"admin1"}},
		{"Non-existent prefix", "nonexistent", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.ListUsers(tt.prefix)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}