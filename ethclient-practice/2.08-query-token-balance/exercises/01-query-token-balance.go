package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/dapp-learning/ethclient/query-token-balance/erc20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
	// Infura: https://sepolia.infura.io/v3/YOUR_API_KEY
	// Alchemy: https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// TODO 3: 设置代币合约地址
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// TODO 4: 创建合约实例
	instance, err := erc20.NewErc20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// TODO 5: 从控制台输入要查询的地址
	var accountAddressHex string
	fmt.Print("请输入要查询的账户地址: ")
	fmt.Scanln(&accountAddressHex)
	if accountAddressHex == "" {
		log.Fatal("错误: 地址不能为空")
	}
	accountAddress := common.HexToAddress(accountAddressHex)

	// TODO 6: 查询代币余额
	balance, err := instance.BalanceOf(nil, accountAddress)
	if err != nil {
		log.Fatal(err)
	}

	// TODO 7: 查询代币信息（名称、符号、精度）
	name, _ := instance.Name(nil)
	symbol, _ := instance.Symbol(nil)
	decimals, _ := instance.Decimals(nil)

	// TODO 8: 将余额转换为可读格式
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	// TODO 9: 输出结果
	fmt.Printf("代币名称: %s\n", name)
	fmt.Printf("代币符号: %s\n", symbol)
	fmt.Printf("代币精度: %d\n", decimals)
	fmt.Printf("余额: %s %s\n", value.String(), symbol)
}
