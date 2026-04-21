---
name: nman
description: Read documentation directly from the terminal using nman. Provides access to both public open-source and private/internal organization documentation. Use when you need to consult project documentation, understand internal systems, or look up API references.
---

# nman - Documentation Reader

Read project documentation directly from the terminal. Use this skill to discover, navigate, and consult documentation installed via nman.

nman provides access to both **public open-source documentation** and **private/internal organization documentation**. This makes it especially useful when working with internal systems, proprietary codebases, or integrating internal services with one another.

## Discovery

Start here to understand what documentation is available:

```sh
# List all installed documentation
$ nman list

# List all projects from a specific author/organization  
$ nman list <author>

# List all pages within a project (with titles and paths)
$ nman list <author/repo>
```

Output is YAML format. Parse the `docs:` array for entries.

## Reading Documentation

```sh
# Read a specific page by path (recommended for AI agents)
$ nman view <author/repo> "<path>"
```

**Note**: The command `nman <author/repo>` opens an interactive TUI (Terminal User Interface) which is not suitable for AI agents. Always use `nman view` with a specific path instead.

## Searching

```sh
# Search within a project's documentation
$ nman query <author/repo> "<search terms>"
```

## Recommended Workflow

1. Run `nman list` to discover available documentation
2. Run `nman list <author/repo>` to see specific pages with their paths
3. Use `nman view <author/repo> "<path>"` to read a specific page
4. Use `nman query` if you need to search for specific content

## Retrieving Documentation

If documentation is not yet installed, use the `--add-only` flag to install it without opening the TUI:

```sh
# Install documentation (AI agent friendly - no TUI)
$ nman <author/repo> --add-only
```

This allows AI agents to autonomously install documentation as needed. Documentation must be installed before it can be listed or viewed.

**Alternative**: The command `nman <author/repo>` (without `--add-only`) downloads documentation but also opens an interactive TUI, which is not suitable for AI agents.

This works for any git-hosted documentation, including private repositories the user has access to.
