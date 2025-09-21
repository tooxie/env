package env

import (
	"os"
	"testing"
)

func TestCustomEnvVarName(t *testing.T) {
	// Set up environment variables with custom names
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/mydb")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DEBUG_MODE", "true")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DEBUG_MODE")
	}()

	// Test struct with custom environment variable names
	type Config struct {
		DatabaseURL string `env:"required,name='DATABASE_URL'"`
		Port        int    `env:"optional,default='3000',name='SERVER_PORT'"`
		Debug       bool   `env:"optional,default='false',name='DEBUG_MODE'"`
	}

	// Test validation
	var config Config
	envConfig, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Verify the values are correctly populated
	if envConfig.DatabaseURL != "postgres://localhost:5432/mydb" {
		t.Errorf("Expected DatabaseURL 'postgres://localhost:5432/mydb', got '%s'", envConfig.DatabaseURL)
	}

	if envConfig.Port != 8080 {
		t.Errorf("Expected Port 8080, got %d", envConfig.Port)
	}

	if envConfig.Debug != true {
		t.Errorf("Expected Debug true, got %v", envConfig.Debug)
	}
}

func TestCustomEnvVarNameMissing(t *testing.T) {
	// Test with missing required environment variable
	type Config struct {
		DatabaseURL string `env:"required,name='MISSING_DATABASE_URL'"`
	}

	var config Config
	_, err := Assert(config)
	if err == nil {
		t.Fatal("Expected error for missing environment variable, but got none")
	}

	// Check that the error message contains the environment variable name, not the field name
	errorStr := err.Error()
	if !contains(errorStr, "MISSING_DATABASE_URL") {
		t.Errorf("Expected error to contain 'MISSING_DATABASE_URL', got: %s", errorStr)
	}
}

func TestCustomEnvVarNameInvalid(t *testing.T) {
	// Set up environment variable with invalid value
	os.Setenv("INVALID_PORT", "not-a-number")
	defer os.Unsetenv("INVALID_PORT")

	// Test struct with custom environment variable name
	type Config struct {
		Port int `env:"required,name='INVALID_PORT'"`
	}

	var config Config
	_, err := Assert(config)
	if err == nil {
		t.Fatal("Expected error for invalid environment variable, but got none")
	}

	// Check that the error message contains the environment variable name, not the field name
	errorStr := err.Error()
	if !contains(errorStr, "INVALID_PORT") {
		t.Errorf("Expected error to contain 'INVALID_PORT', got: %s", errorStr)
	}
}

func TestCustomEnvVarNameWithSlices(t *testing.T) {
	// Set up environment variables with custom names
	os.Setenv("ALLOWED_HOSTS", "localhost|127.0.0.1|example.com")
	os.Setenv("SERVER_PORTS", "80,443,8080")
	defer func() {
		os.Unsetenv("ALLOWED_HOSTS")
		os.Unsetenv("SERVER_PORTS")
	}()

	// Test struct with custom environment variable names for slices
	type Config struct {
		Hosts []string `env:"required,separator='|',name='ALLOWED_HOSTS'"`
		Ports []int    `env:"required,separator=',',name='SERVER_PORTS'"`
	}

	var config Config
	envConfig, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Verify the slice values are correctly populated
	expectedHosts := []string{"localhost", "127.0.0.1", "example.com"}
	if len(envConfig.Hosts) != len(expectedHosts) {
		t.Errorf("Expected %d hosts, got %d", len(expectedHosts), len(envConfig.Hosts))
	}
	for i, expected := range expectedHosts {
		if i < len(envConfig.Hosts) && envConfig.Hosts[i] != expected {
			t.Errorf("Expected host %d to be '%s', got '%s'", i, expected, envConfig.Hosts[i])
		}
	}

	expectedPorts := []int{80, 443, 8080}
	if len(envConfig.Ports) != len(expectedPorts) {
		t.Errorf("Expected %d ports, got %d", len(expectedPorts), len(envConfig.Ports))
	}
	for i, expected := range expectedPorts {
		if i < len(envConfig.Ports) && envConfig.Ports[i] != expected {
			t.Errorf("Expected port %d to be %d, got %d", i, expected, envConfig.Ports[i])
		}
	}
}

func TestCustomEnvVarNameWithCustomTypes(t *testing.T) {
	// Set up environment variables with custom names
	os.Setenv("SERVER_IP", "192.168.1.1")
	os.Setenv("ALLOWED_IPS", "10.0.0.1,172.16.0.1")
	defer func() {
		os.Unsetenv("SERVER_IP")
		os.Unsetenv("ALLOWED_IPS")
	}()

	// Test struct with custom environment variable names for custom types
	type Config struct {
		ServerIP   IPv4   `env:"required,name='SERVER_IP'"`
		AllowedIPs []IPv4 `env:"required,separator=',',name='ALLOWED_IPS'"`
	}

	var config Config
	envConfig, err := Assert(config)
	if err != nil {
		t.Fatalf("Assert failed: %v", err)
	}

	// Verify the custom type values are correctly populated
	if envConfig.ServerIP != "192.168.1.1" {
		t.Errorf("Expected ServerIP '192.168.1.1', got '%s'", envConfig.ServerIP)
	}

	expectedIPs := []IPv4{"10.0.0.1", "172.16.0.1"}
	if len(envConfig.AllowedIPs) != len(expectedIPs) {
		t.Errorf("Expected %d IPs, got %d", len(expectedIPs), len(envConfig.AllowedIPs))
	}
	for i, expected := range expectedIPs {
		if i < len(envConfig.AllowedIPs) && envConfig.AllowedIPs[i] != expected {
			t.Errorf("Expected IP %d to be '%s', got '%s'", i, expected, envConfig.AllowedIPs[i])
		}
	}
}
