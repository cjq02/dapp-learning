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

	// 为每个网络启动一个监听 goroutine
	for _, network := range networks {
		wg.Add(1)
		go func(net Network) {
			defer wg.Done()
			if err := subscribeToNetwork(net, blockCh); err != nil {
				log.Printf("订阅 %s 失败: %v", net.Name, err)
			}
		}(network)
	}

	// 在另一个 goroutine 中处理和打印区块信息
	go func() {
		for block := range blockCh {
			fmt.Printf("[%s] 区块 #%d: %s\n", block.Network, block.Number, block.Hash)
		}
	}()

	// 等待退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	fmt.Println("\n收到退出信号，正在关闭...")
	close(blockCh)
	wg.Wait()
	fmt.Println("已关闭所有订阅")
}

// subscribeToNetwork 订阅指定网络的新区块
func subscribeToNetwork(network Network, blockCh chan<- BlockInfo) error {
	// 连接到 WebSocket 节点
	client, err := ethclient.Dial(network.URL)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	defer client.Close()

	// 创建订阅通道
	headers := make(chan *types.Header)

	// 订阅新区块
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("订阅失败: %w", err)
	}
	defer sub.Unsubscribe()

	log.Printf("[%s] 订阅成功，开始监听...", network.Name)

	// 监听新区块事件
	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("订阅错误: %w", err)
		case header := <-headers:
			// 发送区块信息到通道
			blockCh <- BlockInfo{
				Network: network.Name,
				Number:  header.Number.Uint64(),
				Hash:    header.Hash().Hex(),
			}
		}
	}
}
