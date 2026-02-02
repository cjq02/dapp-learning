# Ethclient 查询代币余额学习指南

> **预计学习时间：** 40 分钟
>
> **难度：** 中等

本指南介绍如何使用 Go 语言的 `go-ethereum` 库（ethclient）查询 ERC20 代币余额。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [ERC20 接口](#erc20-接口)
- [生成 Go 绑定代码](#生成-go-绑定代码)
- [查询代币余额](#查询代币余额)
- [查询代币信息](#查询代币信息)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解 ERC20 代币标准
- 使用 `abigen` 生成合约 Go 绑定代码
- 查询 ERC20 代币余额
- 查询代币元数据（名称、符号、精度）
- 处理不同精度的代币余额转换

## 前置条件

- Go 语言基础
- 已完成 [2.07 查询账户余额](../2.07-query-balance/) 模块
- 了解 Solidity 基础
- 已安装 Go 环境（1.18+）
- 拥有以太坊节点访问地址
- 已安装 `abigen` 工具（go-ethereum 自带）

## 核心概念

### 什么是 ERC20？

```
ERC20 (Ethereum Request for Comment 20)
    └─> 以太坊上代币的技术标准
    └─> 确保不同代币在以太坊网络上的互操作性
    └─> 定义了代币必须实现的一组接口
```

### ERC20 核心功能

```
┌─────────────────────────────────────┐
│            ERC20 接口                │
├─────────────────────────────────────┤
│  totalSupply()      → 代币总供应量    │
│  balanceOf(address) → 地址余额        │
│  transfer(to, amount) → 转账         │
│  approve(spender, amount) → 授权     │
│  allowance(owner, spender) → 授权额度│
├─────────────────────────────────────┤
│  元数据 (可选)                       │
│  name()        → 代币名称             │
│  symbol()      → 代币符号             │
│  decimals()    → 小数位数             │
└─────────────────────────────────────┘
```

### 代币精度（Decimals）

```
代币精度决定了代币的最小单位：

decimals = 18  (最常见)
    1 Token = 10^18 SmallestUnit
    1.5 Token = 1500000000000000000

decimals = 6  (USDC 等稳定币常用)
    1 Token = 10^6 SmallestUnit
    1.5 Token = 1500000
```

---

## ERC20 接口

### IERC20.sol

这是 OpenZeppelin 提供的标准 ERC20 接口：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 value) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 value) external returns (bool);
    function transferFrom(address from, address to, uint256 value) external returns (bool);
}
```

### IERC20Metadata.sol

扩展接口，包含代币元数据：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IERC20} from "./IERC20.sol";

interface IERC20Metadata is IERC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}
```

---

## 生成 Go 绑定代码

### 1. 创建 Solidity 接口文件

创建 `IERC20Metadata.sol` 文件（内容见上文）。

### 2. 编译生成 ABI

```bash
# 使用 solc 编译，生成 ABI 文件
solc --abi IERC20Metadata.sol -o build
```

### 3. 使用 abigen 生成 Go 代码

```bash
# 从 ABI 文件生成 Go 代码
abigen --abi=build/IERC20Metadata.abi --pkg=erc20 --out=erc20.go

# 或者直接从 Solidity 文件生成
abigen --sol=IERC20Metadata.sol --pkg=erc20 --out=erc20.go
```

### 4. 项目结构

```
2.08-query-token-balance/
├── go.mod
├── query-token-balance.md
├── contracts/
│   └── IERC20Metadata.sol
├── erc20/
│   └── erc20.go              # 生成的 Go 绑定代码
├── exercises/
│   ├── 01-query-token-balance.go
│   ├── 02-query-token-info.go
│   └── 03-batch-query.go
└── solutions/
    ├── 01-query-token-balance.go
    ├── 02-query-token-info.go
    └── 03-batch-query.go
```

---

## 查询代币余额

### 步骤 1：导入生成的包

```go
import (
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/dapp-learning/ethclient/query-token-balance/erc20"
)
```

### 步骤 2：实例化合约

```go
// ERC20 代币合约地址（以 USDC 为例）
tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

// 创建合约实例
instance, err := erc20.NewIERC20Metadata(tokenAddress, client)
if err != nil {
    log.Fatal(err)
}
```

### 步骤 3：查询余额

```go
// 要查询的地址
address := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")

// 调用 balanceOf 方法（使用 nil 作为 CallOpts）
bal, err := instance.BalanceOf(nil, address)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("余额: %s (最小单位)\n", bal) // 例如: 74605500647408739782407023
```

---

## 查询代币信息

### 查询代币元数据

```go
// 查询代币名称
name, err := instance.Name(nil)
if err != nil {
    log.Fatal(err)
}

// 查询代币符号
symbol, err := instance.Symbol(nil)
if err != nil {
    log.Fatal(err)
}

// 查询代币精度
decimals, err := instance.Decimals(nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("名称: %s\n", name)         // 例如: "USD Coin"
fmt.Printf("符号: %s\n", symbol)       // 例如: "USDC"
fmt.Printf("精度: %d\n", decimals)     // 例如: 6
```

### 查询总供应量

```go
totalSupply, err := instance.TotalSupply(nil)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("总供应量: %s\n", totalSupply)
```

---

## 单位转换

### 将余额转换为可读格式

```go
import (
    "math"
    "math/big"
)

// 获取余额
bal, err := instance.BalanceOf(nil, address)
if err != nil {
    log.Fatal(err)
}

// 获取精度
decimals, err := instance.Decimals(nil)
if err != nil {
    log.Fatal(err)
}

// 转换为可读格式
fbal := new(big.Float)
fbal.SetString(bal.String())
value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

fmt.Printf("余额: %s %s\n", value.String(), symbol)
```

### 完整示例

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
    erc20 "github.com/dapp-learning/ethclient/query-token-balance/erc20"
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

    // 代币合约地址
    tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

    // 创建合约实例
    instance, err := erc20.NewIERC20Metadata(tokenAddress, client)
    if err != nil {
        log.Fatal(err)
    }

    // 要查询的地址
    address := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")

    // 查询余额
    bal, err := instance.BalanceOf(nil, address)
    if err != nil {
        log.Fatal(err)
    }

    // 查询代币信息
    name, _ := instance.Name(nil)
    symbol, _ := instance.Symbol(nil)
    decimals, _ := instance.Decimals(nil)

    // 转换单位
    fbal := new(big.Float)
    fbal.SetString(bal.String())
    value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

    fmt.Printf("代币: %s (%s)\n", name, symbol)
    fmt.Printf("精度: %d\n", decimals)
    fmt.Printf("余额: %s %s\n", value.String(), symbol)
}
```

---

## 常见问题

### Q1: 如何获取代币合约地址？

**答：** 可以通过以下方式获取：

1. **Etherscan**：在 Etherscan 搜索代币名称，查看合约地址
2. **CoinGecko/CoinMarketCap**：查看代币详情页面
3. **项目官网**：查看官方文档

### Q2: `abigen` 工具在哪里？

**答：** `abigen` 是 `go-ethereum` 的一部分，安装方式：

```bash
# 安装 go-ethereum（包含 abigen）
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 验证安装
abigen --version
```

### Q3: 为什么使用 `nil` 作为 CallOpts？

**答：** `nil` 表示使用默认的调用选项：

```go
// 这两种写法等价
instance.BalanceOf(nil, address)
instance.BalanceOf(&bind.CallOpts{}, address)

// 自定义选项
opts := &bind.CallOpts{
    Pending: false,                 // 是否使用 pending 状态
    Context: context.Background(),   // 上下文
    BlockNumber: big.NewInt(12345), // 指定区块号
}
instance.BalanceOf(opts, address)
```

### Q4: 如何处理精度不是 18 的代币？

**答：** 先查询 `decimals()`，再进行转换：

```go
decimals, _ := instance.Decimals(nil)

// 使用实际的精度进行转换
value := new(big.Float).Quo(
    new(big.Float).SetInt(balance),
    big.NewFloat(math.Pow10(int(decimals))),
)
```

### Q5: 查询代币余额需要 Gas 吗？

**答：** 不需要。调用 `balanceOf`、`name`、`symbol` 等 view 函数都是**只读操作**，不消耗 Gas。

---

## 练习作业

开始练习前，请先准备：

```bash
# 1. 生成 Go 绑定代码（如果还没有）
abigen --sol=contracts/IERC20Metadata.sol --pkg=erc20 --out=erc20/erc20.go

# 2. 安装依赖
go mod tidy
```

### 作业 1：查询代币余额（基础）

练习文件：[exercises/01-query-token-balance.go](exercises/01-query-token-balance.go)

编写一个程序，实现以下功能：

1. 连接到测试网
2. 实例化 ERC20 合约
3. 查询指定地址的代币余额
4. 将余额转换为可读格式
5. 输出：
   - 代币名称
   - 代币符号
   - 余额（可读格式）

**测试代币地址：**
```
Sepolia 测试网 USDC: 0x94a9D9AC8a22534E3FaCa9F4e7F2E2cf85d5E4C8
```

**运行练习：**
```bash
go run exercises/01-query-token-balance.go
```

**参考答案：** [solutions/01-query-token-balance.go](solutions/01-query-token-balance.go)

---

### 作业 2：代币信息查询器（进阶）

练习文件：[exercises/02-query-token-info.go](exercises/02-query-token-info.go)

编写一个程序，实现以下功能：

1. 查询并显示完整的代币信息：
   - 名称
   - 符号
   - 精度
   - 总供应量
2. 查询多个地址的余额
3. 以表格形式输出

**运行练习：**
```bash
go run exercises/02-query-token-info.go
```

**参考答案：** [solutions/02-query-token-info.go](solutions/02-query-token-info.go)

---

### 作业 3：多代币余额查询（挑战）

练习文件：[exercises/03-batch-query.go](exercises/03-batch-query.go)

编写一个程序，实现以下功能：

1. 支持查询多个代币的余额
2. 输出格式化的表格：
   ```
   地址                           USDC       USDT       DAI
   ─────────────────────────────────────────────────────────
   0x1234...abcd                  100.50     0.00       50.25
   0x5678...efgh                  0.00       500.00     0.00
   ```
3. 计算总资产价值（假设每个代币价格）

**提示：**
- 使用 map 存储代币合约地址
- 并发查询提高效率

**运行练习：**
```bash
go run exercises/03-batch-query.go
```

**参考答案：** [solutions/03-batch-query.go](solutions/03-batch-query.go)

---

## 下一步学习

- [代币转账](../2.06-transfer-token/)
- [订阅新区块](../2.09-subscribe-block/)
- [部署合约](../2.10-deploy-contract/)
