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
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入区块号 (输入 'latest' 获取最新): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var blockNumber *big.Int
	if input == "latest" {
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		blockNumber = header.Number
	} else {
		num, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			log.Fatal("无效的区块号")
		}
		blockNumber = big.NewInt(int64(num))
	}

	// 获取区块
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	printBlockInfo(block)
}

func printBlockInfo(block *types.Block) {
	blockTime := time.Unix(int64(block.Time()), 0)
	timeFormatted := blockTime.UTC().Format("2006-01-02 15:04:05 UTC")

	fmt.Printf("\n=== 区块 %d ===\n", block.Number().Uint64())

	fmt.Println("\n基本信息:")
	fmt.Printf("  - 时间戳: %s\n", timeFormatted)
	fmt.Printf("  - 矿工: %s\n", block.Coinbase().Hex())
	fmt.Printf("  - Gas 使用: %d / %d (%.1f%%)\n",
		block.GasUsed(),
		block.GasLimit(),
		float64(block.GasUsed())/float64(block.GasLimit())*100)

	transactions := block.Transactions()
	fmt.Printf("\n交易列表 (共 %d 笔):\n", len(transactions))

	totalGas := uint64(0)
	totalGasPrice := big.NewInt(0)

	for i, tx := range transactions {
		gasPrice := tx.GasPrice()
		totalGas += tx.Gas()
		totalGasPrice.Add(totalGasPrice, gasPrice)

		fmt.Printf("  [%2d] %s - Gas: %d - Price: %s Gwei\n",
			i+1,
			shortenHash(tx.Hash().Hex()),
			tx.Gas(),
			formatGasPrice(gasPrice))
	}

	avgGasPrice := big.NewInt(0)
	if len(transactions) > 0 {
		avgGasPrice.Div(totalGasPrice, big.NewInt(int64(len(transactions))))
	}

	fmt.Println("\n统计:")
	fmt.Printf("  - 总交易: %d\n", len(transactions))
	fmt.Printf("  - 总 Gas 使用: %s\n", formatNumber(totalGas))
	fmt.Printf("  - 平均 Gas 价格: %s Gwei\n", formatGasPrice(avgGasPrice))
}

func shortenHash(hash string) string {
	if len(hash) < 16 {
		return hash
	}
	return hash[:10] + "..." + hash[len(hash)-4:]
}

func formatGasPrice(price *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(price), big.NewFloat(1e9))
	return fmt.Sprintf("%.2f", gwei)
}

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
