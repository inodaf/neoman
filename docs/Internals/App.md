# App

The frontend/TUI for Neoman is currently under development in the `/internal/app/` package.

## Current Status

**Work in Progress** - This section is a placeholder for the upcoming frontend implementation.

## Vision

The app will be a terminal-based user interface (TUI) that:
- Communicates with the daemon via HTTP API
- Provides an interactive interface for browsing and searching documentation
- Renders rich terminal output
- Manages user sessions and preferences
- Handles keyboard navigation and search functionality

## Expected Structure

```
/internal/app/
├── view/        # Terminal rendering and UI components
├── controller/  # Input handling and view orchestration
└── models/      # Client-side state and data models
```

## Communication with Daemon

The app will communicate with the daemon through its HTTP API:
- Base URL: `http://localhost:9090` (or Unix socket at `/tmp/nman.sock`)
- All endpoints defined in [Daemon.md](./Daemon.md)
- Request/response handling abstracted in a client package

## Future Documentation

Once development begins, this file will be updated with:
- Detailed package structure and responsibilities
- View rendering patterns
- State management approach
- Navigation and command handling
- Configuration and preferences

## See Also

- [Architecture.md](./Architecture.md) - Overall system design
- [Daemon.md](./Daemon.md) - Backend API and operations
- [CodingConventions.md](./CodingConventions.md) - Code style and patterns
