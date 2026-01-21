package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接 Sepolia 测试网络
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 获取区块
	blockNumber := big.NewInt(10077132)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 获取链 ID，用于恢复发送者地址
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 使用 CancunSigner 支持所有交易类型（包括 Blob 交易）
	signer := types.NewCancunSigner(chainID)

	// 获取所有交易
	transactions := block.Transactions()
	txCount := len(transactions)

	// 先打印交易数量
	fmt.Printf("=== 区块 %d 包含 %d 笔交易 ===\n\n", blockNumber.Uint64(), txCount)

	// 遍历并显示所有交易
	fmt.Println("=== 交易列表 ===")
	for i, tx := range transactions {
		sender, err := types.Sender(signer, tx)
		if err != nil {
			log.Fatal(err)
		}

		// 转换金额为 Ether
		valueInEther := new(big.Float).Quo(
			new(big.Float).SetInt(tx.Value()),
			big.NewFloat(1e18),
		)

		// 格式化 Gas Price
		gasPriceStr := formatGasPrice(tx)

		// 显示交易信息
		fmt.Printf("[%d] Hash: %s\n", i+1, tx.Hash().Hex())
		fmt.Printf("    From: %s\n", sender.Hex())
		to := "<合约创建>"
		if tx.To() != nil {
			to = tx.To().Hex()
		}
		fmt.Printf("    To: %s\n", to)
		fmt.Printf("    Value: %.6f Ether\n", valueInEther)
		fmt.Printf("    Gas Price: %s\n", gasPriceStr)
		fmt.Println()
	}
}

// 格式化 Gas Price 显示
func formatGasPrice(tx *types.Transaction) string {
	weiToGwei := big.NewFloat(1e9)

	if tx.Type() >= 2 {
		// EIP-1559 或 Blob 交易
		maxFee := new(big.Float).Quo(new(big.Float).SetInt(tx.GasFeeCap()), weiToGwei)
		priority := new(big.Float).Quo(new(big.Float).SetInt(tx.GasTipCap()), weiToGwei)
		result := fmt.Sprintf("MaxFee: %.2f, Priority: %.2f Gwei", maxFee, priority)
		if tx.Type() == 3 {
			result += " (Blob Tx)"
		}
		return result
	}

	// Legacy 或 EIP-2930 交易
	gasPrice := new(big.Float).Quo(new(big.Float).SetInt(tx.GasPrice()), weiToGwei)
	return fmt.Sprintf("%.2f Gwei", gasPrice)
}
