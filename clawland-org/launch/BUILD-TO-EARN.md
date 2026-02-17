# Clawland: Build the Edge AI Future. Get Paid for It.

**TL;DR**: We're building an open-source ecosystem of AI Agents that run on $10 hardware and replace $50,000/year monitoring jobs. Core contributors share 20% of all product revenue. We're looking for our first 10 Core Maintainers.

---

## The Problem

A night-shift data center operator costs $48,000/year to watch screens and restart crashed services. A fish pond patrol worker earns $30,000/year to check oxygen levels at 3am. A greenhouse attendant makes $18,000/year to open vents when it gets hot.

These are important jobs. They're also repetitive, exhausting, and ripe for automation.

But current IoT solutions are expensive ($5,000-50,000), require system integrators, and do dumb if-then alerting.

## The Solution: Clawland

We built a family of AI Agents that run on hardware as cheap as $10:

| Agent | Language | RAM | Hardware | Purpose |
|-------|----------|-----|----------|---------|
| **picclaw** | Go | <10MB | $10 RISC-V board | Edge sensing + local decision + instant action |
| **moltclaw** | TypeScript | >1GB | Cloud server | Multi-agent routing + global orchestration |
| **nanoclaw** | Python | >100MB | $50 SBC | Rich ecosystem, mid-weight tasks |
| **microclaw** | C/Rust | <1MB | $2 MCU | Sensor-level micro agent |

Combined with $5-50 sensors, these agents can:

- **Sense**: Read temperature, humidity, dissolved oxygen, vibration, smoke, motion (via `exec` + Python/GPIO)
- **Think**: LLM-powered multi-sensor analysis with trend prediction (8+ LLM providers)
- **Act**: Control relays, valves, switches (via `exec` + GPIO)
- **Report**: Alert via Telegram, Discord, Feishu, DingTalk, WhatsApp, QQ (7 channels)
- **Learn**: Persistent memory that gets smarter over time
- **Schedule**: Cron-based autonomous monitoring loops

## Real Numbers

| Kit | Hardware Cost | Replaces | Traditional Cost | Savings |
|-----|-------------|----------|-----------------|---------|
| Datacenter Sentinel | $88 | Night shift operator | $48,000/yr | 99% |
| Pond Guardian | $133/pond | Patrol worker | $30,000/yr | 99% |
| Greenhouse Manager | $95/greenhouse | Attendant | $18,000/yr | 98% |
| Safety Guardian | $30 | Elderly care subscription | $900/yr | 90% |
| Equipment Doctor | $32/machine | Predictive maintenance service | $3,200/yr | 95% |

Full analysis with BOM, implementation logic, and ROI: [clawland-kits](https://github.com/Clawland-AI/clawland-kits)

## Why Open Source? Why Revenue Share?

Because the best edge AI system will be built by people who actually deploy it in fish ponds, factories, and server rooms — not by a startup in a San Francisco office.

**We're not asking for donations. We're offering a deal:**

> You build it. It sells. You get 20%.

### How It Works

1. **20% of net product revenue** goes into a Contributor Revenue Pool
2. **Quarterly payouts** based on a transparent contribution points system
3. **Points** earned for: code, reviews, skills, hardware designs, docs, community support
4. **Open Collective** for full financial transparency
5. **First 10 Core Maintainers** get an additional 0.5% lifetime revenue share

### Example

If Clawland generates $50,000 net revenue in Q3:
- Contributor Pool: $10,000
- Developer who contributed 500 points out of 1000 total: **$5,000 for the quarter**

Read the full terms: [CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)

## What We Need

### Right Now (First 10 Core Maintainers)

- **Go developers**: Help build picclaw (edge agent, tools, channels)
- **TypeScript developers**: Help build moltclaw (cloud gateway, fleet management)
- **Python developers**: Help build nanoclaw + sensor driver libraries
- **Hardware hackers**: Design and test sensor kits, write driver scripts
- **Technical writers**: Documentation, tutorials, deployment guides

### Good First Issues

Every repo maintains `good-first-issue` labels. Start there.

### Bounty Board

We maintain a [Bounty Board](https://github.com/orgs/Clawland-AI/projects/1) with priced tasks. Claim one, ship it, get paid.

## Governance

Clawland is governed by a Technical Steering Committee (TSC) with transparent RFC-based decision making. See [GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md).

## Get Started

```bash
# Clone picclaw
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw

# Build
make build

# Run
./picclaw agent -m "Hello, I'm your new edge AI assistant"
```

## Links

- GitHub: [github.com/Clawland-AI](https://github.com/Clawland-AI)
- Revenue Share Terms: [CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)
- Governance: [GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md)
- Contributing: [CONTRIBUTING.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTING.md)

---

*"Let $10 hardware + AI Agents do the watching, so humans can do the thinking."*

**We're building the future where every fish pond, every server room, every greenhouse has an AI guardian. Join us.**
