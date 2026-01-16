package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// 获取区块
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signer := types.NewEIP155Signer(chainID)

	transactions := block.Transactions()

	// 定义统计变量
	totalTxs := len(transactions)
	totalValue := big.NewInt(0)
	totalGasUsed := uint64(0)
	totalGasPrice := big.NewInt(0)
	successCount := 0
	failCount := 0
	contractCreations := 0

	txFees := make([]*TxFee, 0)

	// 遍历所有交易进行统计
	for _, tx := range transactions {
		// 累加总金额
		totalValue.Add(totalValue, tx.Value())

		// 累加 Gas Price
		totalGasPrice.Add(totalGasPrice, tx.GasPrice())

		// 查询收据获取状态和 Gas Used
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			totalGasUsed += receipt.GasUsed
			if receipt.Status == 1 {
				successCount++
			} else {
				failCount++
			}

			// 计算实际交易费用
			actualFee := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GasPrice())
			txFees = append(txFees, &TxFee{
				Hash: tx.Hash().Hex(),
				Fee:  actualFee,
			})
		}

		// 检查是否为合约创建（tx.To() == nil）
		if tx.To() == nil {
			contractCreations++
		}
	}

	// 排序找出最贵的交易
	sort.Slice(txFees, func(i, j int) bool {
		return txFees[i].Fee.Cmp(txFees[j].Fee) > 0
	})

	// 计算平均值
	var avgGasPrice *big.Float
	if totalTxs > 0 {
		avgGasPrice = new(big.Float).Quo(
			new(big.Float).SetInt(totalGasPrice),
			big.NewFloat(int64(totalTxs)),
		)
	}

	// 转换总金额为 Ether
	totalValueEther := new(big.Float).Quo(
		new(big.Float).SetInt(totalValue),
		big.NewFloat(1e18),
	)

	// 输出统计结果
	fmt.Println("=== 区块交易统计 ===")
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("总交易数: %d\n\n", totalTxs)

	fmt.Println("金额统计:")
	fmt.Printf("  总转账: %.6f Ether\n", totalValueEther)

	fmt.Println("\nGas 统计:")
	fmt.Printf("  总 Gas 使用: %d\n", totalGasUsed)
	fmt.Printf("  平均 Gas 价格: %.2f Gwei\n", new(big.Float).Quo(avgGasPrice, big.NewFloat(1e9)))

	fmt.Println("\n交易状态:")
	fmt.Printf("  成功: %d\n", successCount)
	fmt.Printf("  失败: %d\n", failCount)
	fmt.Printf("  合约创建: %d\n", contractCreations)

	fmt.Println("\n最贵的 3 笔交易:")
	for i, tx := range txFees {
		if i >= 3 {
			break
		}
		feeEther := new(big.Float).Quo(
			new(big.Float).SetInt(tx.Fee),
			big.NewFloat(1e18),
		)
		fmt.Printf("  %d. %s - %.6f Ether\n", i+1, shortenHash(tx.Hash), feeEther)
	}
}

func shortenHash(hash string) string {
	if len(hash) < 16 {
		return hash
	}
	return hash[:10] + "..." + hash[len(hash)-4:]
}

// TxFee 用于记录交易费用
type TxFee struct {
	Hash string
	Fee  *big.Int
}
