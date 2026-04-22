# Viewing Pages

The `nman view` command displays the content of a specific documentation page. Use it when you know exactly which page you want to read, or after discovering pages with `nman list`.

## Designed for AI Agents

The view command outputs **plain markdown** with no additional formatting, making it ideal for AI coding agents to:

- Retrieve specific documentation content for context
- Parse and analyze documentation programmatically
- Build knowledge from targeted pages rather than entire projects

Combine with `nman list <author/repo>` to discover available pages, then use `nman view` to retrieve the ones you need.

## Usage

### Basic Usage

```sh
$ nman view <author/repo> "<path>"
```

**Example:**

```sh
$ nman view inodaf/neoman "Authoring Docs/4. Publishing.md"
```

**Output:**

```markdown
# Publishing

Once you have completed your documentation, you can publish it...
```

### Paths Without Extension

You can omit the file extension and nman will resolve it automatically:

```sh
$ nman view inodaf/neoman "Authoring Docs/4. Publishing"
```

### Paths With Spaces

When the path contains spaces, wrap it in quotes:

```sh
$ nman view facebook/react "Getting Started/Installation"
```

## Understanding the Output

The output is the raw markdown content of the page, rendered exactly as stored in the source repository. No YAML frontmatter or metadata is included—just the document content.

## Common Workflows

### Reading a Specific Page

First discover what pages are available, then view one:

```sh
# List available pages
$ nman list inodaf/neoman

# View a specific page
$ nman view inodaf/neoman "Querying/Listing Documentation.md"
```

### Quick Navigation

When you already know the path:

```sh
$ nman view golang/go "doc/effective_go"
```

## Error Messages

**Documentation not found**

```
documentation not found
```

The specified `author/repo` doesn't have documentation installed. Add it first:

```sh
$ nman inodaf/neoman
```

**Page not found**

```
page not found
```

The specified path doesn't exist in the documentation. Check the path with:

```sh
$ nman list <author/repo>
```

## See Also

- **Listing Pages**: Use `nman list <author/repo>` to see all available pages
- **Reading Documentation**: Use `nman <author/repo>` to open documentation in browser
- **Searching**: Use `nman query <author/repo> "<search>"` to search within documentation
