package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
	// Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入区块号 (输入 'latest' 获取最新): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// 练习：解析用户输入，获取区块号
	// 提示：如果输入是 "latest"，获取最新区块号；否则解析为数字
	// var blockNumber *big.Int
	// if input == "latest" {
	//     // header, err := client.HeaderByNumber(???, ???)
	//     // blockNumber = header.Number
	// } else {
	//     // num, err := strconv.ParseUint(???, 10, 64)
	//     // blockNumber = big.NewInt(int64(num))
	// }

	// 练习：获取区块数据
	// block, err := client.BlockByNumber(???, ???)

	// 练习：调用 printBlockInfo 函数显示区块信息
	// printBlockInfo(block)
}

func printBlockInfo(block *types.Block) {
	// 练习：格式化时间戳
	// blockTime := time.Unix(???, ???)
	// timeFormatted := blockTime.UTC().Format("2006-01-02 15:04:05 UTC")

	fmt.Printf("\n=== 区块 %d ===\n", block.Number().Uint64())

	fmt.Println("\n基本信息:")
	fmt.Printf("  - 时间戳: %s\n", timeFormatted)
	fmt.Printf("  - 矿工: %s\n", block.Coinbase().Hex())
	// 练习：计算 Gas 使用率百分比
	// fmt.Printf("  - Gas 使用: %d / %d (%.1f%%)\n", ..., ..., ???)

	transactions := block.Transactions()
	fmt.Printf("\n交易列表 (共 %d 笔):\n", len(transactions))

	// 练习：遍历所有交易，计算总 Gas 和总 Gas 价格
	// totalGas := uint64(0)
	// totalGasPrice := big.NewInt(0)
	//
	// for i, tx := range transactions {
	//     gasPrice := tx.GasPrice()
	//     // TODO: 累加 totalGas
	//     // TODO: 累加 totalGasPrice
	//
	//     fmt.Printf("  [%2d] %s - Gas: %d - Price: %s Gwei\n",
	//         i+1,
	//         shortenHash(tx.Hash().Hex()),
	//         tx.Gas(),
	//         formatGasPrice(gasPrice))
	// }

	// 练习：计算平均 Gas 价格
	// avgGasPrice := big.NewInt(0)
	// if len(transactions) > 0 {
	//     // avgGasPrice.Div(???, ???)
	// }

	fmt.Println("\n统计:")
	fmt.Printf("  - 总交易: %d\n", len(transactions))
	fmt.Printf("  - 总 Gas 使用: %s\n", formatNumber(totalGas))
	fmt.Printf("  - 平均 Gas 价格: %s Gwei\n", formatGasPrice(avgGasPrice))
}

// 辅助函数：缩短哈希显示
func shortenHash(hash string) string {
	if len(hash) < 16 {
		return hash
	}
	return hash[:10] + "..." + hash[len(hash)-4:]
}

// 辅助函数：格式化 Gas 价格（Wei 转 Gwei）
func formatGasPrice(price *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(price), big.NewFloat(1e9))
	return fmt.Sprintf("%.2f", gwei)
}

// 辅助函数：格式化数字（添加千位分隔符）
func formatNumber(n uint64) string {
	str := strconv.FormatUint(n, 10)
	var result []byte
	for i := len(str); i > 0; i-- {
		if i < len(str) && (len(str)-i)%3 == 0 {
			result = append([]byte{','}, result...)
		}
		result = append([]byte{str[i-1]}, result...)
	}
	return string(result)
}
