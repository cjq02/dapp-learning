package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点（请替换为你的 API Key）
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 获取最新区块头
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 获取完整区块以获取交易数量
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	// 格式化时间戳
	blockTime := time.Unix(int64(block.Time()), 0)
	timeFormatted := blockTime.UTC().Format("2006-01-02 15:04:05 UTC")

	// 输出区块信息
	fmt.Println("=== 最新区块信息 ===")
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("时间: %s\n", timeFormatted)
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
}
