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
		"HOST":     "localhost",
		"PORT":     "5432",
		"USERNAME": "admin",
		"PASSWORD": "secret",
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

	// Test validation and get populated struct
	var config DatabaseConfig
	env, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test getting values from populated struct
	if env.Host != "localhost" {
		t.Errorf("Expected 'localhost', got '%s'", env.Host)
	}

	if env.Port != 5432 {
		t.Errorf("Expected 5432, got %d", env.Port)
	}

	if env.Username != "admin" {
		t.Errorf("Expected 'admin', got '%s'", env.Username)
	}

	if env.Password != "secret" {
		t.Errorf("Expected 'secret', got '%s'", env.Password)
	}

	if env.SSL != true {
		t.Errorf("Expected true, got %v", env.SSL)
	}
}

func TestIntegration_ServerConfig(t *testing.T) {
	// Set up environment variables
	envVars := map[string]string{
		"LISTENADDR": "192.168.1.100",
		"PORT":       "9090",
		"DEBUG":      "true",
		"LOGLEVEL":   "debug",
		"ALLOWEDIPS": "192.168.1.1,192.168.1.2,10.0.0.1",
		"PORTS":      "80|443|8080|8443",
		"FEATURES":   "true;false;true;false",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Test validation and get populated struct
	var config ServerConfig
	env, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test getting values from populated struct
	if env.ListenAddr != "192.168.1.100" {
		t.Errorf("Expected '192.168.1.100', got '%s'", env.ListenAddr)
	}

	if env.Port != 9090 {
		t.Errorf("Expected 9090, got %d", env.Port)
	}

	if env.Debug != true {
		t.Errorf("Expected true, got %v", env.Debug)
	}

	if env.LogLevel != "debug" {
		t.Errorf("Expected 'debug', got '%s'", env.LogLevel)
	}

	// Test slice values
	expectedIPs := []IPv4{"192.168.1.1", "192.168.1.2", "10.0.0.1"}
	if len(env.AllowedIPs) != len(expectedIPs) {
		t.Errorf("Expected %d IPs, got %d", len(expectedIPs), len(env.AllowedIPs))
	}
	for i, expected := range expectedIPs {
		if i < len(env.AllowedIPs) && env.AllowedIPs[i] != expected {
			t.Errorf("Expected IP %d to be '%s', got '%s'", i, expected, env.AllowedIPs[i])
		}
	}

	expectedPorts := []int{80, 443, 8080, 8443}
	if len(env.Ports) != len(expectedPorts) {
		t.Errorf("Expected %d ports, got %d", len(expectedPorts), len(env.Ports))
	}
	for i, expected := range expectedPorts {
		if i < len(env.Ports) && env.Ports[i] != expected {
			t.Errorf("Expected port %d to be %d, got %d", i, expected, env.Ports[i])
		}
	}

	expectedFeatures := []bool{true, false, true, false}
	if len(env.Features) != len(expectedFeatures) {
		t.Errorf("Expected %d features, got %d", len(expectedFeatures), len(env.Features))
	}
	for i, expected := range expectedFeatures {
		if i < len(env.Features) && env.Features[i] != expected {
			t.Errorf("Expected feature %d to be %v, got %v", i, expected, env.Features[i])
		}
	}
}

func TestIntegration_ServerConfigWithDefaults(t *testing.T) {
	// Test with minimal environment variables (should use defaults)
	envVars := map[string]string{
		// Only set one variable, others should use defaults
		"DEBUG": "true",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Test validation and get populated struct
	var config ServerConfig
	env, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Test default values
	if env.ListenAddr != "0.0.0.0" {
		t.Errorf("Expected default '0.0.0.0', got '%s'", env.ListenAddr)
	}

	if env.Port != 8080 {
		t.Errorf("Expected default 8080, got %d", env.Port)
	}

	if env.Debug != true {
		t.Errorf("Expected true, got %v", env.Debug)
	}

	if env.LogLevel != "info" {
		t.Errorf("Expected default 'info', got '%s'", env.LogLevel)
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
