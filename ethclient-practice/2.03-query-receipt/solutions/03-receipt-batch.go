package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	// 从环境变量读取 API Key
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY\n例如: export INFURA_API_KEY=your-key-here")
	}
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	blockNumber := big.NewInt(5671744)

	// 查询指定区块的所有收据
	receipts, err := client.BlockReceipts(
		context.Background(),
		rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 区块收据统计 ===")

	// 定义统计变量
	totalCount := len(receipts)
	successCount := 0
	failCount := 0
	totalGasUsed := uint64(0)
	contractCreations := 0

	// 定义 Gas 使用记录（用于找出最大值）
	type GasRecord struct {
		TxHash  string
		GasUsed uint64
	}
	gasRecords := make([]*GasRecord, 0)

	// 遍历所有收据进行统计
	for _, receipt := range receipts {
		totalGasUsed += receipt.GasUsed

		// 统计成功/失败
		if receipt.Status == 1 {
			successCount++
		} else {
			failCount++
		}

		// 统计合约创建
		if receipt.ContractAddress != (common.Address{}) {
			contractCreations++
		}

		// 记录 Gas 使用
		gasRecords = append(gasRecords, &GasRecord{
			TxHash:  receipt.TxHash.Hex(),
			GasUsed: receipt.GasUsed,
		})
	}

	// 排序找出 Gas 使用最多的交易
	sort.Slice(gasRecords, func(i, j int) bool {
		return gasRecords[i].GasUsed > gasRecords[j].GasUsed
	})

	// 输出统计结果
	fmt.Printf("区块号: %d\n", blockNumber.Uint64())
	fmt.Printf("总收据数: %d\n\n", totalCount)

	fmt.Println("交易状态:")
	fmt.Printf("  成功: %d\n", successCount)
	fmt.Printf("  失败: %d\n\n", failCount)

	fmt.Println("Gas 统计:")
	fmt.Printf("  总 Gas 使用: %d\n", totalGasUsed)
	if totalCount > 0 {
		avgGas := totalGasUsed / uint64(totalCount)
		fmt.Printf("  平均 Gas 使用: %d\n", avgGas)
	}
	fmt.Printf("  合约创建: %d\n\n", contractCreations)

	fmt.Println("Gas 使用最多的 3 笔交易:")
	for i, record := range gasRecords {
		if i >= 3 {
			break
		}
		fmt.Printf("  %d. %s - %d Gas\n", i+1, shortenHash(record.TxHash), record.GasUsed)
	}
}

func shortenHash(hash string) string {
	if len(hash) < 16 {
		return hash
	}
	return hash[:10] + "..." + hash[len(hash)-4:]
}
