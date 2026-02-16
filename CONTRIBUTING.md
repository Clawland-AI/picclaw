# Contributing to PicoClaw

Thank you for your interest in contributing to PicoClaw! 🦐

This guide will help you get started with development, testing, and submitting contributions.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Building and Running](#building-and-running)
- [Running Tests](#running-tests)
- [Code Style Guide](#code-style-guide)
- [Pull Request Process](#pull-request-process)
- [Project Structure](#project-structure)

## Development Environment Setup

### Prerequisites

- **Go 1.21+** (Go 1.22+ recommended for latest features)
- **Make** (for build automation)
- **Git** (for version control)

### Installation

1. **Install Go**

   Download and install Go from [golang.org](https://golang.org/dl/):

   ```bash
   # Verify installation
   go version
   ```

2. **Clone the repository**

   ```bash
   git clone https://github.com/Clawland-AI/picclaw.git
   cd picclaw
   ```

3. **Install dependencies**

   ```bash
   make deps
   ```

4. **Set up configuration**

   ```bash
   # Initialize PicoClaw
   picoclaw onboard
   
   # Edit config at ~/.picoclaw/config.json
   # Add your API keys for testing
   ```

## Building and Running

### Build from source

```bash
# Build for your current platform
make build

# Build for all platforms
make build-all

# Build and install to $GOPATH/bin
make install
```

The binary will be created in the `bin/` directory.

### Run locally

```bash
# Interactive chat mode
./bin/picoclaw agent

# Single message
./bin/picoclaw agent -m "Hello, world!"

# Start gateway (for chat apps)
./bin/picoclaw gateway
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/agent

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...
```

## Code Style Guide

### General Guidelines

1. **Follow Go conventions**
   - Use `gofmt` to format code
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use meaningful variable and function names

2. **Code formatting**

   ```bash
   # Format all code
   gofmt -w .
   
   # Run linter (if installed)
   golangci-lint run
   ```

3. **Comments**
   - Add comments for exported functions and types
   - Use complete sentences
   - Explain *why*, not just *what*

4. **Error handling**
   - Always check errors
   - Provide context when wrapping errors
   - Use descriptive error messages

### Example

```go
// Good
func ProcessMessage(msg string) error {
    if msg == "" {
        return fmt.Errorf("message cannot be empty")
    }
    
    result, err := someOperation(msg)
    if err != nil {
        return fmt.Errorf("failed to process message: %w", err)
    }
    
    return nil
}

// Avoid
func ProcessMessage(msg string) error {
    result, _ := someOperation(msg) // Don't ignore errors
    return nil
}
```

### Project Conventions

- **Package naming**: Use short, lowercase, single-word names
- **File naming**: Use lowercase with underscores (e.g., `message_bus.go`)
- **Test files**: Named `*_test.go` in the same package
- **Constants**: Use `CamelCase` for exported, `camelCase` for unexported

## Pull Request Process

### Before submitting

1. **Create an issue first** (if one doesn't exist)
   - Describe the problem or feature
   - Discuss your approach
   - Get feedback before investing time

2. **Fork and create a branch**

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

3. **Make your changes**
   - Write clear, focused commits
   - Add tests for new functionality
   - Update documentation if needed

4. **Test your changes**

   ```bash
   # Run tests
   go test ./...
   
   # Test locally
   make build
   ./bin/picoclaw agent -m "test message"
   ```

5. **Format and lint**

   ```bash
   gofmt -w .
   go vet ./...
   ```

### Commit messages

Use clear, descriptive commit messages following conventional commits:

```
<type>: <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build/tooling changes

**Examples:**

```
feat: add health check endpoint /healthz

docs: update installation instructions for Windows

fix: handle empty messages in Telegram handler

test: add unit tests for message bus
```

### Submitting the PR

1. **Push your branch**

   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request**
   - Use a descriptive title
   - Reference the issue: "Closes #123"
   - Describe what changed and why
   - Add screenshots if relevant (UI changes)

3. **PR checklist**
   - [ ] Code follows style guidelines
   - [ ] Tests pass locally
   - [ ] Added/updated tests for changes
   - [ ] Updated documentation if needed
   - [ ] Commit messages are clear
   - [ ] No merge conflicts

4. **Review process**
   - Address reviewer feedback
   - Push updates to the same branch
   - Be responsive and collaborative

## Project Structure

```
picclaw/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── agent/       # Agent logic
│   ├── config/      # Configuration
│   ├── gateway/     # Gateway server
│   ├── tools/       # Tool implementations
│   └── providers/   # LLM providers
├── assets/          # Images, logos
├── docs/            # Documentation
├── Makefile         # Build automation
└── README.md        # Project overview
```

## Getting Help

- **Discord**: Join our [Discord server](https://discord.gg/V4sAZ9XWpN)
- **Issues**: Check existing [GitHub issues](https://github.com/Clawland-AI/picclaw/issues)
- **WeChat**: See QR code in README

## Code of Conduct

Be respectful, collaborative, and constructive. We're building something cool together! 🦐

---

**Ready to contribute?** Pick a [good first issue](https://github.com/Clawland-AI/picclaw/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) and get started!
