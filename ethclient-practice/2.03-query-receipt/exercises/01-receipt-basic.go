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
	txHash := common.HexToHash("0xb7cedb112cb9b246faed99ffdfd8bfcdca9215e32dc2b4d006d7afc529a0c625")

	// 练习：查询交易收据
	// 提示：使用 client.TransactionReceipt 方法
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	// var receipt *types.Receipt
	// receipt, err = ???

	// 练习：显示收据信息
	// 提示：receipt.Status, receipt.GasUsed, receipt.TransactionIndex
	fmt.Printf("=== 交易收据信息 ===")
	fmt.Printf("交易哈希: %s\n", receipt.TxHash.Hex())
	fmt.Printf("区块号：%d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("区块哈希：%s\n", receipt.BlockHash.Hex())

	fmt.Printf("状态: %d\n", receipt.Status)
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("交易索引: %d\n", receipt.TransactionIndex)
	fmt.Printf("合约地址: %s\n", receipt.ContractAddress.Hex())

	// 练习：判断是否为合约创建交易
	// 提示：receipt.ContractAddress 是否为零地址
	if receipt.ContractAddress != (common.Address{}) {
		fmt.Printf("创建的合约地址: %s\n", receipt.ContractAddress.Hex())
	} else {
		fmt.Println("非合约创建交易")
	}

}
