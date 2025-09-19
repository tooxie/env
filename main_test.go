package env

import (
	"os"
	"strings"
	"testing"
)

// Test structs for validation
type TestConfig struct {
	DatabaseURL string `env:"required"`
	Port        int    `env:"optional,default='8080'"`
	Debug       bool   `env:"optional,default='true'"`
	IPAddress   IPv4   `env:"optional,default='127.0.0.1'"`
}

type TestConfigWithSlices struct {
	Hosts   []string `env:"optional,separator=',',default='localhost,127.0.0.1'"`
	Ports   []int    `env:"optional,separator='|',default='80|443'"`
	Enabled []bool   `env:"optional,separator=';',default='true;false'"`
}

type TestConfigAllRequired struct {
	RequiredString string `env:"required"`
	RequiredInt    int    `env:"required"`
	RequiredBool   bool   `env:"required"`
	RequiredIPv4   IPv4   `env:"required"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name            string
		config          interface{}
		envVars         map[string]string
		expectedMissing []string
		expectedInvalid []string
		shouldPanic     bool
	}{
		{
			name:   "valid config with all required fields",
			config: TestConfigAllRequired{},
			envVars: map[string]string{
				"RequiredString": "test",
				"RequiredInt":    "42",
				"RequiredBool":   "true",
				"RequiredIPv4":   "192.168.1.1",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{},
		},
		{
			name:   "missing required fields",
			config: TestConfigAllRequired{},
			envVars: map[string]string{
				"RequiredString": "test",
				// Missing RequiredInt, RequiredBool, RequiredIPv4
			},
			expectedMissing: []string{"RequiredInt", "RequiredBool", "RequiredIPv4"},
			expectedInvalid: []string{},
		},
		{
			name:   "invalid field values",
			config: TestConfigAllRequired{},
			envVars: map[string]string{
				"RequiredString": "test",
				"RequiredInt":    "not-a-number",
				"RequiredBool":   "maybe",
				"RequiredIPv4":   "not-an-ip",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{"RequiredInt", "RequiredBool", "RequiredIPv4"},
		},
		{
			name:   "valid config with optional fields and defaults",
			config: TestConfig{},
			envVars: map[string]string{
				"DatabaseURL": "postgres://localhost:5432/mydb",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{},
		},
		{
			name:        "invalid parameter type",
			config:      "not a struct",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but didn't get one")
					}
				}()
			}

			missing, invalid := Validate(tt.config)

			if tt.shouldPanic {
				return
			}

			// Check missing fields
			if len(missing) != len(tt.expectedMissing) {
				t.Errorf("Expected %d missing fields, got %d", len(tt.expectedMissing), len(missing))
			}

			for _, expected := range tt.expectedMissing {
				found := false
				for _, actual := range missing {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected missing field '%s' not found in %v", expected, missing)
				}
			}

			// Check invalid fields
			if len(invalid) != len(tt.expectedInvalid) {
				t.Errorf("Expected %d invalid fields, got %d", len(tt.expectedInvalid), len(invalid))
			}

			for _, expected := range tt.expectedInvalid {
				found := false
				for _, actual := range invalid {
					if actual.name == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected invalid field '%s' not found in %v", expected, invalid)
				}
			}
		})
	}
}

func TestAssert(t *testing.T) {
	tests := []struct {
		name        string
		config      interface{}
		envVars     map[string]string
		expectError bool
	}{
		{
			name:   "valid config",
			config: TestConfig{},
			envVars: map[string]string{
				"DatabaseURL": "postgres://localhost:5432/mydb",
			},
			expectError: false,
		},
		{
			name:   "missing required field",
			config: TestConfigAllRequired{},
			envVars: map[string]string{
				"RequiredString": "test",
				// Missing other required fields
			},
			expectError: true,
		},
		{
			name:   "invalid field value",
			config: TestConfigAllRequired{},
			envVars: map[string]string{
				"RequiredString": "test",
				"RequiredInt":    "not-a-number",
				"RequiredBool":   "true",
				"RequiredIPv4":   "192.168.1.1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			err := Assert(tt.config)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Set up a valid environment
	os.Setenv("TestString", "hello")
	os.Setenv("TestInt", "42")
	os.Setenv("TestBool", "true")
	os.Setenv("TestIPv4", "192.168.1.1")
	defer func() {
		os.Unsetenv("TestString")
		os.Unsetenv("TestInt")
		os.Unsetenv("TestBool")
		os.Unsetenv("TestIPv4")
	}()

	// Validate first to populate envMap
	config := struct {
		TestString string `env:"required"`
		TestInt    int    `env:"required"`
		TestBool   bool   `env:"required"`
		TestIPv4   IPv4   `env:"required"`
	}{}

	missing, invalid := Validate(config)
	if len(missing) > 0 || len(invalid) > 0 {
		t.Fatalf("Validation failed: missing=%v, invalid=%v", missing, invalid)
	}

	tests := []struct {
		name     string
		key      string
		expected interface{}
		panic    bool
	}{
		{
			name:     "get string",
			key:      "TestString",
			expected: "hello",
		},
		{
			name:     "get int",
			key:      "TestInt",
			expected: 42,
		},
		{
			name:     "get bool",
			key:      "TestBool",
			expected: true,
		},
		{
			name:     "get IPv4",
			key:      "TestIPv4",
			expected: "192.168.1.1",
		},
		{
			name:  "get non-existent key",
			key:   "NonExistentKey",
			panic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but didn't get one")
					}
				}()
				// Call Get directly for panic test
				Get[string](tt.key)
				return
			}

			switch expected := tt.expected.(type) {
			case string:
				result := Get[string](tt.key)
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			case int:
				result := Get[int](tt.key)
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			case bool:
				result := Get[bool](tt.key)
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestValidateAndParseSlice(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		fieldType   string
		value       string
		sep         string
		expectError bool
	}{
		{
			name:        "valid string slice",
			fieldName:   "Hosts",
			fieldType:   "String",
			value:       "host1,host2,host3",
			sep:         ",",
			expectError: false,
		},
		{
			name:        "valid int slice",
			fieldName:   "Ports",
			fieldType:   "Int",
			value:       "80|443|8080",
			sep:         "|",
			expectError: false,
		},
		{
			name:        "valid bool slice",
			fieldName:   "Enabled",
			fieldType:   "Bool",
			value:       "true;false;true",
			sep:         ";",
			expectError: false,
		},
		{
			name:        "invalid int slice",
			fieldName:   "Ports",
			fieldType:   "Int",
			value:       "80|not-a-number|8080",
			sep:         "|",
			expectError: true,
		},
		{
			name:        "invalid bool slice",
			fieldName:   "Enabled",
			fieldType:   "Bool",
			value:       "true;maybe;false",
			sep:         ";",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateAndParseSlice(tt.fieldName, tt.fieldType, tt.value, tt.sep)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError {
				expectedCount := len(strings.Split(tt.value, tt.sep))
				if len(result) != expectedCount {
					t.Errorf("Expected %d items, got %d", expectedCount, len(result))
				}
			}
		})
	}
}

func TestParseVariable(t *testing.T) {
	tests := []struct {
		name        string
		fieldName   string
		fieldType   string
		value       string
		expectError bool
	}{
		{
			name:        "valid bool",
			fieldName:   "Debug",
			fieldType:   "Bool",
			value:       "true",
			expectError: false,
		},
		{
			name:        "valid string",
			fieldName:   "Name",
			fieldType:   "String",
			value:       "test",
			expectError: false,
		},
		{
			name:        "valid int",
			fieldName:   "Port",
			fieldType:   "Int",
			value:       "8080",
			expectError: false,
		},
		{
			name:        "valid IPv4",
			fieldName:   "IP",
			fieldType:   "IPv4",
			value:       "192.168.1.1",
			expectError: false,
		},
		{
			name:        "invalid bool",
			fieldName:   "Debug",
			fieldType:   "Bool",
			value:       "maybe",
			expectError: true,
		},
		{
			name:        "invalid int",
			fieldName:   "Port",
			fieldType:   "Int",
			value:       "not-a-number",
			expectError: true,
		},
		{
			name:        "invalid IPv4",
			fieldName:   "IP",
			fieldType:   "IPv4",
			value:       "not-an-ip",
			expectError: true,
		},
		{
			name:        "unknown type",
			fieldName:   "Unknown",
			fieldType:   "Unknown",
			value:       "test",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fieldType == "Unknown" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic for unknown type but didn't get one")
					}
				}()
			}

			_, err := parseVariable(tt.fieldName, tt.fieldType, tt.value)

			if tt.fieldType == "Unknown" {
				return
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
