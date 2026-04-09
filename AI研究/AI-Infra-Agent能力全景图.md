# AI Infra + Agent 能力全景图

> 基于腾讯 JD 分析整理，2026-04-08

---

## 🏗️ 第一层：工程基础

| 能力项 | 描述 | 现状 |
|--------|------|------|
| Python / C++ | 模型开发主力语言 | ⚠️ 有 Java 底子，Python 需补 |
| Linux / Shell | 线上环境、服务器操作 | ⚠️ 需确认 |
| Git / Docker | 代码管理、容器化部署 | ✅ 应该有 |
| CUDA 编程 | GPU 算子开发 | ❌ 缺失 |
| 分布式系统原理 | 多机通信、任务调度 | ✅ Java 后端有基础 |

---

## 🧠 第二层：ML / NLP 基础

| 能力项 | 描述 | 现状 |
|--------|------|------|
| 深度学习基础 | 反向传播、梯度、注意力机制 | ✅ NLP 硕士 |
| Transformer / GPT 架构 | 理解 LLM 底层 | ✅ 应该有 |
| 常用模型结构 | BERT、T5、LLaMA、Qwen 等 | ⚠️ 需系统化 |
| 主流框架 | PyTorch、HuggingFace | ⚠️ 需深入 |
| 训练技巧 | 学习率调度、warmup、梯度裁剪 | ⚠️ 需加强 |

---

## 🚀 第三层：LLM 后训练（Post-Training）

| 能力项 | 描述 | 现状 |
|--------|------|------|
| SFT 数据构造 | 指令数据格式化、质量筛选 | ⚠️ 需实践 |
| LoRA / QLoRA | 高效微调方法 | ⚠️ 需深入 |
| DPO / KTO | 偏好训练（免 RM） | ❌ 缺失 |
| PPO / GRPO | 强化学习训练范式 | ❌ 缺失 |
| Reward Model | 奖励模型训练 | ❌ 缺失 |
| 分布式训练 | DeepSpeed / FSDP 多卡 | ❌ 缺失 |
| Megatron-LM | 张量并行训练框架 | ❌ 缺失 |

---

## 🤖 第四层：AI Agent

| 能力项 | 描述 | 现状 |
|--------|------|------|
| ReAct / PlanAct / CodeAct | Agent 基本模式 | ⚠️ 理论有，实践少 |
| Multi-Agent 架构 | 多智能体协作 | ❌ 缺失 |
| 工具调用（Function Call） | 外部工具接入 | ⚠️ 需加强 |
| Prompt Engineering | 复杂任务 Prompt 设计 | ⚠️ 需加强 |
| Memory / Context 管理 | 上下文工程 | ⚠️ 需加强 |
| CoT / ToT / GoT | 推理增强方法 | ❌ 缺失 |
| RAG（检索增强） | 外部知识接入 | ⚠️ 需深入 |
| Agent 评测 | 评测框架与指标 | ❌ 缺失 |

---

## ⚡ 第五层：推理优化（Infra 核心）

| 能力项 | 描述 | 现状 |
|--------|------|------|
| vLLM / SGLang | 推理框架原理 | ❌ 缺失 |
| PagedAttention / KV Cache | 显存优化机制 | ❌ 缺失 |
| INT4 / INT8 量化 | 模型压缩加速 | ⚠️ 理论有 |
| FlashAttention | 注意力高效实现 | ⚠️ 需了解原理 |
| Speculative Decoding | 投机解码加速 | ❌ 缺失 |
| TensorRT / ONNX | 推理引擎 | ❌ 缺失 |
| PD 分离（Prefill-Decode） | 分离式推理架构 | ❌ 缺失 |
| CUDA / Triton 算子优化 | 自定义核函数 | ❌ 缺失 |

---

## 🛠️ 第六层：RL 训练框架

| 能力项 | 描述 | 现状 |
|--------|------|------|
| verl / OpenRLHF | RLHF 训练框架 | ❌ 缺失 |
| trlx / TRL | RLHF 工具库 | ❌ 缺失 |
| GRPO（DeepSeek） | 新一代 RL 方法 | ❌ 缺失 |
| 数据流水线 | 偏好数据构造、分布 | ❌ 缺失 |
| 在线 RL / 离线 RL | 训练范式选择 | ❌ 缺失 |

---

## 现状总结

```
✅ 有基础：  Linux操作、Git/Docker、分布式系统概念、NLP/ML理论
⚠️ 需加强： Python工程化、PyTorch/HuggingFace、SFT、RAG、Function Call
❌ 完全缺失：RL(DPO/PPO/GRPO)、推理优化(vLLM/KV Cache/量化)、分布式训练
```

## 最核心缺失（按优先级）

1. **RL 后训练（DPO → GRPO）** — AI Agent + 推理能力的分水岭
2. **推理框架原理（vLLM 源码、PagedAttention）** — AI Infra 核心战场
3. **分布式训练实战（DeepSpeed + 多卡）** — 工程底子好，补起来快
