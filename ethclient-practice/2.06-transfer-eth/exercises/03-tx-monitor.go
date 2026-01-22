// 03-tx-monitor.go - 交易监控器练习
//
// 任务：
// 1. 发送一笔交易
// 2. 实时监控交易状态
// 3. 显示交易是否在 mempool、是否被打包、是否成功
// 4. 显示实际 Gas 消耗
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TO_ADDRESS=0x... && go run exercises/03-tx-monitor.go

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 交易监控器 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}

	toAddressHex := os.Getenv("TO_ADDRESS")
	if toAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 TO_ADDRESS")
	}

	// TODO 1: 连接并加载私钥
	var client *ethclient.Client
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	{
		// 在这里填写代码
	}
	defer client.Close()

	// TODO 2: 获取 Nonce 和 Gas Price
	var nonce uint64
	var gasPrice *big.Int
	{
		// 在这里填写代码
	}

	// TODO 3: 构建并发送交易
	value := big.NewInt(1000000000000000) // 0.001 ETH
	toAddress := common.HexToAddress(toAddressHex)
	var signedTx *types.Transaction
	{
		// 在这里填写代码
	}

	txHash := signedTx.Hash()
	fmt.Printf("\n交易已发送: %s\n", txHash.Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n\n", txHash.Hex())

	// TODO 4: 开始监控
	fmt.Println("开始监控交易状态...")
	fmt.Println("────────────────────────────────────────")

	// 状态标志
	inMempool := false
	isPending := true
	isConfirmed := false
	var receipt *types.Receipt

	// TODO 5: 轮询监控
	for {
		// 检查交易是否在 mempool
		isInMempool, _ := isTransactionInMempool(client, txHash)
		if isInMempool && !inMempool {
			inMempool = true
			fmt.Printf("[%s] ✅ 交易在 Mempool 中\n", formatTime())
		}

		// 检查交易是否已确认
		if isPending {
			receipt, _ = client.TransactionReceipt(context.Background(), txHash)
			if receipt != nil {
				isPending = false
				isConfirmed = true

				fmt.Printf("\n[%s] 🎉 交易已打包！\n", formatTime())
				fmt.Printf("  区块号: %d\n", receipt.BlockNumber.Uint64())
				fmt.Printf("  区块哈希: %s\n", receipt.BlockHash.Hex())
				fmt.Printf("  交易索引: %d\n", receipt.TransactionIndex)

				// 检查状态
				if receipt.Status == 1 {
					fmt.Printf("\n[%s] ✅ 交易成功！\n", formatTime())
					fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
					fmt.Printf("  Gas Limit: %d\n", gasLimit)

					// 计算实际费用
					// TODO: 填写计算代码

				} else {
					fmt.Printf("\n[%s] ❌ 交易失败\n", formatTime())
				}

				break
			}
		}

		// 等待一段时间再检查
		time.Sleep(5 * time.Second)
		fmt.Printf("[%s] ⏳ 等待确认...\n", formatTime())
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Println("=== 监控结束 ===")
}

// 辅助函数：检查交易是否在 mempool 中
func isTransactionInMempool(client *ethclient.Client, txHash common.Hash) (bool, error) {
	// TODO: 实现检查逻辑
	// 提示：可以使用 client.TransactionInPool() 或查询交易来检查
	// 如果交易还未被打包，TransactionReceipt 会返回错误
	return false, nil
}

// 辅助函数：格式化时间
func formatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 辅助函数：Wei 转 ETH
func weiToEth(wei *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18),
	)
}
