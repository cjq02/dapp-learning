# Ethclient 查询收据学习指南

> **预计学习时间：** 50 分钟
>
> **难度：** 中等

## 目录
- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [查询收据方法](#查询收据方法)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

## 学习目标
完成本指南后，你将能够：
- 理解交易收据的概念和作用
- 使用 `TransactionReceipt` 查询单笔交易收据
- 使用 `BlockReceipts` 批量查询区块收据
- 解析收据中的状态、日志、合约地址等信息

## 前置条件
- 完成 [2.01 查询区块](../2.01-query-block/query-block.md)
- 完成 [2.02 查询交易](../2.02-query-tx/query-tx.md)
- 了解交易的基本结构（Hash、Value、Gas 等）

## 核心概念

### 什么是交易收据？

交易收据（Transaction Receipt）是交易执行后生成的执行结果证明，包含以下关键信息：

```
┌─────────────────────────────────────────────────────────────┐
│                      Transaction Receipt                     │
├─────────────────────────────────────────────────────────────┤
│  Status              │  1 = 成功, 0 = 失败                   │
│  TxHash              │  交易哈希                             │
│  TransactionIndex    │  交易在区块中的索引                   │
│  BlockHash           │  所属区块哈希                         │
│  BlockNumber         │  所属区块号                           │
│  From                │  发送者地址                           │
│  To                  │  接收者地址（合约创建时为空）         │
│  ContractAddress     │  创建的合约地址（仅合约创建交易）     │
│  GasUsed             │  实际消耗的 Gas                       │
│  CumulativeGasUsed   │  累计 Gas 消耗                        │
│  Logs                │  事件日志数组                         │
│  LogsBloom           │  Bloom 过滤器（用于日志索引）         │
└─────────────────────────────────────────────────────────────┘
```

### 收据的作用

1. **交易状态确认**：通过 `Status` 判断交易是否成功执行
2. **Gas 费用计算**：`GasUsed * GasPrice` = 实际交易费用
3. **合约交互记录**：Logs 记录合约触发的事件
4. **合约创建证明**：`ContractAddress` 存储新创建的合约地址

### 查询方式对比

| 方法 | 用途 | 参数 | 返回值 |
|------|------|------|--------|
| `TransactionReceipt` | 查询单笔交易收据 | 交易哈希 | `*types.Receipt` |
| `BlockReceipts` | 查询区块所有收据 | 区块哈希/高度 | `[]*types.Receipt` |

## 查询收据方法

### 1. 查询单笔交易收据

使用 `TransactionReceipt` 方法通过交易哈希查询收据：

```go
txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")
receipt, err := client.TransactionReceipt(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("状态: %d\n", receipt.Status)                // 1 = 成功
fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
fmt.Printf("交易索引: %d\n", receipt.TransactionIndex)
```

### 2. 批量查询区块收据

使用 `BlockReceipts` 方法查询指定区块的所有收据：

```go
import "github.com/ethereum/go-ethereum/rpc"

// 方式 1: 通过区块高度
blockNumber := big.NewInt(5671744)
receipts, err := client.BlockReceipts(
    context.Background(),
    rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())),
)

// 方式 2: 通过区块哈希
blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
receipts, err := client.BlockReceipts(
    context.Background(),
    rpc.BlockNumberOrHashWithHash(blockHash, false),
)

for _, receipt := range receipts {
    fmt.Printf("交易: %s, 状态: %d\n", receipt.TxHash.Hex(), receipt.Status)
}
```

### 3. 解析收据信息

```go
// 判断交易状态
if receipt.Status == 1 {
    fmt.Println("交易成功")
} else {
    fmt.Println("交易失败")
}

// 检查是否为合约创建
if receipt.ContractAddress != (common.Address{}) {
    fmt.Printf("创建的合约地址: %s\n", receipt.ContractAddress.Hex())
}

// 遍历事件日志
for _, log := range receipt.Logs {
    fmt.Printf("合约: %s, 事件: %s\n", log.Address.Hex(), log.Topics[0].Hex())
}
```

### 4. 计算实际交易费用

```go
// 获取交易（需要 Gas Price）
tx, _, err := client.TransactionByHash(context.Background(), txHash)

// 计算实际费用
actualFee := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GasPrice())

// 转换为 Ether
feeInEther := new(big.Float).Quo(
    new(big.Float).SetInt(actualFee),
    big.NewFloat(1e18),
)
fmt.Printf("实际交易费用: %.6f Ether\n", feeInEther)
```

## 常见问题

### Q1: 收据查询返回 "not found"？
**A:** 可能原因：
- 交易还在 pending 状态，尚未被打包进区块
- 交易哈希输入错误
- 连接的网络与交易所在网络不一致

### Q2: Status 为 0 表示什么？
**A:** `Status = 0` 表示交易执行失败，常见原因：
- Gas 不足（out of gas）
- revert 触发
- require/assert 条件不满足

### Q3: ContractAddress 什么时候有值？
**A:** 仅当交易是合约创建交易（`tx.To() == nil`）且创建成功时，`ContractAddress` 才有值。

### Q4: Logs 是什么？
**A:** Logs 是合约执行时触发的事件记录，包含：
- `Address`: 触发事件的合约地址
- `Topics`: 事件签名和索引参数
- `Data`: 非索引参数数据

## 练习作业

### 作业 1：查询交易收据（基础）
练习文件：[exercises/01-receipt-basic.go](exercises/01-receipt-basic.go)

**任务：** 给定交易哈希，查询其收据并显示基本信息

**要求：**
- 使用 `TransactionReceipt` 查询收据
- 显示交易状态、Gas 使用量、交易索引
- 判断是否为合约创建交易

**参考答案：** [solutions/01-receipt-basic.go](solutions/01-receipt-basic.go)

### 作业 2：解析事件日志（进阶）
练习文件：[exercises/02-receipt-logs.go](exercises/02-receipt-logs.go)

**任务：** 查询收据并解析其中的事件日志

**要求：**
- 遍历 `receipt.Logs` 数组
- 显示每个日志的合约地址、主题数量、数据长度
- 统计总共有多少条日志

**参考答案：** [solutions/02-receipt-logs.go](solutions/02-receipt-logs.go)

### 作业 3：批量收据分析（挑战）
练习文件：[exercises/03-receipt-batch.go](exercises/03-receipt-batch.go)

**任务：** 查询指定区块的所有收据并进行统计分析

**要求：**
- 使用 `BlockReceipts` 批量查询
- 统计成功/失败交易数量
- 计算总 Gas 使用量和平均 Gas 使用
- 找出 Gas 使用最多的交易
- 统计创建了多少个新合约

**参考答案：** [solutions/03-receipt-batch.go](solutions/03-receipt-batch.go)

---

**下一步学习：** [2.04 创建新钱包](../2.04-create-wallet/create-wallet.md)
