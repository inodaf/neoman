# Neoman

[![Static Badge](https://img.shields.io/badge/Docs-%24_nman_inodaf%2Fneoman-black)](https://github.com/inodaf/neoman)

A modern documentation reader inspired by Unix `man` pages, designed to make software documentation accessible, searchable, and maintainable.

> [!NOTE]
> Neoman is currently in early development. Expect frequent updates and changes as we refine the experience.

## Features

**Agentic Engineering** - Query documentation semantically, retrieve context for coding agents, and track token usage across your docs. Built for AI-assisted development workflows.

**Local & Secure** - Everything runs locally on your machine. No data leaves your device, perfect for private organizational documentation.

**Zero-Deployment Documentation** - Focus on writing great docs, not managing infrastructure. No servers to maintain, no hosting costs. Just write Markdown and push to Git.

**Git-Native Workflow** - Documentation stays in sync with your code automatically. Push to any Git provider (GitHub, GitLab, etc.) and readers get updates instantly.

**Convention Over Configuration** - Create a `/docs` directory with an `index.md` file and you're ready. No complex setup required.

**Unified Documentation Hub** - All your organization's documentation in one searchable place. Search across all docs or filter by specific projects.

## Requirements

- Go 1.22 or later
- macOS or Linux

## Installation

Build from source:

```sh
git clone https://github.com/inodaf/neoman.git
cd neoman
make all
```

The binaries will be created in `./bin/`:
- `./bin/nman` - CLI client for reading documentation
- `./bin/nmand` - Daemon server for serving documentation

Add to your PATH:

```sh
export PATH="$PATH:$(pwd)/bin"
```

## Quick Start

Start using Neoman immediately:

```sh
nman inodaf/neoman
```

This opens the Neoman documentation using Neoman itself.

## Usage

### Open Documentation by Project

```sh
nman author/project
```

Opens the documentation for the specified project (e.g., `nman inodaf/neoman`).

### Open Documentation from Current Directory

```sh
nman .
# or just
nman
```

Looks for a `docs` directory in the current Git repository and opens it.

### List Available Documentation

```sh
# List all available projects
nman list

# List all projects in an organization
nman list author

# List all documents in a specific project
nman list author/project
```

### View Specific Document

```sh
nman view author/project "path/to/document.md"
```

## Adopting Neoman

To make your project's documentation available through Neoman:

1. Create a `/docs` directory in your repository root
2. Add an `index.md` file with your documentation index
3. Add additional `.md` files for your documentation pages

Example structure:

```
your-repo/
├── docs/
│   ├── index.md
│   ├── getting-started.md
│   ├── api-reference.md
│   └── guides/
│       ├── authentication.md
│       └── deployment.md
├── src/
└── README.md
```

That's it! Your documentation will be accessible via:

```sh
nman your-username/your-repo
```

For detailed authoring guidelines, see [docs/Authoring Docs/](docs/Authoring%20Docs/).

## Contributing

We welcome contributions! Here's how to get started:

1. Review the [Architecture documentation](docs/Internals/Architecture.md) to understand the codebase
2. Check out [Coding Conventions](docs/Internals/CodingConventions.md) for style guidelines
3. Submit pull requests with meaningful commit messages

For more details, see the [docs/Internals/](docs/Internals/) directory.

## License

MIT License - see LICENSE file for details
