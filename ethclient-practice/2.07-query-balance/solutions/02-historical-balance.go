package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	account := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// 查询当前区块余额
	currentBalance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 查询历史区块余额（区块号 5532993）
	blockNumber := big.NewInt(5532993)
	historicalBalance, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 转换为 ETH 并计算变化
	currentEth := new(big.Float).Quo(
		new(big.Float).SetInt(currentBalance),
		big.NewFloat(math.Pow10(18)),
	)
	historicalEth := new(big.Float).Quo(
		new(big.Float).SetInt(historicalBalance),
		big.NewFloat(math.Pow10(18)),
	)

	change := new(big.Float).Sub(currentEth, historicalEth)

	// 输出结果
	fmt.Printf("地址: %s\n", account.Hex())
	fmt.Printf("当前余额: %.18f ETH\n", currentEth)
	fmt.Printf("历史余额 (区块 %d): %.18f ETH\n", blockNumber, historicalEth)
	fmt.Printf("余额变化: %.18f ETH\n", change)
}
