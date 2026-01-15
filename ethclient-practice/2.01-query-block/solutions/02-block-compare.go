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

	// 获取最新区块号
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := header.Number.Uint64()
	if blockNumber == 0 {
		fmt.Println("Genesis 区块，没有前一个区块")
		return
	}

	// 获取当前区块和前一个区块
	currentBlock, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	prevBlockNumber := big.NewInt(int64(blockNumber - 1))
	prevBlock, err := client.BlockByNumber(context.Background(), prevBlockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 计算变化
	gasUsedChange := int64(currentBlock.GasUsed()) - int64(prevBlock.GasUsed())
	gasUsedPercent := float64(gasUsedChange) / float64(prevBlock.GasUsed()) * 100

	timeDiff := int64(currentBlock.Time()) - int64(prevBlock.Time())

	txCountChange := len(currentBlock.Transactions()) - len(prevBlock.Transactions())
	txCountPercent := float64(txCountChange) / float64(len(prevBlock.Transactions())) * 100

	// 输出对比结果
	fmt.Println("=== 区块对比分析 ===")
	fmt.Printf("区块 N-1: %d → 区块 N: %d\n\n", prevBlock.Number().Uint64(), currentBlock.Number().Uint64())

	fmt.Printf("Gas 使用: %d → %d (%+.2f%%)\n",
		prevBlock.GasUsed(), currentBlock.GasUsed(), gasUsedPercent)
	fmt.Printf("时间间隔: %d 秒\n", timeDiff)

	fmt.Printf("交易数量: %d → %d (%+.2f%%)\n",
		len(prevBlock.Transactions()), len(currentBlock.Transactions()), txCountPercent)

	diffCurrent := currentBlock.Difficulty().Uint64()
	diffPrev := prevBlock.Difficulty().Uint64()
	if diffCurrent == diffPrev {
		fmt.Println("难度: 无变化")
	} else {
		fmt.Printf("难度: %d → %d\n", diffPrev, diffCurrent)
	}
}
