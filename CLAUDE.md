# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based rule engine that provides hierarchical rule processing with support for JSON, YAML, and custom data formats. The engine supports lazy evaluation, instant mode, and tree-based rule organization.

## Architecture

### Core Components

- **Forest**: Manages multiple rule trees and provides tree lifecycle management
- **Tree**: Hierarchical structure for organizing rules with path-based access
- **Rule**: Individual rule with path and processors
- **Driver**: Handles data format processing (JSON, YAML, Tile, etc.)
- **Processor**: Rule processing operations (JSON manipulation, HTTP calls, etc.)

### Key Interfaces

- `Forest`: Tree collection management with refresh capabilities
- `Tree`: Hierarchical rule storage with path-based operations
- `Rule`: Individual rule definition with processors
- `Driver`: Data format handling and path parsing
- `Processor`: Rule transformation operations

## Development Commands

### Building and Testing

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./driver

# Build the server binary
go build -o rule-server ./cmd/serve
```

### Running the Server

```bash
# Run the server with default configuration
go run ./cmd/serve

# Run with custom rules file
RULES_FILE=path/to/rules.json go run ./cmd/serve

# Run with custom shutdown timeout
SHUTDOWN_TIMEOUT=10s go run ./cmd/serve
```

### API Documentation

After starting the server, Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## Key Design Patterns

### Tree Modes

- **Lazy Mode**: Nodes are created and calculated only when accessed
- **Instant Mode**: Nodes are recalculated on every access for real-time data
- **Standard Mode**: Full tree built during initialization

### Driver System

Multiple driver implementations:
- `JSONDriver`: JSON data format processing
- `YAMLDriver`: YAML data format processing
- `TileDriver`: Custom tile-based processing
- `webDriver`: Web API integration driver

### Processor Types

- `JSONProcessor`: JSON data manipulation
- `YAMLProcessor`: YAML data manipulation
- `CURLProcessor`: HTTP request processing

## File Structure

```
├── cmd/serve/           # Server entry point
├── driver/              # Data format drivers and processors
├── web/                 # HTTP server and API handlers
├── rule.go              # Rule interface and implementation
├── tree.go              # Tree structure and operations
├── forest.go            # Forest management
├── factory.go           # Factory methods for tree creation
└── export.go            # Public API interfaces
```

## Configuration

### Environment Variables

- `RULES_FILE`: Path to rules configuration file (default: `../../conf/rules.json`)
- `SHUTDOWN_TIMEOUT`: Server shutdown timeout (default: `3s`)

### Rules File Format

Rules are defined in JSON format:
```json
[
  {
    "path": "/api/v1/users",
    "Processors": [
      {
        "type": "json",
        "data": {"operation": "merge", "value": {"enabled": true}}
      }
    ]
  }
]
```

## Testing Notes

- Tests include integration tests that make HTTP calls (may fail without network)
- Some tests require specific file paths that may not exist in all environments
- Test failures related to network connectivity or missing files are expected in isolated environments

## Common Development Tasks

### Adding New Processor Types

1. Implement the `driver.Processor` interface
2. Add processor type handling in `cmd/serve/serve.go:load()`
3. Update tests to cover the new processor

### Creating Custom Drivers

1. Implement `driver.Driver` interface
2. Add factory methods in `factory.go`
3. Update tree creation functions to support the new driver

### Extending Tree Functionality

- Modify `tree.go` for new tree operations
- Update `export.go` interfaces if API changes are needed
- Ensure thread safety with appropriate mutex usage