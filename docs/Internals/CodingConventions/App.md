# App Conventions

This document describes coding patterns for the app/CLI layer (`/internal/app/`). All patterns use generic placeholders to remain stable as the codebase evolves.

## Command Layer Overview

Location: `/internal/app/command/`

The command layer implements CLI commands that communicate with the daemon via HTTP over Unix socket. Each command:

1. Parses and validates user input
2. Makes HTTP requests to daemon endpoints
3. Handles responses and errors
4. Outputs results to stdout/stderr
5. Exits with appropriate status codes

## File Organization

- **`command.go`**: Contains `Command` struct, `New()` constructor, and shared helper methods
- **`<commandname>.go`**: Each command in its own file with public method, private helper, and error constants

## Command Struct and Constructor

```go
// command.go

func New() *Command {
    daemonClient := http.Client{
        Transport: &http.Transport{
            DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
                return net.Dial("unix", config.AppSockPath)
            },
        },
    }

    return &Command{daemonHttpClient: daemonClient}
}

type Command struct {
    daemonHttpClient http.Client
}
```

## Shared Helper Methods

Common parsing/validation logic lives in `command.go`:

```go
func (c *Command) parse<InputType>(arg string) (<field1>, <field2> string, err error) {
    // Validation logic
    if <invalid> {
        return "", "", fmt.Errorf("<error message>")
    }

    // Parse and return fields
    return <field1>, <field2>, nil
}
```

## Command Method Pattern

Each command has a public entry point and a private HTTP helper.

### Public Method (Entry Point)

```go
// <commandname>.go

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
```

### Private HTTP Helper

```go
func (c *Command) <httpHelper>(ctx context.Context, <field1>, <field2> string) (string, error) {
    resource := url.URL{
        Host:   "unix",
        Scheme: "http",
        Path:   fmt.Sprintf("/<endpoint>/%s/%s", <field1>, <field2>),
    }

    req, err := http.NewRequestWithContext(ctx, "<METHOD>", resource.String(), nil)
    if err != nil {
        return "", err
    }

    res, err := c.daemonHttpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    // Handle response status codes
    if res.StatusCode == http.StatusNotFound {
        return "", Err<CommandName>NotFound
    }

    if res.StatusCode == http.StatusNoContent {
        return "", Err<CommandName>NoContent
    }

    if res.StatusCode != http.StatusOK {
        return "", Err<CommandName>Unexpected
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}
```

### Error Constants

Define command-specific errors at the end of the file:

```go
var (
    Err<CommandName>NotFound   = fmt.Errorf("<Resource> not found for this project")
    Err<CommandName>NoContent  = fmt.Errorf("No <items> found for this project")
    Err<CommandName>Unexpected = fmt.Errorf("Unable to <action>")
)
```

## Exit Code Conventions

Commands must exit with appropriate status codes:

| Scenario | Exit Code | Output |
|----------|-----------|--------|
| Success | `0` | Result to stdout |
| User/input error | `1` | Error message to stdout with "neoman:" prefix |
| Daemon/server error | `1` | Error message to stdout with "neoman:" prefix |

```go
// Success
fmt.Print(output)
os.Exit(0)

// Error
fmt.Printf("neoman: %s.\n", err.Error())
os.Exit(1)
```

## Wiring Commands in main.go

Location: `/cmd/nman/main.go`

Commands are wired in the main switch statement:

```go
func main() {
    cmd := command.New()

    // ... daemon connectivity check ...

    switch os.Args[1] {
    case "<command>":
        if <validation> {
            fmt.Println("neoman: Usage: nman <command> <args>")
            return
        }
        cmd.<CommandName>(context.TODO(), os.Args[2])
        return
    // ... more commands ...
    default:
        fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
        return
    }
}
```

## Error Message Format

User-facing error messages use the "neoman:" prefix:

```go
// In command methods
fmt.Printf("neoman: %s.\n", err.Error())

// In main.go for usage errors
fmt.Println("neoman: Usage: nman <command> <args>")

// In main.go for invalid commands
fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
```

**Note**: The "neoman:" prefix is only used at the CLI output level. Internal errors (in error constants) should not include this prefix.

## See Also

- [Overview.md](./Overview.md) - General principles and naming conventions
- [Daemon.md](./Daemon.md) - Daemon layer patterns
- [FeatureChecklist.md](./FeatureChecklist.md) - End-to-end feature implementation guide
