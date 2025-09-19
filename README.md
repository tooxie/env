# env

A robust Go library for environment variable validation and parsing that prevents application startup unless all required environment variables are properly configured.

## Why Use This Library?

Environment variables are typically evaluated only when they're accessed, which means configuration errors can remain hidden until your application tries to use a specific variable. This can lead to:

- **Runtime failures** in production
- **Silent bugs** that surface only under specific conditions
- **Difficult debugging** when the application fails unexpectedly
- **Poor user experience** due to unclear error messages

This library solves these problems by **validating all environment variables at startup**, ensuring your application won't run unless it has everything it needs to function correctly.

## Features

- **Startup Validation**: Fail fast if required environment variables are missing or invalid
- **Type Safety**: Strong typing with generics for compile-time safety
- **Multiple Types**: Support for strings, integers, booleans, IPv4 addresses, and slices
- **Optional Fields**: Mark fields as optional with default values
- **Slice Support**: Parse comma-separated (or custom separator) lists
- **Custom Types**: Define your own types with validation logic
- **Comprehensive Error Messages**: Clear feedback about what's missing or invalid

## Installation

```bash
go get github.com/tooxie/env
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/tooxie/env"
)

type Config struct {
    DatabaseURL string `env:"required"`
    Port        int    `env:"optional,default='8080'"`
    Debug       bool   `env:"optional,default='false'"`
    AllowedIPs  []string `env:"optional,separator=',',default='127.0.0.1,192.168.1.1'"`
}

func main() {
    var config Config

    // Validate all environment variables at startup
    if err := env.Assert(config); err != nil {
        log.Fatalf("Configuration error: %v", err)
    }

    // Now you can safely access validated values
    dbURL := env.Get[string]("DatabaseURL")
    port := env.Get[int]("Port")
    debug := env.Get[bool]("Debug")
    ips := env.Get[[]string]("AllowedIPs")

    log.Printf("Starting server on port %d", port)
    log.Printf("Database: %s", dbURL)
    log.Printf("Debug mode: %v", debug)
    log.Printf("Allowed IPs: %v", ips)
}
```

## Supported Types

### Basic Types

| Type     | Example         | Valid Values                           |
|----------|-----------------|----------------------------------------|
| `string` | `"hello"`       | Any string                             |
| `int`    | `42`            | Valid integers                         |
| `bool`   | `true`          | `true`, `false`, `yes`, `no`, `1`, `0` |
| `IPv4`   | `"192.168.1.1"` | Valid IPv4 addresses                   |

### Slices

| Type       | Example                  | Separator       |
|------------|--------------------------|-----------------|
| `[]string` | `"a,b,c"`                | Comma (default) |
| `[]int`    | `"1,2,3"`                | Comma (default) |
| `[]bool`   | `"true,false"`           | Comma (default) |
| `[]IPv4`   | `"192.168.1.1,10.0.0.1"` | Comma (default) |

## Tag Options

### Required Fields
```go
type Config struct {
    DatabaseURL string `env:"required"`
}
```

### Optional Fields with Defaults
```go
type Config struct {
    Port  int    `env:"optional,default='8080'"`
    Debug bool   `env:"optional,default='false'"`
}
```

### Slice Fields with Custom Separators
```go
type Config struct {
    // Using pipe separator
    Hosts []string `env:"optional,separator='|',default='localhost|127.0.0.1'"`

    // Using semicolon separator
    Ports []int `env:"optional,separator=';',default='80;443'"`

    // Using space separator
    Features []bool `env:"optional,separator=' ',default='true false true'"`

    // Using custom character separator
    IPs []string `env:"optional,separator='#',default='192.168.1.1#10.0.0.1'"`
}
```

### Slice Examples with Different Separators

```go
// Environment variables:
// HOSTS="server1|server2|server3"
// PORTS="80;443;8080"
// FEATURES="true false true"
// ALLOWED_IPS="192.168.1.1#10.0.0.1#172.16.0.1"

type ServerConfig struct {
    Hosts      []string `env:"required,separator='|'"`
    Ports      []int    `env:"required,separator=';'"`
    Features   []bool   `env:"optional,separator=' ',default='true false'"`
    AllowedIPs []string `env:"optional,separator='#',default='127.0.0.1'"`
}

// Usage:
hosts := env.Get[[]string]("Hosts")      // ["server1", "server2", "server3"]
ports := env.Get[[]int]("Ports")         // [80, 443, 8080]
features := env.Get[[]bool]("Features")  // [true, false, true]
ips := env.Get[[]string]("AllowedIPs")   // ["192.168.1.1", "10.0.0.1", "172.16.0.1"]
```

## API Reference

### `env.Assert(config interface{}) error`

Validates all environment variables defined in the struct and returns an error if any are missing or invalid.

```go
var config MyConfig
if err := env.Assert(config); err != nil {
    log.Fatal(err)
}
```

### `env.Get[T](name string) T`

Retrieves a validated environment variable with type safety.

```go
dbURL := env.Get[string]("DatabaseURL")
port := env.Get[int]("Port")
debug := env.Get[bool]("Debug")
```

## Advanced Usage

### Custom Types

You can define custom types with their own validation logic:

```go
type IPv4 string

func (ip IPv4) Validate() error {
    // Custom validation logic
    return nil
}

type Config struct {
    ServerIP IPv4 `env:"required"`
}
```

### Complex Configuration

```go
type DatabaseConfig struct {
    Host     string `env:"required"`
    Port     int    `env:"required"`
    Username string `env:"required"`
    Password string `env:"required"`
    SSL      bool   `env:"optional,default='true'"`
}

type ServerConfig struct {
    ListenAddr string   `env:"optional,default='0.0.0.0'"`
    Port       int      `env:"optional,default='8080'"`
    Debug      bool     `env:"optional,default='false'"`
    AllowedIPs []string `env:"optional,separator=',',default='127.0.0.1'"`
}

type AppConfig struct {
    Database DatabaseConfig
    Server   ServerConfig
    AppName  string `env:"required"`
    Version  string `env:"optional,default='1.0.0'"`
}
```

## Error Handling

The library provides clear error messages for different failure scenarios:

```go
// Missing required field
// Error: Missing: [DatabaseURL, Port]

// Invalid field values
// Error: Invalid: [Port, Debug]

// Both missing and invalid
// Error: Missing: [DatabaseURL]
// Invalid: [Port]
```

## Running Tests

The library includes comprehensive tests covering all functionality:

```bash
# Run all tests
go test -v

# Run specific test packages
go test -v ./bool_test.go
go test -v ./int_test.go
go test -v ./ipv4_test.go

# Run with coverage
go test -v -cover

# Run integration tests only
go test -v -run="TestIntegration"
```

## Test Coverage

The test suite includes:

- **Unit Tests**: Individual parser and validation functions
- **Integration Tests**: Real-world usage scenarios
- **Error Handling**: Missing fields, invalid values, type mismatches
- **Edge Cases**: Boundary conditions and error scenarios
- **Type Safety**: Generic function usage and type conversion

## Best Practices

1. **Validate Early**: Call `env.Assert()` as early as possible in your application
2. **Use Required Fields**: Mark critical configuration as required
3. **Provide Defaults**: Use default values for optional configuration
4. **Handle Errors**: Always check the error returned by `env.Assert()`
5. **Type Safety**: Use the generic `Get[T]()` function for type safety

## Contributing

Contributions are welcome! Please ensure all tests pass before submitting a pull request:

```bash
go test -v
go fmt ./...
go vet ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Why This Approach?

Traditional environment variable handling often looks like this:

```go
// Silent failures, unclear errors
host := os.Getenv("HOST")
if host == "" {
    // This might not be caught until much later, and even if the environment
    // variable is set, there's no guarantee that the value is a valid IP
    log.Fatal("HOST not set")
}
```

With this library:

```go
// Good: Fail fast with clear errors
type Config struct {
    Host IPv4 `env:"required"`
    Port Int `env:"required"`
}

var config Config
if err := env.Assert(config); err != nil {
    log.Fatalf("Configuration error: %v", err)
}
// All environment variables are validated and ready to use
```

This ensures your application **won't start** unless it has everything it needs to run successfully, preventing runtime failures and improving reliability.
