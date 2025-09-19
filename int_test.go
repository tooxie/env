package env

import (
	"testing"
)

func TestIntParser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    int
		expectError bool
	}{
		{
			name:        "positive integer",
			value:       "42",
			expected:    42,
			expectError: false,
		},
		{
			name:        "negative integer",
			value:       "-42",
			expected:    -42,
			expectError: false,
		},
		{
			name:        "zero",
			value:       "0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "large positive integer",
			value:       "2147483647",
			expected:    2147483647,
			expectError: false,
		},
		{
			name:        "large negative integer",
			value:       "-2147483648",
			expected:    -2147483648,
			expectError: false,
		},
		{
			name:        "invalid string",
			value:       "not-a-number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "float number",
			value:       "3.14",
			expected:    0,
			expectError: true,
		},
		{
			name:        "hexadecimal",
			value:       "0xFF",
			expected:    0,
			expectError: true,
		},
		{
			name:        "binary",
			value:       "0b1010",
			expected:    0,
			expectError: true,
		},
		{
			name:        "octal",
			value:       "0777",
			expected:    777,
			expectError: false,
		},
		{
			name:        "string with leading zeros",
			value:       "007",
			expected:    7,
			expectError: false,
		},
		{
			name:        "string with spaces",
			value:       " 42 ",
			expected:    0,
			expectError: true,
		},
		{
			name:        "string with plus sign",
			value:       "+42",
			expected:    42,
			expectError: false,
		},
		{
			name:        "very large number",
			value:       "999999999999999999999999999999",
			expected:    9223372036854775807, // Max int64 (strconv.Atoi returns this even on error)
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := intParser(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d for value '%s'", tt.expected, result, tt.value)
			}
		})
	}
}
