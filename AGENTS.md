## Project Overview

**Neoman** is a modern documentation reader inspired by Unix `man` pages, written in Go. It consists of:
- `nman`: CLI client for reading documentation
- `nmand`: Daemon for serving documentation

**Language**: Go 1.22+

## Project Structure

The project is organized with clear separation of concerns using clean architecture:

**Current Architecture**:
- `/internal/daemon/` - Backend service implementing clean architecture (domain, usecase, controller, repo, worker)
- `/internal/app/` - CLI client (command layer)
- `/internal/{management,models,operations}/` - Legacy packages (being gradually migrated)
- `/cmd/{daemon,nman}/` - CLI entry points
- `/pkg/` - Public/reusable packages (config, browser, git)
- `/docs/` - Documentation

**For detailed information, see**:
- [docs/Internals/Architecture.md](docs/Internals/Architecture.md) - System overview and layer design
- [docs/Internals/CodingConventions/](docs/Internals/CodingConventions/) - Implementation patterns
- [docs/Internals/Legacy.md](docs/Internals/Legacy.md) - Legacy package information

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
