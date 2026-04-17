# Feature Implementation Checklist

This document provides a step-by-step checklist for implementing new features end-to-end, from daemon to CLI.

## Planning

Before writing code:

- [ ] Define the feature scope and expected behavior
- [ ] Identify which layers need changes (Domain, UseCase, Controller, Repository, Worker, Command)
- [ ] Determine the API endpoint path and HTTP method
- [ ] Define input/output data structures

---

## Daemon Implementation

### Step 1: Domain Model (if needed)

Location: `/internal/daemon/domain/<modelname>.go`

- [ ] Create model struct with fields
- [ ] Implement constructor `New<ModelName>()` with validation
- [ ] Add business logic methods if needed

```go
func New<ModelName>(<field1>, <field2> string) (*<ModelName>, error) {
    if <field1> == "" {
        return nil, fmt.Errorf("<field1> must be provided")
    }
    return &<ModelName>{<Field1>: <field1>, <Field2>: <field2>}, nil
}

type <ModelName> struct {
    <Field1> string
    <Field2> string
}
```

### Step 2: Repository Interface (if needed)

Location: `/internal/daemon/repo/ports.go`

- [ ] Add new method(s) to existing interface, OR
- [ ] Create new interface for new resource type

```go
type <Resource>Repository interface {
    // ... existing methods ...
    <NewMethod>(<params>...) (<returnType>, error)
}
```

### Step 3: Repository Implementation

Location: `/internal/daemon/repo/sql<resource>repository.go`

- [ ] Implement the new interface method(s)

```go
func (r *<resource>Repository) <NewMethod>(<params>...) (<returnType>, error) {
    query := `SELECT ... FROM <table> WHERE ...`
    // Execute query and return results
}
```

### Step 4: UseCase

Location: `/internal/daemon/usecase/<operationname>.go`

- [ ] Create Input struct
- [ ] Create Output struct (if returning data)
- [ ] Implement operation method
- [ ] Define error constants

```go
type <OperationName>Input struct {
    <Field1> string
    <Field2> string
}

type <OperationName>Output struct {
    <Items> []domain.<ModelName>
}

func (u *UseCase) <OperationName>(input <OperationName>Input) (<OperationName>Output, error) {
    // Validation
    // Business logic
    // Return output or error
}

var (
    Err<OperationName><Failure1> = fmt.Errorf("<description>")
    Err<OperationName><Failure2> = fmt.Errorf("<description>")
)
```

### Step 5: Controller Endpoint

Location: `/internal/daemon/controller/<operationname>.go`

- [ ] Implement HTTP handler method
- [ ] Extract path parameters
- [ ] Call UseCase method
- [ ] Handle errors with switch statement
- [ ] Format and return response

```go
func (c *controller) <HandlerName>(w http.ResponseWriter, r *http.Request) {
    <param1> := r.PathValue("<param1>")
    <param2> := r.PathValue("<param2>")

    output, err := c.useCase.<OperationName>(usecase.<OperationName>Input{
        <Field1>: <param1>,
        <Field2>: <param2>,
    })

    switch err {
    case nil:
        w.Header().Set("Content-Type", "<content-type>")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(<response>))
    case usecase.Err<OperationName>NotFound:
        http.Error(w, "<Resource> not found", http.StatusNotFound)
    // ... more cases ...
    default:
        slog.Error("failed to <operation>", "error", err)
        http.Error(w, "Unable to <operation>", http.StatusInternalServerError)
    }
}
```

### Step 6: Register Route

Location: `/internal/daemon/controller/controller.go`

- [ ] Add route to `NewHttpController()`

```go
mux.HandleFunc("<METHOD> /<path>/{<param1>}/{<param2>}", ctrl.<HandlerName>)
```

### Step 7: Worker Task (if async needed)

Location: `/internal/daemon/worker/<taskname>.go`

- [ ] Create Input struct
- [ ] Implement task method
- [ ] Trigger from controller after success

---

## App/CLI Implementation

### Step 8: Command Method

Location: `/internal/app/command/<commandname>.go`

- [ ] Create public command method
- [ ] Create private HTTP helper method
- [ ] Define error constants
- [ ] Handle exit codes

```go
func (c *Command) <CommandName>(ctx context.Context, arg string) {
    <field1>, <field2>, err := c.parse<InputType>(arg)
    if err != nil {
        fmt.Printf("neoman: %s.\n", err.Error())
        os.Exit(1)
        return
    }

    output, err := c.<httpHelper>(ctx, <field1>, <field2>)
    if err != nil {
        fmt.Printf("neoman: %s.\n", err.Error())
        os.Exit(1)
        return
    }

    fmt.Print(output)
    os.Exit(0)
}

func (c *Command) <httpHelper>(ctx context.Context, <field1>, <field2> string) (string, error) {
    resource := url.URL{Host: "unix", Scheme: "http", Path: fmt.Sprintf("/<endpoint>/%s/%s", <field1>, <field2>)}
    // Make request, handle response
}

var (
    Err<CommandName>NotFound  = fmt.Errorf("<Resource> not found")
    Err<CommandName>Unexpected = fmt.Errorf("Unable to <action>")
)
```

### Step 9: Wire Command in main.go

Location: `/cmd/nman/main.go`

- [ ] Add case to switch statement
- [ ] Add input validation
- [ ] Call command method

```go
case "<command>":
    if <validation> {
        fmt.Println("neoman: Usage: nman <command> <args>")
        return
    }
    cmd.<CommandName>(context.TODO(), os.Args[2])
    return
```

---

## Testing

### Step 10: UseCase Unit Tests

Location: `/internal/daemon/usecase/<operationname>_test.go`

- [ ] Create mock repositories
- [ ] Test success cases
- [ ] Test validation errors
- [ ] Test repository errors

### Step 11: Command Unit Tests

Location: `/internal/app/command/<commandname>_test.go`

- [ ] Mock HTTP responses
- [ ] Test success output
- [ ] Test error handling
- [ ] Test exit codes

### Step 12: Integration Tests (if applicable)

- [ ] Test endpoint with real database (test instance)
- [ ] Test CLI command with running daemon

---

## Verification

Before submitting for review:

- [ ] Run `go fmt ./...`
- [ ] Run `go vet ./...`
- [ ] Run `go test ./...`
- [ ] Run `make` to build both binaries
- [ ] Test daemon endpoint manually with cURL:
  ```bash
  curl --unix-socket /tmp/nman.sock -X <METHOD> http://localhost/<endpoint>
  ```
- [ ] Test CLI command manually:
  ```bash
  ./bin/nman <command> <args>
  ```

---

## See Also

- [Overview.md](./Overview.md) - General principles and naming conventions
- [Daemon.md](./Daemon.md) - Daemon layer patterns
- [App.md](./App.md) - CLI/Command layer patterns
