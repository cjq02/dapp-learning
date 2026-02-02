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

	// 你的 Test Token 合约地址
	tokenAddress := common.HexToAddress("0x8087EcA92385db7a72e7Afbe0Eb6e2338cB17BDA")

	// 创建合约实例
	instance, err := erc20.NewIERC20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 指定要查询的地址（这里使用部署合约的地址）
	address := common.HexToAddress("0x8087EcA92385db7a72e7Afbe0Eb6e2338cB17BDA")

	// 查询代币余额
	bal, err := instance.BalanceOf(nil, address)
	if err != nil {
		log.Fatal(err)
	}

	// 查询代币信息
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

	// 将余额转换为可读格式
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	// 输出结果
	fmt.Printf("代币名称: %s\n", name)
	fmt.Printf("代币符号: %s\n", symbol)
	fmt.Printf("代币精度: %d\n", decimals)
	fmt.Printf("余额: %s %s\n", value.String(), symbol)
}
