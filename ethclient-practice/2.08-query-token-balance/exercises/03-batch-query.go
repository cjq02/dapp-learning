package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/dapp-learning/ethclient/query-token-balance/erc20"
)

// TokenInfo 代币信息结构
type TokenInfo struct {
	Address  common.Address
	Name     string
	Symbol   string
	Decimals uint8
}

// BalanceResult 余额查询结果
type BalanceResult struct {
	Address common.Address
	Balance *big.Int
	Error   error
}

// TODO: 实现将余额转换为可读格式的函数
func formatTokenBalance(balance *big.Int, decimals uint8) *big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	return new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
}

func main() {
	// TODO 1: 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	if tokenAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 TOKEN_ADDRESS")
	}

	// TODO 2: 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// TODO 3: 定义多个代币合约（可添加你的 TST 代币）
	// tokens := []TokenInfo{
	//     {common.HexToAddress("0x..."), "Token Name", "SYMBOL", 18},
	// }

	// TODO 4: 定义要查询的地址列表
	// addresses := []common.Address{
	//     common.HexToAddress("0x..."),
	//     common.HexToAddress("0x..."),
	// }

	// TODO 5: 使用 goroutine 并发查询所有代币余额
	// 提示：使用 sync.WaitGroup 和 sync.Mutex

	// TODO 6: 输出格式化的余额表格
	// fmt.Println("========== 多代币余额查询 ==========")

	// TODO 7: 计算总资产价值
	// fmt.Println("\n========== 总资产价值 ==========")

	// 移除下面的占位符代码，完成你的作业
	_ = sync.Mutex{}
	_ = sync.WaitGroup{}
	_ = fmt.Sprintf("")
	_ = common.HexToAddress("0x0")
	_ = erc20.IERC20{}
}
