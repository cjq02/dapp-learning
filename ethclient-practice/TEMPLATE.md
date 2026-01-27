# Ethclient 学习模块创建模板

## 快速创建指令

```
参考 ~/web3/MetaNodeAcademy/Advanced1-backend-upgrade/ethclient实战/2.xx xxx
再参考 2.01和 TEMPLATE.md 创建 2.xx-xxx-name
```

## 标准目录结构

```
2.xx-xxx-name/
├── go.mod
├── xxx-name.md              # 教程（与目录同名）
├── .golangci.yml            # Lint 配置
├── exercises/
│   ├── 01-xxx-basic.go
│   ├── 02-xxx-advanced.go
│   └── 03-xxx-challenge.go
└── solutions/
    ├── 01-xxx-basic.go
    ├── 02-xxx-advanced.go
    └── 03-xxx-challenge.go
```

## go.mod 模板

```go
module github.com/dapp-learning/ethclient/xxx-name

go 1.21

require github.com/ethereum/go-ethereum v1.13.14
```

## xxx-name.md 模板

```markdown
# Ethclient XXX 学习指南

## 目录
- [学习目标](#学习目标)
- [前置条件](#前置条件)
- [核心概念](#核心概念)
- [XXX 方法](#xxx-方法)
- [常见问题](#常见问题)
- [练习作业](#练习作业)

## 学习目标
完成本指南后，你将能够：
- 目标1
- 目标2
- 目标3

## 前置条件
- 前置条件1
- 前置条件2

## 核心概念
[概念图解]

## XXX 方法
[代码示例]

## 练习作业
### 作业 1：XXX（基础）
练习文件：[exercises/01-xxx-basic.go](exercises/01-xxx-basic.go)
**参考答案：** [solutions/01-xxx-basic.go](solutions/01-xxx-basic.go)

### 作业 2：XXX（进阶）
练习文件：[exercises/02-xxx-advanced.go](exercises/02-xxx-advanced.go)
**参考答案：** [solutions/02-xxx-advanced.go](solutions/02-xxx-advanced.go)

### 作业 3：XXX（挑战）
练习文件：[exercises/03-xxx-challenge.go](exercises/03-xxx-challenge.go)
**参考答案：** [solutions/03-xxx-challenge.go](solutions/03-xxx-challenge.go)
```

## 练习文件模板

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
    // Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 练习：TODO 描述
    // 提示：使用 XXX 方法
    // result, err := ???
}
```

## 答案文件模板

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 连接到以太坊节点（替换为你的 API Key）
    // Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
    // Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 完整实现代码
    result, err := client.SomeMethod(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

## .golangci.yml 模板

```yaml
linters:
  disable:
    - errcheck
    - unused
    - gosimple
    - staticcheck

issues:
  exclude-rules:
    - path: exercises/
      linters:
        - errcheck
        - unused
        - gosimple
        - staticcheck

run:
  skip-dirs:
    - solutions
```

## 创建步骤

1. 创建目录：`mkdir -p 2.xx-xxx-name/exercises 2.xx-xxx-name/solutions`
2. 复制模板文件并修改
3. 运行：`cd 2.xx-xxx-name && go mod tidy`

## 注意事项

- 文件名全部小写，用连字符分隔
- MD 文件名与目录名一致（不用 README.md）
- 每个练习有对应的 TODO 提示
- 答案代码完整可运行
- API Key 支持两种格式注释
