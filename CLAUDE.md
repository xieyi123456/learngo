# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go learning project (`learngo`) that implements a JSON difference comparison service. The project provides functionality to compare two JSON objects and identify additions, removals, and changes between them.

## Project Structure

```
learngo/
├── go.mod                          # Go module file (Go 1.18)
├── main.go                         # Main entry point with demonstration
└── service/
    ├── json_diff_service.go        # Core JSON diff service implementation
    └── json_diff_service_test.go   # Comprehensive test suite
```

## Architecture

### Core Components

1. **JSONDiffService Interface** (`service/json_diff_service.go:17-19`)
   - Defines the contract for JSON comparison functionality
   - Single method: `CompareJSON(json1, json2 string) (JSONDiffResult, error)`

2. **JSONDiffResult Structure** (`service/json_diff_service.go:10-14`)
   - `Added []string`: Paths of fields added in the second JSON
   - `Removed []string`: Paths of fields removed from the first JSON  
   - `Changed map[string]string`: Paths and descriptions of changed values

3. **Implementation** (`service/json_diff_service.go:22`)
   - `jsonDiffServiceImpl`: Concrete implementation of the service
   - Handles nested objects, arrays, and primitive types recursively
   - Supports complex path notation (e.g., `person.address.city`, `hobbies[1]`)

### Key Features

- Deep comparison of nested JSON objects and arrays
- Type change detection (e.g., string to number)
- Null value handling
- Array length and element comparison
- Detailed path-based change reporting

## Commands

### Development Commands

```bash
# Run the main application
go run main.go

# Build the application  
go build

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./service

# Run a specific test
go test ./service -run TestCompareJSON_IdenticalJSON

# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Clean build cache
go clean
```

### Module Management

```bash
# Initialize module (already done)
go mod init learngo

# Download dependencies
go mod download

# Clean up unused dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Testing

The project includes comprehensive tests (`service/json_diff_service_test.go`) covering:

- Identical JSON comparison
- Value differences
- Field additions and removals  
- Nested object comparison
- Array comparison
- Invalid JSON handling
- Complex multi-level nested structures

Run `go test -v ./service` to execute the full test suite with detailed output.

## Development Notes

- The main application demonstrates the service with sample JSON objects
- All text output uses Chinese language for user-facing messages
- The service uses Go's `reflect` package for deep value comparison
- Path building uses dot notation for nested objects and bracket notation for arrays
- Error handling includes descriptive Chinese error messages