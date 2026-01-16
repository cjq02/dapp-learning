package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入交易哈希: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	txHash := common.HexToHash(input)

	// 练习：查询交易
	// tx, isPending, err := client.TransactionByHash(???, ???)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 练习：显示交易基本信息
	// fmt.Printf("=== 交易信息 ===\n")
	// fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	// fmt.Printf("Value: %s Wei\n", tx.Value().String())
	// fmt.Printf("Gas: %d\n", tx.Gas())
	// fmt.Printf("Gas Price: %s Gwei\n", ???)
	// fmt.Printf("Nonce: %d\n", tx.Nonce())
	// fmt.Printf("To: %s\n", tx.To().Hex())
	// fmt.Printf("Pending: %v\n", isPending)

	// 练习：获取链 ID 并恢复发送者地址
	// chainID, err := ???
	// signer := types.NewEIP155Signer(chainID)
	// sender, err := types.Sender(signer, tx)

	// 练习：查询交易收据
	// receipt, err := client.TransactionReceipt(???, ???)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 练习：显示交易状态
	// fmt.Printf("\n=== 执行状态 ===\n")
	// status := "失败"
	// if receipt.Status == 1 {
	//     status = "成功"
	// }
	// fmt.Printf("状态: %s\n", status)
	// fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
	// fmt.Printf("日志数量: %d\n", len(receipt.Logs))
}
