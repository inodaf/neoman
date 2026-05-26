# Daemon Conventions

This document describes coding patterns for the daemon layer (`/internal/daemon/`). All patterns use generic placeholders to remain stable as the codebase evolves.

## Domain Layer

Location: `/internal/daemon/domain/`

### Model Construction with Validation

Domain models validate all required fields at construction time and return errors on failure.

```go
func New<ModelName>(<field1>, <field2> string) (*<ModelName>, error) {
    if <field1> == "" {
        return nil, fmt.Errorf("<field1> must be provided")
    }

    if <field2> == "" {
        return nil, fmt.Errorf("<field2> must be provided")
    }

    return &<ModelName>{
        <Field1>: <field1>,
        <Field2>: <field2>,
    }, nil
}

type <ModelName> struct {
    <Field1> string
    <Field2> string
}
```

### Business Logic in Models

Simple business logic that operates on a model's own data belongs in that model as methods:

```go
func (m *<ModelName>) Set<DerivedField>() error {
    if m.<SourceField> == "" {
        return fmt.Errorf("unable to derive <DerivedField> without <SourceField>")
    }

    m.<DerivedField> = // derive from m.<SourceField>
    return nil
}
```

**Note**: Complex orchestration involving multiple models or external services belongs in UseCase, not Domain.

## UseCase Layer

Location: `/internal/daemon/usecase/`

### File Organization

- **`usecase.go`**: Contains the `UseCase` struct and `NewUseCase()` constructor
- **`<operationname>.go`**: Each operation in its own file with Input/Output structs and error constants

### UseCase Struct and Constructor

```go
// usecase.go

func NewUseCase(
    <dependency1> repo.<Interface1>,
    <dependency2> repo.<Interface2>,
) *UseCase {
    return &UseCase{
        <field1>: <dependency1>,
        <field2>: <dependency2>,
    }
}

type UseCase struct {
    <field1> repo.<Interface1>
    <field2> repo.<Interface2>
}
```

### Input Structures

Each operation defines its input as a separate struct:

```go
type <OperationName>Input struct {
    <Field1> string
    <Field2> string
}
```

### Output Structures

Operations that return data define an output struct:

```go
type <OperationName>Output struct {
    <Items> []domain.<ModelName>
}
```

### Operation Method (Input Only)

For operations that don't return data:

```go
func (u *UseCase) <OperationName>(input <OperationName>Input) error {
    if input.<Field1> == "" {
        return Err<OperationName><Field1>Required
    }

    if input.<Field2> == "" {
        return Err<OperationName><Field2>Required
    }

    // Business logic using u.<dependency>
    
    return nil
}
```

### Operation Method (Input and Output)

For operations that return data:

```go
func (u *UseCase) <OperationName>(input <OperationName>Input) (<OperationName>Output, error) {
    if input.<Field1> == "" {
        return <OperationName>Output{}, Err<OperationName><Field1>Required
    }

    // Business logic
    result, err := u.<repository>.<Method>(input.<Field1>, input.<Field2>)
    if err != nil {
        return <OperationName>Output{}, err
    }

    return <OperationName>Output{<Items>: result}, nil
}
```

### Error Constants

Define operation-specific errors at the end of the file:

```go
var (
    Err<OperationName><Field1>Required = fmt.Errorf("<field1> is required")
    Err<OperationName><Field2>Required = fmt.Errorf("<field2> is required")
    Err<OperationName>NotFound         = fmt.Errorf("<resource> not found")
)
```

## Controller Layer

Location: `/internal/daemon/controller/`

### File Organization

- **`controller.go`**: Contains `controller` struct, `NewHttpController()`, and route registration
- **`<operationname>.go`**: Each endpoint handler in its own file

### Controller Struct and Initialization

```go
// controller.go

func NewHttpController(useCase *usecase.UseCase, <otherDeps>...) *http.ServeMux {
    mux := http.NewServeMux()
    ctrl := &controller{useCase: useCase, <otherFields>...}

    mux.HandleFunc("<METHOD> /<path>/{<param1>}/{<param2>}", ctrl.<HandlerName>)
    // Register more routes...

    return mux
}

type controller struct {
    useCase       *usecase.UseCase
    <otherFields> <types>
}
```

### Handler Method

```go
// <operationname>.go

func (c *controller) <HandlerName>(w http.ResponseWriter, r *http.Request) {
    <param1> := r.PathValue("<param1>")
    <param2> := r.PathValue("<param2>")

    output, err := c.useCase.<OperationName>(usecase.<OperationName>Input{
        <Field1>: <param1>,
        <Field2>: <param2>,
    })

    switch err {
    case nil:
        // Success handling
        w.Header().Set("Content-Type", "<content-type>")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(<formattedResponse>))
        return
    case usecase.Err<OperationName>NotFound:
        http.Error(w, "<Resource> not found", http.StatusNotFound)
        return
    case usecase.Err<OperationName><Field1>Required, usecase.Err<OperationName><Field2>Required:
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    default:
        slog.Error("failed to <operation>", "error", err)
        http.Error(w, "Unable to <operation>", http.StatusInternalServerError)
        return
    }
}
```

### Content Negotiation with Accept Header

When supporting multiple response formats:

```go
func (c *controller) <HandlerName>(w http.ResponseWriter, r *http.Request) {
    accept := r.Header.Get("Accept")
    if accept != "" && accept != "<supported-type>" && accept != "*/*" {
        http.Error(w, "Unsupported media type. Use Accept: <supported-type>", http.StatusNotAcceptable)
        return
    }

    // Handle request...
}
```

### Custom Response Formatting

For simple formats, use a private helper function instead of external libraries:

```go
func format<ResponseType>(<params>...) string {
    var sb strings.Builder

    sb.WriteString(fmt.Sprintf("<key>: %s\n", <value>))
    sb.WriteString("<section>:\n")

    for _, item := range <items> {
        sb.WriteString(fmt.Sprintf("  - %s: '%s'\n", item.<Field1>, item.<Field2>))
    }

    return sb.String()
}
```

## Repository Layer

Location: `/internal/daemon/repo/`

### Interface Definitions (Ports)

The repository package contains **only interfaces and supporting types**, no implementations.

```go
// ports.go

type <Resource>Repository interface {
    Save(entry domain.<ModelName>) error
    Exists(<identifier1>, <identifier2> string) (bool, error)
    GetOne(<identifier1>, <identifier2> string) (domain.<ModelName>, error)
    FindAll(<identifier1>, <identifier2> string) ([]domain.<ModelName>, error)
}
```

### Supporting Types

Define types needed by interfaces:

```go
type <ContentType> struct {
    <Field1> string
    <Field2> string
}
```

### Implementation Files

Implementation files live alongside ports but implement the interfaces:

```go
// sql<resource>repository.go

type <resource>Repository struct {
    db *sql.DB
}

func New<Resource>Repository(db *sql.DB) <Resource>Repository {
    return &<resource>Repository{db: db}
}

func (r *<resource>Repository) Save(entry domain.<ModelName>) error {
    query := `INSERT INTO <table> (...) VALUES (...)`
    _, err := r.db.Exec(query, ...)
    if err != nil {
        return fmt.Errorf("failed to save <resource>: %w", err)
    }
    return nil
}
```

## Worker Layer

Location: `/internal/daemon/worker/`

### Input Structures

Like UseCase, Worker operations define input structs:

```go
type <TaskName>Input struct {
    <Field1> string
    <Field2> string
}
```

### Worker Struct and Constructor

```go
func NewWorker(
    <repository1> repo.<Interface1>,
    <repository2> repo.<Interface2>,
) *Worker {
    return &Worker{
        <field1>: <repository1>,
        <field2>: <repository2>,
    }
}

type Worker struct {
    <field1> repo.<Interface1>
    <field2> repo.<Interface2>
}
```

### Task Method

```go
func (w *Worker) <TaskName>(input <TaskName>Input) error {
    // Background task logic
    // Uses w.<repository> for data access
    
    return nil
}
```

## See Also

- [Overview.md](./Overview.md) - General principles and naming conventions
- [App.md](./App.md) - CLI/Command layer patterns
- [FeatureChecklist.md](./FeatureChecklist.md) - End-to-end feature implementation guide
