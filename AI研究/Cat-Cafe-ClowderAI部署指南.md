# Cat Cafe / Clowder AI 部署指南（macOS 公司机器）

> **目标**：在公司 Mac 上完整部署 Cat Cafe（Clowder AI）多猫协作平台
> **项目地址**：https://github.com/zts212653/clowder-ai
> **更新时间**：2026-04-16

---

## 1. 项目简介

Cat Cafe（Clowder AI）是一个多 AI Agent 协作平台，核心特点：

- **多猫协作**：宪宪（Ragdoll/Claude）、砚砚（Maine Coon/Codex）、烁烁（Siamese/Gemini）、金渐层（opencode）
- **A2A 通信**：@mention 路由、线程隔离、结构化交接
- **持久记忆**：Evidence Store（SQLite）、共享知识库
- **Skills 框架**：按需加载专业技能（TDD、审查、调试等）
- **MCP 集成**：跨 Agent 工具共享
- **多平台**：Web + 飞书 + Telegram（规划中）
- **开箱即用**：只需 Node.js + pnpm + Redis（可选）

---

## 2. 部署方案选择

Cat Cafe 有两种部署架构：

| 架构 | 说明 | 适用场景 |
|------|------|---------|
| **Worktree 架构（默认）** | 开发目录(clowder-ai) + 运行时目录(cat-cafe-runtime) 分开 | 推荐 — 保持开发目录干净 |
| **Direct 架构** | 直接在当前目录运行，不创建 worktree | 简单部署、无需自动更新 |

**推荐公司机器使用 Direct 架构**，更简单、更稳定。

---

## 3. 前置要求

### 3.1 必须安装

| 工具 | 版本 | 安装命令 |
|------|------|---------|
| **Git** | 任意近期版本 | macOS 通常自带 |
| **Node.js** | >= 20.0.0 | 官网下载或 `brew install node` |
| **pnpm** | >= 9.0.0 | `npm install -g pnpm` |

```bash
# 验证安装
node --version   # 需要 >= 20.0.0
pnpm --version  # 需要 >= 9.0.0
git --version
```

### 3.2 可选安装

| 工具 | 说明 | 安装命令 |
|------|------|---------|
| **Redis** | 数据持久化（建议安装） | `brew install redis` |
| **ffmpeg** | 语音功能必需 | `brew install ffmpeg` |

> **没有 Redis 也能运行**：加 `--memory` 参数，纯内存模式（重启数据丢失）

---

## 4. 安装步骤（公司机器推荐：Direct 架构）

### Step 1：克隆代码

```bash
# 在你喜欢的位置克隆（建议 ~/AIAgent/ 下）
mkdir -p ~/AIAgent
cd ~/AIAgent
git clone https://github.com/zts212653/clowder-ai.git
cd clowder-ai
```

### Step 2：安装依赖

```bash
pnpm install
```

这一步会安装所有 pnpm workspace 的依赖，耗时约 2-5 分钟。

### Step 3：构建项目

```bash
pnpm build
```

构建所有包（api / mcp-server / shared / web），生成 `dist/` 目录。

### Step 4：配置环境变量

```bash
cp .env.example .env
```

`.env` 文件已包含合理的默认值，**公司机器通常只需要修改端口配置**。

### Step 5：启动服务

**方式 A：有 Redis（推荐，数据持久化）**

```bash
# 先启动 Redis
brew services start redis

# 确认 Redis 运行
redis-cli ping
# 应返回：PONG

# 启动 Cat Cafe
pnpm start:direct
```

**方式 B：无 Redis（纯内存，重启丢数据）**

```bash
pnpm start:direct -- --memory
```

**方式 C：后台运行（推荐公司机器使用）**

```bash
pnpm start:direct --daemon

# 查看状态
pnpm start:status

# 查看日志
tail -f cat-cafe-daemon.log

# 停止
pnpm stop
```

### Step 6：访问 Web UI

打开浏览器访问：**http://localhost:3003**

首次使用需要配置模型 API Key。

---

## 5. 配置模型（Hub → 系统配置 → 账号配置）

启动后，打开 http://localhost:3003，进入 **Hub → 系统配置 → 账号配置**。

### 5.1 支持的模型类型

| 类型 | 配置方式 | 适用 Provider |
|------|---------|--------------|
| **内置（OAuth/CLI）** | 通过 CLI 工具认证，无需 API Key | Claude（需安装 claude code）、GPT/Codex（需安装 codex）、Gemini（需安装 gemini-cli） |
| **API Key** | 输入 Key + Base URL | 任何 OpenAI/Anthropic 协议端点 |

### 5.2 添加 API Key 示例

以 Kimi（Moonshot）为例：

1. 点击 **"添加账号"**
2. 选择或添加自定义 provider
3. 填写：
   - API Key：`sk-xxx`（从 Moonshot API 平台获取）
   - Base URL：`https://api.moonshot.cn/v1`
4. 点击 **测试** 验证
5. 保存

### 5.3 国产 Provider 配置参考

| Provider | Base URL | 备注 |
|----------|----------|------|
| Kimi（Moonshot） | `https://api.moonshot.cn/v1` | |
| GLM（智谱） | `https://open.bigmodel.cn/api/paas/v4` | |
| MiniMax | `https://api.minimax.chat/v1` | |
| Qwen（通义） | `https://dashscope.aliyuncs.com/compatible-mode/v1` | |
| OpenRouter | `https://openrouter.ai/api/v1` | 聚合多个模型 |

详细配置说明见：[第三方 AI Provider 配置指南](https://github.com/zts212653/clowder-ai/blob/main/docs/guides/provider-configuration.md)

---

## 6. 完整命令参考

### 6.1 启动命令

```bash
# === Direct 架构（公司机器推荐）===
pnpm start:direct              # 前台启动
pnpm start:direct -- --memory # 无 Redis，内存模式
pnpm start:direct -- --quick  # 跳过重编译

# === 后台运行 ===
pnpm start:direct --daemon     # 后台运行
pnpm start:status             # 查看 daemon 状态
pnpm stop                      # 停止 daemon

# === Worktree 架构（可选）===
pnpm start                     # 自动创建 ../cat-cafe-runtime worktree
pnpm runtime:init              # 只创建 worktree，不启动
pnpm runtime:sync              # 同步 worktree 到最新
pnpm runtime:status            # 查看 worktree 状态
```

### 6.2 构建和开发

```bash
pnpm build    # 构建所有包
pnpm dev      # 所有包并行开发模式
pnpm test     # 运行测试
```

### 6.3 代码质量

```bash
pnpm check           # Lint + 类型检查 + Feature 文档检查
pnpm check:fix       # 自动修复 lint 问题
pnpm lint            # TypeScript 类型检查
```

---

## 7. 端口说明

| 服务 | 默认端口 | 说明 |
|------|---------|------|
| 前端（Next.js） | **3003** | Web UI 访问地址 |
| API 后端 | **3004** | 前后端通信 |
| MCP Server | **3011** | Agent 工具调用 |
| Redis | **6399** | 数据持久化 |

启动后访问 **http://localhost:3003** 即可。

---

## 8. 飞书接入（可选）

如果需要从飞书和猫猫聊天：

### 8.1 创建飞书应用

1. 前往 [飞书开放平台](https://open.feishu.cn/app) → 创建自建应用
2. 在权限管理中添加：
   - `im:message` — 读取消息
   - `im:message:send_as_bot` — 以机器人身份发消息
   - `im:resource` — 读取媒体资源
   - `im:resource:upload` — 上传媒体（语音/图片原生显示必需）

### 8.2 配置事件订阅

- 请求地址：`http://<你的公网IP或域名>:3004/api/connectors/feishu/webhook`
- 订阅事件：`im.message.receive_v1`

### 8.3 配置环境变量

在 `.env` 中添加：

```bash
FEISHU_APP_ID=cli_xxx
FEISHU_APP_SECRET=xxx
FEISHU_VERIFICATION_TOKEN=xxx
```

### 8.4 启用机器人

在飞书应用控制台 → **机器人**，启用机器人能力。

> ⚠️ **公司机器注意事项**：飞书需要公网可访问的 Webhook 地址（`:3004` 端口需对外暴露）。公司内网机器需要配置内网穿透（如 ngrok）或 VPN。

---

## 9. 语音功能（可选）

### 9.1 前置依赖

```bash
brew install ffmpeg
```

### 9.2 启动语音服务

```bash
# TTS（文字转语音）— 默认 Qwen3-TTS（三猫声线）
./scripts/tts-server.sh

# ASR（语音转文字）— Qwen3-ASR
./scripts/qwen3-asr-server.sh
```

### 9.3 配置环境变量

```bash
ASR_ENABLED=1
TTS_ENABLED=1
LLM_POSTPROCESS_ENABLED=1

WHISPER_URL=http://localhost:9876
NEXT_PUBLIC_WHISPER_URL=http://localhost:9876
TTS_URL=http://localhost:9879
```

---

## 10. 公司机器注意事项

### 10.1 网络限制

- **无公网 IP**：飞书 Webhook 需要公网地址，可考虑内网穿透方案
- **代理限制**：如果公司网络需要代理，在终端设置 `http_proxy` / `https_proxy`
- **端口占用**：确保 3003/3004/6399 端口未被占用

### 10.2 权限问题

- Node.js 编译需要一定权限，确保终端有足够权限
- Redis 数据目录 `~/.cat-cafe/redis-dev/` 需要写入权限

### 10.3 更新升级

```bash
# 方式一：Stay on current version（推荐公司机器）
# 不运行 pnpm start，直接用 pnpm start:direct
# 不自动更新，保持当前版本稳定

# 方式二：更新到新版本
git fetch
git checkout v0.x.x  # 某个 release tag
pnpm install
pnpm build
pnpm start:direct
```

### 10.4 systemd 开机自启（如果公司机器是 Linux）

```ini
# /etc/systemd/system/cat-cafe.service
[Unit]
Description=Cat Cafe / Clowder AI
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/clowder-ai
ExecStart=/usr/bin/pnpm start:direct
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now cat-cafe
sudo journalctl -u cat-cafe -f
```

---

## 11. 常见问题排查

### Q1：`pnpm install` 失败

```bash
# 清理后重试
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

### Q2：Redis 启动失败

```bash
# 检查端口占用
lsof -i :6399

# 手动启动
redis-server --port 6399

# 或使用内存模式
pnpm start:direct -- --memory
```

### Q3：前端连不上 API

确认 `.env` 中有：
```bash
NEXT_PUBLIC_API_URL=http://localhost:3004
```

### Q4：Agent 无响应

- 检查 **Hub → 系统配置 → 账号配置** 是否已添加模型账号
- 确认 API Key 有效且余额充足
- 查看终端 API 日志是否有认证错误

### Q5：端口被占用

```bash
# 查找占用进程
lsof -i :3003
lsof -i :3004
lsof -i :6399

# 杀掉进程或修改 .env 中的端口
```

---

## 12. 推荐：公司机器快速部署清单

```bash
# 1. 安装前置工具
brew install node pnpm redis

# 2. 克隆代码
mkdir -p ~/AIAgent && cd ~/AIAgent
git clone https://github.com/zts212653/clowder-ai.git
cd clowder-ai

# 3. 安装 + 构建
pnpm install
pnpm build

# 4. 配置
cp .env.example .env

# 5. 启动 Redis（可选，有 Redis 更好）
brew services start redis

# 6. 启动服务（后台运行）
pnpm start:direct --daemon

# 7. 访问配置
# 浏览器打开 http://localhost:3003
# Hub → 系统配置 → 账号配置 → 添加模型 API Key
```

---

## 13. 相关资源

- **GitHub**：https://github.com/zts212653/clowder-ai
- **官方文档**：https://github.com/zts212653/clowder-ai/blob/main/SETUP.zh-CN.md
- **第三方 Provider 配置**：https://github.com/zts212653/clowder-ai/blob/main/docs/guides/provider-configuration.md
- **Cat Cafe Skills**：https://github.com/zts212653/cat-cafe-skills
