# Ethclient 部署合约学习指南

> **预计学习时间：** 50 分钟
>
> **难度：** 进阶

本指南介绍如何使用 Go 语言的 `go-ethereum` 库部署智能合约到以太坊区块链。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [Solidity 合约](#solidity-合约)
- [编译合约](#编译合约)
- [生成 Go 绑定代码](#生成-go-绑定代码)
- [部署合约（方法一：使用 bind）](#部署合约方法一使用-bind)
- [部署合约（方法二：纯 ethclient）](#部署合约方法二纯-ethclient)
- [Remix 部署](#remix-部署)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 编写简单的 Solidity 智能合约
- 使用 `solc` 编译智能合约
- 使用 `abigen` 生成 Go 绑定代码
- 使用 ethclient 部署智能合约
- 使用纯 ethclient（不用 bind）部署合约
- 验证合约部署是否成功

## 前置条件

- Go 语言基础
- Solidity 基础（了解合约结构、函数、事件）
- 已完成 [ETH 转账](../2.06-transfer-eth/) 模块
- 已安装 Go 环境（1.18+）
- Node.js 环境（用于安装 solc）
- 拥有以太坊节点访问地址
- 拥有测试网 ETH 和私钥

## 核心概念

### 智能合约部署流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Store.sol   │────>│ solc 编译    │────>│ Store.bin   │
│ (源代码)     │     │             │     │ (字节码)     │
└─────────────┘     └─────────────┘     └─────────────┘
                                                  │
                                                  ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ 部署到链上   │<────│ 签名并发送   │<────│ abigen/bind │
│             │     │ 交易         │     │ 生成部署代码 │
└─────────────┘     └─────────────┘     └─────────────┘
```

### 部署即交易

部署合约本质上是一笔特殊的交易：

```
普通交易                        合约部署交易
├── to: 接收地址                 ├── to: 0x0000...0000 (空)
├── value: ETH 数量              ├── value: ETH 数量
├── data: (空)                  ├── data: 合约字节码 + 构造参数
└── gasLimit                    └── gasLimit
                               结果：
                               └── contractAddress (合约地址)
```

---

## Solidity 合约

### 示例合约：Store.sol

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract Store {
    event ItemSet(bytes32 indexed key, bytes32 value);

    string public version;
    mapping(bytes32 => bytes32) public items;

    // 构造函数：部署时执行一次
    constructor(string memory _version) {
        version = _version;
    }

    // 写入函数
    function setItem(bytes32 key, bytes32 value) external {
        items[key] = value;
        emit ItemSet(key, value);
    }

    // 读取函数
    function getItem(bytes32 key) external view returns (bytes32) {
        return items[key];
    }
}
```

**合约文件位置：** [contract/Store.sol](contract/Store.sol)

---

## 编译合约

### 安装 solc 编译器

```bash
npm install -g solc
```

验证安装：

```bash
solcjs --version
# 输出：0.8.26+commit.8a97fa7a.Emscripten.clang
```

### 编译合约

在 `contract/` 目录下执行：

```bash
cd contract
solcjs --bin Store.sol
```

生成文件：`Store_sol_Store.bin`（字节码）

生成 ABI：

```bash
solcjs --abi Store.sol
```

生成文件：`Store_sol_Store.abi`（ABI 接口）

### 编译输出

```
contract/
├── Store.sol                  # 源代码
├── Store_sol_Store.bin        # 字节码（部署用）
└── Store_sol_Store.abi        # ABI（交互用）
```

---

## 生成 Go 绑定代码

### 安装 abigen 工具

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

### 生成 Go 代码

```bash
abigen \
  --bin=Store_sol_Store.bin \
  --abi=Store_sol_Store.abi \
  --pkg=store \
  --out=store.go
```

**参数说明：**

| 参数 | 说明 |
|------|------|
| `--bin` | 字节码文件路径 |
| `--abi` | ABI 文件路径 |
| `--pkg` | 生成的 Go 包名 |
| `--out` | 输出文件名 |

生成的 `store.go` 包含：
- `Store` 结构体（合约实例）
- `DeployStore` 函数（部署合约）
- 每个 Solidity 函数对应的 Go 方法

---

## 部署合约（方法一：使用 bind）

这是推荐的方法，使用 abigen 生成的代码。

### 完整代码

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊节点
    client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // 加载私钥
    privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // 获取链 ID
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // 创建交易认证器
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    if err != nil {
        log.Fatal(err)
    }

    // 设置 Gas 相关参数
    gasPrice, _ := client.SuggestGasPrice(context.Background())
    auth.GasLimit = uint64(300000)  // 部署通常需要较多 Gas
    auth.GasPrice = gasPrice

    // 调用生成的部署函数
    contractAddr, tx, _, err := DeployStore(auth, client, "1.0")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("合约地址: %s\n", contractAddr.Hex())
    fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
}
```

---

## 部署合约（方法二：纯 ethclient）

不使用 abigen，直接用字节码部署。

### 合约字节码常量

```go
const contractBytecode = "608060405234801561000f575f80fd5b50..."
```

### 完整代码

```go
package main

import (
    "context"
    "encoding/hex"
    "log"
    "math/big"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // 获取发送者地址和 nonce
    publicKey := privateKey.Public()
    publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        log.Fatal(err)
    }

    // 获取 Gas 价格
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // 解码合约字节码
    data, err := hex.DecodeString(contractBytecode)
    if err != nil {
        log.Fatal(err)
    }

    // 创建合约部署交易
    // 注意：to 为 nil，表示部署合约
    tx := types.NewContractCreation(
        nonce,
        big.NewInt(0),      // value
        3000000,            // gasLimit
        gasPrice,           // gasPrice
        data,               // 合约字节码
    )

    // 获取链 ID 并签名
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }

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

    // 等待交易确认
    receipt, err := waitForReceipt(client, signedTx.Hash())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("合约部署成功: %s\n", receipt.ContractAddress.Hex())
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
    for {
        receipt, err := client.TransactionReceipt(context.Background(), txHash)
        if err == nil {
            return receipt, nil
        }
        time.Sleep(1 * time.Second)
    }
}
```

### 两种方法对比

| 特性 | bind 方法 | 纯 ethclient 方法 |
|------|-----------|-------------------|
| 代码复杂度 | 低 | 高 |
| 类型安全 | ✅ 是 | ❌ 否 |
| 后续调用 | 方便 | 需要手动构造 |
| 适用场景 | 生产推荐 | 学习底层原理 |

---

## Remix 部署

Remix 是一个在线 Solidity IDE，可以直接部署合约。

### 步骤

1. **打开 Remix**：https://remix.ethereum.org/

2. **创建合约**：
   - 在文件浏览器中创建 `Store.sol`
   - 复制 [contract/Store.sol](contract/Store.sol) 内容

3. **编译合约**：
   - 选择 "Solidity Compiler"
   - 点击 "Compile"

4. **连接钱包**：
   - 选择 "Deploy & Run Transactions"
   - Environment 选择 "Injected Web3"
   - 使用 MetaMask 连接到 Sepolia 测试网

5. **部署**：
   - 点击 "Deploy"
   - 在 MetaMask 中确认交易

6. **验证**：
   - 在 [Sepolia Etherscan](https://sepolia.etherscan.io/) 搜索合约地址
   - 在 Remix 的 "Deployed Contracts" 中与合约交互

---

## 常见问题

### Q1: 如何判断合约是否部署成功？

```go
receipt, err := client.TransactionReceipt(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}

// 检查状态
if receipt.Status == 0 {
    log.Fatal("合约部署失败")
}

// 合约地址
if receipt.ContractAddress == (common.Address{}) {
    log.Fatal("这不是合约部署交易")
}

fmt.Printf("合约地址: %s\n", receipt.ContractAddress.Hex())
```

### Q2: 部署合约需要多少 Gas？

部署合约的 Gas 消耗取决于：

| 因素 | 影响 |
|------|------|
| 合约代码大小 | 代码越多，Gas 越多 |
| 构造函数复杂度 | 逻辑越复杂，Gas 越多 |
| 初始化存储数据 | 存储越多，Gas 越多 |

一般估算：
- 简单合约：~500,000 Gas
- 中等合约：~1,000,000 Gas
- 复杂合约：~2,000,000+ Gas

### Q3: 合约地址是如何确定的？

合约地址由发送者地址和 nonce 决定：

```
contract_address = keccak256(rlp.encode(sender, nonce))[12:]
```

**重要特性：**
- 相同地址，相同 nonce → 相同合约地址
- 这意味着可以在部署前计算合约地址

```go
// 计算预期合约地址
fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
nonce := uint64(0)  // 当前 nonce

contractAddress := crypto.CreateAddress(fromAddress, nonce)
fmt.Printf("预期合约地址: %s\n", contractAddress.Hex())
```

### Q4: 如何在测试网获取 ETH？

参见 [查询区块](../2.01-query-block/) 模块中的测试网资源部分。

---

## 练习作业

开始练习前，请先准备：

```bash
# 安装 solc
npm install -g solc

# 安装 abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 编译合约
cd contract
solcjs --bin Store.sol
solcjs --abi Store.sol

# 生成 Go 绑定
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=store.go

# 安装依赖
cd ..
go mod tidy
```

### 作业 1：使用 bind 部署合约（基础）

练习文件：[exercises/01-deploy-with-bind.go](exercises/01-deploy-with-bind.go)

使用 abigen 生成的代码部署 Store 合约：
1. 连接到测试网
2. 加载私钥
3. 调用 DeployStore 函数
4. 等待交易确认并输出合约地址

**运行练习：**
```bash
export PRIVATE_KEY=your_private_key_here
export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
go run exercises/01-deploy-with-bind.go
```

**参考答案：** [solutions/01-deploy-with-bind.go](solutions/01-deploy-with-bind.go)

---

### 作业 2：纯 ethclient 部署（进阶）

练习文件：[exercises/02-deploy-raw.go](exercises/02-deploy-raw.go)

不使用 bind，直接用合约字节码部署：
1. 手动构造合约部署交易
2. 签名并发送
3. 等待确认
4. 输出合约地址和交易收据

**运行练习：**
```bash
go run exercises/02-deploy-raw.go
```

**参考答案：** [solutions/02-deploy-raw.go](solutions/02-deploy-raw.go)

---

### 作业 3：计算预期合约地址（挑战）

练习文件：[exercises/03-predict-address.go](exercises/03-predict-address.go)

在部署前计算合约地址：
1. 获取当前 nonce
2. 使用 `crypto.CreateAddress` 计算预期地址
3. 部署合约
4. 验证实际地址与预期地址是否一致

**运行练习：**
```bash
go run exercises/03-predict-address.go
```

**参考答案：** [solutions/03-predict-address.go](solutions/03-predict-address.go)

---

## 测试网资源

### 测试网节点

| 提供商 | URL |
|--------|-----|
| Alchemy | `https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` |
| Infura | `https://sepolia.infura.io/v3/YOUR_API_KEY` |

### 测试网水龙头

- [Alchemy Faucet](https://www.alchemy.com/faucets/ethereum-sepolia)
- [Infura Faucet](https://www.infura.io/faucet/sepolia)
- [QuickNode Faucet](https://faucet.quicknode.com/ethereum/sepolia)

---

## 下一步学习

- [加载合约](../2.10-load-contract/)
- [调用合约函数](../2.11-call-contract/)
- [监听合约事件](../2.13-contract-events/)
