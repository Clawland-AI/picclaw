<div align="center">
<img src="assets/logo.jpg" alt="PicClaw" width="512">

<h1>PicClaw</h1>

<h3>用 Go 编写的自进化边缘 AI Agent</h3>

<p>
<img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go&logoColor=white" alt="Go 1.24">
<img src="https://img.shields.io/badge/RAM-<10MB-brightgreen" alt="<10MB RAM">
<img src="https://img.shields.io/badge/Hardware-$10+-orange" alt="$10+ Hardware">
<img src="https://img.shields.io/badge/Arch-x86__64%20|%20ARM64%20|%20RISC--V-blue" alt="Arch">
<img src="https://img.shields.io/badge/license-MIT-green" alt="MIT License">
</p>

[English](README.md) | [简体中文](README.zh-CN.md)

</div>

---

PicoClaw 是一个用 Go 编写的超轻量 AI Agent，可以在 10 美元级硬件上运行，内存占用低于 10MB。它可以连接多种 LLM 服务，执行本地工具，通过聊天渠道交互，并通过内置的 Gene Evolution Protocol 在长期运行中逐步沉淀和复用有效的监控策略。

PicoClaw 最初受到 [nanobot](https://github.com/HKUDS/nanobot) 启发，但已经通过一次自举式迁移从头重构为 Go 版本：迁移过程本身由 AI Agent 驱动完成。

<table align="center">
  <tr align="center">
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/picoclaw_mem.gif" width="360" height="240">
      </p>
    </td>
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/licheervnano.png" width="400" height="240">
      </p>
    </td>
  </tr>
</table>

## 它实际能做什么

PicoClaw 是一个单文件 Go 二进制程序，核心能力包括：

1. **连接 LLM**：支持 Zhipu、OpenRouter、Anthropic、OpenAI、Gemini、Groq，以及任意兼容 vLLM 的 endpoint。
2. **执行工具**：内置 12 个工具，包括文件读写、shell 执行、网页搜索/抓取、定时任务、消息发送和 gene 上报。
3. **接入消息渠道**：支持 Telegram、Discord、QQ、钉钉、飞书、WhatsApp 和 MaixCam。
4. **从经验中学习**：通过 Gene Evolution Protocol，把有效的监控策略固化为可复用的 genes。
5. **接入边缘集群**：可选 Edge API Server 支持多节点协调和 fleet 状态上报。

|  | OpenClaw (TS) | NanoBot (Python) | **PicoClaw (Go)** |
|---|---|---|---|
| **内存** | >1 GB | >100 MB | **< 10 MB** |
| **启动时间** (0.8 GHz) | >500 s | >30 s | **< 1 s** |
| **最低硬件** | Mac Mini $599 | SBC ~$50 | **任意 $10 Linux 设备** |
| **交付形态** | Node runtime | Python runtime | **单个静态二进制文件** |

## 架构

```
cmd/picoclaw/main.go          CLI 入口 (onboard|agent|gateway|status|cron|skills|gene|version)
pkg/
  agent/                       Agent 循环、上下文构建、记忆管理
  bus/                         内部消息总线
  channels/                    7 个消息渠道 (telegram, discord, qq, dingtalk, feishu, whatsapp, maixcam)
  config/                      JSON + 环境变量配置 (agents, channels, providers, gateway, edge, gene)
  cron/                        定时任务服务
  edge/                        Edge API Server + fleet heartbeat reporter
  gene/                        Gene Evolution Protocol (CGEP) 引擎
  heartbeat/                   周期性心跳服务
  logger/                      结构化 JSON 日志
  providers/                   统一 HTTP LLM provider
  session/                     会话管理
  skills/                      Skill 加载和安装
  tools/                       12 个内置工具
  utils/                       字符串工具
  voice/                       Groq Whisper 语音转写
skills/                        6 个内置 skills
```

## 快速开始

**1. 构建**

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
make build          # 在 build/ 目录生成单个二进制文件
# 或者：make build-all  为 linux/darwin x amd64/arm64/riscv64 构建
```

**2. 初始化**

```bash
picoclaw onboard
```

**3. 配置** (`~/.picoclaw/config.json`)

```json
{
  "agents": {
    "defaults": {
      "model": "glm-4.7",
      "max_tokens": 8192
    }
  },
  "providers": {
    "zhipu": {
      "api_key": "YOUR_API_KEY"
    }
  }
}
```

最低配置只需要选择一个 provider 并填写 API key；其他选项都有默认值。

**4. 使用**

```bash
picoclaw agent -m "What is 2+2?"     # 单次问答
picoclaw agent                       # 交互式聊天
picoclaw gateway                     # 启动 channels + cron + edge
```

## LLM Provider

所有 provider 使用统一 HTTP 接口。可以任选一种：

| Provider | Model prefix | API Key |
|----------|-------------|---------|
| **Zhipu** | `glm-*` | [bigmodel.cn](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **OpenRouter** | `openrouter/*`, `anthropic/*`, `openai/*` 等 | [openrouter.ai](https://openrouter.ai/keys) |
| **Anthropic** | `claude-*` | [console.anthropic.com](https://console.anthropic.com) |
| **OpenAI** | `gpt-*` | [platform.openai.com](https://platform.openai.com) |
| **Gemini** | `gemini-*` | [aistudio.google.com](https://aistudio.google.com) |
| **Groq** | `groq/*`，也可启用语音转写 | [console.groq.com](https://console.groq.com) |
| **vLLM** | 任意模型，需设置 `api_base` | 自托管 |

## 消息渠道

启动 gateway 后即可接入消息渠道：

```bash
picoclaw gateway
```

| Channel | Config keys | Notes |
|---------|-------------|-------|
| **Telegram** | `token`, `allow_from` | 推荐使用；配合 Groq 支持语音消息 |
| **Discord** | `token`, `allow_from` | 需要启用 MESSAGE CONTENT INTENT |
| **QQ** | `app_id`, `app_secret` | QQ Open Platform |
| **DingTalk** | `client_id`, `client_secret` | 内部应用 |
| **Feishu** | `app_id`, `app_secret`, `encrypt_key`, `verification_token` | 飞书/Lark bot |
| **WhatsApp** | `bridge_url` | 通过 WhatsApp bridge 接入 |
| **MaixCam** | `host`, `port` | 直接连接硬件 |

## 内置工具

| Tool | Description |
|------|-------------|
| `read_file` | 读取文件内容 |
| `write_file` | 创建或覆盖文件 |
| `edit_file` | 搜索并替换文件内容 |
| `append_file` | 追加文件内容 |
| `list_dir` | 列出目录内容 |
| `exec` | 执行 shell 命令 |
| `spawn` | 启动后台进程 |
| `web_search` | 搜索网页 |
| `web_fetch` | 抓取并提取网页内容 |
| `message` | 通过已配置渠道发送消息 |
| `cron` | 创建和管理定时任务 |
| `report_gene` | 向 Gene Evolution 系统上报经验 |

## Gene Evolution Protocol (CGEP)

PicoClaw 内置一套自进化策略系统。Agent 会从经验中提取信号，选择合适的 gene 注入系统提示，并把成功经验固化为可复用的策略。

工作流程：

1. 从传感器数据、记忆和每日笔记中提取信号。
2. 将活跃信号与 gene pool 匹配，选出最适合当前情境的 gene。
3. 任务处理结束后，把经验固化为 Capsule；有价值的新经验会生成新的 Gene。
4. Gene 的置信度会随成功和失败动态调整。
5. 高置信度 gene 可以发布到 Fleet，在多个边缘节点之间共享。

默认包含 6 个 seed genes：

| Gene | Category | Scenario | What it does |
|------|----------|----------|-------------|
| `gene_sensor_repair_from_errors` | repair | generic | 传感器失败时重试并切换备用方案 |
| `gene_threshold_optimize_from_history` | optimize | generic | 用 7 天统计数据调整阈值 |
| `gene_cross_sensor_innovate` | innovate | generic | 发现跨传感器相关性 |
| `gene_dc_temp_spike_crosscheck` | repair | datacenter | 告警前交叉检查机柜温度 |
| `gene_pond_do_night_emergency` | repair | aquaculture | 夜间溶解氧紧急处理 |
| `gene_greenhouse_ventilation_schedule` | optimize | greenhouse | 基于滚动数据优化通风时间 |

## Edge Server

PicoClaw 可以作为分布式架构中的 L1 边缘节点，向上游 fleet manager 上报状态：

```json
{
  "edge": {
    "enabled": true,
    "port": 9090,
    "node_id": "dc-rack-a1",
    "node_name": "Datacenter Rack A1",
    "cloud_endpoint": "http://nanoclaw:8080",
    "cloud_token": "...",
    "heartbeat_seconds": 30
  }
}
```

启用后，PicoClaw 会暴露：

- `GET /healthz`：健康检查
- `GET /api/v1/status`：节点状态
- `POST /api/v1/command`：接收 fleet 下发的命令

它还会定期向配置的 fleet manager 发送包含 gene 统计信息的 heartbeat。

## Skills

Skills 是一组 Markdown 文件，用来给 Agent 注入领域知识。

| Skill | Description |
|-------|-------------|
| `datacenter-monitoring` | 机柜温度监控，包含 mock sensor 和 Gene Evolution 演示 |
| `github` | GitHub 工作流辅助 |
| `summarize` | 文本摘要 |
| `weather` | 查询天气，无需 API key |
| `tmux` | 终端复用器管理 |
| `skill-creator` | 创建新的 skill |

```bash
picoclaw skills list             # 列出已安装 skills
picoclaw skills install <url>    # 从 Git 仓库安装 skill
picoclaw skills install-builtin  # 安装所有内置 skills
```

## 数据中心监控演示

`datacenter-monitoring` skill 包含完整的端到端 demo：

```bash
# 使用 mock sensors 测试
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --all --summary

# 模拟温度尖峰
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --rack A1 --spike

# 模拟传感器故障
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --rack B2 --fail

# 随机混沌模式：10% 故障，15% 尖峰
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --all --chaos --summary
```

完整的真实硬件部署说明见 [skills/datacenter-monitoring/DEPLOY.md](skills/datacenter-monitoring/DEPLOY.md)。

## CLI 参考

| Command | Description |
|---------|-------------|
| `picoclaw onboard` | 初始化配置和工作区 |
| `picoclaw agent -m "..."` | 单次聊天 |
| `picoclaw agent` | 交互式聊天模式 |
| `picoclaw gateway` | 启动 gateway，包括 channels、cron、edge 和 heartbeat |
| `picoclaw status` | 查看配置状态 |
| `picoclaw cron list` | 列出定时任务 |
| `picoclaw cron add` | 新增定时任务 |
| `picoclaw skills list` | 列出已安装 skills |
| `picoclaw skills install <url>` | 从 Git 安装 skill |
| `picoclaw gene list` | 查看本地 genes |
| `picoclaw gene stats` | 查看 gene pool 统计 |
| `picoclaw gene export` | 导出 genes 为 JSON |
| `picoclaw version` | 查看版本 |

## 硬件目标

PicoClaw 可以运行在任何带网络连接的 Linux 设备上：

| Device | Price | Arch | Use case |
|--------|-------|------|----------|
| [LicheeRV Nano](https://www.aliexpress.com/item/1005006519668532.html) | $10 | RISC-V | 最小边缘 Agent |
| [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) | $30-50 | RISC-V | 服务器监控 |
| [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) | $50-100 | RISC-V | 视觉 + AI 监控 |
| Raspberry Pi Zero 2W | $15 | ARM64 | 通用场景 |
| Any Linux x86/ARM server | - | x86/ARM | 云端或本地服务器 |

## 工作区结构

```
~/.picoclaw/
├── config.json               # 主配置文件
└── workspace/
    ├── genes/                # Gene Evolution 数据
    │   ├── genes.json        #   策略定义
    │   ├── capsules.json     #   固化经验
    │   └── events.jsonl      #   审计日志
    ├── sessions/             # 对话历史
    ├── memory/               # 长期记忆 (MEMORY.md + daily notes)
    ├── cron/                 # 定时任务
    ├── skills/               # 已安装 skills
    ├── AGENTS.md             # Agent 行为指南
    ├── IDENTITY.md           # Agent 身份
    └── USER.md               # 用户偏好
```

## Clawland 生态

PicoClaw 是 [Clawland-AI](https://github.com/Clawland-AI) 开源生态中的 **L1 边缘 Agent**：

| Layer | Repo | Language | Role |
|-------|------|----------|------|
| **L0 MCU** | [microclaw](https://github.com/Clawland-AI/microclaw) | C/Rust | 运行在 2 美元 MCU 上的传感器级 Agent |
| **L1 Edge** | **picclaw** (当前仓库) | Go | 运行在 10 美元硬件上的边缘 AI Agent |
| **L2 Regional** | [nanoclaw](https://github.com/Clawland-AI/nanoclaw) | Python | 运行在 50 美元 SBC 上的区域网关 |
| **L3 Cloud** | [moltclaw](https://github.com/Clawland-AI/moltclaw) | TypeScript | 云端 AI gateway + fleet management |

**Build to Earn**：贡献者可以通过 Contributor Revenue Pool 分享 20% 净收入。详情见 [Clawland-AI](https://github.com/Clawland-AI)。

## 贡献

欢迎提交 PR。代码库约 7500 行 Go，分布在 15 个 package 中，设计目标是保持小而清晰。

```bash
make test        # 运行全部测试
make lint        # 运行 golangci-lint
make build-all   # 构建所有平台
```

## 常见问题

**Web search 返回 "API configuration error"**：没有配置 Brave Search API key 时这是正常现象。可以在 [brave.com/search/api](https://brave.com/search/api) 获取免费 key（每月 2000 次免费查询），然后写入 `tools.web.search.api_key`。

**"Conflict: terminated by other getUpdates"**：另一个 Telegram bot 实例正在运行。同一个 bot token 只能运行一个 `picoclaw gateway`。

**内容过滤错误**：部分 provider（例如 Zhipu）可能触发内容过滤。可以尝试改写提示词，或切换到其他模型。

## License

MIT License. See [LICENSE](LICENSE).

Based on [nanobot](https://github.com/HKUDS/nanobot) by HKUDS.
