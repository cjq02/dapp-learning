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

	// 一个包含事件日志的交易哈希（USDT 转账）
	txHash := common.HexToHash("0x6c0a9e58ec065b56c1ddf4cfeae63c5719bed2ae4b46c3c64fe1ca9e8e8987a6")

	// 练习：查询交易收据
	// var receipt *types.Receipt
	// receipt, err = ???

	fmt.Println("=== 事件日志分析 ===")

	// 练习：遍历 receipt.Logs 数组
	// for i, log := range ??? {
	//     fmt.Printf("\n[日志 %d]\n", i+1)
	//     fmt.Printf("  合约地址: %s\n", log.Address.Hex())
	//     fmt.Printf("  主题数量: %d\n", len(log.Topics))
	//     fmt.Printf("  数据长度: %d bytes\n", len(log.Data))
	//
	//     // 练习：显示第一个主题（通常是事件签名）
	//     if len(log.Topics) > 0 {
	//         fmt.Printf("  事件签名: %s\n", log.Topics[0].Hex())
	//     }
	// }

	// 练习：统计日志总数
	// fmt.Printf("\n总日志数量: %d\n", ???)
}
