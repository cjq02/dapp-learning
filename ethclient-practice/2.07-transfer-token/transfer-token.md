# Ethclient ERC20 代币转账学习指南

> **预计学习时间：** 60 分钟
>
> **难度：** 中等

本指南介绍如何使用 Go 语言的 `go-ethereum` 库发送 ERC20 代币转账。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [ERC20 标准](#erc20-标准)
- [构建交易数据](#构建交易数据)
- [完整示例](#完整示例)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解 ERC20 代币标准
- 掌握如何构建合约调用数据
- 使用 Method ID 编码函数调用
- 发送 ERC20 代币转账
- 估算合约调用 Gas

## 前置条件

- Go 语言基础
- 已完成 [2.06 ETH 转账](../2.06-transfer-eth/) 模块
- 了解智能合约基础
- 了解十六进制和字节操作

## 核心概念

### 什么是 ERC20？

```
ERC20 = Ethereum Request for Comment 20
       └─> 以太坊代币标准
       └─> 定义了代币必须实现的接口
       └─> 所有 ERC20 代币遵循相同规则，可以互相兼容
```

### ERC20 vs ETH 转账

| 特性 | ETH 转账 | ERC20 转账 |
|------|----------|------------|
| 接收地址 | 直接转账给接收方 | 转账给**合约地址** |
| Value 字段 | 转账金额 | **0**（不转 ETH） |
| Data 字段 | 空 | **包含函数调用数据** |
| Gas Limit | 固定 21000 | 需要估算 |

```
ETH 转账流程:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
发送方 → 接收方
        (直接转账)


ERC20 转账流程:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
发送方 → 代币合约 → 调用 transfer() → 接收方收到代币
                 (在合约内部执行)
```

---

## ERC20 标准

### ERC20 必须实现的函数

```solidity
// 查询余额
function balanceOf(address account) external view returns (uint256)

// 转账
function transfer(address to, uint256 amount) external returns (bool)

// 授权
function approve(address spender, uint256 amount) external returns (bool)

// 查询授权额度
function allowance(address owner, address spender) external view returns (uint256)

// 代理转账
function transferFrom(address from, address to, uint256 amount) external returns (bool)
```

### Transfer 函数详解

```solidity
function transfer(address to, uint256 amount) external returns (bool)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `to` | address | 接收代币的地址 |
| `amount` | uint256 | 转账数量（以代币最小单位计） |
| 返回值 | bool | 转账是否成功 |

---

## 构建交易数据

### 数据字段结构

调用智能合约需要构建交易的数据字段：

```
Data 字段 = Method ID (4 字节) + 参数 (每个 32 字节)
           └─> Method ID = 函数签名的 Keccak256 哈希的前 4 字节
           └─> 参数需要左填充到 32 字节
```

### Transfer 函数的数据编码

```
transfer(address to, uint256 amount)
           ↓
┌────────────────────────────────────────────────────────┐
│ Method ID (4 bytes)                                     │
│ 0xa9059cbb                                             │
├────────────────────────────────────────────────────────┤
│ to (32 bytes, 左填充)                                  │
│ 0x0000000000000000000000004592...ac79d                 │
├────────────────────────────────────────────────────────┤
│ amount (32 bytes, 左填充)                              │
│ 0x00000000000000000000000000000...dea00000             │
└────────────────────────────────────────────────────────┘
```

### 步骤 1：生成 Method ID

```go
import "golang.org/x/crypto/sha3"

// 函数签名（函数名 + 参数类型，无空格）
transferFnSignature := []byte("transfer(address,uint256)")

// 计算 Keccak256 哈希
hash := sha3.NewLegacyKeccak256()
hash.Write(transferFnSignature)

// 取前 4 字节作为 Method ID
methodID := hash.Sum(nil)[:4]
fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))
// 输出: Method ID: 0xa9059cbb
```

**什么是函数签名？**

```
函数签名 = 函数名 + (参数类型1,参数类型2,...)
         └─> 不包含参数名称
         └─> 无空格

示例:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
transfer(address,uint256)
balanceOf(address)
approve(address,uint256)
transferFrom(address,address,uint256)
```

### 步骤 2：填充地址到 32 字节

```go
toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

// 左填充到 32 字节
paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
fmt.Printf("Padded Address: %s\n", hexutil.Encode(paddedAddress))
// 输出: 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d
```

**为什么要填充？**

```
EVM (以太坊虚拟机) 要求每个参数占 32 字节
         └─> 地址只有 20 字节
         └─> 需要在左边补 0 到 32 字节

未填充: 0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d  (20 字节)
填充后:  0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d  (32 字节)
         └────────补12字节零───┘        └────20字节地址────┘
```

### 步骤 3：填充金额到 32 字节

```go
// 转账金额（假设代币有 18 位小数，1000 个代币）
amount := new(big.Int)
amount.SetString("1000000000000000000000", 10) // 1000 * 10^18

// 左填充到 32 字节
paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
fmt.Printf("Padded Amount: %s\n", hexutil.Encode(paddedAmount))
// 输出: 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
```

### 步骤 4：组合数据

```go
var data []byte
data = append(data, methodID...)      // 4 字节
data = append(data, paddedAddress...) // 32 字节
data = append(data, paddedAmount...)  // 32 字节

// 总共 68 字节
fmt.Printf("Data: %s\n", hexutil.Encode(data))
// 输出: 0xa9059cbb0000000000000000000000004592...c5dea00000
```

### 数据可视化

```
完整 Data 字段 (68 字节 = 136 个十六进制字符):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

0xa9059cbb                                    ← Method ID (4 字节)
0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d  ← to (32 字节)
00000000000000000000000000000000000000000000003635c9adc5dea00000  ← amount (32 字节)
```

---

## 估算 Gas

ERC20 转账的 Gas 消耗不固定，需要估算：

```go
import "github.com/ethereum/go-ethereum"

gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
    To:   &tokenAddress,  // 注意：是代币合约地址
    From: fromAddress,
    Data: data,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Estimated Gas: %d\n", gasLimit)
// 输出: Estimated Gas: 51698
```

**为什么需要估算？**

```
ETH 转账:
    └─> 简单转账，固定消耗 21000 Gas

ERC20 转账:
    └─> 需要执行合约代码
    └─> 更新合约存储
    └─> 触发 Transfer 事件
    └─> Gas 消耗不定，需要估算
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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 加载私钥
	privateKey, err := crypto.HexToECDSA("your-private-key")
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

	// 设置转账金额（ERC20 转账不转 ETH）
	value := big.NewInt(0)

	// 获取 Gas Price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收地址（你要把代币转给谁）
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// 设置代币合约地址
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// === 构建 transfer 函数调用数据 ===

	// 1. 生成 Method ID
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

	// 2. 填充地址
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Printf("Padded Address: %s\n", hexutil.Encode(paddedAddress))

	// 3. 设置金额并填充（1000 个代币）
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("Padded Amount: %s\n", hexutil.Encode(paddedAmount))

	// 4. 组合数据
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// === 估算 Gas ===
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress, // 注意：是代币合约地址
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Estimated Gas: %d\n", gasLimit)

	// === 构建交易 ===
	// 注意：to 字段是代币合约地址，不是接收代币的地址
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

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

## 常见问题

### Q1: 为什么交易中 `to` 是合约地址而不是接收方地址？

**答：** ERC20 转账本质上是**调用代币合约的 `transfer` 函数**：

```
ETH 转账:
    tx.To = 接收方地址

ERC20 转账:
    tx.To = 代币合约地址
    tx.Data = transfer(接收方地址, 金额)
```

### Q2: 如何确定代币的小数位数？

**答：** 查询代币合约的 `decimals()` 函数：

```go
// decimals 函数签名
decimalsFnSignature := []byte("decimals()")

// 生成 Method ID
hash := sha3.NewLegacyKeccak256()
hash.Write(decimalsFnSignature)
methodID := hash.Sum(nil)[:4]

// 组合数据（无参数）
var data []byte
data = append(data, methodID...)

// 调用合约（不发送交易，只读取）
result, err := client.CallContract(context.Background(), ethereum.CallMsg{
    To:   &tokenAddress,
    Data: data,
}, nil)

// 解析结果（uint256 = 32 字节）
decimals := new(big.Int).SetBytes(result)
fmt.Printf("Decimals: %d\n", decimals.Uint64()) // 通常是 18
```

### Q3: 如何查询代币余额？

**答：** 调用代币合约的 `balanceOf(address)` 函数：

```go
// 构建调用数据
balanceOfSignature := []byte("balanceOf(address)")
hash := sha3.NewLegacyKeccak256()
hash.Write(balanceOfSignature)
methodID := hash.Sum(nil)[:4]

paddedAddress := common.LeftPadBytes(targetAddress.Bytes(), 32)

var data []byte
data = append(data, methodID...)
data = append(data, paddedAddress...)

// 调用合约
result, err := client.CallContract(context.Background(), ethereum.CallMsg{
    To:   &tokenAddress,
    Data: data,
}, nil)

// 解析结果
balance := new(big.Int).SetBytes(result)
fmt.Printf("Balance: %s\n", balance.String())
```

### Q4: 转账失败常见原因？

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| `insufficient balance` | 代币余额不足 | 检查余额是否足够 |
| `transfer amount exceeds allowance` | 未授权或授权额度不足 | 先调用 `approve()` |
| `execution reverted` | 合约执行失败 | 检查合约逻辑 |
| `gas required exceeds allowance` | Gas 不足 | 提高 Gas Limit |

### Q5: 常见 ERC20 代币的 Method ID

| 函数 | 签名 | Method ID |
|------|------|-----------|
| transfer | `transfer(address,uint256)` | `0xa9059cbb` |
| approve | `approve(address,uint256)` | `0x095ea7b3` |
| transferFrom | `transferFrom(address,address,uint256)` | `0x23b872dd` |
| balanceOf | `balanceOf(address)` | `0x70a08231` |
| allowance | `allowance(address,address)` | `0xdd62ed3e` |

---

## 练习作业

开始练习前，请先进入目录并安装依赖：

```bash
cd ethclient-practice/2.07-transfer-token
go mod tidy
```

⚠️ **重要提醒：**
- 你需要持有一些 ERC20 代币（如 USDT、USDC 或测试代币）
- 或者在 Sepolia 上部署测试代币合约

### 作业 1：发送 ERC20 代币转账（基础）

练习文件：[exercises/01-send-token.go](exercises/01-send-token.go)

编写一个程序，实现以下功能：

1. 从环境变量读取配置（私钥、代币合约地址、接收地址）
2. 构建 `transfer` 函数调用数据
3. 估算 Gas 并发送交易
4. 等待交易确认

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export PRIVATE_KEY=your-private-key
export TOKEN_ADDRESS=0x...  # 代币合约地址
export TO_ADDRESS=0x...     # 接收代币的地址
export TOKEN_AMOUNT=1000000000000000000000  # 转账数量（Wei 单位）
go run exercises/01-send-token.go
```

**参考答案：** [solutions/01-send-token.go](solutions/01-send-token.go)

---

### 作业 2：查询代币余额（进阶）

练习文件：[exercises/02-query-token-balance.go](exercises/02-query-token-balance.go)

编写一个程序，实现以下功能：

1. 调用 `balanceOf(address)` 查询代币余额
2. 调用 `decimals()` 获取代币小数位数
3. 将余额转换为人类可读格式
4. 输出结果

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export TOKEN_ADDRESS=0x...
export TARGET_ADDRESS=0x...
go run exercises/02-query-token-balance.go
```

**参考答案：** [solutions/02-query-token-balance.go](solutions/02-query-token-balance.go)

---

### 作业 3：授权并代理转账（挑战）

练习文件：[exercises/03-approve-transfer.go](exercises/03-approve-transfer.go)

编写一个程序，实现以下功能：

1. 调用 `approve(spender, amount)` 授权某个地址使用你的代币
2. 等待授权交易确认
3. 使用 `transferFrom(from, to, amount)` 进行代理转账
4. 验证转账前后余额变化

**提示：**
- `approve` 和 `transferFrom` 是两个独立的交易
- 需要等待第一笔交易确认后再发送第二笔

**运行练习：**
```bash
export INFURA_API_KEY=your-key
export PRIVATE_KEY=your-private-key
export TOKEN_ADDRESS=0x...
export SPENDER_ADDRESS=0x...  # 被授权的地址
export TO_ADDRESS=0x...
export AMOUNT=1000000000000000000000
go run exercises/03-approve-transfer.go
```

**参考答案：** [solutions/03-approve-transfer.go](solutions/03-approve-transfer.go)

---

## 测试代币合约

在 Sepolia 上部署以下测试代币合约，获取测试币：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestToken is ERC20 {
    uint256 public constant RATE = 100000000; // 1 ETH = 100M tokens

    constructor() ERC20("Test Token", "TST") {}

    function mint() public payable {
        require(msg.value >= 0.001 ether, "Min 0.001 ETH");
        uint256 tokensToMint = (msg.value * RATE);
        _mint(msg.sender, tokensToMint);
    }

    // 可以发送 0.001 ETH 到合约获取测试代币
}
```

---

## 下一步学习

- [智能合约部署](../2.10-deploy-contract/)
- [智能合约交互](../2.11-contract-interact/)
- [监听合约事件](../2.13-contract-events/)
