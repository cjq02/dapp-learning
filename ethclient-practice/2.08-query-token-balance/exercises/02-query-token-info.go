package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dapp-learning/ethclient-practice/util"

	"github.com/dapp-learning/ethclient/query-token-balance/erc20"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO: 实现将余额转换为可读格式的函数
// 提示：使用 big.Float 进行除法运算

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

	// TODO 3: 设置代币合约地址
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// TODO 4: 创建合约实例
	instance, err := erc20.NewErc20(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// TODO 5: 查询代币信息（名称、符号、精度）
	name, _ := instance.Name(nil)
	symbol, _ := instance.Symbol(nil)
	decimals, _ := instance.Decimals(nil)

	// TODO 6: 查询总供应量
	totalSupply, _ := instance.TotalSupply(nil)

	// TODO 7: 定义要查询的地址列表
	addresses := []common.Address{
		common.HexToAddress("0x085f0145202298585e699371eb3CFb1441f65110"),
		common.HexToAddress("0x354C393Daf549Da43485FFe85Be464d82149B0e8"),
		common.HexToAddress("0x539F128f1Cd5877cA3712ba33f9E78EdFcC7eFAD"),
		common.HexToAddress("0xd2957b711c6B2291F937471b40ee1272dE171beD"),
	}

	// TODO 8: 输出代币信息
	fmt.Println("========== 代币信息 ==========")
	fmt.Printf("名称: %s\n", name)
	fmt.Printf("符号: %s\n", symbol)
	fmt.Printf("精度: %d\n", decimals)
	fmt.Printf("总供应量: %s %s\n\n", util.FormatTokenBalance(totalSupply, decimals).String(), symbol)

	// TODO 9: 遍历地址列表，查询并输出每个地址的余额
	for _, addr := range addresses {
		balance, err := instance.BalanceOf(nil, addr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s %s\n", addr.Hex(), util.FormatTokenBalance(balance, decimals).String(), symbol)
	}
}
