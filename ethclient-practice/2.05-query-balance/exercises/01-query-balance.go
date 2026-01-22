// 01-query-balance.go - 查询地址余额练习
//
// 任务：
// 1. 连接到 Sepolia 测试网
// 2. 查询指定地址的 ETH 余额
// 3. 将余额从 Wei 转换为 ETH
// 4. 输出结果
//
// 运行：export INFURA_API_KEY=your-key && go run exercises/01-query-balance.go

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 查询账户余额 ===")

	// TODO 1: 从环境变量读取 API Key
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY\n例如: export INFURA_API_KEY=your-key-here")
	}

	// TODO 2: 连接到以太坊节点
	var client *ethclient.Client
	{
		// 在这里填写代码
		// 提示：使用 ethclient.Dial()
		// Sepolia RPC URL: https://sepolia.infura.io/v3/ + apiKey
	}

	defer func() {
		// TODO: 关闭客户端连接
	}()

	// 要查询的地址（Vitalik 的地址，示例）
	address := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	// TODO 3: 查询余额（Wei）
	var balanceWei *big.Int
	{
		// 在这里填写代码
		// 提示：使用 client.BalanceAt(context.Background(), address, nil)
		// nil 表示查询最新区块
	}

	// TODO 4: 将 Wei 转换为 ETH
	var balanceEth *big.Float
	{
		// 在这里填写代码
		// 提示：使用 big.Float 进行除法运算
		// balanceEth = balanceWei / 1e18
		// 使用 new(big.Float).Quo(被除数, 除数)
	}

	// TODO 5: 输出结果
	fmt.Printf("地址: %s\n", address.Hex())
	fmt.Printf("余额: %s Wei\n", balanceWei.String())
	fmt.Printf("余额: %.6f ETH\n", balanceEth)

	fmt.Println("=== 完成 ===")
}
