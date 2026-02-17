# Clawland Contributor Revenue Share Agreement

This document defines how Clawland shares product revenue with its contributors. It is a binding commitment from Clawland to the developer community.

---

## 1. Core Commitment

**Clawland commits to distributing 20% of net product revenue to contributors every quarter.**

This is not a donation, a tip, or a discretionary bonus. It is a structured, transparent, formula-based distribution that contributors can rely on.

---

## 2. Revenue Sources

The Contributor Revenue Pool is funded by the following revenue streams:

| Source | Description |
|--------|-------------|
| Hardware Kit Sales | Sensor kits (e.g., Datacenter Sentinel, Pond Guardian, Greenhouse Manager) |
| SaaS Subscriptions | Clawland Fleet cloud services (dashboard, analytics, multi-node management) |
| Support & Consulting | Paid deployment assistance, custom Skill development, enterprise support |
| Licensing Fees | Commercial licenses for BSL-licensed components (clawland-fleet) |

---

## 3. Revenue Allocation

| Category | Percentage | Description |
|----------|------------|-------------|
| Cost of Goods Sold (COGS) | ~40-50% | Hardware procurement, servers, LLM API costs, shipping |
| Company Reserve | 25-30% | R&D, marketing, legal, operational expenses |
| **Contributor Revenue Pool** | **20%** | Distributed to eligible contributors |
| Community Fund | 5% | Bounties, events, documentation, community programs |

The Pool percentage (20%) is calculated on **net revenue** (gross revenue minus COGS).

---

## 4. Eligibility

### 4.1 Who is Eligible

| Tier | Requirement | Pool Access |
|------|-------------|-------------|
| **Core Maintainer** | TSC-approved, 6+ months sustained core contributions | Full Pool share |
| **Active Contributor** | Merged PRs in 3+ consecutive months | Full Pool share |
| **Community Contributor** | Any merged PR | Bounty rewards only (from Community Fund) |

### 4.2 Eligibility Period

- Contributions are counted per **calendar quarter** (Q1: Jan-Mar, Q2: Apr-Jun, Q3: Jul-Sep, Q4: Oct-Dec)
- A contributor must have at least **1 merged PR** in the quarter to be eligible for that quarter's distribution
- Core Maintainers retain eligibility for **1 additional quarter** after their last contribution (grace period)

---

## 5. Contribution Points System

### 5.1 Point Categories

| Contribution Type | Base Points | Multipliers |
|-------------------|-------------|-------------|
| Code merged (PR) | 10 / PR | Core module: x2, Bug fix: x1.5, Docs: x0.8 |
| Code review | 3 / Review | In-depth review with suggestions: x2 |
| Issue report (valid) | 2 / Issue | With reproduction steps: x1.5, Security vulnerability: x3 |
| Skill contribution | 15 / Skill | With test cases: x1.5 |
| Hardware kit design | 20 / Kit | With test report: x1.5 |
| Documentation / Tutorial | 5 / Article | Multi-language: x1.5 |
| Community support | 1 / Answer | High-quality (Core Maintainer certified): x3 |

### 5.2 Scoring Rules

- Points are calculated per quarter
- A single PR can earn points in only one category (highest applicable)
- Multipliers stack multiplicatively (e.g., Core module bug fix = 10 x 2 x 1.5 = 30 points)
- TSC may apply a **quality adjustment** of 0.5x to 2.0x based on impact assessment
- Automated/trivial contributions (typo fixes, formatting) are capped at 2 points each

### 5.3 Tracking

- Primary data source: GitHub API (PRs, reviews, issues)
- Supplementary: Manual tracking for community support, hardware testing
- Points ledger published in `reports/` directory each quarter
- 7-day dispute window after ledger publication

---

## 6. Distribution Formula

```
Individual Payout = Pool Total x (Individual Points / Sum of All Eligible Points)
```

### Example

Quarterly net revenue: $50,000
Pool (20%): $10,000

| Contributor | Points | Share | Payout |
|-------------|--------|-------|--------|
| Developer A | 500 | 50% | $5,000 |
| Developer B | 300 | 30% | $3,000 |
| Developer C | 200 | 20% | $2,000 |

---

## 7. Payment Mechanics

### 7.1 Payment Channels

| Contributor Location | Payment Method | Currency |
|----------------------|----------------|----------|
| International | Open Collective | USD |
| China | Chinese company bank transfer | RMB |
| Either | GitHub Sponsors (for smaller amounts) | USD |

### 7.2 Payment Schedule

- Quarterly, within **30 days** of quarter end
- Q1 payout: by April 30
- Q2 payout: by July 31
- Q3 payout: by October 31
- Q4 payout: by January 31

### 7.3 Minimum Payout

- Minimum payout threshold: $50 USD (or RMB equivalent)
- Amounts below threshold roll over to next quarter
- Contributors may opt to donate their share to the Community Fund

---

## 8. First 10 Core Maintainers Bonus

The first 10 individuals to achieve Core Maintainer status receive:

- **0.5% lifetime revenue share** (in addition to Pool distribution)
- This share is personal, non-transferable, and not deducted from the Pool
- It continues as long as the individual remains an Active Contributor or Core Maintainer
- Total maximum commitment: 5% (10 x 0.5%)

This incentive rewards early believers who help build the foundation.

---

## 9. Transparency

### 9.1 Quarterly Report

Every quarter, a transparency report is published in `.github/reports/YYYY-QN.md` containing:

- Gross revenue by product line
- COGS breakdown
- Net revenue
- Pool total
- Points ledger (every eligible contributor, their points, and payout)
- Community Fund expenditures
- Next quarter goals and budget

### 9.2 Public Audit

- Open Collective transactions are publicly visible
- Annual financial summary reviewed by TSC
- Any community member may request clarification via GitHub Discussions

---

## 10. Adjustments and Edge Cases

### 10.1 Pool Percentage Changes

- The 20% rate may only be changed via RFC with TSC supermajority (2/3) vote
- Minimum 14-day community comment period
- Cannot be reduced below 10% without unanimous TSC vote
- Changes take effect next quarter (not retroactive)

### 10.2 Zero Revenue Quarter

- If net revenue is zero or negative, no distribution occurs
- Points still accumulate and carry forward with 50% weight to next quarter

### 10.3 Disputes

- Contributors have 7 days after report publication to dispute points
- Disputes are reviewed by TSC
- TSC decision is final

### 10.4 Departure

- Contributors who leave retain payout rights for earned quarters
- Unclaimed payouts expire after 12 months
- Core Maintainer lifetime bonus pauses if contributor is inactive for 2+ consecutive quarters

---

## 11. Legal

- This agreement is supplementary to the Contributor License Agreement (CLA)
- Clawland reserves the right to modify this agreement via the RFC/TSC process described above
- All payouts are subject to applicable tax laws in the contributor's jurisdiction
- Contributors are responsible for their own tax reporting

---

## FAQ

**Q: Do I need to sign anything to participate?**
A: Yes, you need to sign the CLA (automated via GitHub bot on your first PR). The revenue share is automatic for eligible contributors.

**Q: What if I contribute to multiple repos?**
A: Points from all Clawland repos are combined into a single score.

**Q: Can I contribute hardware designs and earn points?**
A: Absolutely. Hardware kit designs earn 20 base points per kit — one of the highest-value contributions.

**Q: Is the 20% guaranteed?**
A: The 20% rate is committed in this public document and can only be changed through the governance process (TSC supermajority + community comment period). It cannot go below 10%.

**Q: When does this start?**
A: The first qualifying quarter begins when the first hardware kit or SaaS subscription generates revenue. Points accumulate from day one.

---

*Last updated: 2026-02-13*
*Version: 1.0*
