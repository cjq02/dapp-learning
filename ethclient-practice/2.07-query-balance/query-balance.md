# Ethclient 查询账户余额学习指南

> **预计学习时间：** 20 分钟
>
> **难度：** 基础

本指南介绍如何使用 Go 语言的 `go-ethereum` 库（ethclient）查询以太坊地址的 ETH 余额。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [查询最新余额](#查询最新余额)
- [查询指定区块余额](#查询指定区块余额)
- [查询待处理余额](#查询待处理余额)
- [单位转换](#单位转换)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解以太坊余额的存储单位（Wei）
- 使用 `BalanceAt` 方法查询地址余额
- 查询指定区块的余额
- 查询待处理的余额
- 将 Wei 转换为 ETH

## 前置条件

- Go 语言基础（变量、函数、错误处理）
- 已完成 [2.04 创建钱包](../2.04-create-wallet/) 模块
- 已安装 Go 环境（1.18+）
- 拥有以太坊节点访问地址

## 核心概念

### ETH 单位体系

以太坊使用 Wei 作为最小单位，不同单位之间的换算关系：

```
1 ETH     = 10^18 Wei
1 Gwei    = 10^9 Wei   (常用于 Gas 价格)
1 Szabo   = 10^12 Wei
1 Finney  = 10^15 Wei
```

```
数值示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1 ETH     = 1,000,000,000,000,000,000 Wei
0.5 ETH   =   500,000,000,000,000,000 Wei
1 Gwei    =               1,000,000,000 Wei
```

### 为什么使用 Wei？

```
原因: 以太坊虚拟机 (EVM) 不支持小数
     └─> 所有金额都用整数表示
     └─> Wei 是最小单位，避免浮点数精度问题

类比: 人民币用"分"作为最小单位
     1 元 = 100 分
     1.5 元 = 150 分（整数）
```

---

## 查询最新余额

### 方法：`BalanceAt`

```go
func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `ctx` | `context.Context` | 上下文，用于超时控制 |
| `account` | `common.Address` | 要查询的地址 |
| `blockNumber` | `*big.Int` | 区块号（`nil` 表示最新区块） |
| 返回 | `*big.Int` | 余额（单位：Wei） |

### 示例：查询地址余额

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common"
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

    // 要查询的地址
    account := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")

    // 查询最新余额（单位：Wei）
    balance, err := client.BalanceAt(context.Background(), account, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("余额: %s Wei\n", balance.String())
}
```

---

## 查询指定区块余额

传入区块号可以查询该区块时的账户余额，区块号必须是 `big.Int` 类型。

```go
import "math/big"

blockNumber := big.NewInt(5532993)
balance, err := client.BalanceAt(context.Background(), account, blockNumber)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("区块 %d 时的余额: %s Wei\n", blockNumber, balance)
```

---

## 查询待处理余额

有时候你想知道待处理的账户余额，例如在提交或等待交易确认后。

### 方法：`PendingBalanceAt`

```go
pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("待处理余额: %s Wei\n", pendingBalance)
```

---

## 单位转换

### Wei 转 ETH

由于余额返回的是 Wei（`*big.Int`），需要转换才能显示为 ETH：

```go
import (
    "math"
    "math/big"
)

// 方法：使用 big.Float 进行除法
fbalance := new(big.Float)
fbalance.SetString(balance.String())
ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

fmt.Printf("余额: %f ETH\n", ethValue)
```

### 完整示例：带单位转换的余额查询

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
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

    account := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")

    // 1. 查询最新余额（Wei）
    balance, err := client.BalanceAt(context.Background(), account, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("余额: %s Wei\n", balance.String())

    // 2. 查询指定区块余额
    blockNumber := big.NewInt(5532993)
    balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("区块 %d 时的余额: %s Wei\n", blockNumber, balanceAt)

    // 3. 转换为 ETH
    fbalance := new(big.Float)
    fbalance.SetString(balanceAt.String())
    ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
    fmt.Printf("余额: %.18f ETH\n", ethValue)

    // 4. 查询待处理余额
    pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("待处理余额: %s Wei\n", pendingBalance)
}
```

---

## 常见问题

### Q1: 为什么余额返回的是 `*big.Int` 而不是 `float64`？

**答：** 以太坊余额可能非常大（超过 `float64` 精度范围），且需要精确计算。`*big.Int` 是任意精度整数，可以精确表示任意大小的金额。

```go
// ❌ 错误：使用 float64 会丢失精度
balance := float64(1000000000000000000) // 精度损失

// ✅ 正确：使用 big.Int 保持精确
balance := big.NewInt(1000000000000000000) // 精确
```

### Q2: 查询余额需要消耗 Gas 吗？

**答：** 不需要。查询是**只读操作**，不修改区块链状态，所以免费：

| 操作类型 | 示例 | 消耗 Gas |
|----------|------|----------|
| 只读操作 | `BalanceAt`, `PendingBalanceAt` | ❌ 免费 |
| 写入操作 | `SendTransaction`, `Commit` | ✅ 需要 Gas |

### Q3: `nil` 作为 blockNumber 参数是什么意思？

**答：** `nil` 表示查询**最新区块**的余额：

```go
// 等价写法
client.BalanceAt(ctx, address, nil)           // 推荐
client.BalanceAt(ctx, address, big.NewInt(-1)) // 也表示最新
```

### Q4: 如何查询历史余额？

**答：** 指定过去的区块号：

```go
// 查询区块 #10000000 时的余额
oldBlock := big.NewInt(10000000)
balance, err := client.BalanceAt(context.Background(), address, oldBlock)
```

---

## 练习作业

开始练习前，请先进入目录并安装依赖：

```bash
cd ethclient-practice/2.07-query-balance
go mod tidy
```

### 作业 1：查询地址余额（基础）

练习文件：[exercises/01-query-balance.go](exercises/01-query-balance.go)

编写一个程序，实现以下功能：

1. 连接到 Sepolia 测试网
2. 查询指定地址的 ETH 余额
3. 将余额从 Wei 转换为 ETH
4. 输出结果：
   - 地址（带 0x 前缀）
   - 余额（Wei）
   - 余额（ETH，保留 18 位小数）

**运行练习：**
```bash
go run exercises/01-query-balance.go
```

**参考答案：** [solutions/01-query-balance.go](solutions/01-query-balance.go)

---

### 作业 2：历史余额查询（进阶）

练习文件：[exercises/02-historical-balance.go](exercises/02-historical-balance.go)

编写一个程序，实现以下功能：

1. 查询地址在当前区块和历史区块的余额
2. 对比两个时间点的余额变化
3. 输出：
   - 当前余额
   - 历史余额
   - 余额变化量

**运行练习：**
```bash
go run exercises/02-historical-balance.go
```

**参考答案：** [solutions/02-historical-balance.go](solutions/02-historical-balance.go)

---

### 作业 3：余额监控器（挑战）

练习文件：[exercises/03-balance-monitor.go](exercises/03-balance-monitor.go)

编写一个程序，实现以下功能：

1. 每隔 10 秒查询一次地址余额
2. 检测余额变化
3. 当余额变化时打印通知：
   ```
   [2024-01-22 12:34:56] 余额变化！
   旧余额: 1.000000 ETH
   新余额: 2.500000 ETH
   变化: +1.500000 ETH
   ```
4. 按 Ctrl+C 退出程序

**提示：**
- 使用 `time.Tick()` 定时查询
- 使用 `signal.Notify()` 捕获退出信号
- 比较两次查询的余额是否相同

**运行练习：**
```bash
go run exercises/03-balance-monitor.go
```

**参考答案：** [solutions/03-balance-monitor.go](solutions/03-balance-monitor.go)

---

## 下一步学习

- [ETH 转账](../2.05-transfer-eth/)
- [代币转账](../2.06-transfer-token/)
- [查询代币余额](../2.08-query-token-balance/)
