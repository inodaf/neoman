# AGENTS.md - Neoman Development Guide

Guidance for agentic coding tools (AI agents, code generators, etc.) working on this repository.

## Project Overview

**Neoman** is a modern documentation reader inspired by Unix `man` pages, written in Go. It consists of:
- `nman`: CLI client for reading documentation
- `nmand`: Daemon for serving documentation
- Core packages for Git integration, configuration, and browser operations

**Language**: Go 1.22+  
**Key Dependencies**: github.com/mattn/go-sqlite3

## Build, Lint & Test Commands

### Building

```bash
# Build both CLI and daemon
make all

# Build only CLI
go build -o ./bin/nman cmd/nman/main.go

# Build only daemon
go build -o ./bin/nmand cmd/daemon/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/management

# Run a single test
go test -run TestName ./path/to/package

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Linting & Formatting

```bash
# Format all Go files
go fmt ./...

# Check code with go vet (identifies suspicious constructs)
go vet ./...

# Gofmt (already used by go fmt)
gofmt -s -w .
```

## Code Style Guidelines

### Imports

- Organize imports in three groups: standard library, external packages, internal packages
- Use `goimport` style organization (automatic with `go fmt`)
- Example:
  ```go
  import (
      "fmt"
      "os"
      
      "github.com/mattn/go-sqlite3"
      
      "github.com/inodaf/neoman/pkg/config"
  )
  ```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `management`, `operations`, `config`)
- **Functions**: PascalCase for exported, camelCase for unexported (e.g., `GetDocs`, `openFile`)
- **Variables**: camelCase for unexported, PascalCase for exported constants
- **Constants**: PascalCase (e.g., `AppName`, `ShortAppName`, `PrimaryDocsDirName`)
- **Interfaces**: End with "er" suffix when describing behavior (e.g., convention from stdlib)
- **Errors**: Prefix with "Err" (e.g., `ErrGetWd`, `ErrReadDocsDir`)

### Formatting

- Standard Go indentation: tabs (enforced by gofmt)
- Line length: No hard limit, but keep readable (Go convention is flexible)
- Use meaningful comments for exported functions and packages
- Comment blocks should start with function/package name for godoc

### Types & Error Handling

- **Error Variables**: Define as package-level vars using `errors.New()`:
  ```go
  var (
      ErrGetWd = errors.New("could not get current working directory")
      ErrReadDocsDir = errors.New("no 'docs/' directory in this workspace")
  )
  ```
- **Error Wrapping**: Use `fmt.Errorf()` with context:
  ```go
  return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s'", owner, repo)
  ```
- **Error Checking**: Always handle errors explicitly, don't ignore them
- **Type Definitions**: Use structs for complex data; prefer value receivers for small types
- Example struct:
  ```go
  type RegistryEntry struct {
      Scope   string
      Owner   string
      Project string
  }
  ```

### Project Structure

- `/cmd`: CLI entry points (nman, daemon)
- `/internal`: Private packages (management, operations, models)
- `/pkg`: Public/reusable packages (config, browser, git)
- `/docs`: Documentation (Markdown files)

### Common Patterns

- **Configuration**: Use package-level constants in `pkg/config`
- **Error Messages**: Start with "neoman:" prefix for user-facing errors
- **Logging**: Use `fmt` for output; consider stderr for errors
- **System Integration**: Unix socket for IPC at `/tmp/nman.sock`

## Before Committing

1. Run `go fmt ./...` to format code
2. Run `go vet ./...` to check for issues
3. Run `go test ./...` to ensure tests pass
4. Verify builds: `make all`
5. Check that changes align with existing code style in relevant package
