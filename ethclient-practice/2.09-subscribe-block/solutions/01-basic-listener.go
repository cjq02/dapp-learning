package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取 WebSocket URL
	wsURL := os.Getenv("SEPOLIA_WS_URL")
	if wsURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_WS_URL\n" +
			"例如: export SEPOLIA_WS_URL=wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
	}

	// 连接到以太坊 WebSocket 节点
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 创建用于接收区块头的通道
	headers := make(chan *types.Header)

	// 订阅新区块事件
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	fmt.Println("开始监听新区块... (按 Ctrl+C 退出)")

	// 监听新区块事件
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			// 打印区块头信息
			fmt.Printf("\n=== 新区块 ===\n")
			fmt.Printf("区块号: %d\n", header.Number.Uint64())
			fmt.Printf("区块哈希: %s\n", header.Hash().Hex())
			fmt.Printf("时间戳: %d\n", header.Time)

			// 获取完整区块以获取交易数量
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("获取完整区块失败: %v", err)
				continue
			}

			fmt.Printf("交易数量: %d\n", len(block.Transactions()))
			fmt.Printf("Gas 使用: %d\n", block.GasUsed())
			fmt.Printf("Gas 上限: %d\n", block.GasLimit())
		}
	}
}
