# Contributing to PicoClaw

Thank you for your interest in contributing to PicoClaw! 🦐

We welcome contributions of all kinds: bug reports, feature requests, documentation improvements, and code contributions.

## 📋 Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Building from Source](#building-from-source)
- [Running Tests](#running-tests)
- [Code Style Guide](#code-style-guide)
- [Pull Request Process](#pull-request-process)
- [Project Structure](#project-structure)

## 🛠️ Development Environment Setup

### Prerequisites

- **Go 1.21 or higher** (we recommend Go 1.22+)
- **Git** for version control
- **Make** (optional, for convenient build commands)

### Installation

1. **Install Go**

   Download and install from [https://go.dev/dl/](https://go.dev/dl/)

   Verify installation:
   ```bash
   go version
   ```

2. **Clone the repository**

   ```bash
   git clone https://github.com/sipeed/picoclaw.git
   cd picoclaw
   ```

3. **Install dependencies**

   ```bash
   make deps
   # or manually:
   go get -u ./...
   go mod tidy
   ```

## 🔨 Building from Source

### Build for your platform

```bash
make build
```

This creates a binary in `build/picoclaw-<platform>-<arch>` and a symlink at `build/picoclaw`.

### Build for all platforms

```bash
make build-all
```

This builds for:
- Linux (amd64, arm64, riscv64)
- Windows (amd64)

### Run directly without installing

```bash
make run ARGS="agent -m 'Hello!'"
```

### Install to system

```bash
make install
```

This installs:
- Binary to `~/.local/bin/picoclaw` (or `$INSTALL_PREFIX/bin`)
- Builtin skills to `~/.picoclaw/workspace/skills/`

You can customize the install prefix:
```bash
make install INSTALL_PREFIX=/usr/local
```

### Uninstall

```bash
make uninstall          # Remove binary only
make uninstall-all      # Remove binary + workspace + config
```

## 🧪 Running Tests

```bash
make test
# or
go test ./...
```

Run tests with race detector:
```bash
go test -race ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## 📝 Code Style Guide

### General Guidelines

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code (or `make fmt`)
- Write clear, descriptive commit messages
- Add comments for exported functions and complex logic
- Keep functions small and focused

### Formatting

Before submitting a PR, format your code:

```bash
make fmt
# or
go fmt ./...
```

### Naming Conventions

- **Packages**: lowercase, single-word names (e.g., `agent`, `tools`, `channels`)
- **Files**: lowercase with underscores (e.g., `agent_loop.go`, `message_bus.go`)
- **Types**: PascalCase (e.g., `AgentLoop`, `MessageBus`)
- **Functions**: camelCase for private, PascalCase for exported (e.g., `processMessage()`, `NewAgentLoop()`)
- **Constants**: UPPER_CASE or PascalCase depending on visibility

### Error Handling

- Always check errors explicitly
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Return early on errors

### Logging

Use the `logger` package for consistent logging:

```go
logger.InfoC("component", "message")
logger.ErrorC("component", "error occurred", err)
```

## 🚀 Pull Request Process

### 1. Fork and Create a Branch

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR_USERNAME/picoclaw.git
cd picoclaw
git checkout -b feat/your-feature-name
```

Branch naming conventions:
- `feat/feature-name` for new features
- `fix/bug-name` for bug fixes
- `docs/topic` for documentation
- `test/scope` for tests
- `refactor/scope` for refactoring

### 2. Make Your Changes

- Write clear, focused commits
- Follow the code style guide
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
make build
make test
./build/picoclaw agent -m "Test message"
```

### 4. Commit Your Changes

Write clear commit messages:

```
feat: add health check endpoint /healthz

- Implement GET /healthz route
- Return JSON with status, uptime, version
- Add tests for health check

Fixes #123
```

Commit message format:
- First line: `<type>: <short description>` (max 50 chars)
- Blank line
- Detailed description (wrap at 72 chars)
- Reference issues/PRs

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `chore`

### 5. Push and Create Pull Request

```bash
git push origin feat/your-feature-name
```

Then create a PR on GitHub with:
- Clear title and description
- Reference related issues (e.g., "Closes #123")
- Screenshots/examples if applicable
- Test results

### 6. Code Review

- Address review feedback promptly
- Keep discussions respectful and constructive
- Update your branch if main has changed:
  ```bash
  git fetch upstream
  git rebase upstream/main
  git push --force-with-lease
  ```

## 📁 Project Structure

```
picoclaw/
├── cmd/
│   └── picoclaw/
│       └── main.go          # CLI entry point
├── pkg/
│   ├── agent/              # Agent loop and logic
│   ├── bus/                # Message bus
│   ├── channels/           # Chat integrations (Telegram, Discord, etc.)
│   ├── config/             # Configuration management
│   ├── cron/               # Scheduled tasks
│   ├── logger/             # Logging utilities
│   ├── providers/          # LLM provider integrations
│   ├── skills/             # Skill system
│   ├── tools/              # Tool implementations
│   └── voice/              # Voice transcription
├── skills/                 # Builtin skills
├── assets/                 # Images and media
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Build automation
├── README.md               # Main documentation
└── CONTRIBUTING.md         # This file

```

### Key Packages

- **`cmd/picoclaw`**: Main CLI application
- **`pkg/agent`**: Core agent loop and message processing
- **`pkg/tools`**: Tool implementations (exec, web, file operations)
- **`pkg/channels`**: Chat platform integrations
- **`pkg/providers`**: LLM provider adapters
- **`pkg/skills`**: Skill loading and management

## 🐛 Reporting Bugs

When reporting bugs, include:

1. **PicoClaw version**: `picoclaw version`
2. **Go version**: `go version`
3. **Operating system**: `uname -a` (Linux/macOS) or OS info (Windows)
4. **Steps to reproduce**
5. **Expected behavior**
6. **Actual behavior**
7. **Logs** (if applicable)

## 💡 Feature Requests

We welcome feature requests! Please:

1. Check if the feature has already been requested
2. Describe the use case clearly
3. Explain why it would benefit PicoClaw users
4. Provide examples or mockups if possible

## 📜 License

By contributing to PicoClaw, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You

Every contribution helps make PicoClaw better. Thank you for being part of the community! 🦐

---

**Questions?** Open an issue or join our [Discord](https://discord.gg/V4sAZ9XWpN)
