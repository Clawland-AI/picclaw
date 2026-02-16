# Contributing to PicoClaw

Thank you for your interest in contributing to PicoClaw! This guide will help you get started with development.

---

## 📋 Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Building and Running Locally](#building-and-running-locally)
- [Running Tests](#running-tests)
- [Code Style Guide](#code-style-guide)
- [Pull Request Process](#pull-request-process)
- [Project Structure](#project-structure)
- [Common Tasks](#common-tasks)

---

## 🛠️ Development Environment Setup

### Prerequisites

**Required:**
- **Go 1.24+** - [Download](https://go.dev/dl/)
- **Git** - For version control
- **Make** - For build automation

**Optional:**
- **Docker** - For containerized deployment testing
- **golangci-lint** - For linting (install via `make deps`)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
```

2. **Install dependencies**

```bash
make deps
```

This installs:
- `golangci-lint` - Code linter
- Go module dependencies

3. **Verify installation**

```bash
go version  # Should be 1.24 or higher
make --version
```

---

## 🏗️ Building and Running Locally

### Quick Start

```bash
# Build the binary
make build

# Binary will be in build/ directory
./build/picoclaw --version
```

### Run from source (no build)

```bash
# Initialize configuration
make init

# Edit config
nano ~/.picoclaw/config.json

# Run directly
go run cmd/picoclaw/main.go
```

### Build for multiple platforms

```bash
# Build for all supported platforms
make build-all

# Outputs:
# - build/picoclaw-linux-amd64
# - build/picoclaw-linux-arm64
# - build/picoclaw-linux-riscv64
# - build/picoclaw-darwin-amd64
# - build/picoclaw-darwin-arm64
```

### Install system-wide

```bash
# Install to ~/.local/bin (default)
make install

# Or specify prefix
make install INSTALL_PREFIX=/usr/local
```

---

## 🧪 Running Tests

### Run all tests

```bash
make test
```

### Run specific package tests

```bash
# Test message bus
go test -v ./pkg/bus

# With coverage
go test -v -cover ./pkg/bus

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run tests with race detection

```bash
go test -race ./...
```

### Benchmarks

```bash
go test -bench=. ./pkg/bus
```

---

## 📝 Code Style Guide

### Go Conventions

Follow standard Go conventions:
- Use `gofmt` for formatting (automatic via `make fmt`)
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Write comments for exported functions

### Linting

```bash
# Run linter
make lint

# Auto-fix issues
make lint-fix
```

### Code Organization

```
pkg/
  agent/     - Core agent logic
  bus/       - Message bus
  channels/  - Chat integrations (Telegram, Discord, etc.)
  config/    - Configuration management
  providers/ - LLM providers
  tools/     - Tool implementations
  ...
cmd/
  picoclaw/  - Main entry point
skills/      - Built-in agent skills
```

### Naming Conventions

- **Packages**: lowercase, single word (`bus`, `agent`)
- **Files**: lowercase with underscores (`message_bus.go`)
- **Types**: PascalCase (`MessageBus`, `InboundMessage`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase or ALL_CAPS for package-level

### Error Handling

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to publish message: %w", err)
}

// ❌ Bad
if err != nil {
    panic(err) // Don't panic in libraries
}
```

### Testing

- Test files: `*_test.go` in same package
- Table-driven tests preferred
- Use `t.Helper()` for test helpers
- Aim for >80% coverage on new code

Example:
```go
func TestMessageBus(t *testing.T) {
    tests := []struct {
        name string
        want int
    }{
        {"case 1", 1},
        {"case 2", 2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

---

## 🔄 Pull Request Process

### 1. Fork and Clone

```bash
# Fork on GitHub, then:
git clone https://github.com/YOUR_USERNAME/picoclaw.git
cd picoclaw
git remote add upstream https://github.com/sipeed/picoclaw.git
```

### 2. Create a Branch

Use descriptive branch names:

```bash
# Feature
git checkout -b feat/add-whatsapp-channel

# Bug fix
git checkout -b fix/memory-leak-in-bus

# Documentation
git checkout -b docs/update-readme

# Tests
git checkout -b test/add-agent-tests
```

### 3. Make Changes

- Write clear, focused commits
- Follow code style guide
- Add tests for new functionality
- Update documentation

### 4. Test Locally

```bash
# Lint
make lint

# Test
make test

# Build
make build
```

### 5. Commit

Use [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat: add WhatsApp channel integration"
git commit -m "fix: resolve memory leak in message bus"
git commit -m "docs: update contributing guide"
git commit -m "test: add unit tests for agent package"
```

### 6. Push and Open PR

```bash
git push origin feat/add-whatsapp-channel
```

Then open a PR on GitHub with:
- **Title**: Clear, descriptive (e.g., "feat: add WhatsApp channel")
- **Description**: What changed, why, how to test
- **Link to issue**: "Closes #123" if applicable

### PR Checklist

- [ ] Code compiles without warnings
- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow convention
- [ ] Branch is up to date with `main`

---

## 📁 Project Structure

```
picoclaw/
├── cmd/
│   └── picoclaw/         # Main application entry point
│       └── main.go
├── pkg/                  # Public packages (importable)
│   ├── agent/            # Agent orchestration
│   ├── bus/              # Message bus
│   ├── channels/         # Chat integrations
│   │   ├── telegram.go
│   │   ├── discord.go
│   │   └── ...
│   ├── config/           # Configuration
│   ├── providers/        # LLM providers
│   ├── tools/            # Built-in tools
│   └── ...
├── skills/               # Agent skills (scripts, configs)
├── assets/               # Images, assets
├── Makefile              # Build automation
├── go.mod                # Go module definition
├── config.example.json   # Example configuration
├── Dockerfile            # Container image
└── README.md             # Project documentation
```

---

## 🔧 Common Tasks

### Add a new channel

1. Create `pkg/channels/yourchannel.go`
2. Implement `Channel` interface
3. Register in `pkg/channels/manager.go`
4. Add config schema
5. Write tests in `yourchannel_test.go`
6. Update docs

### Add a new tool

1. Create `pkg/tools/yourtool.go`
2. Implement tool logic
3. Register in tool manager
4. Write tests
5. Update tool documentation

### Add a new skill

1. Create directory: `skills/yourskill/`
2. Add `SKILL.md` with description
3. Add scripts/configs
4. Test with `picoclaw skill run yourskill`

### Update dependencies

```bash
# Update all
go get -u ./...
go mod tidy

# Update specific package
go get -u github.com/example/package@latest
```

### Debug

```bash
# Enable debug logs
export PICOCLAW_LOG_LEVEL=debug
./build/picoclaw

# Or in config.json
{
  "logging": {
    "level": "debug"
  }
}
```

---

## 🐛 Reporting Issues

Found a bug? Please [open an issue](https://github.com/sipeed/picoclaw/issues) with:

- **Description**: What happened vs. what you expected
- **Steps to reproduce**
- **Environment**: OS, Go version, PicoClaw version
- **Logs**: Relevant error messages

---

## 💬 Community

- **Discord**: [Join our server](https://discord.gg/clawd)
- **GitHub Discussions**: For Q&A and feature requests
- **Issue Tracker**: For bug reports

---

## 📜 License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

---

**Thank you for contributing to PicoClaw! 🦐**
