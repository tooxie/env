package env

import (
	"os"
	"testing"
)

// Test structs for values functionality
type TestConfigWithValues struct {
	Port        int    `env:"values='8000,8080,9000',default='8080'"`
	Environment string `env:"values='dev,staging,prod'"`
	Debug       bool   `env:"values='true,false'"`
}

type TestConfigWithValuesOptional struct {
	Port        int    `env:"optional,values='8000,8080,9000',default='8080'"`
	Environment string `env:"optional,values='dev,staging,prod',default='dev'"`
	Debug       bool   `env:"optional,values='true,false',default='false'"`
}

type TestConfigWithValuesAndDefault struct {
	Port        int    `env:"optional,values='8000,8080,9000',default='8080'"`
	Environment string `env:"optional,values='dev,staging,prod',default='dev'"`
	Debug       bool   `env:"optional,values='true,false',default='false'"`
}

func TestValuesValidation(t *testing.T) {
	tests := []struct {
		name            string
		config          interface{}
		envVars         map[string]string
		expectedMissing []string
		expectedInvalid []string
		shouldPanic     bool
	}{
		{
			name:   "valid values for all fields",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{},
		},
		{
			name:   "invalid port value",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "3000", // Not in allowed values
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{"Port (PORT)"},
		},
		{
			name:   "invalid environment value",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "production", // Not in allowed values
				"DEBUG":       "true",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{"Environment (ENVIRONMENT)"},
		},
		{
			name:   "invalid debug value",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "staging",
				"DEBUG":       "maybe", // Not in allowed values
			},
			expectedMissing: []string{},
			expectedInvalid: []string{"Debug (DEBUG)"},
		},
		{
			name:   "missing required field with values",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT": "8080",
				// Missing ENVIRONMENT and DEBUG
			},
			expectedMissing: []string{"Environment (ENVIRONMENT)", "Debug (DEBUG)"},
			expectedInvalid: []string{},
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

			missing, invalid := Validate(tt.config)

			// Check missing variables
			if len(missing) != len(tt.expectedMissing) {
				t.Errorf("Expected %d missing variables, got %d: %v", len(tt.expectedMissing), len(missing), missing)
			} else {
				for i, expected := range tt.expectedMissing {
					if i >= len(missing) || missing[i] != expected {
						t.Errorf("Expected missing[%d] to be '%s', got '%s'", i, expected, missing[i])
					}
				}
			}

			// Check invalid variables
			if len(invalid) != len(tt.expectedInvalid) {
				t.Errorf("Expected %d invalid variables, got %d: %v", len(tt.expectedInvalid), len(invalid), invalid)
			} else {
				for i, expected := range tt.expectedInvalid {
					if i >= len(invalid) || invalid[i].name != expected {
						t.Errorf("Expected invalid[%d].name to be '%s', got '%s'", i, expected, invalid[i].name)
					}
				}
			}
		})
	}
}

func TestValuesWithOptionalFields(t *testing.T) {
	tests := []struct {
		name            string
		config          interface{}
		envVars         map[string]string
		expectedMissing []string
		expectedInvalid []string
	}{
		{
			name:   "all optional fields with valid values",
			config: TestConfigWithValuesOptional{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectedMissing: []string{},
			expectedInvalid: []string{},
		},
		{
			name:            "all optional fields with no values (should use defaults)",
			config:          TestConfigWithValuesOptional{},
			envVars:         map[string]string{},
			expectedMissing: []string{},
			expectedInvalid: []string{},
		},
		{
			name:   "some optional fields with invalid values",
			config: TestConfigWithValuesOptional{},
			envVars: map[string]string{
				"PORT":        "3000", // Invalid
				"ENVIRONMENT": "staging",
				"DEBUG":       "maybe", // Invalid
			},
			expectedMissing: []string{},
			expectedInvalid: []string{"Port (PORT)", "Debug (DEBUG)"},
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

			missing, invalid := Validate(tt.config)

			// Check missing variables
			if len(missing) != len(tt.expectedMissing) {
				t.Errorf("Expected %d missing variables, got %d: %v", len(tt.expectedMissing), len(missing), missing)
			}

			// Check invalid variables
			if len(invalid) != len(tt.expectedInvalid) {
				t.Errorf("Expected %d invalid variables, got %d: %v", len(tt.expectedInvalid), len(invalid), invalid)
			}
		})
	}
}

func TestValuesWithDefaultValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      interface{}
		shouldPanic bool
		panicMsg    string
	}{
		{
			name:        "valid default values",
			config:      TestConfigWithValuesAndDefault{},
			shouldPanic: false,
		},
		{
			name: "invalid default value should panic",
			config: struct {
				Port int `env:"optional,values='8000,8080,9000',default='3000'"` // 3000 not in allowed values
			}{},
			shouldPanic: true,
			panicMsg:    "Default value '3000' for field 'Port' is not in allowed values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but got none")
					} else if panicMsg, ok := r.(string); !ok || panicMsg[:len(tt.panicMsg)] != tt.panicMsg {
						t.Errorf("Expected panic message to contain '%s', got '%s'", tt.panicMsg, panicMsg)
					}
				}()
			}

			// This should not panic for valid defaults, or panic for invalid defaults
			_, _ = Validate(tt.config)
		})
	}
}

func TestValuesAssertIntegration(t *testing.T) {
	tests := []struct {
		name      string
		config    interface{}
		envVars   map[string]string
		expectErr bool
	}{
		{
			name:   "successful assertion with valid values",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectErr: false,
		},
		{
			name:   "failed assertion with invalid values",
			config: TestConfigWithValues{},
			envVars: map[string]string{
				"PORT":        "3000", // Invalid
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectErr: true,
		},
		{
			name:   "successful assertion with optional fields",
			config: TestConfigWithValuesOptional{},
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "staging",
				"DEBUG":       "true",
			},
			expectErr: false,
		},
		{
			name:      "successful assertion with no optional fields (use defaults)",
			config:    TestConfigWithValuesOptional{},
			envVars:   map[string]string{},
			expectErr: false,
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

			_, err := Assert(tt.config)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValuesWithDefaults(t *testing.T) {
	// Test that optional fields with defaults work correctly
	config := TestConfigWithValuesOptional{}

	// No environment variables set - should use defaults
	result, err := Assert(config)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Type assert to access fields
	cfg := result

	// Check that default values are set correctly
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %v", cfg.Port)
	}
	if cfg.Environment != "dev" {
		t.Errorf("Expected Environment to be 'dev', got %v", cfg.Environment)
	}
	if cfg.Debug != false {
		t.Errorf("Expected Debug to be false, got %v", cfg.Debug)
	}
}
