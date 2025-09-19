package env

import (
	"testing"
)

func TestIPv4Validator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid IPv4 address",
			value:       "192.168.1.1",
			expectError: false,
		},
		{
			name:        "valid localhost",
			value:       "127.0.0.1",
			expectError: false,
		},
		{
			name:        "valid private IP",
			value:       "10.0.0.1",
			expectError: false,
		},
		{
			name:        "valid public IP",
			value:       "8.8.8.8",
			expectError: false,
		},
		{
			name:        "valid broadcast IP",
			value:       "255.255.255.255",
			expectError: false,
		},
		{
			name:        "valid zero IP",
			value:       "0.0.0.0",
			expectError: false,
		},
		{
			name:        "invalid IPv6 address",
			value:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expectError: true,
		},
		{
			name:        "invalid short IPv6",
			value:       "::1",
			expectError: true,
		},
		{
			name:        "invalid format - too many octets",
			value:       "192.168.1.1.1",
			expectError: true,
		},
		{
			name:        "invalid format - too few octets",
			value:       "192.168.1",
			expectError: true,
		},
		{
			name:        "invalid format - non-numeric",
			value:       "192.168.1.a",
			expectError: true,
		},
		{
			name:        "invalid format - empty octet",
			value:       "192.168..1",
			expectError: true,
		},
		{
			name:        "invalid format - leading zero",
			value:       "192.168.01.1",
			expectError: true,
		},
		{
			name:        "invalid range - octet too large",
			value:       "192.168.1.256",
			expectError: true,
		},
		{
			name:        "invalid range - negative octet",
			value:       "192.168.1.-1",
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "not an IP at all",
			value:       "not-an-ip",
			expectError: true,
		},
		{
			name:        "string with spaces",
			value:       " 192.168.1.1 ",
			expectError: true,
		},
		{
			name:        "string with newlines",
			value:       "192.168.1.1\n",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ipv4Validator(tt.value)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for value '%s' but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for value '%s' but got: %v", tt.value, err)
			}
		})
	}
}

func TestIPv4Parser(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid IPv4 address",
			value:       "192.168.1.1",
			expected:    "192.168.1.1",
			expectError: false,
		},
		{
			name:        "valid localhost",
			value:       "127.0.0.1",
			expected:    "127.0.0.1",
			expectError: false,
		},
		{
			name:        "invalid IPv6 address",
			value:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid format",
			value:       "192.168.1.256",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ipv4Parser(tt.value)

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
