# Clawland GitHub Organization Setup Guide

This directory contains all the files needed to set up the `github.com/Clawland-AI` organization.

---

## Step 1: Create GitHub Organization

1. Go to https://github.com/organizations/plan
2. Create organization: **Clawland-AI**
3. Add a description: "Edge AI Agents for $10 Hardware. Build to Earn."

## Step 2: Create the `.github` Repository

This is the special org-level repo for profiles, templates, and governance.

```bash
# Create the repo on GitHub first, then:
git init Clawland-AI-.github
cd Clawland-AI-.github

# Copy all files from this directory's .github/ folder:
# .github/
#   profile/README.md        → Organization landing page
#   GOVERNANCE.md             → TSC charter and decision process
#   CONTRIBUTOR-REVENUE-SHARE.md → Revenue share rules
#   CLA.md                    → Contributor License Agreement
#   CODE_OF_CONDUCT.md        → Community standards
#   CONTRIBUTING.md            → How to contribute
#   ISSUE_TEMPLATE/            → Issue templates (shared across org)
#     bug_report.md
#     feature_request.md
#     sensor_kit_proposal.md
#   PULL_REQUEST_TEMPLATE.md   → PR template (shared across org)
#   reports/
#     TEMPLATE.md              → Quarterly report template

git add -A
git commit -m "feat: initialize Clawland organization governance and templates"
git remote add origin https://github.com/Clawland-AI/.github.git
git push -u origin main
```

## Step 3: Transfer/Fork picclaw Repository

Option A — Transfer existing repo:
1. Go to https://github.com/sipeed/picoclaw/settings (or your fork)
2. Scroll to "Danger Zone" → "Transfer ownership"
3. Transfer to `Clawland-AI` organization
4. Rename from `picoclaw` to `picclaw`

Option B — Fork and rename:
1. Fork https://github.com/sipeed/picoclaw to `Clawland-AI`
2. Rename to `picclaw` in repo settings
3. Update LICENSE to Apache 2.0 (see `licenses/APACHE-2.0-picclaw.txt`)

## Step 4: Create Additional Repositories

```bash
# On GitHub, create these repos under Clawland-AI org:
# - moltclaw          (TypeScript, Apache 2.0)
# - nanoclaw           (Python, Apache 2.0)
# - microclaw          (C/Rust, Apache 2.0)
# - clawland-fleet     (Go/TypeScript, BSL 1.1)
# - clawland-kits      (Hardware designs, CERN-OHL-S-2.0)
# - clawland-skills    (Markdown skills, MIT)
# - clawland.github.io (Docs site, CC BY-SA 4.0)
```

### License files

Use the license files in `licenses/` directory:
- `APACHE-2.0-picclaw.txt` → For picclaw, moltclaw, nanoclaw, microclaw
- `BSL-1.1-fleet.txt` → For clawland-fleet

## Step 5: Configure CLA Bot

1. Install [CLA Assistant](https://github.com/apps/cla-assistant) on the Clawland-AI organization
2. Configure it to use `Clawland-AI/.github/CLA.md` as the agreement
3. Enable for all repositories

## Step 6: Set Up Open Collective

1. Go to https://opencollective.com/create
2. Create collective: **Clawland-AI**
3. Choose fiscal host (Open Source Collective recommended)
4. Link GitHub organization
5. Set up payout methods (PayPal, Wise/TransferWise)

## Step 7: Create Bounty Board

1. Go to https://github.com/orgs/Clawland-AI/projects
2. Create project: "Bounty Board"
3. Add columns: Unclaimed, Claimed, In Progress, Review, Paid
4. Add initial bounty issues to each repo with price labels

## Step 8: Launch Announcement

Use `launch/BUILD-TO-EARN.md` as the template for:
- GitHub Discussions announcement (pin it)
- Hacker News post
- Reddit r/opensource, r/golang, r/homelab, r/selfhosted
- V2EX / 掘金 (Chinese developer communities)
- Discord server creation

---

## Directory Structure Reference

```
clawland-org/
├── .github/                          ← Files for .github repo
│   ├── profile/
│   │   └── README.md                 ← Org landing page
│   ├── GOVERNANCE.md                 ← TSC charter
│   ├── CONTRIBUTOR-REVENUE-SHARE.md  ← Revenue share rules
│   ├── CLA.md                        ← Contributor agreement
│   ├── CODE_OF_CONDUCT.md            ← Community standards
│   ├── CONTRIBUTING.md               ← How to contribute
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   ├── feature_request.md
│   │   └── sensor_kit_proposal.md
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── reports/
│       └── TEMPLATE.md               ← Quarterly report template
├── launch/
│   └── BUILD-TO-EARN.md              ← Launch announcement draft
├── licenses/
│   ├── APACHE-2.0-picclaw.txt        ← For core agent repos
│   └── BSL-1.1-fleet.txt             ← For clawland-fleet
└── SETUP-GUIDE.md                    ← This file
```
