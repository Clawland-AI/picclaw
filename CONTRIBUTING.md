# Contributing to PicoClaw

PRs are welcome. The codebase is intentionally small (~7500 lines of Go).

## Requirements

- Go 1.24+
- `make` for common workflows

## Development

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
make test
make lint
make build-all
```

## Pull requests

- One logical change per PR when possible
- Run `make test` and `make lint` before submitting
- Describe how you tested gateway, tools, or channel changes
