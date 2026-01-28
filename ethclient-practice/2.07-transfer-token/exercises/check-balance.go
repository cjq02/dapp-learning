// check-balance.go - 查询 ERC20 代币余额
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	if apiKey == "" || tokenAddress == "" || privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, TOKEN_ADDRESS, PRIVATE_KEY")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal("错误: 连接失败", err)
	}
	defer client.Close()

	tokenContract := common.HexToAddress(tokenAddress)

	// 获取您的地址
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("错误: 解析私钥失败", err)
	}
	userAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 调用 ERC20 合约的 balanceOf(address) 函数
	// 1. 计算 balanceOf(address) 的方法 ID
	hash := crypto.Keccak256([]byte("balanceOf(address)"))
	methodID := hash[:4]

	// 2. 将地址填充为 32 字节
	paddedAddress := common.LeftPadBytes(userAddress.Bytes(), 32)

	// 3. 组合数据
	data := append(methodID, paddedAddress...)

	// 4. 调用合约
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContract,
		Data: data,
	}, nil)
	if err != nil {
		log.Fatal("错误: 查询代币余额失败", err)
	}

	// 解析结果
	balance := new(big.Int).SetBytes(result)

	fmt.Printf("代币地址: %s\n", tokenAddress)
	fmt.Printf("您的地址: %s\n", userAddress.Hex())
	fmt.Printf("代币余额 (wei): %s\n", balance.String())

	// 转换为代币数量 (假设 18 位小数)
	balanceFloat := new(big.Float).SetInt(balance)
	balanceFloat.Quo(balanceFloat, big.NewFloat(1e18))
	fmt.Printf("代币余额: %s 代币\n", balanceFloat.String())
}
