# Clawland：构建边缘 AI 未来，为之获得报酬

**TL;DR**：我们正在构建一个开源 AI Agent 生态系统，在 $10 硬件上运行，替代每年 $50,000 的监控工作。核心贡献者共享所有产品收入的 20%。我们正在招募前 10 名核心维护者。

---

## 问题

一名夜班数据中心操作员每年成本 $48,000，负责盯屏幕并重启崩溃的服务。一名鱼塘巡塘工每年收入 $30,000，凌晨 3 点检查溶氧水平。一名大棚管理员每年收入 $18,000，负责在天热时开通风口。

这些工作很重要。但它们重复、令人筋疲力尽，而且非常适合自动化。

然而，现有的 IoT 解决方案昂贵（$5,000–50,000），需要系统集成商，且只能进行简单的 if-then 告警。

## 解决方案：Clawland

我们构建了一系列能在低至 $10 的硬件上运行的 AI Agent：

| Agent | 语言 | 内存 | 硬件 | 用途 |
|-------|------|------|------|------|
| **picclaw** | Go | <10MB | $10 RISC-V 开发板 | 边缘感知 + 本地决策 + 即时动作 |
| **moltclaw** | TypeScript | >1GB | 云服务器 | 多 Agent 路由 + 全局编排 |
| **nanoclaw** | Python | >100MB | $50 SBC | 丰富生态系统，中等负载任务 |
| **microclaw** | C/Rust | <1MB | $2 MCU | 传感器级微 Agent |

结合 $5–50 的传感器，这些 Agent 可以：

- **感知**：读取温度、湿度、溶解氧、振动、烟雾、运动（通过 `exec` + Python/GPIO）
- **思考**：基于 LLM 的多传感器分析与趋势预测（支持 8+ 个 LLM 提供商）
- **行动**：控制继电器、阀门、开关（通过 `exec` + GPIO）
- **报告**：通过 Telegram、Discord、飞书、钉钉、WhatsApp、QQ 发送警报（7 个渠道）
- **学习**：持久记忆，随时间变得更智能
- **调度**：基于 Cron 的自主监控循环

## 实际数字

| 套件 | 硬件成本 | 替代 | 传统成本 | 节省 |
|------|---------|------|---------|------|
| 数据中心哨兵 | $88 | 夜班操作员 | $48,000/年 | 99% |
| 池塘守护者 | $133/池塘 | 巡塘工人 | $30,000/年 | 99% |
| 大棚管理员 | $95/大棚 | 管理员 | $18,000/年 | 98% |
| 安全守护者 | $30 | 老年护理订阅 | $900/年 | 90% |
| 设备医生 | $32/机器 | 预测性维护服务 | $3,200/年 | 95% |

包含 BOM、实现逻辑和 ROI 的完整分析：[clawland-kits](https://github.com/Clawland-AI/clawland-kits)

## 为什么开源？为什么收入分成？

因为最好的边缘 AI 系统将由那些实际在鱼塘、工厂和机房部署的人构建，而不是旧金山办公室的一家初创公司。

**我们不是在请求捐款。我们提供一笔交易：**

> 你构建它。它销售出去。你获得 20%。

### 工作原理

1. **20% 的净产品收入**进入贡献者收入池
2. **季度付款**基于透明的贡献积分系统
3. **积分**为以下工作获取：代码、审查、技能、硬件设计、文档、社区支持
4. **Open Collective** 提供完全的财务透明度
5. **前 10 名核心维护者**额外获得 0.5% 的终身收入分成

### 示例

如果 Clawland 在 Q3 产生 $50,000 净收入：
- 贡献者池：$10,000
- 贡献了 1000 总积分中 500 积分的开发者：**本季度获得 $5,000**

阅读完整条款：[CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)

## 我们需要什么

### 当下（前 10 名核心维护者）

- **Go 开发者**：协助构建 picclaw（边缘 Agent、工具、渠道）
- **TypeScript 开发者**：协助构建 moltclaw（云网关、舰队管理）
- **Python 开发者**：协助构建 nanoclaw + 传感器驱动库
- **硬件极客**：设计和测试传感器套件，编写驱动脚本
- **技术文档撰写者**：文档、教程、部署指南

### 适合新手的 Issue

每个仓库都维护 `good-first-issue` 标签。从这些开始。

### 赏金看板

我们维护一个带有定价任务的[赏金看板](https://github.com/orgs/Clawland-AI/projects/1)。认领一项，交付，获得报酬。

## 治理

Clawland 由技术指导委员会（TSC）治理，采用透明的基于 RFC 的决策流程。参见 [GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md)。

## 开始使用

```bash
# 克隆 picclaw
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw

# 构建
make build

# 运行
./picclaw agent -m "Hello, I'm your new edge AI assistant"
```

## 链接

- GitHub：[github.com/Clawland-AI](https://github.com/Clawland-AI)
- 收入分成条款：[CONTRIBUTOR-REVENUE-SHARE.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTOR-REVENUE-SHARE.md)
- 治理：[GOVERNANCE.md](https://github.com/Clawland-AI/.github/blob/main/GOVERNANCE.md)
- 贡献：[CONTRIBUTING.md](https://github.com/Clawland-AI/.github/blob/main/CONTRIBUTING.md)

---

*"让 $10 硬件 + AI Agent 完成监控，让人类专注思考。"*

**我们正在构建一个未来：每个鱼塘、每个服务器机房、每个大棚都有一个 AI 守护者。加入我们。**
