# Ethclient 查询区块学习指南

> **预计学习时间：** 35 分钟
>
> **难度：** 基础

本指南介绍如何使用 Go 语言的 `go-ethereum` 库（ethclient）查询以太坊区块链上的区块信息。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [查询区块头](#查询区块头)
- [查询完整区块](#查询完整区块)
- [常见问题](#常见问题)
- [练习作业](#练习作业)
- [标准答案](#标准答案)

---

## 学习目标

完成本指南后，你将能够：
- 理解区块头和完整区块的区别
- 使用 ethclient 查询最新区块头
- 使用 ethclient 按区块号查询完整区块
- 获取区块的各种元数据（时间戳、难度、哈希等）
- 统计区块中的交易数量

## 前置条件

- Go 语言基础（变量、函数、错误处理）
- 以太坊基本概念（区块、交易、哈希）
- 已安装 Go 环境（1.18+）
- 拥有以太坊节点访问地址（测试网或主网）

## 核心概念

### 区块结构

```
┌─────────────────────────────────────┐
│            区块 (Block)              │
├─────────────────────────────────────┤
│  区块头 (Header)                     │
│  - 父区块哈希                        │
│  - 区块号                            │
│  - 时间戳                            │
│  - 难度                              │
│  - 状态根                            │
│  - 交易根                            │
│  - 收据根                            │
├─────────────────────────────────────┤
│  交易列表 (Transactions)             │
│  - Transaction 1                     │
│  - Transaction 2                     │
│  - ...                               │
├─────────────────────────────────────┤
│  叔块 (Uncles)                       │
└─────────────────────────────────────┘
```

### 为什么要区分区块头和完整区块？

| 类型 | 包含内容 | 使用场景 | 性能 |
|------|----------|----------|------|
| 区块头 (Header) | 仅区块元数据 | 快速获取区块信息 | 快 |
| 完整区块 (Block) | 元数据 + 交易列表 | 需要分析交易时 | 慢 |

---

## 查询区块头

### 方法：`HeaderByNumber`

```go
func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
```

### 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| `ctx` | `context.Context` | 上下文，用于超时控制 |
| `number` | `*big.Int` | 区块号，`nil` 表示最新区块 |

### 示例：获取最新区块头

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊节点
    // Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
    // Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 获取最新区块头（传入 nil）
    header, err := client.HeaderByNumber(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("区块号: %s\n", header.Number.String())
    fmt.Printf("时间戳: %d\n", header.Time)
    fmt.Printf("哈希: %s\n", header.Hash().Hex())
}
```

### 示例：获取指定区块头

```go
import "math/big"

blockNumber := big.NewInt(5671744)
header, err := client.HeaderByNumber(context.Background(), blockNumber)
```

### Header 可用字段

```go
type Header struct {
    Number        *big.Int        // 区块号
    Time          uint64          // 时间戳（Unix 时间）
    Difficulty    *big.Int        // 难度值
    Hash          common.Hash     // 区块哈希
    ParentHash    common.Hash     // 父区块哈希
    GasUsed       uint64          // 使用的 Gas
    GasLimit      uint64          // Gas 上限
    BaseFee       *big.Int        // 基础费用（EIP-1559）
    // ... 更多字段
}
```

---

## 查询完整区块

### 方法：`BlockByNumber`

```go
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
```

### 示例：查询完整区块

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
    // Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    blockNumber := big.NewInt(5671744)

    // 获取完整区块
    block, err := client.BlockByNumber(context.Background(), blockNumber)
    if err != nil {
        log.Fatal(err)
    }

    // 输出区块信息
    fmt.Printf("区块号: %d\n", block.Number().Uint64())       // 5671744
    fmt.Printf("时间戳: %d\n", block.Time())                   // 1712798400
    fmt.Printf("难度: %d\n", block.Difficulty().Uint64())      // 0
    fmt.Printf("哈希: %s\n", block.Hash().Hex())               // 0xae71...
    fmt.Printf("交易数量: %d\n", len(block.Transactions()))    // 70
    fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())       // 矿工地址
}
```

### 方法：`TransactionCount`

仅获取区块的交易数量（不获取完整交易数据）

```go
count, err := client.TransactionCount(context.Background(), block.Hash())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("交易数量: %d\n", count) // 70
```

### 性能对比

| 操作 | 数据量 | 推荐场景 |
|------|--------|----------|
| `HeaderByNumber` | 小 | 仅需要区块元数据 |
| `TransactionCount` | 小 | 仅需要交易数量 |
| `BlockByNumber` | 大 | 需要交易详情 |

---

## 常见问题

### Q1: 如何判断一个区块是否存在？

```go
block, err := client.BlockByNumber(ctx, blockNumber)
if err != nil {
    if err == ethereum.NotFound {
        fmt.Println("区块不存在")
    } else {
        log.Fatal(err)
    }
}
```

### Q2: 如何获取最新区块号？

```go
header, err := client.HeaderByNumber(context.Background(), nil)
if err != nil {
    log.Fatal(err)
}
latestBlockNumber := header.Number.Uint64()
```

### Q3: Block 和 Header 有什么区别？

```
Header (区块头)           Block (完整区块)
    ├── Number               ├── Number
    ├── Time                 ├── Time
    ├── Hash                 ├── Hash
    ├── Difficulty           ├── Difficulty
    ├── GasUsed              ├── GasUsed
    └── ...                  └── Transactions()  ← 交易列表
```

---

## 练习作业

开始练习前，请先安装依赖：

```bash
go mod tidy
```

### 作业 1：区块信息查询器（基础）

练习文件：[exercises/01-basic-info.go](exercises/01-basic-info.go)

编写一个程序，实现以下功能：

1. 连接到以太坊测试网
2. 获取最新区块号
3. 输出该区块的以下信息：
   - 区块号
   - 时间戳（转换为可读格式）
   - 区块哈希
   - 父区块哈希
   - 交易数量

**运行练习：**
```bash
# 编辑 exercises/01-basic-info.go，填充 TODO 部分
# 然后运行
go run exercises/01-basic-info.go
```

**参考答案：** [solutions/01-basic-info.go](solutions/01-basic-info.go)

---

### 作业 2：区块对比分析（进阶）

练习文件：[exercises/02-block-compare.go](exercises/02-block-compare.go)

编写一个程序，对比连续两个区块的差异：

1. 获取区块 N 和区块 N-1
2. 对比并输出：
   - Gas 使用量的变化
   - 时间间隔（秒）
   - 交易数量的变化
   - 难度变化（如果适用）

**运行练习：**
```bash
go run exercises/02-block-compare.go
```

**参考答案：** [solutions/02-block-compare.go](solutions/02-block-compare.go)

---

### 作业 3：区块浏览器（挑战）

练习文件：[exercises/03-block-explorer.go](exercises/03-block-explorer.go)

编写一个交互式命令行工具：

1. 用户输入区块号
2. 程序显示：
   - 区块基本信息
   - 所有交易的哈希列表
   - 统计信息（总 Gas 使用、平均 Gas 价格）

**运行练习：**
```bash
go run exercises/03-block-explorer.go
```

**参考答案：** [solutions/03-block-explorer.go](solutions/03-block-explorer.go)

---

## 测试网资源

### 测试网节点获取

| 提供商 | 网址 | 备注 |
|--------|------|------|
| [Alchemy](https://www.alchemy.com) | https://www.alchemy.com | 推荐，稳定 |
| [Infura](https://www.infura.io/) | https://www.infura.io/ | 老牌服务商 |
| [QuickNode](https://www.quicknode.com/) | https://www.quicknode.com/ | 需绑定信用卡 |
| [Chainstack](https://www.chainstack.com) | https://www.chainstack.com | 企业级 |
| [PublicNode](https://ethereum.publicnode.com/?sepolia) | - | 无需注册，偶尔不稳定 |

### 测试网代币水龙头

| 水龙头 | 网址 | 备注 |
|--------|------|------|
| Alchemy Faucet | https://www.alchemy.com/faucets/ethereum-sepolia | 需 Alchemy 账号 |
| Infura Faucet | https://www.infura.io/faucet/sepolia | 需 Infura 账号 |
| QuickNode Faucet | https://faucet.quicknode.com/ethereum/sepolia | 需 QuickNode 账号 |
| Optimism Faucet | https://console.optimism.io/faucet | 不需要注册账号 |

---

## 下一步学习

- [查询交易详情](./2.02-查询交易.md)
- [监听新区块](./2.03-监听新区块.md)
- [查询账户余额](./2.04-查询账户余额.md)
