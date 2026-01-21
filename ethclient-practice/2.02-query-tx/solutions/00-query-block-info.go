package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取 API Key
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY\n例如: export INFURA_API_KEY=your-key-here")
	}

	// 连接到 Sepolia 测试网络
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal("连接失败: ", err)
	}
	defer client.Close()

	fmt.Println("成功连接到 Sepolia 测试网络")

	// 指定要查询的区块号
	blockNumber := big.NewInt(10077132)

	// 获取区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal("获取区块失败: ", err)
	}

	// 格式化时间戳为可读格式
	blockTime := time.Unix(int64(block.Time()), 0)
	timeFormatted := blockTime.UTC().Format("2006-01-02 15:04:05 UTC")

	// 输出区块信息
	fmt.Println("\n=== 区块信息 ===")
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("时间戳: %d (%s)\n", block.Time(), timeFormatted)
	fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())

	// Gas 信息
	fmt.Printf("Gas 限制: %d\n", block.GasLimit())
	fmt.Printf("Gas 使用: %d (%.2f%%)\n",
		block.GasUsed(),
		float64(block.GasUsed())/float64(block.GasLimit())*100)

	// 交易信息
	txCount := len(block.Transactions())
	fmt.Printf("交易数量: %d\n", txCount)

	// 难度信息
	fmt.Printf("难度: %s\n", block.Difficulty().String())

	// 区块大小
	fmt.Printf("区块大小: %d bytes\n", block.Size())

	// 其他信息
	fmt.Printf("状态根: %s\n", block.Root().Hex())
	fmt.Printf("交易根: %s\n", block.TxHash().Hex())
	fmt.Printf("收据根: %s\n", block.ReceiptHash().Hex())

	// 如果有基础费用（EIP-1559 后）
	if block.BaseFee() != nil {
		baseFeeGwei := new(big.Float).Quo(
			new(big.Float).SetInt(block.BaseFee()),
			big.NewFloat(1e9),
		)
		fmt.Printf("基础费用: %.2f Gwei (%s Wei)\n", baseFeeGwei, block.BaseFee().String())
	}
}
