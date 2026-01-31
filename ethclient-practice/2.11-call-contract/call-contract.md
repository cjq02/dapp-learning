# Ethclient 调用合约学习指南

> **预计学习时间：** 45 分钟
>
> **难度：** 中等

本指南介绍如何使用 Go 语言的 `go-ethereum` 库调用已部署的智能合约，包括读取数据和发送交易。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [调用 view/pure 函数](#调用-viewpure-函数)
- [调用写函数（发送交易）](#调用写函数发送交易)
- [不使用 abigen 调用合约](#不使用-abigen-调用合约)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解合约调用的两种方式（eth_call 和 eth_sendRawTransaction）
- 使用生成的 Go 绑定代码调用合约函数
- 调用 view/pure 函数读取合约数据
- 调用写函数发送交易修改合约状态
- 手动构造调用数据（不使用 abigen）
- 处理交易回执和等待交易确认

## 前置条件

- Go 语言基础
- Solidity 基础（了解 view/pure 函数与写函数的区别）
- 已完成 [加载合约](../2.10-load-contract/) 模块
- 已安装 Go 环境（1.18+）
- 拥有以太坊节点访问地址
- 拥有已部署合约的地址和 ABI

## 核心概念

### 合约调用的两种方式

```
┌─────────────────────────────────────────────────────────────┐
│                    合约调用方式                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────┐      ┌─────────────────────┐     │
│  │   eth_call          │      │ eth_sendRawTransaction│   │
│  │   (读取数据)         │      │   (发送交易)          │   │
│  ├─────────────────────┤      ├─────────────────────┤     │
│  │ • view/pure 函数    │      │ • 写函数              │   │
│  │ • 不消耗 Gas        │      │ • 消耗 Gas            │   │
│  │ • 不修改状态        │      │ • 修改链上状态         │   │
│  │ • 不需要私钥        │      │ • 需要私钥签名         │   │
│  │ • 立即返回          │      │ • 需要等待确认         │   │
│  └─────────────────────┘      └─────────────────────┘     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 读取 vs 写入

| 特性 | 读取（eth_call） | 写入（eth_sendRawTransaction） |
|------|------------------|-------------------------------|
| 函数类型 | view / pure | 非 view/pure |
| Gas 消耗 | 0 | 需要 Gas |
| 私钥 | 不需要 | 需要 |
| 状态改变 | 无 | 修改链上状态 |
| 返回方式 | 立即返回 | 返回交易哈希，需等待确认 |

---

## 调用 view/pure 函数

view/pure 函数不会修改合约状态，仅读取数据，调用时不需要 Gas，也不需要私钥签名。

### 方法：直接调用生成的函数

```go
// 使用 nil 作为 CallOpts 参数
result, err := contract.MethodName(nil, arg1, arg2)
```

### 示例：读取合约版本

```go
package main

import (
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/dapp-learning/ethclient/call-contract/store"
)

const (
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

    // 调用 view 函数，传入 nil 作为 CallOpts
    version, err := storeContract.Version(nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("合约版本: %s\n", version)
}
```

### 示例：读取存储的值

```go
// 准备 key
var key [32]byte
copy(key[:], []byte("my_key"))

// 调用 GetItem 函数
value, err := storeContract.GetItem(nil, key)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("存储值: %x\n", value)
```

### CallOpts 参数

虽然通常使用 `nil`，但也可以自定义 `CallOpts`：

```go
import "github.com/ethereum/go-ethereum/accounts/abi/bind"

opts := &bind.CallOpts{
    Pending: false,        // 是否使用 pending 状态
    Context: context.Background(),  // 上下文
    BlockNumber: big.NewInt(12345), // 指定区块号
}

result, err := storeContract.Version(opts)
```

---

## 调用写函数（发送交易）

写函数会修改合约状态，需要发送交易，消耗 Gas，并且需要私钥签名。

### 准备工作

```go
import (
    "context"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/crypto"
)

// 1. 从私钥创建交易认证器
privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
if err != nil {
    log.Fatal(err)
}

// 2. 创建 TransactOpts（需要指定 ChainID）
chainID := big.NewInt(11155111) // Sepolia
auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
if err != nil {
    log.Fatal(err)
}
```

### 示例：调用 SetItem 函数

```go
package main

import (
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/dapp-learning/ethclient/call-contract/store"
)

const (
    contractAddress = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
)

func main() {
    // 连接到以太坊节点
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

    // 从私钥创建交易认证器
    privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
    if err != nil {
        log.Fatal(err)
    }

    chainID := big.NewInt(11155111) // Sepolia
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    if err != nil {
        log.Fatal(err)
    }

    // 准备数据
    var key [32]byte
    var value [32]byte
    copy(key[:], []byte("demo_key"))
    copy(value[:], []byte("demo_value"))

    // 调用写函数 SetItem
    tx, err := storeContract.SetItem(auth, key, value)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())

    // 等待交易确认
    receipt, err := bind.WaitMined(context.Background(), client, tx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
    fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
}
```

### TransactOpts 常用选项

```go
auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
if err != nil {
    log.Fatal(err)
}

// 设置 Gas 价格
auth.GasPrice = big.NewInt(30000000000) // 30 Gwei

// 设置 Gas 限制
auth.GasLimit = uint64(100000)

// 设置 Value（发送 ETH）
auth.Value = big.NewInt(1000000000000000000) // 1 ETH
```

---

## 不使用 abigen 调用合约

如果不使用 `abigen` 生成代码，需要手动构造调用数据。这种方式更底层，但更灵活。

### 使用 ABI 调用合约

```go
package main

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "log"
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

const (
    contractAddress = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
    contractABI = `[{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"getItem","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"}]`
)

func main() {
    client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // 解析 ABI
    parsedABI, err := abi.JSON(strings.NewReader(contractABI))
    if err != nil {
        log.Fatal(err)
    }

    // 准备数据
    var key [32]byte
    var value [32]byte
    copy(key[:], []byte("my_key"))
    copy(value[:], []byte("my_value"))

    // 打包函数调用数据
    data, err := parsedABI.Pack("setItem", key, value)
    if err != nil {
        log.Fatal(err)
    }

    // 获取发送地址和 nonce
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

    // 创建交易
    chainID := big.NewInt(11155111)
    tx := types.NewTransaction(
        nonce,
        common.HexToAddress(contractAddress),
        big.NewInt(0),
        300000,
        gasPrice,
        data,
    )

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

### 不使用 ABI 手动构造调用数据

```go
import "github.com/ethereum/go-ethereum/crypto"

// 1. 计算函数选择器（函数签名的 Keccak256 哈希的前 4 字节）
methodSignature := []byte("setItem(bytes32,bytes32)")
methodSelector := crypto.Keccak256(methodSignature)[:4]

// 2. 准备参数
var key [32]byte
var value [32]byte
copy(key[:], []byte("my_key"))
copy(value[:], []byte("my_value"))

// 3. 组合调用数据：选择器 + 参数
var inputData []byte
inputData = append(inputData, methodSelector...)
inputData = append(inputData, key[:]...)
inputData = append(inputData, value[:]...)

// inputData 现在可以用于交易的 Data 字段
```

---

## 常见问题

### Q1: 调用合约函数时传入 `nil` 是什么意思？

`nil` 表示使用默认的 `CallOpts`：

```go
// 这两种写法等价
result, err := contract.MethodName(nil, args)
result, err := contract.MethodName(&bind.CallOpts{}, args)
```

### Q2: 如何判断一个函数是 view/pure 还是写函数？

查看合约函数定义：

```solidity
// view 函数 - 不修改状态
function version() public view returns (string) {}

// pure 函数 - 不读取也不修改状态
function add(uint a, uint b) public pure returns (uint) {}

// 写函数 - 修改状态
function setItem(bytes32 key, bytes32 value) public {}
```

### Q3: 为什么调用写函数没有立即返回结果？

写函数发送的是交易，需要：
1. 被矿工打包进区块
2. 区块被确认

使用 `bind.WaitMined` 等待交易确认：

```go
receipt, err := bind.WaitMined(context.Background(), client, tx)
```

### Q4: 如何获取交易的回执信息？

```go
receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
if err != nil {
    log.Fatal(err)
}

fmt.Printf("状态: %d\n", receipt.Status)     // 1 = 成功, 0 = 失败
fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
```

### Q5: 交易失败了怎么办？

检查 `receipt.Status`：

```go
receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
if err != nil {
    log.Fatal(err)
}

if receipt.Status == 0 {
    log.Fatal("交易失败")
}
```

---

## 练习作业

开始练习前，请先准备：

```bash
# 生成 Go 绑定代码（如果还没有）
abigen --abi=Store_sol_Store.abi --pkg=store --out=store.go

# 安装依赖
go mod tidy
```

### 作业 1：读取合约数据（基础）

练习文件：[exercises/01-read-contract.go](exercises/01-read-contract.go)

读取 Store 合约的数据：
1. 连接到测试网
2. 加载合约实例
3. 调用 `Version()` 函数获取版本
4. 调用 `GetItem()` 函数获取存储的值

**运行练习：**
```bash
export CONTRACT_ADDRESS=0xYourContractAddress
export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
go run exercises/01-read-contract.go
```

**参考答案：** [solutions/01-read-contract.go](solutions/01-read-contract.go)

---

### 作业 2：写入合约数据（进阶）

练习文件：[exercises/02-write-contract.go](exercises/02-write-contract.go)

向 Store 合约写入数据：
1. 从环境变量读取私钥
2. 创建交易认证器
3. 调用 `SetItem()` 函数存储数据
4. 等待交易确认并打印结果

**运行练习：**
```bash
export PRIVATE_KEY=your_private_key
export CONTRACT_ADDRESS=0xYourContractAddress
export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
go run exercises/02-write-contract.go
```

**参考答案：** [solutions/02-write-contract.go](solutions/02-write-contract.go)

---

### 作业 3：使用手动 ABI 调用合约（挑战）

练习文件：[exercises/03-manual-call.go](exercises/03-manual-call.go)

不使用 abigen 生成的代码，手动构造调用：
1. 使用 ABI 字符串解析合约
2. 手动打包函数调用数据
3. 创建并签名交易
4. 发送交易并等待确认

**运行练习：**
```bash
export PRIVATE_KEY=your_private_key
export CONTRACT_ADDRESS=0xYourContractAddress
export SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
go run exercises/03-manual-call.go
```

**参考答案：** [solutions/03-manual-call.go](solutions/03-manual-call.go)

---

## 测试网资源

### 测试网节点

| 提供商 | URL |
|--------|-----|
| Alchemy | `https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` |
| Infura | `https://sepolia.infura.io/v3/YOUR_API_KEY` |

### 示例合约地址

```
0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e
```

### 测试网水龙头

- [Alchemy Faucet](https://www.alchemy.com/faucets/ethereum-sepolia)
- [Infura Faucet](https://www.infura.io/faucet/sepolia)

---

## 下一步学习

- [监听合约事件](../2.13-contract-events/)（如果存在）
- [订阅新区块](../2.08-subscribe-block/)
- [代币转账](../2.07-transfer-token/)
