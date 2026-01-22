// 02-batch-query.go - 批量查询余额 - 答案

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

	// 定义多个地址
	addresses := []common.Address{
		common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), // Vitalik
		common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F"), // 示例地址1
		common.HexToAddress("0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359"), // 示例地址2
	}

	// 账户信息结构
	type AccountInfo struct {
		Address    common.Address
		BalanceWei *big.Int
		BalanceEth *big.Float
	}

	var accounts []AccountInfo
	totalBalanceWei := big.NewInt(0)

	// 遍历地址，查询余额
	for _, addr := range addresses {
		balance, err := client.BalanceAt(context.Background(), addr, nil)
		if err != nil {
			log.Printf("查询 %s 失败: %v\n", addr.Hex(), err)
			continue
		}

		balanceEth := new(big.Float).Quo(
			new(big.Float).SetInt(balance),
			big.NewFloat(1e18),
		)

		accounts = append(accounts, AccountInfo{
			Address:    addr,
			BalanceWei: balance,
			BalanceEth: balanceEth,
		})

		totalBalanceWei.Add(totalBalanceWei, balance)
	}

	// 输出表格
	fmt.Println("\n地址                           余额 (ETH)")
	fmt.Println(strings.Repeat("─", 50))

	for _, acc := range accounts {
		// 缩短地址显示
		shortAddr := acc.Address.Hex()
		if len(shortAddr) > 20 {
			shortAddr = shortAddr[:6] + "..." + shortAddr[len(shortAddr)-4:]
		}
		fmt.Printf("%-30s  %12.6f\n", shortAddr, acc.BalanceEth)
	}

	// 输出总余额
	fmt.Println(strings.Repeat("─", 50))
	totalBalanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(totalBalanceWei),
		big.NewFloat(1e18),
	)
	fmt.Printf("总计: %.6f ETH\n", totalBalanceEth)

	fmt.Println("=== 完成 ===")
}
