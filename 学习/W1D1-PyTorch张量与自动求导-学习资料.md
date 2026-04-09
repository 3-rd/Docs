# W1 周一｜PyTorch 张量操作 + 自动求导

> 学习日期：待定
> 目标：掌握 PyTorch 核心 API，理解自动求导机制

---

## 📚 推荐学习资料

### 1. PyTorch 官方教程（首选）

**张量操作**：
- https://pytorch.org/tutorials/beginner/basics/tensorqs_tutorial.html
- https://pytorch.org/docs/stable/generated/torch.Tensor.requires_grad.html
- https://pytorch.org/docs/stable/tensors.html

**自动求导**：
- https://pytorch.org/tutorials/beginner/basics/autogradqs_tutorial.html
- https://pytorch.org/docs/stable/autograd.html

### 2. 核心概念速查

#### requires_grad
```python
# 创建时开启
x = torch.tensor([1., 2., 3.], requires_grad=True)

# 后续开启
x.requires_grad_(True)

# 关闭（节省显存）
x = x.detach()
```

#### view vs reshape
```python
x = torch.randn(2, 3)
y = x.view(6)      # 不复制数据，共享底层存储
z = x.reshape(6)    # 可能复制（当维度不连续时）

# 验证
assert y.storage().data_ptr() == x.storage().data_ptr()  # view 共享
```

#### backward 求导
```python
x = torch.tensor([1., 2., 3.], requires_grad=True)
y = (x ** 2).sum()  # y = 1+4+9 = 14
y.backward()
print(x.grad)  # tensor([2., 4., 6.])
```

#### 梯度累积
```python
# 每次迭代前清零
optimizer.zero_grad()

# 三种方式的区别：
optimizer.zero_grad()           # 推荐：只清零叶子节点
optimizer.zero_grad(set_to_none=True)  # 更彻底：设为 None
model.zero_grad()              # 只清零模型参数
```

---

## 🧪 必须完成的练习

### 练习 1：view vs reshape 区别
```python
import torch

x = torch.randn(2, 3)
y = x.view(6)
z = x.reshape(6)

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

---

## 🔑 面试常考点

1. **`requires_grad` 的作用是什么？**
   - 开启后，PyTorch 会记录该张量的操作历史，用于反向传播

2. **`torch.no_grad()` 和 `torch.inference_mode()` 的区别？**
   - `no_grad()`：所有张量 `requires_grad=False`，节省显存
   - `inference_mode()`：更快，但禁止梯度计算，禁止原地操作

3. **`backward()` 背后的计算图是什么？**
   - PyTorch 会构建一个 DAG（有向无环图）
   - 叶子节点是 `requires_grad=True` 的原始张量

4. **为什么梯度会累加？**
   - 如果不清零 `zero_grad()`，梯度会累加到 `.grad` 上

---

## 📖 扩展阅读

- [PyTorch Autograd 机制详解](https://pytorch.org/docs/stable/notes/autograd.html)
- [PyTorch 张量 API 文档](https://pytorch.org/docs/stable/tensors.html)
- [动手学深度学习 - PyTorch 版](https://zh.d2l.ai/chapter_preliminaries/ndarray.html)
