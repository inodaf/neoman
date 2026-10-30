---
name: nman
description: Retrieve and read documentation directly from the terminal using `nman`. Provides access to both public open-source and private/internal organization documentation. Use when you need to consult project documentation and understand internal systems.
---

# nman - Documentation Retrieval

Retrieve and read documentation directly from the terminal.

`nman` provides access to both **public open-source documentation** and **private/internal organization documentation**. This makes it especially useful when working with internal systems, proprietary codebases, or integrating internal services with one another.

## Heuristics

- <owner/repo> follows the format of git-hosted repositories (e.g. `facebook/react`, `microsoft/typescript`, `myorg/myproject`)
  - E.g: https://github.com/**facebook/react**, https://gitlab.com/**myorg/myproject**
- Repositories must have a `docs/` directory at the root level to be recognized as documentation sources;
- Markdown is the supported format for documentation files;
- Mermaid diagrams and code blocks within markdown files are supported.

## Retrieving Documentation

If documentation is not yet installed, use the following:

```sh
# Retrieve documentation from a specific repository (public or private)
$ nman add <owner/repo>
```

Works for any git-hosted (GitHub) documentation, including private repositories the user has access to.

## Discovery

Start here to understand what documentation is available:

```sh
# List all installed documentation
$ nman list

# List all projects from a specific owner/organization
$ nman list <owner>

# List all pages within a project (with titles and paths)
$ nman list <owner/repo>
```

Output is YAML format.

## Reading Documentation

```sh
# Read a specific page by path (recommended for AI agents)
$ nman view <owner/repo> "<path>"
```

## Searching

```sh
# Search within a project's documentation
$ nman query <owner/repo> "<search terms>"
```

## Recommended Workflow

1. Use `nman list` to discover available documentation
2. Use `nman list <owner/repo>` to see specific pages with their paths
3. Use `nman view <owner/repo> "<path>"` to read a specific page
4. Use `nman query` to search for specific content/keywords
