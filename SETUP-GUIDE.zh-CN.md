# Clawland GitHub 组织搭建指南

本目录包含搭建 `github.com/Clawland-AI` 组织所需的所有文件。

---

## 步骤 1：创建 GitHub 组织

1. 访问 https://github.com/organizations/plan
2. 创建组织：**Clawland-AI**
3. 添加描述："Edge AI Agents for $10 Hardware. Build to Earn."

---

## 步骤 2：创建 `.github` 仓库

这是用于存放组织简介、模板和治理文件的特殊组织级仓库。

```bash
# 先在 GitHub 上创建仓库，然后：
git init Clawland-AI-.github
cd Clawland-AI-.github

# 复制本目录 .github/ 文件夹中的所有文件：
# .github/
#   profile/README.md        → 组织简介页面
#   GOVERNANCE.md             → TSC 章程和决策流程
#   CONTRIBUTOR-REVENUE-SHARE.md → 收入分成规则
#   CLA.md                    → 贡献者许可协议
#   CODE_OF_CONDUCT.md        → 社区行为准则
#   CONTRIBUTING.md            → 如何贡献
#   ISSUE_TEMPLATE/            → Issue 模板（组织内共享）
#     bug_report.md
#     feature_request.md
#     sensor_kit_proposal.md
#   PULL_REQUEST_TEMPLATE.md   → PR 模板（组织内共享）
#   reports/
#     TEMPLATE.md              → 季度报告模板

git add -A
git commit -m "feat: initialize Clawland organization governance and templates"
git remote add origin https://github.com/Clawland-AI/.github.git
git push -u origin main
```

---

## 步骤 3：转移/复刻 picclaw 仓库

选项 A — 转移现有仓库：
1. 访问 https://github.com/sipeed/picoclaw/settings（或你的 Fork）
2. 滚动到“危险区”→“转移所有权”
3. 转移至 `Clawland-AI` 组织
4. 从 `picoclaw` 重命名为 `picclaw`

选项 B — 复刻并重命名：
1. 将 https://github.com/sipeed/picoclaw 复刻到 `Clawland-AI`
2. 在仓库设置中重命名为 `picclaw`
3. 将 LICENSE 更新为 Apache 2.0（见 `licenses/APACHE-2.0-picclaw.txt`）

---

## 步骤 4：创建其他仓库

在 GitHub 上，在 `Clawland-AI` 组织下创建以下仓库：
- `moltclaw` （TypeScript，Apache 2.0）
- `nanoclaw` （Python，Apache 2.0）
- `microclaw` （C/Rust，Apache 2.0）
- `clawland-fleet` （Go/TypeScript，BSL 1.1）
- `clawland-kits` （硬件设计，CERN-OHL-S-2.0）
- `clawland-skills` （Markdown 技能，MIT）
- `clawland.github.io` （文档站点，CC BY-SA 4.0）

### 许可证文件

使用 `licenses/` 目录中的许可证文件：
- `APACHE-2.0-picclaw.txt` → 用于 picclaw、moltclaw、nanoclaw、microclaw
- `BSL-1.1-fleet.txt` → 用于 clawland-fleet

---

## 步骤 5：配置 CLA 机器人

1. 安装 [CLA Assistant](https://github.com/apps/cla-assistant) 到 Clawland-AI 组织
2. 配置其使用 `Clawland-AI/.github/CLA.md` 作为协议
3. 启用所有仓库

---

## 步骤 6：设置 Open Collective

1. 访问 https://opencollective.com/create
2. 创建集体：**Clawland-AI**
3. 选择财政主机（推荐 Open Source Collective）
4. 链接 GitHub 组织
5. 设置付款方式（PayPal、Wise/TransferWise）

---

## 步骤 7：创建赏金看板

1. 访问 https://github.com/orgs/Clawland-AI/projects
2. 创建项目："赏金看板"
3. 添加列：未认领、已认领、进行中、审查、已付款
4. 在各仓库中向每个列添加初始赏金 Issue，并加上价格标签

---

## 步骤 8：发布公告

使用 `launch/BUILD-TO-EARN.md` 作为模板，发布：
- GitHub Discussions 公告（置顶）
- Hacker News 帖子
- Reddit r/opensource、r/golang、r/homelab、r/selfhosted
- V2EX / 掘金（中文开发者社区）
- Discord 服务器创建

---

## 目录结构参考

```
clawland-org/
├── .github/                          ← 用于 .github 仓库的文件
│   ├── profile/
│   │   └── README.md                 ← 组织简介页面
│   ├── GOVERNANCE.md                 ← TSC 章程
│   ├── CONTRIBUTOR-REVENUE-SHARE.md  ← 收入分成规则
│   ├── CLA.md                        ← 贡献者协议
│   ├── CODE_OF_CONDUCT.md            ← 社区行为准则
│   ├── CONTRIBUTING.md               ← 如何贡献
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   ├── feature_request.md
│   │   └── sensor_kit_proposal.md
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── reports/
│       └── TEMPLATE.md               ← 季度报告模板
├── launch/
│   └── BUILD-TO-EARN.md              ← 发布公告草稿
├── licenses/
│   ├── APACHE-2.0-picclaw.txt        ← 用于核心代理仓库
│   └── BSL-1.1-fleet.txt             ← 用于 clawland-fleet
└── SETUP-GUIDE.md                    ← 本文件
```
