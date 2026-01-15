package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// TODO: 替换为你的 API Key
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 练习：获取最新区块头
	// 提示：使用 HeaderByNumber，传入 nil 获取最新区块
	// var header *types.Header
	// header, err = ???

	// 练习：获取完整区块以获取交易数量
	// 提示：使用 BlockByNumber，传入 header.Number
	// var block *types.Block
	// block, err = ???

	// 练习：格式化时间戳为可读格式
	// 提示：使用 time.Unix() 和 Format()
	// var blockTime time.Time
	// var timeFormatted string
	// blockTime = ???
	// timeFormatted = ???

	// 输出区块信息
	fmt.Println("=== 最新区块信息 ===")
	// fmt.Printf("区块号: %d\n", block.Number().Uint64())
	// fmt.Printf("时间: %s\n", timeFormatted)
	// fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	// fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	// fmt.Printf("交易数量: %d\n", len(block.Transactions()))
}
