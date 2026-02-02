package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
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

func formatTokenBalance(balance *big.Int, decimals uint8) *big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	return new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
}

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// TODO: 定义多个代币合约
	tokens := []TokenInfo{
		{common.HexToAddress("0x94a9D9AC8a22534E3FaCa9F4e7F2E2cf85d5E4C8"), "USD Coin", "USDC", 6},
		// 可以添加更多代币
	}

	// 要查询的地址列表
	addresses := []common.Address{
		common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"),
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}

	// TODO: 并发查询所有代币余额
	results := make(map[string]map[common.Address]*big.Int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, token := range tokens {
		// 为每个代币创建查询任务
		token := token // 避免闭包问题
		results[token.Symbol] = make(map[common.Address]*big.Int)

		for _, addr := range addresses {
			addr := addr // 避免闭包问题
			wg.Add(1)

			go func(t TokenInfo, a common.Address) {
				defer wg.Done()

				instance, err := erc20.NewIERC20(t.Address, client)
				if err != nil {
					log.Printf("创建合约实例失败: %v", err)
					return
				}

				balance, err := instance.BalanceOf(nil, a)
				if err != nil {
					log.Printf("查询 %s 余额失败: %v", a.Hex(), err)
					return
				}

				mu.Lock()
				results[t.Symbol][a] = balance
				mu.Unlock()
			}(token, addr)
		}
	}

	wg.Wait()

	// TODO: 输出格式化的余额表格
	fmt.Println("========== 多代币余额查询 ==========\n")

	// 打印表头
	fmt.Printf("%-45s", "Address")
	for _, token := range tokens {
		fmt.Printf("%-15s", token.Symbol)
	}
	fmt.Println("\n" + string(make([]byte, 45+len(tokens)*15)))

	// 打印每个地址的余额
	for _, addr := range addresses {
		fmt.Printf("%-45s", addr.Hex())
		for _, token := range tokens {
			if balance, ok := results[token.Symbol][addr]; ok {
				fmt.Printf("%-15s", formatTokenBalance(balance, token.Decimals).String())
			} else {
				fmt.Printf("%-15s", "N/A")
			}
		}
		fmt.Println()
	}

	// TODO: 计算总资产价值（假设 USDC 价格 = $1）
	fmt.Println("\n========== 总资产价值 ==========")
	totalValue := big.NewFloat(0)
	for _, addr := range addresses {
		addrValue := big.NewFloat(0)
		for _, token := range tokens {
			if balance, ok := results[token.Symbol][addr]; ok {
				value := formatTokenBalance(balance, token.Decimals)
				addrValue.Add(addrValue, value)
			}
		}
		fmt.Printf("%s: $%.2f USD\n", addr.Hex(), addrValue)
		totalValue.Add(totalValue, addrValue)
	}
	fmt.Printf("\n总计: $%.2f USD\n", totalValue)
}
