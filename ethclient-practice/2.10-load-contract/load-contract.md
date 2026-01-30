# Ethclient 加载合约学习指南

> **预计学习时间：** 30 分钟
>
> **难度：** 中等

本指南介绍如何使用 Go 语言的 `go-ethereum` 库加载已部署的智能合约并与合约进行交互。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [生成 Go 绑定代码](#生成-go-绑定代码)
- [加载合约实例](#加载合约实例)
- [Remix 加载合约](#remix-加载合约)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解合约加载的概念和原理
- 使用 `abigen` 仅根据 ABI 文件生成 Go 绑定代码
- 使用 ethclient 加载已部署的合约实例
- 使用 Remix 在线 IDE 加载已部署的合约
- 验证合约加载是否成功

## 前置条件

- Go 语言基础
- Solidity 基础（了解合约结构、函数、事件）
- 已完成 [部署合约](../2.09-deploy-contract/) 模块
- 已安装 Go 环境（1.18+）
- 拥有以太坊节点访问地址
- 拥有已部署合约的地址和 ABI

## 核心概念

### 合约加载流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ 已部署合约   │────>│ ABI 文件    │────>│ abigen      │
│ 合约地址     │     │             │     │ 生成 Go 代码 │
└─────────────┘     └─────────────┘     └─────────────┘
                                                  │
                                                  ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ 调用合约函数 │<────│ 合约实例     │<────│ NewStore()  │
│ 发送交易     │     │             │     │ 加载合约     │
└─────────────┘     └─────────────┘     └─────────────┘
```

### 部署 vs 加载

| 特性 | 部署合约 | 加载合约 |
|------|----------|----------|
| 输入 | 源代码/字节码 | 合约地址 + ABI |
| 生成 | 新合约地址 | 合约实例 |
| 操作 | 上传合约到链上 | 连接已存在的合约 |
| 需要私钥 | ✅ 是 | ❌ 否（仅需读取） |

---

## 生成 Go 绑定代码

加载合约需要使用 `abigen` 工具生成 Go 绑定代码。与部署合约不同，加载合约**仅需要 ABI 文件**，不需要字节码。

### 安装 abigen 工具

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

### 生成 Go 代码

仅使用 ABI 文件生成绑定代码：

```bash
abigen \
  --abi=Store_sol_Store.abi \
  --pkg=store \
  --out=store.go
```

**参数说明：**

| 参数 | 说明 |
|------|------|
| `--abi` | ABI 文件路径 |
| `--pkg` | 生成的 Go 包名 |
| `--out` | 输出文件名 |

**注意：** 仅使用 ABI 文件生成的代码中**不包含**部署合约的函数（如 `DeployStore`），但包含所有与合约交互的方法。

### 生成的代码结构

生成的 `store.go` 文件包含：
- `Store` 结构体（合约实例）
- `NewStore` 函数（加载合约实例）
- 每个 public 函数对应的 Go 方法
- 事件对应的结构体和过滤方法

---

## 加载合约实例

加载合约实例需要两个参数：
1. **ethclient 实例**：用于与以太坊网络通信
2. **合约地址**：已部署合约的地址

### 完整代码示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/dapp-learning/ethclient/load-contract/store"
)

const (
    // 替换为你的合约地址
    contractAddress = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
)

func main() {
    // 连接到以太坊节点
    // Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
    // Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
    client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 加载合约实例
    storeContract, err := store.NewStore(common.HexToAddress(contractAddress), client)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("合约加载成功: %s\n", contractAddress)

    // 现在可以调用合约的方法
    version, err := storeContract.Version(nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("合约版本: %s\n", version)
}
```

### 代码解析

| 步骤 | 说明 |
|------|------|
| `ethclient.Dial()` | 连接到以太坊节点 |
| `common.HexToAddress()` | 将十六进制地址字符串转换为 `common.Address` 类型 |
| `store.NewStore()` | 创建合约实例 |
| `storeContract.Version(nil)` | 调用合约的 view 函数（不需要交易） |

### 调用合约函数

**读取函数（view/pure）：**

```go
// 调用不需要 gas 的函数
result, err := storeContract.Version(nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

**写入函数（需要交易）：**

```go
import (
    "context"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/crypto"
)

// 创建交易认证器（需要私钥）
privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

// 调用需要 gas 的函数
tx, err := storeContract.SetItem(auth, [32]byte{}, [32]byte{})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
```

---

## Remix 加载合约

Remix 是一个在线 Solidity IDE，可以方便地加载已部署的合约。

### 步骤

1. **打开 Remix**：https://remix.ethereum.org/

2. **创建合约文件**：
   - 在文件浏览器中创建 `Store.sol`
   - 复制合约源代码到文件中

3. **编译合约**：
   - 选择 "Solidity Compiler"
   - 点击 "Compile" 按钮

4. **连接钱包**：
   - 选择 "Deploy & Run Transactions"
   - Environment 选择 "Injected Web3"
   - 使用 MetaMask 连接到测试网

5. **加载合约**：
   - 在 "Contract" 下拉菜单中选择 Store 合约
   - 在 "At Address" 输入框中输入合约地址
   - 点击 "At Address" 按钮

6. **交互**：
   - 在 "Deployed/Unpinned Contracts" 栏中可以看到加载的合约
   - 展开合约即可调用函数或发送交易

### Remix 加载示意图

```
┌─────────────────────────────────────────────────────┐
│  Deploy & Run Transactions                          │
├─────────────────────────────────────────────────────┤
│  Environment: [Injected Web3 - MetaMask ▼]         │
│  Contract:    [Store ▼]                             │
│                                                     │
│  At Address:  [0x8D41...EC24e]    [At Address]     │
├─────────────────────────────────────────────────────┤
│  Deployed/Unpinned Contracts:                       │
│  ┌─ Store at 0x8D41...EC24e ▼                      │
│  │  version    [Call]                              │
│  │  getItem    [key] [Call]                        │
│  │  setItem    [key, value] [setItem]              │
│  └───────────────────────────────────────────────  │
└─────────────────────────────────────────────────────┘
```

---

## 常见问题

### Q1: 加载合约需要私钥吗？

**不需要。** 加载合约实例只是创建一个 Go 对象，用于与合约交互。只有在调用需要修改状态的函数（发送交易）时才需要私钥。

```go
// ✅ 只需读取数据，不需要私钥
version, err := storeContract.Version(nil)

// ❌ 修改状态，需要私钥和交易认证
tx, err := storeContract.SetItem(auth, key, value)
```

### Q2: 如何验证合约地址是否正确？

使用 `ethclient.CodeAt` 检查地址是否有合约代码：

```go
code, err := client.CodeAt(context.Background(), contractAddress, nil)
if err != nil {
    log.Fatal(err)
}

if len(code) == 0 {
    log.Fatal("该地址没有合约代码")
}

fmt.Printf("合约代码长度: %d 字节\n", len(code))
```

### Q3: 合约地址可以是 EOAccount 地址吗？

**不可以。** 只有智能合约地址才有代码。普通账户（EOAccount）没有合约代码，加载会失败。

### Q4: 如何获取合约 ABI？

有几种方式获取 ABI：

1. **从源码生成**（如果有合约源码）：
   ```bash
   solcjs --abi Store.sol
   ```

2. **从 Etherscan 获取**：
   - 访问合约地址的 Etherscan 页面
   - 点击 "Contract" → "Contract ABI"
   - 复制 ABI 内容

3. **从 Remix 导出**：
   - 在 Remix 中编译合约后
   - 点击编译图标旁边的 "ABI" 按钮
   - 复制或下载 ABI 文件

### Q5: 加载合约和部署合约时 abigen 的区别？

| 特性 | 部署合约 | 加载合约 |
|------|----------|----------|
| abigen 参数 | `--bin + --abi` | 仅 `--abi` |
| 生成代码 | 包含 `DeployStore` | 不含部署函数 |
| 使用场景 | 首次部署 | 连接已部署合约 |

---

## 练习作业

开始练习前，请先准备：

```bash
# 安装 abigen（如果还未安装）
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 生成 Go 绑定代码
abigen --abi=Store_sol_Store.abi --pkg=store --out=store.go

# 安装依赖
go mod tidy
```

### 作业 1：加载合约并调用 view 函数（基础）

练习文件：[exercises/01-load-contract.go](exercises/01-load-contract.go)

加载已部署的 Store 合约：
1. 连接到测试网
2. 使用合约地址加载合约实例
3. 调用 `Version()` 函数获取合约版本
4. 调用 `GetItem()` 函数获取存储的值

**运行练习：**
```bash
export CONTRACT_ADDRESS=0xYourContractAddress
export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
go run exercises/01-load-contract.go
```

**参考答案：** [solutions/01-load-contract.go](solutions/01-load-contract.go)

---

### 作业 2：验证合约是否存在（进阶）

练习文件：[exercises/02-verify-contract.go](exercises/02-verify-contract.go)

验证合约地址是否有效：
1. 连接到测试网
2. 检查地址是否有合约代码
3. 检查合约代码长度是否合理
4. 尝试加载合约并捕获错误

**运行练习：**
```bash
go run exercises/02-verify-contract.go
```

**参考答案：** [solutions/02-verify-contract.go](solutions/02-verify-contract.go)

---

### 作业 3：加载多个合约实例（挑战）

练习文件：[exercises/03-load-multiple.go](exercises/03-load-multiple.go)

加载多个合约实例：
1. 定义多个合约地址（有效和无效的）
2. 批量加载这些合约
3. 验证每个合约是否成功加载
4. 输出统计信息（成功/失败数量）

**运行练习：**
```bash
go run exercises/03-load-multiple.go
```

**参考答案：** [solutions/03-load-multiple.go](solutions/03-load-multiple.go)

---

## 测试网资源

### 测试网节点

| 提供商 | URL |
|--------|-----|
| Alchemy | `https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` |
| Infura | `https://sepolia.infura.io/v3/YOUR_API_KEY` |

### 示例合约地址

如果你没有部署过合约，可以使用以下测试地址（Sepolia 测试网）：

```
0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e
```

### 测试网水龙头

- [Alchemy Faucet](https://www.alchemy.com/faucets/ethereum-sepolia)
- [Infura Faucet](https://www.infura.io/faucet/sepolia)
- [QuickNode Faucet](https://faucet.quicknode.com/ethereum/sepolia)

---

## 下一步学习

- [调用合约函数](../2.11-call-contract/)（如果存在）
- [监听合约事件](../2.13-contract-events/)（如果存在）
- [ETH 转账](../2.06-transfer-eth/)
- [代币转账](../2.07-transfer-token/)
