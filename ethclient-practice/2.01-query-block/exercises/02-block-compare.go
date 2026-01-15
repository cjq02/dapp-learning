package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
	// Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 练习：获取最新区块号
	// header, err := ???

	// 练习：获取区块号并检查是否为创世区块
	// blockNumber := header.Number.Uint64()
	// if blockNumber == 0 { ... }

	// 练习：获取当前区块和前一个区块
	// currentBlock, err := client.BlockByNumber(???, ???)
	// prevBlockNumber := big.NewInt(int64(blockNumber - 1))
	// prevBlock, err := client.BlockByNumber(???, ???)

	// 练习：计算 Gas 使用量变化和百分比
	// gasUsedChange := ???
	// gasUsedPercent := ???

	// 练习：计算时间间隔
	// timeDiff := ???

	// 练习：计算交易数量变化和百分比
	// txCountChange := ???
	// txCountPercent := ???

	// 输出对比结果（已实现）
	fmt.Println("=== 区块对比分析 ===")
	// fmt.Printf("区块 N-1: %d → 区块 N: %d\n\n", ..., ...)

	// fmt.Printf("Gas 使用: %d → %d (%+.2f%%)\n", ..., ..., ...)
	// fmt.Printf("时间间隔: %d 秒\n", ...)
	// fmt.Printf("交易数量: %d → %d (%+.2f%%)\n", ..., ..., ...)

	// 练习：判断难度是否变化
	// if ... { ... } else { ... }
}
