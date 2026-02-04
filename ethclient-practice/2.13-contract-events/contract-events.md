# Ethclient 合约事件学习指南

> **预计学习时间：** 40 分钟
>
> **难度：** 进阶

本指南介绍如何使用 Go 语言的 `go-ethereum` 库查询和订阅智能合约事件。

## 目录

- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [查询事件](#查询事件)
- [订阅事件](#订阅事件)
- [从交易收据获取事件](#从交易收据获取事件)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

---

## 学习目标

完成本指南后，你将能够：
- 理解以太坊事件（日志）的结构与 Topics
- 使用 `FilterQuery` 和 `FilterLogs` 查询历史事件
- 使用 `SubscribeFilterLogs` 订阅实时事件
- 从交易收据的 Logs 中解析合约事件
- 使用 ABI 解码事件数据

## 前置条件

- Go 语言基础
- Solidity 基础（了解 event、indexed）
- 已完成 [加载合约](../2.11-load-contract/) 或 [调用合约](../2.12-call-contract/)
- 已部署 Store 合约（可参考 [部署合约](../2.10-deploy-contract/)）

## 核心概念

### 事件即日志

智能合约在执行时可以“发出”事件，事件在以太坊中存储为**日志**，是交易收据的一部分：

```
交易 (Transaction)
└── 收据 (Receipt)
    └── Logs[]  ← 事件列表
        ├── Address   (合约地址)
        ├── Topics[]  (主题，最多 4 个)
        └── Data      (非 indexed 参数编码)
```

### Topics 说明

| Topic       | 说明 |
|------------|------|
| `topics[0]` | 事件签名哈希：`keccak256("EventName(type1,type2,...)")`，必定存在 |
| `topics[1]` | 第 1 个 `indexed` 参数的值（若有） |
| `topics[2]` | 第 2 个 `indexed` 参数的值（若有） |
| `topics[3]` | 第 3 个 `indexed` 参数的值（若有） |

**注意：** 被 `indexed` 修饰的字段不会出现在 `Data` 中，仅出现在 Topics 里。

### Store 合约事件

本模块使用与 [2.10 部署合约](../2.10-deploy-contract/) 相同的 Store 合约：

```solidity
event ItemSet(bytes32 indexed key, bytes32 value);
```

- `key`：indexed，出现在 `topics[1]`
- `value`：非 indexed，出现在 `Data` 中

---

## 查询事件

### FilterQuery

构造过滤条件：

```go
query := ethereum.FilterQuery{
    FromBlock: big.NewInt(6920583),  // 起始区块
    ToBlock:   big.NewInt(6920600),  // 结束区块，可选
    Addresses: []common.Address{contractAddress},
    Topics:    [][]common.Hash{},    // 可选，按 topic 过滤
}
```

### FilterLogs

```go
logs, err := client.FilterLogs(context.Background(), query)
if err != nil {
    log.Fatal(err)
}
```

### 解码事件

使用合约 ABI 解码 `vLog.Data`：

```go
contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
// ...
event := struct {
    Key   [32]byte
    Value [32]byte
}{}
err = contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
```

---

## 订阅事件

订阅需要 **WebSocket RPC**（与订阅区块相同）。

```go
client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
```

```go
logsCh := make(chan types.Log)
sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
if err != nil {
    log.Fatal(err)
}
defer sub.Unsubscribe()

for {
    select {
    case err := <-sub.Err():
        log.Fatal(err)
    case vLog := <-logsCh:
        // 处理新事件
    }
}
```

---

## 从交易收据获取事件

除 FilterLogs 和 SubscribeFilterLogs 外，也可从单笔交易的收据中获取事件：

```go
receipt, err := client.TransactionReceipt(context.Background(), txHash)
if err != nil {
    log.Fatal(err)
}
for _, vLog := range receipt.Logs {
    // 解析 vLog
}
```

---

## 常见问题

### Q1: 为什么需要 WebSocket 才能订阅事件？

HTTP RPC 只能请求-响应，无法持续推送。订阅是长连接推送模型，必须使用 WebSocket（或 IPC）。

### Q2: Topics 过滤怎么写？

若只关心 `ItemSet` 事件，可设置：

```go
eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))
query.Topics = [][]common.Hash{{eventSig}}
```

### Q3: 如何按 indexed 参数过滤？

将对应 topic 填入 `query.Topics` 的对应位置，例如按 `key` 过滤时，`Topics[1]` 为 `[]common.Hash{keyHash}`。

---

## 练习作业

练习使用**已部署的 Store 合约**。若尚未部署，请先完成 [2.10 部署合约](../2.10-deploy-contract/)，并准备好合约地址与 RPC URL。

### 作业 1：查询历史事件（基础）

练习文件：[exercises/01-query-events.go](exercises/01-query-events.go)

使用 `FilterQuery` 和 `FilterLogs` 查询指定区块范围内某合约的 ItemSet 事件，并解码输出 key、value 与 topics。

**提示：** 使用 `ethereum.FilterQuery`、`client.FilterLogs`、`abi.JSON`、`UnpackIntoInterface`。

**参考答案：** [solutions/01-query-events.go](solutions/01-query-events.go)

---

### 作业 2：订阅事件（进阶）

练习文件：[exercises/02-subscribe-events.go](exercises/02-subscribe-events.go)

使用 **WebSocket** 连接和 `SubscribeFilterLogs` 订阅指定合约的日志，收到 ItemSet 事件时解码并打印。

**提示：** RPC URL 需为 `wss://...`，使用 `SubscribeFilterLogs` 和 channel 循环读取。

**参考答案：** [solutions/02-subscribe-events.go](solutions/02-subscribe-events.go)

---

### 作业 3：从收据解析事件（挑战）

练习文件：[exercises/03-events-from-receipt.go](exercises/03-events-from-receipt.go)

给定一笔交易哈希，通过 `TransactionReceipt` 获取收据，遍历 `receipt.Logs`，识别并解码 ItemSet 事件（可比较 `log.Topics[0]` 与事件签名哈希）。

**参考答案：** [solutions/03-events-from-receipt.go](solutions/03-events-from-receipt.go)

---

## 测试网资源

| 提供商 | HTTP RPC | WebSocket |
|--------|----------|-----------|
| Alchemy | `https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` | `wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY` |
| Infura   | `https://sepolia.infura.io/v3/YOUR_API_KEY`         | `wss://sepolia.infura.io/ws/v3/YOUR_API_KEY`       |

---

## 下一步学习

- [订阅区块](../2.09-subscribe-block/)：订阅新区块与区块头
- [查询收据](../2.03-query-receipt/)：收据结构与状态
