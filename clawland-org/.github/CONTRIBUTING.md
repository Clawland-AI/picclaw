# Contributing to Clawland

Thank you for your interest in contributing to Clawland! This guide applies to all repositories under the `github.com/Clawland-AI` organization.

---

## Quick Start

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests (if applicable)
5. Commit with a clear message: `git commit -m "feat: add temperature sensor support for SHT40"`
6. Push and open a Pull Request
7. Sign the CLA (automated bot will prompt you on first PR)

---

## Ways to Contribute

### Code

- Fix bugs (look for [`bug`](https://github.com/search?q=org%3AClawland-AI+label%3Abug&type=issues) labels)
- Implement features (look for [`enhancement`](https://github.com/search?q=org%3AClawland-AI+label%3Aenhancement&type=issues) labels)
- New contributors: look for [`good first issue`](https://github.com/search?q=org%3AClawland-AI+label%3A%22good+first+issue%22&type=issues)

### Skills

- Create new Skills for [clawland-skills](https://github.com/Clawland-AI/clawland-skills)
- Each Skill is a directory with a `SKILL.md` file (YAML frontmatter + markdown body)
- See existing Skills for examples

### Hardware Kit Designs

- Design new sensor kits for [clawland-kits](https://github.com/Clawland-AI/clawland-kits)
- Include: BOM, wiring diagram, driver scripts, PicClaw Skill, cost analysis
- See existing kits for the expected format

### Documentation

- Improve README files, add tutorials, translate docs
- Fix typos, clarify confusing sections

### Community Support

- Answer questions in GitHub Discussions and Discord
- Help triage issues

### Bounties

- Check the [Bounty Board](https://github.com/orgs/Clawland-AI/projects/1) for paid tasks
- Claim a bounty by commenting on the issue
- Submit your work as a PR referencing the bounty issue

---

## Contribution Points

Every contribution earns points that determine your share of the quarterly Contributor Revenue Pool. See [CONTRIBUTOR-REVENUE-SHARE.md](CONTRIBUTOR-REVENUE-SHARE.md) for the full points table and distribution formula.

---

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): description

feat(agent): add sub-agent timeout configuration
fix(channels): handle Telegram reconnection on network loss
docs(skills): add weather skill tutorial
chore(ci): update Go version to 1.24
test(tools): add exec tool safety guard tests
```

Types: `feat`, `fix`, `docs`, `chore`, `test`, `refactor`, `perf`, `style`, `ci`

---

## Pull Request Guidelines

1. **One PR per feature/fix** — Keep PRs focused and reviewable
2. **Description** — Explain what, why, and how. Reference related issues
3. **Tests** — Add tests for new features, verify existing tests pass
4. **Documentation** — Update relevant docs if behavior changes
5. **Size** — Aim for <400 lines changed. Split larger changes into a series of PRs

### PR Template

Every repository has a PR template. Fill it out completely.

---

## Code Review

- All PRs require at least 1 approving review from a Core Maintainer
- Reviewers should respond within 72 hours
- Be constructive and specific in review comments
- Approve when satisfied; request changes with clear explanation

---

## Development Setup

### Picclaw (Go)

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
make build
./picclaw status
```

### Moltclaw (TypeScript)

```bash
git clone https://github.com/Clawland-AI/moltclaw.git
cd moltclaw
npm install
npm run dev
```

### Skills

```bash
# Install a skill for testing
picclaw skills install /path/to/your-skill

# Test it
picclaw agent -m "test your skill functionality"
```

---

## Reporting Issues

- Use the issue templates provided in each repository
- Include: steps to reproduce, expected behavior, actual behavior, environment details
- For security vulnerabilities, **do NOT open a public issue** — email security@clawland.dev

---

## Code of Conduct

All contributors must follow our [Code of Conduct](CODE_OF_CONDUCT.md). Be kind, be constructive, be welcoming.

---

## Questions?

- [GitHub Discussions](https://github.com/orgs/Clawland-AI/discussions) — For questions and ideas
- [Discord](https://discord.gg/clawland-ai) — For real-time chat

---

*Thank you for helping build the Edge AI Future!*
