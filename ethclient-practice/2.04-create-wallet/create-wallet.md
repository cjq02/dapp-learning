# Ethclient 创建新钱包学习指南

本指南介绍如何使用 Go 语言的 `go-ethereum` 库创建以太坊钱包，包括私钥、公钥和地址的生成与转换。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [生成新钱包](#生成新钱包)
- [从私钥恢复钱包](#从私钥恢复钱包)
- [地址生成原理](#地址生成原理)
- [常见问题](#常见问题)
- [练习作业](#练习作业)
- [标准答案](#标准答案)

---

## 学习目标

完成本指南后，你将能够：
- 理解以太坊钱包的组成（私钥、公钥、地址）
- 使用 crypto 包生成新的随机私钥
- 从私钥派生公钥和地址
- 从十六进制私钥恢复钱包
- 理解以太坊地址的生成原理（Keccak-256 哈希）

## 前置条件

- Go 语言基础（变量、函数、错误处理、类型断言）
- 加密学基础概念（非对称加密、哈希）
- 已安装 Go 环境（1.18+）
- 了解以太坊账户系统

## 核心概念

### 钱包结构

```
┌─────────────────────────────────────────┐
│              以太坊钱包                   │
├─────────────────────────────────────────┤
│  私钥 (Private Key)                      │
│  - 32 字节随机数                         │
│  - 用于签名交易                          │
│  - 绝对不能泄露！                        │
├─────────────────────────────────────────┤
│  公钥 (Public Key)                       │
│  - 64 字节（未压缩）                     │
│  - 从私钥通过椭圆曲线算法派生             │
│  - 可以公开分享                          │
├─────────────────────────────────────────┤
│  地址 (Address)                          │
│  - 20 字节                               │
│  - 公钥的 Keccak-256 哈希的后 20 字节    │
│  - 用作接收资金的账户标识                │
└─────────────────────────────────────────┘
```

### 密钥派生关系

```
私钥 (32 字节)
    │
    ├─> 椭圆曲线算法 (secp256k1)
    │       │
    │       └─> 公钥 (64 字节)
    │               │
    │               └─> Keccak-256 哈希
    │                       │
    │                       └─> 取后 20 字节
    │                               │
    │                               └─> 地址 (20 字节)
```

### 为什么需要这样设计？

| 组件 | 保密性 | 用途 |
|------|--------|------|
| 私钥 | 🔴 绝对保密 | 签名交易，证明所有权 |
| 公钥 | 🟢 可以公开 | 验证签名 |
| 地址 | 🟢 可以公开 | 接收资金，标识账户 |

---

## 生成新钱包

### 导入必要的包

```go
import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
    "golang.org/x/crypto/sha3"
)
```

### 方法：`GenerateKey`

```go
func GenerateKey() (*ecdsa.PrivateKey, error)
```

生成一个随机的 ECDSA 私钥。

### 示例：生成新钱包

```go
package main

import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {
    // 生成新的私钥
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        log.Fatal(err)
    }

    // 将私钥转换为字节
    privateKeyBytes := crypto.FromECDSA(privateKey)

    // 转换为十六进制字符串（去掉 '0x' 前缀）
    privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
    fmt.Printf("私钥: %s\n", privateKeyHex)
    // 输出: 类似 ccec5314acec3d18eae81b6bd988b844fc4f7f7d3c828b351de6d0fede02d3f2
}
```

### 示例：从私钥派生公钥

```go
// 从私钥获取公钥
publicKey := privateKey.Public()

// 类型断言：转换为 *ecdsa.PublicKey
publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
if !ok {
    log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
}

// 将公钥转换为字节
publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

// 转换为十六进制（去掉 '0x' 和 '04' 前缀）
publicKeyHex := hexutil.Encode(publicKeyBytes)[4:]
fmt.Printf("公钥: %s\n", publicKeyHex)
```

### 示例：从公钥生成地址

```go
// 方法一：使用内置函数（推荐）
address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
fmt.Printf("地址: %s\n", address)
// 输出: 类似 0x96Ded6f7d6E5B45b5a3B3E5B8C1d2C3D4E5F6A7B

// 方法二：手动计算（理解原理）
hash := sha3.NewLegacyKeccak256()
hash.Write(publicKeyBytes[1:])  // 跳过第一个字节
addressBytes := hash.Sum(nil)[12:]  // 取后 20 字节
addressManual := hexutil.Encode(addressBytes)
fmt.Printf("地址（手动）: %s\n", addressManual)
```

### 完整示例：生成新钱包

```go
package main

import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
    "golang.org/x/crypto/sha3"
)

func main() {
    // 1. 生成私钥
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        log.Fatal(err)
    }

    privateKeyBytes := crypto.FromECDSA(privateKey)
    fmt.Printf("私钥 (Hex): %s\n", hexutil.Encode(privateKeyBytes)[2:])

    // 2. 派生公钥
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
    fmt.Printf("公钥 (Hex): %s\n", hexutil.Encode(publicKeyBytes)[4:])

    // 3. 生成地址
    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
    fmt.Printf("地址: %s\n", address)

    // 4. 验证地址生成原理
    hash := sha3.NewLegacyKeccak256()
    hash.Write(publicKeyBytes[1:])
    fmt.Printf("完整哈希: %s\n", hexutil.Encode(hash.Sum(nil)[:]))
    fmt.Printf("后20字节: %s\n", hexutil.Encode(hash.Sum(nil)[12:]))
}
```

---

## 从私钥恢复钱包

### 方法：`HexToECDSA`

```go
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error)
```

从十六进制私钥字符串恢复 ECDSA 私钥。

### 示例：从已有私钥恢复钱包

```go
package main

import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {
    // 已有的私钥（十六进制格式）
    privateKeyHex := "ccec5314acec3d18eae81b6bd988b844fc4f7f7d3c828b351de6d0fede02d3f2"

    // 从十六进制恢复私钥
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        log.Fatal(err)
    }

    // 派生公钥和地址
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
    fmt.Printf("从私钥恢复的地址: %s\n", address)
}
```

---

## 地址生成原理

### 为什么地址是 20 字节？

以太坊地址是公钥的 Keccak-256 哈希值的**后 20 字节**。

```
公钥 (64 字节未压缩)
    │
    ├─> Keccak-256 哈希
    │       │
    │       └─> 32 字节哈希值
    │               │
    │               └─> 取后 20 字节
    │                       │
    │                       └─> 以太坊地址
```

### 为什么要跳过公钥的第一个字节？

未压缩的公钥以 `0x04` 开头（EC 前缀），表示这是一个未压缩的公钥。

```go
// 公钥格式: 0x04 + X坐标 (32字节) + Y坐标 (32字节)
publicKeyBytes[1:]  // 跳过 0x04
```

### 为什么要取后 20 字节？

这是以太坊的设计选择，20 字节（160 位）的地址长度在：
- 安全性：足够大，几乎不可能碰撞
- 效率：足够小，便于存储和传输
之间取得了平衡。

---

## 常见问题

### Q1: 私钥丢失了怎么办？

**答：** 私钥丢失后无法恢复。没有私钥就无法访问钱包中的资产。这就是为什么必须：
- 妥善备份私钥
- 使用助记词（通过 BIP39 标准派生私钥）

### Q2: 不同私钥会生成相同地址吗？

**答：** 理论上可能，但概率极低。需要 2^160 次尝试才会发生碰撞，实际上不可能。

### Q3: 为什么公钥有两种格式？

**答：** 公钥有两种表示方式：
- **未压缩格式**：65 字节（0x04 + X + Y），以太坊使用
- **压缩格式**：33 字节（0x02/03 + X），比特币使用

### Q4: 十六进制中的 '0x' 前缀是什么？

**答：** '0x' 是十六进制的前缀，表示后面的字符是十六进制数。
- `hexutil.Encode()` 返回带 `0x` 的字符串
- `[2:]` 或 `[4:]` 用于去掉这个前缀

---

## 练习作业

开始练习前，请先进入目录并安装依赖：

```bash
cd ethclient-practice/2.04-create-wallet
go mod init 2.04-create-wallet
go get github.com/ethereum/go-ethereum
go mod tidy
```

### 作业 1：生成新钱包（基础）

练习文件：[exercises/01-generate-wallet.go](exercises/01-generate-wallet.go)

编写一个程序，实现以下功能：

1. 生成一个新的随机私钥
2. 从私钥派生公钥
3. 从公钥生成地址
4. 输出所有信息：
   - 私钥（十六进制，不带 '0x'）
   - 公钥（十六进制，不带 '0x04'）
   - 地址（带 '0x'）

**运行练习：**
```bash
go run exercises/01-generate-wallet.go
```

**参考答案：** [solutions/01-generate-wallet.go](solutions/01-generate-wallet.go)

---

### 作业 2：钱包恢复器（进阶）

练习文件：[exercises/02-restore-wallet.go](exercises/02-restore-wallet.go)

编写一个程序，从已有的私钥恢复钱包：

1. 从命令行参数或硬编码的私钥恢复钱包
2. 输出恢复的地址
3. 验证：从同一私钥多次恢复，地址应该一致

**运行练习：**
```bash
go run exercises/02-restore-wallet.go
```

**参考答案：** [solutions/02-restore-wallet.go](solutions/02-restore-wallet.go)

---

### 作业 3：地址生成验证器（挑战）

练习文件：[exercises/03-address-validator.go](exercises/03-address-validator.go)

编写一个程序，手动实现地址生成过程：

1. 生成新钱包
2. 使用 `crypto.PubkeyToAddress()` 生成地址（方法 A）
3. 手动使用 Keccak-256 哈希生成地址（方法 B）
4. 验证两种方法生成的地址一致
5. 输出详细步骤，帮助理解地址生成过程

**运行练习：**
```bash
go run exercises/03-address-validator.go
```

**参考答案：** [solutions/03-address-validator.go](solutions/03-address-validator.go)

---

## 安全提醒

⚠️ **重要安全注意事项：**

1. **永远不要分享私钥**：私钥是你资产的唯一凭证
2. **不要在代码中硬编码私钥**：使用环境变量或配置文件
3. **测试网私钥 ≠ 主网私钥**：不要在主网使用测试网私钥
4. **妥善备份**：将私钥/助记词写在纸上，存放在安全的地方

---

## 下一步学习

- [查询账户余额](../2.05-query-balance/)
- [ETH 转账](../2.06-transfer-eth/)
- [签名交易](../2.07-sign-transaction/)
