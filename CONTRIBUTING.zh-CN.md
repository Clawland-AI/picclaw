# 向 Clawland 贡献

感谢你对 Clawland 的关注！本指南适用于 `github.com/Clawland-AI` 组织下的所有仓库。

---

## 快速入门

1. Fork 仓库
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 进行修改
4. 运行测试（如适用）
5. 提交清晰的提交信息：`git commit -m "feat: add temperature sensor support for SHT40"`
6. 推送并打开 Pull Request
7. 签署 CLA（首次 PR 时自动化机器人会提示你）

---

## 贡献方式

### 代码

- 修复 Bug（查看 [`bug`](https://github.com/search?q=org%3AClawland-AI+label%3Abug&type=issues) 标签）
- 实现功能（查看 [`enhancement`](https://github.com/search?q=org%3AClawland-AI+label%3Aenhancement&type=issues) 标签）
- 新手贡献者：查看 [`good first issue`](https://github.com/search?q=org%3AClawland-AI+label%3A%22good+first+issue%22&type=issues)

### 技能（Skills）

- 为 [clawland-skills](https://github.com/Clawland-AI/clawland-skills) 创建新技能
- 每个技能是一个包含 `SKILL.md` 文件（YAML 前置元数据 + Markdown 正文）的目录
- 参考现有技能的示例

### 硬件套件设计

- 为 [clawland-kits](https://github.com/Clawland-AI/clawland-kits) 设计新的传感器套件
- 包含：BOM、接线图、驱动脚本、PicClaw 技能、成本分析
- 参考现有套件的预期格式

### 文档

- 改进 README 文件，添加教程，翻译文档
- 修正拼写错误，澄清令人困惑的章节

### 社区支持

- 在 GitHub Discussions 和 Discord 中回答问题
- 协助 Issue 分类

### 赏金

- 查看 [赏金看板](https://github.com/orgs/Clawland-AI/projects/1) 中的付费任务
- 通过在 Issue 中评论认领赏金
- 提交引用赏金 Issue 的 PR 作为你的工作

---

## 贡献积分

每项贡献都会获得积分，积分决定你对季度贡献者收入池的分配。完整积分表和分配公式请参见 [CONTRIBUTOR-REVENUE-SHARE.md](CONTRIBUTOR-REVENUE-SHARE.md)。

---

## 提交约定

我们遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

```
type(scope): description

feat(agent): add sub-agent timeout configuration
fix(channels): handle Telegram reconnection on network loss
docs(skills): add weather skill tutorial
chore(ci): update Go version to 1.24
test(tools): add exec tool safety guard tests
```

类型：`feat`、`fix`、`docs`、`chore`、`test`、`refactor`、`perf`、`style`、`ci`

---

## Pull Request 指南

1. **每个功能/修复一个 PR** —— 保持 PR 专注且易于审查
2. **描述** —— 解释是什么、为什么、怎么做。引用相关 Issues
3. **测试** —— 为新功能添加测试，验证现有测试通过
4. **文档** —— 如果行为发生变化，更新相关文档
5. **大小** —— 目标为少于 400 行更改。将较大的更改拆分为一系列 PR

### PR 模板

每个仓库都有 PR 模板。请完整填写。

---

## 代码审查

- 所有 PR 都需要至少一名 Core Maintainer 的批准性审查
- 审核者应在 72 小时内回应
- 审查意见应建设性且具体
- 满意时批准；要求修改时给出清晰解释

---

## 开发环境搭建

### Picclaw（Go）

```bash
git clone https://github.com/Clawland-AI/picclaw.git
cd picclaw
make build
./picclaw status
```

### Moltclaw（TypeScript）

```bash
git clone https://github.com/Clawland-AI/moltclaw.git
cd moltclaw
npm install
npm run dev
```

### 技能

```bash
# 安装技能进行测试
picclaw skills install /path/to/your-skill

# 测试它
picclaw agent -m "test your skill functionality"
```

---

## 报告问题

- 使用每个仓库提供的 Issue 模板
- 包括：重现步骤、预期行为、实际行为、环境详情
- 对于安全漏洞，**请不要公开 Issue** — 发邮件至 security@clawland.dev

---

## 行为准则

所有贡献者必须遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。友善、建设性、热情欢迎。

---

## 疑问？

- [GitHub Discussions](https://github.com/orgs/Clawland-AI/discussions) — 用于提问和建议
- [Discord](https://discord.gg/clawland-ai) — 用于实时聊天

---

*感谢你帮助构建边缘 AI 的未来！*
