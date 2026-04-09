# Clowder AI 项目核心概念分析

> 项目地址：https://github.com/zts212653/clowder-ai
> 克隆位置：~/AIAgent/clowder-ai/
> 调研时间：2026-04-09
> 定位：多 Agent 团队协作平台（vs CodeMind 是代码理解工具）

---

## 项目概述

**Slogan**："Build AI teams, not just agents. Hard rails, soft power, shared mission."

**核心技术栈**：TypeScript + Node.js 20+ + pnpm + Redis

**规模**：394 stars, 150 forks, 2026-03-12 创建（很新）

---

## 核心概念清单

### 1. Multi-Agent Orchestration（多 Agent 编排）

**定义**：协调多个 AI Agent 协同工作，而非单独使用一个 Agent。

**Clowder 实现**：
- 集成 Claude Code、Codex (GPT)、Gemini CLI、opencode 等多种 Agent
- 每只"猫"有独立身份：XianXian（宪宪/Claude）、YanYan（砚砚/Codex）、ShuoShuo（烁烁/Gemini）
- 通过 @mention 路由任务到合适的 Agent

**关键设计**：
```
┌──────────────────────────────────────────────────┐
│              Clowder Platform Layer               │
│   Identity    A2A Router    Skills Framework     │
│   Manager     & Threads     & Manifest          │
└──────────────────────────────────────────────────┘
```

**学习优先级**：⭐⭐⭐⭐⭐（必学）

---

### 2. A2A (Agent-to-Agent) Communication

**定义**：Agent 之间异步消息传递协议，支持 @mention 路由和线程隔离。

**核心功能**：
- `@mention` 路由：@opus 发给架构猫，@codex 发给审查猫
- Thread 隔离：每个 feat/feature 独立线程，上下文不泄露
- Cross-post message：跨线程通知
- Handoff：接力棒传递（上一个猫完成交给下一个猫）

**External Agent Contract v1（EAC）**：
- **Invocation Contract**：支持异步消息流、取消信号
- **Stream Contract**：JSONL/NDJSON/SSE 输出格式
- **Session Contract**：跨轮 resume 能力
- **Capability Contract**：MCP 编排兼容
- **Collaboration Contract**：A2A 协议（AgentCard + tasks/send）

**学习优先级**：⭐⭐⭐⭐⭐（必学）

---

### 3. MCP (Model Context Protocol)

**定义**：Anthropic 提出的模型上下文协议，用于 Agent 与外部工具交互。

**Clowder 实现**：
- MCP Callback Bridge：让 Claude/Codex/Gemini 都能调用 Clowder 平台工具
- MCP Server 拆分：
  - `cat-cafe-collab`（协作核心，~14 tools）：三猫必装
  - `cat-cafe-memory`（记忆回溯，~9 tools）：按需
  - `cat-cafe-signals`（信号猎手，5 tools）：按需
- Prompt 瘦身：按需加载 vs 全量注入（减少 50% tool schema）

**关键工具**：
- `search_messages`：按猫/关键词过滤消息
- `list_threads`：线程发现
- `feat_index`：feat→thread 映射
- `cross_post_message`：跨线程发消息
- `list_tasks`：全局任务视图

**学习优先级**：⭐⭐⭐⭐⭐（必学，MCP 是行业标准）

---

### 4. Persistent Identity（持久化身份）

**定义**：每个 Agent 在会话和上下文压缩后保持角色、性格、记忆。

**Clowder 实现**：
- 猫有独特名字（XianXian、YanYan、ShuoShuo），不是标签
- 系统提示词动态注入身份和家规
- 跨会话记忆持久化

**Anti-Compression 设计**：
- 防止上下文压缩丢失重要信息
- Evidence Store 存储决策依据

**学习优先级**：⭐⭐⭐⭐（重要）

---

### 5. Skills Framework（Skills 框架）

**定义**：按需加载 specialized prompts，Agent 加载 specialized skills（TDD、debugging、review）只在需要时。

**Clowder 实现**：
- `manifest.yaml`：定义所有 skill 和触发条件
- 30+ 个 skill：`feat-lifecycle`、`worktree`、`quality-gate`、`request-review`、`receive-review`、`tdd`、`debugging`、`self-evolution` 等
- Skill 触发器：基于上下文的自动加载

**Self-Evolution Skill（F100）**：
- Mode A: Scope Guard（防御）- 发现 scope 发散时提醒
- Mode B: Process Evolution（防御→改进）- SOP/skills 演化
- Mode C: Knowledge Evolution（进攻→成长）- 从经验中沉淀知识

**学习优先级**：⭐⭐⭐⭐⭐（必学，skill 是复用知识的好方式）

---

### 6. SOP (Standard Operating Procedures) Auto-Guardian

**定义**：自动化的标准操作流程守护者。

**Clowder SOP 流程**：
```
⓪ Design Gate → ① worktree → ② quality-gate → ③ review → ④ merge → ⑤ 愿景守护
```

| Step | 做什么 | Skill |
|------|--------|-------|
| ⓪ | 设计确认（UX→铲屎官/后端→猫猫/架构→两边） | `feat-lifecycle` |
| ① | 创建 worktree，配置 Redis 6398 | `worktree` |
| ② | 愿景对照 + spec 合规 + 跑测试 | `quality-gate` |
| ③a | 发 review 请求 | `request-review` |
| ③b | 处理 review 反馈 | `receive-review` |
| ④ | 门禁 → PR → merge → 清理 | `merge-gate` |
| ⑤ | 愿景守护 + feat close | `feat-lifecycle` |

**设计门禁**：UX 没确认不准开 worktree

**愿景守护**：非作者非 reviewer 的猫做愿景三问

**学习优先级**：⭐⭐⭐⭐（重要，理解如何让 Agent 遵守流程）

---

### 7. Self-Evolution（自我进化）

**定义**：Agent 从错误和经验中主动学习和改进。

**F100 Self-Evolution 机制**：

**五级知识成熟度**：
| Level | 形态 | 说明 |
|-------|------|------|
| L0 | Episode | 原始记录 |
| L1 | Pattern | 草稿（≥2 相似 episode） |
| L2 | Draft | Method Card / Skill Draft（smoke gate ≥3 cases） |
| L3 | Validated | 正式 method/skill（≥6 uses, ≥80%） |
| L4 | Standard | 团队标准（≥12 uses, ≥90%） |

**三机制闭环**：
```
Episode Card（原料）→ Dual Distillation（蒸馏）→ Eval Ledger（验证）
```

**Knowledge Object Contract**：
```yaml
knowledge:
  artifact_type: episode | method | skill | proposal | eval | lesson
  domain: development | medical | legal | product | ops | general
  knowledge_type: declarative | procedural | analytical | metacognitive
  trust_level: experimental | tested | validated | production
```

**学习优先级**：⭐⭐⭐⭐（重要，Agent 自我改进是高级能力）

---

### 8. Memory & Evidence Store

**定义**：持久化存储 + 证据追踪。

**Clowder 实现**：
- Redis：会话状态、消息、任务
- Evidence SQLite + FTS5 + sqlite-vec：项目知识检索
- `search_evidence`：统一检索入口（scope/mode/depth 三维参数）
- 自动 edges 提取 + memory invalidation

**记忆分层**：
- Session Memory：单次会话
- Thread Memory：线程级
- Team Memory：跨线程共享
- Evidence Store：结构化知识

**学习优先级**：⭐⭐⭐⭐（重要，记忆是 Agent 连续性的基础）

---

### 9. Cross-Model Review（跨模型审查）

**定义**：Claude 写代码，GPT 审查，Gemini 设计，不同模型协同。

**实现**：
- Agent 配对规则：跨 family 优先、peer-reviewer 角色、available
- 同一猫不能 review 自己代码
- 降级规则：无跨 family reviewer → 同 family 不同个体 → 铲屎官

**学习优先级**：⭐⭐⭐⭐（重要，质量保障机制）

---

### 10. CLI Integration（CLI 集成架构）

**定义**：如何对接多种厂商的 AI Agent CLI。

**Clowder 支持的 CLI**：
| Agent CLI | Model Family | Output Format | MCP |
|-----------|-------------|---------------|-----|
| Claude Code | Claude | stream-json | Yes |
| Codex CLI | GPT / Codex | json | Yes |
| Gemini CLI | Gemini | stream-json | Yes |
| opencode | Multi-model | ndjson | Yes |

**NDJSON 流解析**：
```typescript
export async function* parseNDJSON(stream: Readable): AsyncGenerator<unknown> {
  // 每行一个 JSON 对象
}
```

**CLI 进程管理器**：
- 超时控制（默认 30 分钟）
- 优雅终止（SIGTERM → 3秒后 SIGKILL）
- 僵尸进程清理
- AbortSignal 支持

**学习优先级**：⭐⭐⭐⭐（重要，理解如何集成多种 Agent）

---

### 11. Hard Rails + Soft Power（架构哲学）

**定义**：硬约束 + 软实力。

**Hard Rails（硬规则）**：
- "We don't delete our own databases."
- "We don't kill our parent process."
- "Runtime config is read-only to us."
- "We don't touch each other's ports."

**Soft Power（软实力）**：
- 自由判断，结构化交付
- 自主协调，自我审查，自我改进

**Five Principles**：
| # | Principle | Meaning |
|---|-----------|---------|
| P1 | Face the final state | 每一步是基础，不是脚手架 |
| P2 | Co-creators, not puppets | 硬约束是地板，上面释放自主权 |
| P3 | Direction > speed | 不确定就停→搜→问→确认→执行 |
| P4 | Single source of truth | 每个概念只定义一次 |
| P5 | Verified = done | 证据说话，不是自信说话 |

**学习优先级**：⭐⭐⭐⭐（理解设计理念）

---

### 12. CVO (Chief Vision Officer) 模式

**定义**：人类作为首席愿景官，不是 manager，是 co-creator。

**CVO 职责**：
- 表达愿景："我希望用户做 Y 时感觉 X"
- 在关键节点做决策：设计批准、优先级、冲突解决
- 通过反馈塑造文化
- 共创：构建世界、讲故事、玩游戏

**学习优先级**：⭐⭐⭐（理解人机协作模式）

---

## 概念关联图

```
                    ┌─────────────────────────────────────────┐
                    │         Clowder Platform                │
                    │                                         │
                    │  ┌─────────┐  ┌─────────┐  ┌────────┐ │
                    │  │Identity │  │   A2A   │  │ Skills │ │
                    │  │Manager  │  │ Router  │  │Framework│ │
                    │  └────┬────┘  └────┬────┘  └───┬────┘ │
                    │       │            │           │       │
                    │  ┌────▼────────────▼───────────▼────┐ │
                    │  │     SOP Auto-Guardian             │ │
                    │  │  (Design→Worktree→Quality→Review)│ │
                    │  └────────────────┬──────────────────┘ │
                    │                   │                    │
                    │  ┌────────────────▼────────────────┐ │
                    │  │      Memory & Evidence Store     │ │
                    │  │  (Redis + SQLite + Vector)      │ │
                    │  └───────────────────────────────────┘ │
                    └─────────────────────────────────────────┘
                                        │
          ┌──────────────────────────────┼──────────────────────────────┐
          │                              │                              │
    ┌─────▼─────┐                ┌──────▼──────┐                ┌─────▼─────┐
    │  Claude    │                │    GPT      │                │  Gemini   │
    │  (XianXian)│                │  (YanYan)   │                │ (ShuoShuo) │
    └────────────┘                └─────────────┘                └───────────┘
```

---

## 与 CodeMind 对比

| 维度 | Clowder AI | CodeMind |
|------|-----------|----------|
| **目标** | 多 Agent 团队协作平台 | Java 微服务代码理解工具 |
| **交互** | 聊天 + @mention 路由 | CLI 工具 |
| **Agent** | Claude/Codex/Gemini/opencode | Python Agent + JavaParser |
| **记忆** | Redis + Evidence Store | 本地 KnowledgeStore |
| **Skills** | 30+ skill 框架 | 工具函数 |
| **用例** | 通用软件开发、陪伴、游戏 | 专注文本代码理解 |

**相似点**：多 Agent 架构、持久化记忆、Skills 框架、SOP 流程

**不同点**：Clowder 是平台层，CodeMind 是工具层

---

## 关键参考文档

| 文档 | 位置 | 内容 |
|------|------|------|
| README.md | 项目根目录 | 项目概览 |
| VISION.md | docs/VISION.md | 愿景：Cats & U |
| SOP.md | docs/SOP.md | 开发全流程 |
| cli-integration.md | docs/architecture/cli-integration.md | CLI 集成架构 |
| F002-agent-to-agent.md | docs/features/ | A2A 通信设计 |
| F043-mcp-unification.md | docs/features/ | MCP 归一化 |
| F050-a2a-external-agent-onboarding.md | docs/features/ | 外部 Agent 接入契约 |
| F100-self-evolution.md | docs/features/ | 自我进化机制 |
| cat-cafe-skills/ | cat-cafe-skills/ | 30+ skill 定义 |
| manifest.yaml | cat-cafe-skills/manifest.yaml | skill 触发器 |

---

## 涉及的技术点

### 语言/框架
- TypeScript / Node.js 20+
- pnpm workspace (monorepo)
- React + Tailwind (前端)
- Redis / SQLite + FTS5 + sqlite-vec

### 协议/标准
- MCP (Model Context Protocol)
- A2A (Agent-to-Agent)
- NDJSON 流式输出
- JSONL / SSE

### 架构模式
- Multi-Agent Orchestration
- CLI Adapter Pattern
- Skills Framework
- SOP Automation
- Evidence Store
- Self-Evolution

---

## 对 CodeMind 的借鉴

1. **A2A 通信**：CodeMind 的 3 个 Agent（Discovery/Structure/Trace）可以借鉴 @mention 路由
2. **Skills 框架**：把工具函数封装成 skill，按需加载
3. **SOP 流程**：Discovery→Structure→Trace 的流程可以显式化
4. **MCP 集成**：让 CodeMind 可以被其他 Agent 调用
5. **自我进化**：从分析结果中学习，改进分析质量

---

## 后续行动

- [ ] 把本文档同步到 clowder-ai 学习笔记
- [ ] 深入研究 MCP 协议规范
- [ ] 研究 Skills manifest 格式
- [ ] 理解 Evidence Store 实现
