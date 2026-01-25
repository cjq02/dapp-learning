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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {
	fmt.Println("=== 查询 ERC20 代币余额 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	targetAddressHex := os.Getenv("TARGET_ADDRESS")

	if apiKey == "" || tokenAddressHex == "" || targetAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, TOKEN_ADDRESS, TARGET_ADDRESS")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
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
	}

	// TODO 2: 查询代币小数位数
	// 函数签名: decimals()
	var decimals uint64
	{
		// 在这里填写代码
		// 提示：构建 decimals 函数调用数据
		// 提示：返回的 result 是 32 字节
	}

	// TODO 3: 转换余额为人类可读格式
	// 余额 / 10^decimals
	var readableBalance *big.Float
	{
		// 在这里填写代码
	}

	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("查询地址: %s\n", targetAddress.Hex())
	fmt.Printf("\n原始余额: %s\n", balance.String())
	fmt.Printf("小数位数: %d\n", decimals)
	fmt.Printf("可读余额: %s\n", readableBalance.String())

	fmt.Println("=== 完成 ===")
}

// 辅助函数：构建函数调用数据
func buildCallData(signature string, args ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(signature))
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)

	for _, arg := range args {
		padded := common.LeftPadBytes(arg, 32)
		data = append(data, padded...)
	}

	return data
}
