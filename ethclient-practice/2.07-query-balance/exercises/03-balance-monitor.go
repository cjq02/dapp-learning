package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	account := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// TODO: 获取初始余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	lastBalance := balance

	// 转换为 ETH 的辅助函数
	weiToEth := func(wei *big.Int) *big.Float {
		return new(big.Float).Quo(
			new(big.Float).SetInt(wei),
			big.NewFloat(math.Pow10(18)),
		)
	}

	fmt.Printf("开始监控地址: %s\n", account.Hex())
	fmt.Printf("初始余额: %.18f ETH\n", weiToEth(lastBalance))
	fmt.Println("按 Ctrl+C 退出...")

	// 设置信号捕获
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// 定时器
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-sigCh:
				fmt.Println("\n收到退出信号，停止监控...")
				return
			case <-ticker.C:
				// TODO: 查询当前余额
				currentBalance, err := client.BalanceAt(context.Background(), account, nil)
				if err != nil {
					log.Printf("查询余额失败: %v", err)
					continue
				}

				// TODO: 检查余额是否变化
				if currentBalance.Cmp(lastBalance) != 0 {
					timestamp := time.Now().Format("2006-01-02 15:04:05")
					oldEth := weiToEth(lastBalance)
					newEth := weiToEth(currentBalance)
					change := new(big.Float).Sub(newEth, oldEth)

					fmt.Printf("\n[%s] 余额变化！\n", timestamp)
					fmt.Printf("旧余额: %.18f ETH\n", oldEth)
					fmt.Printf("新余额: %.18f ETH\n", newEth)
					fmt.Printf("变化: %.18f ETH\n\n", change)

					lastBalance = currentBalance
				} else {
					fmt.Printf("[%s] 余额未变化\n", time.Now().Format("2006-01-02 15:04:05"))
				}
			}
		}
	}()

	wg.Wait()
	fmt.Println("监控已停止")
}
