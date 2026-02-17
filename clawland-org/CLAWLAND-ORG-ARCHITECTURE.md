# Clawland GitHub 组织完整架构报告

> 数据截止: 2026-02-16 | 组织: [github.com/Clawland-AI](https://github.com/Clawland-AI)
> 参考模板: OpenClaw 组织架构 | 来源: clawdland_profit_sharing_plan + ARCHITECTURE-CLOUD-EDGE + SENSOR-SOLUTIONS + JOBS-REPLACEMENT-ANALYSIS

---

## 一、组织概览

| 项目 | 详情 |
|------|------|
| **组织名** | Clawland-AI |
| **定位** | Edge AI Agents for $10 Hardware. Build to Earn. |
| **口号** | "Let $10 hardware + AI Agents do the watching, so humans can do the thinking." |
| **官网** | [clawland.ai](https://clawland.ai) (规划) |
| **联系邮箱** | admin@clawland.ai |
| **Discord** | discord.gg/clawland-ai |
| **创建时间** | 2026-02 (规划中) |
| **公开仓库数** | 16 (含 1 个 upstream fork) |
| **许可证策略** | 双轨制 (Apache 2.0 核心 + BSL 1.1 商业保护) |
| **利润分享** | 净收入 20% → Contributor Revenue Pool |
| **治理模式** | TSC (Technical Steering Committee) + RFC |

---

## 二、全部仓库清单

### 第一层：核心 Agent 运行时 (Claw Family)

| # | 仓库 | 描述 | 语言 | 内存 | 硬件成本 | License | 优先级 |
|---|------|------|------|------|----------|---------|--------|
| 1 | [**picclaw**](https://github.com/Clawland-AI/picclaw) | 极致轻量边缘 AI Agent — 传感器感知 + 本地决策 + 即时执行 | Go | <10MB | $10 RISC-V 板 | Apache 2.0 | P0 |
| 2 | [**moltclaw**](https://github.com/Clawland-AI/moltclaw) | 全功能云端 AI Gateway — 多 Agent 路由 + 全局编排 | TypeScript | >1GB | 云服务器 / Mac Mini | Apache 2.0 | P0 |
| 3 | [**nanoclaw**](https://github.com/Clawland-AI/nanoclaw) | 中等量级 Agent — Python 生态 + SBC 部署 | Python | >100MB | $50 SBC (树莓派) | Apache 2.0 | P1 |
| 4 | [**microclaw**](https://github.com/Clawland-AI/microclaw) | MCU 级超微 Agent — 传感器嵌入式 | C / Rust | <1MB | $2-5 MCU | Apache 2.0 | P2 |

### 第二层：云边协同基础设施

| # | 仓库 | 描述 | 语言 | License | 优先级 |
|---|------|------|------|---------|--------|
| 5 | [**clawland-fleet**](https://github.com/Clawland-AI/clawland-fleet) | Cloud-Edge 编排系统 — Fleet Manager + Edge API + Reporter + 事件总线 | Go / TypeScript | BSL 1.1 → Apache 2.0 (4年) | P0 |
| 6 | [**clawland-deploy**](https://github.com/Clawland-AI/clawland-deploy) | 一键部署工具 — Ansible Playbook + Docker Compose + 预装镜像 | Shell / YAML | Apache 2.0 | P1 |

### 第三层：硬件 + 技能生态

| # | 仓库 | 描述 | 语言 | License | 优先级 |
|---|------|------|------|---------|--------|
| 7 | [**clawland-kits**](https://github.com/Clawland-AI/clawland-kits) | 硬件传感器套件设计 — BOM + 接线图 + 驱动脚本 + 场景方案 | Python / Markdown | CERN-OHL-S-2.0 | P0 |
| 8 | [**clawland-skills**](https://github.com/Clawland-AI/clawland-skills) | 社区技能市场 — 每个 Skill 一个目录, 按场景分类, 即插即用 | Markdown / YAML | MIT | P0 |
| 9 | [**clawland-drivers**](https://github.com/Clawland-AI/clawland-drivers) | 传感器驱动库 — GPIO/I2C/SPI/Serial 统一接口 + Python 驱动脚本集 | Python / C | Apache 2.0 | P1 |

### 第四层：网站 / 品牌 / 文档

| # | 仓库 | 描述 | 语言 | License | 优先级 |
|---|------|------|------|---------|--------|
| 10 | [**clawland.github.io**](https://github.com/Clawland-AI/clawland.github.io) | 组织官网 + 文档站 — 场景展示 + 快速入门 + API 文档 | Astro / MDX | CC BY-SA 4.0 | P1 |
| 11 | [**clawland-dashboard**](https://github.com/Clawland-AI/clawland-dashboard) | Fleet 可视化仪表盘 — 节点状态 + 告警面板 + 历史分析 | TypeScript (React) | BSL 1.1 → Apache 2.0 (4年) | P1 |

### 第五层：平台集成

| # | 仓库 | 描述 | 语言 | License | 优先级 |
|---|------|------|------|---------|--------|
| 12 | [**clawland-homebridge**](https://github.com/Clawland-AI/clawland-homebridge) | 智能家居桥接 — HomeKit / Home Assistant 集成 | TypeScript | MIT | P2 |
| 13 | [**clawland-grafana**](https://github.com/Clawland-AI/clawland-grafana) | Grafana 数据源插件 — 节点指标 + 传感器数据可视化 | Go | Apache 2.0 | P2 |

### 第六层：治理 / 信任 / 元仓库

| # | 仓库 | 描述 | 语言 | License | 优先级 |
|---|------|------|------|---------|--------|
| 14 | [**.github**](https://github.com/Clawland-AI/.github) | 组织级配置 — profile, GOVERNANCE, CLA, 模板, 季度报告 | Markdown | CC BY-SA 4.0 | P0 |
| 15 | [**clawland-cla**](https://github.com/Clawland-AI/clawland-cla) | CLA 签署记录 + CLA Assistant 配置 | Config | - | P0 |

### 第七层：生态协同 (Upstream Fork)

| # | 仓库 | 描述 | 语言 | License | 定位 |
|---|------|------|------|---------|------|
| 16 | [**openclaw**](https://github.com/Clawland-AI/openclaw) _(fork)_ | OpenClaw 上游跟踪 Fork — 生态桥梁 + 贡献通道 + Edge 集成补丁 | TypeScript | MIT (继承上游) | 生态桥梁 |

> **定位说明**: `openclaw` 是 [openclaw/openclaw](https://github.com/openclaw/openclaw) (196k+ Stars) 的 tracked fork。
> - **不替代 moltclaw** — moltclaw 是 Clawland 自主的 L3 云端 Gateway，包含 Fleet 集成等差异化功能
> - **用途**: (1) 跟踪上游更新 (2) 向上游贡献 Edge 相关功能 (3) 展示 Clawland 与 OpenClaw 的生态协同关系
> - **同步策略**: 每周自动 sync upstream，仅添加少量 edge-integration patch
> - **许可证**: 保持上游 MIT 不变，Clawland 新增代码单独标注 Apache 2.0

---

## 三、组织架构全景图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Clawland-AI 组织架构                                  │
│                 "Edge AI Agents for $10 Hardware"                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────── 核心 Agent 运行时 ───────────────────────┐        │
│  │                        (Claw Family)                            │        │
│  │                                                                 │        │
│  │  picclaw (Go)           moltclaw (TypeScript)                   │        │
│  │  ┌─────────────┐       ┌───────────────────┐                   │        │
│  │  │ Edge Agent   │       │ Cloud AI Gateway   │                  │        │
│  │  │ <10MB RAM    │       │ Multi-Agent Router  │                 │        │
│  │  │ $10 RISC-V   │◄────►│ Fleet Coordinator   │                 │        │
│  │  │              │ 上报  │ Decision Engine     │                  │        │
│  │  │ • 传感器采集 │ ═══► │ • 全局决策          │                  │        │
│  │  │ • 本地决策   │ ◄═══ │ • 跨节点关联        │                  │        │
│  │  │ • 即时执行   │ 指令  │ • 多渠道通知        │                  │        │
│  │  │ • 告警上报   │       │ • 知识库管理        │                  │        │
│  │  └─────────────┘       └───────────────────┘                   │        │
│  │                                                                 │        │
│  │  nanoclaw (Python)      microclaw (C/Rust)                      │        │
│  │  ┌─────────────┐       ┌───────────────────┐                   │        │
│  │  │ Mid-weight   │       │ MCU-level Micro    │                  │        │
│  │  │ >100MB RAM   │       │ <1MB RAM           │                  │        │
│  │  │ $50 SBC      │       │ $2-5 MCU           │                  │        │
│  │  │ Python 生态  │       │ 传感器级嵌入       │                  │        │
│  │  └─────────────┘       └───────────────────┘                   │        │
│  └─────────────────────────────────────────────────────────────────┘        │
│            │                        │                                       │
│            ▼                        ▼                                       │
│  ┌─────────────────── 云边协同基础设施 ───────────────────────────┐         │
│  │  clawland-fleet      Cloud-Edge 编排 (Fleet Manager +          │        │
│  │                      Edge API + Reporter + 事件总线)           │         │
│  │  clawland-deploy     一键部署 (Ansible + Docker + 预装镜像)    │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌─────────────────── 硬件 + 技能生态 ───────────────────────────┐         │
│  │  clawland-kits       传感器套件设计 (BOM + 接线 + 驱动)        │         │
│  │  clawland-skills     社区技能市场 (按场景分类)                  │         │
│  │  clawland-drivers    传感器驱动库 (GPIO/I2C/SPI/Serial)        │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌─────────────────── 网站 & 仪表盘 ─────────────────────────────┐         │
│  │  clawland.github.io  官网 + 文档 (Astro/MDX)                  │         │
│  │  clawland-dashboard  Fleet 可视化仪表盘 (React)                │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌─────────────────── 平台集成 ──────────────────────────────────┐         │
│  │  clawland-homebridge 智能家居桥接 (HomeKit/HA)                 │         │
│  │  clawland-grafana    Grafana 数据源插件                        │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌─────────────────── 治理 & 元仓库 ─────────────────────────────┐         │
│  │  .github             组织 profile + 模板 + GOVERNANCE          │         │
│  │  clawland-cla        CLA 签署记录                              │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  ┌─────────────────── 生态协同 (Upstream Fork) ─────────────────┐         │
│  │  openclaw (fork)     OpenClaw 上游跟踪 — 生态桥梁              │         │
│  │  ┌────────────────────────────────────────────────────┐       │         │
│  │  │ openclaw/openclaw (196k+ ⭐ MIT)                    │       │         │
│  │  │         │                                          │       │         │
│  │  │         │  fork + sync                             │       │         │
│  │  │         ▼                                          │       │         │
│  │  │ Clawland-AI/openclaw ──► moltclaw (自主 Gateway)   │       │         │
│  │  │         │                    │                     │       │         │
│  │  │         │  contribute back   │  Fleet/Edge 集成    │       │         │
│  │  │         ▲                    ▼                     │       │         │
│  │  │ openclaw/openclaw ◄── edge integration PRs        │       │         │
│  │  └────────────────────────────────────────────────────┘       │         │
│  └────────────────────────────────────────────────────────────────┘         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 四、三层分布式技术架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                     ☁️  L3 云端 (moltclaw 集群)                       │
│                                                                     │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────┐  ┌───────────┐ │
│  │ moltclaw    │  │ Fleet        │  │ Event Hub  │  │ Decision  │ │
│  │ AI Gateway  │  │ Manager      │  │ 事件总线   │  │ Engine    │ │
│  │             │  │              │  │            │  │           │ │
│  │ - 多Agent   │  │ - 节点注册   │  │ - 告警聚合 │  │ - 升级    │ │
│  │   路由      │  │ - 心跳监控   │  │ - 指令下发 │  │ - 聚合    │ │
│  │ - 复杂推理  │  │ - 能力声明   │  │ - 事件缓存 │  │ - 关联    │ │
│  │ - 知识库    │  │ - 分组管理   │  │            │  │ - 调度    │ │
│  └──────┬──────┘  └──────┬───────┘  └─────┬──────┘  └─────┬─────┘ │
│         │               │                │               │        │
│  ───────┴───────────────┴────────────────┴───────────────┴─────── │
│                    统一 API 层 (REST + WebSocket)                   │
│                    clawland-dashboard 仪表盘                        │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                    ══════════════╪══════════════
                    互联网 / 内网 / 4G / LoRa
                    ══════════════╪══════════════
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
       ┌──────┴──────┐    ┌──────┴──────┐    ┌──────┴──────┐
       │ L2 区域网关 A │    │ L2 区域网关 B │    │ L2 区域网关 C │
       │ nanoclaw     │    │ nanoclaw     │    │ nanoclaw     │
       │ (树莓派 4)   │    │ (树莓派 4)   │    │ (树莓派 4)   │
       │              │    │              │    │              │
       │ • 区域聚合   │    │ • 离线自治   │    │ • 跨节点关联 │
       │ • 断网缓存   │    │ • 本地路由   │    │ • 批量上报   │
       └──────┬───────┘    └──────┬───────┘    └──────┬───────┘
              │                   │                   │
        ┌─────┼─────┐      ┌─────┼─────┐      ┌─────┼─────┐
        │     │     │      │     │     │      │     │     │
       ┌┴┐  ┌┴┐  ┌┴┐    ┌┴┐  ┌┴┐  ┌┴┐    ┌┴┐  ┌┴┐  ┌┴┐
       │01│  │02│  │03│  │04│  │05│  │06│  │07│  │08│  │09│
       └──┘  └──┘  └──┘  └──┘  └──┘  └──┘  └──┘  └──┘  └──┘
       $10   $10   $10   $10   $10   $10   $10   $10   $10
       picclaw 末端节点 (LicheeRV-Nano / MaixCAM / Milk-V Duo)
       + microclaw 传感器芯片 ($2-5)

┌─────────────────────────────────────────────────────┐
│ 层级     │ Agent       │ 硬件             │ 职责      │
│──────────│─────────────│──────────────────│──────────│
│ L1 末端  │ picclaw     │ $10 LicheeRV    │ 采集+执行 │
│          │ + microclaw │ + $2-5 MCU      │ +即时告警 │
│ L2 区域  │ nanoclaw    │ $50 树莓派      │ 聚合+自治 │
│          │             │                  │ +关联分析 │
│ L3 云端  │ moltclaw    │ 云服务器/Mac    │ 全局决策  │
│          │             │ Mini            │ +知识库   │
└─────────────────────────────────────────────────────┘
```

---

## 五、与 OpenClaw 的架构对比

### 5.1 组织层级对比

| 架构层 | OpenClaw (18 仓库) | Clawland-AI (15 仓库) | 差异分析 |
|--------|-------------------|----------------------|----------|
| **核心运行时** | openclaw (CLI+Gateway), lobster (工作流), clawgo (Go 节点) | picclaw (边缘Go), moltclaw (云端TS), nanoclaw (中端Python), microclaw (MCU C/Rust) | Clawland 按硬件算力分级为 4 个 Agent，OpenClaw 单体 + 工作流引擎 |
| **技能生态** | clawhub (市场平台), skills (存档), barnacle (Bot 扩展) | clawland-skills (社区市场), clawland-drivers (驱动库) | Clawland 增加硬件驱动层，无独立市场 Web App (技能直接以仓库目录形式) |
| **基础设施** | openclaw-ansible, clawdinators (Nix), nix-openclaw, nix-steipete-tools, homebrew-tap | clawland-fleet (编排), clawland-deploy (部署) | OpenClaw 偏 DevOps 多方式部署; Clawland 核心是云边编排 Fleet |
| **硬件生态** | 无 | clawland-kits (传感器套件) | Clawland 独有的开源硬件层，CERN-OHL-S 许可 |
| **平台集成** | casa (HomeKit), voice-community | clawland-homebridge, clawland-grafana | 类似思路，Clawland 增加运维生态 (Grafana) |
| **网站/品牌** | openclaw.ai, butter.bot, flawd-bot | clawland.github.io, clawland-dashboard | Clawland 合并为文档站 + 仪表盘 |
| **治理** | trust, .github | .github, clawland-cla | 类似，Clawland 增加 CLA 管理 + 利润分享机制 |

### 5.2 核心理念对比

| 维度 | OpenClaw | Clawland-AI |
|------|----------|-------------|
| **核心产品** | 个人 AI 助手 (聊天/工具/多渠道) | 边缘 AI Agent (传感器+执行器+硬件) |
| **运行环境** | 桌面/服务器/手机 | $2-$600 全谱系硬件 (MCU → RISC-V → SBC → 云) |
| **数据理念** | "Own your data" — 本地优先 | "Own your edge" — 边缘自治 + 云端增强 |
| **商业模式** | 无明确商业模式 | 硬件套件销售 + SaaS (Fleet) + 咨询服务 |
| **开发者激励** | Star-driven / Contribution-driven | **20% 利润分享** (Contributor Revenue Pool) |
| **技术栈** | TypeScript 为主 | 多语言 (Go + TS + Python + C/Rust) 按场景分层 |
| **许可证** | MIT 为主 (9/18) | 双轨制 (Apache 2.0 核心 + BSL 1.1 商业) |

### 5.3 从 OpenClaw 借鉴的最佳实践

| 实践 | OpenClaw 做法 | Clawland 采纳 |
|------|-------------|---------------|
| 组织 profile | `.github` 仓库含 profile/README.md | ✅ 已创建 |
| 信任/安全 | 独立 `trust` 仓库 | ✅ 集成在 `.github/GOVERNANCE.md` |
| 工作流引擎 | 独立 `lobster` 仓库 | ❌ 集成在 `clawland-fleet` 编排系统中 |
| 技能市场 Web | 独立 `clawhub` (2000+ Stars) | 🔄 Phase 2 考虑，初期用仓库目录式 |
| Nix 打包 | 独立 `nix-openclaw` (360 Stars) | 🔄 Phase 3 考虑，初期用 Docker + Ansible |
| 多平台原生应用 | iOS/Android/macOS 内置 | ❌ 不适用 (硬件终端无 App Store 需求) |
| 品牌趣味仓库 | `flawd-bot` (邪恶双胞胎) | 🔄 可选，如 `crabclaw` (调皮版本) |

---

## 六、开源许可证策略 (双轨制)

```
核心原则: 核心代码足够开放以吸引开发者, 商业增值部分受保护以支撑利润分享

┌──────────────────────────────────────────────────────────────┐
│                    许可证分布                                  │
│                                                              │
│  Apache 2.0 (核心开放)           BSL 1.1 (商业保护)           │
│  ┌────────────────────┐         ┌────────────────────┐      │
│  │ picclaw            │         │ clawland-fleet     │      │
│  │ moltclaw           │         │ clawland-dashboard │      │
│  │ nanoclaw           │         │                    │      │
│  │ microclaw          │         │ ★ 4年后自动转为    │      │
│  │ clawland-deploy    │         │   Apache 2.0       │      │
│  │ clawland-drivers   │         └────────────────────┘      │
│  │ clawland-grafana   │                                      │
│  └────────────────────┘                                      │
│                                                              │
│  CERN-OHL-S-2.0 (开源硬件)      MIT (最宽松)                 │
│  ┌────────────────────┐         ┌────────────────────┐      │
│  │ clawland-kits      │         │ clawland-skills    │      │
│  │                    │         │ clawland-homebridge│      │
│  └────────────────────┘         └────────────────────┘      │
│                                                              │
│  CC BY-SA 4.0 (文档)                                         │
│  ┌────────────────────┐                                      │
│  │ clawland.github.io │                                      │
│  │ .github            │                                      │
│  └────────────────────┘                                      │
└──────────────────────────────────────────────────────────────┘
```

| 仓库 | 许可证 | 理由 |
|------|--------|------|
| picclaw | Apache 2.0 | 最大限度吸引开发者, 允许商用, 专利保护 |
| moltclaw | Apache 2.0 | 同上 |
| nanoclaw | Apache 2.0 | 同上 |
| microclaw | Apache 2.0 | 同上 |
| clawland-fleet | BSL 1.1 (4年转 Apache 2.0) | 保护云端编排/SaaS 收入 |
| clawland-dashboard | BSL 1.1 (4年转 Apache 2.0) | 保护商业仪表盘产品 |
| clawland-kits | CERN-OHL-S-2.0 | 开源硬件标准许可 |
| clawland-skills | MIT | 最宽松, 鼓励社区贡献 |
| clawland-drivers | Apache 2.0 | 驱动库需要开放以促进生态 |
| clawland-deploy | Apache 2.0 | 部署工具开放以降低门槛 |
| clawland.github.io | CC BY-SA 4.0 | 文档标准许可 |
| clawland-homebridge | MIT | 集成层最宽松 |
| clawland-grafana | Apache 2.0 | 插件开放促进生态 |

**BSL 1.1 附加使用授权**: 允许非竞争性生产环境使用。禁止将本软件作为托管服务直接与 Clawland 官方 SaaS 竞争。4年后自动转为 Apache 2.0。

---

## 七、技术栈分析

### 编程语言分布

| 语言 | 仓库数 | 代表仓库 | 用途 |
|------|--------|----------|------|
| **Go** | 2 | picclaw, clawland-grafana | 边缘 Agent + 高性能插件 |
| **TypeScript** | 3 | moltclaw, clawland-dashboard, clawland-homebridge | 云端网关 + Web 仪表盘 + 集成 |
| **Python** | 3 | nanoclaw, clawland-drivers, clawland-kits | 中端 Agent + 传感器驱动 + 硬件方案 |
| **C / Rust** | 1 | microclaw | MCU 级超微 Agent |
| **Shell / YAML** | 1 | clawland-deploy | 部署自动化 |
| **Astro / MDX** | 1 | clawland.github.io | 文档站 |
| **Markdown** | 3 | clawland-skills, .github, clawland-cla | 技能定义 + 治理 |

### 运行时 & 工具链

| 层级 | 运行时 | 构建工具 | 测试框架 | CI/CD |
|------|--------|----------|----------|-------|
| picclaw (Go) | Go 1.24+ | make / go build | go test | GitHub Actions |
| moltclaw (TS) | Node 22+ / Bun | pnpm / esbuild | Vitest | GitHub Actions |
| nanoclaw (Python) | Python 3.11+ | uv / pip | pytest | GitHub Actions |
| microclaw (C/Rust) | GCC / Cargo | CMake / Cargo | Unity / cargo test | GitHub Actions |
| clawland-fleet | Go + TS | make / pnpm | go test + Vitest | GitHub Actions |

### 硬件适配矩阵

| 硬件 | 价格 | CPU | RAM | 适配 Agent | 关键能力 |
|------|------|-----|-----|------------|----------|
| ESP32-C3 | $2 | RISC-V 160MHz | 400KB | microclaw | WiFi/BLE 传感器采集 |
| ESP32-S3 | $3-5 | Xtensa 240MHz | 512KB | microclaw | 带 AI 加速的 MCU |
| LicheeRV-Nano | $9.9 | RISC-V 1GHz | 256MB | picclaw | WiFi6/以太网, Linux |
| Milk-V Duo | $9 | RISC-V 1GHz | 64MB | picclaw | 极简 Linux, GPIO |
| NanoKVM | $30-50 | RISC-V | 256MB | picclaw | 远程 KVM, HDMI |
| MaixCAM | $50 | RISC-V + NPU | 256MB | picclaw | 摄像头, AI 推理 |
| 树莓派 4/5 | $35-80 | ARM A72/A76 | 4-8GB | nanoclaw | 通用 Linux, GPIO |
| Mac Mini / 云服务器 | $599+ | ARM/x86 | 8-16GB+ | moltclaw | 全能力云端 |

---

## 八、利润分享模型 (Contributor Revenue Pool)

### 8.1 收入分配架构

```
 ┌──────────────────────────────────────────────────────────┐
 │                     Revenue Sources                       │
 │  ┌──────────┐  ┌──────────────┐  ┌───────────────────┐  │
 │  │ Hardware  │  │    SaaS      │  │   Support &       │  │
 │  │ Kit Sales │  │ (Fleet/Dash) │  │   Consulting      │  │
 │  │           │  │ Subscriptions│  │   + License Fees  │  │
 │  └─────┬────┘  └──────┬───────┘  └─────────┬─────────┘  │
 │        │              │                     │            │
 │        └──────────────┼─────────────────────┘            │
 │                       ▼                                   │
 │               Gross Revenue                               │
 │                       │                                   │
 │              - COGS (~40-50%)                             │
 │                       │                                   │
 │                       ▼                                   │
 │               Net Revenue                                 │
 │                       │                                   │
 │        ┌──────────────┼──────────────┬────────────┐      │
 │        ▼              ▼              ▼            ▼      │
 │   ┌─────────┐  ┌───────────┐  ┌──────────┐  ┌───────┐  │
 │   │ Company │  │ ★ Pool    │  │ Community│  │ First │  │
 │   │ Reserve │  │ 20%       │  │ Fund 5%  │  │ 10    │  │
 │   │ 25-30%  │  │           │  │          │  │ 0.5%  │  │
 │   │         │  │ Quarterly │  │ Bounties │  │ each  │  │
 │   │ R&D     │  │ Payout    │  │ Events   │  │       │  │
 │   │ Legal   │  │           │  │ Docs     │  │ Life- │  │
 │   │ Market  │  │           │  │          │  │ time  │  │
 │   └─────────┘  └─────┬─────┘  └──────────┘  └───────┘  │
 │                       │                                   │
 │            ┌──────────┴──────────┐                        │
 │            ▼                     ▼                        │
 │     ┌────────────┐       ┌────────────┐                  │
 │     │ Open       │       │ Chinese    │                  │
 │     │ Collective │       │ Company    │                  │
 │     │ (USD)      │       │ (RMB)      │                  │
 │     │ 海外开发者  │       │ 国内开发者  │                  │
 │     └────────────┘       └────────────┘                  │
 └──────────────────────────────────────────────────────────┘
```

### 8.2 贡献积分体系

| 贡献类型 | 基础分 | 乘数规则 |
|----------|--------|----------|
| 代码合并 (PR merged) | 10 分/PR | 核心模块 x2, Bug fix x1.5, 文档 x0.8 |
| 代码审查 (Review) | 3 分/Review | 深度审查 (带建议的) x2 |
| Issue 报告 (有效) | 2 分/Issue | 含复现步骤 x1.5, 安全漏洞 x3 |
| Skill 贡献 | 15 分/Skill | 含测试用例 x1.5 |
| **硬件套件设计** | **20 分/Kit** | 含测试报告 x1.5 |
| 传感器驱动 | 12 分/Driver | 含兼容性矩阵 x1.5 |
| 文档/教程 | 5 分/篇 | 多语言 x1.5 |
| 社区支持 | 1 分/次 | 高质量回答 x3 |

**季度分配公式**: `个人分配 = Pool 总额 × (个人积分 / 所有参与者积分总和)`

### 8.3 贡献者层级

| 层级 | 条件 | 权益 |
|------|------|------|
| **Core Maintainer** | TSC 投票通过, 持续 6+ 个月核心贡献 | Revenue Pool + 技术决策权 + 路线图投票 |
| **Active Contributor** | 连续 3 个月有合并的 PR | Revenue Pool (按贡献分计算) |
| **Community Contributor** | 任何被合并的 PR | Bounty 奖励 + 社区基金 |
| **First 10 Bonus** | 前 10 位 Core Maintainer | 额外 0.5% 终身固定分成 (不占 Pool) |

---

## 九、治理结构

### 9.1 TSC (Technical Steering Committee)

```
┌─────────────────────────────────────────────────────┐
│                    TSC 组织架构                       │
│                                                     │
│  ┌─────────────────────────────────────────────┐    │
│  │            TSC (3-7 人)                      │    │
│  │                                             │    │
│  │  • 技术路线图决策                            │    │
│  │  • Core Maintainer 提名/投票                 │    │
│  │  • Revenue Pool 争议仲裁                     │    │
│  │  • RFC 审核与投票                            │    │
│  │  • 许可证变更 (2/3 超多数)                   │    │
│  │  • Pool 比例变更 (2/3 超多数, 不低于 10%)    │    │
│  └─────────────────────┬───────────────────────┘    │
│                        │                            │
│           ┌────────────┼────────────┐               │
│           ▼            ▼            ▼               │
│     ┌───────────┐ ┌──────────┐ ┌──────────┐       │
│     │ picclaw   │ │ moltclaw │ │ kits &   │       │
│     │ Lead      │ │ Lead     │ │ skills   │       │
│     │ Maintainer│ │ Maintainer│ │ Lead     │       │
│     └─────┬─────┘ └────┬─────┘ └────┬─────┘       │
│           │            │            │               │
│           ▼            ▼            ▼               │
│     Core Maintainers (各仓库)                       │
│           │                                         │
│           ▼                                         │
│     Active Contributors                             │
│           │                                         │
│           ▼                                         │
│     Community Contributors                          │
└─────────────────────────────────────────────────────┘
```

### 9.2 决策流程

| 决策类型 | 所需投票 | 投票期 |
|----------|----------|--------|
| 日常 (Bug 分诊, 发布时机) | 任意 1 名 Core Maintainer | 即时 |
| 标准 (功能审批, 依赖变更) | TSC 简单多数 (>50%) | 7 天 |
| 重大 (治理变更, 许可证, Revenue Share 调整) | TSC 超多数 (2/3) | 14 天 |
| 紧急 (安全漏洞, 服务宕机) | 任意 2 名 TSC 成员 | 即时 |

### 9.3 关键治理文件

| 文件 | 位置 | 内容 |
|------|------|------|
| `profile/README.md` | .github | 组织介绍 + 项目矩阵 + "Build to Earn" |
| `GOVERNANCE.md` | .github | TSC 章程 + 决策流程 + 晋升标准 |
| `CONTRIBUTOR-REVENUE-SHARE.md` | .github | 利润分享规则 + 积分体系 + FAQ |
| `CLA.md` | .github | 贡献者许可协议 |
| `CODE_OF_CONDUCT.md` | .github | 行为准则 (Contributor Covenant v2.1) |
| `CONTRIBUTING.md` | .github | 贡献指南 (适用全组织) |
| `reports/YYYY-QN.md` | .github | 季度透明度报告 |
| `rfcs/NNNN-title.md` | .github | RFC 提案 |

---

## 十、各仓库详细设计

### 10.1 picclaw — 边缘 AI Agent (旗舰)

```
picclaw/
├── cmd/
│   └── picclaw/
│       └── main.go              # 入口: CLI + Gateway 模式
├── pkg/
│   ├── agent/                   # Agent Loop (LLM 调用 + 工具执行)
│   │   ├── loop.go
│   │   ├── prompt.go
│   │   └── tools.go
│   ├── bus/                     # 消息总线 (Channel ↔ Agent 桥梁)
│   │   └── message_bus.go
│   ├── channels/                # 多渠道消息适配
│   │   ├── telegram/
│   │   ├── discord/
│   │   ├── feishu/
│   │   ├── dingtalk/
│   │   ├── whatsapp/
│   │   ├── qq/
│   │   └── maixcam/             # 硬件设备 TCP/JSON
│   ├── tools/                   # 内建工具
│   │   ├── exec.go              # Shell 执行 (传感器/执行器)
│   │   ├── web_fetch.go         # HTTP 请求
│   │   ├── read_file.go
│   │   ├── write_file.go
│   │   ├── cron.go              # 定时任务
│   │   └── spawn.go             # 子代理
│   ├── edge/                    # ★ 云边协同模块 (新增)
│   │   ├── server.go            # Edge HTTP Server (接收云端指令)
│   │   ├── reporter.go          # 事件上报器 (向云端汇报)
│   │   ├── config.go            # Edge 配置
│   │   └── tool_report.go       # report_event 工具
│   ├── memory/                  # 持久化记忆
│   │   ├── memory.go            # MEMORY.md 管理
│   │   └── daily_notes.go       # 每日笔记
│   ├── skills/                  # 技能加载
│   │   └── loader.go
│   └── config/
│       └── config.go
├── workspace/                   # 运行时数据目录
│   ├── SOUL.md                  # Agent 人格定义
│   ├── memory/
│   │   └── MEMORY.md
│   └── skills/
├── Makefile
├── go.mod
├── LICENSE                      # Apache 2.0
├── README.md
└── CONTRIBUTING.md
```

**关键指标**: <10MB RAM, 单二进制, RISC-V 原生支持, <1s 启动时间

### 10.2 moltclaw — 云端 AI Gateway

```
moltclaw/
├── src/
│   ├── gateway/                 # AI Gateway 核心
│   │   ├── router.ts            # 多 Agent 路由
│   │   ├── llm-providers.ts     # 8+ LLM 提供者适配
│   │   └── session.ts           # 会话管理
│   ├── fleet/                   # Fleet Manager
│   │   ├── manager.ts           # 节点管理 (注册/心跳/分组)
│   │   ├── event-hub.ts         # 事件总线 (聚合/去重/升级)
│   │   ├── command.ts           # 指令下发
│   │   └── decision-engine.ts   # 决策引擎 (关联分析/自动调度)
│   ├── api/                     # REST + WebSocket API
│   │   ├── routes/
│   │   │   ├── fleet.ts         # /fleet/* 边缘节点管理
│   │   │   ├── dashboard.ts     # /dashboard/* 仪表盘数据
│   │   │   └── skills.ts        # /skills/* 技能管理
│   │   └── middleware/
│   ├── channels/                # 通知渠道
│   │   ├── feishu.ts
│   │   ├── dingtalk.ts
│   │   ├── telegram.ts
│   │   └── webhook.ts
│   ├── knowledge/               # 知识库
│   │   └── store.ts
│   └── config/
├── package.json
├── tsconfig.json
├── LICENSE                      # Apache 2.0
└── README.md
```

### 10.3 clawland-fleet — 云边编排系统

```
clawland-fleet/
├── edge-api/                    # 边缘节点 API SDK
│   ├── server.go                # Edge HTTP Server 可嵌入库
│   ├── reporter.go              # Reporter 可嵌入库
│   └── discovery.go             # mDNS 节点发现
├── fleet-manager/               # 云端编排服务
│   ├── registry.ts              # 设备注册中心
│   ├── heartbeat.ts             # 心跳监控
│   ├── event-aggregator.ts      # 事件聚合
│   ├── command-dispatcher.ts    # 指令分发
│   └── scheduler.ts             # 任务调度
├── protocol/                    # 通信协议定义
│   ├── heartbeat.proto          # 心跳协议 (JSON Schema)
│   ├── event.proto              # 事件上报协议
│   ├── command.proto            # 指令下发协议
│   └── status.proto             # 状态查询协议
├── database/                    # 数据库 Schema
│   ├── edge_nodes.sql
│   ├── edge_events.sql
│   ├── edge_commands.sql
│   └── edge_groups.sql
├── LICENSE                      # BSL 1.1
└── README.md
```

### 10.4 clawland-kits — 硬件传感器套件

```
clawland-kits/
├── datacenter-sentinel/         # 方案一: 机房/数据中心无人值守
│   ├── BOM.md                   # 硬件清单 ($88 总成本)
│   ├── WIRING.md                # 接线图
│   ├── drivers/                 # 驱动脚本
│   │   ├── read_dht22.py        # 温湿度传感器
│   │   ├── read_ups.py          # UPS 状态
│   │   └── relay_control.py     # 继电器控制
│   ├── skills/                  # 预配置 Skill
│   │   └── datacenter-monitor.skill.md
│   ├── COST-ANALYSIS.md         # 替代 $48,000/yr 值班员
│   └── README.md
├── pond-guardian/               # 方案二: 水产养殖
│   ├── BOM.md                   # $133/塘
│   ├── WIRING.md
│   ├── drivers/
│   │   ├── read_do_sensor.py    # 溶氧传感器
│   │   └── aerator_control.py   # 增氧机控制
│   ├── skills/
│   └── README.md
├── greenhouse-manager/          # 方案三: 温室大棚
│   ├── BOM.md                   # $95/棚
│   ├── ...
├── safety-guardian/             # 方案四: 独居老人安全
│   ├── BOM.md                   # $30
│   ├── ...
├── equipment-doctor/            # 方案五: 设备预测性维护
│   ├── BOM.md                   # $32/台
│   ├── ...
├── SHARED-COMPONENTS.md         # 通用组件清单
├── LICENSE                      # CERN-OHL-S-2.0
└── README.md
```

### 10.5 clawland-skills — 社区技能市场

```
clawland-skills/
├── industrial/                  # 工业场景
│   ├── datacenter-monitoring/
│   │   └── SKILL.md
│   ├── predictive-maintenance/
│   │   └── SKILL.md
│   └── production-line-qc/
│       └── SKILL.md
├── agriculture/                 # 农业场景
│   ├── aquaculture-guardian/
│   │   └── SKILL.md
│   ├── greenhouse-automation/
│   │   └── SKILL.md
│   └── soil-monitoring/
│       └── SKILL.md
├── smart-home/                  # 智能家居
│   ├── elderly-care/
│   │   └── SKILL.md
│   ├── security-camera/
│   │   └── SKILL.md
│   └── energy-monitor/
│       └── SKILL.md
├── devops/                      # DevOps
│   ├── server-health-check/
│   │   └── SKILL.md
│   ├── ssl-certificate-watcher/
│   │   └── SKILL.md
│   └── container-monitor/
│       └── SKILL.md
├── SKILL-TEMPLATE.md            # Skill 开发模板
├── SKILL-REVIEW-CHECKLIST.md    # Skill 审核清单
├── LICENSE                      # MIT
└── README.md
```

---

## 十一、核心场景 × 仓库映射

| 场景 | 替代岗位 | 硬件成本 | 年节省 | 核心仓库 | 涉及仓库 |
|------|----------|----------|--------|----------|----------|
| 机房无人值守 | 值班员 $48k/yr | $88 | 99% | picclaw + datacenter-sentinel kit | picclaw, kits, skills, fleet, drivers |
| 水产养殖监控 | 巡塘员 $30k/yr | $133/塘 | 99% | picclaw + pond-guardian kit | picclaw, kits, skills, fleet, drivers |
| 温室大棚管理 | 种植员 $18k/yr | $95/棚 | 98% | picclaw + greenhouse-manager kit | picclaw, kits, skills, fleet, drivers |
| 独居老人看护 | 护理订阅 $900/yr | $30 | 90% | picclaw + safety-guardian kit | picclaw, kits, skills, homebridge |
| 设备预测维护 | 维保服务 $3.2k/yr | $32/台 | 95% | picclaw + equipment-doctor kit | picclaw, kits, skills, fleet, grafana |
| 全局编排管理 | 运维团队 | 云端 | - | moltclaw + fleet | moltclaw, fleet, dashboard |

---

## 十二、数据流全景

```
┌─────────────────────────────────────────────────────────────────────┐
│                          数据流全景图                                 │
│                                                                     │
│  物理世界                                                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐              │
│  │温度 42°C│  │溶氧 3mg│  │震动异常 │  │画面移动 │               │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘               │
│       │            │            │            │                      │
│       ▼            ▼            ▼            ▼                      │
│  ┌──────────────────────────────────────────────┐                   │
│  │ clawland-drivers (Python)                     │                  │
│  │ GPIO / I2C / SPI / Serial / USB / TCP         │                  │
│  │ read_dht22.py  read_do.py  read_vibration.py │                  │
│  └─────────────────────┬────────────────────────┘                   │
│                        │                                            │
│                        ▼                                            │
│  ┌──────────────────────────────────────────────┐                   │
│  │ picclaw (Go) — 边缘 Agent                     │                  │
│  │                                               │                  │
│  │  exec (采集) → LLM (分析) → exec (执行)       │                  │
│  │       │              │              │         │                  │
│  │       │         memory (学习)       │         │                  │
│  │       │              │              │         │                  │
│  │       │     clawland-skills         │         │                  │
│  │       │     (领域知识 SOP)          │         │                  │
│  │       │                             │         │                  │
│  │  本地判断: 42°C > 40°C 阈值                    │                  │
│  │  本地执行: 开启备用空调                         │                  │
│  │  本地记录: memory/202602/20260212.md           │                  │
│  │                                               │                  │
│  │  ┌─ 需要升级? ─────────────────────┐          │                  │
│  │  │                                 │          │                  │
│  │  ▼                                 ▼          │                  │
│  │  report_event (上报)     message (本地告警)    │                  │
│  └──────────┬──────────────────────┬─────────────┘                  │
│             │                      │                                │
│             │                ┌─────┴─────────────┐                  │
│             │                │ 飞书/钉钉/Telegram │                  │
│             │                │ (即时通知值班人)    │                  │
│             │                └───────────────────┘                  │
│             ▼                                                       │
│  ┌──────────────────────────────────────────────┐                   │
│  │ clawland-fleet (Cloud-Edge Protocol)          │                  │
│  │                                               │                  │
│  │  POST /fleet/events                           │                  │
│  │  POST /fleet/heartbeat                        │                  │
│  │  POST {node}/api/command                      │                  │
│  └─────────────────────┬────────────────────────┘                   │
│                        │                                            │
│                        ▼                                            │
│  ┌──────────────────────────────────────────────┐                   │
│  │ moltclaw (TypeScript) — 云端 AI Gateway       │                  │
│  │                                               │                  │
│  │  Fleet Manager: 聚合 rack-01 + rack-03 告警   │                  │
│  │  Decision Engine: 判断为区域空调故障            │                  │
│  │  Command: 向 rack-03 下发 enable_backup_cooling│                 │
│  │  Notify: 飞书群 @运维组 P1 告警 + 钉钉工单     │                 │
│  └─────────────────────┬────────────────────────┘                   │
│                        │                                            │
│                        ▼                                            │
│  ┌──────────────────────────────────────────────┐                   │
│  │ clawland-dashboard (React)                    │                  │
│  │                                               │                  │
│  │  实时节点地图 + 告警时间线 + 传感器图表         │                  │
│  │  + clawland-grafana 深度指标面板               │                  │
│  └──────────────────────────────────────────────┘                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 十三、实施路线图

### Phase 1: 基础建立 (第 1-2 周) — P0 仓库

```
Week 1:
├── 创建 GitHub 组织 Clawland-AI
├── 创建 .github 仓库 (全套治理文档已就绪)
│   ├── profile/README.md          ✅ 已编写
│   ├── GOVERNANCE.md              ✅ 已编写
│   ├── CONTRIBUTOR-REVENUE-SHARE.md ✅ 已编写
│   ├── CLA.md                     ✅ 已编写
│   ├── CODE_OF_CONDUCT.md         ✅ 已编写
│   ├── CONTRIBUTING.md            ✅ 已编写
│   ├── ISSUE_TEMPLATE/            ✅ 已编写
│   └── PULL_REQUEST_TEMPLATE.md   ✅ 已编写
├── 迁移 picclaw 仓库到组织
└── 配置 CLA Assistant bot

Week 2:
├── 创建 clawland-kits 仓库 (放入 datacenter-sentinel 方案)
├── 创建 clawland-skills 仓库 (放入首批 Skills)
├── 创建 clawland-fleet 仓库 (Edge API + Reporter 协议定义)
└── 所有仓库添加正确 LICENSE 文件
```

### Phase 2: 生态扩展 (第 3-6 周) — P1 仓库

```
Week 3-4:
├── 创建 moltclaw 仓库 (基础 AI Gateway 骨架)
├── 创建 clawland-drivers 仓库 (首批 5 个传感器驱动)
├── 创建 clawland-deploy 仓库 (Docker Compose + Ansible)
├── picclaw 实现 Edge API Server + Reporter
└── 注册 Open Collective

Week 5-6:
├── 创建 clawland.github.io 文档站 (Astro)
├── 创建 clawland-dashboard 仓库 (React 基础框架)
├── 发布 "Build to Earn" 公告 (HN/Reddit/V2EX/掘金)
├── 创建 Bounty Board (GitHub Projects)
└── 设置 GitHub Actions CI/CD
```

### Phase 3: 产品完善 (第 7-12 周) — P2 仓库

```
Week 7-9:
├── 创建 nanoclaw 仓库 (Python Agent)
├── 创建 clawland-homebridge (Home Assistant 集成)
├── 创建 clawland-grafana (Grafana 插件)
├── clawland-fleet 实现完整 Fleet Manager
├── clawland-kits 增加 pond-guardian + greenhouse-manager 方案
└── clawland-dashboard 实现节点地图 + 告警面板

Week 10-12:
├── microclaw 仓库 (C/Rust MCU Agent 原型)
├── L2 区域网关功能 (nanoclaw 增强)
├── 离线自治 + 断网缓存 + 恢复上传
├── 第一次季度报告 (reports/2026-Q1.md)
└── 发展首批 3-5 名 Active Contributors
```

### Phase 4: 持续迭代 (持续)

```
每季度:
├── 发布透明度报告
├── 维护 Bounty Board
├── 发展 Core Maintainer 团队
├── 新增场景套件 (clawland-kits)
├── 社区 Skills 评审与发布
└── 技术路线图 RFC 讨论
```

---

## 十四、成功指标

### 6 个月目标

| 指标 | 目标值 |
|------|--------|
| GitHub Stars (picclaw) | 500+ |
| 公开仓库数 | 10+ |
| Active Contributors | 10+ |
| Core Maintainers | 3+ |
| 硬件套件方案 | 5+ |
| 社区 Skills | 20+ |
| 传感器驱动 | 15+ |

### 12 个月目标

| 指标 | 目标值 |
|------|--------|
| GitHub Stars (全组织) | 2,000+ |
| 公开仓库数 | 15 |
| Active Contributors | 30+ |
| Core Maintainers | 10 (First 10 全部就位) |
| 首次 Revenue Pool 分配 | $5,000+ / 季度 |
| 硬件套件销售 | 100+ 套 |
| SaaS 订阅用户 | 50+ |
| 部署节点数 | 500+ |

---

## 十五、总结

Clawland-AI 是一个以 **边缘 AI Agent + 开源硬件 + 利润分享** 为核心的新型开源生态系统。参考 OpenClaw 的六层架构模型，结合自身"$10 硬件替代万元人力"的独特定位，组织架构设计为六个层次：

1. **核心 Agent 运行时** (picclaw + moltclaw + nanoclaw + microclaw) — 覆盖 $2 MCU 到云端服务器的全谱系
2. **云边协同基础设施** (clawland-fleet + clawland-deploy) — L1-L2-L3 三层分布式架构
3. **硬件 + 技能生态** (clawland-kits + clawland-skills + clawland-drivers) — OpenClaw 没有的硬件层
4. **网站 + 仪表盘** (clawland.github.io + clawland-dashboard) — 文档 + 可视化管理
5. **平台集成** (clawland-homebridge + clawland-grafana) — 连接现有生态
6. **治理** (.github + clawland-cla) — TSC + 20% Revenue Pool + CLA

### 与 OpenClaw 的关键差异

| 维度 | OpenClaw | Clawland-AI |
|------|----------|-------------|
| 核心价值 | 个人 AI 助手 | 边缘 AI 替代人力 |
| 硬件层 | 无 | **独有** (CERN-OHL-S 开源硬件) |
| 商业模式 | Star-driven 社区 | **利润分享** (20% Revenue Pool) |
| Agent 分级 | 单体 + 工作流 | **4 级按算力分层** ($2 → $10 → $50 → Cloud) |
| 许可证 | MIT 为主 | **双轨制** (Apache 2.0 + BSL 1.1) |
| 治理 | 信任仓库 | **TSC + RFC + 季度财报** |

**一句话定位**: Clawland 不是又一个 AI 框架，而是第一个**将完整 AI Agent 能力塞进 $10 硬件、并让开发者通过利润分享直接获利**的开源生态系统。

---

*Last updated: 2026-02-16*
*Version: 1.0*
