# Contributing to octo-linter

Thanks for your interest in improving **octo-linter**! This guide outlines how to get started with local development and contribute effectively to the project.

## Getting Started with Development

### Required Tools

To contribute to the project, ensure you have the following installed:

1. **Go** (v1.22 or newer) – Primary language used throughout the codebase
2. **Docker** – Used for container builds and integration testing
3. **Git** – For source control and collaboration

### Initial Setup

To get your development environment ready:

1. Clone the repository:

   ```bash
   git clone https://github.com/keenbytes/octo-linter.git
   cd octo-linter
   ```

2. Fetch Go module dependencies:

   ```bash
   go mod download
   ```

3. *(Optional)* Set up the documentation tooling:

   ```bash
   python3 -m venv venv
   source venv/bin/activate
   pip install mkdocs-material
   ```

## Repository Layout

Here’s a quick overview of the project structure:

```
octo-linter/
├── cmd/               # CLI entry points
├── pkg/               # Core logic and utilities
├── tests/             # Integration and end-to-end tests
├── docs/              # Project documentation
├── argocd-config/     # Argo CD setup files
└── examples/          # Sample use cases and test inputs
```

## Building the Project

### Compiling the Go Binary

To build the executable:

```bash
cd cmd/octo-linter
go build .
```

This will produce a binary under the `bin/` directory.

### Building a Docker Image

To create a Docker image:

```bash
docker build -t octo-linter .
```

## Running Tests

### Executing Unit Tests

Unit tests are written using Go’s standard testing tools.

To run them all:

```bash
go test ./...
```

Add verbose output with:

```bash
go test -v ./...
```

## Documentation Workflow

The project’s documentation is powered by MkDocs. To preview it locally:

```bash
mkdocs serve
```

Your browser should automatically open the local site after that.

## Licensing

By contributing code or content to this repository, you agree that your contributions will be licensed under the terms outlined in the [`LICENSE`](./LICENSE) file.

---

Let me know if you'd like to adjust the tone (e.g., more casual or more corporate) or include a code of conduct section.
