// 查询区块：连接 Sepolia，按区块号查询并输出区块信息。
//
// 环境变量：SEPOLIA_RPC_URL（必填）
// 可选：BLOCK_NUMBER  不设则查最新区块
//
// 运行：在 task1 目录下执行 go run ./blockchain-read-write/query-block
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func getRPCURL() string {
	u := strings.TrimSpace(os.Getenv("SEPOLIA_RPC_URL"))
	if u == "" {
		log.Fatal("请设置环境变量 SEPOLIA_RPC_URL")
	}
	return u
}

func main() {
	client, err := ethclient.Dial(getRPCURL())
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	ctx := context.Background()
	blockNumberStr := strings.TrimSpace(os.Getenv("BLOCK_NUMBER"))
	if blockNumberStr == "" {
		log.Fatal("请设置环境变量 BLOCK_NUMBER")
	}
	blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
	if err != nil {
		log.Fatal("无效的 BLOCK_NUMBER:", err)
	}
	block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Fatal("查询区块失败:", err)
	}

	blockTime := time.Unix(int64(block.Time()), 0).UTC()
	timeStr := blockTime.Format(time.DateTime + " MST") // 用标准常量，不必记 2006-01-02；MST 会按实际时区显示（如 UTC）

	fmt.Println("=== 区块信息 (Sepolia) ===")
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("时间戳: %s\n", timeStr)
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("Gas 使用: %d\n", block.GasUsed())
	fmt.Printf("Gas 上限: %d\n", block.GasLimit())
	fmt.Println("=== 完成 ===")
}
