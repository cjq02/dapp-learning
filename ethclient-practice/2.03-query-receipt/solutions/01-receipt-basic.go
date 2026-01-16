package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
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

	// 交易哈希示例
	txHash := common.HexToHash("0x20294a03e8766e9aeab58327fc4112756017c6c28f6f99c7722f4a29075601c5")

	// 查询交易收据
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 交易收据信息 ===")
	fmt.Printf("交易哈希: %s\n", receipt.TxHash.Hex())
	fmt.Printf("区块号: %d\n", receipt.BlockNumber)
	fmt.Printf("区块哈希: %s\n", receipt.BlockHash.Hex())

	// 显示收据信息
	status := "失败"
	if receipt.Status == 1 {
		status = "成功"
	}
	fmt.Printf("状态: %s\n", status)
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("累计 Gas: %d\n", receipt.CumulativeGasUsed)
	fmt.Printf("交易索引: %d\n", receipt.TransactionIndex)

	// 判断是否为合约创建交易
	if receipt.ContractAddress != (common.Address{}) {
		fmt.Printf("创建的合约地址: %s\n", receipt.ContractAddress.Hex())
	} else {
		fmt.Println("非合约创建交易")
	}

	fmt.Printf("From: %s\n", receipt.From.Hex())
	if receipt.To != nil {
		fmt.Printf("To: %s\n", receipt.To.Hex())
	}
}
