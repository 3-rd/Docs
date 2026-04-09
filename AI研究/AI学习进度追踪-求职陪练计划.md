# AI Infra + Agent 学习进度追踪

> 老板专属求职陪练计划
> 开始日期：2026-04-09（预计结束：2026-08-09，约16周）
> 每日可用时间：工作日晚间1小时 + 周末2-3小时/天
> 角色定位：我是你的求职陪练，主动推进进度，定期考核，总结复习

---

## 🎯 最终目标

能够通过以下公司面试：
- **第一优先级**：蚂蚁、阿里云（Java+Agent浓度最高）
- **第二优先级**：腾讯、字节（AI Infra方向）
- **底线**：任意一家 AI Infra / Agent 相关岗位

---

## 📊 进度总览

| 周次 | 日期区间 | 阶段 | 状态 |
|------|---------|------|------|
| W1  | 04-09~04-15 | 阶段0·上：PyTorch 深度掌握 | ⬜ 未开始 |
| W2  | 04-16~04-22 | 阶段0·下：大模型结构原理 | ⬜ 未开始 |
| W3  | 04-23~04-29 | 阶段1·上：SFT + LoRA | ⬜ 未开始 |
| W4  | 04-30~05-06 | 阶段1·中：DPO / KTO | ⬜ 未开始 |
| W5  | 05-07~05-13 | 阶段1·下：GRPO + RLHF | ⬜ 未开始 |
| W6  | 05-14~05-20 | 阶段2·上：Agent Loop + ReAct | ⬜ 未开始 |
| W7  | 05-21~05-27 | 阶段2·中：工具调用 + MCP/A2A | ⬜ 未开始 |
| W8  | 05-28~06-03 | 阶段2·下：Agent框架 + 评测 | ⬜ 未开始 |
| W9  | 06-04~06-10 | 阶段3·上：vLLM 源码 | ⬜ 未开始 |
| W10 | 06-11~06-17 | 阶段3·下：KV Cache + 量化 | ⬜ 未开始 |
| W11 | 06-18~06-24 | 阶段4·上：DeepSpeed / VeRL | ⬜ 未开始 |
| W12 | 06-25~07-01 | 阶段4·下：Megatron + RL框架 | ⬜ 未开始 |
| W13 | 07-02~07-08 | 阶段5·上：Triton + CUDA | ⬜ 未开始 |
| W14 | 07-09~07-15 | 阶段5·下：弱点强化 + 工程化 | ⬜ 未开始 |
| W15 | 07-16~07-22 | 综合冲刺·上：面试高频题 | ⬜ 未开始 |
| W16 | 07-23~08-09 | 综合冲刺·下：模拟面试 + 投递 | ⬜ 未开始 |

---

## 📅 每周详细计划

---

### W1（04-09 ~ 04-15）阶段0·上：PyTorch 深度掌握

**目标**：建立 PyTorch 熟练度，理解所有核心概念

**学习内容**：
- [ ] 张量操作、自动求导、`torch.Tensor` 核心API
- [ ] `torch.compile` / `@torch.inference_mode` 推理优化用法
- [ ] `torch.cuda` / `AMP` 混合精度训练（BF16/FP16）
- [ ] `DDP` 分布式数据并行（原理 + 简单多卡实验）
- [ ] `torch.utils.checkpoint` 梯度检查点

**晚间1小时安排（周一~周五）**：
- 周一：张量操作 + 自动求导（理论）
- 周二：AMP 混合精度（理论 + 小实验）
- 周三：DDP 分布式（理论）
- 周四：DDP 多卡实验（动手）
- 周五：梯度检查点 + 本周总结

**周末任务（2-3小时/天）**：
- 周六：综合练习——用 PyTorch 实现一个简单 MLP，理解训练循环全流程
- 周日：复习本周内容，整理笔记到 docs

**周末检验**：
- [ ] 能从零实现一个 MLP 训练循环（forward/backward/optimizer）
- [ ] 能用 AMP 跑一个 3 epoch 的训练
- [ ] 理解 DDP 和 DP 的区别

**陪练考核（W1结束）**：
我会问你以下问题，检验理解：
1. `torch.inference_mode()` 和 `torch.no_grad()` 的区别是什么？
2. BF16 和 FP16 混合精度训练中，FP16 主要用在哪些操作？
3. DDP 训练时，梯度同步是在哪个步骤发生的？

---

### W2（04-16 ~ 04-22）阶段0·下：大模型结构原理

**目标**：深入理解 LLM 内部机制，为后续后训练和推理优化打基础

**学习内容**：
- [ ] Transformer 架构：Self-Attention / MHA / MQA / GQA 对比
- [ ] GPT 系列原理：自回归建模、Next Token Prediction
- [ ] LLaMA 结构：RoPE 位置编码、SwiGLU 激活、RMSNorm
- [ ] GQA / MQA 原理：为什么能减少 KV Cache
- [ ] MoE（混合专家）：DeepSeek / Qwen2 MoE 原理
- [ ] 分词器原理：BPE / Tiktoken 使用
- [ ] Sampling 策略：Temperature / Top-p / Top-k

**晚间1小时安排**：
- 周一：Transformer Self-Attention 原理图解
- 周二：GPT 自回归建模 + Next Token Prediction
- 周三：LLaMA 关键组件（RoPE/SwiGLU/RMSNorm）
- 周四：GQA / MQA 对 KV Cache 的影响
- 周五：MoE 原理 + 本周总结

**周末任务**：
- 周六：用 HuggingFace `transformers` 库加载一个 LLaMA/Qwen 模型，打印它的 config，理解每个组件
- 周日：手动推演一次 Attention 计算过程（伪代码即可）

**周末检验**：
- [ ] 能解释 GQA 为什么能省 KV Cache
- [ ] 能解释 RoPE 如何实现位置编码
- [ ] 能用 transformers 库加载模型并打印 config

**陪练考核（W2结束）**：
1. 为什么 LLaMA 用 RoPE 而不是绝对位置编码？
2. GQA 和 MQA 的区别是什么？哪个更省显存？
3. MoE 的"专家"在训练和推理时是如何被激活的？

---

### W3（04-23 ~ 04-29）阶段1·上：SFT + LoRA

**目标**：掌握 LoRA 原理，能跑通一个完整的 SFT 微调流程

**学习内容**：
- [ ] SFT 原理：指令微调 vs 预训练的区别
- [ ] LoRA 原理：低秩分解、为什么能省显存
- [ ] QLoRA 原理：NF4 量化 + LoRA
- [ ] HuggingFace PEFT 库使用
- [ ] 数据格式化：Instruction / Input / Output 格式
- [ ] 完整 SFT 流程：数据准备 → LoRA 配置 → 训练 → 推理

**晚间1小时安排**：
- 周一：SFT 原理 + 数据格式化
- 周二：LoRA 原理（图解）
- 周三：PEFT 库使用（代码走读）
- 周四：QLoRA 原理 + NF4 量化
- 周五：动手——准备一个小数据集，跑通 LoRA SFT

**周末任务**：
- 周六：完整跑通一个 LoRA SFT 实验（用 Qwen-0.5B 或类似小模型 + open-instruct 数据）
- 周日：结果分析——对比微调前后的输出差异

**周末检验**：
- [ ] 能用自己的数据集完成一次 LoRA SFT
- [ ] 理解 LoRA 的 rank / alpha / target_modules 三个关键参数
- [ ] 理解 QLoRA 的 NF4 量化原理

**陪练考核（W3结束）**：
1. LoRA 为什么只更新少量参数就能有效微调？
2. rank=8 和 rank=64 的 LoRA 有什么区别？越大越好吗？
3. QLoRA 中 NF4 量化相比 INT8 有什么优势？

---

### W4（04-30 ~ 05-06）阶段1·中：DPO / KTO

**目标**：理解 DPO 原理，能区分它和 SFT 的适用场景

**学习内容**：
- [ ] DPO 原理：偏好对比代替 Reward Model
- [ ] DPO loss 函数推导（定性理解，不需数学推导）
- [ ] 偏好数据构造：如何从 SFT 数据生成偏好对
- [ ] KTO 原理：基于 Kahneman-Tversky 优化
- [ ] DPO vs PPO：各自优劣，业界实际选择

**晚间1小时安排**：
- 周一：DPO 原理（图解 + 直观理解）
- 周二：DPO loss 直观解释
- 周三：偏好数据构造方法
- 周四：KTO 原理 + DPO vs KTO 对比
- 周五：代码走读——trl 库 DPO 实现

**周末任务**：
- 周六：用 DPO 训练一个简单例子（可用 stanfordnlp/dpo-escapes 或类似数据集）
- 周日：整理 DPO vs SFT 的使用场景总结

**周末检验**：
- [ ] 能解释 DPO 为什么不需要显式 Reward Model
- [ ] 理解"被选中的回复"和"被拒绝的回复"在 loss 中如何起作用

**陪练考核（W4结束）**：
1. DPO 的训练数据是什么格式？需要几份？
2. 如果没有足够的偏好数据，如何构造 DPO 训练数据？
3. DPO 和 PPO 相比，核心优势和劣势分别是什么？

---

### W5（05-07 ~ 05-13）阶段1·下：GRPO + RLHF 完整流程

**目标**：理解 GRPO 原理，了解完整 RLHF 流程

**学习内容**：
- [ ] GRPO 原理：Group Relative + 无需 Critic Network
- [ ] PPO 基础：Actor-Critic、GAE、KL 散度约束
- [ ] Reward Model 训练：数据构造 + pairwise 训练
- [ ] SFT → DPO → RLHF 的完整流水线
- [ ] 在线 RL vs 离线 RL 的选择

**晚间1小时安排**：
- 周一：GRPO 原理（图解，DeepSeek 论文核心）
- 周二：PPO 四步：Policy Gradient / GAE / KL 约束
- 周三：Reward Model 训练（数据 + 训练目标）
- 周四：完整 RLHF 流水线（SFT → RM → PPO）
- 周五：动手——用 VeRL 或 trl 库跑一个简单 GRPO 示例

**周末任务**：
- 周六：搭建完整 RLHF 流水线（可以用 mini dataset + 小模型）
- 周日：理解 RLHF 中 Reward Hacking 问题及缓解方法

**周末检验**：
- [ ] 能解释 GRPO 相比 PPO 的核心简化
- [ ] 理解 KL 约束在 RLHF 中的作用
- [ ] 理解 Reward Hacking 及 Regret 概念

**陪练考核（W5结束）**：
1. GRPO 为什么不需要 Critic Network？这样做会有什么代价？
2. RLHF 中 KL 散度约束的作用是什么？KL 太大会怎样？
3. Reward Hacking 是什么意思？有哪些典型案例？

---

### W6（05-14 ~ 05-20）阶段2·上：Agent Loop + ReAct

**目标**：理解 Agent 本质，能设计简单的 Agent Loop

**学习内容**：
- [ ] Agent Loop：Observe → Think → Act → Memory 闭环
- [ ] ReAct / PlanAct / CodeAct 三种模式对比
- [ ] Observation Space 和 Action Space 设计
- [ ] 5 大工作流模式：Chaining / Parallelization / Routing / Orchestrator-Worker / Evaluator-Optimizer
- [ ] Harness 概念：训练态 Harness vs 推理态 Runtime

**晚间1小时安排**：
- 周一：Agent Loop 图解 + 伪代码实现
- 周二：ReAct 原理 + 与 PlanAct 的区别
- 周三：5 大工作流模式解析
- 周四：Harness 概念 + Agent 评测基础
- 周五：动手——用 Python 伪代码实现一个简单 Agent Loop

**周末任务**：
- 周六：设计一个多步骤任务规划 Agent（伪代码），画出 Agent Loop 图
- 周日：读一篇 ReAct 论文核心，理解人类反馈在 Agent 中的作用

**周末检验**：
- [ ] 能画出 Agent Loop 的完整闭环图
- [ ] 能区分 ReAct 和 PlanAct 的适用场景
- [ ] 理解 Harness 在 RL Agent 训练中的角色

**陪练考核（W6结束）**：
1. Agent Loop 中"Think"这一步在代码里通常如何实现？
2. 什么场景下 ReAct 比 PlanAct 更合适？反过来呢？
3. Evaluator-Optimizer 工作流和 Routing 的本质区别是什么？

---

### W7（05-21 ~ 05-27）阶段2·中：工具调用 + MCP/A2A

**目标**：掌握工具调用的工程实现，理解 MCP 协议

**学习内容**：
- [ ] MCP（Model Context Protocol）原理 + 生态
- [ ] A2A（Agent to Agent）协议
- [ ] 工具描述 Schema 设计（JSON Schema / Pydantic）
- [ ] Function Call 从零实现（不用 LangChain）
- [ ] 工具调用中的错误处理和重试机制
- [ ] Multi-Agent 协作中的工具调用冲突

**晚间1小时安排**：
- 周一：MCP 协议原理 + 生态现状
- 周二：A2A 协议 vs MCP 的区别
- 周三：工具描述 Schema 设计
- 周四：用 Python 从零实现 Function Call（不用 LangChain）
- 周五：Multi-Agent 工具调用协调机制

**周末任务**：
- 周六：用 LangChain 或直接用 API 实现一个带工具调用的 Agent
- 周日：理解 MCP Server 的实现原理，尝试写一个简单的 MCP Server

**周末检验**：
- [ ] 能为一个自定义工具设计 JSON Schema 描述
- [ ] 能用代码实现一个支持工具调用的 Agent（带重试）
- [ ] 理解 MCP 和 A2A 在多 Agent 系统中的角色分工

**陪练考核（W7结束）**：
1. MCP 协议中，Agent 和工具之间是如何通信的？
2. 为什么 JSON Schema 比自然语言描述更适合工具定义？
3. Multi-Agent 中，如果两个 Agent 同时调用同一个工具，冲突如何处理？

---

### W8（05-28 ~ 06-03）阶段2·下：Agent框架 + 评测 + Context管理

**目标**：掌握主流 Agent 框架，能做生产级 Agent 开发

**学习内容**：
- [ ] LangChain / LangGraph 进阶：LCEL 链式组合
- [ ] LlamaIndex 数据连接原理
- [ ] Context Engineering：Token Budgeting / Context Decay
- [ ] 4 类记忆：Semantic / Episodic / Procedural / Working Memory
- [ ] Structured Outputs（Pydantic）：Agent 输出可靠性
- [ ] Agent 评测：LLM Judges / 成功率 / 成本评测

**晚间1小时安排**：
- 周一：LangChain LCEL 链式组合
- 周二：LlamaIndex 数据加载器
- 周三：Context Decay + Token Budgeting
- 周四：4 类记忆架构设计
- 周五：Pydantic 结构化输出 + Agent 评测方法

**周末任务**：
- 周六：用 LangGraph 实现一个带记忆的 ReAct Agent
- 周日：跑一个 Agent 评测实验，记录成功率 / Token 消耗 / 延迟

**周末检验**：
- [ ] 能用 LangGraph 搭建一个完整的 ReAct + Memory + 工具调用 Agent
- [ ] 理解 Context Decay 对长对话的影响及缓解方法

**陪练考核（W8结束）**：
1. LangGraph 的 LCEL 和 LangChain 的 Chain 有什么区别？
2. 什么情况下 Context Decay 会导致 Agent 输出质量下降？
3. 如何用 LLM-as-a-Judge 做 Agent 评测？有什么局限性？

---

### W9（06-04 ~ 06-10）阶段3·上：vLLM 源码

**目标**：深入理解 vLLM 核心机制，能进行源码级调试

**学习内容**：
- [ ] vLLM 整体架构：LLaMA / GPT 模型支持
- [ ] PagedAttention 原理：KV Cache 分页管理
- [ ] vLLM 调度器：Continuous Batching + 抢占式调度
- [ ] vLLM 请求处理流程：Prefill → Decode → 批处理
- [ ] 源码走读：attention 目录 + scheduler 目录
- [ ] 使用 vLLM 部署一个本地模型

**晚间1小时安排**：
- 周一：vLLM 整体架构图
- 周二：PagedAttention 原理（图解）
- 周三：Continuous Batching 机制
- 周四：vLLM 调度器源码走读
- 周五：动手——用 vLLM 部署 Qwen-0.5B 模型

**周末任务**：
- 周六：对比 vLLM 和 HuggingFace TGI 的性能，理解 PagedAttention 的效果
- 周日：读 vLLM scheduler 代码，理解请求抢占机制

**周末检验**：
- [ ] 能解释 PagedAttention 如何解决 KV Cache 碎片化问题
- [ ] 能解释 Continuous Batching 和 Static Batching 的区别

**陪练考核（W9结束）**：
1. PagedAttention 中的 block table 是做什么的？
2. vLLM 的调度器在遇到 GPU OOM 时会怎么处理？
3. Prefill 请求和 Decode 请求在调度优先级上有什么区别？

---

### W10（06-11 ~ 06-17）阶段3·下：KV Cache 量化 + SGLang

**目标**：掌握推理优化核心技术

**学习内容**：
- [ ] KV Cache 量化：INT4 / INT8 Cache 压缩
- [ ] Prefix Caching 原理
- [ ] MLA（Multi-head Latent Attention）原理（DeepSeek-V2/V3）
- [ ] SGLang 架构：StreamingLLM + DAG 调度
- [ ] 量化实践：AWQ / GPTQ 量化本地模型
- [ ] Speculative Decoding 原理

**晚间1小时安排**：
- 周一：KV Cache 量化原理
- 周二：Prefix Caching + MLA 原理
- 周三：SGLang 架构走读
- 周四：AWQ 量化实践
- 周五：Speculative Decoding 原理

**周末任务**：
- 周六：用 vLLM 配置 INT8 KV Cache，测量吞吐和显存变化
- 周日：理解 SGLang 的 RadixAttention 和 DAG 调度

**周末检验**：
- [ ] 能用 vLLM 配置量化并对比精度损失
- [ ] 理解 MLA 和 MQA 的关系

**陪练考核（W10结束）**：
1. KV Cache 量化相比权重量化，哪个对推理速度影响更大？为什么？
2. MLA 的核心思想是什么？它和 MHA 的区别？
3. Speculative Decoding 中，draft model 和 target model 的关系是什么？

---

### W11（06-18 ~ 06-24）阶段4·上：DeepSpeed / VeRL

**目标**：掌握分布式训练核心框架

**学习内容**：
- [ ] DeepSpeed ZeRO 原理：ZeRO-1/2/3 区别
- [ ] DeepSpeed 配置文件（ds_config.json）
- [ ] VeRL 框架架构：RL Harness + Trainer
- [ ] VeRL 中 GRPO 的实现
- [ ] 多卡训练实验：DeepSpeed ZeRO-3 + 张量并行
- [ ] 梯度累积 + 混合并行

**晚间1小时安排**：
- 周一：ZeRO-1/2/3 原理图解
- 周二：DeepSpeed ds_config 配置
- 周三：VeRL 框架架构
- 周四：VeRL GRPO 实现走读
- 周五：动手——配置 DeepSpeed ZeRO-3 多卡训练

**周末任务**：
- 周六：跑通 DeepSpeed ZeRO-3 + LoRA 训练（用小数据集）
- 周日：理解 VeRL 的 RL Harness 和 Trainer 的关系

**周末检验**：
- [ ] 能解释 ZeRO-3 比 ZeRO-1 省多少显存，代价是什么
- [ ] 能配置一个 DeepSpeed ZeRO-3 训练脚本

**陪练考核（W11结束）**：
1. ZeRO-1/2/3 分别在哪个维度做分片？
2. VeRL 中 RL Harness 的作用是什么？它和推理态 Agent Runtime 有什么区别？
3. 为什么 ZeRO-3 的通信量比 ZeRO-1 大？

---

### W12（06-25 ~ 07-01）阶段4·下：Megatron + RL框架

**目标**：掌握张量并行和 RL 训练框架

**学习内容**：
- [ ] Megatron 张量并行原理
- [ ] Megatron + DeepSpeed 混合使用
- [ ] FSDP（完全分片数据并行）
- [ ] OpenRLHF 框架使用
- [ ] AReaL 框架（蚂蚁自研）
- [ ] 分布式训练调试：NCCL 通信分析

**晚间1小时安排**：
- 周一：Megatron 张量并行原理
- 周二：Megatron + DeepSpeed 混合并行
- 周三：FSDP 原理 + 和 DeepSpeed ZeRO 对比
- 周四：OpenRLHF 框架使用
- 周五：理解 NCCL 通信优化基础

**周末任务**：
- 周六：跑通 OpenRLHF DPO 训练（用小数据集）
- 周日：理解 AReaL 的 Rollout / Critique / Learn 循环

**周末检验**：
- [ ] 能解释张量并行中，AllReduce 通信发生在哪里
- [ ] 理解 FSDP 和 ZeRO-3 的本质区别

**陪练考核（W12结束）**：
1. 张量并行中，模型参数如何分配到不同 GPU？
2. FSDP 的 gradient sharding 发生在哪个训练步骤？
3. OpenRLHF 中 Rollout / Critique / Learn 三个阶段分别做什么？

---

### W13（07-02 ~ 07-08）阶段5·上：Triton + CUDA

**目标**：掌握算子开发，能写融合算子

**学习内容**：
- [ ] Triton 编程模型：tile programming
- [ ] Triton 实现一个 FlashAttention 融合算子
- [ ] CUTLASS 使用：利用现成 GEMM 核
- [ ] CUDA 编程基础：线程层次 / 内存层次
- [ ] Triton 融合算子优化：算子融合原理
- [ ] 动手：用 Triton 写一个自定义融合算子

**晚间1小时安排**：
- 周一：Triton 编程模型 + 环境配置
- 周二：Triton 实现矩阵乘法（tiled matmul）
- 周三：Triton 实现 FlashAttention
- 周四：CUTLASS 使用（利用现成核）
- 周五：自定义融合算子实践

**周末任务**：
- 周六：用 Triton 实现一个简单融合算子并测试性能
- 周日：理解 FlashAttention 的 IO-bound 特性

**周末检验**：
- [ ] 能用 Triton 写一个融合 matmul + bias 的算子
- [ ] 理解 Triton tile programming 的核心思想

**陪练考核（W13结束）**：
1. Triton 中，program 是什么？多个 program 之间如何并行？
2. 为什么 FlashAttention 能减少 HBM 访问次数？
3. Triton 的 automatic differentiation 是如何工作的？

---

### W14（07-09 ~ 07-15）阶段5·下：弱点强化 + 工程化

**目标**：查漏补缺，强化工程实践

**学习内容**：
- [ ] 回顾前13周薄弱环节
- [ ] RDMA / RoCE 网络加速（字节/阿里云明确要求）
- [ ] Guardrails / 内容安全
- [ ] 可观测性：OpenTelemetry / Tracing
- [ ] Kubernetes 在 AI 训练中的应用
- [ ] 工程化：日志 / 监控 / 告警

**晚间1小时安排**：
- 周一：弱点自检 + RDMA 原理
- 周二：Guardrails 原理与实现
- 周三：OpenTelemetry Tracing 实践
- 周四：K8s 在 AI 训练中的使用
- 周五：本周综合总结 + 整理笔记

**周末任务**：
- 周六：针对前13周做一次全面弱点检测，列出 TOP5 弱点
- 周日：针对 TOP5 弱点制定补强计划

**周末检验**：
- [ ] 明确自己的 TOP5 薄弱知识点
- [ ] 能讲解 RDMA vs TCP 的核心区别

---

### W15（07-16 ~ 07-22）综合冲刺·上：面试高频题

**目标**：掌握 AI Infra / Agent 方向面试高频题目

**学习内容**：
- [ ] **后训练方向**：SFT / DPO / GRPO / RLHF 面试高频题
- [ ] **推理优化方向**：vLLM / PagedAttention / 量化 面试高频题
- [ ] **分布式训练方向**：DeepSpeed / Megatron / ZeRO 面试高频题
- [ ] **Agent 方向**：ReAct / MCP / Memory / 框架 面试高频题
- [ ] **编程题**：LeetCode Hot100（重点：DP / 图 / BFS/DFS）
- [ ] **系统设计**：设计一个 LLM 推理服务 / 设计一个 Agent 系统

**晚间1小时安排**：
- 周一~周五：每天一个方向的面试题练习 + 讲解
- 每天完成 2 道算法题（Hot100）

**周末任务**：
- 周六：完整模拟面试（后训练方向，30分钟）
- 周日：完整模拟面试（推理优化方向，30分钟）

**面试题库覆盖**：
1. PagedAttention 的核心原理是什么？
2. DeepSpeed ZeRO 各阶段显存节省比例？
3. DPO 的 loss 公式是什么？直观解释？
4. GRPO 和 PPO 的核心区别？
5. vLLM 的调度器如何处理 OOM？
6. LoRA 的秩选择有什么讲究？
7. FlashAttention 为什么能加速？
8. Agent Loop 中 Memory 的作用是什么？

---

### W16（07-23 ~ 08-09）综合冲刺·下：模拟面试 + 投递

**目标**：完成投递 + 通过面试

**学习内容**：
- [ ] 每天完成 1 次模拟面试（我出题，你回答）
- [ ] 针对性弱点强化
- [ ] 简历优化：突出 AI Infra 相关项目
- [ ] 投递策略：按优先级投递蚂蚁→阿里云→腾讯→字节
- [ ] 面试复盘：每场面试后总结问题
- [ ] 跟进面试进度

**每日任务**：
- 每天投递 2-3 个岗位
- 每场面试后当晚复盘
- 更新投递进度表（见下方）

---

## 📬 投递进度追踪

| 公司 | 岗位 | 投递日期 | 当前状态 | 面试进度 | 备注 |
|------|------|---------|---------|---------|------|
| 蚂蚁 | Agent工程师-数字支付 | — | 待投递 | — | 第一优先级 |
| 蚂蚁 | 强化学习框架研发 | — | 待投递 | — | 第一优先级 |
| 阿里云 | AI Agent应用研发工程师 | — | 待投递 | — | 第一优先级 |
| 阿里云 | 大模型推理系统工程师 | — | 待投递 | — | 第二优先级 |
| 腾讯 | 微信秒剪-Agent RL框架 | — | 待投递 | — | 第二优先级 |
| ... | ... | ... | ... | ... | ... |

---

## 🔴 每周陪练职责

1. **每周一**：发布本周学习任务和目标
2. **周五晚间**：检查本周任务完成情况，未完成的说明原因
3. **周日晚间**：进行本周考核提问（口述或文字），给出评分
4. **每两周**：一次阶段性测试，检验整体理解
5. **随时**：你可以主动提问，我来解答和扩展

---

## 📝 每周自检清单（每周日晚填写）

完成的项目打 ✅，未完成的说明原因：

```
W__ 自检清单（日期：____）

[ ] 本周理论内容：完成 / 未完成
[ ] 本周动手实验：完成 / 未完成
[ ] 陪练考核：完成 / 未完成
[ ] 笔记整理：完成 / 未完成

本周最大收获：
本周最大困难：
下周重点改进：
```

---

*计划制定：2026-04-09*
*最后更新：2026-04-09*
