package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockStats 区块统计信息
type BlockStats struct {
	Number       uint64
	Hash         string
	TxCount      int
	TotalGasUsed uint64
	AvgGasPrice  *big.Int
}

func main() {
	wsURL := os.Getenv("SEPOLIA_WS_URL")
	if wsURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_WS_URL")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 创建通道并订阅
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// 用于存储最近的区块统计（最多10个）
	var stats []BlockStats
	var mu sync.Mutex

	fmt.Println("开始监听并分析新区块...")
	fmt.Println("================================")

	// 监听新区块并分析
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			// 获取完整区块
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("获取区块失败: %v", err)
				continue
			}

			// 计算统计数据
			txCount := len(block.Transactions())
			totalGasUsed := block.GasUsed()
			avgGasPrice := calculateAverageGasPrice(block)

			// 创建统计信息
			stat := BlockStats{
				Number:       block.Number().Uint64(),
				Hash:         block.Hash().Hex(),
				TxCount:      txCount,
				TotalGasUsed: totalGasUsed,
				AvgGasPrice:  avgGasPrice,
			}

			// 更新 stats 列表（保持最多10个）
			mu.Lock()
			stats = append(stats, stat)
			if len(stats) > 10 {
				stats = stats[1:]
			}
			mu.Unlock()

			// 打印当前区块的统计信息
			fmt.Printf("\n区块 #%d\n", block.Number().Uint64())
			fmt.Printf("  交易数量: %d\n", txCount)
			fmt.Printf("  总 Gas 使用: %d\n", totalGasUsed)
			if avgGasPrice != nil {
				fmt.Printf("  平均 Gas 价格: %s Gwei\n", avgGasPrice.String())
			}
			fmt.Printf("  区块哈希: %s\n", block.Hash().Hex())

			// 打印统计列表
			mu.Lock()
			printStats(stats)
			mu.Unlock()
		}
	}
}

// calculateAverageGasPrice 计算平均 Gas 价格
func calculateAverageGasPrice(block *types.Block) *big.Int {
	if len(block.Transactions()) == 0 {
		return big.NewInt(0)
	}

	total := new(big.Int)
	for _, tx := range block.Transactions() {
		total.Add(total, tx.GasPrice())
	}

	avg := new(big.Int).Div(total, big.NewInt(int64(len(block.Transactions()))))
	return avg
}

// printStats 打印统计列表
func printStats(stats []BlockStats) {
	fmt.Println("\n=== 最近区块统计 ===")
	for i, s := range stats {
		fmt.Printf("%d. 区块#%d: %d tx, Gas:%d, Avg:%s Gwei\n",
			i+1, s.Number, s.TxCount, s.TotalGasUsed, s.AvgGasPrice.String())
	}
}
