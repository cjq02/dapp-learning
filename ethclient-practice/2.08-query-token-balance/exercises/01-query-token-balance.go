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

func main() {
	// 连接到以太坊节点
	// Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
	// Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
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

	// TODO: 指定要查询的地址
	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")

	// TODO: 查询代币余额
	bal, err := instance.BalanceOf(nil, address)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: 查询代币信息
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

	// TODO: 将余额转换为可读格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	// 输出结果
	fmt.Printf("代币名称: %s\n", name)
	fmt.Printf("代币符号: %s\n", symbol)
	fmt.Printf("代币精度: %d\n", decimals)
	fmt.Printf("余额: %s %s\n", value.String(), symbol)
}
