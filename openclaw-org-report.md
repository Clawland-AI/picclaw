# OpenClaw GitHub 组织完整架构报告

> 数据截止: 2026-02-16 | 来源: [github.com/openclaw](https://github.com/orgs/openclaw/repositories)

---

## 一、组织概览

| 项目 | 详情 |
|------|------|
| **组织名** | openclaw |
| **定位** | Your personal, open source AI assistant |
| **官网** | [openclaw.ai](https://openclaw.ai) |
| **联系邮箱** | peter@openclaw.ai |
| **Twitter** | @openclaw |
| **创建时间** | 2026-01-04 |
| **公开仓库数** | 18 |
| **组织关注者** | 8,391 |
| **总 Stars** | ~201,200+ |
| **总 Forks** | ~35,400+ |
| **总 Open Issues** | ~6,570+ |

---

## 二、全部仓库清单

### 核心产品层

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 1 | [**openclaw**](https://github.com/openclaw/openclaw) | 主仓库 - 个人 AI 助手，跨平台 CLI + Gateway | TypeScript | **196,610** | 34,098 | 6,331 | MIT | Feb 15 |
| 2 | [**lobster**](https://github.com/openclaw/lobster) | Clawdbot 原生工作流 shell，可组合 pipeline 自动化引擎 | TypeScript | 464 | 113 | 8 | MIT | Feb 6 |
| 3 | [**clawgo**](https://github.com/openclaw/clawgo) | Go 语言实现的 Clawd 节点 | Go | 27 | 11 | 0 | MIT | Jan 5 |

### 技能/插件生态层

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 4 | [**clawhub**](https://github.com/openclaw/clawhub) | Skill 目录/市场平台 (clawhub.ai) | TypeScript | 2,031 | 433 | 178 | MIT | Feb 15 |
| 5 | [**skills**](https://github.com/openclaw/skills) | clawdhub.com 上所有 Skill 的存档仓库 | Python | 1,016 | 312 | 0 | MIT | Feb 15 |
| 6 | [**barnacle**](https://github.com/openclaw/barnacle) | "粘"在身边的实用 Bot 扩展 | TypeScript | 12 | 15 | 1 | - | Feb 15 |

### 基础设施 / 部署层

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 7 | [**openclaw-ansible**](https://github.com/openclaw/openclaw-ansible) | Ansible 自动化部署，集成 Tailscale VPN + UFW + Docker | Shell | 317 | 139 | 7 | MIT | Feb 13 |
| 8 | [**clawdinators**](https://github.com/openclaw/clawdinators) | 声明式基础设施 + NixOS 模块 (CLAWTINATOR hosts) | Nix | 106 | 26 | 2 | MIT | Feb 6 |
| 9 | [**nix-openclaw**](https://github.com/openclaw/nix-openclaw) | 将 OpenClaw 打包为 Nix package | Nix | 360 | 102 | 31 | Other | Feb 15 |
| 10 | [**nix-steipete-tools**](https://github.com/openclaw/nix-steipete-tools) | 核心维护者个人 Nix 工具集 | Nix | 18 | 20 | 3 | - | Feb 15 |

### 分发 / 安装渠道

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 11 | [**homebrew-tap**](https://github.com/openclaw/homebrew-tap) | macOS Homebrew 安装 tap | - | 27 | 8 | 0 | - | Jan 4 |

### 平台集成层

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 12 | [**casa**](https://github.com/openclaw/casa) | 智能家居集成 (HomeKit / Home base) | Swift | 24 | 9 | 0 | - | Feb 7 |
| 13 | [**voice-community**](https://github.com/openclaw/voice-community) | 语音社区功能 (新仓库) | - | 4 | 0 | 0 | - | Feb 7 |

### 网站 / 品牌 / 文档层

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 14 | [**openclaw.ai**](https://github.com/openclaw/openclaw.ai) | 组织官网 (Astro 框架) | Astro | 167 | 105 | 7 | - | Feb 13 |
| 15 | [**butter.bot**](https://github.com/openclaw/butter.bot) | butter.bot 域名网站 | HTML | 7 | 10 | 0 | MIT | Jan 2 |
| 16 | [**flawd-bot**](https://github.com/openclaw/flawd-bot) | "Clawd 的邪恶双胞胎" 趣味/品牌项目 | - | 26 | 10 | 0 | CC0 | Jan 24 |

### 治理 / 信任 / 元仓库

| # | 仓库 | 描述 | 语言 | Stars | Forks | Open Issues | License | 最近更新 |
|---|------|------|------|-------|-------|-------------|---------|----------|
| 17 | [**trust**](https://github.com/openclaw/trust) | OpenClaw 信任/安全资源 | - | 12 | 7 | 5 | - | Feb 8 |
| 18 | [**.github**](https://github.com/openclaw/.github) | GitHub 组织级配置文件 (profile, templates) | - | 3 | 5 | 1 | - | Jan 28 |

---

## 三、架构全景图

```
┌─────────────────────────────────────────────────────────────────┐
│                     OpenClaw 组织架构                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────────── 核心运行时 ────────────────────┐         │
│  │  openclaw (TypeScript)     主 CLI + Gateway        │         │
│  │  ├── src/cli               命令行交互              │         │
│  │  ├── src/commands          命令实现                │         │
│  │  ├── src/channels          多渠道消息 (Telegram,   │         │
│  │  │                         Discord, Slack, Signal, │         │
│  │  │                         iMessage, WhatsApp Web) │         │
│  │  ├── src/routing           消息路由                │         │
│  │  ├── src/media             媒体管道                │         │
│  │  ├── extensions/           插件系统 (msteams,      │         │
│  │  │                         matrix, zalo, voice-call│         │
│  │  └── apps/                 原生应用 (iOS, Android, │         │
│  │                            macOS)                  │         │
│  └────────────────────────────────────────────────────┘         │
│            │                         │                          │
│            ▼                         ▼                          │
│  ┌─────────────────┐    ┌──────────────────────┐               │
│  │  lobster         │    │  clawgo              │               │
│  │  工作流自动化     │    │  Go 节点实现          │               │
│  │  Pipeline Shell  │    │  (实验性)             │               │
│  └─────────────────┘    └──────────────────────┘               │
│                                                                 │
│  ┌──────────────────── 技能生态 ──────────────────────┐         │
│  │  clawhub          Skill 目录 + 市场 Web App       │         │
│  │  skills           Skill 版本存档                   │         │
│  │  barnacle         附属 Bot 扩展                    │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                 │
│  ┌──────────────────── 部署 & 分发 ──────────────────┐         │
│  │  openclaw-ansible  Ansible playbook (Docker+VPN)  │         │
│  │  clawdinators      NixOS 声明式基础设施            │         │
│  │  nix-openclaw      Nix 打包                       │         │
│  │  nix-steipete-tools 维护者工具链                    │         │
│  │  homebrew-tap      macOS Homebrew 安装             │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                 │
│  ┌──────────────────── 平台集成 ──────────────────────┐         │
│  │  casa              智能家居 (Swift/HomeKit)         │         │
│  │  voice-community   语音社区                        │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                 │
│  ┌──────────────────── 网站 & 品牌 ──────────────────┐         │
│  │  openclaw.ai       组织官网 (Astro)                │         │
│  │  butter.bot        域名网站 (HTML)                 │         │
│  │  flawd-bot         品牌/趣味项目                   │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                 │
│  ┌──────────────────── 治理 ─────────────────────────┐         │
│  │  trust             信任/安全资源                    │         │
│  │  .github           组织级 GitHub 配置              │         │
│  └────────────────────────────────────────────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 四、技术栈分析

### 编程语言分布

| 语言 | 仓库数 | 代表仓库 |
|------|--------|----------|
| **TypeScript** | 4 | openclaw, clawhub, lobster, barnacle |
| **Nix** | 3 | nix-openclaw, clawdinators, nix-steipete-tools |
| **Python** | 1 | skills |
| **Shell** | 1 | openclaw-ansible |
| **Astro** | 1 | openclaw.ai |
| **Swift** | 1 | casa |
| **Go** | 1 | clawgo |
| **HTML** | 1 | butter.bot |
| 无代码/配置 | 5 | trust, voice-community, .github, flawd-bot, homebrew-tap |

**核心语言: TypeScript** -- 主产品 + 技能平台 + 工作流引擎均基于 TypeScript (ESM)，运行时要求 Node 22+，开发优先使用 Bun。

### 许可证分布

| License | 数量 |
|---------|------|
| MIT | 9 |
| 无 License | 7 |
| CC0-1.0 | 1 |
| Other | 1 |

---

## 五、关键指标

### Stars 排名 Top 5

1. **openclaw** -- 196,610 (占全组织 97.7%)
2. **clawhub** -- 2,031
3. **skills** -- 1,016
4. **lobster** -- 464
5. **nix-openclaw** -- 360

### 活跃度 (最近 7 天内有推送)

- `openclaw` (Feb 15)
- `skills` (Feb 15)
- `clawhub` (Feb 15)
- `nix-openclaw` (Feb 15)
- `nix-steipete-tools` (Feb 15)
- `openclaw-ansible` (Feb 13)
- `openclaw.ai` (Feb 13)
- `trust` (Feb 8)

### Open Issues 分布

| 仓库 | Open Issues |
|------|-------------|
| openclaw | 6,331 (96.3%) |
| clawhub | 178 |
| nix-openclaw | 31 |
| lobster | 8 |
| openclaw-ansible | 7 |
| openclaw.ai | 7 |
| trust | 5 |
| 其余 11 个 | 0-3 |

---

## 六、产品架构解读

### 1. 核心产品: Clawdbot (`openclaw/openclaw`)

这是整个组织的旗舰项目 -- 一个**开源、本地优先的个人 AI 助手**，核心理念是 "own your data"。

**关键特性:**

- **CLI + Gateway 架构**: 通过 `clawdbot` CLI 运行本地 gateway，支持多种消息渠道
- **多渠道支持**: 核心支持 Telegram、Discord、Slack、Signal、iMessage、WhatsApp Web；通过插件支持 MS Teams、Matrix、Zalo、语音通话等
- **多平台**: macOS App (menubar)、iOS、Android 原生应用，以及 CLI (Any OS)
- **插件系统**: `extensions/` 目录下的 workspace packages
- **技能系统**: 可扩展的 Skill 框架，与 ClawdHub 生态互通
- **技术栈**: TypeScript (ESM), Node 22+, Bun, pnpm, Vitest, Oxlint

### 2. 工作流引擎: Lobster (`openclaw/lobster`)

Lobster 是一个 **Clawdbot 原生的工作流 shell** -- 将 skills/tools 编排为可组合的 pipeline，支持类型化、本地优先的自动化。本质上是一个 "macro engine"，让 Clawdbot 能一步调用复杂工作流。

### 3. 技能生态: ClawdHub (`openclaw/clawhub` + `openclaw/skills`)

- **clawhub**: 技能目录/市场的 Web 应用 (clawhub.ai)，社区可以发布、发现、安装 Skills
- **skills**: 所有已发布到 clawdhub.com 的 Skill 的版本化存档，使用 Python

### 4. 部署矩阵

OpenClaw 提供了极为丰富的部署选项:

| 方式 | 仓库 | 适用场景 |
|------|------|----------|
| npm 全局安装 | openclaw | 任何 Node.js 环境 |
| Homebrew | homebrew-tap | macOS |
| Nix Flake | nix-openclaw | NixOS / nix-darwin |
| Ansible Playbook | openclaw-ansible | VPS / 服务器自动化部署 |
| NixOS Module | clawdinators | 声明式服务器基础设施 |
| macOS App | openclaw (内置) | macOS menubar 应用 |
| Docker | openclaw-ansible | 容器化部署 |

### 5. 平台集成

- **casa**: Swift 编写的智能家居桥接，将 HomeKit 等家庭设备暴露给 Clawdbot
- **voice-community**: 语音交互社区功能 (最新仓库，2026-02-07 创建)

### 6. Go 实现: `clawgo`

Go 语言的 Clawd 节点实现，看起来是一个实验性/早期项目（仅 27KB，最近一次提交在 Jan 5），可能用于嵌入式或高性能场景。

---

## 七、组织健康度评估

### 优势

- **极高关注度**: 主仓库近 20 万 Stars，社区规模巨大
- **组织年轻但高产**: 成立仅 6 周 (2026-01-04 至今) 就建立了 18 个仓库的完整生态
- **多语言多平台**: TypeScript/Go/Swift/Nix/Python/Ansible 覆盖全栈
- **部署友好**: 提供 npm/Homebrew/Nix/Ansible/Docker 等多种安装方式
- **开源透明**: 大部分核心代码 MIT 许可
- **信任治理**: 专设 `trust` 仓库管理安全/信任资源

### 潜在关注点

- **Issue 积压**: 主仓库 6,331 个 open issues，需要持续分诊
- **许可证不统一**: 7 个仓库无 License，可能影响企业采用
- **Go 实现停滞**: `clawgo` 自 Jan 5 后无更新
- **品牌仓库目的模糊**: `flawd-bot`、`butter.bot` 的定位和维护优先级不明确

---

## 八、总结

OpenClaw 是一个以 **Clawdbot** 为核心的开源 AI 助手生态系统，围绕"个人数据自主、本地优先、跨平台"的理念构建。组织架构清晰地分为六个层次:

1. **核心运行时** (openclaw + lobster + clawgo)
2. **技能生态** (clawhub + skills + barnacle)
3. **基础设施/部署** (ansible + nix + clawdinators + homebrew)
4. **平台集成** (casa + voice-community)
5. **网站/品牌** (openclaw.ai + butter.bot + flawd-bot)
6. **治理** (trust + .github)

整个组织虽然成立仅 6 周，但已经建立了从核心产品到部署分发、从技能市场到平台集成的完整生态，并且以 ~197k Stars 的数据来看，这是一个在 AI 助手领域具有极高影响力的开源项目。
