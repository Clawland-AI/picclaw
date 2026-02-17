# Clawland-AI GitHub 组织更新日志

> 更新日期: 2026-02-16 | 操作者: @Tonyfudecai | 组织: [github.com/Clawland-AI](https://github.com/Clawland-AI)

---

## 一、更新概要

基于 [CLAWLAND-ORG-ARCHITECTURE.md](./CLAWLAND-ORG-ARCHITECTURE.md) 六层架构方案，对 Clawland-AI GitHub 组织进行全面扩展：新建 6 个仓库（从 **9 个扩展至 15 个**），为全部 15 个仓库添加 Topics 标签（共计约 100 个标签），更新 9 个现有仓库的 Description 和 Homepage，重写组织 Profile README 展示完整六层架构与 L1-L2-L3 分布式拓扑图，并补全组织级元信息（Description / Website / Twitter / Email）。

---

## 二、组织级变更

| 字段 | 更新前 | 更新后 |
|------|--------|--------|
| **Description** | *(未设置)* | Edge AI Agents for $10 Hardware. Build to Earn. 20% Revenue Pool for contributors. |
| **Website** | http://clawland.ai | https://clawland.ai |
| **Twitter** | *(未设置)* | @clawland_ai |
| **Email** | hello@clawland.ai | hello@clawland.ai *(不变)* |
| **Location** | United States of America | United States of America *(不变)* |
| **公开仓库数** | 9 | **15** |

---

## 三、新增 6 个仓库

### 3.1 clawland-deploy — 一键部署工具

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-deploy |
| **架构层级** | 第二层：云边协同基础设施 |
| **定位** | Ansible Playbooks、Docker Compose、预装镜像，覆盖 $10 开发板到云服务器的全谱系部署 |
| **License** | Apache 2.0 |
| **Topics** | `deployment` `ansible` `docker` `automation` `devops` `infrastructure` `edge-ai` |
| **推送内容** | README.md — 部署方法矩阵、Quick Start（Docker / Ansible / SD Card）、目录结构、硬件支持表 |

### 3.2 clawland-drivers — 统一传感器驱动库

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-drivers |
| **架构层级** | 第三层：硬件 + 技能生态 |
| **定位** | GPIO/I2C/SPI/Serial 统一接口的 Python 驱动脚本集，让 picclaw 通过 `exec` 一行命令读取任意传感器 |
| **License** | Apache 2.0 |
| **Topics** | `sensors` `drivers` `gpio` `i2c` `spi` `python` `hardware` `iot` `edge-ai` |
| **推送内容** | README.md — 已支持传感器目录（温湿度/水质/运动/振动/环境/执行器共 15+ 型号）、统一 JSON 输出规范、pip 安装示例、目录结构 |

### 3.3 clawland-dashboard — Fleet 可视化仪表盘

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-dashboard |
| **架构层级** | 第四层：网站 + 仪表盘 |
| **定位** | React + TypeScript 实时 Web 仪表盘：节点地图、告警时间线、传感器图表、指令中心 |
| **License** | BSL 1.1（4 年后转 Apache 2.0） |
| **Topics** | `dashboard` `visualization` `react` `fleet-management` `monitoring` `edge-ai` `real-time` |
| **推送内容** | README.md — 功能特性、技术栈（React 19 / Zustand / Recharts / Leaflet / Tailwind / Vite）、目录结构；LICENSE — BSL 1.1 全文 |

### 3.4 clawland-homebridge — 智能家居桥接

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-homebridge |
| **架构层级** | 第五层：平台集成 |
| **定位** | 将 Claw Agent 传感器数据暴露给 HomeKit / Home Assistant / Google Home / Alexa |
| **License** | MIT |
| **Topics** | `homekit` `home-assistant` `smart-home` `iot` `integration` `edge-ai` |
| **推送内容** | README.md — 平台支持矩阵、架构图（Smart Home ↔ Bridge ↔ Claw Agent）、设备映射表（温度→TemperatureSensor 等） |

### 3.5 clawland-grafana — Grafana 数据源插件

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-grafana |
| **架构层级** | 第五层：平台集成 |
| **定位** | Grafana Plugin SDK (Go + React) 构建的数据源插件，接入 clawland-fleet API 查询传感器指标和 Fleet 健康数据 |
| **License** | Apache 2.0 |
| **Topics** | `grafana` `monitoring` `visualization` `metrics` `datasource` `edge-ai` `iot` |
| **推送内容** | README.md — 功能特性、4 个预置仪表盘（Fleet Overview / Datacenter / Aquaculture / Sensor Deep Dive）、技术栈、目录结构 |

### 3.6 clawland-cla — CLA 签署管理

| 项目 | 详情 |
|------|------|
| **URL** | https://github.com/Clawland-AI/clawland-cla |
| **架构层级** | 第六层：治理 |
| **定位** | 存储 CLA Assistant Bot 的签署记录与配置，所有贡献者首次 PR 时自动触发签署 |
| **License** | *(无 — 治理类仓库)* |
| **Topics** | `cla` `contributor-license-agreement` `governance` `open-source` |
| **推送内容** | README.md — 签署流程说明、CLA 核心条款摘要（保留版权 + 授权许可 + 20% Revenue Pool 承诺）、相关链接 |

---

## 四、现有 9 个仓库更新

### 4.1 Description 变更

| 仓库 | 更新前 | 更新后 |
|------|--------|--------|
| **moltclaw** | Full-featured TypeScript AI Gateway with multi-Agent routing, session management, and cloud orchestration. | Full-featured TypeScript AI Gateway — multi-Agent routing, Fleet coordination, decision engine, and cloud orchestration. |
| **nanoclaw** | Mid-weight Python AI Agent with rich ecosystem support. Runs on $50 SBCs like Raspberry Pi. | Mid-weight Python AI Agent — rich ecosystem support, runs on $50 SBCs like Raspberry Pi. Regional gateway for L2 edge networks. |
| **microclaw** | Sensor-level micro AI Agent in C/Rust. <1MB RAM, runs on $2-5 MCU hardware. | Sensor-level micro AI Agent in C/Rust — <1MB RAM, runs on $2-5 MCU hardware like ESP32. The smallest member of the Claw family. |
| **clawland-fleet** | Cloud-Edge orchestration platform — Fleet Manager, Edge API Server, and Edge Reporter for managing distributed Claw agents. | Cloud-Edge orchestration platform — Fleet Manager, Edge API Server, Edge Reporter, Event Hub for managing distributed Claw agent fleets. |
| **clawland-kits** | Open hardware sensor kit designs — BOMs, wiring diagrams, drivers, and pre-built images for edge AI monitoring scenarios. | Open hardware sensor kit designs — BOMs, wiring diagrams, drivers, pre-built images. Replace $48k/yr jobs with $88 kits. |
| **clawland-skills** | Community skill marketplace — plug-and-play AI capabilities for all Claw family agents. | Community skill marketplace — plug-and-play AI capabilities for all Claw family agents. Industrial, agriculture, smart home, and DevOps skills. |
| **clawland-ai.github.io** | Documentation website for the Clawland ecosystem — Edge AI for Everyone | Documentation website for the Clawland ecosystem — quick start guides, API docs, scenario showcases, and deployment tutorials. |
| **.github** | Clawland-AI organization governance, templates, and profile | Clawland-AI organization governance — TSC charter, 20% Revenue Pool rules, CLA, Code of Conduct, contributing guide, and templates. |
| **picclaw** | *(内容不变)* | *(新增 homepage URL)* |

### 4.2 Homepage URL 添加

| 仓库 | Homepage |
|------|----------|
| picclaw | https://clawland.ai |
| moltclaw | https://clawland.ai |
| nanoclaw | https://clawland.ai |
| microclaw | https://clawland.ai |
| clawland-fleet | https://clawland.ai |
| clawland-kits | https://clawland.ai |
| clawland-skills | https://clawland.ai |
| .github | https://clawland.ai |
| clawland-ai.github.io | https://clawland-ai.github.io |

### 4.3 Topics 标签添加

| 仓库 | Topics |
|------|--------|
| **picclaw** | `ai-agent` `edge-computing` `iot` `go` `risc-v` `sensors` `automation` `edge-ai` `lightweight` `monitoring` |
| **moltclaw** | `ai-agent` `cloud-gateway` `typescript` `multi-agent` `orchestration` `ai-gateway` `llm` `edge-ai` |
| **nanoclaw** | `ai-agent` `python` `raspberry-pi` `sbc` `edge-computing` `iot` `edge-ai` `automation` |
| **microclaw** | `ai-agent` `mcu` `embedded` `c` `rust` `esp32` `sensors` `edge-ai` `microcontroller` `iot` |
| **clawland-fleet** | `fleet-management` `cloud-edge` `orchestration` `edge-computing` `iot-platform` `distributed-systems` `edge-ai` |
| **clawland-kits** | `hardware` `sensors` `open-hardware` `iot` `bom` `wiring-diagrams` `edge-ai` `diy` `cern-ohl` |
| **clawland-skills** | `skills` `marketplace` `plugins` `automation` `community` `edge-ai` `ai-agent` |
| **clawland-ai.github.io** | `documentation` `website` `astro` `edge-ai` `docs` |
| **.github** | `governance` `organization` `templates` `code-of-conduct` `revenue-sharing` `open-source` |

---

## 五、Profile README 更新

文件: `.github/profile/README.md`

Git 提交:
```
commit 41d5474
Message: docs: update org profile with complete 15-repo architecture and L1-L2-L3 diagram
```

### 变更内容

| 区块 | 变更 |
|------|------|
| **The Claw Family** | nanoclaw 描述增加 "Regional gateway" 定位 |
| **Infrastructure** | 新增 `clawland-deploy` 和 `clawland-dashboard` 两行 |
| **Hardware & Skills Ecosystem** | **新增区块** — `clawland-kits`, `clawland-skills`, `clawland-drivers` |
| **Integrations & Tools** | **新增区块** — `clawland-homebridge`, `clawland-grafana`, `clawland-ai.github.io` |
| **What Can You Build?** | 新增第 5 条场景: "$3,200/yr predictive maintenance — $32 equipment doctor" |
| **Architecture** | **新增区块** — L1-L2-L3 三层分布式 ASCII 拓扑图 |
| **License** | 从简单文字列表升级为六行分层许可证表格 |

---

## 六、完整 15 仓库架构总览

```
┌──────────────────────────────────────────────────────────┐
│              Clawland-AI 六层架构 (15 仓库)                │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  第一层: 核心 Agent 运行时                                 │
│  ├── picclaw          Go         Apache 2.0    ★ 旗舰    │
│  ├── moltclaw         TypeScript Apache 2.0              │
│  ├── nanoclaw         Python     Apache 2.0              │
│  └── microclaw        C/Rust     Apache 2.0              │
│                                                          │
│  第二层: 云边协同基础设施                                   │
│  ├── clawland-fleet   Go/TS      BSL 1.1                 │
│  ├── clawland-deploy  Shell      Apache 2.0    ★ 新建    │
│  └── clawland-dashboard React    BSL 1.1       ★ 新建    │
│                                                          │
│  第三层: 硬件 + 技能生态                                   │
│  ├── clawland-kits    Markdown   CERN-OHL-S              │
│  ├── clawland-skills  Markdown   MIT                     │
│  └── clawland-drivers Python     Apache 2.0    ★ 新建    │
│                                                          │
│  第四层: 网站 + 文档                                       │
│  └── clawland-ai.github.io  Astro  CC BY-SA 4.0         │
│                                                          │
│  第五层: 平台集成                                          │
│  ├── clawland-homebridge  TypeScript  MIT      ★ 新建    │
│  └── clawland-grafana     Go          Apache   ★ 新建    │
│                                                          │
│  第六层: 治理                                              │
│  ├── .github          Markdown   CC BY-SA 4.0            │
│  └── clawland-cla     Config     —            ★ 新建    │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

---

## 七、许可证分布

| License | 仓库数 | 仓库 |
|---------|--------|------|
| **Apache 2.0** | 7 | picclaw, moltclaw, nanoclaw, microclaw, clawland-deploy, clawland-drivers, clawland-grafana |
| **BSL 1.1** (4年后转 Apache 2.0) | 2 | clawland-fleet, clawland-dashboard |
| **MIT** | 2 | clawland-skills, clawland-homebridge |
| **CERN-OHL-S-2.0** | 1 | clawland-kits |
| **CC BY-SA 4.0** | 2 | clawland-ai.github.io, .github |
| **无 License** | 1 | clawland-cla |

---

## 八、Topics 热度统计

| Topic | 出现次数 | 覆盖仓库 |
|-------|----------|----------|
| `edge-ai` | 14 | 除 clawland-cla 外所有仓库 |
| `iot` | 7 | picclaw, nanoclaw, microclaw, clawland-kits, clawland-drivers, clawland-homebridge, clawland-grafana |
| `ai-agent` | 5 | picclaw, moltclaw, nanoclaw, microclaw, clawland-skills |
| `automation` | 4 | picclaw, nanoclaw, clawland-skills, clawland-deploy |
| `sensors` | 4 | picclaw, microclaw, clawland-kits, clawland-drivers |
| `monitoring` | 3 | picclaw, clawland-dashboard, clawland-grafana |
| `visualization` | 3 | clawland-dashboard, clawland-grafana, clawland-kits (wiring-diagrams) |
| `edge-computing` | 3 | picclaw, nanoclaw, clawland-fleet |
| `python` | 3 | nanoclaw, clawland-drivers, clawland-kits (driver scripts) |
| `orchestration` | 2 | moltclaw, clawland-fleet |
| `governance` | 2 | .github, clawland-cla |
| `open-source` | 2 | .github, clawland-cla |

---

## 九、操作日志

| # | 时间 (UTC) | 操作 | 状态 |
|---|------------|------|------|
| 1 | 21:26 | 安装 GitHub CLI — `winget install GitHub.cli` (v2.86.0) | Done |
| 2 | 21:27 | GitHub CLI Web 认证 — `gh auth login --web` | Done |
| 3 | 21:28 | 查询现有仓库 — `gh repo list Clawland-AI` (9 个) | Done |
| 4 | 21:29 | 创建 `clawland-deploy` — `gh repo create` Apache 2.0 | Done |
| 5 | 21:29 | 创建 `clawland-drivers` — `gh repo create` Apache 2.0 | Done |
| 6 | 21:29 | 创建 `clawland-dashboard` — `gh repo create` | Done |
| 7 | 21:29 | 创建 `clawland-homebridge` — `gh repo create` MIT | Done |
| 8 | 21:29 | 创建 `clawland-grafana` — `gh repo create` Apache 2.0 | Done |
| 9 | 21:29 | 创建 `clawland-cla` — `gh repo create` | Done |
| 10 | 21:30 | 为 picclaw 添加 10 个 Topics | Done |
| 11 | 21:30 | 为 moltclaw 添加 8 个 Topics | Done |
| 12 | 21:30 | 为 nanoclaw 添加 8 个 Topics | Done |
| 13 | 21:30 | 为 microclaw 添加 10 个 Topics | Done |
| 14 | 21:30 | 为 clawland-fleet 添加 7 个 Topics | Done |
| 15 | 21:31 | 为 clawland-deploy 添加 7 个 Topics | Done |
| 16 | 21:31 | 为 clawland-kits 添加 9 个 Topics | Done |
| 17 | 21:31 | 为 clawland-skills 添加 7 个 Topics | Done |
| 18 | 21:31 | 为 clawland-drivers 添加 9 个 Topics | Done |
| 19 | 21:31 | 为 clawland-ai.github.io 添加 5 个 Topics | Done |
| 20 | 21:31 | 为 clawland-dashboard 添加 7 个 Topics | Done |
| 21 | 21:31 | 为 clawland-homebridge 添加 6 个 Topics | Done |
| 22 | 21:32 | 为 clawland-grafana 添加 7 个 Topics | Done |
| 23 | 21:32 | 为 clawland-cla 添加 4 个 Topics | Done |
| 24 | 21:32 | 为 .github 添加 6 个 Topics | Done |
| 25 | 21:32 | 编写 6 个新仓库的 README.md 文件 | Done |
| 26 | 21:33 | Clone 6 个新仓库到本地临时目录 | Done |
| 27 | 21:34 | Push README → `clawland-deploy` (commit ea732af) | Done |
| 28 | 21:34 | Push README → `clawland-drivers` (commit 087ac52) | Done |
| 29 | 21:34 | Push README → `clawland-homebridge` (commit 48798f4) | Done |
| 30 | 21:34 | Push README → `clawland-grafana` (commit b1410f2) | Done |
| 31 | 21:34 | Push README → `clawland-dashboard` (commit 8204d54, 空仓库 `git checkout -b main`) | Done |
| 32 | 21:34 | Push README → `clawland-cla` (commit 346b714, 空仓库 `git checkout -b main`) | Done |
| 33 | 21:35 | Push BSL 1.1 LICENSE → `clawland-dashboard` (commit 91603ab) | Done |
| 34 | 21:35 | 更新 picclaw description + homepage | Done |
| 35 | 21:35 | 更新 moltclaw description + homepage | Done |
| 36 | 21:35 | 更新 nanoclaw description + homepage | Done |
| 37 | 21:35 | 更新 microclaw description + homepage | Done |
| 38 | 21:35 | 更新 clawland-fleet description + homepage | Done |
| 39 | 21:35 | 更新 clawland-kits description + homepage | Done |
| 40 | 21:35 | 更新 clawland-skills description + homepage | Done |
| 41 | 21:36 | 更新 clawland-ai.github.io description + homepage | Done |
| 42 | 21:36 | 更新 .github description + homepage | Done |
| 43 | 21:36 | Clone `.github` 仓库，重写 `profile/README.md` | Done |
| 44 | 21:36 | Push 更新后的 Profile README (commit 41d5474) | Done |
| 45 | 21:36 | 更新组织 Description — `gh api PATCH /orgs/Clawland-AI` | Done |
| 46 | 21:36 | 更新组织 Website / Twitter / Email — `gh api PATCH` | Done |
| 47 | 21:37 | 修复组织 Description 中 `$10` 的 PowerShell 变量转义问题 | Done |
| 48 | 21:37 | 修复 `.github` Description 中 `%%` 双百分号问题 | Done |
| 49 | 21:37 | 最终验证 — `gh repo list` 确认 15 仓库全部就位 | Done |
| 50 | 21:37 | 浏览器截图验证组织首页展示效果 | Done |

---

## 十、组织现状扫描 (2026-02-16 05:45 UTC+8)

> 以下为 `gh api` 全量扫描结果，用于判断哪些待办已被完成。

### 10.1 组织级状态

| 项目 | 状态 | 值 |
|------|------|-----|
| Description | OK | Edge AI Agents for $10 Hardware. Build to Earn. 20% Revenue Pool for contributors. |
| Website | OK | https://clawland.ai |
| Twitter | OK | @clawland_ai |
| Email | OK | hello@clawland.ai |
| Location | OK | United States of America |
| 公开仓库数 | OK | 15 |
| 组织 Projects | 存在 | 页面显示 "Projects 1" |
| 公开成员 | 缺失 | 0 人可见（Tonyfudecai 为成员但未设为 public） |

### 10.2 各仓库状态矩阵

| 仓库 | 文件数 | License | Topics | Homepage | Discussions | Issues | Labels | CODEOWNERS |
|------|--------|---------|--------|----------|-------------|--------|--------|------------|
| **picclaw** | 11 | Other (Apache 2.0) | 10 | clawland.ai | No | 0 | 0 | Yes |
| **moltclaw** | 2 | Apache 2.0 | 8 | clawland.ai | No | 0 | 0 | Yes |
| **nanoclaw** | 2 | Apache 2.0 | 8 | clawland.ai | No | 0 | 0 | — |
| **microclaw** | 2 | Apache 2.0 | 10 | clawland.ai | No | 0 | 0 | — |
| **clawland-fleet** | 2 | Other (BSL 1.1) | 7 | clawland.ai | No | 0 | 0 | Yes |
| **clawland-deploy** | 2 | Apache 2.0 | 7 | — | No | 0 | 0 | — |
| **clawland-dashboard** | 2 | Other (BSL 1.1) | 7 | — | No | 0 | 0 | — |
| **clawland-kits** | 2 | Other (CERN-OHL) | 9 | clawland.ai | No | 0 | 0 | — |
| **clawland-skills** | 2 | MIT | 7 | clawland.ai | No | 0 | 0 | — |
| **clawland-drivers** | 2 | Apache 2.0 | 9 | — | No | 0 | 0 | — |
| **clawland-homebridge** | 2 | MIT | 6 | — | No | 0 | 0 | — |
| **clawland-grafana** | 2 | Apache 2.0 | 7 | — | No | 0 | 0 | — |
| **clawland-cla** | 1 | — | 4 | — | No | 0 | 0 | — |
| **.github** | 9 | — | 6 | clawland.ai | **Yes** | 0 | 0 | Yes |
| **clawland-ai.github.io** | 3 | — | 5 | github.io | **Yes** | 0 | 0 | Yes |

### 10.3 关键发现

- **CODEOWNERS**: 5 个核心仓库已有（picclaw, moltclaw, clawland-fleet, .github, clawland-ai.github.io），其余 10 个仓库缺失
- **GitHub Discussions**: 仅 `.github` 和 `clawland-ai.github.io` 两个仓库启用，其余未启用（组织级讨论通过 `.github` 仓库承载，已满足需求）
- **Labels**: 全部 15 个仓库均无任何 Label（包括 GitHub 默认标签），无法使用 `good-first-issue` 分类
- **Issues**: 全部 15 个仓库均无 Open Issues
- **GitHub Actions**: 全部仓库 0 个 Workflow，未配置任何 CI/CD
- **CLA Assistant**: 未检测到已安装的 GitHub App
- **picclaw 内容**: 最完整（11 个文件/目录，含 cmd/, pkg/, skills/, Makefile, go.mod 等完整项目结构）
- **其他仓库**: 大多仅有 README.md + LICENSE 两个文件，骨架代码尚未创建
- **Open Collective**: 无法从 GitHub 侧验证，需单独确认

---

## 十一、后续待办（基于扫描结果更新）

> 最近更新: 2026-02-16 08:00 (UTC+8)

### Phase 1 补完（本周）— 状态更新

- [x] ~~为所有仓库初始化 Labels~~ — **已完成** (08:00): 全部 15 个非 fork 仓库已有 GitHub 默认 9 个 Label + 新增 `bounty` 标签，共 10 个/仓库
- [x] ~~创建首批 `good-first-issue` 类型的 Issue~~ — **已完成** (08:00): 5 个核心仓库各 5 个 Issue + 3 个生态仓库各 1 个 Issue，共 **28 个 Issue**
- [ ] **配置 CLA Assistant GitHub App** — 需手动在浏览器中安装 GitHub App，指向 `Clawland-AI/clawland-cla`（无法通过 CLI 完成）
- [x] ~~为全部 15 个非 fork 仓库设置 `CODEOWNERS`~~ — **已完成** (08:00): 全部 15 个仓库已通过 API 创建 CODEOWNERS，owner 设为 @Tonyfudecai
- [x] ~~启用组织级 GitHub Discussions~~ — **已完成** (之前): `.github` 仓库已启用 Discussions
- [x] ~~为缺失仓库设置 Homepage URL~~ — **已完成** (08:00): clawland-deploy, clawland-dashboard, clawland-drivers, clawland-homebridge, clawland-grafana, clawland-cla 已全部设置为 https://clawland.ai

### Phase 2（第 2-3 周）— 状态更新

- [x] ~~创建 GitHub Projects "Bounty Board"~~ — **已完成** (之前): 组织页面显示 "Projects 1"
- [x] ~~向 Bounty Board 添加 Issue~~ — **已完成** (08:15): 全部 28 个 Issue 已关联到 Bounty Board Project
- [x] ~~注册 Open Collective 账户~~ — **已完成** (用户手动完成): opencollective.com/Clawland-AI
- [x] ~~为 `clawland-ai.github.io` 初始化 Astro 框架~~ — **已完成** (08:20): 完整 Astro 5 项目骨架 (package.json + astro.config + landing page + favicon)
- [x] ~~picclaw 实现 Edge API Server + Edge Reporter~~ — **已完成** (08:15): `pkg/edge/` 下新增 config.go, server.go, reporter.go
- [x] ~~设置 GitHub Actions CI/CD~~ — **已完成** (08:20): picclaw 已有 Go CI workflow (lint + test + cross-build); openclaw fork 已有 upstream sync workflow
- [x] ~~为骨架仓库创建初始代码~~ — **已完成** (08:20): moltclaw (TypeScript), nanoclaw (Python + FastAPI), microclaw (ESP32/PlatformIO), clawland-fleet (Go), clawland-deploy (Ansible + Docker Compose) 均已初始化
- [x] ~~为 openclaw fork 设置自动 upstream sync~~ — **已完成** (08:10): 每周一 06:00 UTC 自动同步 + 支持手动触发

### Phase 3（第 4 周）— 状态更新

- [x] ~~将成员设为 public~~ — **已完成** (09:00): Tonyfudecai 已设为组织公开成员 (`gh api PUT /orgs/Clawland-AI/public_members/Tonyfudecai`)
- [x] ~~更新组织 Profile README~~ — **已完成** (09:00): Infrastructure 部分扩充 clawland-deploy、clawland-drivers、clawland-dashboard
- [x] ~~发布 "Build to Earn" 公告（GitHub Discussions）~~ — **已完成** (09:05): https://github.com/Clawland-AI/.github/discussions/2
- [x] ~~准备多平台发布文案~~ — **已完成** (09:10): HN / Reddit / V2EX / 掘金 四套文案已就绪
  - `clawland-org/launch/POST-HN.md` — Hacker News Show HN 格式
  - `clawland-org/launch/POST-REDDIT.md` — Reddit r/opensource 等子版块格式
  - `clawland-org/launch/POST-V2EX.md` — V2EX 中文社区格式
  - `clawland-org/launch/POST-JUEJIN.md` — 掘金长文格式（含架构图 + 场景表 + 技术栈详解）
- [ ] **在各平台发布公告** — 需手动在 HN / Reddit / V2EX / 掘金 发帖（文案已就绪）
- [ ] **创建 Discord 社区服务器** — 需在 discord.com 手动创建
- [ ] **置顶 Discussion #2** — 需在 GitHub Web UI 手动操作（API 不支持 pinDiscussion）

### Phase 4（持续）

- [ ] 每季度发布透明度报告 `reports/YYYY-QN.md`（报告模板已在 `.github/reports/` 就绪）
- [ ] 维护 Bounty Board 并追踪贡献积分
- [ ] 发展首批 10 名 Core Maintainer

---

## 十二、Phase 1 补完执行 (2026-02-16 08:00 UTC+8)

### 12.1 重新扫描发现

| 项目 | 之前记录 | 实际状态 |
|------|---------|---------|
| Labels | 全部 0 个 | 全部已有 GitHub 默认 9 个 (bug, documentation, duplicate, enhancement, good first issue, help wanted, invalid, question, wontfix) |
| CODEOWNERS | 5 个仓库有 | **全部 16 个仓库均无** CODEOWNERS（之前记录有误） |
| Homepage | 6 个缺失 | 确认 6 个缺失（与记录一致） |
| Issues | 全部 0 个 | 确认全部 0 个 |

### 12.2 执行操作日志

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 59 | 08:00 | 为全部 15 个非 fork 仓库添加 `bounty` label (金色 #E8A317) | Done |
| 60 | 08:00 | 为 6 个缺失仓库设置 Homepage URL → https://clawland.ai | Done |
| 61 | 08:01 | 为全部 15 个非 fork 仓库创建 CODEOWNERS 文件 (owner: @Tonyfudecai) | Done |
| 62 | 08:02 | picclaw: 创建 5 个 good-first-issue (架构图、健康检查、Dockerfile、测试、贡献指南) | Done |
| 63 | 08:02 | moltclaw: 创建 5 个 good-first-issue (项目初始化、HTTP 网关、架构文档、Docker、贡献指南) | Done |
| 64 | 08:03 | nanoclaw: 创建 5 个 good-first-issue (Python 项目、数据聚合、RPi 部署、Docker、离线决策) | Done |
| 65 | 08:03 | microclaw: 创建 5 个 good-first-issue (ESP32 初始化、DHT22 驱动、MQTT、接线图、深度睡眠) | Done |
| 66 | 08:04 | clawland-fleet: 创建 5 个 good-first-issue (API 规范、注册心跳、事件中心、部署图、指令下发) | Done |
| 67 | 08:04 | clawland-kits: 创建 1 个 bounty issue (数据中心监控套件 BOM) | Done |
| 68 | 08:04 | clawland-skills: 创建 1 个 bounty issue (温度告警技能模板) | Done |
| 69 | 08:04 | clawland-deploy: 创建 1 个 issue (Ansible 单节点部署 playbook) | Done |

### 12.3 当前各仓库状态

| 仓库 | Issues | Labels | CODEOWNERS | Homepage |
|------|--------|--------|------------|----------|
| picclaw | 5 | 10 | Yes | clawland.ai |
| moltclaw | 5 | 10 | Yes | clawland.ai |
| nanoclaw | 5 | 10 | Yes | clawland.ai |
| microclaw | 5 | 10 | Yes | clawland.ai |
| clawland-fleet | 5 | 10 | Yes | clawland.ai |
| clawland-kits | 1 | 10 | Yes | clawland.ai |
| clawland-skills | 1 | 10 | Yes | clawland.ai |
| clawland-deploy | 1 | 10 | Yes | clawland.ai |
| clawland-drivers | 0 | 10 | Yes | clawland.ai |
| clawland-dashboard | 0 | 10 | Yes | clawland.ai |
| clawland-homebridge | 0 | 10 | Yes | clawland.ai |
| clawland-grafana | 0 | 10 | Yes | clawland.ai |
| clawland-cla | 0 | 10 | Yes | clawland.ai |
| .github | 0 | 10 | Yes | clawland.ai |
| clawland-ai.github.io | 0 | 10 | Yes | clawland-ai.github.io |
| openclaw _(fork)_ | — | _(inherited)_ | _(upstream)_ | openclaw.ai |
| **合计** | **28** | **10/仓库** | **15/15** | **16/16** |

### 12.4 Issue 清单速览

#### picclaw (Go Edge Agent)
1. [#1](https://github.com/Clawland-AI/picclaw/issues/1) docs: add architecture diagram to README
2. [#2](https://github.com/Clawland-AI/picclaw/issues/2) feat: add health check endpoint /healthz
3. [#3](https://github.com/Clawland-AI/picclaw/issues/3) feat: add Dockerfile for containerized deployment
4. [#4](https://github.com/Clawland-AI/picclaw/issues/4) test: add unit tests for message bus
5. [#5](https://github.com/Clawland-AI/picclaw/issues/5) docs: create CONTRIBUTING.md with development setup guide

#### moltclaw (TypeScript Cloud Gateway)
1. [#1](https://github.com/Clawland-AI/moltclaw/issues/1) feat: initialize TypeScript project with tsconfig and package.json
2. [#2](https://github.com/Clawland-AI/moltclaw/issues/2) feat: implement basic HTTP gateway with Express/Hono
3. [#3](https://github.com/Clawland-AI/moltclaw/issues/3) docs: add architecture overview and API design document
4. [#4](https://github.com/Clawland-AI/moltclaw/issues/4) feat: add Dockerfile and docker-compose.yml
5. [#5](https://github.com/Clawland-AI/moltclaw/issues/5) docs: create CONTRIBUTING.md with dev setup guide

#### nanoclaw (Python SBC Agent)
1. [#1](https://github.com/Clawland-AI/nanoclaw/issues/1) feat: initialize Python project with pyproject.toml
2. [#2](https://github.com/Clawland-AI/nanoclaw/issues/2) feat: implement sensor data aggregation from PicoClaw nodes
3. [#3](https://github.com/Clawland-AI/nanoclaw/issues/3) docs: add README with Raspberry Pi deployment guide
4. [#4](https://github.com/Clawland-AI/nanoclaw/issues/4) feat: add Dockerfile for ARM64 deployment
5. [#5](https://github.com/Clawland-AI/nanoclaw/issues/5) feat: implement local decision engine for offline operation

#### microclaw (C/Rust MCU Agent)
1. [#1](https://github.com/Clawland-AI/microclaw/issues/1) feat: initialize ESP32 project with PlatformIO
2. [#2](https://github.com/Clawland-AI/microclaw/issues/2) feat: implement DHT22 temperature/humidity sensor driver
3. [#3](https://github.com/Clawland-AI/microclaw/issues/3) feat: implement MQTT client for edge reporting
4. [#4](https://github.com/Clawland-AI/microclaw/issues/4) docs: add hardware wiring guide for common sensor setups
5. [#5](https://github.com/Clawland-AI/microclaw/issues/5) feat: implement deep sleep with periodic wake for battery operation

#### clawland-fleet (Cloud-Edge Orchestration)
1. [#1](https://github.com/Clawland-AI/clawland-fleet/issues/1) feat: define Fleet API protocol specification (OpenAPI)
2. [#2](https://github.com/Clawland-AI/clawland-fleet/issues/2) feat: implement node registration and heartbeat service
3. [#3](https://github.com/Clawland-AI/clawland-fleet/issues/3) feat: implement event hub for sensor data collection
4. [#4](https://github.com/Clawland-AI/clawland-fleet/issues/4) docs: create deployment architecture diagram
5. [#5](https://github.com/Clawland-AI/clawland-fleet/issues/5) feat: add command dispatch to edge nodes

#### 生态仓库
- clawland-kits [#1](https://github.com/Clawland-AI/clawland-kits/issues/1) docs: create BOM and wiring guide for data center monitoring kit `bounty`
- clawland-skills [#1](https://github.com/Clawland-AI/clawland-skills/issues/1) feat: create temperature-alert skill template `bounty`
- clawland-deploy [#1](https://github.com/Clawland-AI/clawland-deploy/issues/1) feat: create Ansible playbook for single-node PicoClaw deployment

---

## 十三、Phase 2 执行 (2026-02-16 08:10-08:25 UTC+8)

### 13.1 执行操作日志

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 70 | 08:10 | `gh auth refresh -s project` — 获取 GitHub Projects API 权限 | Done |
| 71 | 08:12 | `gh project item-add` — 将 28 个 Issue 逐一关联到 Bounty Board Project | Done |
| 72 | 08:15 | picclaw: 创建 `pkg/edge/config.go` — Edge API 配置结构体 | Done |
| 73 | 08:15 | picclaw: 创建 `pkg/edge/server.go` — Edge API HTTP 服务器 (/healthz, /status, /command) | Done |
| 74 | 08:16 | picclaw: 创建 `pkg/edge/reporter.go` — Fleet 心跳上报 + 事件上报 + 节点注册 | Done |
| 75 | 08:17 | picclaw: 通过 git push 创建 `.github/workflows/ci.yml` — Go CI (lint + test + cross-build) | Done |
| 76 | 08:17 | openclaw fork: 创建 `.github/workflows/sync-upstream.yml` — 每周自动 sync | Done |
| 77 | 08:18 | moltclaw: 创建初始 TypeScript 骨架 (package.json, tsconfig.json, src/index.ts, .gitignore) | Done |
| 78 | 08:18 | nanoclaw: 创建初始 Python 骨架 (pyproject.toml, src/nanoclaw/, FastAPI server, Dockerfile) | Done |
| 79 | 08:19 | microclaw: 创建初始 ESP32 骨架 (platformio.ini, src/main.cpp, .gitignore) | Done |
| 80 | 08:19 | clawland-fleet: 创建初始 Go 骨架 (go.mod, cmd/fleet/main.go, pkg/fleet/registry.go) | Done |
| 81 | 08:20 | clawland-deploy: 创建初始 Ansible 骨架 (ansible.cfg, inventory, playbooks, templates, docker-compose.yml) | Done |
| 82 | 08:21 | clawland-ai.github.io: 创建 Astro 5 官网骨架 (package.json, astro.config, landing page, favicon) | Done |

### 13.2 骨架仓库初始化详情

| 仓库 | 技术栈 | 创建文件 | 文件总数 |
|------|--------|---------|---------|
| **picclaw** | Go | +3 edge files, +1 CI workflow | 13 → 17 |
| **moltclaw** | TypeScript (ESM, Node 22+) | package.json, tsconfig.json, src/index.ts, .gitignore | 3 → 7 |
| **nanoclaw** | Python 3.11+ (FastAPI + uvicorn) | pyproject.toml, src/nanoclaw/{__init__,cli,server}.py, Dockerfile, .gitignore | 3 → 9 |
| **microclaw** | C++ (PlatformIO, ESP32) | platformio.ini, src/main.cpp, .gitignore | 3 → 6 |
| **clawland-fleet** | Go 1.22 | go.mod, cmd/fleet/main.go, pkg/fleet/registry.go, .gitignore | 3 → 7 |
| **clawland-deploy** | Ansible + Docker | ansible.cfg, inventory/hosts.yml, playbooks/, templates/, docker-compose.yml, .gitignore | 3 → 9 |
| **clawland-ai.github.io** | Astro 5 + TypeScript | package.json, astro.config.mjs, tsconfig.json, src/pages/index.astro, public/favicon.svg, .gitignore | 4 → 10 |

### 13.3 picclaw Edge API Server 技术详情

新增 `pkg/edge/` 目录，包含 3 个文件:

| 文件 | 功能 |
|------|------|
| `config.go` | Edge 配置结构体 (端口、节点ID、云端连接设置、心跳间隔) |
| `server.go` | HTTP API 服务器 — `GET /healthz` (健康检查), `GET /api/v1/status` (节点状态), `POST /api/v1/command` (指令接收) |
| `reporter.go` | Fleet 心跳上报器 — `POST /fleet/register` (节点注册), `POST /fleet/heartbeat` (定期心跳), `ReportEvent()` (事件上报) |

### 13.4 GitHub Actions CI/CD

| 仓库 | Workflow | 触发条件 | Jobs |
|------|----------|---------|------|
| **picclaw** | `ci.yml` | push/PR to main | lint (golangci-lint), test (race + coverage), build (linux/darwin × amd64/arm64) |
| **openclaw** (fork) | `sync-upstream.yml` | 每周一 06:00 UTC + 手动触发 | fetch upstream → merge → push |

### 13.5 Bounty Board 状态

- Project 名称: "Paid tasks"
- 关联 Issue: **28 个** (5×picclaw + 5×moltclaw + 5×nanoclaw + 5×microclaw + 5×clawland-fleet + 1×clawland-kits + 1×clawland-skills + 1×clawland-deploy)
- 其中 2 个标记 `bounty` label: clawland-kits#1, clawland-skills#1

---

## 十四、OpenClaw 生态协同 — Fork 集成 (2026-02-16 07:30 UTC+8)

### 14.1 决策背景

为增强 Clawland 与 OpenClaw (196k+ Stars) 知名项目的生态协同融合，采用 **Hybrid 策略**：Fork `openclaw/openclaw` 作为生态桥梁，同时保留 `moltclaw` 作为 Clawland 自主的 L3 云端 Gateway。

**定位划分**:
- `Clawland-AI/openclaw` (fork) = **外交大使** — 上游跟踪、贡献通道、生态信号
- `Clawland-AI/moltclaw` (独立) = **自主产品** — Clawland 专属 Cloud Gateway，含 Fleet 集成等差异化功能

### 14.2 执行操作日志

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 51 | 07:26 | `gh repo fork openclaw/openclaw --org Clawland-AI` — Fork 创建 | Done |
| 52 | 07:26 | `gh repo edit` — 更新 fork description: "Clawland tracked fork of OpenClaw — edge integration patches, upstream sync, ecosystem bridge for moltclaw cloud gateway" | Done |
| 53 | 07:27 | `gh api PUT /repos/.../topics` — 添加 7 个 topics: edge-ai, clawland, openclaw, ai-agent, cloud-edge, iot, moltclaw-upstream | Done |
| 54 | 07:28 | Clone fork → 在 README.md 顶部添加 Clawland 生态桥梁 banner (含仓库关系表) | Done |
| 55 | 07:28 | Push README 更新到 Clawland-AI/openclaw (commit 5fbdc3f) | Done |
| 56 | 07:29 | 更新 `.github/profile/README.md` — 新增 "Ecosystem — OpenClaw Integration" 区块 (含 ASCII 架构图) | Done |
| 57 | 07:29 | Push Profile README 更新到 Clawland-AI/.github (commit 44328bf) | Done |
| 58 | 07:30 | 更新 `CLAWLAND-ORG-ARCHITECTURE.md` — 新增第七层"生态协同"、更新仓库数 15→16、全景图加入 fork | Done |

### 14.3 变更详情

#### 新增仓库

| 仓库 | 类型 | 描述 | License | Topics |
|------|------|------|---------|--------|
| [**openclaw**](https://github.com/Clawland-AI/openclaw) | Fork | OpenClaw 上游跟踪 Fork — 生态桥梁 + 贡献通道 + Edge 集成补丁 | MIT (继承上游) | edge-ai, clawland, openclaw, ai-agent, cloud-edge, iot, moltclaw-upstream |

#### 组织 Profile README 新增内容

- 新增 **"Ecosystem — OpenClaw Integration"** 区块，位于 "The Claw Family" 和 "Infrastructure" 之间
- 包含 ASCII 架构图，展示 OpenClaw (Cloud Brain) ↔ Clawland (Edge Hands & Feet) 关系
- 列出 `openclaw` fork 和 `moltclaw` 的定位对比表

#### 架构文档更新

- 组织概览: 仓库数从 "15 (规划)" 更新为 "16 (含 1 个 upstream fork)"
- 新增 **"第七层：生态协同 (Upstream Fork)"** 仓库列表
- 全景图: 在治理层下方新增生态协同区块，展示 fork → sync → contribute back 的完整流程

#### Fork README Banner

在 `openclaw/openclaw` 原始 README 顶部添加引用块 banner，说明：
- Fork 用途（跟踪上游、贡献回馈、桥接生态）
- 相关仓库链接表（moltclaw、picclaw、clawland-fleet、原始上游）

### 14.4 后续同步策略

| 项目 | 计划 |
|------|------|
| Upstream Sync | 每周一次自动同步 (可通过 GitHub Actions 或 GitHub 内置 "Sync fork" 实现) |
| Edge Patches | 仅添加少量 edge-integration 补丁，避免长期 diverge |
| 贡献回馈 | 开发 Edge API Server 协议、Fleet 心跳协议后，向上游 openclaw/openclaw 提交 PR |
| 许可证 | Fork 保持 MIT 不变，Clawland 新增代码单独标注 Apache 2.0 |

### 14.5 组织仓库总数

更新后 Clawland-AI 组织共 **16 个公开仓库**:

| 层级 | 仓库 | 数量 |
|------|------|------|
| 核心 Agent | picclaw, moltclaw, nanoclaw, microclaw | 4 |
| 云边协同 | clawland-fleet, clawland-deploy | 2 |
| 硬件 + 技能 | clawland-kits, clawland-skills, clawland-drivers | 3 |
| 网站 & 仪表盘 | clawland-ai.github.io, clawland-dashboard | 2 |
| 平台集成 | clawland-homebridge, clawland-grafana | 2 |
| 治理 | .github, clawland-cla | 2 |
| 生态协同 | openclaw _(fork)_ | 1 |
| **合计** | | **16** |

---

## 十四、Phase 3 执行 (2026-02-16 09:00 UTC+8)

### 14.1 Open Collective 注册

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 82 | 09:00 | 用户手动注册 Open Collective 账户 (opencollective.com/Clawland-AI) | Done |

### 14.2 成员公开化

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 83 | 09:00 | `gh api PUT /orgs/Clawland-AI/public_members/Tonyfudecai` — 设为公开成员 | Done |
| 84 | 09:00 | 验证: `gh api /orgs/Clawland-AI/public_members` → `Tonyfudecai` 已显示 | Done |

### 14.3 组织 Profile README 更新

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 85 | 09:02 | 更新 `.github/profile/README.md` — Infrastructure 部分扩充 3 个仓库 (clawland-deploy, clawland-drivers, clawland-dashboard) | Done |
| 86 | 09:02 | `git push origin main` → Clawland-AI/.github | Done |

### 14.4 Build to Earn 公告

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 87 | 09:05 | 通过 GraphQL `createDiscussion` 在 `.github` 仓库 Announcements 分类下创建公告 | Done |
| 88 | 09:05 | Discussion URL: https://github.com/Clawland-AI/.github/discussions/2 | Done |
| 89 | 09:05 | 公告内容: Build to Earn 完整介绍 (Claw Family 表格 + 收入分享模式 + 招募需求 + 链接) | Done |

### 14.5 多平台发布文案

| # | 时间 | 操作 | 状态 |
|---|------|------|------|
| 90 | 09:10 | 创建 `launch/POST-HN.md` — Hacker News "Show HN" 格式 | Done |
| 91 | 09:10 | 创建 `launch/POST-REDDIT.md` — Reddit 格式 (含推荐 subreddits: r/opensource, r/selfhosted, r/homelab, r/IOT, r/embedded, r/golang, r/Python, r/rust) | Done |
| 92 | 09:10 | 创建 `launch/POST-V2EX.md` — V2EX 中文社区格式 (RMB 价格换算) | Done |
| 93 | 09:10 | 创建 `launch/POST-JUEJIN.md` — 掘金长文格式 (含 ASCII 架构图 + 场景表 + 技术栈详解) | Done |

### 14.6 Phase 3 完成总结

**已自动化完成:**
- ✅ Open Collective 账户注册（用户手动完成）
- ✅ 成员公开化 (Tonyfudecai → public member)
- ✅ 组织 Profile README 扩充
- ✅ GitHub Discussions 公告帖创建
- ✅ 四平台发布文案准备 (HN / Reddit / V2EX / 掘金)

**需手动操作:**
- ⏳ 在 HN / Reddit / V2EX / 掘金 各平台发布帖子（文案已就绪于 `clawland-org/launch/` 目录）
- ⏳ 创建 Discord 社区服务器
- ⏳ 在 GitHub Web UI 置顶 Discussion #2
- ⏳ 安装 CLA Assistant GitHub App (Phase 1 遗留)

---

## 后续待办总览

### 手动操作（不可自动化）

| 序号 | 任务 | 平台/位置 | 文案/指南 |
|------|------|-----------|-----------|
| M1 | 安装 CLA Assistant App | https://github.com/apps/cla-assistant | 指向 `Clawland-AI/clawland-cla` |
| M2 | Hacker News 发帖 | https://news.ycombinator.com/submit | `launch/POST-HN.md` |
| M3 | Reddit 发帖 | r/opensource, r/selfhosted 等 | `launch/POST-REDDIT.md` |
| M4 | V2EX 发帖 | https://v2ex.com/new | `launch/POST-V2EX.md` |
| M5 | 掘金发文 | https://juejin.cn/editor/drafts/new | `launch/POST-JUEJIN.md` |
| M6 | 创建 Discord 服务器 | https://discord.com | 建议频道: #general, #dev, #hardware, #bounties |
| M7 | 置顶 Discussion #2 | GitHub Web UI → Discussions | 点击 Pin discussion |

### Phase 4（持续运营）

- [ ] 每季度发布透明度报告 `reports/YYYY-QN.md`（报告模板已在 `.github/reports/` 就绪）
- [ ] 维护 Bounty Board 并追踪贡献积分
- [ ] 发展首批 10 名 Core Maintainer
- [ ] 建立 Open Collective 季度分红流程
- [ ] 持续为各仓库添加 Issue 和功能开发

---

*初次生成: 2026-02-16 05:40 (UTC+8)*
*扫描更新: 2026-02-16 05:45 (UTC+8)*
*Fork 集成: 2026-02-16 07:30 (UTC+8)*
*Phase 1 补完: 2026-02-16 08:00 (UTC+8)*
*Phase 2 执行: 2026-02-16 08:25 (UTC+8)*
*Phase 3 执行: 2026-02-16 09:10 (UTC+8)*
