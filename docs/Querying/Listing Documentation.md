# Listing Documentation

The `nman list` command helps you discover what documentation is installed on your local machine. Whether you want to see everything available, explore a specific author's projects, or view all pages within a project, the list command provides three levels of granularity.

## Designed for AI Agents

The list command outputs **structured YAML** format, making it ideal for AI coding agents to:

- Discover what documentation is available to consult
- Understand the scope of installed knowledge bases
- Programmatically navigate documentation to gain context

Combine with `nman view` to retrieve specific pages, or `nman query` to search within documentation.

## Usage

### List All Documentation

```sh
$ nman list
```

Shows all documentation installed across all authors and organizations. Perfect for discovering what's available on your system.

**Example output:**

```yaml
docs:
  - facebook/react
  - golang/go
  - inodaf/neoman
  - microsoft/typescript
```

Each entry follows the `author/repo` format and represents a documentation project you can explore.

### List by Author

```sh
$ nman list <author>
```

Shows all documentation projects from a specific author or organization.

**Example:**

```sh
$ nman list golang
```

**Output:**

```yaml
author: golang
docs:
  - go
  - tools
  - tour
```

This lists all the projects under the `golang` organization that have documentation installed.

### List Project Pages

```sh
$ nman list <author/repo>
```

Shows all documentation pages available within a specific project, including their titles and paths.

**Example:**

```sh
$ nman list inodaf/neoman
```

**Output:**

```yaml
documentation: inodaf/neoman
pages:
  - Table of Contents: "Authoring Docs/1. Table of Contents.md"
  - Best Practices: "Authoring Docs/2. Best Practices.md"
  - Publishing: "Authoring Docs/4. Publishing.md"
  - Private Repositories: "Teams/1. Private Repositories.md"
  - Privacy and Security: "Teams/2. Privacy and Security.md"
```

Each entry shows the document title followed by its path within the project.

## Understanding the Output

All list commands return data in YAML format for easy readability and potential scripting use.

The output directly maps to other `nman` commands:

- **From `nman list`**: Use `author/repo` with `nman <author/repo>` to read documentation
- **From `nman list <author>`**: Combine with author to form `nman <author/repo>`
- **From `nman list <author/repo>`**: Use the path with `nman view <author/repo> "<path>"`

## Common Workflows

### Discovering What's Installed

Start broad and narrow down:

```sh
# See everything
$ nman list

# Pick an author
$ nman list golang

# Pick a project
$ nman golang/go
```

### Finding a Specific Document

When you know the project but not the exact page:

```sh
# List all pages
$ nman list facebook/react

# View a specific page
$ nman view facebook/react "Getting Started/Installation.md"
```

### Exploring New Documentation

After adding new documentation:

```sh
# Add documentation
$ nman microsoft/typescript

# Verify it's installed
$ nman list microsoft

# Browse available pages
$ nman list microsoft/typescript
```

## Error Messages

**No documentation found**

```
No documentation found. Add by running 'nman author/repo'
```

This means you don't have any documentation installed yet. Add some using:

```sh
$ nman golang/go
```

**Author not found**

```
No documentation found for author "unknown-author"
```

The specified author/organization doesn't have any installed documentation. Double-check the spelling or add documentation for that author.

**Project not found**

```
No documentation found for project "author/unknown-repo"
```

The specified project isn't installed. Verify the project name or add it:

```sh
$ nman author/repo-name
```

## See Also

- **Viewing Documentation**: Use `nman <author/repo>` to read a project's main documentation
- **Viewing Specific Pages**: Use `nman view <author/repo> "<path>"` to read a specific page
- **Adding Documentation**: Simply run `nman <author/repo>` to add new documentation
- **Searching**: Use `nman query <author/repo> "<search terms>"` to search within documentation
