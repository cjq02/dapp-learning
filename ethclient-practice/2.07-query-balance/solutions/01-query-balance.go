// 01-query-balance.go - 查询地址余额 - 答案
//
// 运行：export INFURA_API_KEY=your-key && go run solutions/01-query-balance.go

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

	// 从环境变量读取 API Key
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY\n例如: export INFURA_API_KEY=your-key-here")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 要查询的地址
	address := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	// 查询余额（Wei）
	balanceWei, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 将 Wei 转换为 ETH
	balanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(balanceWei),
		big.NewFloat(1e18),
	)

	// 输出结果
	fmt.Printf("地址: %s\n", address.Hex())
	fmt.Printf("余额: %s Wei\n", balanceWei.String())
	fmt.Printf("余额: %.6f ETH\n", balanceEth)

	fmt.Println("=== 完成 ===")
}
