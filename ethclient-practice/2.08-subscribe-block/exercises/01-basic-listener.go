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

	// 练习：使用 WebSocket 连接到以太坊节点
	// 提示：使用 ethclient.Dial，注意 URL 以 wss:// 开头
	// var client *ethclient.Client
	// client, err = ???

	// 练习：创建用于接收区块头的通道
	// 提示：使用 make 创建 chan *types.Header
	// var headers chan *types.Header
	// headers = ???

	// 练习：订阅新区块事件
	// 提示：使用 client.SubscribeNewHead
	// var sub ethereum.Subscription
	// sub, err = ???

	fmt.Println("开始监听新区块... (按 Ctrl+C 退出)")

	// 练习：监听新区块事件
	// 提示：使用 select 语句监听 sub.Err() 和 headers 通道
	// for {
	//     select {
	//     case err := <-sub.Err():
	//         log.Fatal(err)
	//     case header := <-headers:
	//         // TODO: 处理新区块
	//         // 1. 打印区块号 header.Number.Uint64()
	//         // 2. 打印区块哈希 header.Hash().Hex()
	//         // 3. 获取完整区块 client.BlockByHash()
	//         // 4. 打印交易数量
	//     }
	// }
}
