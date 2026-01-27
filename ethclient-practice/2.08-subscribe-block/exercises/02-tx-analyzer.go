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
	Number      uint64
	Hash        string
	TxCount     int
	TotalGasUsed uint64
	AvgGasPrice *big.Int
}

func main() {
	wsURL := os.Getenv("SEPOLIA_WS_URL")
	if wsURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_WS_URL")
	}

	// 练习：连接到以太坊节点
	// var client *ethclient.Client
	// client, err = ???

	// 练习：创建通道并订阅
	// var headers chan *types.Header
	// var sub ethereum.Subscription

	// 用于存储最近的区块统计（最多10个）
	var stats []BlockStats
	var mu sync.Mutex

	fmt.Println("开始监听并分析新区块...")
	fmt.Println("================================")

	// 练习：监听新区块并分析
	// for {
	//     select {
	//     case err := <-sub.Err():
	//         log.Fatal(err)
	//     case header := <-headers:
	//         // TODO: 获取完整区块
	//         // block, err := ???
	//
	//         // TODO: 计算统计数据
	//         // 1. 交易数量: len(block.Transactions())
	//         // 2. 总 Gas 使用: block.GasUsed()
	//         // 3. 平均 Gas 价格: 遍历所有交易，计算平均值
	//
	//         // TODO: 更新 stats 列表（保持最多10个）
	//         // 使用 mu.Lock() 和 mu.Unlock() 保护并发访问
	//
	//         // TODO: 打印当前区块的统计信息
	//     }
	// }
}

// 辅助函数：计算平均 Gas 价格
func calculateAverageGasPrice(block *types.Block) *big.Int {
	// TODO: 遍历 block.Transactions()，计算平均 Gas Price
	// 提示：
	// total := new(big.Int)
	// for _, tx := range block.Transactions() {
	//     total.Add(total, tx.GasPrice())
	// }
	// return new(big.Int).Div(total, big.NewInt(int64(len(block.Transactions()))))
	return nil
}

// 辅助函数：打印统计列表
func printStats(stats []BlockStats) {
	// TODO: 打印格式化的统计信息
	fmt.Println("\n=== 最近区块统计 ===")
	for _, s := range stats {
		fmt.Printf("区块 #%d: %d 笔交易, Gas 使用: %d, 平均 Gas 价格: %s\n",
			s.Number, s.TxCount, s.TotalGasUsed, s.AvgGasPrice.String())
	}
}
