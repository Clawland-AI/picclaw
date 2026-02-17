# Clawland Governance

This document describes the governance model for the Clawland organization and all projects under `github.com/Clawland-AI`.

---

## 1. Principles

- **Transparency**: All decisions, finances, and processes are public.
- **Meritocracy**: Influence is earned through sustained, high-quality contributions.
- **Sustainability**: The project must generate revenue to fund contributors fairly.
- **Openness**: Anyone can contribute; barriers to entry are kept minimal.

---

## 2. Roles

### 2.1 Community Contributor

Anyone who has had at least one PR merged into any Clawland repository.

**Rights**:
- Participate in GitHub Discussions
- Submit RFCs
- Receive bounty rewards for completed bounty tasks
- Listed in CONTRIBUTORS.md

### 2.2 Active Contributor

A contributor with merged PRs in at least 3 consecutive months.

**Rights**:
- All Community Contributor rights
- Participate in Contributor Revenue Pool quarterly distribution
- Nominate candidates for Core Maintainer
- Vote on non-binding community polls

### 2.3 Core Maintainer

An Active Contributor nominated and approved by the TSC after 6+ months of sustained core contributions.

**Rights**:
- All Active Contributor rights
- Merge permissions on assigned repositories
- Weighted share in Contributor Revenue Pool
- Vote on RFCs
- Eligible for TSC membership

**Responsibilities**:
- Review PRs within 72 hours
- Participate in at least 75% of TSC votes
- Mentor new contributors
- Uphold Code of Conduct

### 2.4 Technical Steering Committee (TSC) Member

Elected from Core Maintainers by existing TSC members.

**Rights**:
- All Core Maintainer rights
- Binding vote on RFCs, roadmap, and governance changes
- Revenue Pool dispute arbitration
- Core Maintainer nomination and approval

---

## 3. Technical Steering Committee (TSC)

### 3.1 Composition

- Minimum 3 members, maximum 7
- Initial TSC: Founder + first 2 Core Maintainers
- New members elected by existing TSC (2/3 majority required)

### 3.2 Terms

- 1-year terms, renewable
- Members may step down at any time with 30 days notice
- Removal requires 2/3 TSC vote for cause (inactivity, CoC violation)

### 3.3 Meetings

- Monthly sync (async-first via GitHub Discussions, optional video call)
- Meeting notes published in `.github/meetings/`
- All meetings are open for community observation

### 3.4 Decision Making

| Decision Type | Required Vote | Voting Period |
|---------------|---------------|---------------|
| Routine (bug triage, release timing) | Any 1 Core Maintainer | Immediate |
| Standard (feature approval, dependency changes) | Simple majority (>50%) of TSC | 7 days |
| Major (governance changes, license changes, revenue share adjustments) | Supermajority (2/3) of TSC | 14 days |
| Emergency (security vulnerabilities, service outages) | Any 2 TSC members | Immediate |

---

## 4. RFC Process

Major changes require a Request for Comments (RFC):

1. **Propose**: Create a markdown file in `.github/rfcs/NNNN-title.md` via PR
2. **Discuss**: Community discussion period (minimum 7 days)
3. **Vote**: TSC votes in the PR thread
4. **Merge or Close**: Based on vote outcome
5. **Implement**: Assigned to contributor(s)

### RFC Template

```markdown
# RFC-NNNN: Title

## Summary
One paragraph explanation.

## Motivation
Why are we doing this?

## Detailed Design
Technical details of the proposal.

## Alternatives Considered
What other approaches were evaluated?

## Impact
- Revenue impact:
- Developer experience impact:
- Breaking changes:
```

---

## 5. Repository Governance

Each repository has:

- **1-2 Lead Maintainers**: Primary decision makers for that repo
- **Repository-specific CONTRIBUTING.md**: Links to org-level guide with repo-specific additions
- **CODEOWNERS file**: Defines review requirements per directory

---

## 6. Conflict Resolution

1. **Step 1**: Direct discussion between parties in the relevant GitHub issue/PR
2. **Step 2**: Mediation by a Core Maintainer not involved in the dispute
3. **Step 3**: TSC vote (binding)
4. **Step 4**: If TSC is split, Founder has tie-breaking vote (limited to deadlocks only)

---

## 7. Amendments

This governance document can be amended through the RFC process with a 2/3 TSC supermajority vote and a 14-day community comment period.

---

*Last updated: 2026-02-13*
*Version: 1.0*
