# Ethclient 查询账户余额学习指南

本指南介绍如何使用 Go 语言的 `go-ethereum` 库查询以太坊地址的 ETH 余额。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [查询余额](#查询余额)
- [单位转换](#单位转换)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解以太坊余额的存储单位（Wei）
- 使用 `BalanceAt` 方法查询指定区块的余额
- 使用 `PendingNonceAt` 查询账户交易序列号
- 将 Wei 转换为 ETH 等更易读的单位
- 处理大数运算（`math/big`）

## 前置条件

- Go 语言基础（变量、函数、错误处理）
- 已完成 [2.04 创建钱包](../2.04-create-wallet/) 模块
- 了解以太坊账户系统
- 已安装 Go 环境（1.18+）

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

### 账户状态

以太坊地址在以下情况余额为 0：
- 新创建的地址
- 从未接收过转账的地址
- 余额已全部转出的地址

```
地址状态示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
0x1234...0001  → 1.5 ETH   (有余额)
0xabcd...ffff  → 0 ETH    (无余额/新地址)
```

---

## 查询余额

### 方法：`BalanceAt`

```go
func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `ctx` | `context.Context` | 上下文，用于超时控制 |
| `account` | `common.Address` | 要查询的地址 |
| `blockNumber` | `*big.Int` | 区块号（nil 表示最新区块） |
| 返回 | `*big.Int` | 余额（单位：Wei） |

### 区块号参数

```go
// 查询最新区块余额
balance, err := client.BalanceAt(context.Background(), address, nil)

// 查询指定区块余额
blockNumber := big.NewInt(12345678)
balance, err := client.BalanceAt(context.Background(), address, blockNumber)

// 使用常量查询最新区块
balance, err := client.BalanceAt(context.Background(), address, big.NewInt(-1))
```

### 示例：查询地址余额

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 要查询的地址
	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// 查询最新余额（单位：Wei）
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("余额: %s Wei\n", balance.String()) // 余额: 1000000000000000000 Wei
}
```

---

## 单位转换

### Wei 转 ETH

由于余额返回的是 Wei（`*big.Int`），需要转换才能显示为 ETH：

```go
// 方法：使用 big.Float 进行除法
func weiToEth(wei *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18), // 1 ETH = 10^18 Wei
	)
}

// 使用示例
balanceWei := big.NewInt(1000000000000000000) // 1 ETH
balanceEth := weiToEth(balanceWei)
fmt.Printf("余额: %f ETH\n", balanceEth) // 余额: 1.000000 ETH
```

### 完整示例：带单位转换的余额查询

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// 查询余额（Wei）
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 转换为 ETH
	balanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(1e18),
	)

	fmt.Printf("地址: %s\n", address.Hex())
	fmt.Printf("余额: %s Wei\n", balance.String())
	fmt.Printf("余额: %.6f ETH\n", balanceEth)
	// 输出:
	// 地址: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
	// 余额: 1000000000000000000 Wei
	// 余额: 1.000000 ETH
}
```

### 格式化输出

```go
// 根据余额大小自动选择单位
func formatBalance(wei *big.Int) string {
	if wei.Cmp(big.NewInt(1e18)) >= 0 {
		// >= 1 ETH，显示 ETH
		eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
		return fmt.Sprintf("%.4f ETH", eth)
	} else if wei.Cmp(big.NewInt(1e9)) >= 0 {
		// >= 1 Gwei，显示 Gwei
		gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e9))
		return fmt.Sprintf("%.2f Gwei", gwei)
	} else {
		// < 1 Gwei，显示 Wei
		return fmt.Sprintf("%s Wei", wei.String())
	}
}
```

---

## 查询交易序列号（Nonce）

### 什么是 Nonce？

```
Nonce (Number Used Once)
    └─> 账户的交易计数器
    └─> 每笔交易都需要唯一的 Nonce
    └─> 防止重放攻击
```

```
账户 Nonce 示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
新账户        → Nonce = 0
发送第 1 笔交易 → Nonce = 0
发送第 2 笔交易 → Nonce = 1
发送第 3 笔交易 → Nonce = 2
...
```

### 方法：`PendingNonceAt`

```go
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
```

返回该账户的**下一个** Nonce 值。

### 示例：查询账户 Nonce

```go
nonce, err := client.PendingNonceAt(context.Background(), address)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("当前 Nonce: %d\n", nonce)
// 下笔交易应该使用这个 Nonce
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
| 只读操作 | `BalanceAt`, `Call` | ❌ 免费 |
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
cd ethclient-practice/2.05-query-balance
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
   - 余额（ETH，保留 6 位小数）

**运行练习：**
```bash
export INFURA_API_KEY=your-key-here
go run exercises/01-query-balance.go
```

**参考答案：** [solutions/01-query-balance.go](solutions/01-query-balance.go)

---

### 作业 2：批量查询余额（进阶）

练习文件：[exercises/02-batch-query.go](exercises/02-batch-query.go)

编写一个程序，实现以下功能：

1. 定义多个地址（至少 3 个）
2. 批量查询每个地址的余额
3. 统计总余额
4. 输出格式化的表格：
   ```
   地址                           余额 (ETH)
   ──────────────────────────────────────────
   0x1234...abcd                    1.500000
   0x5678...efgh                    0.000000
   0x9abc...def0                    2.750000
   ──────────────────────────────────────────
   总计:                            4.250000 ETH
   ```

**运行练习：**
```bash
go run exercises/02-batch-query.go
```

**参考答案：** [solutions/02-batch-query.go](solutions/02-batch-query.go)

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

- [ETH 转账](../2.06-transfer-eth/)
- [签名交易](../2.07-sign-transaction/)
- [ERC20 代币查询](../2.08-query-token-balance/)
