package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"

	"github.com/ethereum/go-ethereum/common"
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

	// 练习：获取区块
	// block, err := ???

	// 练习：获取链 ID
	// chainID, err := ???
	// signer := types.NewEIP155Signer(chainID)

	transactions := block.Transactions()

	// 练习：定义统计变量
	// totalTxs := len(transactions)
	// totalValue := big.NewInt(0)
	// totalGasUsed := uint64(0)
	// totalGasPrice := big.NewInt(0)
	// successCount := 0
	// failCount := 0
	// contractCreations := 0
	//
	// txFees := make([]*TxFee, 0)

	// 练习：遍历所有交易进行统计
	// for _, tx := range transactions {
	//     // TODO: 累加总金额
	//     totalValue.Add(totalValue, tx.Value())
	//
	//     // TODO: 累加 Gas Price
	//     totalGasPrice.Add(totalGasPrice, tx.GasPrice())
	//
	//     // TODO: 计算交易费用 (Gas × GasPrice)
	//     fee := new(big.Int).Mul(big.NewInt(int64(tx.Gas())), tx.GasPrice())
	//
	//     // TODO: 查询收据获取状态和 Gas Used
	//     receipt, err := client.TransactionReceipt(???, ???)
	//     if err == nil {
	//         totalGasUsed += receipt.GasUsed
	//         if receipt.Status == 1 {
	//             successCount++
	//         } else {
	//             failCount++
	//         }
	//     }
	//
	//     // TODO: 检查是否为合约创建（tx.To() == nil）
	//     if ??? {
	//         contractCreations++
	//     }
	//
	//     // TODO: 记录交易费用用于排序
	//     txFees = append(txFees, &TxFee{
	//         Hash: tx.Hash().Hex(),
	//         Fee:  fee,
	//     })
	// }

	// 练习：排序找出最贵的交易
	// sort.Slice(txFees, func(i, j int) bool {
	//     return txFees[i].Fee.Cmp(txFees[j].Fee) > 0
	// })

	// 练习：输出统计结果
	fmt.Println("=== 区块交易统计 ===")
	// fmt.Printf("区块号: %d\n", block.Number().Uint64())
	// fmt.Printf("总交易数: %d\n", totalTxs)
	// fmt.Printf("\n金额统计:\n")
	// fmt.Printf("  总转账: %s Ether\n", ???)
	// fmt.Printf("\nGas 统计:\n")
	// fmt.Printf("  总 Gas 使用: %d\n", totalGasUsed)
	// fmt.Printf("  平均 Gas 价格: %s Gwei\n", ???)
	// fmt.Printf("\n交易状态:\n")
	// fmt.Printf("  成功: %d\n", successCount)
	// fmt.Printf("  失败: %d\n", failCount)
	// fmt.Printf("  合约创建: %d\n", contractCreations)
	// fmt.Printf("\n最贵的 3 笔交易:\n")
	// for i, tx := range txFees {
	//     if i >= 3 {
	//         break
	//     }
	//     fmt.Printf("  %d. %s - %s Wei\n", i+1, tx.Hash[:10]+"...", tx.Fee.String())
	// }
}

// TxFee 用于记录交易费用
type TxFee struct {
	Hash string
	Fee  *big.Int
}
