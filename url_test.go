package env

import (
	"testing"
)

func TestURLValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		// Valid URLs with protocol
		{
			name:        "valid HTTP URL",
			value:       "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			value:       "https://example.com",
			expectError: false,
		},
		{
			name:        "valid URL with port",
			value:       "http://example.com:8080",
			expectError: false,
		},
		{
			name:        "valid URL with path",
			value:       "https://example.com/path/to/resource",
			expectError: false,
		},
		{
			name:        "valid URL with query parameters",
			value:       "https://example.com?param1=value1&param2=value2",
			expectError: false,
		},
		{
			name:        "valid URL with subdomain",
			value:       "https://api.example.com",
			expectError: false,
		},
		{
			name:        "valid URL with IP address",
			value:       "http://192.168.1.1",
			expectError: false,
		},
		{
			name:        "valid URL with IP address and port",
			value:       "http://192.168.1.1:8080",
			expectError: false,
		},
		{
			name:        "valid URL with localhost",
			value:       "http://localhost",
			expectError: false,
		},
		{
			name:        "valid URL with localhost and port",
			value:       "http://localhost:3000",
			expectError: false,
		},
		{
			name:        "valid URL with FTP protocol",
			value:       "ftp://example.com",
			expectError: false,
		},
		{
			name:        "valid URL with custom protocol",
			value:       "custom://example.com",
			expectError: false,
		},
		{
			name:        "invalid URL without protocol",
			value:       "example.com",
			expectError: true,
		},
		{
			name:        "valid URL without protocol with port",
			value:       "example.com:8080",
			expectError: false,
		},
		{
			name:        "invalid URL without protocol with path",
			value:       "example.com/path",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with IP",
			value:       "192.168.1.1",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with localhost",
			value:       "localhost",
			expectError: true,
		},
		// Invalid URLs
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			value:       "   ",
			expectError: true,
		},
		{
			name:        "invalid format - just protocol",
			value:       "http://",
			expectError: false,
		},
		{
			name:        "invalid format - protocol with no host",
			value:       "https://",
			expectError: false,
		},
		{
			name:        "invalid format - just colon",
			value:       ":",
			expectError: true,
		},
		{
			name:        "invalid format - just slash",
			value:       "/",
			expectError: false,
		},
		{
			name:        "invalid format - malformed protocol",
			value:       "htp://example.com",
			expectError: false,
		},
		{
			name:        "invalid format - spaces in host",
			value:       "http://exam ple.com",
			expectError: true,
		},
		{
			name:        "invalid format - invalid characters",
			value:       "http://exam[ple.com",
			expectError: false,
		},
		{
			name:        "invalid format - multiple colons",
			value:       "http://example.com:8080:9090",
			expectError: false,
		},
		{
			name:        "invalid format - port without host",
			value:       "http://:8080",
			expectError: false,
		},
		{
			name:        "invalid format - just port",
			value:       ":8080",
			expectError: true,
		},
		{
			name:        "invalid format - host with whitespace",
			value:       " example.com ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := urlValidator(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}
		})
	}
}

func TestURLParser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid HTTP URL",
			value:       "http://example.com",
			expected:    "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			value:       "https://example.com",
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "invalid URL without protocol",
			value:       "example.com",
			expected:    "",
			expectError: true,
		},
		{
			name:        "valid URL without protocol with port",
			value:       "example.com:8080",
			expected:    "example.com:8080",
			expectError: false,
		},
		{
			name:        "invalid URL without protocol with path",
			value:       "example.com/path",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid URL with whitespace",
			value:       "  example.com  ",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid URL with IP address",
			value:       "192.168.1.1",
			expected:    "",
			expectError: true,
		},
		{
			name:        "valid URL with localhost",
			value:       "localhost:3000",
			expected:    "localhost:3000",
			expectError: false,
		},
		{
			name:        "invalid empty URL",
			value:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "valid protocol only",
			value:       "http://",
			expected:    "http://",
			expectError: false,
		},
		{
			name:        "valid malformed URL",
			value:       "htp://example.com",
			expected:    "htp://example.com",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := urlParser(tt.value)

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

func TestHTTPURLValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		// Valid HTTP/HTTPS URLs
		{
			name:        "valid HTTP URL",
			value:       "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			value:       "https://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTP URL with port",
			value:       "http://example.com:8080",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL with port",
			value:       "https://example.com:443",
			expectError: false,
		},
		{
			name:        "valid HTTP URL with path",
			value:       "http://example.com/path",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL with query parameters",
			value:       "https://example.com?param1=value1&param2=value2",
			expectError: false,
		},
		{
			name:        "valid HTTP URL with subdomain",
			value:       "http://api.example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL with IP address",
			value:       "https://192.168.1.1",
			expectError: false,
		},
		{
			name:        "valid HTTP URL with localhost",
			value:       "http://localhost:3000",
			expectError: false,
		},
		// Invalid URLs (HTTPURL is stricter than URL)
		{
			name:        "invalid URL without protocol",
			value:       "example.com",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with port",
			value:       "example.com:8080",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with path",
			value:       "example.com/path",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with IP",
			value:       "192.168.1.1",
			expectError: true,
		},
		{
			name:        "invalid URL without protocol with localhost",
			value:       "localhost",
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			value:       "   ",
			expectError: true,
		},
		// Invalid protocols for HTTPURL
		{
			name:        "invalid FTP protocol",
			value:       "ftp://example.com",
			expectError: true,
		},
		{
			name:        "invalid WS protocol",
			value:       "ws://example.com",
			expectError: true,
		},
		{
			name:        "invalid WSS protocol",
			value:       "wss://example.com",
			expectError: true,
		},
		{
			name:        "invalid FILE protocol",
			value:       "file:///path/to/file",
			expectError: true,
		},
		{
			name:        "invalid DATA protocol",
			value:       "data:text/plain,Hello World",
			expectError: true,
		},
		{
			name:        "invalid MAILTO protocol",
			value:       "mailto:test@example.com",
			expectError: true,
		},
		{
			name:        "invalid SSH protocol",
			value:       "ssh://user@example.com",
			expectError: true,
		},
		{
			name:        "invalid custom protocol",
			value:       "custom://example.com",
			expectError: true,
		},
		// Other invalid formats (same as URL type)
		{
			name:        "invalid format - just colon",
			value:       ":",
			expectError: true,
		},
		{
			name:        "invalid format - just port",
			value:       ":8080",
			expectError: true,
		},
		{
			name:        "invalid format - host with whitespace",
			value:       " example.com ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := httpURLValidator(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}
		})
	}
}

func TestHTTPURLParser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid HTTP URL",
			value:       "http://example.com",
			expected:    "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			value:       "https://example.com",
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTP URL with port",
			value:       "http://example.com:8080",
			expected:    "http://example.com:8080",
			expectError: false,
		},
		{
			name:        "invalid URL without protocol with port",
			value:       "example.com:8080",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid URL with localhost",
			value:       "localhost:3000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid empty URL",
			value:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid FTP protocol",
			value:       "ftp://example.com",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid WS protocol",
			value:       "ws://example.com",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid custom protocol",
			value:       "custom://example.com",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := httpURLParser(tt.value)

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
