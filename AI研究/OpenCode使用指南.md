# OpenCode 使用指南

> **更新时间**：2026-04-16
> **官网**：https://opencode.ai

---

## 1. OpenCode 是什么

OpenCode 是一个**开源、provider-agnostic 的 AI 编程 Agent**，核心特点：

- **MIT 开源协议**，代码完全透明
- **75+ 模型 Provider 支持**（Anthropic / OpenAI / Google / 本地模型等）
- **三合一产品形态**：桌面客户端 + CLI 工具 + AI 模型聚合 API
- **原生 MCP 支持**，内置 MCP Client
- **Oh My OpenCode (OMOC)** 插件生态：多专家编排、自循环、上下文智能管理
- **强 TUI/主题生态**，社区活跃

**产品定位对比**：

| 维度 | OpenCode | Claude Code | Cursor |
|------|----------|-------------|--------|
| 开源 | ✅ MIT | ❌ 闭源 | ❌ 闭源 |
| 桌面客户端 | ✅ | ❌ | ✅ |
| CLI 模式 | ✅ | ✅ | ❌ |
| 多 Provider | ✅ 75+ | ❌ 仅 Anthropic | ❌ 仅 OpenAI |
| MCP 原生 | ✅ | ✅ | ✅ |
| 插件生态 | ✅ OMOC | ✅ Plugin | ✅ |
| 平台 | macOS/Linux/Windows | macOS/Linux | macOS/Windows |

---

## 2. 安装

### 2.1 桌面客户端（推荐）

```bash
brew install --cask opencode-desktop
```

或从 [opencode.ai](https://opencode.ai) 下载安装。

### 2.2 CLI 工具

```bash
# via npm
npm install -g opencode-ai

# via npx（免安装）
npx opencode-ai --version

# via bun
bun install -g opencode-ai
```

### 2.3 验证安装

```bash
opencode --version
```

---

## 3. 桌面客户端使用

### 3.1 界面布局

```
┌─────────────────────────────────────────┐
│  侧边栏    │     主对话区                 │
│  ─────    │  ─────────────────────────   │
│  会话列表  │  AI 回复 + 代码块            │
│  文件浏览  │  工具调用展示                │
│  工具面板  │                              │
└─────────────────────────────────────────┘
```

### 3.2 核心功能区

- **Chat 对话**：主要交互区域，支持多轮对话
- **File Tree**：项目文件浏览
- **Terminal**：内置终端，可直接执行命令
- **MCP Server**：MCP 工具调用面板
- **Theme**：切换 TUI 主题

---

## 4. CLI 命令详解

### 4.1 基础命令

```bash
# 启动交互式对话（当前目录）
opencode

# 指定项目目录
opencode /path/to/project

# 指定模型
opencode --model anthropic/claude-sonnet-4-6

# 查看帮助
opencode --help

# 查看版本
opencode --version
```

### 4.2 运行模式

```bash
# 交互式 TUI 模式（默认）
opencode

# 单次请求模式（非交互）
opencode "帮我写一个快速排序"

# JSON 输出模式（适合脚本集成）
opencode run --format json "分析这个项目的结构"

# Headless 模式（后台运行）
opencode run --headless "审查代码安全"

# 流式输出
opencode "写一个 HTTP 服务器" --stream
```

### 4.3 模型选择

```bash
# 列出可用模型
opencode models list

# 指定 provider + 模型
opencode --provider anthropic --model claude-sonnet-4-6

# 通过 OpenCode API 使用各厂商模型
opencode --model opencode/claude-opus-4-6    # Zen 目录
opencode --model opencode-go/kimi-k2.5        # Go 目录
```

### 4.4 项目操作

```bash
# 初始化新项目
opencode init

# 扫描并理解项目
opencode scan

# 指定工作目录
opencode --cwd /path/to/project

# 分析依赖
opencode deps analyze
```

### 4.5 MCP 相关

```bash
# 查看已配置的 MCP 服务器
opencode mcp list

# 添加 MCP 服务器
opencode mcp add github

# 移除 MCP 服务器
opencode mcp remove github

# 检查 MCP 工具
opencode mcp tools
```

### 4.6 配置管理

```bash
# 打开配置文件
opencode config edit

# 查看当前配置
opencode config show

# 重置配置
opencode config reset
```

### 4.7 会话管理

```bash
# 列出历史会话
opencode sessions list

# 恢复会话
opencode resume <session-id>

# 导出会话
opencode export <session-id>

# 清除会话
opencode sessions clear
```

---

## 5. Oh My OpenCode（OMOC）插件生态

OMOC 是 OpenCode 的杀手级插件系统，类比 VSCode 的插件市场。

### 5.1 安装插件

```bash
# 安装 OMOC
opencode plugin install oh-my-opencode

# 安装社区插件
opencode plugin install sisyphus        # 多专家编排器
opencode plugin install ralph-loop     # 自循环执行器
opencode plugin install context-manager # 上下文管理器
```

### 5.2 OMOC 核心组件

| 组件 | 功能 |
|------|------|
| **Sisyphus** | 多专家编排器，内部管理 Oracle/Librarian/Frontend 等子 Agent |
| **Ralph Loop** | 自循环执行器，支持 Architect 级验证的持久执行循环 |
| **Context Manager** | 70% 预警 / 85% 自动压缩的上下文管理 |

### 5.3 插件配置

在 `opencode.json` 中配置：

```jsonc
{
  "plugin": ["oh-my-opencode"],
  "model": "anthropic/claude-sonnet-4-6",
  "provider": {
    "anthropic": {
      "options": {
        "apiKey": "{env:ANTHROPIC_API_KEY}",
        "baseURL": "https://chat.nuoda.vip/claudecode"  // 可用代理
      }
    }
  }
}
```

---

## 6. OpenCode API：Zen 与 Go 两大目录

OpenCode 提供两个模型目录，通过同一个 API Key 访问：

### 6.1 Zen 目录（`opencode/...`）

| 模型 | 说明 |
|------|------|
| `opencode/claude-opus-4-6` | Claude Opus 4 |
| `opencode/gpt-5.2` | GPT-5.2 |
| `opencode/gemini-3-pro` | Gemini 3 Pro |

### 6.2 Go 目录（`opencode-go/...`）

| 模型 | 说明 |
|------|------|
| `opencode-go/kimi-k2.5` | Kimi K2.5 |
| `opencode-go/glm-5` | GLM-5 |
| `opencode-go/minimax-m2.5` | MiniMax M2.5 |

### 6.3 OpenClaw 接入配置

OpenCode 已接入 OpenClaw，使用方式：

```bash
# Zen 目录接入
openclaw onboard --auth-choice opencode-zen
openclaw onboard --opencode-zen-api-key "$OPENCODE_API_KEY"

# Go 目录接入
openclaw onboard --auth-choice opencode-go
openclaw onboard --opencode-go-api-key "$OPENCODE_API_KEY"
```

配置示例：

```jsonc
{
  "env": { "OPENCODE_API_KEY": "sk-..." },
  "agents": {
    "defaults": {
      "model": { "primary": "opencode-go/kimi-k2.5" }
    }
  }
}
```

---

## 7. 常用场景

### 7.1 代码编写

```bash
# 直接请求
opencode "用 Python 写一个 LRU 缓存"

# 交互式协作
opencode
> 帮我写一个 WebSocket 服务器
> 支持心跳检测
> 添加断线重连
```

### 7.2 代码审查

```bash
# 安全审查
opencode run "审查这段代码的安全性" --file ./auth.py

# 性能审查
opencode run "分析性能瓶颈" --dir ./src
```

### 7.3 项目分析

```bash
# 理解项目结构
opencode scan

# 分析调用链路
opencode "画出 user-service 的调用链路图"
```

### 7.4 Git 操作

```bash
# 生成 commit message
opencode git commit-message

# 代码审查
opencode "审查 main 分支与 feature 分支的差异"
```

### 7.5 自定义 Provider

OpenCode 支持接入任意兼容 OpenAI/Anthropic 格式的 API：

```jsonc
{
  "provider": {
    "custom": {
      "options": {
        "apiKey": "{env:CUSTOM_API_KEY}",
        "baseURL": "https://your-proxy.com/v1"
      }
    }
  }
}
```

---

## 8. 配置文件参考

主配置文件：`~/.opencode/config.json` 或项目根目录 `opencode.json`

```jsonc
{
  "model": "anthropic/claude-sonnet-4-6",
  "provider": {
    "anthropic": {
      "options": {
        "apiKey": "{env:ANTHROPIC_API_KEY}"
      }
    }
  },
  "plugin": [],
  "mcpServers": {},
  "theme": "catppuccin-mocha",
  "context": {
    "warningAt": 0.7,
    "compactAt": 0.85
  }
}
```

---

## 9. Cat Cafe（金渐层）集成

OpenCode 已作为**金渐层（Golden Chinchilla）**接入 Cat Cafe 多猫协作体系：

- **定位**：开源多模型编码猫
- **句柄**：`@opencode`、`@金渐层`、`@golden`
- **颜色**：Primary `#C8A951`（金色）/ Secondary `#F5EDDA`（奶白）
- **独有能力**：OMOC 多专家内部编排 + LSP + 主题生态
- **编排隔离**：Sisyphus 只编排自己的子 Agent，不编排其他猫

---

## 10. 与其他工具对比

| 特性 | OpenCode | Claude Code | Cursor | Copilot |
|------|----------|-------------|--------|---------|
| **开源协议** | MIT | 闭源 | 闭源 | 闭源 |
| **Provider 数量** | 75+ | 仅 Anthropic | 仅 OpenAI | 仅 Azure |
| **桌面客户端** | ✅ | ❌ | ✅ | ❌ |
| **CLI** | ✅ | ✅ | ❌ | ❌ |
| **MCP 支持** | ✅ | ✅ | ✅ | 有限 |
| **OMOC 插件** | ✅ | ❌ | ✅ | ❌ |
| **Ralph Loop** | ✅ | ❌ | ❌ | ❌ |
| **多模型聚合 API** | ✅ | ❌ | ❌ | ❌ |

---

## 附录：命令速查

| 命令 | 说明 |
|------|------|
| `opencode` | 启动交互式对话 |
| `opencode <prompt>` | 单次请求 |
| `opencode run --format json` | JSON 输出模式 |
| `opencode models list` | 列出可用模型 |
| `opencode mcp list` | 查看 MCP 服务器 |
| `opencode plugin install <name>` | 安装插件 |
| `opencode config edit` | 编辑配置 |
| `opencode sessions list` | 列出历史会话 |
| `opencode resume <id>` | 恢复会话 |
| `opencode scan` | 扫描项目结构 |
