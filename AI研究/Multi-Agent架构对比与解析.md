# Multi-Agent 架构对比与解析

> **来源**：2026-04-17 与付晨阳的讨论
> **主题**：Clowder AI 多猫协作机制、OpenClaw subagent 架构、无状态 vs 有状态 Agent 设计

---

## 一、Clowder AI 多猫协作机制

### 1.1 架构概述

Clowder AI 是一个多 Agent 协作平台，核心定位是**平台层**（不是模型层），在现有 Agent CLI 之上加了一层协作编排。

```
Node.js 进程（Clowder 平台）
├── 协调层（路由、A2A、Evidence Store）
├── Skills 框架
└── 子进程（独立 CLI）
    ├── 布偶猫 #1（Claude CLI）→ 独立进程
    ├── 布偶猫 #2（Claude CLI）→ 独立进程
    ├── 缅因猫 #1（Codex CLI）→ 独立进程
    └── 暹罗猫 #1（Gemini CLI）→ 独立进程
```

**关键发现**：
- Clowder 平台本身是**单 Node.js 进程**
- 每只"猫"是**独立的 CLI 子进程**（通过 `spawnCli()` 启动）
- 这与 OpenClaw 的 session 机制有本质区别

### 1.2 品种 vs Variant（品种 vs 变种）

**品种（Breed）** = 不同种类的猫 = 不同 Agent CLI = 不同专业分工

| 品种 | CLI | Provider | 团队角色 |
|------|-----|---------|---------|
| 布偶猫（Ragdoll） | claude | Anthropic | 主架构师、深度思考 |
| 缅因猫（Maine Coon） | codex | OpenAI | 代码审查、bug 定位 |
| 暹罗猫（Siamese） | gemini | Google | UI/UX 设计、创意 |
| 孟加拉猫（Bengal） | antigravity | 多模型 | 浏览器自动化、截图 |
| 狸花猫（Dragon Li） | dare | dare | 确定性执行、零信任验证 |
| 金渐层（Golden Chinchilla） | opencode | 多模型 | 全能型、多专家编排 |

**Variant（变种）** = 同品种的不同版本 = 同一 CLI + 不同模型档位

```
布偶猫-Opus     → claude-opus-4-6（最强、最贵）
布偶猫-Sonnet   → claude-sonnet-4-6（中等）
布偶猫-Opus4.5  → claude-opus-4-5（稍弱）
缅因猫-Spark    → gpt-5.3-codex-spark（最快、context 最小）
```

### 1.3 区分品种的目的

**品种区分 = 能力分工**（不是随机分配）

- 架构猫（Claude）→ 深度思考、系统设计
- 开发猫（opencode）→ 代码落地、多专家编排
- 审查猫（Codex）→ 找 bug、安全分析
- 设计猫（Gemini）→ UI/UX、创意表达

每个品种调用**不同的 Agent CLI**，背后是不同公司的模型，各有擅长领域。

**Variant 区分 = 经济考量 + 上下文大小**

- Opus 180k tokens vs Spark 64k tokens
- 能用 Sonnet 搞定的简单任务，不浪费 Opus 的算力

### 1.4 运行时隔离与共享机制

| 资源 | 是否共享 |
|------|---------|
| CLI 工具本身（二进制） | ✅ 同一份 |
| 运行时进程/上下文 | ❌ 每只猫独立 |
| Evidence Store（知识库） | ✅ 显式设计共享 |
| Skills 框架 | ✅ 共享资源池 |
| MCP 工具 | ✅ 共享（部分品种不支持） |
| 线程上下文 | ❌ 完全隔离 |

**Evidence Store** 是所有猫共享的 SQLite 数据库，弥补了 CLI 无状态的问题。每次任务的结果、决策、教训都沉淀到这里，供其他猫查询。

**A2A 通信**：猫猫之间**不直接读对方的上下文**，而是通过结构化消息传递：
```
猫猫 A：完成架构设计 → 写入 Evidence Store，通知 @codex
猫猫 B：收到通知 → 主动拉取 Evidence Store → review 后写回结果
```

**跨品种 Review 规则**（`cat-config.json`）：
```json
"reviewPolicy": {
  "requireDifferentFamily": true  // 审查必须跨品种，布偶猫写的代码必须让缅因猫来审
}
```

---

## 二、CLI 无状态 vs OpenClaw 有状态

### 2.1 设计哲学的根本差异

| | Claude CLI / opencode / Codex | OpenClaw |
|--|------------------------------|---------|
| **设计哲学** | 无状态，用完即走 | 有状态，长期陪伴 |
| **启动方式** | 每次全新进程 | 长期驻留进程 |
| **上下文** | 每次运行从 0 开始 | session 内持久 |
| **记忆** | 不保留（需要自己搭 Evidence Store） | 内置 memory 系统 |
| **多实例** | 原生支持（直接起多个） | 不支持（单进程单端口） |
| **进程占用** | 无端口占用 | 占端口（3000/3001） |

### 2.2 无状态 CLI 的实际行为

```bash
$ claude
# 会话内：上下文保留
> 帮我实现登录模块
> 帮我实现支付模块
# 关掉

# 第二天，重新启动
$ claude
> 继续完善订单模块
# Claude：不记得之前做了什么（除非你手动提示）
```

**会话内的上下文是保留的，但会话结束就丢了。**

### 2.3 有状态 OpenClaw 的优势

```
每日 memory → 自动存档到 memory/YYYY-MM-DD.md
长期 memory → 自动更新 MEMORY.md
Skills → 持续加载
```

OpenClaw 把"记忆"做成了**平台内置功能**，不需要自己搭记忆层。

### 2.4 为什么无状态的 CLI 能起多个进程

因为每次启动都是**全新的独立进程**，不共享任何内存：
```
终端 A → claude 进程 #1（项目 A）
终端 B → claude 进程 #2（项目 B）
终端 C → claude 进程 #3（项目 C）
```

而 OpenClaw 是长期服务，一旦启动就占端口，再启动第二个会报端口冲突。

---

## 三、OpenClaw Subagent 机制详解

### 3.1 两种 subagent 模式

OpenClaw 支持两种 subagent 模式：

| 模式 | 说明 |
|------|------|
| `runtime="subagent"` | 轻量隔离上下文，共享主进程 tools/skills |
| `runtime="acp"` | 独立 OS 进程，更彻底隔离 |

### 3.2 subagent 模式（轻量）

```
OpenClaw 进程
├── Main Session
│   ├── memory 系统 ✅ 共享
│   ├── skills ✅ 共享
│   └── tools ✅ 共享
└── Subagent Context（隔离执行上下文）
    ├── 独立 prompt/指令
    ├── 独立工具调用栈
    ├── 无法直接读写 main session 的变量
    └── 通过消息传递结果
```

| 维度 | 说明 |
|------|------|
| **进程** | 同一进程，隔离的是执行上下文 |
| **memory** | ❌ 不自动共享 |
| **tools/skills** | ✅ 共享（主进程的） |
| **启动速度** | 快（毫秒级） |
| **隔离程度** | 中（上下文级别） |

### 3.3 ACP 模式（独立进程）

```
OpenClaw 进程                              ACP 进程
├── Main Session                          ├── 独立执行上下文
│   └── tools/skills/memory              │   ├── 独立 workspace
└── ACP 进程（独立 OS 进程）              │   ├── 独立 tools/skills
    ├── 独立 Node.js 运行时              │   └── 通过 sessions_send 通信
    └── 独立 workspace
```

| 维度 | 说明 |
|------|------|
| **进程** | 独立 OS 进程 |
| **memory** | ❌ 完全隔离 |
| **workspace** | ❌ 独立（默认继承主 workspace） |
| **启动速度** | 慢（秒级） |
| **隔离程度** | 高（真正进程级） |
| **通信方式** | sessions_send |

**ACP 不是另一个 OpenClaw Gateway**，而是一个可以被 OpenClaw 远程调度的独立执行器，不占 Gateway 端口。

### 3.4 对比 subagent vs ACP

| | subagent | acp |
|--|---------|-----|
| **隔离级别** | 上下文级 | 进程级 |
| **进程** | 共享主进程 | 独立进程 |
| **memory** | 不共享 | 不共享 |
| **workspace** | 共享 | 可配置 |
| **tools/skills** | 共享主进程 | 独立副本 |
| **通信方式** | 子 agent 消息 | sessions_send |
| **启动速度** | 快 | 慢 |
| **适用场景** | 临时任务、并行思考 | 需要彻底隔离的任务 |

---

## 四、Clowder AI vs OpenClaw 架构对比

### 4.1 核心架构对比

| | Clowder AI | OpenClaw |
|--|-----------|---------|
| **平台定位** | 多 Agent 协作平台 | 个人 AI 陪伴助手 |
| **协调层进程** | 单 Node.js 进程 | 单 Node.js Gateway |
| **多任务执行** | 每只猫 = 独立 CLI 进程 | session 在进程内隔离 |
| **隔离单位** | OS 进程级（每只猫） | 线程/上下文级（session） |
| **记忆机制** | 外置 Evidence Store（SQLite） | 内置 memory 系统（文件） |
| **记忆共享** | 显式通过 Evidence Store | 同一进程内自然共享 |
| **启动开销** | 高（要起完整 CLI） | 低（线程切换） |
| **多实例** | 每只猫天然多实例 | Gateway 单实例 |

### 4.2 设计权衡

**Clowder 的选择**：
- 优点：进程级隔离，一只猫崩了不影响其他；天然支持多猫并行
- 缺点：每次任务要重新初始化上下文（依赖 Evidence Store）；启动慢

**OpenClaw 的选择**：
- 优点：自带 memory 系统，记忆自动持久化；轻量快速
- 缺点：单进程，多 session 无法真正进程级隔离

### 4.3 鱼和熊掌

```
无状态 CLI（Claude/opencode）
├── ✅ 天然支持多进程并行
├── ✅ 进程级隔离
├── ✅ 一只猫崩了不影响其他
└── ❌ 需要自己搭记忆层（Evidence Store）

有状态单进程（OpenClaw）
├── ✅ 自带 memory 系统，开箱即用
├── ✅ 记忆自动持久化
└── ❌ 单进程，多 session 无法真正进程级隔离
```

---

## 五、关键概念区分

### 5.1 多猫 vs 多 session

| | Clowder 多猫 | OpenClaw 多 session |
|--|------------|-------------------|
| **本质** | 不同 CLI 进程 | 同一进程内不同会话 |
| **隔离** | 进程级 | 上下文级 |
| **记忆** | 外置（Evidence Store） | 内置（memory 文件） |
| **启动** | 慢（起进程） | 快（线程切换） |

### 5.2 opencode vs OpenClaw

| | opencode（工具） | OpenClaw（平台） |
|--|----------------|-----------------|
| **性质** | 无状态 CLI 工具 | 有状态服务进程 |
| **多实例** | ✅ 随便开多少 | ❌ 端口唯一 |
| **记忆** | 无（自己想办法） | 内置 memory |
| **端口** | 无 | 占 3000/3001 |

### 5.3 Clowder 平台 vs Clowder 的猫

```
Clowder 平台进程（1个）
├── 协调层（路由、A2A、Evidence Store）
├── 每只猫 = 独立 CLI 进程（不是平台的一部分）
└── Skills 框架
```

Clowder 平台本身也是单进程，和 OpenClaw 一样。区别在于 Clowder 把"多猫"外包给独立 CLI 进程，而 OpenClaw 用 session 在进程内模拟多任务。

---

## 六、总结

1. **Clowder 多猫的本质**：单平台进程 + 多独立 CLI 进程，Evidence Store 是共享记忆层
2. **品种 vs Variant**：品种 = 专业分工（不同 CLI），Variant = 能力档位（同一 CLI 不同模型）
3. **无状态 vs 有状态**：CLI 用完即走不保留记忆，OpenClaw 自动持久化记忆
4. **OpenClaw subagent 两种模式**：subagent 轻量共享主进程资源，ACP 独立进程更彻底隔离
5. **架构选择没有绝对优劣**：Clowder 适合真正的多 Agent 团队协作，OpenClaw 适合个人陪伴场景

---

## 相关文档

- `Cat-Cafe-ClowderAI部署指南-Windows.md` — Clowder AI 部署指南
- `OpenCode使用指南.md` — opencode CLI 使用指南
