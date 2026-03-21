# Coding Conventions

This document describes coding patterns and standards established by the clean architecture in the daemon layer. Follow these conventions when writing new code in the daemon and future app packages.

## Architecture Principles

**Dependency Flow**: Dependencies flow inward toward the domain.
- Controller depends on UseCase
- UseCase depends on Domain and Repository interfaces
- Domain has no external dependencies
- Repository interfaces have no implementation dependencies

This ensures business logic is independent of infrastructure and easy to test.

## Domain Layer Patterns

### Model Construction with Validation

Domain models must validate all required fields at construction time and return errors if validation fails.

**Pattern**:
```go
func New<ModelName>(requiredField1, requiredField2 string) (*<ModelName>, error) {
    if requiredField1 == "" {
        return nil, fmt.Errorf("requiredField1 must be provided")
    }
    
    // Additional validation
    
    model := &<ModelName>{
        Field1: requiredField1,
        Field2: requiredField2,
    }
    
    return model, nil
}

type <ModelName> struct {
    Field1 string
    Field2 string
    // ... other fields
}
```

**Benefits**:
- Models are self-validating
- Impossible states cannot be represented
- Errors caught early at boundaries

**Example**: `domain/DocsPage.go` - `NewDocsPage()` validates author, repository, and content before creating instance.

### Business Logic in Models

Simple business logic that operates on a model's own data belongs in that model as methods.

**Example**: DocsPage automatically extracts the title from markdown content using `SetTitle()`.

**Note**: Complex orchestration belongs in UseCase, not Domain.

## UseCase Layer Patterns

### Input Structures

Each UseCase operation defines its input as a separate struct with a suffix of "Input".

**Pattern**:
```go
type <OperationName>Input struct {
    Field1 string
    Field2 string
    // Required fields for the operation
}

func (u *UseCase) <OperationName>(input <OperationName>Input) error {
    // Validate input
    // Execute business logic
    // Return error or nil
}
```

**Benefits**:
- Clear contract for what's needed
- Easy to add parameters in future without changing function signature
- Scalable as operations grow more complex

**Example**: `usecase/AddDocsInput` structure contains `Author` and `Repository` fields needed to add documentation.

### Error Constants

Define operation-specific error constants as package-level variables. Each error should be descriptive and specific to its operation.

**Pattern**:
```go
var (
    Err<OperationName><SpecificFailure> = fmt.Errorf("description of failure")
    Err<OperationName><AnotherFailure>  = fmt.Errorf("different failure reason")
)
```

**Benefits**:
- Callers can use type switches or direct comparison
- Different failures can be handled differently
- Errors are documented by their constant names
- No "magic strings" scattered through code

**Example**:
```go
var (
    ErrAddRemoteDocsAuthorRequired      = fmt.Errorf("author is required")
    ErrAddRemoteDocsRepositoryRequired  = fmt.Errorf("repository is required")
    ErrAddRemoteDocsAlreadyExist        = fmt.Errorf("docs already exist")
    ErrAddRemoteDocsDirNotFound         = fmt.Errorf("docs dir not found in repo")
)
```

**Controller Integration**: Controllers use switch statements to handle each error type with appropriate HTTP status codes.

### Dependency Injection via Constructor

UseCase struct receives all dependencies through a constructor function, not globals or singletons.

**Pattern**:
```go
func NewUseCase(
    dependency1 InterfaceType1,
    dependency2 InterfaceType2,
) *UseCase {
    return &UseCase{
        field1: dependency1,
        field2: dependency2,
    }
}

type UseCase struct {
    field1 InterfaceType1
    field2 InterfaceType2
}
```

**Benefits**:
- Dependencies are explicit and visible
- Easy to mock for testing
- No hidden global state
- Clear how this UseCase will behave

**Example**: `usecase/NewUseCase()` accepts `DocsRepository`, `SourceRegistry`, and `git.GitRemote` dependencies.

### File Organization Pattern

The `usecase/` directory follows a specific organization pattern:
- **`usecase.go`**: Contains the single `UseCase` struct definition and the `NewUseCase()` constructor
- **`<operationname>.go`**: Each additional file contains methods that extend the `UseCase` struct (e.g., `addremotedocs.go` contains the `AddRemoteDocs()` method)

**Pattern**:

`usecase.go`:
```go
type UseCase struct {
    docsRepository  repo.DocsRepository
    sourceRegistry  repo.SourceRegistry
    gitRemoteClient git.GitRemote
}

func NewUseCase(
    docsRepository repo.DocsRepository,
    gitRemoteClient git.GitRemote,
    sourceRegistry repo.SourceRegistry,
) *UseCase {
    return &UseCase{
        docsRepository:  docsRepository,
        sourceRegistry:  sourceRegistry,
        gitRemoteClient: gitRemoteClient,
    }
}
```

`addremotedocs.go`:
```go
type AddDocsInput struct {
    Author     string
    Repository string
}

func (u *UseCase) AddRemoteDocs(input AddDocsInput) error {
    // Implementation
}

var (
    ErrAddRemoteDocsAuthorRequired = fmt.Errorf("author is required")
    // ... other error constants for this operation
)
```

**Benefits**:
- Struct definition is centralized and easy to find
- Each file is focused on a single operation
- Easy to scale: add new operations by creating new files
- Clear which methods belong to which struct
- Reduces file size and improves maintainability

## Controller Layer Patterns

### HTTP Handler Error Switching

Controllers handle errors from UseCase operations using a switch statement on the error value, mapping to HTTP status codes.

**Pattern**:
```go
err := c.useCase.<Operation>(input)

switch err {
case nil:
    // Operation succeeded
    http.Error(w, "success message", http.StatusOK)
    return
case usecase.ErrValidationFailure, usecase.ErrInvalidInput:
    http.Error(w, "descriptive message", http.StatusBadRequest)
    return
case usecase.ErrNotFound:
    http.Error(w, "descriptive message", http.StatusNotFound)
    return
default:
    slog.Error("unexpected error", "error", err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
}
```

**Benefits**:
- Each error type maps to appropriate HTTP status
- Exhaustive - each case is explicit
- Logging for unexpected errors
- Clear error messages to clients

**Example**: `controller/AddDocs()` - Handles author/repo required (400), docs not found (404), docs already exist (ignored), and unexpected errors (500).

### Path Value Extraction

Extract path parameters using `r.PathValue()`, which is Go 1.22+ standard library.

**Pattern**:
```go
func (c *controller) Handler(w http.ResponseWriter, r *http.Request) {
    field1 := r.PathValue("field1")
    field2 := r.PathValue("field2")
    
    // Use values
}
```

**Note**: Requires route to be registered with path patterns like `"POST /path/{field1}/{field2}"` in the handler registration.

### Controller Initialization

Create a single function to initialize the HTTP controller with all routes registered.

**Pattern**:
```go
func NewHttpController(useCase *usecase.UseCase) *http.ServeMux {
    mux := http.NewServeMux()
    ctrl := &controller{useCase: useCase}
    
    mux.HandleFunc("POST /endpoint1/{id}", ctrl.Handler1)
    mux.HandleFunc("GET /endpoint2", ctrl.Handler2)
    
    return mux
}

type controller struct {
    useCase *usecase.UseCase
    worker  *worker.Worker  // Optional: if controller triggers workers
}
```

**Benefits**:
- All routes defined in one place
- Easy to see full API surface
- Dependency injection of UseCase
- Controller is unexported (lowercase), methods are unexported too

### File Organization Pattern

The `controller/` directory follows a specific organization pattern similar to UseCase:
- **`controller.go`**: Contains the single `controller` struct definition, the `NewHttpController()` initializer, and route registration
- **`<operationname>.go`**: Each additional file contains handler methods that extend the `controller` struct (e.g., `adddocs.go` contains the `AddDocs()` HTTP handler)

**Pattern**:

`controller.go`:
```go
func NewHttpController(useCase *usecase.UseCase) *http.ServeMux {
    mux := http.NewServeMux()
    ctrl := &controller{useCase: useCase}
    
    mux.HandleFunc("POST /add/{author}/{repo}", ctrl.AddDocs)
    
    return mux
}

type controller struct {
    useCase *usecase.UseCase
    worker  *worker.Worker
}
```

`adddocs.go`:
```go
func (c *controller) AddDocs(w http.ResponseWriter, r *http.Request) {
    author := r.PathValue("author")
    repo := r.PathValue("repo")
    
    err := c.useCase.AddRemoteDocs(usecase.AddDocsInput{
        Author:     author,
        Repository: repo,
    })
    
    switch err {
    case nil:
        // Success handling
    case usecase.ErrAddRemoteDocsAuthorRequired:
        http.Error(w, "Missing author or repository", http.StatusBadRequest)
    // ... more error cases
    }
}
```

**Benefits**:
- Central struct and initialization point is easy to find
- Each file focuses on a single endpoint/handler
- Easy to scale: add new endpoints by creating new files and registering in `controller.go`
- Clear separation of concerns within the controller
- Improves file organization and readability

## Worker Layer Patterns

### Input Structures

Like UseCase, Worker operations define input as separate structs with "Input" suffix.

**Pattern**:
```go
type <TaskName>Input struct {
    Field1 string
    Field2 string
}

func (w *Worker) <TaskName>(input <TaskName>Input) error {
    // Async task logic
    return nil // Or error if critical
}
```

### Dependency Injection

Worker receives repository and service dependencies via constructor.

**Pattern**:
```go
func NewWorker(
    repository1 InterfaceType1,
    repository2 InterfaceType2,
) *Worker {
    return &Worker{
        field1: repository1,
        field2: repository2,
    }
}
```

**Benefits**:
- Same as UseCase - explicit, testable, no globals
- Can be tested independently with mock repositories

## Repository Layer Patterns

### Interface Definitions (Ports)

Define clean interfaces that represent data access operations needed by higher layers.

**Pattern**:
```go
type <ResourceRepository> interface {
    Save(entity <DomainEntity>) error
    GetOne(id string) (<DomainEntity>, error)
    GetAll() ([]<DomainEntity>, error)
    Delete(id string) error
}

type ExternalServiceRegistry interface {
    Download(source <DomainEntity>) error
    GetContents(source <DomainEntity>) ([]<ContentType>, error)
}
```

**Benefits**:
- Defines contract without implementation
- UseCase depends on interface, not concrete implementation
- Easy to mock in tests
- Enables multiple implementations (database, file system, API, etc.)

**Key Principle**: Repository package contains only interfaces and supporting types - no implementation logic.

## Package Organization

### Naming Conventions

- **Packages**: Lowercase, single word (e.g., `domain`, `usecase`, `controller`, `repo`, `worker`)
- **Functions**: PascalCase for exported (e.g., `NewUseCase`, `AddRemoteDocs`)
- **Types**: PascalCase for exported (e.g., `UseCase`, `DocsRepository`, `AddDocsInput`)
- **Variables**: camelCase for unexported (e.g., `docsRepository`)
- **Constants**: PascalCase (e.g., `PrimaryDocsDirName`)
- **Errors**: Prefix with `Err` and operation name (e.g., `ErrAddRemoteDocsAuthorRequired`)

### Comments

- Export functions and types with a comment starting with the name
- Example: `// UseCase orchestrates business operations for documentation management`
- Keep comments focused on "why" not "what" - code shows what it does

## Testing Patterns

### Mocking Repositories

Tests should mock repository interfaces, not implement real databases.

**Pattern**:
```go
// Mock implementation of DocsRepository for testing
type mockDocsRepository struct {}

func (m *mockDocsRepository) Save(entry domain.RemoteDocs) error {
    return nil // Or return test error
}

func (m *mockDocsRepository) Exists(author, repo string) (bool, error) {
    return false, nil // Or return test values
}

// Use in test
useCase := usecase.NewUseCase(
    &mockDocsRepository{},
    &mockSourceRegistry{},
    &mockGitRemote{},
)
```

**Benefits**:
- Fast tests without database
- Deterministic - no flakiness
- Easy to test error paths

## Error Handling

### Error Variable Definition

Define errors at package level for reuse and consistency.

**Pattern**:
```go
var (
    ErrValidationFailed = fmt.Errorf("validation failed")
    ErrNotFound        = fmt.Errorf("resource not found")
)
```

### Error Wrapping

When wrapping errors with context, use `fmt.Errorf()` with context:

```go
if err != nil {
    return fmt.Errorf("failed to save docs: %w", err)
}
```

### No "neoman:" Prefix in Errors

Internal errors (in UseCase, Repository, Domain) should not include "neoman:" prefix. That prefix is reserved for user-facing error messages at the CLI level.

## Import Organization

Organize imports in three groups:
1. Standard library
2. External packages
3. Internal packages

**Example**:
```go
import (
    "fmt"
    "net/http"
    
    "github.com/mattn/go-sqlite3"
    
    "github.com/inodaf/neoman/internal/daemon/domain"
    "github.com/inodaf/neoman/pkg/config"
)
```

## See Also

- [Architecture.md](./Architecture.md) - How layers interact
- [Daemon.md](./Daemon.md) - Layer-by-layer detailed documentation
- AGENTS.md - General code style guidelines (imports, formatting, naming)
