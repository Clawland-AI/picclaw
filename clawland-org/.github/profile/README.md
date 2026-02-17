# Clawland

**Build the Edge AI Future. Get Paid for It.**

Clawland is an open-source ecosystem of AI Agents designed to run on hardware ranging from $2 microcontrollers to cloud servers — replacing human monitoring, maintenance, and service jobs at 1/100th the cost.

---

## The Claw Family

| Project | Language | RAM | Hardware Cost | Purpose |
|---------|----------|-----|---------------|---------|
| [**picclaw**](https://github.com/Clawland-AI/picclaw) | Go | <10MB | $10 board | Ultra-lightweight edge AI Agent |
| [**moltclaw**](https://github.com/Clawland-AI/moltclaw) | TypeScript | >1GB | Cloud / Mac Mini | Full-featured cloud AI Gateway with multi-Agent routing |
| [**nanoclaw**](https://github.com/Clawland-AI/nanoclaw) | Python | >100MB | $50 SBC | Mid-weight Agent with rich Python ecosystem |
| [**microclaw**](https://github.com/Clawland-AI/microclaw) | C / Rust | <1MB | $2-5 MCU | Sensor-level micro Agent |

## Ecosystem — OpenClaw Integration

Clawland is built to extend the [**OpenClaw**](https://github.com/openclaw/openclaw) open-source AI assistant ecosystem. OpenClaw serves as the **cloud brain**, while Clawland agents act as the **edge hands and feet**.

```
  ☁️ OpenClaw (Cloud Brain)          🤖 Clawland (Edge Hands & Feet)
  ┌──────────────────────┐          ┌──────────────────────────────┐
  │ openclaw/openclaw     │  ◄────► │ picclaw   — $10 Edge Agent   │
  │ 196k+ ⭐ MIT          │  sync   │ nanoclaw  — $50 SBC Gateway  │
  │ Multi-channel Gateway │  ────►  │ microclaw — $2 MCU Sensor    │
  └──────────────────────┘          │ moltclaw  — Cloud Gateway    │
         ▲                          └──────────────────────────────┘
         │ fork
  ┌──────┴───────────────┐
  │ Clawland-AI/openclaw  │ ← Tracked fork for upstream sync
  │ Edge integration      │   and contribution bridge
  └──────────────────────┘
```

| Project | Role |
|---------|------|
| [**openclaw**](https://github.com/Clawland-AI/openclaw) _(fork)_ | Upstream tracking & contribution bridge to [openclaw/openclaw](https://github.com/openclaw/openclaw) |
| [**moltclaw**](https://github.com/Clawland-AI/moltclaw) | Clawland's own L3 Cloud Gateway — Fleet integration, Dashboard, Edge API |

## Infrastructure

| Project | Purpose |
|---------|---------|
| [**clawland-fleet**](https://github.com/Clawland-AI/clawland-fleet) | Cloud-Edge orchestration — Fleet Manager, Edge API, Reporter |
| [**clawland-deploy**](https://github.com/Clawland-AI/clawland-deploy) | One-click deployment — Ansible Playbooks, Docker Compose |
| [**clawland-kits**](https://github.com/Clawland-AI/clawland-kits) | Hardware sensor kit designs — BOMs, wiring, drivers |
| [**clawland-skills**](https://github.com/Clawland-AI/clawland-skills) | Community skill marketplace — plug-and-play AI capabilities |
| [**clawland-drivers**](https://github.com/Clawland-AI/clawland-drivers) | Unified sensor driver library — GPIO/I2C/SPI/Serial |
| [**clawland-dashboard**](https://github.com/Clawland-AI/clawland-dashboard) | Fleet visualization dashboard — real-time node map + alerts |

## What Can You Build?

Deploy a **$10 board + $5 sensor + Picclaw** and replace:

- **$48,000/yr** data center night shift operator — with an $88 monitoring kit
- **$30,000/yr** fish pond patrol worker — with a $133 aquaculture guardian
- **$18,000/yr** greenhouse attendant — with a $95 smart farming kit
- **$900/yr** elderly care subscription — with a $30 safety guardian

> Every scenario above is documented with full hardware BOM, wiring diagrams, PicoClaw implementation logic, and cost comparison. See [clawland-kits](https://github.com/Clawland-AI/clawland-kits).

## Build to Earn

Clawland is not donation-ware. **Core contributors share 20% of all product revenue** — hardware kit sales, SaaS subscriptions, and consulting income.

- Transparent quarterly payouts via [Open Collective](https://opencollective.com/Clawland-AI)
- Contribution-based scoring system (code, reviews, skills, hardware designs, docs)
- First 10 Core Maintainers get an additional 0.5% lifetime revenue share

Read the full terms: [CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)

## Get Involved

1. Pick a repo and look for [`good first issue`](https://github.com/search?q=org%3AClawland-AI+label%3A%22good+first+issue%22&type=issues) labels
2. Read our [Contributing Guide](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTING.md)
3. Join the conversation in [GitHub Discussions](https://github.com/orgs/Clawland-AI/discussions)
4. Check the [Bounty Board](https://github.com/orgs/Clawland-AI/projects/1) for paid tasks

## Governance

Clawland is governed by a Technical Steering Committee (TSC). Major decisions go through an RFC process. See [GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md).

## License

Core agents (picclaw, moltclaw, nanoclaw, microclaw): **Apache 2.0**
Cloud orchestration (clawland-fleet): **BSL 1.1** (converts to Apache 2.0 after 4 years)
Hardware designs (clawland-kits): **CERN-OHL-S-2.0**
Skills (clawland-skills): **MIT**

---

*"Let $10 hardware + AI Agents do the watching, so humans can do the thinking."*
