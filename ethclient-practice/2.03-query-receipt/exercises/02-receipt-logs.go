package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"

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

	// 一个包含事件日志的交易哈希（USDT 转账）
	txHash := common.HexToHash("0xb7cedb112cb9b246faed99ffdfd8bfcdca9215e32dc2b4d006d7afc529a0c625")

	// 练习：查询交易收据
	var receipt *types.Receipt
	receipt, err = client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 事件日志分析 ===")

	// 练习：遍历 receipt.Logs 数组
	for i, log := range receipt.Logs {
		fmt.Printf("\n[日志 %d]\n", i+1)
		fmt.Printf("  合约地址: %s\n", log.Address.Hex())
		fmt.Printf("  区块号: %d\n", log.BlockNumber)
		fmt.Printf("  交易索引: %d\n", log.TxIndex)
		fmt.Printf("  日志索引: %d\n", log.Index)

		// 练习：显示第一个主题（通常是事件签名）
		if len(log.Topics) > 0 {
			fmt.Printf("  事件签名: %s\n", log.Topics[0].Hex())
		}

		fmt.Printf("  主题数量: %d\n", len(log.Topics))
		for j, topic := range log.Topics {
			fmt.Printf("    Topic[%d]: %s\n", j, topic.Hex())
		}

		dataHex := hex.EncodeToString(log.Data)
		if len(dataHex) > 0 {
			fmt.Printf("  数据: %s\n", dataHex)
			fmt.Printf("  数据长度: %d bytes\n", len(log.Data))
		}

		// 练习：统计日志总数
		fmt.Printf("\n总日志数量: %d\n", len(receipt.Logs))
	}
}
