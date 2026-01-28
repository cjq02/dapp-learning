// 03-approve-transfer.go - 授权并代理转账练习
//
// 任务：
// 1. 调用 approve(spender, amount) 授权
// 2. 等待授权交易确认
// 3. 使用 transferFrom(from, to, amount) 进行代理转账
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TOKEN_ADDRESS=0x... && export SPENDER_ADDRESS=0x... && export TO_ADDRESS=0x... && export TOKEN_AMOUNT=... && go run exercises/03-approve-transfer.go

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// tokenAmountToWei 将人类可读的代币数量转换为最小单位（Wei）
func tokenAmountToWei(amount float64, decimals uint64) *big.Int {
	decimalsBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountFloat := big.NewFloat(amount)
	wei := new(big.Float).Mul(amountFloat, new(big.Float).SetInt(decimalsBig))
	result, _ := wei.Int(nil)
	return result
}

func main() {
	fmt.Println("=== ERC20 授权并代理转账 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	spenderAddressHex := os.Getenv("SPENDER_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")
	amountStr := os.Getenv("TOKEN_AMOUNT")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || spenderAddressHex == "" || toAddressHex == "" || amountStr == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, SPENDER_ADDRESS, TO_ADDRESS, TOKEN_AMOUNT")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tokenAddress := common.HexToAddress(tokenAddressHex)
	spenderAddress := common.HexToAddress(spenderAddressHex)
	toAddress := common.HexToAddress(toAddressHex)

	// 查询代币小数位数（默认 18）
	var decimals uint64 = 18
	{
		// TODO: 可选，查询代币的 decimals()
		// 如果查询失败，使用默认值 18
	}

	// 将人类可读的数量转换为最小单位
	var amountFloat float64
	_, err := fmt.Sscanf(amountStr, "%f", &amountFloat)
	if err != nil {
		log.Fatalf("错误: 无法解析代币数量 %s: %v", amountStr, err)
	}
	amount := tokenAmountToWei(amountFloat, decimals)
	fmt.Printf("代币数量: %s (decimals: %d) = %s\n", amountStr, decimals, amount.String())

	// TODO 1: 构建 approve 函数调用数据
	// 函数签名: approve(address,uint256)
	var approveData []byte
	{
		// 在这里填写代码
		// 提示：生成 methodID，填充 spenderAddress 和 amount
	}

	fmt.Printf("授权地址: %s\n", spenderAddress.Hex())
	fmt.Printf("授权金额: %s\n", amount.String())

	// TODO 2: 发送 approve 交易
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, _ := client.SuggestGasPrice(context.Background())

	value := big.NewInt(0)
	gasLimit, _ := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: approveData,
	})

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, approveData)

	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n授权交易已发送: %s\n", signedTx.Hash().Hex())

	// TODO 3: 等待授权交易确认
	fmt.Println("等待授权交易确认...")
	// 提示：轮询查询交易收据，直到确认

	fmt.Println("\n✅ 授权已确认！")

	// TODO 4: 构建 transferFrom 函数调用数据
	// 函数签名: transferFrom(address,address,uint256)
	var transferFromData []byte
	{
		// 在这里填写代码
		// 提示：fromAddress, toAddress, amount 各填充 32 字节
	}

	// TODO 5: 发送 transferFrom 交易
	nonce, _ = client.PendingNonceAt(context.Background(), fromAddress)
	// 提示：类似 approve 交易的发送流程

	fmt.Printf("\n代理转账交易已发送: %s\n", signedTx.Hash().Hex())

	// TODO 6: 等待转账交易确认
	// 提示：轮询查询交易收据

	fmt.Println("\n✅ 代理转账完成！")
	fmt.Println("=== 完成 ===")
}
