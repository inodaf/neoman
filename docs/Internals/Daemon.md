# Daemon

The daemon (`/internal/daemon/`) is the backend service that manages documentation. It implements clean architecture with five distinct layers that handle different responsibilities.

## Package Structure

```
/internal/daemon/
├── domain/           # Core business models and rules
├── usecase/          # Business logic operations
├── controller/       # HTTP endpoints and request handling
├── repo/             # Data access abstraction (ports/interfaces)
└── worker/           # Background tasks and async operations
```

## Domain Layer (`domain/`)

Defines the core business entities that represent documentation management concepts.

### Key Models

**RemoteDocs**
- Represents documentation from a remote source (e.g., GitHub repository)
- Properties: `Author`, `Repository`, `Source` (RemoteSource), `Indexing` (bool)
- Constructor: `NewRemoteDocs(author, repository, source)` with validation
- Immutable once created

**DocsPage**
- Represents a single documentation page/file
- Properties: `Title`, `Author`, `Repository`, `Content`, `RelativePath`, `LastModifiedAt`, `Vector` (for embeddings)
- Constructor: `NewDocsPage(author, repo, content, relPath)` with validation
- Automatically infers title from content
- Validates all required fields are present

**RemoteSource**
- Enum representing the source platform (currently GitHub)
- Used to identify where documentation originates

### Business Rules
- RemoteDocs and DocsPage require all essential fields to be provided at construction
- DocsPage automatically extracts title from markdown content
- Models self-validate and return errors during construction
- No getters/setters; direct field access after validation

## UseCase Layer (`usecase/`)

Implements business operations that orchestrate domain models and repositories.

### Core Operations

**AddRemoteDocs**
- Input: `AddDocsInput{Author, Repository}`
- Flow:
  1. Validate input fields are provided
  2. Check if docs already exist
  3. Verify repository has a `docs/` directory
  4. Create RemoteDocs domain model
  5. Download documentation source
  6. Save RemoteDocs to persistence
- Errors: Author/Repository required, docs already exist, directory not found, download failures

### Error Handling Pattern
- Each operation defines its own error constants at the package level
- Example: `var ErrAddRemoteDocsAuthorRequired = fmt.Errorf("author is required")`
- Errors are specific to allow callers to handle different failure modes
- Errors never contain "neoman:" prefix (that's for user-facing output)

### Dependency Injection
UseCase struct holds references to required dependencies:
```go
type UseCase struct {
    docsRepository  repo.DocsRepository      // Data persistence
    sourceRegistry  repo.SourceRegistry      // External source access
    gitRemoteClient git.GitRemote           // Git operations
}
```

**Constructor**: `NewUseCase(docsRepository, gitRemoteClient, sourceRegistry)` initializes with dependencies.

## Controller Layer (`controller/`)

HTTP endpoints that expose usecase operations to external clients.

### Request Handling Pattern

**AddDocs Endpoint** (`POST /add/{author}/{repo}`)
- Extracts path parameters using `r.PathValue()`
- Calls corresponding usecase method
- Error handling via switch statement on returned error
- Categorizes errors into HTTP status codes:
  - `StatusBadRequest` (400): Validation failures
  - `StatusNotFound` (404): Resource not found
  - `StatusInternalServerError` (500): Unexpected errors
- Triggers background worker for async processing on success

### Response Strategy
- HTTP status codes communicate success/failure
- Error messages are human-readable and descriptive
- Structured logging for debugging unexpected errors

### Initialization
- `NewHttpController(useCase)` creates an `http.ServeMux` with all endpoints registered
- Accepts fully-constructed UseCase as dependency

## Repository Layer (`repo/`)

Defines interfaces (ports) for data access without implementation details.

### Ports/Interfaces

**DocsRepository**
- `Save(RemoteDocs)` - Persist documentation metadata
- `Exists(author, repository)` - Check if docs already registered
- `GetOne(author, repository)` - Retrieve specific docs
- `StartIndexing(author, repository)` - Mark docs as indexing in progress
- `StopIndexing(author, repository)` - Mark indexing complete

**DocsPageRepository**
- `Save(DocsPage)` - Persist a single page
- `SaveMany([]DocsPage)` - Batch persist multiple pages

**SourceRegistry**
- `Download(RemoteDocs)` - Fetch documentation from source
- `GetAllContents(RemoteDocs)` - List all documentation files with metadata

### Design Notes
- These are contracts/interfaces, not implementations
- Implementations live in concrete packages (database adapters, file system, etc.)
- UseCase and Worker depend only on these interfaces
- Enables easy testing with mock implementations

## Worker Layer (`worker/`)

Handles background operations that shouldn't block HTTP responses.

### Background Tasks

**IndexPages**
- Input: `IndexPagesInput{Author, Repository}`
- Operation: Parse all documentation files and create searchable index
- Flow:
  1. Mark docs as indexing
  2. Fetch all file contents from source registry
  3. Create DocsPage models from each file
  4. Batch save pages to repository
  5. Mark indexing complete

### Async Execution
- Triggered from Controller after successful primary operation
- Runs independently without blocking client response
- Failures are logged but don't affect HTTP response
- Designed for horizontal scaling (can run in separate processes/containers)

### Dependency Injection
Worker struct holds repository references:
```go
type Worker struct {
    sourceRegistry     repo.SourceRegistry
    docsPageRepository repo.DocsPageRepository
    docsRepository     repo.DocsRepository
}
```

**Constructor**: `NewWorker(docsPageRepository, sourceRegistry, docsRepository)` initializes with dependencies.

## Integration with pkg/ Packages

The daemon integrates with shared utility packages:

**git.GitRemote** (`pkg/git/`)
- Verifies repository existence
- Checks for `docs/` directory presence
- Identifies the Git provider (GitHub, GitLab, etc.)

**config** (`pkg/config/`)
- Configuration constants
- Application settings

**browser** (`pkg/browser/`)
- Utilities for opening documentation in browser (future)

## HTTP API

### Endpoints

**POST /add/{author}/{repo}**
- Add remote documentation to the registry
- Parameters: `{author}` (GitHub user/org), `{repo}` (repository name)
- Response: `200 OK` on success, appropriate error status on failure
- Side effect: Triggers background indexing worker

### Communication Mechanism
- Default: Unix domain socket at `/tmp/nman.sock`
- Can also use HTTP over TCP for remote access
- Client authentication and authorization handled separately

## Adding a New Operation

To add a new feature to the daemon:

1. **Domain Model** - If needed, add new types to `domain/`
2. **Interfaces** - Extend `repo/` with new port interfaces if new persistence needed
3. **UseCase** - Implement operation in `usecase/` with input structure and error constants
4. **Controller** - Add HTTP handler in `controller/` that calls the usecase
5. **Worker** - If async processing needed, add task to `worker/`
6. **Integration** - Register endpoint in `NewHttpController()`

## Error Handling Flow

```
Domain Model Construction (validation)
    ↓ (error) → UseCase catches, returns error constant
    ↓ (success)
UseCase Business Logic (validation + operations)
    ↓ (error) → Error constant specific to operation
    ↓ (success)
Controller Receives Error (or nil)
    ├─ Switch on error type
    ├─ Map to HTTP status code
    └─ Return human-readable message
```

## See Also

- [Architecture.md](./Architecture.md) - Overall system design
- [CodingConventions.md](./CodingConventions.md) - Patterns and standards used
- [App.md](./App.md) - Frontend client (planned)
