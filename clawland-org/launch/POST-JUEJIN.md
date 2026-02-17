# 掘金文章

**标签:** 开源, AI, IoT, Go, 边缘计算

**标题:** 开源项目 Clawland：$10 硬件 + AI Agent 替代数万元/年的人力监控岗位，贡献者分享 20% 收入

---

## 前言

一个数据中心夜班运维工程师年薪约 30 万。一个鱼塘巡塘工年薪约 20 万。一个温室管理员年薪约 12 万。

这些工作很重要，但本质上是重复性的巡检和告警。

现有的 IoT 方案要么太贵（动辄数万起步），要么只能做简单的 if-then 告警。

**Clawland** 想要改变这一切。

## 什么是 Clawland？

Clawland 是一个开源的边缘 AI Agent 生态系统，Agent 运行在从 ¥15 单片机到云服务器的全谱系硬件上：

```
    ☁️ Cloud Layer
    ┌──────────────────────┐
    │ moltclaw (TypeScript) │ ← 云端多 Agent 路由网关
    └──────────┬───────────┘
               │
    🌐 Regional Gateway
    ┌──────────┴───────────┐
    │ nanoclaw (Python)     │ ← 区域网关，运行在 ¥350 SBC
    └──────────┬───────────┘
               │
    🤖 Edge Agent
    ┌──────────┴───────────┐
    │ picclaw (Go, <10MB)   │ ← 边缘 Agent，¥70 RISC-V 开发板
    └──────────┬───────────┘
               │
    📡 Sensor Layer
    ┌──────────┴───────────┐
    │ microclaw (C/Rust)    │ ← 传感器级 Agent，¥15 MCU
    └──────────────────────┘
```

## 核心能力

- **感知 (Sense)**：温度、湿度、溶解氧、振动、烟雾、运动（通过 exec + Python/GPIO）
- **思考 (Think)**：LLM 驱动的多传感器分析 + 趋势预测（支持 8+ LLM 提供商）
- **执行 (Act)**：控制继电器、阀门、开关（通过 exec + GPIO）
- **上报 (Report)**：Telegram、Discord、飞书、钉钉、WhatsApp、QQ（7 个通道）
- **学习 (Learn)**：持久化记忆，越用越聪明
- **调度 (Schedule)**：基于 Cron 的自主监控循环

## 真实数字

| 套件 | 硬件成本 | 替代岗位 | 传统成本 | 节省 |
|------|---------|---------|---------|------|
| 数据中心哨兵 | ¥600 | 夜班运维 | ¥30万/年 | 99% |
| 鱼塘守护者 | ¥950/塘 | 巡塘工 | ¥20万/年 | 99% |
| 温室管家 | ¥680/棚 | 管理员 | ¥12万/年 | 98% |
| 安全卫士 | ¥215 | 独居老人看护 | ¥6500/年 | 90% |
| 设备医生 | ¥230/台 | 预测性维护服务 | ¥2.3万/年 | 95% |

## 基因进化协议（CGEP）

Clawland 的核心技术创新是 **Clawland Gene Evolution Protocol (CGEP)**——一套让 AI Agent 自主进化监控策略的协议。

### 工作原理

1. **感知信号**：Agent 从传感器数据中提取信号（传感器错误、阈值越界、交叉传感器异常等 8 种信号类型）
2. **匹配基因**：在本地基因池中查找最匹配的监控策略（基因），注入到 Agent 的系统提示中
3. **执行并学习**：Agent 执行策略后，经验被固化为"胶囊"（Capsule），基因信心值自动调整
4. **跨节点共享**：高信心基因通过 Fleet 管理系统在边缘节点间共享

这意味着：**部署越久，Agent 越聪明。一个节点学到的经验，所有节点都能受益。**

## Build to Earn 模式

Clawland 不走传统捐赠制开源路线。我们提供的是一笔交易：

> **你来建设，它产生收入，你拿 20%。**

### 具体机制

1. **净产品收入的 20%** 进入贡献者收入池（硬件套件销售 + SaaS 订阅 + 咨询收入）
2. **按季度分红**，通过 Open Collective 完全公开透明
3. **积分制度**：代码贡献、Code Review、技能包、硬件设计、文档、社区支持都计分
4. **前 10 名核心维护者** 额外获得 0.5% 终身收入份额

### 举例

假设 Clawland 在 Q3 产生 ¥35 万净收入：
- 贡献者池：¥7 万
- 贡献了 500/1000 总积分的开发者：**当季度获得 ¥3.5 万**

## 技术栈

| 组件 | 语言 | 许可证 |
|------|------|--------|
| picclaw（边缘 Agent） | Go | Apache 2.0 |
| moltclaw（云网关） | TypeScript | Apache 2.0 |
| nanoclaw（区域网关） | Python | Apache 2.0 |
| microclaw（MCU Agent） | C/Rust | Apache 2.0 |
| clawland-fleet（编排） | Go | BSL 1.1 → Apache 2.0 |
| clawland-kits（硬件设计） | — | CERN-OHL-S-2.0 |

## 我们需要你

### 现在就可以开始

1. 浏览各仓库的 `good first issue` 标签
2. 查看 [Bounty Board](https://github.com/orgs/Clawland-AI/projects/1) 领取付费任务
3. 阅读 [贡献指南](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTING.md)
4. 参与 [GitHub 讨论](https://github.com/Clawland-AI/.github/discussions/2)

### 重点招募

- **Go 开发者**：picclaw 边缘 Agent 核心（ARM 交叉编译、GPIO 集成、LLM 对接）
- **TypeScript 开发者**：moltclaw 云网关（多 Agent 路由、Fleet 管理 Dashboard）
- **Python 开发者**：nanoclaw + 传感器驱动库
- **嵌入式开发者**：microclaw MCU Agent（ESP32/STM32）
- **硬件极客**：传感器套件设计、测试、BOM 优化
- **技术写作**：文档、教程、部署指南

## 链接

- **GitHub**: [github.com/Clawland-AI](https://github.com/Clawland-AI)
- **Bounty Board**: [付费任务看板](https://github.com/orgs/Clawland-AI/projects/1)
- **收入分享条款**: [CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)
- **治理**: [GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md)

---

*让 ¥70 的硬件 + AI Agent 来做监控，让人类去做更有价值的事。*

**我们正在构建一个未来：每个鱼塘、每个机房、每个温室都有一个 AI 守护者。加入我们。**
