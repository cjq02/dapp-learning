// 03-balance-monitor.go - 余额监控器 - 答案

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 余额监控器 ===")
	fmt.Println("按 Ctrl+C 退出\n")

	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 要监控的地址
	address := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	// 初始化上一次余额
	var lastBalance *big.Int

	// 设置定时器，每 10 秒查询一次
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// 设置信号监听，捕获 Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// 首次查询
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Printf("查询余额失败: %v\n", err)
	} else {
		lastBalance = balance
		fmt.Printf("[%s] 初始余额: %.6f ETH\n", formatTime(), weiToEth(balance))
	}

	// 主循环
	for {
		select {
		case <-ticker.C:
			// 定时器触发，查询余额
			balance, err := client.BalanceAt(context.Background(), address, nil)
			if err != nil {
				log.Printf("查询余额失败: %v\n", err)
				continue
			}

			// 检查余额是否变化
			if lastBalance == nil || balance.Cmp(lastBalance) != 0 {
				change := new(big.Int).Sub(balance, lastBalance)
				changeEth := weiToEth(change)

				fmt.Printf("\n[%s] 余额变化！\n", formatTime())
				fmt.Printf("  旧余额: %.6f ETH\n", weiToEth(lastBalance))
				fmt.Printf("  新余额: %.6f ETH\n", weiToEth(balance))

				if change.Sign() > 0 {
					fmt.Printf("  变化: +%.6f ETH (收入)\n\n", changeEth)
				} else {
					fmt.Printf("  变化: %.6f ETH (支出)\n\n", changeEth)
				}

				lastBalance = balance
			} else {
				fmt.Printf("[%s] 余额未变化: %.6f ETH\n", formatTime(), weiToEth(balance))
			}

		case <-interrupt:
			// 捕获退出信号
			fmt.Println("\n\n收到退出信号，停止监控...")
			fmt.Printf("最终余额: %.6f ETH\n", weiToEth(lastBalance))
			fmt.Println("=== 监控结束 ===")
			return
		}
	}
}

// 辅助函数：格式化时间
func formatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 辅助函数：Wei 转 ETH
func weiToEth(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	return new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18),
	)
}
