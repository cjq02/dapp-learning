# Ethclient 订阅区块学习指南

> **预计学习时间：** 40 分钟
>
> **难度：** 进阶

本指南介绍如何使用 Go 语言的 `go-ethereum` 库（ethclient）订阅以太坊区块链上的新区块，实现实时监听区块链事件。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [WebSocket 连接](#websocket-连接)
- [订阅新区块](#订阅新区块)
- [获取完整区块](#获取完整区块)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解 WebSocket 与 HTTP 连接的区别
- 使用 ethclient 建立 WebSocket 连接
- 订阅新区块事件
- 实时处理区块数据
- 构建实时区块链监听应用

## 前置条件

- Go 语言基础（channel、select、goroutine）
- 以太坊基本概念（区块、交易、哈希）
- 已完成 [查询区块](../2.01-query-block/query-block.md) 模块
- 已安装 Go 环境（1.18+）
- 拥有以太坊 WebSocket 节点访问地址

## 核心概念

### WebSocket vs HTTP

```
HTTP 连接                          WebSocket 连接
    │                                    │
    ├─ 请求 ─────────────────>           ├─ 握手 ────────────────>
    │                                    │
    │<──── 响应 ───────────────           │<──── 握手响应 ─────────
    │                                    │
    │                                    │◄──────── 保持连接 ─────►
    │                                    │
    │                                    │<───── 推送数据 ─────────
    │                                    │
    └─ 连接关闭                          └─ 持续双向通信
```

| 特性 | HTTP | WebSocket |
|------|------|-----------|
| 连接方式 | 请求-响应 | 持久连接 |
| 数据流向 | 单向 | 双向 |
| 服务器推送 | ❌ 不支持 | ✅ 支持 |
| 适用场景 | 查询历史数据 | 实时事件监听 |

### 订阅机制

```
┌─────────────┐                    ┌─────────────┐
│   Client    │                    │   Ethereum  │
│             │                    │    Node     │
├─────────────┤                    ├─────────────┤
│             │  1. SubscribeNewHead│             │
│   headers   │ ───────────────────>│             │
│  (channel)  │                     │             │
│             │                     │  ╔═════════╗│
│             │                     │  ║ 新区块   ║│
│             │                     │  ║ #12345  ║│
│             │<────── 推送 Header ─│  ╚═════════╝│
│             │                     │             │
│   ┌─────┐   │                     │  ╔═════════╗│
│   │select│   │                     │  ║ 新区块   ║│
│   └─────┘   │<────── 推送 Header ─│  ║ #12346  ║│
│             │                     │  ╚═════════╝│
└─────────────┘                     └─────────────┘
```

---

## WebSocket 连接

### WebSocket URL 格式

不同的服务提供商提供不同的 WebSocket URL：

| 提供商 | WebSocket URL |
|--------|---------------|
| Infura | `wss://sepolia.infura.io/ws/v3/YOUR_API_KEY` |
| Alchemy | `wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` |
| QuickNode | `wss://YOUR_ENDPOINT` |

### 建立连接

```go
package main

import (
    "log"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 使用 WebSocket URL 连接（注意 wss:// 协议）
    client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### 注意事项

⚠️ **重要**：
- 订阅功能**必须**使用 WebSocket 连接
- HTTP 连接（`https://`）不支持订阅
- 生产环境建议使用专业节点服务（Alchemy、Infura 等）

---

## 订阅新区块

### 步骤 1：创建通道

```go
import "github.com/ethereum/go-ethereum/core/types"

// 创建用于接收区块头的通道
headers := make(chan *types.Header)
```

### 步骤 2：订阅新区块

```go
import "context"

// 订阅新区块事件
sub, err := client.SubscribeNewHead(context.Background(), headers)
if err != nil {
    log.Fatal(err)
}
```

### 步骤 3：监听事件

```go
for {
    select {
    case err := <-sub.Err():
        // 订阅出错时退出
        log.Fatal(err)
    case header := <-headers:
        // 收到新区块头
        fmt.Printf("新区块: %s\n", header.Hash().Hex())
        fmt.Printf("区块号: %d\n", header.Number.Uint64())
    }
}
```

### 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊 WebSocket 节点
    client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 创建订阅通道
    headers := make(chan *types.Header)

    // 订阅新区块
    sub, err := client.SubscribeNewHead(context.Background(), headers)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("开始监听新区块...")

    // 监听新区块事件
    for {
        select {
        case err := <-sub.Err():
            log.Fatal(err)
        case header := <-headers:
            fmt.Printf("新区块: %s\n", header.Hash().Hex())
            fmt.Printf("区块号: %d\n", header.Number.Uint64())
        }
    }
}
```

---

## 获取完整区块

区块头（Header）只包含元数据，如需交易列表，需要获取完整区块：

```go
case header := <-headers:
    fmt.Printf("新区块哈希: %s\n", header.Hash().Hex())

    // 根据哈希获取完整区块
    block, err := client.BlockByHash(context.Background(), header.Hash())
    if err != nil {
        log.Printf("获取区块失败: %v", err)
        continue
    }

    fmt.Printf("区块号: %d\n", block.Number().Uint64())
    fmt.Printf("时间戳: %d\n", block.Time().Uint64())
    fmt.Printf("交易数量: %d\n", len(block.Transactions()))
    fmt.Printf("Gas 使用: %d\n", block.GasUsed())
    fmt.Printf("矿工: %s\n", block.Coinbase().Hex())
```

### Header vs Block

```go
Header (区块头)                    Block (完整区块)
├── Number                         ├── Number
├── Hash                           ├── Hash
├── Time                           ├── Time
├── ParentHash                     ├── ParentHash
├── GasUsed                        ├── GasUsed
├── GasLimit                       ├── GasLimit
└── ...                            └── Transactions() ← 交易列表
```

---

## 常见问题

### Q1: 为什么我的订阅不工作？

**检查清单：**

```go
// 1. 确认使用 WebSocket URL
// ❌ 错误：使用 HTTP
client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")

// ✅ 正确：使用 WebSocket
client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/YOUR_API_KEY")

// 2. 确认订阅成功
sub, err := client.SubscribeNewHead(context.Background(), headers)
if err != nil {
    log.Fatal(err)  // 这里会报错
}

// 3. 确认在主循环中监听
for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case header := <-headers:
        // 处理新区块
    }
}
```

### Q2: 如何优雅地关闭订阅？

```go
import (
    "os"
    "os/signal"
    "syscall"
)

// 监听系统信号
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case header := <-headers:
        fmt.Printf("新区块: %d\n", header.Number.Uint64())
    case <-sigCh:
        fmt.Println("\n收到退出信号，关闭订阅...")
        sub.Unsubscribe()
        return
    }
}
```

### Q3: 如何处理订阅断开重连？

```go
func subscribeWithReconnect(client *ethclient.Client) {
    for {
        headers := make(chan *types.Header)
        sub, err := client.SubscribeNewHead(context.Background(), headers)
        if err != nil {
            log.Printf("订阅失败，5秒后重试: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }

        log.Println("订阅成功，开始监听...")

        for {
            select {
            case err := <-sub.Err():
                log.Printf("订阅断开: %v，正在重连...", err)
                sub.Unsubscribe()
                break  // 跳出内层循环，重新订阅
            case header := <-headers:
                fmt.Printf("新区块: %d\n", header.Number.Uint64())
            }
        }
    }
}
```

### Q4: 订阅和轮询有什么区别？

| 特性 | 订阅 (Subscribe) | 轮询 (Polling) |
|------|------------------|----------------|
| 实时性 | 高（毫秒级） | 低（取决于间隔） |
| 资源消耗 | 低 | 高 |
| 服务器压力 | 小 | 大 |
| 复杂度 | 中 | 低 |
| 适用场景 | 实时应用 | 定期检查 |

---

## 练习作业

开始练习前，请先安装依赖：

```bash
go mod tidy
```

### 作业 1：区块监听器（基础）

练习文件：[exercises/01-basic-listener.go](exercises/01-basic-listener.go)

编写一个程序，实现以下功能：

1. 使用 WebSocket 连接到以太坊测试网
2. 订阅新区块事件
3. 输出每个新区块的基本信息：
   - 区块号
   - 区块哈希
   - 时间戳
   - 交易数量

**运行练习：**
```bash
# 编辑 exercises/01-basic-listener.go，填充 TODO 部分
# 然后运行
go run exercises/01-basic-listener.go
```

**参考答案：** [solutions/01-basic-listener.go](solutions/01-basic-listener.go)

---

### 作业 2：交易分析器（进阶）

练习文件：[exercises/02-tx-analyzer.go](exercises/02-tx-analyzer.go)

编写一个程序，实时分析新区块中的交易：

1. 监听新区块
2. 统计每个区块的：
   - 交易总数
   - 总 Gas 使用量
   - 平均 Gas 价格
3. 显示最近 10 个区块的统计数据

**运行练习：**
```bash
go run exercises/02-tx-analyzer.go
```

**参考答案：** [solutions/02-tx-analyzer.go](solutions/02-tx-analyzer.go)

---

### 作业 3：多链监听器（挑战）

练习文件：[exercises/03-multi-chain.go](exercises/03-multi-chain.go)

编写一个程序，同时监听多个测试网：

1. 同时监听 Sepolia 和 Goerli 测试网
2. 使用 goroutine 并发处理
3. 统一输出新区块信息，包含网络标识
4. 实现优雅关闭

**运行练习：**
```bash
# 需要两个 API Key
export SEPOLIA_WS_URL=wss://...
export GOERLI_WS_URL=wss://...
go run exercises/03-multi-chain.go
```

**参考答案：** [solutions/03-multi-chain.go](solutions/03-multi-chain.go)

---

## WebSocket 节点资源

### 测试网 WebSocket 节点

| 提供商 | Sepolia WebSocket URL | 获取方式 |
|--------|----------------------|----------|
| Alchemy | `wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` | https://www.alchemy.com |
| Infura | `wss://sepolia.infura.io/ws/v3/YOUR_API_KEY` | https://www.infura.io |
| QuickNode | `wss://YOUR_ENDPOINT` | https://www.quicknode.com |
| PublicNode | `wss://ethereum.publicnode.com` | 无需注册，可能不稳定 |

### 主网 WebSocket 节点

| 提供商 | Mainnet WebSocket URL |
|--------|----------------------|
| Alchemy | `wss://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY` |
| Infura | `wss://mainnet.infura.io/ws/v3/YOUR_API_KEY` |

---

## 下一步学习

- [部署合约](../2.10-deploy-contract/)
- [监听合约事件](../2.13-contract-events/)
- [构建实时 DApp](../../advanced/)
