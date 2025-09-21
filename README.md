# env

A robust Go library for environment variable validation and parsing that prevents application startup unless all required environment variables are properly configured.

## Why Use This Library?

Environment variables are typically evaluated only when they're accessed, which means configuration errors can remain hidden until your application tries to use a specific variable. This can lead to:

- **Runtime failures** in production
- **Silent bugs** that surface only under specific conditions
- **Difficult debugging** when the application fails unexpectedly
- **Poor user experience** due to unclear error messages

Traditional environment variable handling often looks like this:

```go
// Silent failures, manual casting and validation, no compile-time checks.
host := os.Getenv("HOST")
if host == "" {
	// This might not be caught until much later, and even if the environment
	// variable is set, there's no guarantee that the value is a valid IP.
	log.Fatal("HOST not set")
}
```

With `env`:

```go
// Fail fast with clear errors. Values are automatically validated and casted to their corresponding types.
type EnvConfig struct {
	Host env.IPv4 `env:"required"`
	Port int      `env:"required"`
}

var envConfig EnvConfig
config := env.MustAssert(envConfig) // Panics if validation fails.

// All environment variables are validated and ready to use.
fmt.Printf("Server: %s:%d\n", config.Host, config.Port)
```

This ensures your application **won't start** unless it has everything it needs to run successfully, preventing runtime failures and improving reliability.

This library solves these problems by **validating all environment variables at startup**, ensuring your application won't run unless it has everything it needs to function correctly.

## API Quick Start
```go
// Recommended: MustAssert panics if validation fails - fail fast approach
config := env.MustAssert(envConfig)

// Alternative: Assert returns error for manual handling
config, err := env.Assert(envConfig)
if err != nil {
	log.Fatal(err)
}

// Direct field access - compile-time validation
fmt.Printf("Database: %s\n", config.DatabaseURL)
fmt.Printf("Server IP: %s\n", config.ServerIP)

// Custom types work seamlessly with validation
fmt.Printf("Port: %d\n", config.Port)

// Environment variable names use field names converted to uppercase by default
// `DatabaseURL` field will read the `DATABASEURL` environment variable
// Can be overridden with a custom name in the tag: `name='DB_URL'`
```

## Philosophy: Fail Fast, No Silent Defaults

**Production applications should never continue running with missing or invalid configuration.** This library enforces this principle by:

- **Failing immediately** if required environment variables are missing or invalid
- **Discouraging default values** in production environments
- **Providing clear error messages** about what's wrong

### Why Defaults Are Dangerous in Production

Default values can mask critical configuration issues and create a false sense of security:

- **Silent failures**: Your app might run with wrong settings without you knowing
- **Security risks**: Default credentials or insecure settings could be used
- **Data corruption**: Wrong database URLs or API endpoints could cause data loss
- **Debugging nightmares**: Issues only surface when the default value causes problems

**Best Practice**: Use `env.MustAssert()` and make all production configuration explicit and required.

## Features

- **Startup Validation**: Fail fast if required environment variables are missing or invalid
- **Type Safety**: Strong typing with generics for compile-time safety
- **Multiple Types**: Support for strings, integers, booleans, IPv4 addresses, and slices
- **Optional Fields**: Mark fields as optional with default values (use with caution in production)
- **Slice Support**: Parse comma-separated (or custom separator) lists
- **Custom Types**: Define your own types with validation logic
- **Comprehensive Error Messages**: Clear feedback about what's missing or invalid
- **Production Ready**: Designed for fail-fast configuration validation

## Installation

```bash
go get github.com/tooxie/env
```

## Supported Types

### Basic Types

| Type     | Example         | Valid Values                           |
|----------|-----------------|----------------------------------------|
| `string` | `"hello"`       | Any string                             |
| `int`    | `42`            | Valid integers                         |
| `bool`   | `true`          | `true`, `false`, `yes`, `no`, `1`, `0` |

### Custom Types

| Type   | Example         | Valid Values         |
|--------|-----------------|----------------------|
| `IPv4` | `"192.168.1.1"` | Valid IPv4 addresses |

### Slices

| Type       | Example                  | Separator              |
|------------|--------------------------|------------------------|
| `[]string` | `"a,b,c"`                | `","` (Comma, default) |
| `[]int`    | `"1 2 3"`                | `" "` (Space)          |
| `[]bool`   | `"true;false"`           | `";"` (Semicolon)      |
| `[]IPv4`   | `"192.168.1.1 10.0.0.1"` | `" "` (Space)          |

## Tag Options

| Option      | Description                                                                 | Example                         |
|-------------|-----------------------------------------------------------------------------|---------------------------------|
| `required`  | Field must have a value set (clashes with `optional`)                       | `env:"required"`                |
| `optional`  | Field can be empty, will use `nil` or `default` (clashes with `required`)   | `env:"optional"`                |
| `default`   | Default value if environment variable is not set (only valid with optional) | `env:"optional,default='8080'"` |
| `name`      | Custom environment variable name override                                   | `env:"name='DB_URL'"`           |
| `separator` | Custom separator for slice types (default is comma: `","`)                  | `env:"separator=' '"` (Space)   |

### Required Fields
```go
type EnvConfig struct {
	DatabaseURL string `env:"required"`
}
```

### Optional Fields with Defaults
```go
type EnvConfig struct {
	Port  int    `env:"optional,default='8080'"`
	Debug bool   `env:"optional,default='false'"`
}
```

### Environment Variable Names
By default, the field name is converted to uppercase for the environment variable name. You can override this behavior using the `name` field in the tag:

```go
type EnvConfig struct {
	// Uses field name in uppercase: DatabaseURL -> DATABASEURL
	DatabaseURL string `env:"required"`

	// Uses field name in uppercase: Port -> PORT
	Port        int    `env:"optional,default='8080'"`

	// Custom environment variable name: Debug -> DEBUG_MODE
	Debug       bool   `env:"optional,default='false',name='DEBUG_MODE'"`

	// Custom environment variable name: ServerPort -> SERVER_PORT
	ServerPort  int    `env:"required,name='SERVER_PORT'"`
}
```

### Slice Fields with Custom Separators
```go
type EnvConfig struct {
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

### Environment Variable Names Example

```go
// Environment:
//   * DATABASEURL="postgres://localhost:5432/mydb"
//   * SERVER_PORT="8080"
//   * DEBUG_MODE="true"
//   * ALLOWED_HOSTS="localhost|127.0.0.1|example.com"

type EnvConfig struct {
	DatabaseURL string     `env:"required"`                                    // DATABASEURL
	Port        int        `env:"optional,default='3000',name='SERVER_PORT'"`  // SERVER_PORT
	Debug       bool       `env:"optional,default='false',name='DEBUG_MODE'"`  // DEBUG_MODE
	Hosts       []string   `env:"required,separator='|',name='ALLOWED_HOSTS'"` // HOSTS
}

var envConfig EnvConfig
config, err := env.Assert(envConfig)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("DB_URL: %s\n", config.DatabaseURL) // postgres://localhost:5432/mydb (string)
fmt.Printf("Port: %d\n", config.Port)          // 8080 (int)
fmt.Printf("Debug: %v\n", config.Debug)        // true (bool)
fmt.Printf("Hosts: %v\n", config.Hosts)        // ["localhost", "127.0.0.1", "example.com"] ([]string)
```

### Slice Examples with Different Separators

```go
// Environment variables:
// HOSTS="server1|server2|server3"
// PORTS="80;443;8080"
// FEATURES="true false true"
// ALLOWED_IPS="192.168.1.1#10.0.0.1#172.16.0.1"

type ServerConfig struct {
	Hosts      []string   `env:"required,separator='|'"`
	Ports      []int      `env:"required,separator=';'"`
	Features   []bool     `env:"optional,separator=' ',default='true false'"`
	AllowedIPs []env.IPv4 `env:"optional,separator='#',default='127.0.0.1'"`
}

var serverConfig ServerConfig
config, err := env.Assert(serverConfig)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Hosts: %v\n", config.Hosts)      // ["server1", "server2", "server3"]
fmt.Printf("Ports: %v\n", config.Ports)      // [80, 443, 8080]
fmt.Printf("Features: %v\n", config.Features) // [true, false, true]
fmt.Printf("Allowed IPs: %v\n", config.AllowedIPs) // [IPv4("192.168.1.1"), IPv4("10.0.0.1"), IPv4("172.16.0.1")]
```

## API Reference

### `env.MustAssert[T](config T) T`

Validates all environment variables and returns a populated struct instance. **Panics if validation fails** - this is the recommended approach for production applications.

```go
var myConfig MyConfig
config := env.MustAssert(myConfig) // Panics if validation fails

// Direct field access - no type repetition, no string field names!
fmt.Printf("Host: %s\n", config.Host)
fmt.Printf("Port: %d\n", config.Port)
```

**Why MustAssert is recommended:**
- **Fail fast**: Application stops immediately if configuration is invalid
- **No silent failures**: Forces you to fix configuration issues before deployment
- **Cleaner code**: No need to handle errors manually
- **Production safety**: Prevents running with wrong configuration

### `env.Assert[T](config T) (T, error)`

Alternative function that returns an error instead of panicking. Use only when you need custom error handling.

```go
var myConfig MyConfig
config, err := env.Assert(myConfig)
if err != nil {
	log.Fatal(err) // You'll typically end up doing this anyway
}

// Direct field access
fmt.Printf("Host: %s\n", config.Host)
```

## Error Handling

The library provides clear error messages for different failure scenarios:

```go
// Missing required field
// Error: Missing: ["DatabaseURL (DATABASE_URL)", "Port (PORT)"]

// Invalid field values
// Error: Invalid: ["PORT (PORT)", "Debug (DEBUG)"]

// Both missing and invalid
// Error: Missing: ["DatabaseURL (DATABASE_URL)"]
// Invalid: ["PORT (PORT)"]
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

1. **Use MustAssert**: Prefer `env.MustAssert()` for production applications - fail fast approach
2. **Validate Early**: Call validation as early as possible in your application startup
3. **Make Everything Required**: Mark all production configuration as required - avoid optional values
4. **Explicit Configuration**: Never rely on default values in production environments
5. **Custom Types**: Use custom types like `IPv4` for validation instead of plain strings
6. **Direct Field Access**: Access fields directly from the returned struct instead of using string-based lookups
7. **Environment-Specific Configs**: Use different structs for development vs production if needed

### Production Configuration Philosophy

- **No Silent Defaults**: Every configuration value should be explicitly set
- **Fail Fast**: Application should not start if any required configuration is missing
- **Clear Errors**: Configuration errors should be obvious and actionable
- **Security First**: Default credentials or insecure settings are dangerous
- **Validate Credentials**: Don't just validate data format - verify that credentials work by testing connections (e.g. database, API keys) during startup to fail fast if they're invalid

## Contributing

Contributions are welcome! Please ensure all tests pass before submitting a pull request:

```bash
go test -v
go fmt ./...
go vet ./...
```
