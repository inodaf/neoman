# Coding Conventions Overview

This document describes the general coding principles and standards for the Neoman codebase. For layer-specific patterns, see the dedicated convention files.

## Architecture Principles

### Dependency Flow

Dependencies flow **inward** toward the domain layer:

```
Controller → UseCase → Domain
     ↓          ↓
  Repository (interfaces only)
```

- **Domain** has no external dependencies
- **UseCase** depends on Domain and Repository interfaces (not implementations)
- **Controller** depends on UseCase
- **Repository** package contains only interfaces, no implementations

This ensures business logic is independent of infrastructure and easy to test.

### Layer Responsibilities

| Layer | Responsibility |
|-------|----------------|
| Domain | Business entities, validation, core rules |
| UseCase | Business logic orchestration, input/output contracts |
| Controller | HTTP handling, request/response translation |
| Repository | Data access interfaces (ports) |
| Worker | Background/async operations |
| Command (App) | CLI commands, daemon communication |

## Naming Conventions

### Packages

- Lowercase, single word: `domain`, `usecase`, `controller`, `repo`, `worker`, `command`

### Functions and Types

- **Exported**: PascalCase (e.g., `NewUseCase`, `DocsRepository`)
- **Unexported**: camelCase (e.g., `parseAuthorAndRepo`, `formatResponse`)

### Variables

- **Exported constants**: PascalCase (e.g., `MaxFileSize`)
- **Unexported variables**: camelCase (e.g., `docsRepository`)

### Errors

- Prefix with `Err` followed by operation context
- Pattern: `Err<Operation><SpecificFailure>`
- Example: `ErrAddDocsAuthorRequired`, `ErrListPagesNotFound`

### Files

- Lowercase with no separators: `addremotedocs.go`, `listpages.go`
- Test files: `<filename>_test.go`

## Import Organization

Organize imports in three groups, separated by blank lines:

```go
import (
    // 1. Standard library
    "fmt"
    "net/http"

    // 2. External packages
    "github.com/external/package"

    // 3. Internal packages
    "github.com/inodaf/neoman/internal/daemon/domain"
    "github.com/inodaf/neoman/pkg/config"
)
```

## Error Handling

### Package-Level Error Variables

Define errors as package-level variables for type-safe comparison:

```go
var (
    Err<Operation><Failure1> = fmt.Errorf("<description>")
    Err<Operation><Failure2> = fmt.Errorf("<description>")
)
```

### Error Wrapping

When adding context to errors, use `%w` for wrapping:

```go
if err != nil {
    return fmt.Errorf("failed to <action>: %w", err)
}
```

### No "neoman:" Prefix Internally

The "neoman:" prefix is reserved for user-facing CLI output only. Internal errors (UseCase, Repository, Domain) should not include this prefix.

## Comments

### Exported Functions

All exported functions and types should have a comment starting with the name:

```go
// NewUseCase creates a new UseCase instance with the provided dependencies.
func NewUseCase(...) *UseCase {
```

### Focus on "Why"

Comments should explain *why*, not *what*. The code shows what it does.

```go
// Good: Explains why
// Skip files larger than 10MB to prevent memory issues during indexing
if fileSize > MaxFileSize {

// Bad: Explains what (redundant)
// Check if file size is greater than max file size
if fileSize > MaxFileSize {
```

## See Also

- [Daemon.md](./Daemon.md) - Daemon layer patterns (Domain, UseCase, Controller, Repository, Worker)
- [App.md](./App.md) - App/CLI layer patterns (Command)
- [FeatureChecklist.md](./FeatureChecklist.md) - End-to-end feature implementation guide
