package env

import (
	"testing"
)

func TestToLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "uppercase string",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "mixed case string",
			input:    "HeLLo",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with numbers",
			input:    "Hello123",
			expected: "hello123",
		},
		{
			name:     "string with special characters",
			input:    "Hello-World!",
			expected: "hello-world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toLower(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsOptional(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{
			name:     "contains optional",
			tag:      "optional",
			expected: true,
		},
		{
			name:     "contains OPTIONAL",
			tag:      "OPTIONAL",
			expected: true,
		},
		{
			name:     "contains Optional",
			tag:      "Optional",
			expected: true,
		},
		{
			name:     "contains optional with other tags",
			tag:      "optional,default='test'",
			expected: true,
		},
		{
			name:     "contains optional at end",
			tag:      "required,optional",
			expected: true,
		},
		{
			name:     "does not contain optional",
			tag:      "required",
			expected: false,
		},
		{
			name:     "empty tag",
			tag:      "",
			expected: false,
		},
		{
			name:     "contains optional in word",
			tag:      "optionality",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOptional(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}

func TestHasDefault(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{
			name:     "has default with single quotes",
			tag:      "default='test'",
			expected: true,
		},
		{
			name:     "has default with other tags",
			tag:      "optional,default='test',separator=','",
			expected: true,
		},
		{
			name:     "has default at beginning",
			tag:      "default='test',optional",
			expected: true,
		},
		{
			name:     "has default at end",
			tag:      "optional,default='test'",
			expected: true,
		},
		{
			name:     "no default",
			tag:      "optional",
			expected: false,
		},
		{
			name:     "empty tag",
			tag:      "",
			expected: false,
		},
		{
			name:     "default without quotes",
			tag:      "default=test",
			expected: false,
		},
		{
			name:     "default with double quotes",
			tag:      `default="test"`,
			expected: false,
		},
		{
			name:     "default with empty value",
			tag:      "default=''",
			expected: true,
		},
		{
			name:     "default with special characters",
			tag:      "default='test,value'",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasDefault(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}

func TestGetDefault(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{
			name:     "simple default",
			tag:      "default='test'",
			expected: "test",
		},
		{
			name:     "default with other tags",
			tag:      "optional,default='test',separator=','",
			expected: "test",
		},
		{
			name:     "default at beginning",
			tag:      "default='test',optional",
			expected: "test",
		},
		{
			name:     "default at end",
			tag:      "optional,default='test'",
			expected: "test",
		},
		{
			name:     "default with special characters",
			tag:      "default='test,value'",
			expected: "test,value",
		},
		{
			name:     "default with empty value",
			tag:      "default=''",
			expected: "",
		},
		{
			name:     "default with numbers",
			tag:      "default='123'",
			expected: "123",
		},
		{
			name:     "default with spaces",
			tag:      "default='hello world'",
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDefault(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s' for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}

func TestGetSeparator(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{
			name:     "comma separator",
			tag:      "separator=','",
			expected: ",",
		},
		{
			name:     "pipe separator",
			tag:      "separator='|'",
			expected: "|",
		},
		{
			name:     "semicolon separator",
			tag:      "separator=';'",
			expected: ";",
		},
		{
			name:     "space separator",
			tag:      "separator=' '",
			expected: " ",
		},
		{
			name:     "separator with other tags",
			tag:      "optional,separator=',',default='test'",
			expected: ",",
		},
		{
			name:     "no separator",
			tag:      "optional,default='test'",
			expected: " ",
		},
		{
			name:     "empty tag",
			tag:      "",
			expected: " ",
		},
		{
			name:     "separator with special character",
			tag:      "separator='\t'",
			expected: "\t",
		},
		{
			name:     "separator with newline",
			tag:      "separator='\n'",
			expected: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSeparator(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s' for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}

func TestGetSeparatorPanic(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{
			name: "multiple separators",
			tag:  "separator=',',separator='|'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for tag '%s' but didn't get one", tt.tag)
				}
			}()
			getSeparator(tt.tag)
		})
	}
}

func TestGetSeparatorInvalidFormat(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{
			name:     "separator without quotes",
			tag:      "separator=,",
			expected: " ",
		},
		{
			name:     "separator with double quotes",
			tag:      `separator=","`,
			expected: " ",
		},
		{
			name:     "separator with invalid format",
			tag:      "separator=invalid",
			expected: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSeparator(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s' for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name        string
		tag         string
		expected    string
		shouldPanic bool
	}{
		{
			name:     "no custom name specified",
			tag:      "required",
			expected: "",
		},
		{
			name:     "custom name with single quotes",
			tag:      "required,name='DATABASE_URL'",
			expected: "DATABASE_URL",
		},
		{
			name:     "custom name with other tags",
			tag:      "optional,default='8080',name='SERVER_PORT'",
			expected: "SERVER_PORT",
		},
		{
			name:     "custom name at beginning",
			tag:      "name='DEBUG_MODE',optional,default='false'",
			expected: "DEBUG_MODE",
		},
		{
			name:     "custom name at end",
			tag:      "required,name='HOST_NAME'",
			expected: "HOST_NAME",
		},
		{
			name:     "custom name with special characters",
			tag:      "required,name='API_KEY_123'",
			expected: "API_KEY_123",
		},
		{
			name:     "empty tag",
			tag:      "",
			expected: "",
		},
		{
			name:        "multiple name specifications",
			tag:         "name='FIRST',name='SECOND'",
			shouldPanic: true,
		},
		{
			name:     "name without quotes",
			tag:      "name=NO_QUOTES",
			expected: "",
		},
		{
			name:     "name with double quotes",
			tag:      `name="DOUBLE_QUOTES"`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but got none")
					}
				}()
			}

			result := getName(tt.tag)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s' for tag '%s'", tt.expected, result, tt.tag)
			}
		})
	}
}
