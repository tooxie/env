package env

import (
	"testing"
)

func TestStringParser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:        "normal string",
			value:       "hello world",
			expected:    "hello world",
			expectError: false,
		},
		{
			name:        "empty string",
			value:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "string with special characters",
			value:       "!@#$%^&*()",
			expected:    "!@#$%^&*()",
			expectError: false,
		},
		{
			name:        "string with numbers",
			value:       "123abc456",
			expected:    "123abc456",
			expectError: false,
		},
		{
			name:        "string with whitespace",
			value:       "  hello  world  ",
			expected:    "  hello  world  ",
			expectError: false,
		},
		{
			name:        "unicode string",
			value:       "你好世界",
			expected:    "你好世界",
			expectError: false,
		},
		{
			name:        "string with newlines",
			value:       "line1\nline2\nline3",
			expected:    "line1\nline2\nline3",
			expectError: false,
		},
		{
			name:        "string with tabs",
			value:       "col1\tcol2\tcol3",
			expected:    "col1\tcol2\tcol3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := stringParser(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
