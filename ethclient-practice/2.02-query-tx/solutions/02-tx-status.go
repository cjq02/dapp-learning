package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

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

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入交易哈希: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	txHash := common.HexToHash(input)

	// 查询交易
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	// 显示交易基本信息
	fmt.Println("\n=== 交易信息 ===")
	fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("Value: %s Wei\n", tx.Value().String())
	fmt.Printf("Gas: %d\n", tx.Gas())

	// 转换 Gas Price 为 Gwei
	gasPriceInGwei := new(big.Float).Quo(
		new(big.Float).SetInt(tx.GasPrice()),
		big.NewFloat(1e9),
	)
	fmt.Printf("Gas Price: %.2f Gwei\n", gasPriceInGwei)

	fmt.Printf("Nonce: %d\n", tx.Nonce())
	if tx.To() != nil {
		fmt.Printf("To: %s\n", tx.To().Hex())
	} else {
		fmt.Printf("To: <合约创建>\n")
	}
	fmt.Printf("Pending: %v\n", isPending)

	// 获取链 ID 并恢复发送者地址
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signer := types.NewEIP155Signer(chainID)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("From: %s\n", sender.Hex())

	// 查询交易收据
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 显示交易状态
	fmt.Println("\n=== 执行状态 ===")
	status := "失败"
	if receipt.Status == 1 {
		status = "成功"
	}
	fmt.Printf("状态: %s\n", status)
	fmt.Printf("Gas Used: %d / %d (%.1f%%)\n",
		receipt.GasUsed,
		tx.Gas(),
		float64(receipt.GasUsed)/float64(tx.Gas())*100)
	fmt.Printf("日志数量: %d\n", len(receipt.Logs))

	// 显示交易费用
	txFee := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), tx.GasPrice())
	txFeeEther := new(big.Float).Quo(
		new(big.Float).SetInt(txFee),
		big.NewFloat(1e18),
	)
	fmt.Printf("交易费用: %.6f Ether\n", txFeeEther)
}
