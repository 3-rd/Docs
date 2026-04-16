# Cat Cafe / Clowder AI 部署指南（Windows）

> **目标**：在公司 Windows 机器上完整部署 Cat Cafe（Clowder AI）多猫协作平台
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

---

## 2. Windows 安装方式

Windows 有两种安装方式，推荐使用**一键脚本安装**（PowerShell）。

| 方式 | 说明 | 适用场景 |
|------|------|---------|
| **一键脚本（推荐）** | `install.ps1` 自动安装所有依赖 | 公司机器快速部署 |
| **手动安装** | 手动安装 Node/pnpm/Redis | 了解细节、特殊环境 |

---

## 3. 前置要求

### 必须安装

| 工具 | 版本 | 安装方式 |
|------|------|---------|
| **Git** | 任意近期版本 | [git-scm.com](https://git-scm.com/download/win) |
| **PowerShell** | 5.0+ | Windows 10/11 通常自带 |

### 自动安装的工具

`install.ps1` 脚本会自动检测并安装：

- **Node.js** >= 20.0.0（通过 winget 或手动检测）
- **pnpm** >= 9.0.0（通过 npm 或 corepack 安装）
- **Redis**（通过 winget 或便携版）

---

## 4. 安装步骤（一键脚本，推荐）

### Step 1：下载代码

```powershell
# 在你喜欢的位置克隆（建议 ~/AIAgent/ 下）
mkdir $HOME\AIAgent -Force
cd $HOME\AIAgent
git clone https://github.com/zts212653/clowder-ai.git
cd clowder-ai
```

### Step 2：运行一键安装脚本

```powershell
# 进入项目目录
cd clowder-ai

# 运行安装脚本（自动安装所有依赖）
.\scripts\install.ps1
```

`install.ps1` 会依次执行：

| 步骤 | 内容 |
|------|------|
| Step 1 | 环境检测（PowerShell/Git/Node.js） |
| Step 2 | Node.js >= 20 和 pnpm >= 9 安装 |
| Step 3 | Redis 安装 |
| Step 4 | 生成 `.env` 配置文件 |
| Step 5 | 安装依赖（pnpm install）和构建（pnpm build） |
| Step 6 | Skills 挂载 |
| Step 7 | AI CLI 工具安装（Claude/Codex/Gemini，可选） |
| Step 8 | 认证配置 |
| Step 9 | 验证并启动 |

### Step 3：启动服务

```powershell
# 生产模式启动（推荐）
.\scripts\start-windows.ps1

# 跳过重编译，快速启动
.\scripts\start-windows.ps1 -Quick

# 无 Redis（纯内存，重启数据丢失）
.\scripts\start-windows.ps1 -Memory

# 开发模式（热更新）
.\scripts\start-windows.ps1 -Dev

# 开启调试日志
.\scripts\start-windows.ps1 -Debug
```

### Step 4：访问 Web UI

打开浏览器访问：**http://localhost:3003**

---

## 5. 详细安装说明

### 5.1 `install.ps1` 脚本参数

```powershell
.\scripts\install.ps1           # 安装 + 验证
.\scripts\install.ps1 -Start    # 安装完成后自动启动
.\scripts\install.ps1 -SkipBuild # 跳过构建（已有 build 时使用）
.\scripts\install.ps1 -SkipCli   # 跳过 AI CLI 工具安装
.\scripts\install.ps1 -SkipPreflight # 跳过网络预检
.\scripts\install.ps1 -Debug    # 调试模式
```

### 5.2 `start-windows.ps1` 脚本参数

```powershell
.\scripts\start-windows.ps1       # 标准启动（生产模式）
.\scripts\start-windows.ps1 -Quick   # 跳过重编译
.\scripts\start-windows.ps1 -Memory  # 无 Redis，内存模式
.\scripts\start-windows.ps1 -Dev     # 开发模式（next dev，热更新）
.\scripts\start-windows.ps1 -Debug   # 调试日志
```

### 5.3 停止服务

```powershell
.\scripts\stop-windows.ps1
```

---

## 6. 手动安装（可选，了解细节）

如果一键脚本失败，或需要更精细的控制，按以下步骤手动安装。

### Step 1：安装 Git

下载安装：https://git-scm.com/download/win

验证：
```powershell
git --version
```

### Step 2：安装 Node.js

方式一：winget（推荐）
```powershell
winget install OpenJS.NodeJS.LTS
```

方式二：手动下载 https://nodejs.org/

验证：
```powershell
node --version   # 需要 >= 20.0.0
npm --version
```

### Step 3：安装 pnpm

```powershell
npm install -g pnpm
pnpm --version    # 需要 >= 9.0.0
```

### Step 4：安装 Redis

方式一：winget
```powershell
winget install Redis.Redis
```

方式二：Chocolatey
```powershell
choco install redis-64
```

方式三：便携版（脚本自动下载）

验证：
```powershell
redis-cli ping
# 应返回：PONG
```

### Step 5：克隆并安装项目

```powershell
cd $HOME\AIAgent
git clone https://github.com/zts212653/clowder-ai.git
cd clowder-ai

# 安装依赖
pnpm install

# 构建项目
pnpm build

# 配置环境变量
Copy-Item .env.example .env

# 启动 Redis（如果手动安装）
redis-server --port 6399

# 启动服务
.\scripts\start-windows.ps1
```

---

## 7. 配置模型（Hub → 系统配置 → 账号配置）

启动后，打开 http://localhost:3003，进入 **Hub → 系统配置 → 账号配置**。

### 7.1 支持的模型类型

| 类型 | 配置方式 | 适用 Provider |
|------|---------|--------------|
| **内置（OAuth/CLI）** | 通过 CLI 工具认证，无需 API Key | Claude（需安装 claude code）、GPT/Codex（需安装 codex）、Gemini（需安装 gemini-cli） |
| **API Key** | 输入 Key + Base URL | 任何 OpenAI/Anthropic 协议端点 |

### 7.2 国产 Provider 配置参考

| Provider | Base URL |
|----------|----------|
| Kimi（Moonshot） | `https://api.moonshot.cn/v1` |
| GLM（智谱） | `https://open.bigmodel.cn/api/paas/v4` |
| MiniMax | `https://api.minimax.chat/v1` |
| Qwen（通义） | `https://dashscope.aliyuncs.com/compatible-mode/v1` |
| OpenRouter | `https://openrouter.ai/api/v1` |

---

## 8. 飞书接入（可选）

### 8.1 创建飞书应用

1. 前往 [飞书开放平台](https://open.feishu.cn/app) → 创建自建应用
2. 在权限管理中添加：
   - `im:message` — 读取消息
   - `im:message:send_as_bot` — 以机器人身份发消息
   - `im:resource` — 读取媒体资源
   - `im:resource:upload` — 上传媒体

### 8.2 配置事件订阅

- 请求地址：`http://<你的公网IP>:3004/api/connectors/feishu/webhook`
- 订阅事件：`im.message.receive_v1`

### 8.3 配置环境变量

在 `.env` 中添加：

```powershell
FEISHU_APP_ID=cli_xxx
FEISHU_APP_SECRET=xxx
FEISHU_VERIFICATION_TOKEN=xxx
```

> ⚠️ **公司 Windows 机器注意事项**：飞书需要公网可访问的 Webhook 地址，Windows 机器需要配置内网穿透（如 ngrok）。

---

## 9. 端口说明

| 服务 | 默认端口 | 说明 |
|------|---------|------|
| 前端（Next.js） | **3003** | Web UI 访问地址 |
| API 后端 | **3004** | 前后端通信 |
| MCP Server | **3011** | Agent 工具调用 |
| Redis | **6399** | 数据持久化 |

---

## 10. 公司 Windows 机器注意事项

### 10.1 网络限制

- **代理/防火墙**：如果公司网络需要代理，在 PowerShell 设置 `http_proxy` / `https_proxy`
- **端口占用**：确保 3003/3004/6399 端口未被占用
```powershell
netstat -ano | findstr "3003"
netstat -ano | findstr "3004"
```

### 10.2 权限问题

- **管理员权限**：某些操作可能需要管理员 PowerShell
- **执行策略**：如果脚本被阻止，运行：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 10.3 常见问题

**Q1：脚本无法运行**

```powershell
# 检查执行策略
Get-ExecutionPolicy

# 设置允许本地脚本
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 运行脚本
.\scripts\install.ps1
```

**Q2：Redis 启动失败**

```powershell
# 检查端口占用
netstat -ano | findstr "6399"

# 手动启动 Redis
redis-server --port 6399

# 或使用内存模式
.\scripts\start-windows.ps1 -Memory
```

**Q3：Node 版本不对**

```powershell
# 检查版本
node --version

# 用 winget 升级
winget install OpenJS.NodeJS.LTS

# 或下载安装 https://nodejs.org/
```

**Q4：pnpm 安装失败**

```powershell
# 手动安装
npm install -g pnpm

# 验证
pnpm --version
```

**Q5：端口被占用**

```powershell
# 查找占用进程
netstat -ano | findstr "3003"
netstat -ano | findstr "3004"

# 杀掉进程（替换 PID 为实际值）
taskkill /PID <PID> /F
```

---

## 11. Windows 快速部署清单

```powershell
# 1. 克隆代码
mkdir $HOME\AIAgent -Force
cd $HOME\AIAgent
git clone https://github.com/zts212653/clowder-ai.git
cd clowder-ai

# 2. 一键安装（自动处理所有依赖）
.\scripts\install.ps1

# 3. 启动服务
.\scripts\start-windows.ps1

# 4. 访问配置
# 浏览器打开 http://localhost:3003
# Hub → 系统配置 → 账号配置 → 添加模型 API Key
```

---

## 12. 相关资源

- **GitHub**：https://github.com/zts212653/clowder-ai
- **官方文档（中文）**：https://github.com/zts212653/clowder-ai/blob/main/SETUP.zh-CN.md
- **第三方 Provider 配置**：https://github.com/zts212653/clowder-ai/blob/main/docs/guides/provider-configuration.md
- **Cat Cafe Skills**：https://github.com/zts212653/cat-cafe-skills
