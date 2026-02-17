<div align="center">
<img src="assets/logo.jpg" alt="PicClaw" width="512">

<h1>PicClaw</h1>

<h3>Self-Evolving Edge AI Agent in Go</h3>

<p>
<img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go&logoColor=white" alt="Go 1.24">
<img src="https://img.shields.io/badge/RAM-<10MB-brightgreen" alt="<10MB RAM">
<img src="https://img.shields.io/badge/Hardware-$10+-orange" alt="$10+ Hardware">
<img src="https://img.shields.io/badge/Arch-x86__64%20|%20ARM64%20|%20RISC--V-blue" alt="Arch">
<img src="https://img.shields.io/badge/license-MIT-green" alt="MIT License">
</p>

</div>

---

PicoClaw is an ultra-lightweight AI agent written in Go that runs on $10 hardware with <10MB RAM. It connects to LLM providers, executes tools, talks through messaging channels, and **evolves its monitoring strategies over time** through a built-in Gene Evolution Protocol.

Originally inspired by [nanobot](https://github.com/HKUDS/nanobot), PicoClaw was refactored from the ground up in Go through a self-bootstrapping process where the AI agent drove the entire migration.

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

## What It Actually Does

PicoClaw is a single Go binary that:

1. **Connects to an LLM** (Zhipu, OpenRouter, Anthropic, OpenAI, Gemini, Groq, or any vLLM endpoint)
2. **Executes tools** (12 built-in: file I/O, shell exec, web search/fetch, cron, messaging, gene reporting)
3. **Talks through channels** (Telegram, Discord, QQ, DingTalk, Feishu, WhatsApp, MaixCam)
4. **Learns from experience** via the Gene Evolution Protocol — successful monitoring strategies are solidified into reusable "genes"
5. **Reports to fleet** via an optional Edge API server for multi-node coordination

|  | OpenClaw (TS) | NanoBot (Python) | **PicoClaw (Go)** |
|---|---|---|---|
| **RAM** | >1 GB | >100 MB | **< 10 MB** |
| **Boot** (0.8 GHz) | >500 s | >30 s | **< 1 s** |
| **Min Hardware** | Mac Mini $599 | SBC ~$50 | **Any Linux $10** |
| **Binary** | Node runtime | Python runtime | **Single static binary** |

## Architecture

```
cmd/picoclaw/main.go          CLI entry point (onboard|agent|gateway|status|cron|skills|gene|version)
pkg/
  agent/                       Agent loop, context builder, memory manager
  bus/                         Internal message bus
  channels/                    7 messaging channels (telegram, discord, qq, dingtalk, feishu, whatsapp, maixcam)
  config/                      JSON + env var configuration (agents, channels, providers, gateway, edge, gene)
  cron/                        Cron scheduler service
  edge/                        Edge API server + fleet heartbeat reporter
  gene/                        Gene Evolution Protocol (CGEP) engine
    ├── types.go               Gene, Capsule, EvolutionEvent, StrategyPreset structs
    ├── store.go               JSON file persistence (genes.json, capsules.json, events.jsonl)
    ├── selector.go            Gene scoring, selection, drift intensity computation
    ├── solidify.go            Experience solidification into Capsules and Genes
    ├── signals.go             Signal extraction from text/memory (8 signal types)
    ├── engine.go              Top-level orchestrator
    ├── seeds.go               6 built-in seed genes
    └── hash.go                SHA-256 content-addressable asset IDs
  heartbeat/                   Periodic heartbeat service
  logger/                      Structured JSON logger
  providers/                   Unified HTTP LLM provider (OpenRouter, Anthropic, OpenAI, Gemini, Zhipu, Groq, vLLM)
  session/                     Conversation session manager
  skills/                      Skill loader + installer (from Git repos)
  tools/                       12 tools: read_file, write_file, edit_file, append_file, list_dir,
                                         exec, spawn, web_search, web_fetch, message, cron, report_gene
  utils/                       String utilities
  voice/                       Groq Whisper voice transcription
skills/                        6 built-in skills (datacenter-monitoring, github, summarize, weather, tmux, skill-creator)
```

## Quick Start

**1. Build**

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
make build          # Single binary in build/
# or: make build-all  for linux/darwin x amd64/arm64/riscv64
```

**2. Initialize**

```bash
picoclaw onboard
```

**3. Configure** (`~/.picoclaw/config.json`)

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

Minimum config: pick one provider and set its API key. Everything else has defaults.

**4. Use**

```bash
picoclaw agent -m "What is 2+2?"     # One-shot query
picoclaw agent                         # Interactive chat
picoclaw gateway                       # Start with channels + cron + edge
```

## LLM Providers

All providers use a unified HTTP interface. Pick one:

| Provider | Model prefix | Get API Key |
|----------|-------------|-------------|
| **Zhipu** | `glm-*` | [bigmodel.cn](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **OpenRouter** | `openrouter/*`, `anthropic/*`, `openai/*`, etc. | [openrouter.ai](https://openrouter.ai/keys) |
| **Anthropic** | `claude-*` | [console.anthropic.com](https://console.anthropic.com) |
| **OpenAI** | `gpt-*` | [platform.openai.com](https://platform.openai.com) |
| **Gemini** | `gemini-*` | [aistudio.google.com](https://aistudio.google.com) |
| **Groq** | `groq/*` (also enables voice transcription) | [console.groq.com](https://console.groq.com) |
| **vLLM** | any (set `api_base`) | Self-hosted |

## Messaging Channels

Start the gateway to connect channels:

```bash
picoclaw gateway
```

| Channel | Config keys | Notes |
|---------|------------|-------|
| **Telegram** | `token`, `allow_from` | Recommended. Supports voice messages (with Groq) |
| **Discord** | `token`, `allow_from` | Enable MESSAGE CONTENT INTENT |
| **QQ** | `app_id`, `app_secret` | QQ Open Platform |
| **DingTalk** | `client_id`, `client_secret` | Internal app |
| **Feishu** | `app_id`, `app_secret`, `encrypt_key`, `verification_token` | Lark/Feishu bot |
| **WhatsApp** | `bridge_url` | Via WhatsApp bridge |
| **MaixCam** | `host`, `port` | Direct hardware connection |

<details>
<summary><b>Telegram example config</b></summary>

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456:ABC-DEF...",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```

Get a token from `@BotFather` on Telegram. Get your user ID from `@userinfobot`.

</details>

<details>
<summary><b>Full config example</b></summary>

See [config.example.json](config.example.json) for all available options.

</details>

## Tools (12 built-in)

| Tool | Description |
|------|-------------|
| `read_file` | Read file contents |
| `write_file` | Create or overwrite a file |
| `edit_file` | Search-and-replace edit in a file |
| `append_file` | Append content to a file |
| `list_dir` | List directory contents |
| `exec` | Execute shell commands |
| `spawn` | Run background processes |
| `web_search` | Search the web (Brave Search API) |
| `web_fetch` | Fetch and extract web page content |
| `message` | Send messages through configured channels |
| `cron` | Create/manage scheduled tasks |
| `report_gene` | Report a learned experience to the Gene Evolution system |

## Gene Evolution Protocol (CGEP)

PicoClaw includes a self-evolving strategy system. Agents learn from experience and reuse proven monitoring strategies.

**How it works:**

1. **Signals** are extracted from sensor data, memory, and daily notes (8 types: `sensor_error`, `threshold_breach`, `cross_sensor_anomaly`, `new_pattern_detected`, `response_too_slow`, `strategy_proven`, `unknown_situation`, `time_pattern`)
2. Active signals are **matched against the gene pool** — the best-scoring genes are injected into the LLM system prompt
3. After handling, experiences are **solidified** into Capsules. Successful novel handling creates new Genes.
4. Gene confidence adjusts over time: success increases it, failure decreases it
5. High-confidence genes can be **published to Fleet** for cross-node sharing

**6 seed genes ship by default:**

| Gene | Category | Scenario | What it does |
|------|----------|----------|-------------|
| `gene_sensor_repair_from_errors` | repair | generic | Retry failed sensors, switch to backup |
| `gene_threshold_optimize_from_history` | optimize | generic | Adjust thresholds using 7-day statistics |
| `gene_cross_sensor_innovate` | innovate | generic | Detect cross-sensor correlations |
| `gene_dc_temp_spike_crosscheck` | repair | datacenter | Cross-check rack temps before alerting |
| `gene_pond_do_night_emergency` | repair | aquaculture | Night-time dissolved oxygen emergency |
| `gene_greenhouse_ventilation_schedule` | optimize | greenhouse | Optimize vent timing from rolling data |

**Strategy presets** control which genes are applied:

| Preset | Min Confidence | Min Verified | Use case |
|--------|---------------|-------------|----------|
| `conservative` | 90% | 5 | Production critical systems |
| `balanced` | 70% | 2 | Default |
| `exploratory` | 30% | 0 | Testing new scenarios |
| `repair-only` | 50% | 1 | Emergency mode |

**Configure in `~/.picoclaw/config.json`:**

```json
{
  "gene": {
    "strategy": "balanced",
    "auto_publish": true,
    "min_confidence": 0.7,
    "min_verified_by": 3
  }
}
```

**CLI:**

```bash
picoclaw gene list      # Show genes with confidence scores
picoclaw gene stats     # Pool statistics and drift intensity
picoclaw gene export    # Export genes as JSON
```

## Edge Server

PicoClaw can act as an L1 edge node in a distributed architecture, reporting to upstream fleet managers:

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

When enabled, PicoClaw exposes:
- `GET /healthz` — Health check
- `GET /api/v1/status` — Node status
- `POST /api/v1/command` — Receive commands from fleet

And periodically sends heartbeats (including gene stats) to the configured fleet manager.

## Skills (6 built-in)

Skills are markdown files that teach the agent domain-specific knowledge.

| Skill | Description |
|-------|-------------|
| `datacenter-monitoring` | Rack temperature monitoring with mock sensors and Gene Evolution demo |
| `github` | GitHub workflow assistance |
| `summarize` | Text summarization |
| `weather` | Weather lookup (no API key needed) |
| `tmux` | Terminal multiplexer management |
| `skill-creator` | Create new skills |

```bash
picoclaw skills list           # List installed skills
picoclaw skills install <url>  # Install from Git repo
picoclaw skills install-builtin  # Install all built-in skills
```

## Datacenter Monitoring Demo

The `datacenter-monitoring` skill includes a complete end-to-end demo:

```bash
# Test with mock sensors
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --all --summary

# Simulate a temperature spike
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --rack A1 --spike

# Simulate a sensor failure
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --rack B2 --fail

# Random chaos mode (10% failures, 15% spikes)
python3 skills/datacenter-monitoring/scripts/mock-sensor.py --all --chaos --summary
```

See [skills/datacenter-monitoring/DEPLOY.md](skills/datacenter-monitoring/DEPLOY.md) for full deployment guide with real hardware.

## CLI Reference

| Command | Description |
|---------|-------------|
| `picoclaw onboard` | Initialize config and workspace |
| `picoclaw agent -m "..."` | One-shot chat |
| `picoclaw agent` | Interactive chat mode |
| `picoclaw gateway` | Start gateway (channels + cron + edge + heartbeat) |
| `picoclaw status` | Show configuration status |
| `picoclaw cron list` | List scheduled jobs |
| `picoclaw cron add` | Add a scheduled job |
| `picoclaw skills list` | List installed skills |
| `picoclaw skills install <url>` | Install skill from Git |
| `picoclaw gene list` | Show local genes |
| `picoclaw gene stats` | Gene pool statistics |
| `picoclaw gene export` | Export genes as JSON |
| `picoclaw version` | Show version |

## Hardware Targets

PicoClaw runs on anything with Linux and a network connection:

| Device | Price | Arch | Use case |
|--------|-------|------|----------|
| [LicheeRV Nano](https://www.aliexpress.com/item/1005006519668532.html) | $10 | RISC-V | Minimal edge agent |
| [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) | $30-50 | RISC-V | Server monitoring |
| [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) | $50-100 | RISC-V | Vision + AI monitoring |
| Raspberry Pi Zero 2W | $15 | ARM64 | General purpose |
| Any Linux x86/ARM server | - | x86/ARM | Cloud or on-prem |

## Demonstrations

<table align="center">
  <tr align="center">
    <th>Full-Stack Engineering</th>
    <th>Logging & Planning</th>
    <th>Web Search</th>
  </tr>
  <tr>
    <td align="center"><img src="assets/picoclaw_code.gif" width="240" height="180"></td>
    <td align="center"><img src="assets/picoclaw_memory.gif" width="240" height="180"></td>
    <td align="center"><img src="assets/picoclaw_search.gif" width="240" height="180"></td>
  </tr>
</table>

## Workspace Layout

```
~/.picoclaw/
├── config.json               # Main configuration
└── workspace/
    ├── genes/                # Gene Evolution data
    │   ├── genes.json        #   Strategy definitions
    │   ├── capsules.json     #   Solidified experiences
    │   └── events.jsonl      #   Audit trail
    ├── sessions/             # Conversation history
    ├── memory/               # Long-term memory (MEMORY.md + daily notes)
    ├── cron/                 # Scheduled jobs
    ├── skills/               # Installed skills
    ├── AGENTS.md             # Agent behavior guide
    ├── IDENTITY.md           # Agent identity
    └── USER.md               # User preferences
```

## Clawland Ecosystem

PicoClaw is the **L1 edge agent** in the [Clawland-AI](https://github.com/Clawland-AI) open-source ecosystem:

| Layer | Repo | Language | Role |
|-------|------|----------|------|
| **L0 MCU** | [microclaw](https://github.com/Clawland-AI/microclaw) | C/Rust | Sensor-level agent on $2 MCUs |
| **L1 Edge** | **picclaw** (this repo) | Go | Edge AI agent on $10 hardware |
| **L2 Regional** | [nanoclaw](https://github.com/Clawland-AI/nanoclaw) | Python | Regional gateway on $50 SBCs |
| **L3 Cloud** | [moltclaw](https://github.com/Clawland-AI/moltclaw) | TypeScript | Cloud AI gateway + fleet management |

**Build to Earn**: Contributors share 20% of net revenue through the Contributor Revenue Pool. See [Clawland-AI](https://github.com/Clawland-AI) for details.

## Contributing

PRs welcome. The codebase is ~7500 lines of Go across 15 packages, intentionally small and readable.

```bash
make test       # Run all tests
make lint       # Run golangci-lint
make build-all  # Build for all platforms
```

## Troubleshooting

**Web search returns "API configuration error"** — This is normal without a Brave Search API key. Get a free key at [brave.com/search/api](https://brave.com/search/api) (2000 free queries/month) and add it to `tools.web.search.api_key`.

**"Conflict: terminated by other getUpdates"** — Another instance of the Telegram bot is running. Only one `picoclaw gateway` per bot token.

**Content filtering errors** — Some providers (Zhipu) have content filtering. Try rephrasing or switching models.

## License

MIT License. See [LICENSE](LICENSE).

Based on [nanobot](https://github.com/HKUDS/nanobot) by HKUDS.
