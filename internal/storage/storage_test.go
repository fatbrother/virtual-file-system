package storage

import (
	"reflect"
	"testing"
	"time"
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

func TestStorage_CreateFolder(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")

	tests := []struct {
		name        string
		username    string
		folderName  string
		description string
		wantErr     bool
	}{
		{"Valid folder", "testuser", "documents", "My documents", false},
		{"Duplicate folder", "testuser", "documents", "Another description", true},
		{"Invalid folder name", "testuser", "invalid/name", "Invalid name", true},
		{"Non-existent user", "nonexistent", "folder", "Description", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateFolder(tt.username, tt.folderName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteFolder(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")
	_ = s.CreateFolder("testuser", "documents", "My documents")

	tests := []struct {
		name       string
		username   string
		folderName string
		wantErr    bool
	}{
		{"Existing folder", "testuser", "documents", false},
		{"Non-existent folder", "testuser", "nonexistent", true},
		{"Non-existent user", "nonexistent", "folder", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteFolder(tt.username, tt.folderName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteFolder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_ListFolders(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")
	_ = s.CreateFolder("testuser", "documents", "My documents")
	time.Sleep(1 * time.Second) // 確保創建時間不同
	_ = s.CreateFolder("testuser", "pictures", "My pictures")

	tests := []struct {
		name      string
		username  string
		sortField string
		sortOrder string
		want      []string
		wantErr   bool
	}{
		{"Sort by name asc", "testuser", "name", "asc", []string{"documents", "pictures"}, false},
		{"Sort by name desc", "testuser", "name", "desc", []string{"pictures", "documents"}, false},
		{"Sort by created asc", "testuser", "created", "asc", []string{"documents", "pictures"}, false},
		{"Sort by created desc", "testuser", "created", "desc", []string{"pictures", "documents"}, false},
		{"Non-existent user", "nonexistent", "name", "asc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListFolders(tt.username, tt.sortField, tt.sortOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.ListFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_CreateFile(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")
	_ = s.CreateFolder("testuser", "documents", "My documents")

	tests := []struct {
		name        string
		username    string
		folderName  string
		fileName    string
		description string
		wantErr     bool
	}{
		{"Valid file", "testuser", "documents", "file1.txt", "File description", false},
		{"Duplicate file", "testuser", "documents", "file1.txt", "Another description", true},
		{"Invalid file name", "testuser", "documents", "invalid/file", "Invalid name", true},
		{"Non-existent folder", "testuser", "nonexistent", "file.txt", "Description", true},
		{"Non-existent user", "nonexistent", "documents", "file.txt", "Description", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateFile(tt.username, tt.folderName, tt.fileName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteFile(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")
	_ = s.CreateFolder("testuser", "documents", "My documents")
	_ = s.CreateFile("testuser", "documents", "file1.txt", "File description")

	tests := []struct {
		name       string
		username   string
		folderName string
		fileName   string
		wantErr    bool
	}{
		{"Existing file", "testuser", "documents", "file1.txt", false},
		{"Non-existent file", "testuser", "documents", "nonexistent.txt", true},
		{"Non-existent folder", "testuser", "nonexistent", "file.txt", true},
		{"Non-existent user", "nonexistent", "documents", "file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeleteFile(tt.username, tt.folderName, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_ListFiles(t *testing.T) {
	s := NewStorage()
	_ = s.AddUser("testuser")
	_ = s.CreateFolder("testuser", "documents", "My documents")
	_ = s.CreateFile("testuser", "documents", "file1.txt", "File description")
	time.Sleep(1 * time.Second) // 確保創建時間不同
	_ = s.CreateFile("testuser", "documents", "file2.txt", "Another file description")

	tests := []struct {
		name       string
		username   string
		folderName string
		sortField  string
		sortOrder  string
		want       []string
		wantErr    bool
	}{
		{"Sort by name asc", "testuser", "documents", "name", "asc", []string{"file1.txt", "file2.txt"}, false},
		{"Sort by name desc", "testuser", "documents", "name", "desc", []string{"file2.txt", "file1.txt"}, false},
		{"Sort by created asc", "testuser", "documents", "created", "asc", []string{"file1.txt", "file2.txt"}, false},
		{"Sort by created desc", "testuser", "documents", "created", "desc", []string{"file2.txt", "file1.txt"}, false},
		{"Non-existent folder", "testuser", "nonexistent", "name", "asc", nil, true},
		{"Non-existent user", "nonexistent", "documents", "name", "asc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListFiles(tt.username, tt.folderName, tt.sortField, tt.sortOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.ListFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}