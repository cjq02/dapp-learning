package main

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/dapp-learning/ethclient/query-token-balance/erc20"
)

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

	// Sepolia 测试网 USDC 合约地址
	tokenAddress := common.HexToAddress("0x94a9D9AC8a22534E3FaCa9F4e7F2E2cf85d5E4C8")

	// TODO: 创建合约实例
	instance, err := erc20.NewIERC20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: 查询代币信息（名称、符号、精度）
	name, err := instance.Name(nil)
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(nil)
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(nil)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: 查询总供应量
	totalSupply, err := instance.TotalSupply(nil)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: 查询多个地址的余额
	addresses := []common.Address{
		common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"),
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
	}

	// 输出代币信息
	fmt.Println("========== 代币信息 ==========")
	fmt.Printf("名称: %s\n", name)
	fmt.Printf("符号: %s\n", symbol)
	fmt.Printf("精度: %d\n", decimals)
	fmt.Printf("总供应量: %s %s\n\n", formatTokenBalance(totalSupply, decimals).String(), symbol)

	// 输出余额表格
	fmt.Println("========== 地址余额 ==========")
	for _, addr := range addresses {
		balance, err := instance.BalanceOf(nil, addr)
		if err != nil {
			log.Printf("查询地址 %s 余额失败: %v", addr.Hex(), err)
			continue
		}
		fmt.Printf("%s: %s %s\n", addr.Hex(), formatTokenBalance(balance, decimals).String(), symbol)
	}
}
