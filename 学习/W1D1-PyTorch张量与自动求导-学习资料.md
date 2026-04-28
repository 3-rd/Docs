# W1D1｜PyTorch 张量操作 + 自动求导

> 学习日期：2026-04-10
> 目标：掌握 PyTorch 核心 API，理解自动求导机制，夯实 Day 1 基础

---

## 🔥 核心知识全景图（面试高频）

### 一、张量（Tensor）— 最底层数据结构

| 知识点 | 说明 | 面试高频追问 |
|---|---|---|
| 创建方式 | `torch.tensor()` vs `torch.Tensor()` | 两者区别（是否拷贝数据、类型推断） |
| 数据类型 | `float/int/bool/dtype` | 如何指定 device + dtype 同时创建 |
| GPU 迁移 | `.cuda()` / `.to(device)` | 如何判断 GPU 可用 |
| 形状操作 | `view / reshape / transpose / permute` | view vs reshape 区别（是否连续） |
| 索引切片 | `tensor[mask]` / `torch.masked_select` | 如何按条件筛选 |
| 运算 | `matmul / mm / @` / `torch.sum / mean / max` | 矩阵乘法哪个最快/最安全 |
| 与 NumPy 互转 | `tensor.numpy()` / `torch.from_numpy(nparr)` | 共享内存问题 |

**⚠️ 面试必答题**：
- `torch.Tensor()` 和 `torch.tensor()` 的区别？
- 张量在 GPU 和 CPU 之间转换需要注意什么？

---

### 二、自动求导（Autograd）— PyTorch 灵魂

| 知识点 | 说明 | 面试高频追问 |
|---|---|---|
| `requires_grad` | 默认为 False，设为 True 开启追踪 | 哪些操作会默认开启 |
| `backward()` | 反向传播计算梯度 | 何时调用，梯度会累加还是覆盖 |
| `grad` | 保存梯度值 | 多个 `backward()` 时梯度如何变化 |
| `grad_fn` | 记录创建张量的运算 | 用于什么 |
| `torch.no_grad()` | 前向推理时不追踪梯度 | 与 `eval()` 区别 |
| `detach()` | 截断计算图 | 何时需要 detach |
| `hook` 机制 | 注册前向/反向 hook | 有什么用 |

**⚠️ 面试必答题**：
- PyTorch 反向传播原理？计算图是怎么构建的？
- 梯度消失/爆炸的原因？在 PyTorch 中如何检测和解决？
- `backward()` 掉了梯度会怎样？多次 backward 梯度累加还是覆盖？

---

### 三、`nn.Module` — 模型构建核心

| 知识点 | 说明 | 面试高频追问 |
|---|---|---|
| 继承 `nn.Module` | 必须重写 `__init__` + `forward` | 为什么要继承 |
| `super().__init__()` | 调用父类构造函数 | 不调用会怎样 |
| `named_parameters()` / `parameters()` | 遍历模型参数 | 如何冻结部分层 |
| `state_dict()` / `load_state_dict()` | 模型序列化/加载 | 怎么只加载部分参数 |
| `children()` / `modules()` | 遍历子模块 | 区别是什么 |
| 常见层 | `Linear / Conv2d / BatchNorm / Dropout / LSTM` | 参数含义 |

**⚠️ 面试必答题**：
- `nn.Module` 的 `forward` 为什么只需写前向，反向自动搞定？
- `model(img)` 背后发生了什么？（call → forward → hooks）
- `model.train()` vs `model.eval()` 区别？（BN 和 Dropout 的行为差异）

---

### 四、`torch.optim` — 优化器

| 知识点 | 说明 | 面试高频追问 |
|---|---|---|
| SGD | 随机梯度下降 + momentum | momentum 是什么 |
| Adam / AdamW | 自适应学习率 | Adam 的原理，W 是什么 |
| 学习率调度 | `lr_scheduler` | 常用调度策略 |
| 不同层不同学习率 | optimizer 参数分组 | 怎么配 |
| `zero_grad()` | 清零梯度 | 为什么要手动调用 |

**⚠️ 面试必答题**：
- SGD 和 Adam 的区别？各自适用场景？
- 学习率衰减策略有哪些？
- 为什么梯度要用 `zero_grad()` 清零，不能累加？

---

### 五、`DataLoader` — 数据加载

| 知识点 | 说明 | 面试高频追问 |
|---|---|---|
| `Dataset` | 自定义数据抽象，必须实现 `__getitem__` + `__len__` | 如何自己实现 |
| `DataLoader` | batch / shuffle / num_workers | 各参数含义 |
| `collate_fn` | 自定义 batch 拼接 | 什么时候需要重写 |
| `pin_memory` | 加速 GPU 传输 | 什么原理 |
| `torchvision` | 图像领域数据集 |  |

---

### 六、模型保存与加载

| 方式 | 代码 | 适用场景 |
|---|---|---|
| 保存 state_dict | `torch.save(model.state_dict(), path)` | **推荐**，轻量 |
| 保存整个模型 | `torch.save(model, path)` | 不推荐，依赖类定义 |
| 加载 | `model.load_state_dict(torch.load(path))` | 常用方式 |
| 只加载部分参数 | `strict=False` / 过滤 key | 迁移学习/微调 |
| 保存优化器状态 | `torch.save({'model': ... , 'optimizer': ...}, path)` | 断点续训 |

---

### 七、GPU 加速与分布式

| 知识点 | 说明 |
|---|---|
| 单 GPU | `model.cuda()` / `tensor.to(device)` |
| 多 GPU | `nn.DataParallel(model, device_ids=[0,1,2])` |
| 多 GPU 原理 | 按 batch 维度分割 → 各 GPU 独立 forward → 梯度累加到主 GPU |
| 分布式 DDP | `DistributedDataParallel` — 工业级多机多卡 |
| `torch.cuda.is_available()` | 判断 GPU 是否可用 |

---

### 八、Fine-tuning（微调）

| 方式 | 做法 |
|---|---|
| 局部微调 | 冻结底层参数（`requires_grad=False`），只训练顶层 |
| 全局微调 | 不同层设不同学习率（在 optimizer param_groups 中配置） |
| 加载预训练 | `model = torchvision.models.resnet18(pretrained=True)` |

---

### 九、常见网络层速查

```
卷积:  nn.Conv2d(in, out, kernel, stride, padding)
池化:  nn.MaxPool2d / nn.AvgPool2d / nn.AdaptiveAvgPool2d
BN:    nn.BatchNorm2d(channels)  — train时用batch统计，eval时用全局统计
Dropout: nn.Dropout(p)           — train时随机置0，eval时全部保留
全连接: nn.Linear(in_features, out_features)
激活:  nn.ReLU / nn.Sigmoid / nn.Tanh
```

---

### 十、训练流程全链路（必须能默写）

```
1. 定义 Dataset + DataLoader
2. 定义模型 (继承nn.Module) → 放到GPU
3. 定义损失函数 (CrossEntropyLoss / MSELoss...)
4. 定义优化器 (SGD / Adam)
5. 训练循环:
   for epoch in range(E):
       model.train()
       for batch in train_loader:
           optimizer.zero_grad()       # 清梯度
           output = model(input)       # 前向
           loss = criterion(output, target)
           loss.backward()             # 反向
           optimizer.step()            # 更新参数
       model.eval()
       with torch.no_grad():           # 评估
           ...
```

---

## 📚 推荐学习资料

### 官方教程（首选）

**张量操作**：
- https://pytorch.org/tutorials/beginner/basics/tensorqs_tutorial.html
- https://pytorch.org/docs/stable/generated/torch.Tensor.requires_grad.html
- https://pytorch.org/docs/stable/tensors.html

**自动求导**：
- https://pytorch.org/tutorials/beginner/basics/autogradqs_tutorial.html
- https://pytorch.org/docs/stable/autograd.html

### 扩展阅读
- [PyTorch Autograd 机制详解](https://pytorch.org/docs/stable/notes/autograd.html)
- [PyTorch 张量 API 文档](https://pytorch.org/docs/stable/tensors.html)
- [动手学深度学习 - PyTorch 版](https://zh.d2l.ai/chapter_preliminaries/ndarray.html)

---

## 🧪 必须完成的练习

### 练习 1：view vs reshape 区别
```python
import torch

x = torch.randn(2, 3)
y = x.view(6)      # 不复制数据，共享底层存储
z = x.reshape(6)    # 可能复制（当维度不连续时）

print("view 共享存储:", y.storage().data_ptr() == x.storage().data_ptr())
print("reshape 共享存储:", z.storage().data_ptr() == x.storage().data_ptr())

# 当 tensor 不连续时，reshape 会复制
x = torch.randn(2, 3, 4)
x_t = x.transpose(0, 1)  # 变为 (3, 2, 4)，不连续
y = x_t.view(24)  # 会报错！
z = x_t.reshape(24)  # 可以，自动复制
```

### 练习 2：backward 验证
```python
import torch

x = torch.tensor([1., 2., 3.], requires_grad=True)
y = (x ** 2).sum()  # y = 1+4+9=14
y.backward()
print(x.grad)  # tensor([2., 4., 6.])  = 2x

# 验证：dy/dx = 2x
```

### 练习 3：detach 截断梯度
```python
import torch

x = torch.randn(3, requires_grad=True)
y = x ** 2
z = y.detach()  # 梯度到此截断

print("z.requires_grad:", z.requires_grad)  # False
print("y.requires_grad:", y.requires_grad)  # True
```

### 练习 4：梯度累积机制
```python
import torch

# 验证梯度累加
x = torch.tensor([1., 2., 3.], requires_grad=True)
for _ in range(3):
    y = (x ** 2).sum()
    y.backward()
    print(x.grad)

# 第二次 backward 前必须清零，否则累加
```

### 练习 5：nn.Module 实战
```python
import torch
import torch.nn as nn

class SimpleNet(nn.Module):
    def __init__(self):
        super().__init__()  # 必须调用
        self.fc1 = nn.Linear(10, 5)
        self.relu = nn.ReLU()
        self.fc2 = nn.Linear(5, 2)

    def forward(self, x):
        x = self.fc1(x)
        x = self.relu(x)
        x = self.fc2(x)
        return x

model = SimpleNet()
print("参数数量:", sum(p.numel() for p in model.parameters()))
print("结构:\n", model)
```

---

## 🔑 面试常考点（答案要点）

### Q1: `torch.Tensor()` vs `torch.tensor()` 区别？
- `torch.Tensor()` 是类构造函数，默认类型 `floatTensor`，从数据推断类型
- `torch.tensor()` 是工厂函数，类型由 `dtype` 参数指定，**拷贝数据**

### Q2: `torch.no_grad()` vs `model.eval()` 区别？
- `no_grad()`：**不构建计算图**，节省显存
- `eval()`：切换到评估模式，**BN 用全局统计量，Dropout 不生效**
- 两者可叠加：`model.eval()` + `with torch.no_grad():`

### Q3: 反向传播原理？
- PyTorch 构建**有向无环图（DAG）**，叶子节点是原始张量
- `backward()` 从输出节点反向遍历图，累加梯度到 `.grad`

### Q4: 梯度消失/爆炸的原因和解决？
- 原因：多层链式求导，连乘效应导致梯度趋近 0 或爆炸
- 解决：梯度裁剪（`torch.nn.utils.clip_grad_norm_`）、残差连接、归一化、激活函数选择

### Q5: DataLoader 的 shuffle 和 num_workers？
- `shuffle=True`：每个 epoch 打乱数据顺序，防止过拟合
- `num_workers`：并行加载进程数，过多会导致进程开销

---

> 📌 **Day 1 目标**：学完能回答上面的面试题，动手跑通练习 1-5
