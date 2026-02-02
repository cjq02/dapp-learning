// 03-balance-monitor.go - 余额监控器练习
//
// 任务：
// 1. 每隔 10 秒查询一次地址余额
// 2. 检测余额变化
// 3. 当余额变化时打印通知
// 4. 按 Ctrl+C 退出程序
//
// 运行：export INFURA_API_KEY=your-key && go run exercises/03-balance-monitor.go

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

	// TODO 1: 初始化上一次余额
	var lastBalance *big.Int

	// TODO 2: 设置定时器，每 10 秒查询一次
	// 提示：使用 time.NewTicker(10 * time.Second)

	// TODO 3: 设置信号监听，捕获 Ctrl+C
	// 提示：使用 signal.Notify() 和 os.Interrupt

	// TODO 4: 主循环
	// 提示：使用 select 监听定时器和退出信号
	// 提示：每次定时器触发时查询余额
	// 提示：比较余额是否变化
	// 提示：如果变化，打印通知并更新 lastBalance

	fmt.Println("监控开始...")
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
