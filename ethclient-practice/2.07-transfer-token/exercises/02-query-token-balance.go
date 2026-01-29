// 02-query-token-balance.go - 查询代币余额练习
//
// 任务：
// 1. 调用 balanceOf(address) 查询代币余额
// 2. 调用 decimals() 获取代币小数位数
// 3. 将余额转换为人类可读格式
//
// 运行：export INFURA_API_KEY=your-key && export TOKEN_ADDRESS=0x... && export TARGET_ADDRESS=0x... && go run exercises/02-query-token-balance.go

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/dapp-learning/ethclient/transfer-token/util"
)

func main() {
	fmt.Println("=== 查询 ERC20 代币余额 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	targetAddressHex := os.Getenv("TARGET_ADDRESS")

	if apiKey == "" || tokenAddressHex == "" || targetAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, TOKEN_ADDRESS, TARGET_ADDRESS")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	tokenAddress := common.HexToAddress(tokenAddressHex)
	targetAddress := common.HexToAddress(targetAddressHex)

	// TODO 1: 查询代币余额
	// 函数签名: balanceOf(address)
	var balance *big.Int
	{
		// 在这里填写代码
		// 提示：构建 balanceOf 函数调用数据
		// 提示：使用 client.CallContract() 调用合约
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: util.BuildCallData("balanceOf(address)", targetAddress.Bytes()),
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询代币余额失败", err)
		}
		balance = new(big.Int).SetBytes(result)
	}

	// TODO 2: 查询代币小数位数
	// 函数签名: decimals()
	var decimals uint64
	{
		// 在这里填写代码
		// 提示：构建 decimals 函数调用数据
		// 提示：返回的 result 是 32 字节
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: util.BuildCallData("decimals()"),
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询代币小数位数失败", err)
		}
		decimals = new(big.Int).SetBytes(result).Uint64()
	}

	// TODO 3: 转换余额为人类可读格式
	// 余额 / 10^decimals
	var readableBalance *big.Float
	{
		// 在这里填写代码
		readableBalance = util.WeiToTokenAmount(balance, decimals)
	}

	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("查询地址: %s\n", targetAddress.Hex())
	fmt.Printf("\n原始余额: %s\n", balance.String())
	fmt.Printf("小数位数: %d\n", decimals)
	fmt.Printf("可读余额: %s\n", readableBalance.String())

	fmt.Println("=== 完成 ===")
}
