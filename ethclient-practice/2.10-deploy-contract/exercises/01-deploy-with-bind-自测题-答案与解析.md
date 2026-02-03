# 01-deploy-with-bind 自测题 — 答案与解析

---

## Q1：为什么必须用 ChainID 才能创建 auth？

**原题：**  
代码里先连节点、再加载私钥、再拿 ChainID、再创建 auth。为什么必须用 **ChainID** 才能创建 auth？如果跳过 ChainID 直接用私钥发交易会怎样？

**你的解答：**  
如果不用 chainId，会拿不到 auth，代码会报错。

**参考答案与解析：**  
- `bind.NewKeyedTransactorWithChainID(privateKey, chainID)` 的第二个参数是**必填**的，不传 chainID 就无法调用这个 API，代码会报错，所以「拿不到 auth、代码会报错」是对的。  
- 更完整一点：ChainID 用于 **EIP-155 签名**（防重放），不同链的 chainID 不同，签名时带上 chainID 才能保证交易只在当前链上有效；若用老的不带 chainID 的构造方式（已废弃）发交易，在主网/测试网上会有安全与兼容问题。

**判定：** ✅ 正确（结论对；可补充：ChainID 用于 EIP-155 签名防重放）

**得分：** 1/1

---

## Q2：GasLimit 和 GasPrice 分别表示什么？为什么说设大 GasLimit 不会多付钱？

**原题：**  
`auth.GasLimit = 2_000_000` 和 `auth.GasPrice = gasPrice` 分别表示什么？为什么文档里说「设大 GasLimit 不会多付钱」？实际扣费公式是什么？

**你的解答：**  
gasLimit 规定了部署合约 gas 的使用上限，gasPrice 预估了部署合约需要消耗多少。

**参考答案与解析：**  
- **GasLimit**：表示这笔交易**最多允许消耗的 gas 数量**（上限），你理解成「使用上限」是对的。  
- **GasPrice**：表示**每单位 gas 的单价**（单位：wei），不是「预估消耗多少」；预估消耗多少是 **gasUsed**（实际用掉的 gas）。  
- **实际扣费公式**：`手续费 = gasPrice × gasUsed`（按实际消耗扣费），所以把 GasLimit 设大只是放宽上限，**不会多付钱**，只会避免因上限不够而 out of gas。

**判定：** ⚠️ 部分正确（GasLimit 对；GasPrice 理解有误，应为「每单位 gas 的单价」而非「预估消耗」）

**得分：** 0.5/1

---

## Q3："1.0" 是什么？contractAddr、tx、instance 分别是什么？部署成功指什么？

**原题：**  
`store.DeployStore(auth, client, "1.0")` 里的 **"1.0"** 是什么？返回的 `contractAddr`、`tx`、`instance` 分别有什么用途？部署「成功」是指 RPC 返回了，还是指链上已经确认？

**你的解答：**  
1.0 指的是部署合约的版本号。返回的 contractAddr 指的是合约的地址，tx 指的是交易 hash，instance 指的是实例。部署成功指的是链上已经确认。

**参考答案与解析：**  
- **"1.0"**：是合约**构造函数的参数** `_version`（Store 合约里 `constructor(string memory _version)`），合约会把它存成 `version` 状态，语义上可以当「版本号」理解，没问题。  
- **contractAddr**：合约部署后的**链上地址**，对。  
- **tx**：是**整笔交易对象**（*types.Transaction），**tx.Hash()** 才是交易哈希；说「tx 是交易 hash」不够准确。  
- **instance**：部署后得到的**合约实例**（可用于后续调用合约方法），对。  
- **部署成功**：你说「链上已经确认」是**理想意义上的成功**（即 receipt 已上链且 Status==1）；当前代码里是在 DeployStore 返回后就打印「合约部署成功！」，那只代表**交易已提交到节点**，尚未等待链上确认，所以脚本里的「成功」和「链上确认」要区分开。

**判定：** ⚠️ 部分正确（1.0、contractAddr、instance、部署成功的含义都对；tx 应为「交易对象」，tx.Hash() 才是哈希）

**得分：** 0.5/1

---

## Q4：什么时候 TransactionReceipt 会报错？怎样稳定等到链上确认？

**原题：**  
代码里在 DeployStore 返回后立刻调用了 `client.TransactionReceipt(context.Background(), tx.Hash())`。在什么情况下这里会报错（例如 not found）？要怎样改才能稳定等到「链上确认」再判断成功？

**你的解答：**  
合约部署需要时间来确认。需要轮询。

**参考答案与解析：**  
- **何时报错**：交易刚提交时还在 mempool，尚未被打包进区块，此时 `TransactionReceipt` 会返回 **not found**（或类似错误），所以「需要时间确认」是对的。  
- **怎样稳定等到链上确认**：需要**轮询** `TransactionReceipt`，直到拿到 receipt（或超时）；项目中已有 `util.WaitForReceipt` 就是做这件事的，部署后应用它等 receipt，再根据 `receipt.Status == 1` 判断是否真正成功。

**判定：** ✅ 正确（原因和「需要轮询」都对）

**得分：** 1/1

---

## Q5：store 包从哪来？从 Solidity 到可运行 Go 要几步？bind 与纯 ethclient 部署的对比？

**原题：**  
`store.DeployStore` 和 `store` 包是从哪来的？要得到这个包，在写这段 Go 之前需要做哪几步（从 Solidity 源码到可运行的 Go）？用 bind 部署和「纯 ethclient + 字节码」部署相比，各有什么优缺点？

**你的解答：**  
需要执行 solc 命令得到 bin 和 abi，再通过 abigen 生成 go 代码。bind 部署代码简洁，安全性高，适合在生产环境使用。

**参考答案与解析：**  
- **store 包来源**：由 **abigen** 根据合约的 **bin（字节码）** 和 **abi** 生成的 Go 包（即 contract/store.go）。  
- **从 Solidity 到可运行 Go 的步骤**：① 用 **solc**（或 solcjs）编译 Solidity 得到 **.bin** 和 **.abi**；② 用 **abigen** 指定 bin、abi、pkg、out 生成 **store.go**；③ 在 Go 里 import 该包并调用 `DeployStore` 等。你写的「solc 得到 bin 和 abi，abigen 生成 go 代码」正确。  
- **bind vs 纯 ethclient**：bind 方式代码简洁、类型安全、后续调用合约方便，适合生产；纯 ethclient 要手写字节码、自己构造交易，适合学习底层。你说的「bind 代码简洁、安全性高、适合生产」都对。

**判定：** ✅ 正确

**得分：** 1/1

---

## 统计表与总分

| 题号 | 主题 | 你的判定 | 得分 |
|------|------|----------|------|
| Q1 | ChainID 与 auth | ✅ 正确 | 1/1 |
| Q2 | GasLimit / GasPrice / 扣费 | ⚠️ 部分正确 | 0.5/1 |
| Q3 | 1.0、返回值、部署成功含义 | ⚠️ 部分正确 | 0.5/1 |
| Q4 | Receipt 报错与轮询确认 | ✅ 正确 | 1/1 |
| Q5 | store 来源、工具链、bind 对比 | ✅ 正确 | 1/1 |
| **合计** | — | — | **4/5** |

**总分：4 / 5（80%）**

---

## 小结

- **掌握较好**：ChainID 与 auth、轮询确认、工具链与 bind 对比。  
- **可再抠细一点**：GasPrice 是「每单位 gas 单价」不是「预估消耗」；tx 是交易对象、tx.Hash() 才是哈希；代码里「部署成功」打印时机 vs「链上确认」的区别。
