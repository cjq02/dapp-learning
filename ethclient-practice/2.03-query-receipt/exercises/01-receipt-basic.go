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

	// 练习：查询交易收据
	// 提示：使用 client.TransactionReceipt 方法
	// var receipt *types.Receipt
	// receipt, err = ???

	// 练习：显示收据信息
	// 提示：receipt.Status, receipt.GasUsed, receipt.TransactionIndex
	// fmt.Printf("状态: %d\n", ???)
	// fmt.Printf("Gas 使用: %d\n", ???)
	// fmt.Printf("交易索引: %d\n", ???)

	// 练习：判断是否为合约创建交易
	// 提示：receipt.ContractAddress 是否为零地址
	// if receipt.ContractAddress != (common.Address{}) {
	//     fmt.Printf("创建的合约地址: %s\n", ???)
	// } else {
	//     fmt.Println("非合约创建交易")
	// }
}
