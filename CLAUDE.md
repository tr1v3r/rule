# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based hierarchical content-construction engine (ivy) that supports JSON, YAML, XML, TOML, and custom data formats. The engine supports lazy evaluation, instant mode, cache TTL, and tree-based content organization.

## Architecture

### Core Components

- **Forest**: Manages multiple trees and provides tree lifecycle management
- **Tree**: Hierarchical structure with path-based access and layered content transformation
- **Directive**: Path + processors pair that defines a transformation on the tree
- **Driver**: Handles data format processing (JSON, YAML, XML, TOML, Tile, etc.)
- **Processor**: Content transformation operations (JSON manipulation, HTTP calls, etc.)

### Key Interfaces

- `Forest`: Tree collection management with refresh capabilities
- `Tree`: Hierarchical content storage with path-based operations
- `Directive`: Path and processor definition for content transformation
- `Driver`: Data format handling and path parsing
- `Processor`: Content transformation operations

## Development Commands

### Building and Testing

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./driver

# Build the server binary
go build -o ivy-server ./cmd/serve
```

### Running the Server

```bash
# Run the server with default configuration
go run ./cmd/serve

# Run with custom directives file
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
- **Cache TTL Mode**: Lazy with time-based cache expiration
- **Standard Mode**: Full tree built during initialization

### Driver System

Multiple driver implementations:
- `JSONDriver`: JSON data format processing
- `YAMLDriver`: YAML data format processing
- `XMLDriver`: XML data format processing
- `TOMLDriver`: TOML data format processing
- `TileDriver`: Custom tile-based processing
- `webDriver`: Web API integration driver

### Processor Types

- `JSONProcessor`: JSON data manipulation
- `YAMLProcessor`: YAML data manipulation
- `XMLProcessor`: XML data manipulation
- `TOMLProcessor`: TOML data manipulation
- `CURLProcessor`: HTTP request processing
- `RawProcessor`: Custom transformation function
- `CombinedProcessor`: Combines multiple processors

## File Structure

```
├── cmd/serve/           # Server entry point
├── driver/              # Data format drivers and processors
├── web/                 # HTTP server and API handlers
├── directive.go         # Directive interface and implementation
├── tree.go              # Tree structure and operations
├── forest.go            # Forest management
├── factory.go           # Factory methods for tree creation
└── export.go            # Public API interfaces
```

## Configuration

### Environment Variables

- `RULES_FILE`: Path to directives configuration file (default: `../../conf/rules.json`)
- `SHUTDOWN_TIMEOUT`: Server shutdown timeout (default: `3s`)

### Directives File Format

Directives are defined in JSON format:
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
