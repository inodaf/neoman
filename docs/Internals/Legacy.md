# Legacy

The following packages under `/internal/` predate the clean architecture refactoring and are considered legacy. They are gradually being migrated into the new structure.

## Legacy Packages

### `internal/management/`

**Status**: Legacy - Being phased out  
**Purpose**: Originally provided core daemon functionality including:
- Database initialization and management (`database.go`)
- Cron job scheduling (`cron.go`)
- Documentation registry management (`registry.go`)
- Unix socket communication (`sockets.go`)
- Version information (`version.go`)

**Migration Path**:
- Database operations → `daemon/repo/` (repository implementations)
- Registry management → `daemon/usecase/` (business logic)
- Socket management → HTTP controller in `daemon/controller/`
- Version info → `pkg/config/`

**Files**:
- `database.go` - SQLite database operations
- `cron.go` - Scheduled task execution
- `registry.go` - Documentation registry CRUD
- `sockets.go` - Unix domain socket server
- `version.go` - Version constants

### `internal/models/`

**Status**: Legacy - Being replaced  
**Purpose**: Data models for documentation management

**Migration Path**:
- `docpage.go` → `daemon/domain/DocsPage` model

**Files**:
- `docpage.go` - Documentation page structure

### `internal/operations/`

**Status**: Legacy - Being replaced  
**Purpose**: CLI-level operations for managing documentation

**Migration Path**:
- Authentication logic → `daemon/usecase/` or `pkg/`
- List operations → `daemon/usecase/ListDocs` (planned)
- Sync operations → `daemon/usecase/SyncDocs` (planned)
- Clean operations → `daemon/usecase/CleanDocs` (planned)

**Files**:
- `auth.go` - Authentication logic
- `list.go` - List available documentation
- `sync.go` - Synchronize documentation
- `clean.go` - Clean/remove documentation
- `sync_test.go` - Sync operation tests

## Timeline for Removal

**No fixed timeline** - Legacy packages will be removed as their functionality is fully migrated to the new architecture. Priority depends on:
1. Clean architecture implementation coverage
2. Test coverage of new implementations
3. Verification that no critical functionality is lost

## Migration Approach

When migrating legacy functionality:

1. **Identify Responsibility** - Which architectural layer should own this?
2. **Create New Implementation** - Implement in appropriate daemon layer
3. **Test Thoroughly** - Ensure behavior is equivalent
4. **Update Callers** - Point to new implementation
5. **Deprecate Legacy** - Leave legacy code but mark as deprecated
6. **Remove** - Delete once all callers migrated

## Current Usage

These legacy packages are still actively used by:
- `cmd/daemon/main.go` - Daemon initialization
- `cmd/nman/main.go` - CLI client operations

As the clean architecture matures, these entry points will be refactored to use daemon layer components instead.

## See Also

- [Architecture.md](./Architecture.md) - New clean architecture
- [Daemon.md](./Daemon.md) - Where legacy functionality moves to
- [CodingConventions.md](./CodingConventions.md) - Standards for new code
