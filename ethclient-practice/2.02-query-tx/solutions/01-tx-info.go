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

	// 获取区块和所有交易
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 获取链 ID，用于恢复发送者地址
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 创建 EIP155 签名器
	signer := types.NewEIP155Signer(chainID)

	// 遍历所有交易并显示信息
	fmt.Println("=== 交易列表 ===")
	for i, tx := range block.Transactions() {
		// 恢复发送者地址
		sender, err := types.Sender(signer, tx)
		if err != nil {
			log.Fatal(err)
		}

		// 转换金额为 Ether (1 Ether = 10^18 Wei)
		valueInEther := new(big.Float).Quo(
			new(big.Float).SetInt(tx.Value()),
			big.NewFloat(1e18),
		)

		// 转换 Gas Price 为 Gwei (1 Gwei = 10^9 Wei)
		gasPriceInGwei := new(big.Float).Quo(
			new(big.Float).SetInt(tx.GasPrice()),
			big.NewFloat(1e9),
		)

		fmt.Printf("[%d] Hash: %s\n", i+1, tx.Hash().Hex())
		fmt.Printf("    From: %s\n", sender.Hex())
		fmt.Printf("    To: %s\n", tx.To().Hex())
		fmt.Printf("    Value: %.6f Ether\n", valueInEther)
		fmt.Printf("    Gas Price: %.2f Gwei\n", gasPriceInGwei)
		fmt.Println()
	}
}
