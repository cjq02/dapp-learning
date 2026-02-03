# 02-deploy-raw 自测题 — 答案与解析

---

## Q1：纯 ethclient 与 bind 的流程差异

**原题：**  
「纯 ethclient 部署」和 01 的「bind 部署」在流程上有什么不同？这里为什么没有调用 DeployStore，而是用 types.NewContractCreation？

**你的解答：**  
流程上，bind 部署调用生成的 Deploy 函数；纯 ethclient 部署调用 NewContractCreation。因为纯 ethclient 部署不需要通过 abigen 生成 go 文件，而是拿生成的 bin 字节码去解析。

**参考答案与解析：**  
- **流程差异**：bind 方式依赖 abigen 生成的 Go 代码，直接调 `DeployStore(auth, client, "1.0")`，内部帮你拼好 data、构造并发送交易；纯 ethclient 方式**不用** abigen 生成的部署函数，自己用 **bin 字节码** 构造一笔「合约创建交易」，再签名、发送。  
- **为什么没有 DeployStore**：因为 02  deliberately 不依赖 abigen 生成的 `store.go`，所以没有 `DeployStore` 可调；部署就等价于「发一笔 to 为空、data 为字节码的交易」，所以用 **types.NewContractCreation** 手动构造这笔交易。

**判定：** ✅ 正确  
**得分：** 1/1

---

## Q2：data、to、构造函数参数

**原题：**  
`types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data)` 里，**data** 是什么？**to** 在哪（为什么创建合约时没有 to）？若要把构造函数参数（例如 Store 的 version "1.0"）也编码进去，data 应该怎么得到？

**你的解答：**  
data 是 solc 编译生成的合约 bin 字节码。创建合约时不需要 to，to 是零地址。不知道（构造函数参数如何编码）。

**参考答案与解析：**  
- **data**：就是合约的**部署字节码**（solc 的 bin 产物）；对无参或只有默认值的构造函数，data = bin 即可。  
- **to**：合约创建交易的 **to 为空（nil）**，链上约定「to 为空 + data 非空」表示创建合约，所以没有 to 字段，或者说 to 相当于零地址。  
- **构造函数参数**：若构造函数有参数（如 Store 的 `_version`），**data = 字节码 + ABI 编码的构造参数**（即 `bin + abi.encode("1.0")`）。  
  - bind 方式里 `DeployStore(..., "1.0")` 会内部帮你拼好这段 data；  
  - 纯 ethclient 要自己用 abi 包或手写把参数 ABI 编码后拼到 bin 后面，再传给 `NewContractCreation`。  
  当前 02 里的 Store 若用无参或固定版本，可以只传 bin；若传 "1.0"，就需要在 data 里拼上编码后的参数。

**判定：** ⚠️ 部分正确（data、to 对；构造函数参数未答）  
**得分：** 0.6/1

---

## Q3：StoreContractBytecode 来源、DecodeString、为何带 store_bytecode.go

**原题：**  
StoreContractBytecode 是从哪来的？为什么要先 hex.DecodeString(StoreContractBytecode) 再传给 NewContractCreation？运行 02 时为什么要带上 store_bytecode.go？

**你的解答：**  
是通过 solc --bin 命令生成的文件里面的字节码。因为 NewContractCreation 需要的参数类型是 byte[]，store_bytecode 是通用代码，其它文件只需要调整，而不用单独写。

**参考答案与解析：**  
- **来源**：来自 **solc 编译** 得到的 bin 文件（如 `Store_sol_Store.bin`），把里面的十六进制字符串拷到 `store_bytecode.go` 里做成常量 **StoreContractBytecode**。  
- **为何 DecodeString**：`NewContractCreation` 的 data 参数类型是 **[]byte**；Go 里常量是十六进制**字符串**，所以要 **hex.DecodeString** 转成 **[]byte** 再传进去。  
- **为何带 store_bytecode.go**：  
  - **go run 02-deploy-raw.go** 只编译**当前这一个文件**，不会自动带上同目录的 `store_bytecode.go`，所以 02 里用到的 **StoreContractBytecode** 会 **undefined**；  
  - 写成 **go run store_bytecode.go 02-deploy-raw.go** 会把两个文件一起编译成同一个 main 包，02 才能引用到 **StoreContractBytecode**。  
  「通用、其它文件不用单独写」方向对，但直接原因是「单文件 go run 不会包含其它 .go，必须显式带上定义常量的文件」。

**判定：** ⚠️ 部分正确（来源、类型对；带 store_bytecode.go 的原因可再精确）  
**得分：** 0.7/1

---

## Q4：nonce、chainID、EIP155Signer

**原题：**  
这里为什么要先取 nonce 再构造交易？为什么要用 chainID 和 types.NewEIP155Signer(chainID) 签名，而不是别的签名方式？

**你的解答：**  
每次交易都需要先取 nonce，不然会跳号、太小、重复，导致交易失败。用 chainID 是为了防止重放。NewEIP155Signer 不知道。

**参考答案与解析：**  
- **nonce**：每笔交易必须带**当前账户的 nonce**（递增、不重复）。先取再构造，才能保证 nonce 正确；否则会「重复 / 太小 / 跳号」，链上拒收或顺序错乱。  
- **chainID**：用在签名里是为了 **EIP-155 防重放**，同一签名不能在另一条链（不同 chainID）上复用。  
- **NewEIP155Signer**：以太坊当前标准用 **EIP-155 签名**，把 **chainID** 编码进签名；**types.NewEIP155Signer(chainID)** 返回的就是这个 signer，用它对交易做 **SignTx**，节点才会接受。用别的（如 HomesteadSigner）会不符合当前网络要求，交易会被拒。

**判定：** ⚠️ 部分正确（nonce、chainID 防重放对；EIP155Signer 未答）  
**得分：** 0.65/1

---

## Q5：02 与 01 的对比

**原题：**  
和 01-deploy-with-bind.go 相比，02 的优缺点各是什么？什么场景你更倾向用 02 这种写法？

**你的解答：**  
02 的代码写起来比较啰嗦而且相对复杂。学习的场景。

**参考答案与解析：**  
- **02 缺点**：代码多、要自己处理字节码、nonce、构造 data、签名等；没有类型安全的合约调用；若带构造函数参数还要自己 ABI 编码。  
- **02 优点**：不依赖 abigen 生成的 Go 代码；能看清「部署就是一笔 to=nil、data=字节码 的交易」；适合理解底层、写工具或无法用 abigen 的环境。  
- **场景**：**学习底层、理解交易结构**时用 02；**业务开发、生产环境**更推荐 01（bind）。

**判定：** ✅ 正确（优缺点和「学习场景」都到位）  
**得分：** 1/1

---

## 统计表与总分（细粒度）

| 题号 | 主题 | 你的判定 | 得分 |
|------|------|----------|------|
| Q1 | 纯 ethclient vs bind、为何用 NewContractCreation | ✅ 正确 | 1.00/1 |
| Q2 | data / to / 构造函数参数编码 | ⚠️ 部分正确 | 0.60/1 |
| Q3 | StoreContractBytecode 来源、DecodeString、为何带 store_bytecode.go | ⚠️ 部分正确 | 0.70/1 |
| Q4 | nonce、chainID、EIP155Signer | ⚠️ 部分正确 | 0.65/1 |
| Q5 | 02 与 01 对比、适用场景 | ✅ 正确 | 1.00/1 |
| **合计** | — | — | **3.95/5** |

**总分：3.95 / 5（约 79%）**

---

## 小结

- **掌握较好**：纯 ethclient 与 bind 的流程差异、02 的定位与适用场景。  
- **可再抠细**：构造函数参数如何拼进 data（bin + ABI 编码）；单文件 `go run` 为何要显式带 `store_bytecode.go`；EIP155Signer 是当前链上要求的签名方式。
