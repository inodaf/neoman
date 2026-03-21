## Project Overview

**Neoman** is a modern documentation reader inspired by Unix `man` pages, written in Go. It consists of:
- `nman`: CLI client for reading documentation
- `nmand`: Daemon for serving documentation

**Language**: Go 1.22+

## Project Structure

The project is organized with clear separation of concerns using clean architecture:

**Current Architecture**:
- `/internal/daemon/` - Backend service implementing clean architecture (domain, usecase, controller, repo, worker)
- `/internal/app/` - Frontend/TUI (planned)
- `/internal/{management,models,operations}/` - Legacy packages (being gradually migrated)
- `/cmd/{daemon,nman}/` - CLI entry points
- `/pkg/` - Public/reusable packages (config, browser, git)
- `/docs/` - Documentation

**For detailed information, see**:
- [docs/Internals/Architecture.md](docs/Internals/Architecture.md) - System architecture and layer design
- [docs/Internals/Daemon.md](docs/Internals/Daemon.md) - Daemon layer documentation
- [docs/Internals/App.md](docs/Internals/App.md) - Frontend/TUI documentation (planned)
- [docs/Internals/Legacy.md](docs/Internals/Legacy.md) - Legacy package information
- [docs/Internals/CodingConventions.md](docs/Internals/CodingConventions.md) - Architecture-specific patterns

## Build

```bash
# Build both CLI and daemon
make all

# Build only CLI
go build -o ./bin/nman cmd/nman/main.go

# Build only daemon
go build -o ./bin/nmand cmd/daemon/main.go
```

## Before Committing

1. Always ask for a code review and wait approval
2. Run `go fmt` to format code
3. Run `go vet` to check for issues
4. Run `go test` to ensure tests pass
5. Verify builds: `make all`
6. Check that changes align with existing code style in relevant package
7. Commit with meaningful messages E.g: "[Add|Update|Remove] Daemon/UseCase: List documentation now return doc titles"
