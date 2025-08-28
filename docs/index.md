# Neoman

[![Static Badge](https://img.shields.io/badge/Docs-%24_nman_inodaf%2Fnman-black)](https://nman.local/inodaf/nman)

A modern documentation reader inspired by Unix `man` pages, designed to make software documentation accessible, searchable, and maintainable.

## Why Neoman?

Finding and reading documentation for your software stack shouldn't be a hassle. Neoman brings back the simplicity of Unix `man` pages while adding modern features for today's development workflows.

## Key Features

### 🚀 Zero-Deployment Documentation
Focus on writing great docs, not managing infrastructure. No servers to maintain, no hosting costs, no deployment pipelines. Just write Markdown and push to Git.

### 📚 Git-Native Workflow
Documentation stays in sync with your code automatically. Push to any Git provider (GitHub, GitLab, etc.) and readers get updates instantly.

### 📝 Convention Over Configuration
Simple structure: add a `/docs` directory with an `index.md` file and additional `.md` files. That's it.

### 🔍 Unified Documentation Hub
All your organization's documentation in one searchable place. Use full-text search across all docs or filter by specific projects.

```sh
nman inodaf/nman
```

### 🔒 Local & Secure
Everything runs locally on your machine. No data leaves your device, perfect for private organizational documentation.

### 🔗 Shareable Links
Share documentation with colleagues using local URLs like [https://nman.local/inodaf/nman](https://nman.local/inodaf/nman). Perfect for README badges:

[![Static Badge](https://img.shields.io/badge/Read_Docs-%24_nman_inodaf%2Fnman-black)](https://nman.local/inodaf/nman)

## Installation

```sh
curl https://raw.githubusercontent.com/inodaf/neoman/refs/heads/main/install.sh | bash
```

The installer will request permission to configure the `nman.local` domain in your `/etc/hosts` file for the best experience.

## Quick Start

View this documentation using Neoman itself:

```sh
nman inodaf/nman
```

## Documentation

Complete documentation is available in the [/docs](/docs) directory or through Neoman:

```sh
nman inodaf/nman
```
