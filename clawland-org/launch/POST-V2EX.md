# V2EX 帖子

**节点:** /go/create/share 或 /go/create/programmer

**标题:** [开源] Clawland: $10 硬件上运行的边缘 AI Agent，贡献者分享 20% 收入

**正文:**

大家好，我们正在构建 **Clawland**，一个运行在超低成本硬件上的开源 AI Agent 生态系统。

## 项目简介

一系列不同规模的 AI Agent：

| Agent | 语言 | 内存 | 硬件成本 | 用途 |
|-------|------|------|---------|------|
| picclaw | Go | <10MB | ¥70 开发板 | 边缘感知 + 本地决策 + 即时执行 |
| moltclaw | TypeScript | >1GB | 云服务器 | 多 Agent 路由 + 全局编排 |
| nanoclaw | Python | >100MB | ¥350 SBC | 丰富 Python 生态的中量级 Agent |
| microclaw | C/Rust | <1MB | ¥15 MCU | 传感器级别的微型 Agent |

## 真实场景

- **¥600 套件** → 替代 ¥30 万/年的数据中心夜班运维（节省 99%）
- **¥950 套件** → 替代 ¥20 万/年的鱼塘巡塘工（节省 99%）
- **¥680 套件** → 替代 ¥12 万/年的温室管理员（节省 98%）

每个场景都有完整的硬件 BOM、接线图、实现逻辑和成本对比文档。

## 自进化监控策略（CGEP 基因进化协议）

Clawland 与传统 IoT 的核心区别在于 **基因进化协议 (CGEP)**。当 Agent 成功处理一个新情况时，经验会被固化为"基因"——一个可复用的监控策略模式。随着时间推移，Agent 通过类似生物进化的模型自主优化策略：基因通过反复验证获得信心值，高信心基因可以通过 Fleet 管理在节点间共享。

## Build to Earn 模式

这不是传统的开源捐赠模式，而是一个商业合作：

- **净产品收入的 20%** 进入贡献者收入池
- **按季度分红**，通过 Open Collective 完全透明
- **积分制**：代码、Review、技能包、硬件设计、文档都计分
- **前 10 名核心维护者** 额外获得 0.5% 终身收入份额

## 招募

- **Go 开发者**：picclaw 边缘 Agent（ARM 交叉编译、GPIO 集成）
- **TypeScript 开发者**：moltclaw 云网关（多 Agent 路由）
- **Python 开发者**：nanoclaw + 传感器驱动库
- **硬件玩家**：传感器套件设计与测试
- **技术写作**：文档、教程、部署指南

## 链接

- GitHub 组织：https://github.com/Clawland-AI
- Bounty Board（付费任务）：https://github.com/orgs/Clawland-AI/projects/1
- 收入分享条款：https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md
- 讨论区：https://github.com/Clawland-AI/.github/discussions/2

欢迎讨论和提问！
