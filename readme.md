<p align="center">
  <strong>Neoman</strong><br/>
  <em>Equip your AI Agents with Internal/Enterprise Docs</em>
</p>

<p align="center">
  <a href="https://github.com/inodaf/neoman"><img alt="Go Version" src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white"></a>
  <a href="https://github.com/inodaf/neoman/blob/main/LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
  <a href="https://github.com/inodaf/neoman"><img alt="Docs" src="https://img.shields.io/badge/Docs-%24_nman_inodaf%2Fneoman-black"></a>
</p>

<p align="center">
  A modern documentation reader inspired by Unix <code>man</code> pages.<br/>
  Read, search, and serve documentation locally — built for AI-assisted development workflows.
</p>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#quickstart">Quickstart</a> •
  <a href="#usage">Usage</a> •
  <a href="#adopting-neoman">Adopting Neoman</a> •
  <a href="#contributing">Contributing</a>
</p>

---

> [!NOTE]
> Neoman is currently in early development. Expect frequent updates and changes as we refine the experience.

## Overview

Neoman brings the simplicity of Unix `man` pages to modern software documentation. Docs stays in your Git repositories. Neoman fetches, indexes, and serves it locally. No external servers, no hosting costs, no data leaving your machine.

## Features

### Available Now

- **Local & Safe** — Runs and stays on your device. Perfect for private org documentations.

- **Documentation Hub** — All your org's internal documentations in one local searchable place.

- **Git-Native Workflow** — Documentation syncs with your code. Push to any Git provider (GitHub, GitLab, etc.).

- **Practical Convention** — A `/docs` directory with `.md` files is all you need.

- **Zero-Deployment** — No servers to maintain. Write Markdown, push to a Git provider, done.

### Coming Soon

- **Semantic Search** — Query using natural language. Get relevance-ranked results with preview snippets.

- **Agentic Context Retrieval** — Retrieve relevant up-to-date documentation context for coding agents and AI assistants.


## Installation

### Requirements

- Go 1.22 or later
- macOS or Linux

### Build from Source

```sh
git clone https://github.com/inodaf/neoman.git
cd neoman
make all
```

The binaries will be created in `./bin/`:
- `./bin/nman` — CLI client
- `./bin/nmand` — Daemon server

Add to your PATH:

```sh
export PATH="$PATH:$(pwd)/bin"
```

## Quickstart

Start using Neoman immediately:

```sh
nman inodaf/neoman
```

This launches a terminal interface for browsing the Neoman documentation — using Neoman itself.

## Usage

### Open Documentation

```sh
nman a-github-user/repo-name # E.g. nman vercel/next.js
```

Opens an interactive terminal interface for navigating through the documentation.

### List Documentation

```sh
# All available projects
nman list

# Projects in an organization
nman list <author>

# Documents in a project
nman list <author>/<project>
```

### View a Specific Document

```sh
nman view author/project "path/to/document.md"
```

For detailed usage, run `nman inodaf/neoman` or visit the [documentation](docs/).

## Adopting Neoman

Make your project's documentation available through Neoman in few steps:

1. Create a `/docs` directory in your repository root
3. Add additional `.md` files for your documentation pages
2. Add an `index.md` file as your documentation front-page (optional)
4. Push to your Git provider (GitHub, GitLab, etc.)

**Example structure:**

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

Your documentation is now accessible via:

```sh
nman your-username/your-repo
```

For authoring guidelines, see [docs/Authoring Docs/](docs/Authoring%20Docs/).

## Contributing

We welcome contributions! Here's how to get started:

1. Review the [Architecture documentation](docs/Internals/Architecture.md)
2. Check out [Coding Conventions](docs/Internals/CodingConventions.md)

See [docs/Internals/](docs/Internals/) for more details.

## Get in Touch

- [GitHub Issues](https://github.com/inodaf/neoman/issues) — Report bugs or request features
- [GitHub Discussions](https://github.com/inodaf/neoman/discussions) — Ask questions and share ideas

## License

MIT License — see [LICENSE](LICENSE) for details.
