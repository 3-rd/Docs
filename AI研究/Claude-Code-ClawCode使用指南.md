# OpenCode 使用指南

> **说明**：本目录中的 `claude-code` 是 [instructkr/claw-code](https://github.com/instructkr/claw-code) 的本地克隆 —— 一个对 Anthropic Claude Code Agent Harness 的 Clean-Room Python 重写。原始 TypeScript 代码由 Anthropic 泄漏，claw-code 团队从架构层面进行了 Python 移植，不包含原始专有代码。

---

## 1. 什么是 Claude Code

Claude Code 是 Anthropic 官方推出的 CLI 编程工具，允许用户通过终端与 Claude（Claude 3.5 Sonnet 等）进行交互式协作，完成代码编写、调试、重构、代码审查等开发任务。

**核心特点：**
- 终端内直接对话式编程
- 强大的代码理解与修改能力（读写文件、执行命令）
- 支持多文件项目分析
- 内置安全审计和权限控制
- 支持 Git 操作、PR 创建
- 可扩展的 Plugin 和 Skill 系统

---

## 2. 安装与启动

### 2.1 安装（官方 Claude Code）

```bash
# macOS/Linux
curl -s https://artifacts.claudecode.com/install.sh | sh

# 或通过 npm
npm install -g @anthropic-ai/claude-code
```

### 2.2 启动方式

```bash
# 基础启动（当前目录）
claude

# 指定项目目录
claude /path/to/project

# 指定模型
claude --model claude-3-5-sonnet

# 禁用网络搜索
claude --no-web-search

# 查看帮助
claude --help
```

### 2.3 本地 claw-code 运行

```bash
cd ~/AIAgent/claude-code
python3 -m src.main --help
```

---

## 3. 核心概念

### 3.1 命令（Commands）

Claude Code 拥有 **207 个内置命令**，覆盖从代码编辑到系统配置的方方面面。每个命令以 `/` 开头触发。

### 3.2 工具（Tools）

Claude Code 通过 **184 个工具** 赋予 Agent 实际操作能力，最核心的工具包括：

| 工具 | 作用 |
|------|------|
| `BashTool` | 执行 Shell 命令 |
| `FileReadTool` | 读取文件内容 |
| `FileEditTool` | 编辑文件（精确替换） |
| `FileWriteTool` | 写入/新建文件 |
| `GlobTool` | 文件模式匹配（glob） |
| `GrepTool` | 文本内容搜索 |
| `AgentTool` | 启动子 Agent |
| `WebSearchTool` | 网络搜索 |
| `WebFetchTool` | 抓取网页内容 |
| `TaskCreateTool` | 创建任务 |
| `MCPTool` | 调用 MCP（Model Context Protocol）服务 |
| `LSPTool` | 调用 Language Server Protocol |

### 3.3 工作模式

- **普通模式（Default）**：常规交互式对话
- **Plan 模式**：`/plan` 进入计划模式，先规划再执行
- **Fast 模式**：`/fast` 快速模式，减少确认直接执行
- **Resume 模式**：`/resume` 恢复被中断的会话

### 3.4 权限与安全

Claude Code 有完善的权限体系：
- **允许列表（Allowlist）**：只允许特定目录下的操作
- **只读模式（Readonly）**：禁止修改文件系统
- **沙箱模式（Sandbox）**：隔离执行危险操作
- **Deny Tool**：可禁用特定工具

---

## 4. 命令详解（按功能分类）

### 4.1 文件与代码操作

| 命令 | 说明 |
|------|------|
| `/add-dir` | 添加目录到工作上下文 |
| `/diff` | 显示当前改动与 HEAD 的差异 |
| `/files` | 列出项目中的文件 |
| `/read` | 读取文件（隐式，通过工具） |
| `/edit` | 编辑文件（隐式，通过工具） |
| `/glob` | 按模式搜索文件 |
| `/grep` | 在文件中搜索文本 |
| `/copy` | 复制文本到剪贴板 |
| `/init` | 初始化新项目 |
| `/lsp` | 调用 Language Server Protocol |

### 4.2 Git 操作

| 命令 | 说明 |
|------|------|
| `/commit` | 创建 Git 提交 |
| `/branch` | 创建/切换分支 |
| `/diff` | 查看文件差异 |
| `/rewind` | 回退到指定提交 |
| `/tag` | 管理 Git 标签 |

### 4.3 PR 与代码审查

| 命令 | 说明 |
|------|------|
| `/review` | 对代码进行审查 |
| `/ultrareview` | 深度代码审查（支持超额用量警告） |
| `/commit-push-pr` | 提交并创建 PR |
| `/autofix-pr` | 自动修复 PR 中的问题 |
| `/pr_comments` | 管理 PR 评论 |

### 4.4 项目构建与执行

| 命令 | 说明 |
|------|------|
| `/build` | 执行构建命令（需项目配置） |
| `/test` | 运行测试 |
| `/debug-tool-call` | 调试工具调用 |
| `/security-review` | 安全审查 |
| `/bughunter` | Bug 追踪与修复 |

### 4.5 Agent 与多任务

| 命令 | 说明 |
|------|------|
| `/agents` | 管理多个 Agent 并行工作 |
| `/thinkback` | 回溯思考过程 |
| `/thinkback-play` | 重放思考过程 |
| `/tasks` | 任务管理（创建/更新/列表） |
| `/passes` | 多轮 Pass 模式 |

### 4.6 记忆与上下文

| 命令 | 说明 |
|------|------|
| `/memory` | 访问和管理持久记忆 |
| `/compact` | 压缩对话上下文 |
| `/context` | 管理上下文内容 |
| `/clear` | 清除上下文 |
| `/clear caches` | 清除缓存 |
| `/clear conversation` | 清除对话历史 |

### 4.7 会话管理

| 命令 | 说明 |
|------|------|
| `/session` | 查看当前会话信息 |
| `/sessions` | 列出所有会话 |
| `/resume` | 恢复之前的会话 |
| `/rename` | 重命名会话 |
| `/share` | 分享会话记录 |
| `/export` | 导出会话记录 |
| `/rewind` | 回退会话到之前的状态 |

### 4.8 配置与设置

| 命令 | 说明 |
|------|------|
| `/model` | 选择/切换模型 |
| `/config` | 编辑配置文件 |
| `/env` | 查看环境变量 |
| `/ide` | 配置 IDE 集成 |
| `/theme` | 切换主题 |
| `/color` | 配置颜色方案 |
| `/output-style` | 配置输出样式 |
| `/rate-limit-options` | 配置速率限制 |

### 4.9 工具链集成

| 命令 | 说明 |
|------|------|
| `/mcp` | 管理 MCP（Model Context Protocol）服务器 |
| `/skills` | 管理 Skills（按需加载的技能包） |
| `/plugins` | 管理插件市场与已安装插件 |
| `/hooks` | 配置 Git hooks |

### 4.10 远程与协作

| 命令 | 说明 |
|------|------|
| `/remote-setup` | 配置远程连接 |
| `/remote-env` | 管理远程环境变量 |
| `/bridge` | 桥接远程会话 |
| `/teleport` | 远程跳转执行 |
| `/install-github-app` | 安装 GitHub App 集成 |

### 4.11 计划与审查

| 命令 | 说明 |
|------|------|
| `/plan` | 进入计划模式（先规划后执行） |
| `/ultraplan` | 超强计划模式 |
| `/brief` | 生成项目简报 |
| `/effort` | 评估任务工作量 |

### 4.12 信息查询

| 命令 | 说明 |
|------|------|
| `/status` | 查看当前状态 |
| `/cost` | 查看 Token 消耗 |
| `/stats` | 查看使用统计 |
| `/doctor` | 诊断配置问题 |
| `/help` | 显示帮助信息 |
| `/keybindings` | 查看快捷键绑定 |
| `/release-notes` | 查看版本更新说明 |
| `/insights` | 查看使用洞察 |

### 4.13 账户与认证

| 命令 | 说明 |
|------|------|
| `/login` | 登录 Anthropic 账号 |
| `/logout` | 登出账号 |
| `/upgrade` | 升级订阅计划 |
| `/extra-usage` | 查看额外用量 |
| `/reset-limits` | 重置用量限制 |

### 4.14 其他命令

| 命令 | 说明 |
|------|------|
| `/btw` | 顺便说一句（插入侧注） |
| `/exit` | 退出 Claude Code |
| `/feedback` | 发送反馈 |
| `/privacy-settings` | 隐私设置 |
| `/sticker` | 发送贴纸 |

---

## 5. 工具详解（按功能分类）

### 5.1 文件操作工具

| 工具 | 说明 |
|------|------|
| `BashTool` | 执行 Shell 命令，支持路径验证、只读验证、危险命令警告 |
| `FileReadTool` | 读取文件，支持图片处理、大文件限制 |
| `FileEditTool` | 精确编辑文件（基于 diff/patch 机制） |
| `FileWriteTool` | 写入或新建文件 |
| `GlobTool` | Glob 模式文件匹配 |
| `GrepTool` | 文本搜索 |
| `NotebookEditTool` | Jupyter Notebook 编辑 |

### 5.2 Agent 与子任务

| 工具 | 说明 |
|------|------|
| `AgentTool` | 启动子 Agent（内置 explore/plan/verify/guide 等类型） |
| `TaskCreateTool` | 创建新任务 |
| `TaskListTool` | 列出所有任务 |
| `TaskGetTool` | 获取任务详情 |
| `TaskUpdateTool` | 更新任务状态 |
| `TaskStopTool` | 停止任务 |
| `TaskOutputTool` | 获取任务输出 |
| `TeamCreateTool` | 创建多 Agent 团队 |
| `TeamDeleteTool` | 删除 Agent 团队 |

### 5.3 计划与模式切换

| 工具 | 说明 |
|------|------|
| `EnterPlanModeTool` | 进入计划模式 |
| `ExitPlanModeV2Tool` | 退出计划模式 |
| `EnterWorktreeTool` | 进入 Git Worktree |
| `ExitWorktreeTool` | 退出 Git Worktree |
| `BriefTool` | 生成项目简报（含附件上传） |

### 5.4 网络工具

| 工具 | 说明 |
|------|------|
| `WebSearchTool` | 网络搜索 |
| `WebFetchTool` | 抓取网页内容（含白名单机制） |

### 5.5 MCP 与扩展

| 工具 | 说明 |
|------|------|
| `MCPTool` | 调用 MCP 服务器工具 |
| `ListMcpResourcesTool` | 列出 MCP 资源 |
| `ReadMcpResourceTool` | 读取 MCP 资源 |
| `McpAuthTool` | MCP 认证 |

### 5.6 定时与消息

| 工具 | 说明 |
|------|------|
| `CronCreateTool` | 创建定时任务 |
| `CronListTool` | 列出定时任务 |
| `CronDeleteTool` | 删除定时任务 |
| `SendMessageTool` | 发送消息 |

### 5.7 其他工具

| 工具 | 说明 |
|------|------|
| `AskUserQuestionTool` | 向用户提问 |
| `ConfigTool` | 读写配置项 |
| `TodoWriteTool` | 写 TODO 项 |
| `ToolSearchTool` | 搜索可用工具 |
| `SyntheticOutputTool` | 合成输出 |
| `SleepTool` | 延迟执行 |
| `RemoteTriggerTool` | 远程触发 |
| `PowerShellTool` | PowerShell 命令执行 |

---

## 6. 常用场景与操作指南

### 6.1 日常代码开发

```bash
# 启动项目
claude

# 让 Claude 读取并理解代码
> 帮我理解这个项目的架构

# 编写新功能
> 在 user service 中添加一个获取用户信息的接口

# 编辑现有文件
> 将 auth.py 中的 JWT 替换为 PASETO

# 搜索代码
> 找出所有使用了 mysql_query 的地方
```

### 6.2 Git 工作流

```bash
# 查看改动
/diff

# 提交代码
/commit

# 创建分支并提交
/branch feature/new-api

# 创建 PR
/commit-push-pr

# 代码审查
/review
```

### 6.3 Bug 修复

```bash
# 让 Claude 分析 Bug
> 这个接口返回 500 错误，帮我分析原因

# 进入计划模式逐步排查
/plan

# 安全审查
/security-review

# 自动修复
/autofix-pr
```

### 6.4 代码重构

```bash
# 启动重构计划
/plan

# 指定重构范围
> 将所有 Controller 从 MVC 模式迁移到 DDD 模式

# 查看变更差异
/diff

# 提交
/commit -m "refactor: migrate controllers to DDD"
```

### 6.5 项目初始化

```bash
# 初始化新项目
/init

# 配置项目
> 使用 TypeScript + Express + Prisma 初始化

# 初始化 Git hooks
/hooks setup

# 安装需要的插件
/plugin install eslint
```

### 6.6 多文件分析

```bash
# 让 Claude 扫描整个项目
> 扫描这个微服务项目，列出所有 API 端点

# 分析依赖关系
> 画出 user-service 的调用链路图

# 找特定模式
> 找出所有可能的内存泄漏点（未关闭的流）
```

### 6.7 使用 MCP 扩展

```bash
# 查看已配置的 MCP 服务器
/mcp

# 添加新的 MCP 服务器
> 添加一个 GitHub MCP 服务器

# 使用 MCP 工具
> 用 GitHub MCP 创建一个 issue
```

### 6.8 使用 Skills

```bash
# 查看可用 Skills
/skills

# 安装 Skill
> 安装 react-development skill

# 使用 Skill
> 使用 react skill 帮我写一个组件
```

### 6.9 会话管理

```bash
# 查看当前会话
/session

# 列出所有历史会话
/sessions

# 恢复之前会话
/resume <session-id>

# 重命名会话
/rename <new-name>

# 压缩上下文（省钱）
/compact
```

### 6.10 远程协作

```bash
# 配置远程环境
/remote-setup

# SSH 模式运行
claude --ssh user@server

# 远程触发
/teleport production
```

---

## 7. claw-code CLI 本地命令

claw-code（Python 重写版）提供以下本地 CLI 命令：

```bash
# 查看工作区摘要
python3 -m src.main summary

# 查看项目清单
python3 -m src.main manifest

# 列出所有命令条目
python3 -m src.main commands --limit 200

# 搜索命令
python3 -m src.main commands --query "commit"

# 列出所有工具条目
python3 -m src.main tools --limit 200

# 搜索工具
python3 -m src.main tools --query "bash"

# 查看单个命令详情
python3 -m src.main show-command <name>

# 查看单个工具详情
python3 -m src.main show-tool <name>

# 路由提示词到命令/工具
python3 -m src.main route "帮我提交代码"

# 运行对话循环
python3 -m src.main turn-loop "帮我写一个排序算法" --max-turns 5

# 持久化会话
python3 -m src.main flush-transcript "用户的请求"

# 加载历史会话
python3 -m src.main load-session <session-id>

# 查看命令图谱
python3 -m src.main command-graph

# 查看工具池
python3 -m src.main tool-pool

# 运行单元测试
python3 -m unittest discover -s tests -v
```

---

## 8. 配置参考

### 8.1 主要配置项（通过 `/config` 或 `~/.claude.json`）

```json
{
  "model": "claude-3-5-sonnet-20241022",
  "maxTokens": 8192,
  "temperature": 1,
  "tools": {
    "bash": {
      "allowlist": ["/Users/fucy/projects/**"],
      "readOnly": false
    }
  },
  "mcpServers": {},
  "plugins": [],
  "theme": "dark"
}
```

### 8.2 环境变量

```bash
# API Key
ANTHROPIC_API_KEY=sk-...

# API Base URL（代理）
ANTHROPIC_BASE_URL=https://api.anthropic.com

# 本地模型
ANTHROPIC_BASE_URL=http://localhost:11434
```

---

## 9. Plugin 系统

Claude Code 支持插件扩展：

```bash
# 浏览插件市场
/plugin browse

# 安装插件
/plugin install <plugin-name>

# 查看已安装插件
/plugin list

# 卸载插件
/plugin uninstall <plugin-name>

# 管理插件市场源
/plugin marketplaces
```

---

## 10. 与 OpenClaw 的对比

| 特性 | Claude Code (claw-code) | OpenClaw |
|------|------------------------|----------|
| 定位 | 单用户编程 Agent | 多平台 AI 助手 + Agent 编排 |
| 交互方式 | 终端 REPL | 飞书/微信/Discord 等多平台 |
| 核心能力 | 代码编写/重构/审查 | 消息处理 + 任务自动化 |
| MCP 支持 | ✅ | ✅ |
| Plugin 支持 | ✅ | ✅（飞书等平台插件） |
| 多 Agent | ✅（Team 工具） | ✅ |
| 定时任务 | ✅（Cron 工具） | ✅（Cron 系统） |
| 适用场景 | 日常编程、代码审查 | 跨平台消息、办公自动化 |

---

## 附录：命令速查表

### 高频命令 Top 20

| 排名 | 命令 | 用途 |
|------|------|------|
| 1 | `/help` | 获取帮助 |
| 2 | `/plan` | 进入计划模式 |
| 3 | `/commit` | Git 提交 |
| 4 | `/diff` | 查看变更 |
| 5 | `/review` | 代码审查 |
| 6 | `/session` | 会话信息 |
| 7 | `/compact` | 压缩上下文 |
| 8 | `/model` | 切换模型 |
| 9 | `/config` | 编辑配置 |
| 10 | `/agent` | 启动 Agent |
| 11 | `/mcp` | 管理 MCP |
| 12 | `/skills` | 管理 Skills |
| 13 | `/plugin` | 管理插件 |
| 14 | `/clear` | 清除上下文 |
| 15 | `/resume` | 恢复会话 |
| 16 | `/doctor` | 诊断问题 |
| 17 | `/cost` | 查看花费 |
| 18 | `/init` | 初始化项目 |
| 19 | `/btw` | 插入侧注 |
| 20 | `/exit` | 退出 |
