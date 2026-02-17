# OpenClaw（云端大脑）+ PicoClaw（边缘手脚）架构设计

> 基于 PicoClaw 现有代码（Message Bus、Channel 体系、MaixCAM TCP 协议、Spawn 子代理、web_fetch/exec 工具）的实际能力，设计一套可落地的"监控中心 + 现场人员"分布式 Agent 系统。

---

## 一、现状分析：当前 PicoClaw 能做什么、不能做什么

### 已有能力（直接可用）

```
✅ 7 个消息渠道（Telegram/Discord/飞书/钉钉/QQ/WhatsApp/MaixCAM）
✅ web_fetch —— 可调用任何 HTTP API（GET 请求）
✅ exec —— 可执行本地 shell 命令（带安全防护）
✅ cron —— 可定时触发任务
✅ spawn —— 可创建进程内子代理
✅ MaixCAM TCP 通道 —— 硬件设备 JSON 协议对接
✅ Memory 持久化 —— MEMORY.md + 每日笔记
✅ Session 隔离 —— 每个渠道/设备独立会话
✅ 语音转写 —— Groq Whisper 集成
```

### 缺失能力（需要新增）

```
❌ HTTP Server —— PicoClaw 没有 HTTP 入口，无法被外部系统调用
❌ 跨实例通信 —— 多个 PicoClaw 之间不能直接通信
❌ 上报协议 —— 没有向 OpenClaw 中控汇报状态的机制
❌ 指令下发 —— OpenClaw 无法主动向 PicoClaw 下发命令
❌ 设备注册 —— 没有边缘节点自动注册/发现机制
❌ 共享状态 —— 没有分布式 Session/Memory
```

---

## 二、整体架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                     ☁️  云端（OpenClaw 集群）                        │
│                                                                     │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────┐  ┌───────────┐ │
│  │ OpenClaw    │  │ 设备注册中心  │  │ 事件总线   │  │ 决策引擎  │ │
│  │ AI Gateway  │  │ Registry     │  │ Event Hub  │  │ Decision  │ │
│  │ (多Agent    │  │              │  │ (MQTT/     │  │ Engine    │ │
│  │  路由)      │  │ - 节点列表   │  │  Redis     │  │           │ │
│  │            │  │ - 心跳状态   │  │  Pub/Sub)  │  │ - 升级    │ │
│  │ - 意图分拣 │  │ - 能力声明   │  │            │  │ - 聚合    │ │
│  │ - 复杂推理 │  │ - 分组管理   │  │ - 告警聚合 │  │ - 关联    │ │
│  │ - 知识库   │  │              │  │ - 指令下发 │  │ - 调度    │ │
│  └──────┬──────┘  └──────┬───────┘  └─────┬──────┘  └─────┬─────┘ │
│         │               │                │               │        │
│  ───────┴───────────────┴────────────────┴───────────────┴─────── │
│                         统一 API 层                                │
│                    (REST + WebSocket + MQTT)                       │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                    ══════════════╪══════════════
                    互联网 / 内网 / 4G / LoRa
                    ══════════════╪══════════════
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
        ▼                         ▼                         ▼
┌───────────────┐      ┌───────────────┐      ┌───────────────┐
│ 🏭 Edge Node A │      │ 🏠 Edge Node B │      │ 🌾 Edge Node C │
│ (机房巡检)     │      │ (安防监控)     │      │ (农业监控)     │
│                │      │                │      │                │
│ PicoClaw       │      │ PicoClaw       │      │ PicoClaw       │
│ + LicheeRV     │      │ + MaixCAM2     │      │ + LicheeRV     │
│ + 温湿度传感器 │      │ + 摄像头       │      │ + 溶氧传感器   │
│ + UPS 串口     │      │                │      │ + 继电器       │
└───────────────┘      └───────────────┘      └───────────────┘
```

---

## 三、需要新增的 3 个核心模块

### 模块 1：Edge API Server（PicoClaw 侧）

**目的**：让 PicoClaw 能被外部系统调用，同时主动向云端汇报。

**实现位置**：`pkg/edge/server.go`（新增）

```go
// 新增 Edge HTTP Server，在 gateway 模式下与 Channel Manager 并行启动
// 职责：
//   1. POST /api/command  —— 接收云端下发的指令
//   2. POST /api/message  —— 接收云端转发的用户消息
//   3. GET  /api/status   —— 返回节点状态（CPU/内存/传感器/渠道状态）
//   4. GET  /api/health   —— 健康检查（供注册中心心跳探测）

type EdgeServer struct {
    config     EdgeConfig
    bus        *bus.MessageBus
    agentLoop  *agent.AgentLoop
    nodeID     string
    capabilities []string  // 该节点的能力声明
}

// 指令处理流程：
// 云端 POST /api/command → EdgeServer 解析 → 注入 MessageBus (inbound)
// → AgentLoop 处理 → 结果通过 callback 返回 HTTP response
// → 同时通过 Edge Reporter 上报给云端事件总线
```

**配置扩展**（`config.json`）：

```json
{
  "edge": {
    "enabled": true,
    "node_id": "datacenter-rack-01",
    "node_group": "datacenter-monitoring",
    "api_port": 18791,
    "api_token": "edge-secret-token",
    "cloud": {
      "url": "https://openclaw.example.com",
      "api_key": "cloud-api-key",
      "report_interval_seconds": 30,
      "heartbeat_interval_seconds": 10
    },
    "capabilities": ["temperature", "humidity", "ups_status", "shell_exec"]
  }
}
```

### 模块 2：Edge Reporter（PicoClaw 侧 → 云端上报）

**目的**：PicoClaw 主动向 OpenClaw 上报事件和状态。

**实现位置**：`pkg/edge/reporter.go`（新增）

```go
// Reporter 定期或事件驱动地向云端汇报
// 使用 HTTP POST（利用现有 net/http，无需新增依赖）

type EdgeReporter struct {
    cloudURL   string
    apiKey     string
    nodeID     string
    httpClient *http.Client
}

// 上报类型：
// 1. 心跳（定时）—— 节点存活 + 基础指标
// 2. 告警（事件驱动）—— 传感器阈值触发
// 3. 状态（定时）—— 完整状态快照
// 4. 日志（批量）—— 近期操作日志
```

**上报数据格式**：

```json
{
  "node_id": "datacenter-rack-01",
  "node_group": "datacenter-monitoring",
  "type": "alert",
  "severity": "critical",
  "timestamp": "2026-02-12T15:30:00Z",
  "payload": {
    "sensor": "temperature",
    "value": 42.5,
    "threshold": 40.0,
    "message": "机房温度超过阈值: 42.5°C > 40°C",
    "action_taken": "已自动开启备用空调"
  },
  "capabilities": ["temperature", "humidity", "ups_status"],
  "system": {
    "cpu_percent": 12.5,
    "mem_used_mb": 6.2,
    "uptime_hours": 720,
    "agent_sessions": 3
  }
}
```

### 模块 3：Fleet Manager（OpenClaw 侧）

**目的**：云端统一管理所有 PicoClaw 边缘节点。

**实现位置**：OpenClaw 新增 Edge Function 或插件

```
Fleet Manager 职责：
├── 设备注册 —— 边缘节点启动时自动注册（POST /fleet/register）
├── 心跳监控 —— 检测节点在线状态，掉线告警
├── 指令下发 —— 向指定节点/节点组发送命令（POST node_url/api/command）
├── 事件聚合 —— 收集所有节点上报的告警，去重/关联/升级
├── 决策路由 —— 判断哪些事件需要 AI 深度分析 vs 直接转发
└── 仪表盘 —— 所有节点状态实时可视化
```

---

## 四、完整数据流（以"机房监控"为例）

### 场景：机房温度异常 → 自动处置 → 升级告警 → 人工确认

```
时间轴 ─────────────────────────────────────────────────────────────>

Step 1: 边缘感知
─────────────────
PicoClaw (机房节点)
├── cron 每 60 秒触发
├── exec: python3 read_sensor.py --type temperature
├── 返回: {"temperature": 42.5, "humidity": 65}
└── Agent 判断: 42.5°C > 阈值 40°C → 触发告警流程

Step 2: 边缘自主处置（不等云端，直接行动）
──────────────────────────────────────────
PicoClaw (同一节点)
├── exec: python3 control_relay.py --device backup_ac --action on
│   → 返回: "备用空调已开启"
├── exec: ipmitool -H 192.168.1.100 chassis status
│   → 返回: "System Power: on, Drive Fault: false"
├── append_file: workspace/memory/202602/20260212.md
│   → 记录: "15:30 机房温度 42.5°C，已自动开启备用空调"
└── 判断: 温度仍在上升 → 需要升级

Step 3: 上报云端（Edge Reporter → OpenClaw）
─────────────────────────────────────────────
PicoClaw EdgeReporter
├── POST https://openclaw.example.com/fleet/events
│   Body: {
│     "node_id": "datacenter-rack-01",
│     "type": "alert",
│     "severity": "critical",
│     "payload": {
│       "sensor": "temperature",
│       "value": 42.5,
│       "threshold": 40.0,
│       "action_taken": "已自动开启备用空调",
│       "trend": "still_rising"
│     }
│   }
└── 继续本地监控（不阻塞）

Step 4: 云端 AI 决策（OpenClaw）
────────────────────────────────
OpenClaw Fleet Manager
├── 收到告警事件
├── 关联分析:
│   ├── 查询同机房其他节点状态
│   │   → POST datacenter-rack-02/api/status → 温度正常
│   │   → POST datacenter-rack-03/api/status → 温度也偏高
│   ├── 判断: rack-01 和 rack-03 同区域，可能是该区域空调故障
│   └── 决策: 升级为 P1 告警 + 通知物业检查空调
│
├── 下发关联指令:
│   ├── POST datacenter-rack-03/api/command
│   │   Body: {"action": "enable_backup_cooling"}
│   └── POST datacenter-rack-01/api/command
│       Body: {"action": "increase_monitoring_frequency", "interval": 30}
│
└── 多渠道通知:
    ├── 飞书群: "@运维组 P1告警: A区机房温度异常，rack-01(42.5°C) rack-03(41.2°C)，疑似空调故障，已自动开启备用制冷，请物业现场检查"
    ├── Telegram: 给值班主管发详细报告
    └── 钉钉: 创建工单

Step 5: 持续监控 + 闭环
──────────────────────
PicoClaw (rack-01)
├── 监控频率已提升到 30 秒
├── 15:32 温度 41.8°C → 下降趋势 → 上报云端
├── 15:35 温度 40.2°C → 继续下降
├── 15:40 温度 38.5°C → 恢复正常
└── 上报: {"type": "alert_resolved", "duration_minutes": 10}

OpenClaw
├── 收到恢复事件
├── 飞书群: "✅ A区机房温度已恢复正常 (38.5°C)，持续异常 10 分钟"
├── 指令: POST rack-01/api/command → 恢复 60 秒监控频率
└── 归档工单
```

---

## 五、通信协议详细设计

### 5.1 边缘 → 云端（上行）

#### 心跳协议

```
PicoClaw 每 10 秒 → POST {cloud_url}/fleet/heartbeat

Request:
{
  "node_id": "datacenter-rack-01",
  "node_group": "datacenter-monitoring",
  "timestamp": "2026-02-12T15:30:00Z",
  "uptime_seconds": 2592000,
  "system": {
    "cpu_percent": 12.5,
    "mem_used_mb": 6.2,
    "disk_percent": 45.0
  },
  "channels": {
    "maixcam": { "connected": true, "devices": 2 },
    "telegram": { "connected": true }
  },
  "agent": {
    "active_sessions": 3,
    "pending_cron_jobs": 5,
    "last_activity": "2026-02-12T15:29:45Z"
  }
}

Response:
{
  "ack": true,
  "commands_pending": 1  // 有待执行的指令，边缘节点应拉取
}
```

#### 事件上报协议

```
PicoClaw 事件驱动 → POST {cloud_url}/fleet/events

Request:
{
  "node_id": "datacenter-rack-01",
  "events": [
    {
      "id": "evt-20260212-153000-001",
      "type": "sensor_alert",
      "severity": "critical",     // info / warning / critical / emergency
      "timestamp": "2026-02-12T15:30:00Z",
      "source": "temperature_sensor_1",
      "payload": { ... },
      "action_taken": "已自动开启备用空调",
      "requires_escalation": true
    }
  ]
}

Response:
{
  "received": 1,
  "escalated": ["evt-20260212-153000-001"]
}
```

### 5.2 云端 → 边缘（下行）

#### 指令下发协议

```
OpenClaw → POST {node_url}/api/command

Request:
{
  "command_id": "cmd-20260212-153005-001",
  "type": "exec",           // exec / config / query / update
  "priority": "high",       // low / normal / high / urgent
  "payload": {
    "action": "increase_monitoring_frequency",
    "params": { "interval": 30 }
  },
  "timeout_seconds": 60,
  "callback_url": "https://openclaw.example.com/fleet/command-result"
}

Response (同步):
{
  "command_id": "cmd-20260212-153005-001",
  "status": "accepted",     // accepted / rejected / executing
  "estimated_seconds": 5
}

Callback (异步):
POST {callback_url}
{
  "command_id": "cmd-20260212-153005-001",
  "node_id": "datacenter-rack-01",
  "status": "completed",    // completed / failed / timeout
  "result": { ... },
  "duration_ms": 1200
}
```

### 5.3 断网容错机制

```
┌──────────────────────────────────────────────────────┐
│ PicoClaw 边缘节点的离线自治能力                        │
│                                                      │
│ 正常模式（有网络）:                                    │
│   感知 → 本地判断 → 本地执行 → 上报云端 → 云端决策     │
│                                                      │
│ 离线模式（断网）:                                      │
│   感知 → 本地判断 → 本地执行 → 缓存事件到本地文件       │
│                     ↓                                │
│              Skills 中定义的离线 SOP                   │
│              Memory 中积累的历史经验                    │
│                     ↓                                │
│              自主决策（有限范围内）                      │
│                                                      │
│ 恢复模式（网络恢复）:                                  │
│   批量上传缓存事件 → 云端补充分析 → 同步状态            │
└──────────────────────────────────────────────────────┘
```

实现方式：

```go
// pkg/edge/reporter.go 中的离线缓存逻辑

func (r *EdgeReporter) Report(event Event) error {
    err := r.sendToCloud(event)
    if err != nil {
        // 网络不可达 → 写入本地缓存
        r.cacheLocally(event)
        logger.Warn("Cloud unreachable, event cached locally")
        return nil  // 不阻塞本地流程
    }
    return nil
}

func (r *EdgeReporter) FlushCache() {
    // 网络恢复后批量上传
    events := r.loadCachedEvents()
    for _, event := range events {
        if err := r.sendToCloud(event); err == nil {
            r.removeCachedEvent(event.ID)
        }
    }
}
```

---

## 六、节点分组与多层级拓扑

### 6.1 逻辑分组

```yaml
fleet:
  groups:
    datacenter-monitoring:        # 机房监控组
      nodes: [rack-01, rack-02, rack-03, rack-04]
      alert_channel: feishu://group/ops-dc
      escalation: [ops-lead-telegram, ops-manager-phone]

    warehouse-security:           # 仓库安防组
      nodes: [gate-01, gate-02, indoor-cam-01, indoor-cam-02]
      alert_channel: dingtalk://group/warehouse-security
      escalation: [security-lead-telegram]

    farm-aquaculture:             # 水产养殖组
      nodes: [pond-01, pond-02, pond-03]
      alert_channel: telegram://group/farm-ops
      escalation: [farmer-phone]
      critical_auto_action: true  # 关键告警允许自动执行（如开增氧机）
```

### 6.2 多层级拓扑（大规模部署）

```
                    ┌──────────────┐
                    │  OpenClaw    │
                    │  总控中心    │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────┴──────┐ ┌──┴──────┐ ┌──┴──────────┐
       │ 区域网关 A   │ │ 区域网关 B │ │ 区域网关 C   │
       │ (PicoClaw    │ │ (PicoClaw │ │ (PicoClaw    │
       │  树莓派 4)   │ │  树莓派 4)│ │  树莓派 4)   │
       └──────┬───────┘ └────┬─────┘ └──────┬───────┘
              │              │              │
        ┌─────┼─────┐       │        ┌─────┼─────┐
        │     │     │       │        │     │     │
       ┌┴┐  ┌┴┐  ┌┴┐     ┌┴┐     ┌┴┐  ┌┴┐  ┌┴┐
       │01│  │02│  │03│    │04│    │05│  │06│  │07│
       └──┘  └──┘  └──┘    └──┘    └──┘  └──┘  └──┘
       $10   $10   $10     $10     $10   $10   $10
    (PicoClaw 末端节点，LicheeRV-Nano / MaixCAM)
```

**三层架构**：

| 层级 | 硬件 | 角色 | 实例 |
|------|------|------|------|
| L1 末端 | $10 LicheeRV-Nano / MaixCAM | 传感器采集 + 本地执行 + 即时告警 | PicoClaw |
| L2 区域 | $50 树莓派 4 | 区域聚合 + 跨节点关联 + 离线自治 | PicoClaw (增强) |
| L3 云端 | 云服务器 / Mac Mini | 全局决策 + 知识库 + 多渠道通知 | OpenClaw |

**L2 区域网关的价值**：
- 局域网内通信（节点间延迟 <1ms）
- 断网时仍可区域内协调
- 减少云端 API 调用量（聚合后上报）
- 跨节点关联分析（如：rack-01 和 rack-03 同时温度异常 → 区域空调故障）

---

## 七、实现路线图

### Phase 1：最小可行版（2-3 周）

**目标**：单个 PicoClaw 节点能向 OpenClaw 上报 + 接收指令。

```
新增文件：
├── pkg/edge/
│   ├── server.go      # HTTP Server（接收云端指令）
│   ├── reporter.go    # 事件上报器（向云端汇报）
│   └── config.go      # Edge 配置结构
│
修改文件：
├── cmd/picoclaw/main.go    # gateway 模式增加 Edge Server 启动
├── pkg/config/config.go    # 增加 edge 配置段
└── config.example.json     # 增加 edge 配置示例
```

**关键实现**：

```go
// cmd/picoclaw/main.go — gateway 模式修改
func runGateway(cfg *config.Config) {
    // ... 现有代码 ...

    // 新增: 启动 Edge Server
    if cfg.Edge.Enabled {
        edgeServer := edge.NewServer(cfg.Edge, messageBus, agentLoop)
        go edgeServer.Start(ctx)

        // 新增: 启动 Edge Reporter
        reporter := edge.NewReporter(cfg.Edge, cfg.Edge.Cloud)
        go reporter.StartHeartbeat(ctx)

        // 将 reporter 注入 agent 工具链，让 agent 可以主动上报
        toolRegistry.Register(edge.NewReportTool(reporter))
    }
}
```

### Phase 2：多节点管理（3-4 周）

**目标**：OpenClaw 能管理多个 PicoClaw 节点，跨节点关联。

```
OpenClaw 侧新增：
├── Edge Function: fleet-manager
│   ├── /fleet/register     # 节点注册
│   ├── /fleet/heartbeat    # 心跳接收
│   ├── /fleet/events       # 事件收集
│   ├── /fleet/command      # 指令下发
│   └── /fleet/dashboard    # 状态面板
│
├── 数据库表：
│   ├── edge_nodes          # 节点信息
│   ├── edge_events         # 事件历史
│   ├── edge_commands       # 指令历史
│   └── edge_groups         # 节点分组
```

### Phase 3：区域网关 + 离线自治（4-6 周）

**目标**：支持 L1-L2-L3 三层架构，断网时区域内自治。

```
新增：
├── pkg/edge/gateway.go    # L2 区域网关逻辑
│   ├── 管理下辖 L1 节点
│   ├── 区域内事件聚合
│   ├── 离线事件缓存 + 恢复上传
│   └── 跨节点关联分析
│
├── pkg/edge/discovery.go  # 局域网节点自动发现（mDNS/UDP广播）
└── pkg/edge/sync.go       # 状态同步（网络恢复后）
```

### Phase 4：智能调度 + 自学习（持续迭代）

**目标**：云端 AI 基于历史数据优化边缘策略。

```
- 云端分析告警模式，自动调整边缘阈值
- 故障预测（"每周三下午机房温度偏高" → 提前预冷）
- 跨场景知识迁移（A 机房的经验应用到 B 机房）
- 自动生成新的 Skills 下发到边缘节点
```

---

## 八、具体代码改动（Phase 1 详细设计）

### 8.1 新增 Edge Server

```go
// pkg/edge/server.go
package edge

import (
    "encoding/json"
    "net/http"
    "github.com/sipeed/picoclaw/pkg/bus"
)

type Server struct {
    config    EdgeConfig
    bus       *bus.MessageBus
    mux       *http.ServeMux
}

func NewServer(cfg EdgeConfig, b *bus.MessageBus) *Server {
    s := &Server{config: cfg, bus: b}
    s.mux = http.NewServeMux()
    s.mux.HandleFunc("/api/health", s.handleHealth)
    s.mux.HandleFunc("/api/status", s.authMiddleware(s.handleStatus))
    s.mux.HandleFunc("/api/command", s.authMiddleware(s.handleCommand))
    s.mux.HandleFunc("/api/message", s.authMiddleware(s.handleMessage))
    return s
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", 405)
        return
    }

    var cmd CommandRequest
    json.NewDecoder(r.Body).Decode(&cmd)

    // 注入 MessageBus，让 Agent 处理
    s.bus.PublishInbound(bus.InboundMessage{
        Channel:  "edge-api",
        SenderID: "cloud:" + cmd.CommandID,
        ChatID:   "cloud",
        Content:  cmd.Payload.Action,
        Metadata: map[string]string{
            "command_id": cmd.CommandID,
            "priority":   cmd.Priority,
        },
    })

    json.NewEncoder(w).Encode(map[string]string{
        "command_id": cmd.CommandID,
        "status":     "accepted",
    })
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "Bearer "+s.config.APIToken {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next(w, r)
    }
}
```

### 8.2 新增 Edge Reporter

```go
// pkg/edge/reporter.go
package edge

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "time"
)

type Reporter struct {
    cloudURL   string
    apiKey     string
    nodeID     string
    cacheDir   string
    client     *http.Client
}

func (r *Reporter) StartHeartbeat(ctx context.Context) {
    ticker := time.NewTicker(time.Duration(r.interval) * time.Second)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.sendHeartbeat()
        }
    }
}

func (r *Reporter) ReportEvent(event Event) error {
    data, _ := json.Marshal(event)

    resp, err := r.client.Post(
        r.cloudURL+"/fleet/events",
        "application/json",
        bytes.NewReader(data),
    )
    if err != nil {
        // 离线缓存
        return r.cacheEvent(event)
    }
    defer resp.Body.Close()
    return nil
}

func (r *Reporter) cacheEvent(event Event) error {
    // 写入本地文件，网络恢复后批量上传
    filename := fmt.Sprintf("%s/%s.json", r.cacheDir, event.ID)
    data, _ := json.Marshal(event)
    return os.WriteFile(filename, data, 0644)
}
```

### 8.3 新增 Report Tool（让 Agent 主动上报）

```go
// pkg/edge/tool_report.go
package edge

// 注册为 Agent 工具，让 LLM 决定何时上报

var ReportToolDef = tools.ToolDefinition{
    Name:        "report_event",
    Description: "上报事件到云端监控中心。当检测到需要云端关注的异常时使用。",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "severity": {
                "type": "string",
                "enum": ["info", "warning", "critical", "emergency"],
                "description": "事件严重程度",
            },
            "message": {
                "type": "string",
                "description": "事件描述",
            },
            "sensor_data": {
                "type": "object",
                "description": "传感器原始数据",
            },
            "action_taken": {
                "type": "string",
                "description": "本地已采取的措施",
            },
        },
        "required": ["severity", "message"],
    },
}
```

---

## 九、与现有代码的集成点

### 9.1 复用 MaixCAM 通道模式

MaixCAM 已实现的 TCP + JSON 协议是一个完美的模板：

```
现有 MaixCAM 模式:                    新增 Edge API 模式:
─────────────────                    ─────────────────
MaixCAM → TCP → PicoClaw             OpenClaw → HTTP → PicoClaw
         JSON 消息                              JSON 消息
         person_detected              指令: exec/config/query
         heartbeat                    下行: command/message
         status                       上行: heartbeat/event
```

两者可以并存，MaixCAM 走 TCP（本地低延迟），云端走 HTTP（互联网穿透）。

### 9.2 复用 Spawn 子代理

当云端下发复杂任务时，PicoClaw 可以用 `spawn` 创建子代理并行处理：

```
云端指令: "检查所有服务器的 SSL 证书有效期"

PicoClaw 主代理:
├── spawn("检查 server-01 SSL 证书")
├── spawn("检查 server-02 SSL 证书")
├── spawn("检查 server-03 SSL 证书")
└── 汇总结果 → report_event 上报云端
```

### 9.3 复用 Cron 实现定时上报

```
Agent 可以自己创建 cron 任务:

cron({
  "label": "hourly_status_report",
  "every_seconds": 3600,
  "message": "执行完整系统巡检并上报云端",
  "deliver": false  // 通过 agent 处理
})
```

### 9.4 复用 Memory 实现边缘知识积累

```
workspace/memory/
├── MEMORY.md           # 长期知识："机房空调周三下午容易故障"
├── 202602/
│   ├── 20260210.md     # "15:30 温度告警，开启备用空调后 10 分钟恢复"
│   ├── 20260211.md     # "正常巡检，无异常"
│   └── 20260212.md     # "15:30 再次温度告警，同 2月10日模式"
```

Agent 基于历史 Memory 做出更好的判断：
> "根据 2月10日的记录，同样是 15:30 温度异常，当时开启备用空调后 10 分钟恢复。这次直接开启备用空调并通知物业检查主空调。"

---

## 十、成本模型

### 单个监控场景的完整成本

| 组件 | 硬件成本 | 月度成本 | 备注 |
|------|---------|---------|------|
| PicoClaw × 4 节点 | $40 (4× LicheeRV-Nano) | $0 | 一次性 |
| MaixCAM × 2 摄像头 | $100 | $0 | 一次性 |
| NanoKVM × 1 远程管理 | $50 | $0 | 一次性 |
| 传感器套件 | $30 | $0 | 一次性 |
| OpenClaw 云服务器 | — | $20-50 | 轻量云主机 |
| LLM API 调用 | — | $10-30 | 按量付费 |
| **合计** | **$220** | **$30-80/月** | — |

### 对比人力成本

| 方案 | 初始投入 | 月度成本 | 年度总成本 |
|------|---------|---------|-----------|
| 人力 (1 名值班员 × 7×24) | $0 | $4,000-6,000 | $48,000-72,000 |
| 人力 (3 班倒 × 3 人) | $0 | $12,000-18,000 | $144,000-216,000 |
| **Claw Agent 方案** | **$220** | **$30-80** | **$580-1,180** |
| **节省** | — | — | **99%+** |

---

## 十一、总结

### 核心理念

```
不是"用 AI 替代人"，而是——

"让 $10 的硬件 + AI Agent 组成一支永不疲倦的数字员工团队，
 每个节点是一个专注的哨兵，
 云端是运筹帷幄的指挥中心，
 人类只需要做最终决策。"
```

### 一句话架构

**PicoClaw 在边缘感知和执行（<1 秒响应），OpenClaw 在云端思考和调度（跨全局决策），两者通过事件上报 + 指令下发形成闭环。断网时边缘自治，联网时云端增强。**
