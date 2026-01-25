# Ethclient ETH 转账学习指南

> **预计学习时间：** 45 分钟
>
> **难度：** 中等

本指南介绍如何使用 Go 语言的 `go-ethereum` 库发起 ETH 转账交易。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [交易结构](#交易结构)
- [发送交易流程](#发送交易流程)
- [完整示例](#完整示例)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解以太坊交易的各个组成部分
- 使用私钥对交易进行签名
- 构建并发送 ETH 转账交易
- 查询交易状态和收据
- 处理 Gas 相关参数

## 前置条件

- Go 语言基础
- 已完成 [2.04 创建钱包](../2.04-create-wallet/) 模块
- 已完成 [2.05 查询余额](../2.05-query-balance/) 模块
- 了解 Nonce、Gas、Gas Price 等概念
- 拥有 Sepolia 测试币

## 核心概念

### 什么是交易？

```
交易 = 发送方用私钥签名的数据包
       └─> 包含：转账金额、接收地址、Gas 等信息
       └─> 提交到区块链后，矿工验证并打包进区块
```

### 交易的生命周期

```
1. 构建交易
   ↓
2. 用私钥签名
   ↓
3. 广播到网络
   ↓
4. 矿工验证并打包
   ↓
5. 交易确认（写入区块）
   ↓
6. 可以查询交易收据
```

### 交易的组成部分

```
Transaction
├── Nonce        (交易序号)
├── To           (接收地址)
├── Value        (转账金额，单位：Wei)
├── GasLimit     (Gas 上限)
├── GasPrice     (Gas 价格)
└── Data         (附加数据，ETH 转账为空)
```

---

## 交易结构

### Nonce（交易序号）

```
Nonce = Number Used Once（仅使用一次的数字）
       └─> 每个账户都有一个计数器
       └─> 每发送一笔交易，Nonce +1
       └─> 防止重放攻击
```

```
账户 Nonce 示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
新账户     → Nonce = 0
发送第1笔  → Nonce = 0
发送第2笔  → Nonce = 1
发送第3笔  → Nonce = 2
...
```

**查询 Nonce：**
```go
nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
```

### Gas 相关参数

#### Gas Limit（Gas 上限）

```
Gas Limit = 愿意为这笔交易支付的最大 Gas 单位数
           └─> ETH 转账：固定 21000
           └─> 合约调用：根据合约逻辑而定
```

```
常见 Gas Limit:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ETH 转账      → 21,000
ERC20 转账    → ~50,000
合约部署      → ~2,000,000
```

#### Gas Price（Gas 价格）

```
Gas Price = 每单位 Gas 的价格（单位：Wei）
           └─> 越高，交易被打包越快
           └─> 可以根据网络状况调整
```

```
Gas Price 示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
30 Gwei  = 30,000,000,000 Wei  (正常)
100 Gwei = 100,000,000,000 Wei (拥堵时)
```

#### 交易费用

```
交易费用 = Gas Used × Gas Price
          └─> ETH 转账：21000 × GasPrice
          └─> 实际费用 ≤ GasLimit × GasPrice
```

```
计算示例 (GasPrice = 30 Gwei):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
费用 = 21000 × 30000000000
    = 630000000000000 Wei
    = 0.00063 ETH
```

### Value（转账金额）

```go
// 1 ETH = 10^18 Wei
value := big.NewInt(1000000000000000000) // 1 ETH

// 0.5 ETH
value := new(big.Int).Mul(big.NewInt(5), big.NewInt(1e17))

// 使用字符串构造大数
value, _ := new(big.Int).SetString("1000000000000000000", 10)
```

---

## 发送交易流程

### 步骤 1：加载私钥

```go
privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
if err != nil {
    log.Fatal(err)
}
```

### 步骤 2：获取发送方地址和 Nonce

```go
publicKey := privateKey.Public()
publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
if !ok {
    log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
}

fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

// 获取 Nonce
nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
if err != nil {
    log.Fatal(err)
}
```

### 步骤 3：设置转账金额

```go
// 转账 1 ETH
value := big.NewInt(1000000000000000000) // in wei (1 eth)
```

### 步骤 4：设置 Gas 参数

```go
// ETH 转账固定 Gas Limit
gasLimit := uint64(21000) // in units

// 获取建议的 Gas Price
gasPrice, err := client.SuggestGasPrice(context.Background())
if err != nil {
    log.Fatal(err)
}

// 或手动设置
gasPrice := big.NewInt(30000000000) // 30 gwei
```

### 步骤 5：设置接收地址

```go
toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
```

### 步骤 6：构建未签名交易

```go
// 方法签名
func NewTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction

// ETH 转账（data 为 nil）
tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
```

### 步骤 7：获取 Chain ID 并签名

```go
// 获取网络 ID
chainID, err := client.NetworkID(context.Background())
if err != nil {
    log.Fatal(err)
}

// 签名交易
signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
if err != nil {
    log.Fatal(err)
}
```

**什么是 EIP-155？**

```
EIP-155 = 以太坊改进提案 155
         └─> 在签名中引入 chain ID
         └─> 防止跨链重放攻击
         └─> 主网交易不能在测试网重放
```

### 步骤 8：发送交易

```go
err = client.SendTransaction(context.Background(), signedTx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("交易已发送: %s\n", signedTx.Hash().Hex())
```

---

## 完整示例

```go
package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 加载私钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	// 获取发送方地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取 Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Nonce: %d\n", nonce)

	// 设置转账金额
	value := big.NewInt(1000000000000000000) // 1 ETH in wei

	// 设置 Gas 参数
	gasLimit := uint64(21000)                // ETH 转账固定值
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Gas Price: %s Wei\n", gasPrice.String())

	// 设置接收地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// 构建交易
	var data []byte // ETH 转账不需要数据
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 获取 Chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已发送: %s\n", signedTx.Hash().Hex())
}
```

---

## 查询交易状态

### 等待交易确认

```go
// 交易哈希
txHash := common.HexToHash("0x...")

// 等待交易被打包（超时 5 分钟）
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

receipt, err := bind.WaitMined(ctx, client, signedTx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("交易已确认，区块: %d\n", receipt.BlockNumber.Uint64())
```

### 查询交易收据

```go
receipt, err := client.TransactionReceipt(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("状态: %v\n", receipt.Status) // 1 = 成功, 0 = 失败
fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
```

---

## 常见问题

### Q1: 交易发送失败怎么办？

**常见原因：**
| 错误 | 原因 | 解决方案 |
|------|------|----------|
| `insufficient funds` | 余额不足 | 确保账户有足够 ETH |
| `nonce too low` | Nonce 太小 | 使用 `PendingNonceAt` 获取最新值 |
| `replacement transaction underpriced` | Gas Price 太低 | 提高 Gas Price |
| `exceeds block gas limit` | Gas Limit 太高 | 检查 Gas Limit 设置 |

### Q2: 如何计算实际交易费用？

```go
receipt, _ := client.TransactionReceipt(context.Background(), txHash)

// 实际费用 = GasUsed × GasPrice
actualFee := new(big.Int).Mul(receipt.GasUsed, gasPrice)
actualFeeEth := new(big.Float).Quo(
    new(big.Float).SetInt(actualFee),
    big.NewFloat(1e18),
)
fmt.Printf("实际费用: %.6f ETH\n", actualFeeEth)
```

### Q3: Nonce 管理不当会怎样？

```
问题场景:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Nonce 跳号    → 交易排队等待
Nonce 重复    → "replacement transaction" 错误
Nonce 太小    → "nonce too low" 错误
```

**最佳实践：** 每次发送交易前都从链上获取最新 Nonce。

### Q4: 为什么需要 Chain ID？

```
Chain ID 的作用:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
主网   → Chain ID = 1
Sepolia → Chain ID = 11155111

└─> 防止在一条链上签名的交易在另一条链上重放
└─> 保护用户资产安全
```

### Q5: 如何加速交易？

```go
// 使用相同的 Nonce 和更高的 Gas Price 发送新交易
acceleratedTx := types.NewTransaction(
    nonce,
    toAddress,
    value,
    gasLimit,
    new(big.Int).Mul(gasPrice, big.NewInt(110)), // 提高 10%
    nil,
)
```

---

## 练习作业

开始练习前，请先进入目录并安装依赖：

```bash
cd ethclient-practice/2.06-transfer-eth
go mod tidy
```

⚠️ **重要提醒：** 这些练习需要**真实的 Sepolia 测试币**。请确保你的钱包有足够的测试 ETH。

### 作业 1：发送 ETH 转账（基础）

练习文件：[exercises/01-send-eth.go](exercises/01-send-eth.go)

编写一个程序，实现以下功能：

1. 从环境变量读取私钥
2. 连接到 Sepolia 测试网
3. 获取发送方地址和 Nonce
4. 构建并发送 0.001 ETH 的转账交易
5. 等待交易确认
6. 输出交易哈希和收据信息

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export PRIVATE_KEY=your-private-key
export TO_ADDRESS=recipient-address
go run exercises/01-send-eth.go
```

**参考答案：** [solutions/01-send-eth.go](solutions/01-send-eth.go)

---

### 作业 2：批量转账（进阶）

练习文件：[exercises/02-batch-transfer.go](exercises/02-batch-transfer.go)

编写一个程序，实现以下功能：

1. 定义多个接收地址和转账金额
2. 按顺序依次发送多笔转账
3. 每笔交易使用递增的 Nonce
4. 输出每笔交易的哈希和状态
5. 统计总转账金额和总 Gas 费用

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export PRIVATE_KEY=your-private-key
go run exercises/02-batch-transfer.go
```

**参考答案：** [solutions/02-batch-transfer.go](solutions/02-batch-transfer.go)

---

### 作业 3：交易监控器（挑战）

练习文件：[exercises/03-tx-monitor.go](exercises/03-tx-monitor.go)

编写一个程序，实现以下功能：

1. 发送一笔交易
2. 实时监控交易状态
3. 显示以下信息：
   - 交易是否在 mempool 中
   - 交易是否被打包
   - 交易是否成功
   - 实际 Gas 消耗
4. 当交易确认后显示详细信息

**提示：**
- 使用 `client.TransactionInPool()` 检查交易是否在 mempool
- 使用 `client.TransactionReceipt()` 查询交易收据
- 轮询查询直到交易确认

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export PRIVATE_KEY=your-private-key
export TO_ADDRESS=recipient-address
go run exercises/03-tx-monitor.go
```

**参考答案：** [solutions/03-tx-monitor.go](solutions/03-tx-monitor.go)

---

## 安全提醒

⚠️ **安全注意事项：**

1. **私钥安全**
   - ❌ 不要把私钥硬编码在代码中
   - ❌ 不要提交私钥到 Git
   - ✅ 使用环境变量或密钥管理服务

2. **测试网 vs 主网**
   - ✅ 在 Sepolia 测试网练习
   - ❌ 不要在主网测试代码

3. **交易验证**
   - ✅ 先用小额测试
   - ✅ 验证接收地址正确
   - ✅ 检查余额和 Gas 费用

---

## 下一步学习

- [ERC20 代币转账](../2.07-transfer-token/)
- [智能合约交互](../2.10-contract-interact/)
- [事件监听](../2.13-contract-events/)
