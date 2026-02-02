// 03-query-allowance.go - 查询 ERC20 授权额度练习
//
// 任务：
// 1. 调用 allowance(owner, spender) 查询授权额度
// 2. 将额度转换为人类可读格式
//
// 运行：export INFURA_API_KEY=your-key && export TOKEN_ADDRESS=0x... && go run exercises/03-query-allowance.go
//
// 说明：
// - TOKEN_ADDRESS: ERC20 代币合约地址
// - OWNER_ADDRESS: 从控制台输入代币持有者地址
// - SPENDER_ADDRESS: 从控制台输入被授权者地址

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

	"github.com/dapp-learning/ethclient/util"
)

func main() {
	fmt.Println("=== 查询 ERC20 授权额度 ===")

	// 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")

	if apiKey == "" || tokenAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, TOKEN_ADDRESS")
	}

	// 从控制台输入代币持有者地址
	var ownerAddressHex string
	fmt.Print("请输入代币持有者地址 (owner): ")
	fmt.Scanln(&ownerAddressHex)

	if ownerAddressHex == "" {
		log.Fatal("错误: 代币持有者地址不能为空")
	}

	// 从控制台输入被授权者地址
	var spenderAddressHex string
	fmt.Print("请输入被授权者地址 (spender): ")
	fmt.Scanln(&spenderAddressHex)

	if spenderAddressHex == "" {
		log.Fatal("错误: 被授权者地址不能为空")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	tokenAddress := common.HexToAddress(tokenAddressHex)
	ownerAddress := common.HexToAddress(ownerAddressHex)
	spenderAddress := common.HexToAddress(spenderAddressHex)

	// TODO 1: 查询代币小数位数
	// 函数签名: decimals()
	var decimals uint64
	{
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: util.BuildCallData("decimals()"),
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询代币小数位数失败", err)
		}
		decimals = new(big.Int).SetBytes(result).Uint64()
	}

	// TODO 2: 查询授权额度
	// 函数签名: allowance(address owner, address spender)
	var allowance *big.Int
	{
		// 在这里填写代码
		// 提示：构建 allowance 函数调用数据
		// 提示：使用 client.CallContract() 调用合约
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: util.BuildCallData("allowance(address,address)", ownerAddress.Bytes(), spenderAddress.Bytes()),
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询授权额度失败", err)
		}
		allowance = new(big.Int).SetBytes(result)
	}

	// TODO 3: 转换授权额度为人类可读格式
	var readableAllowance *big.Float
	{
		// 在这里填写代码
		// 提示：使用 util.WeiToTokenAmount() 转换
		readableAllowance = util.WeiToTokenAmount(allowance, decimals)
	}

	fmt.Printf("\n=== 查询结果 ===\n")
	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("代币持有者 (owner): %s\n", ownerAddress.Hex())
	fmt.Printf("被授权者 (spender): %s\n", spenderAddress.Hex())
	fmt.Printf("\n原始授权额度: %s wei\n", allowance.String())
	fmt.Printf("小数位数: %d\n", decimals)
	fmt.Printf("可读授权额度: %s\n", readableAllowance.String())

	fmt.Println("\n=== 完成 ===")
}
