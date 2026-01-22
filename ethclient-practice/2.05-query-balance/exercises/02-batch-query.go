// 02-batch-query.go - 批量查询余额练习
//
// 任务：
// 1. 定义多个地址（至少 3 个）
// 2. 批量查询每个地址的余额
// 3. 统计总余额
// 4. 输出格式化的表格
//
// 运行：export INFURA_API_KEY=your-key && go run exercises/02-batch-query.go

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 批量查询账户余额 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// TODO 1: 定义多个地址
	addresses := []common.Address{
		common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), // Vitalik
		// 添加更多地址...
	}

	// TODO 2: 批量查询余额并累加
	type AccountInfo struct {
		Address  common.Address
		Balance  *big.Int
		BalanceEth *big.Float
	}

	var accounts []AccountInfo
	totalBalanceWei := big.NewInt(0)

	// TODO 3: 遍历地址，查询余额
	// 提示：使用 for range 循环
	// 提示：使用 client.BalanceAt() 查询
	// 提示：将结果添加到 accounts 切片
	// 提示：累加到 totalBalanceWei

	// TODO 4: 输出表格
	fmt.Println("\n地址                           余额 (ETH)")
	fmt.Println(strings.Repeat("─", 50))

	// 提示：遍历 accounts，格式化输出每个地址的余额

	// TODO 5: 输出总余额
	totalBalanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(totalBalanceWei),
		big.NewFloat(1e18),
	)
	fmt.Printf("\n总计: %.6f ETH\n", totalBalanceEth)

	fmt.Println("=== 完成 ===")
}
