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

**计算事件签名时用的是什么字符串？**

对 `topics[0]` 做哈希时，用的是事件的**规范原型**字符串，规则与函数选择器一致：

- **格式**：`EventName(type1,type2,...)`，括号里只写**参数类型**，类型之间用英文逗号分隔，**无空格**。
- **不包含**：参数名、`indexed` 关键字都不写。
- **类型用规范名**：`uint` → `uint256`，`int` → `int256`，`byte` → `bytes1` 等（与 [ABI 规范](https://docs.soliditylang.org/en/latest/abi-spec.html) 一致）。

例如 Solidity 里定义为 `event ItemSet(bytes32 indexed key, bytes32 value)`，用来算哈希的字符串是 **`ItemSet(bytes32,bytes32)`**，而不是带 `key`/`value` 或 `indexed` 的写法。再如 `event Transfer(address indexed from, address indexed to, uint256 value)` 对应 **`Transfer(address,address,uint256)`**。

**为什么是这样的格式？**

- **和函数选择器一致**：以太坊 ABI 里，函数用 `keccak256("函数名(类型列表)")` 的前 4 字节做 selector；事件用同一套“规范原型”的哈希做 `topics[0]`，这样合约的“对外形状”用同一规则描述，工具链也统一。
- **只认“形状”，不认名字**：参数名、`indexed` 只影响可读性和存储方式，不改变“哪个事件、参数类型是什么”。签名只描述“事件名 + 参数类型列表”，这样同一事件在不同合约、不同命名下仍得到相同 topic，便于过滤和跨合约识别。
- **规范类型名保证唯一、可互操作**：用规范名（如 `uint256`）而不是别名（如 `uint`），保证所有语言、编译器算出的字符串一致，哈希才能对齐；参数名各语言可能不同，不参与签名更稳妥。

**注意：** 被 `indexed` 修饰的字段不会出现在 `Data` 中，仅出现在 Topics 里。

以上是**普通事件**的规则。声明为 `anonymous` 的**匿名事件**不把事件签名写入日志，布局不同，见下节。

### 匿名事件与 Topics

Solidity 中可用 `anonymous` 修饰事件：

```solidity
event ItemSet(bytes32 indexed key, bytes32 value) anonymous;
```

匿名事件在日志中的行为：

- **不**把事件签名的哈希写入任何 topic，相当于省掉一个 topic。
- `topics[0]` **不是**事件签名，而是第 1 个 `indexed` 参数。
- 因此最多可以有 **4 个** `indexed` 参数（4 个 topic 全用于参数）；普通事件只有 3 个（1 个给签名 + 3 个给参数）。

| 类型       | topics[0]            | topics[1..3]              | indexed 数量 |
|------------|----------------------|---------------------------|--------------|
| 普通事件   | 事件签名哈希（固定） | 第 1～3 个 indexed 参数   | 最多 3 个    |
| 匿名事件   | 第 1 个 indexed 参数 | 第 2～4 个 indexed 参数   | 最多 4 个    |

**匿名事件的特点：**

- **优点**：少 32 字节（不存签名）、多一个 indexed 槽位，适合对 gas 敏感或需要 4 个 indexed 的场景。
- **缺点**：链上无法通过“事件签名”区分事件类型，只能结合合约地址和 ABI 解析，按事件类型过滤会麻烦一些。

### Data 说明

日志中的 **Data** 字段存放的是事件里**未**用 `indexed` 修饰的参数，按 ABI 编码规则顺序拼接成一段字节（与合约调用的参数编码方式一致）。

- **谁进 Data**：只有**非 indexed** 参数会进入 Data；所有 `indexed` 参数都在 Topics 里，不会出现在 Data 中。
- **编码方式**：按事件参数声明顺序，对非 indexed 参数做 ABI 编码（与 [Contract ABI](https://docs.soliditylang.org/en/latest/abi-spec.html) 一致）。例如多个 `uint256`、`bytes32`、动态类型等，都按标准规则编码。
- **空 Data**：若事件没有任何非 indexed 参数，Data 即为 `0x`（空）。
- **解析**：链上只存原始字节，解析时需要合约 ABI 才能把 Data 解码成具体类型和字段；`go-ethereum` 的 `abi.UnpackLog` 等即根据事件 ABI 解码 Data（以及从 Topics 解析 indexed 参数）。

**Data 里通常放哪些数据？**

- **不需要按条件过滤、只用来“读”的值**：例如 key-value 里的 value、金额、数量、状态码等。这些用 Topics 存会占掉宝贵的 indexed 槽位，且过滤时很少按 value 查，放在 Data 即可。
- **变长或复杂类型**：`string`、`bytes`、`bytes32` 以外的字节、数组、结构体等。Topics 每个只能是 32 字节，变长类型只能放在 Data 里做 ABI 编码。
- **内容较多、仅作记录的信息**：描述、备注、URL、签名等，只在对某条日志做解析时需要，不需要被 `eth_getLogs` 按 topic 过滤。

反之，**需要按条件筛选**的（例如“某地址发的”“某 ID 的”“某 key 的”）通常放在 **indexed** 里进 Topics，这样 RPC 可按 `topics[1..3]` 高效过滤。

以 `event ItemSet(bytes32 indexed key, bytes32 value)` 为例：`key` 在 `topics[1]` 便于按 key 查；`value` 在 Data 里，Data 即为单个 `bytes32` 的 ABI 编码（32 字节）。

### Store 合约事件

本模块使用与 [2.10 部署合约](../2.10-deploy-contract/) 相同的 Store 合约：

```solidity
event ItemSet(bytes32 indexed key, bytes32 value);
```

- `key`：indexed，出现在 `topics[1]`
- `value`：非 indexed，出现在 `Data` 中

### 怎么判断一条 log 是哪个合约、哪个事件？

- **哪个合约**：看 **`log.Address`**（在 `types.Log` 里即 `vLog.Address`）。每条日志都带发出它的合约地址，直接比较即可，例如 `vLog.Address == contractAddr`。
- **哪个事件（普通事件）**：看 **`log.Topics[0]`**。普通事件会把事件签名的哈希放在第一个 topic，即 `topics[0] == keccak256("EventName(type1,type2,...)")`。在 Go 里可预先算好再比较：
  ```go
  eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))
  if len(vLog.Topics) > 0 && vLog.Topics[0] == eventSig {
      // 是 ItemSet 事件
  }
  ```
  签名字符串要和 Solidity 定义一致（类型用规范名，如 `uint256` 不用 `uint`）。
- **匿名事件**：没有事件签名 topic，无法单从日志区分是哪种事件。只能结合 **合约地址**（先确定是哪个合约）、以及业务上该合约可能发出的匿名事件来试解码，或用 `abi.ParseTopics` 按某事件的 indexed 参数去匹配。

**典型写法**：先按 `vLog.Address` 筛出目标合约的 log，再按 `vLog.Topics[0]` 判断事件类型（非匿名时），最后用对应 ABI 解码 Data/Topics。

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

**区块范围默认值：**

- **不设置 `FromBlock`**：从创世区块（0）开始查。
- **不设置 `ToBlock`**：查到**当前最新区块**（执行 `FilterLogs` 时节点认为的 latest 区块）。

因此若不设 `ToBlock`，每次查询的“终点”是当时的最新块，适合“从某块到现在”的开放式区间。范围两端都是**含边界**的（inclusive）。

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

**`UnpackIntoInterface(&event, "ItemSet", vLog.Data)` 是怎么工作的？**

1. **三个参数**：目标结构体指针 `&event`、事件名 `"ItemSet"`、原始字节 `vLog.Data`（仅包含日志的 Data 字段，不包含 Topics）。
2. **按名字找事件**：在 `contractAbi` 里根据 `"ItemSet"` 查到该事件的 ABI 定义（参数列表、谁 indexed、谁非 indexed）。
3. **只解 Data 段**：第三个参数只传了 `vLog.Data`，所以库只会用这段字节做 ABI 解码。事件里只有**非 indexed** 参数会编码在 Data 里，因此内部等价于用该事件的「非 indexed 参数」类型列表去 `Unpack(data)`，再按**参数名**或顺序填到 `event` 的对应字段里。
4. **结果**：  
   - **非 indexed 参数**（如 `value`）会从 `vLog.Data` 解出并填到 `event.Value`。  
   - **indexed 参数**（如 `key`）在链上存在 `vLog.Topics` 里，不在 Data 里，这次调用**不会**填 `event.Key`；若需要 `key`，要自己从 `vLog.Topics[1]` 取，或用 `abi.ParseTopics(event.Inputs, vLog.Topics)` 等按 ABI 解析 Topics。

因此：**只传 `vLog.Data` 时，UnpackIntoInterface 只负责把 Data 按 ABI 解成结构体里“对应非 indexed 参数”的字段；indexed 参数需从 Topics 另算。**

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

**环境变量（与 [2.12 调用合约](../2.12-call-contract/) 一致）：** `CONTRACT_ADDRESS`、`SEPOLIA_RPC_URL`；作业 2 需 `SEPOLIA_WS_URL`（WebSocket）；作业 3 需 `TX_HASH`（一笔 setItem 交易哈希）。

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
