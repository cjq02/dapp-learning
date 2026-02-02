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
	// Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
	// Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 指定要查询的地址
	account := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// 查询最新余额（Wei）
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 将余额从 Wei 转换为 ETH
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	// 输出结果
	fmt.Printf("地址: %s\n", account.Hex())
	fmt.Printf("余额: %s Wei\n", balance.String())
	fmt.Printf("余额: %.18f ETH\n", ethValue)
}
