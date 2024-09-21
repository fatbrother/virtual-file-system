package validator

import (
	"testing"
)

func TestLengthValidator(t *testing.T) {
	tests := []struct {
		name   string
		min    int
		max    int
		input  string
		expect bool
	}{
		{"valid length", 3, 5, "test", true},
		{"too short", 3, 5, "go", false},
		{"too long", 3, 5, "testing", false},
		{"exact min length", 3, 5, "cat", true},
		{"exact max length", 3, 5, "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewLengthValidator(tt.min, tt.max)
			if result := validator.Validate(tt.input); result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestPatternValidator(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		expect  bool
	}{
		{"valid pattern", `^[a-z]+$`, "test", true},
		{"invalid pattern", `^[a-z]+$`, "Test123", false},
		{"empty pattern", `^$`, "", true},
		{"numeric pattern", `^\d+$`, "12345", true},
		{"alphanumeric pattern", `^[a-zA-Z0-9]+$`, "Test123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewPatternValidator(tt.pattern)
			if result := validator.Validate(tt.input); result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}