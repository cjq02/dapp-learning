# Ethclient 查询交易学习指南

> **预计学习时间：** 40 分钟
>
> **难度：** 基础

本指南介绍如何使用 Go 语言的 `go-ethereum` 库（ethclient）查询以太坊区块链上的交易信息。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [从区块获取交易](#从区块获取交易)
- [获取发送者地址](#获取发送者地址)
- [查询交易收据](#查询交易收据)
- [按索引查询交易](#按索引查询交易)
- [按哈希查询交易](#按哈希查询交易)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 从区块中遍历所有交易
- 获取交易的各种信息（哈希、金额、Gas、Nonce 等）
- 恢复交易发送者地址
- 查询交易收据和状态
- 按索引或哈希直接查询单个交易

## 前置条件

- Go 语言基础（变量、函数、错误处理）
- 以太坊基本概念（交易、哈希、Gas）
- 已完成 [2.01 查询区块](../2.01-query-block/) 章节
- 拥有以太坊节点访问地址（测试网或主网）

## 核心概念

### 交易结构

```
┌─────────────────────────────────────────┐
│            交易 (Transaction)            │
├─────────────────────────────────────────┤
│  基本信息                                │
│  - Hash: 交易哈希                        │
│  - Value: 转账金额（Wei）                │
│  - Gas: Gas 限制                         │
│  - GasPrice: Gas 价格                    │
│  - Nonce: 序列号                         │
├─────────────────────────────────────────┤
│  地址信息                                │
│  - To: 接收地址                          │
│  - From: 发送地址（需要恢复）            │
├─────────────────────────────────────────┤
│  数据                                    │
│  - Data: 附加数据（合约调用时使用）      │
│  - v, r, s: 签名信息                     │
└─────────────────────────────────────────┘
```

### 交易查询方式

| 方法 | 说明 | 使用场景 |
|------|------|----------|
| `block.Transactions()` | 获取区块所有交易 | 遍历区块交易 |
| `TransactionInBlock` | 按索引获取单笔交易 | 精确获取某笔交易 |
| `TransactionByHash` | 按哈希获取交易 | 已知交易哈希 |
| `TransactionReceipt` | 获取交易收据 | 查询交易状态 |

---

## 从区块获取交易

使用 `BlockByNumber` 获取区块后，调用 `Transactions()` 方法遍历交易：

```go
block, err := client.BlockByNumber(context.Background(), blockNumber)
if err != nil {
    log.Fatal(err)
}

for _, tx := range block.Transactions() {
    fmt.Printf("Hash: %s\n", tx.Hash().Hex())
    fmt.Printf("Value: %s Wei\n", tx.Value().String())
    fmt.Printf("Gas: %d\n", tx.Gas())
    fmt.Printf("Gas Price: %s Gwei\n", tx.GasPrice().String())
    fmt.Printf("Nonce: %d\n", tx.Nonce())
    fmt.Printf("To: %s\n", tx.To().Hex())
}
```

### Transaction 可用字段

```go
type Transaction struct {
    Hash     common.Hash      // 交易哈希
    Value    *big.Int         // 转账金额（Wei）
    Gas      uint64           // Gas 限制
    GasPrice *big.Int         // Gas 价格
    Nonce    uint64           // 序列号
    To       *common.Address  // 接收地址
    Data     []byte           // 附加数据
}
```

---

## 获取发送者地址

交易中不直接存储发送者地址，需要通过签名恢复。

### 方法：使用 `types.Sender`

```go
// 1. 获取链 ID
chainID, err := client.ChainID(context.Background())
if err != nil {
    log.Fatal(err)
}

// 2. 创建 EIP155 签名器
signer := types.NewEIP155Signer(chainID)

// 3. 恢复发送者地址
sender, err := types.Sender(signer, tx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("From: %s\n", sender.Hex())
```

### 方法：使用 `AsMessage`

```go
msg, err := tx.AsMessage(types.NewEIP155Signer(chainID))
if err != nil {
    log.Fatal(err)
}

fmt.Printf("From: %s\n", msg.From().Hex())
```

---

## 查询交易收据

交易收据包含执行结果、状态和日志：

```go
receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %d\n", receipt.Status)       // 1=成功, 0=失败
fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
fmt.Printf("Logs: %d\n", len(receipt.Logs))
```

### Receipt 可用字段

```go
type Receipt struct {
    Status    uint64        // 1=成功, 0=失败
    GasUsed   uint64        // 实际使用的 Gas
    Logs      []*Log        // 事件日志
    ContractAddress common.Address  // 创建的合约地址
}
```

---

## 按索引查询交易

使用 `TransactionInBlock` 在区块中按索引获取交易：

```go
blockHash := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")

// 获取交易数量
count, err := client.TransactionCount(context.Background(), blockHash)
if err != nil {
    log.Fatal(err)
}

// 遍历交易
for idx := uint(0); idx < count; idx++ {
    tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("[%d] %s\n", idx, tx.Hash().Hex())
}
```

---

## 按哈希查询交易

使用 `TransactionByHash` 直接查询单个交易：

```go
txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")

tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hash: %s\n", tx.Hash().Hex())
fmt.Printf("Pending: %v\n", isPending)  // true=在内存池中, false=已打包
```

---

## 常见问题

### Q1: 为什么需要 ChainID 来恢复发送者地址？

EIP-155 引入 ChainID 是为了重放攻击保护。不同链的交易即使其他参数相同，由于 ChainID 不同，签名也会不同。

### Q2: Gas 和 GasPrice 的区别？

| 字段 | 说明 |
|------|------|
| `Gas` | 用户愿意为交易支付的最大 Gas 单位 |
| `GasPrice` | 每单位 Gas 的价格（Wei） |
| 实际费用 | `GasUsed × GasPrice` |

### Q3: 如何判断交易是否成功？

检查 `receipt.Status`：
- `1` = 交易成功
- `0` = 交易失败

---

## 练习作业

开始练习前，请先安装依赖：

```bash
go mod tidy
```

### 作业 1：交易信息提取器（基础）

练习文件：[exercises/01-tx-info.go](exercises/01-tx-info.go)

给定一个区块号，提取并显示该区块中所有交易的基本信息。

**要求：**
1. 获取指定区块的所有交易
2. 显示每笔交易的：
   - 交易哈希
   - 发送者地址（恢复）
   - 接收者地址
   - 转账金额（转换为 Ether）
   - Gas 价格（转换为 Gwei）

**运行练习：**
```bash
go run exercises/01-tx-info.go
```

**参考答案：** [solutions/01-tx-info.go](solutions/01-tx-info.go)

---

### 作业 2：交易状态检查器（进阶）

练习文件：[exercises/02-tx-status.go](exercises/02-tx-status.go)

编写一个程序，检查指定交易的执行状态和详细信息。

**要求：**
1. 输入交易哈希
2. 显示：
   - 交易基本信息
   - 交易状态（成功/失败）
   - 实际使用的 Gas
   - 事件日志数量
   - 是否在内存池中

**运行练习：**
```bash
go run exercises/02-tx-status.go
```

**参考答案：** [solutions/02-tx-status.go](solutions/02-tx-status.go)

---

### 作业 3：交易分析器（挑战）

练习文件：[exercises/03-tx-analyzer.go](exercises/03-tx-analyzer.go)

编写一个程序，分析一个区块中的所有交易，输出统计信息。

**要求：**
1. 统计总交易数
2. 计算：
   - 总转账金额
   - 总 Gas 使用量
   - 平均 Gas 价格
   - 成功/失败交易数（需要查询收据）
3. 识别合约创建交易
4. 输出最贵的 3 笔交易

**运行练习：**
```bash
go run exercises/03-tx-analyzer.go
```

**参考答案：** [solutions/03-tx-analyzer.go](solutions/03-tx-analyzer.go)

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

- [查询账户余额](../2.07-query-balance/)
- [ETH 转账](../2.05-eth-transfer/)
- [监听新区块](../2.09-subscribe-blocks/)
