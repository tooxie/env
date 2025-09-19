package env

import (
	"testing"
)

func TestBoolValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid true values",
			value:       "true",
			expectError: false,
		},
		{
			name:        "valid false values",
			value:       "false",
			expectError: false,
		},
		{
			name:        "valid yes values",
			value:       "yes",
			expectError: false,
		},
		{
			name:        "valid no values",
			value:       "no",
			expectError: false,
		},
		{
			name:        "valid 1 values",
			value:       "1",
			expectError: false,
		},
		{
			name:        "valid 0 values",
			value:       "0",
			expectError: false,
		},
		{
			name:        "case sensitive - TRUE",
			value:       "TRUE",
			expectError: true,
		},
		{
			name:        "case sensitive - FALSE",
			value:       "FALSE",
			expectError: true,
		},
		{
			name:        "case sensitive - Yes",
			value:       "Yes",
			expectError: true,
		},
		{
			name:        "case sensitive - No",
			value:       "No",
			expectError: true,
		},
		{
			name:        "invalid value",
			value:       "maybe",
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "numeric 2",
			value:       "2",
			expectError: true,
		},
		{
			name:        "numeric -1",
			value:       "-1",
			expectError: true,
		},
		{
			name:        "whitespace",
			value:       " true ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := boolValidator(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}
		})
	}
}

func TestBoolParser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    bool
		expectError bool
	}{
		{
			name:        "parse true",
			value:       "true",
			expected:    true,
			expectError: false,
		},
		{
			name:        "parse false",
			value:       "false",
			expected:    false,
			expectError: false,
		},
		{
			name:        "parse yes",
			value:       "yes",
			expected:    true,
			expectError: false,
		},
		{
			name:        "parse no",
			value:       "no",
			expected:    false,
			expectError: false,
		},
		{
			name:        "parse 1",
			value:       "1",
			expected:    true,
			expectError: false,
		},
		{
			name:        "parse 0",
			value:       "0",
			expected:    false,
			expectError: false,
		},
		{
			name:        "parse invalid value",
			value:       "maybe",
			expected:    false,
			expectError: true,
		},
		{
			name:        "parse empty string",
			value:       "",
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := boolParser(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for value '%s'", tt.expected, result, tt.value)
			}
		})
	}
}

func TestIn(t *testing.T) {
	haystack := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		needle   string
		expected bool
	}{
		{
			name:     "found in middle",
			needle:   "banana",
			expected: true,
		},
		{
			name:     "found at beginning",
			needle:   "apple",
			expected: true,
		},
		{
			name:     "found at end",
			needle:   "cherry",
			expected: true,
		},
		{
			name:     "not found",
			needle:   "grape",
			expected: false,
		},
		{
			name:     "empty needle",
			needle:   "",
			expected: false,
		},
		{
			name:     "case sensitive",
			needle:   "Apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := in(tt.needle, haystack)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for needle '%s'", tt.expected, result, tt.needle)
			}
		})
	}
}
