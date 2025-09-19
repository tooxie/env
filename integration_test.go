package env

import (
	"os"
	"testing"
)

// Integration test structs
type DatabaseConfig struct {
	Host     string `env:"required"`
	Port     int    `env:"required"`
	Username string `env:"required"`
	Password string `env:"required"`
	SSL      bool   `env:"optional,default='true'"`
}

type ServerConfig struct {
	ListenAddr IPv4   `env:"optional,default='0.0.0.0'"`
	Port       int    `env:"optional,default='8080'"`
	Debug      bool   `env:"optional,default='false'"`
	LogLevel   string `env:"optional,default='info'"`
	AllowedIPs []IPv4 `env:"optional,separator=',',default='127.0.0.1,192.168.1.1'"`
	Ports      []int  `env:"optional,separator='|',default='80|443'"`
	Features   []bool `env:"optional,separator=';',default='true;false;true'"`
}

type ComplexConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	AppName  string `env:"required"`
	Version  string `env:"optional,default='1.0.0'"`
}

func TestIntegration_DatabaseConfig(t *testing.T) {
	// Set up environment variables
	envVars := map[string]string{
		"Host":     "localhost",
		"Port":     "5432",
		"Username": "admin",
		"Password": "secret",
		"SSL":      "true",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Test validation
	var config DatabaseConfig
	err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test getting values
	host := Get[string]("Host")
	if host != "localhost" {
		t.Errorf("Expected 'localhost', got '%s'", host)
	}

	port := Get[int]("Port")
	if port != 5432 {
		t.Errorf("Expected 5432, got %d", port)
	}

	username := Get[string]("Username")
	if username != "admin" {
		t.Errorf("Expected 'admin', got '%s'", username)
	}

	password := Get[string]("Password")
	if password != "secret" {
		t.Errorf("Expected 'secret', got '%s'", password)
	}

	ssl := Get[bool]("SSL")
	if ssl != true {
		t.Errorf("Expected true, got %v", ssl)
	}
}

func TestIntegration_ServerConfig(t *testing.T) {
	// Set up environment variables
	envVars := map[string]string{
		"ListenAddr": "192.168.1.100",
		"Port":       "9090",
		"Debug":      "true",
		"LogLevel":   "debug",
		"AllowedIPs": "192.168.1.1,192.168.1.2,10.0.0.1",
		"Ports":      "80|443|8080|8443",
		"Features":   "true;false;true;false",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Test validation
	var config ServerConfig
	err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test getting values
	listenAddr := Get[string]("ListenAddr")
	if listenAddr != "192.168.1.100" {
		t.Errorf("Expected '192.168.1.100', got '%s'", listenAddr)
	}

	port := Get[int]("Port")
	if port != 9090 {
		t.Errorf("Expected 9090, got %d", port)
	}

	debug := Get[bool]("Debug")
	if debug != true {
		t.Errorf("Expected true, got %v", debug)
	}

	logLevel := Get[string]("LogLevel")
	if logLevel != "debug" {
		t.Errorf("Expected 'debug', got '%s'", logLevel)
	}

	// Test slice values
	allowedIPs := Get[[]string]("AllowedIPs")
	expectedIPs := []string{"192.168.1.1", "192.168.1.2", "10.0.0.1"}
	if len(allowedIPs) != len(expectedIPs) {
		t.Errorf("Expected %d IPs, got %d", len(expectedIPs), len(allowedIPs))
	}
	for i, expected := range expectedIPs {
		if i < len(allowedIPs) && allowedIPs[i] != expected {
			t.Errorf("Expected IP %d to be '%s', got '%s'", i, expected, allowedIPs[i])
		}
	}

	ports := Get[[]int]("Ports")
	expectedPorts := []int{80, 443, 8080, 8443}
	if len(ports) != len(expectedPorts) {
		t.Errorf("Expected %d ports, got %d", len(expectedPorts), len(ports))
	}
	for i, expected := range expectedPorts {
		if i < len(ports) && ports[i] != expected {
			t.Errorf("Expected port %d to be %d, got %d", i, expected, ports[i])
		}
	}

	features := Get[[]bool]("Features")
	expectedFeatures := []bool{true, false, true, false}
	if len(features) != len(expectedFeatures) {
		t.Errorf("Expected %d features, got %d", len(expectedFeatures), len(features))
	}
	for i, expected := range expectedFeatures {
		if i < len(features) && features[i] != expected {
			t.Errorf("Expected feature %d to be %v, got %v", i, expected, features[i])
		}
	}
}

func TestIntegration_ServerConfigWithDefaults(t *testing.T) {
	// Test with minimal environment variables (should use defaults)
	envVars := map[string]string{
		// Only set one variable, others should use defaults
		"Debug": "true",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Test validation
	var config ServerConfig
	err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test default values
	listenAddr := Get[string]("ListenAddr")
	if listenAddr != "0.0.0.0" {
		t.Errorf("Expected default '0.0.0.0', got '%s'", listenAddr)
	}

	port := Get[int]("Port")
	if port != 8080 {
		t.Errorf("Expected default 8080, got %d", port)
	}

	debug := Get[bool]("Debug")
	if debug != true {
		t.Errorf("Expected true, got %v", debug)
	}

	logLevel := Get[string]("LogLevel")
	if logLevel != "info" {
		t.Errorf("Expected default 'info', got '%s'", logLevel)
	}
}

func TestIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		config      interface{}
		envVars     map[string]string
		expectError bool
		errorType   string
	}{
		{
			name:   "missing required field",
			config: DatabaseConfig{},
			envVars: map[string]string{
				"Host": "localhost",
				// Missing Port, Username, Password
			},
			expectError: true,
			errorType:   "missing",
		},
		{
			name:   "invalid field values",
			config: DatabaseConfig{},
			envVars: map[string]string{
				"Host":     "localhost",
				"Port":     "not-a-number",
				"Username": "admin",
				"Password": "secret",
				"SSL":      "maybe",
			},
			expectError: true,
			errorType:   "invalid",
		},
		{
			name:   "invalid IPv4 in slice",
			config: ServerConfig{},
			envVars: map[string]string{
				"ListenAddr": "192.168.1.1",
				"Port":       "8080",
				"Debug":      "true",
				"LogLevel":   "info",
				"AllowedIPs": "192.168.1.1,not-an-ip,10.0.0.1",
			},
			expectError: true,
			errorType:   "invalid",
		},
		{
			name:   "invalid int in slice",
			config: ServerConfig{},
			envVars: map[string]string{
				"ListenAddr": "192.168.1.1",
				"Port":       "8080",
				"Debug":      "true",
				"LogLevel":   "info",
				"Ports":      "80|not-a-number|443",
			},
			expectError: true,
			errorType:   "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Test validation
			err := Assert(tt.config)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if tt.expectError && err != nil {
				errorStr := err.Error()
				if tt.errorType == "missing" && !contains(errorStr, "Missing") {
					t.Errorf("Expected missing field error, got: %s", errorStr)
				}
				if tt.errorType == "invalid" && !contains(errorStr, "Invalid") {
					t.Errorf("Expected invalid field error, got: %s", errorStr)
				}
			}
		})
	}
}

func TestIntegration_ComplexConfig(t *testing.T) {
	// Set up environment variables for complex config
	envVars := map[string]string{
		"AppName": "myapp",
		"Version": "2.0.0",
		// Database config
		"Database.Host":     "db.example.com",
		"Database.Port":     "5432",
		"Database.Username": "dbuser",
		"Database.Password": "dbpass",
		"Database.SSL":      "true",
		// Server config
		"Server.ListenAddr": "0.0.0.0",
		"Server.Port":       "8080",
		"Server.Debug":      "false",
		"Server.LogLevel":   "warn",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Note: This test demonstrates the limitation of the current implementation
	// which doesn't handle nested structs. In a real implementation, you'd need
	// to handle nested struct field names properly.
	t.Skip("Skipping complex config test - nested structs not supported in current implementation")
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr))))
}
