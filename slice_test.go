package env

import (
	"strings"
	"testing"
)

func TestSliceValidation(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		fieldType   string
		value       string
		sep         string
		expectError bool
		description string
	}{
		// String slices
		{
			name:        "valid string slice with comma",
			fieldName:   "Hosts",
			fieldType:   "string",
			value:       "host1,host2,host3",
			sep:         ",",
			expectError: false,
			description: "Basic string slice with comma separator",
		},
		{
			name:        "valid string slice with pipe",
			fieldName:   "Servers",
			fieldType:   "string",
			value:       "server1|server2|server3",
			sep:         "|",
			expectError: false,
			description: "String slice with pipe separator",
		},
		{
			name:        "valid string slice with semicolon",
			fieldName:   "Domains",
			fieldType:   "string",
			value:       "example.com;test.com;demo.com",
			sep:         ";",
			expectError: false,
			description: "String slice with semicolon separator",
		},
		{
			name:        "valid string slice with space",
			fieldName:   "Features",
			fieldType:   "string",
			value:       "feature1 feature2 feature3",
			sep:         " ",
			expectError: false,
			description: "String slice with space separator",
		},
		{
			name:        "valid string slice with hash",
			fieldName:   "Tags",
			fieldType:   "string",
			value:       "tag1#tag2#tag3",
			sep:         "#",
			expectError: false,
			description: "String slice with hash separator",
		},
		{
			name:        "single element string slice",
			fieldName:   "SingleHost",
			fieldType:   "string",
			value:       "localhost",
			sep:         ",",
			expectError: false,
			description: "Single element string slice",
		},
		{
			name:        "empty string slice",
			fieldName:   "EmptyList",
			fieldType:   "string",
			value:       "",
			sep:         ",",
			expectError: false,
			description: "Empty string slice",
		},

		// Int slices
		{
			name:        "valid int slice",
			fieldName:   "Ports",
			fieldType:   "int",
			value:       "80,443,8080",
			sep:         ",",
			expectError: false,
			description: "Basic int slice",
		},
		{
			name:        "valid int slice with negative numbers",
			fieldName:   "Numbers",
			fieldType:   "int",
			value:       "1,-2,3,-4",
			sep:         ",",
			expectError: false,
			description: "Int slice with negative numbers",
		},
		{
			name:        "invalid int slice",
			fieldName:   "Ports",
			fieldType:   "int",
			value:       "80,not-a-number,8080",
			sep:         ",",
			expectError: true,
			description: "Int slice with invalid element",
		},

		// Bool slices
		{
			name:        "valid bool slice",
			fieldName:   "Flags",
			fieldType:   "bool",
			value:       "true,false,true",
			sep:         ",",
			expectError: false,
			description: "Basic bool slice",
		},
		{
			name:        "valid bool slice with various formats",
			fieldName:   "Options",
			fieldType:   "bool",
			value:       "true,false,yes,no,1,0",
			sep:         ",",
			expectError: false,
			description: "Bool slice with various valid formats",
		},
		{
			name:        "invalid bool slice",
			fieldName:   "Flags",
			fieldType:   "bool",
			value:       "true,maybe,false",
			sep:         ",",
			expectError: true,
			description: "Bool slice with invalid element",
		},

		// IPv4 slices
		{
			name:        "valid IPv4 slice",
			fieldName:   "AllowedIPs",
			fieldType:   "ipv4",
			value:       "192.168.1.1,10.0.0.1,172.16.0.1",
			sep:         ",",
			expectError: false,
			description: "Basic IPv4 slice",
		},
		{
			name:        "valid IPv4 slice with localhost",
			fieldName:   "LocalIPs",
			fieldType:   "ipv4",
			value:       "127.0.0.1,127.0.0.1,0.0.0.0",
			sep:         ",",
			expectError: false,
			description: "IPv4 slice with localhost IPs",
		},
		{
			name:        "invalid IPv4 slice",
			fieldName:   "AllowedIPs",
			fieldType:   "ipv4",
			value:       "192.168.1.1,not-an-ip,10.0.0.1",
			sep:         ",",
			expectError: true,
			description: "IPv4 slice with invalid element",
		},
		{
			name:        "invalid IPv4 slice with IPv6",
			fieldName:   "IPs",
			fieldType:   "ipv4",
			value:       "192.168.1.1,::1,10.0.0.1",
			sep:         ",",
			expectError: true,
			description: "IPv4 slice with IPv6 address",
		},

		// URL slices
		{
			name:        "valid URL slice",
			fieldName:   "ApiURLs",
			fieldType:   "url",
			value:       "https://api1.com,http://api2.com,ftp://files.com",
			sep:         ",",
			expectError: false,
			description: "Basic URL slice with various protocols",
		},
		{
			name:        "valid URL slice with ports",
			fieldName:   "Endpoints",
			fieldType:   "url",
			value:       "https://api.com:443,http://web.com:8080",
			sep:         ",",
			expectError: false,
			description: "URL slice with ports",
		},
		{
			name:        "valid URL slice without protocols",
			fieldName:   "Hosts",
			fieldType:   "url",
			value:       "api1.com:8080,api2.com:9090",
			sep:         ",",
			expectError: false,
			description: "URL slice without protocols (ParseRequestURI accepts these)",
		},
		{
			name:        "invalid URL slice",
			fieldName:   "URLs",
			fieldType:   "url",
			value:       "https://api.com,invalid-url,http://web.com",
			sep:         ",",
			expectError: true,
			description: "URL slice with invalid element",
		},
		{
			name:        "invalid URL slice with fragments",
			fieldName:   "URLs",
			fieldType:   "url",
			value:       "https://api.com,https://web.com#section",
			sep:         ",",
			expectError: true,
			description: "URL slice with fragment (not supported)",
		},

		// HTTPURL slices
		{
			name:        "valid HTTPURL slice",
			fieldName:   "WebEndpoints",
			fieldType:   "httpurl",
			value:       "https://api1.com,http://api2.com,https://web.com",
			sep:         ",",
			expectError: false,
			description: "Basic HTTPURL slice with HTTP/HTTPS only",
		},
		{
			name:        "valid HTTPURL slice with ports",
			fieldName:   "Services",
			fieldType:   "httpurl",
			value:       "https://api.com:443,http://web.com:8080",
			sep:         ",",
			expectError: false,
			description: "HTTPURL slice with ports",
		},
		{
			name:        "invalid HTTPURL slice with FTP",
			fieldName:   "Endpoints",
			fieldType:   "httpurl",
			value:       "https://api.com,ftp://files.com,http://web.com",
			sep:         ",",
			expectError: true,
			description: "HTTPURL slice with non-HTTP protocol",
		},
		{
			name:        "invalid HTTPURL slice with WS",
			fieldName:   "Endpoints",
			fieldType:   "httpurl",
			value:       "https://api.com,ws://socket.com,http://web.com",
			sep:         ",",
			expectError: true,
			description: "HTTPURL slice with WebSocket protocol",
		},
		{
			name:        "invalid HTTPURL slice without protocols",
			fieldName:   "Endpoints",
			fieldType:   "httpurl",
			value:       "https://api.com,api.com,http://web.com",
			sep:         ",",
			expectError: true,
			description: "HTTPURL slice with element without protocol",
		},

		// Edge cases
		{
			name:        "slice with empty elements",
			fieldName:   "Items",
			fieldType:   "string",
			value:       "item1,,item3",
			sep:         ",",
			expectError: false,
			description: "Slice with empty elements (should be handled gracefully)",
		},
		{
			name:        "slice with whitespace elements",
			fieldName:   "Items",
			fieldType:   "string",
			value:       " item1 , item2 , item3 ",
			sep:         ",",
			expectError: false,
			description: "Slice with whitespace around elements",
		},
		{
			name:        "slice with only separators",
			fieldName:   "Empty",
			fieldType:   "string",
			value:       ",,,",
			sep:         ",",
			expectError: false,
			description: "Slice with only separators",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateAndParseSlice(tt.fieldName, tt.fieldType, tt.value, tt.sep)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s but got none", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s but got: %v", tt.description, err)
			}

			// For valid cases, verify we got some result
			if !tt.expectError && err == nil {
				if result == nil {
					t.Errorf("Expected non-nil result for %s", tt.description)
				}
			}
		})
	}
}

func TestSliceWithCustomSeparators(t *testing.T) {
	tests := []struct {
		name        string
		fieldType   string
		value       string
		sep         string
		expectError bool
	}{
		{
			name:        "URL slice with pipe separator",
			fieldType:   "url",
			value:       "https://api1.com|https://api2.com|http://web.com",
			sep:         "|",
			expectError: false,
		},
		{
			name:        "HTTPURL slice with semicolon separator",
			fieldType:   "httpurl",
			value:       "https://api1.com;https://api2.com;http://web.com",
			sep:         ";",
			expectError: false,
		},
		{
			name:        "IPv4 slice with hash separator",
			fieldType:   "ipv4",
			value:       "192.168.1.1#10.0.0.1#172.16.0.1",
			sep:         "#",
			expectError: false,
		},
		{
			name:        "Int slice with space separator",
			fieldType:   "int",
			value:       "80 443 8080 8443",
			sep:         " ",
			expectError: false,
		},
		{
			name:        "Bool slice with custom separator",
			fieldType:   "bool",
			value:       "true@false@true",
			sep:         "@",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateAndParseSlice("TestField", tt.fieldType, tt.value, tt.sep)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError && err == nil && result == nil {
				t.Errorf("Expected non-nil result")
			}
		})
	}
}

func TestSliceErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		fieldType     string
		value         string
		sep           string
		expectError   bool
		errorContains string
	}{
		{
			name:          "mixed valid and invalid ints",
			fieldType:     "int",
			value:         "1,2,invalid,4",
			sep:           ",",
			expectError:   true,
			errorContains: "invalid slice",
		},
		{
			name:          "mixed valid and invalid bools",
			fieldType:     "bool",
			value:         "true,false,maybe,true",
			sep:           ",",
			expectError:   true,
			errorContains: "invalid slice",
		},
		{
			name:          "mixed valid and invalid IPv4s",
			fieldType:     "ipv4",
			value:         "192.168.1.1,not-an-ip,10.0.0.1",
			sep:           ",",
			expectError:   true,
			errorContains: "invalid slice",
		},
		{
			name:          "mixed valid and invalid URLs",
			fieldType:     "url",
			value:         "https://api.com,invalid-url,http://web.com",
			sep:           ",",
			expectError:   true,
			errorContains: "invalid slice",
		},
		{
			name:          "mixed valid and invalid HTTPURLs",
			fieldType:     "httpurl",
			value:         "https://api.com,ftp://files.com,http://web.com",
			sep:           ",",
			expectError:   true,
			errorContains: "invalid slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateAndParseSlice("TestField", tt.fieldType, tt.value, tt.sep)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result == nil {
					t.Errorf("Expected non-nil result")
				}
			}
		})
	}
}
