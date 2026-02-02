package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Network 网络配置
type Network struct {
	Name string
	URL  string
}

// BlockInfo 区块信息
type BlockInfo struct {
	Network string
	Number  uint64
	Hash    string
}

func main() {
	// 从环境变量读取 WebSocket URL
	networks := []Network{
		{Name: "Sepolia", URL: os.Getenv("SEPOLIA_WS_URL")},
		{Name: "Goerli", URL: os.Getenv("GOERLI_WS_URL")},
	}

	// 检查环境变量
	for _, n := range networks {
		if n.URL == "" {
			log.Fatalf("错误: 请设置环境变量 %s_WS_URL", n.Name)
		}
	}

	// 用于接收区块信息的通道
	blockCh := make(chan BlockInfo, 100)
	var wg sync.WaitGroup

	// 练习：为每个网络启动一个监听 goroutine
	// for _, network := range networks {
	//     wg.Add(1)
	//     go func(net Network) {
	//         defer wg.Done()
	//
	//         // TODO: 在 goroutine 中:
	//         // 1. 连接到 WebSocket
	//         // 2. 订阅新区块
	//         // 3. 将区块信息发送到 blockCh
	//
	//     }(network)
	// }

	// 练习：在另一个 goroutine 中处理和打印区块信息
	// go func() {
	//     for block := range blockCh {
	//         // TODO: 打印区块信息，包含网络名称
	//     }
	// }()

	// 练习：实现优雅关闭
	// 提示：监听 SIGINT 和 SIGTERM 信号
	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//
	// <-sigCh
	// fmt.Println("\n收到退出信号，正在关闭...")
	// close(blockCh)
	// wg.Wait()
	// fmt.Println("已关闭所有订阅")
}

// subscribeToNetwork 订阅指定网络的新区块
func subscribeToNetwork(network Network, blockCh chan<- BlockInfo) error {
	// TODO: 实现订阅逻辑
	// 提示：
	// 1. 使用 ethclient.Dial 连接
	// 2. 创建 headers 通道
	// 3. 使用 client.SubscribeNewHead 订阅
	// 4. 在 for 循环中监听事件
	// 5. 将 BlockInfo 发送到 blockCh
	//
	// BlockInfo 结构:
	// BlockInfo{
	//     Network: network.Name,
	//     Number:  header.Number.Uint64(),
	//     Hash:    header.Hash().Hex(),
	// }

	return nil
}
