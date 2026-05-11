# Contributing to PicoClaw

Thanks for helping improve PicoClaw. This guide covers the local workflow expected for small fixes, tests, docs, and feature PRs.

## Development setup

Requirements:

- Go 1.24 or newer (matching `go.mod`)
- Git
- Make on Linux/macOS, or a shell that can run the equivalent `go` commands on Windows
- Docker, only if you are changing container packaging

Clone the repository and download dependencies:

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
go mod download
```

## Build and run locally

Build the PicoClaw binary:

```bash
go build ./cmd/picoclaw
# or
make build
```

Initialize a local workspace and config:

```bash
./picoclaw onboard
```

Run a one-shot agent request or the gateway:

```bash
./picoclaw agent -m "What is 2+2?"
./picoclaw gateway
```

On Windows PowerShell, use `picoclaw.exe` after building.

## Tests and checks

Run the full test suite before opening a PR:

```bash
go test ./...
go vet ./...
go build ./...
```

For bus changes, also check coverage for the bus package:

```bash
go test ./pkg/bus -cover
```

If you modify the Dockerfile, build the image locally when Docker is available:

```bash
docker build -t picoclaw:dev .
```

## Code style

- Format Go code with `gofmt` or `go fmt ./...`.
- Keep changes focused on one issue or behavior at a time.
- Prefer clear package-level tests over broad rewrites.
- Do not commit local config, API keys, `.env` files, build artifacts, or generated coverage files.
- Keep public APIs small unless an issue explicitly asks for a new API.

## Pull request process

1. Create a branch from the latest `main`.
2. Make the smallest change that satisfies the issue requirements.
3. Add or update tests for behavior changes.
4. Run the checks listed above and include the results in the PR body.
5. Link the issue with `Closes #<issue-number>` when the PR fully resolves it.
6. Keep the PR description concise: what changed, why, and how it was verified.

Maintainers may ask for follow-up changes. Please address review feedback with focused commits rather than force-pushing unrelated rewrites.