// debug-transfer.go - 调试 transfer 调用
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
	"golang.org/x/crypto/sha3"
)

func main() {
	fmt.Println("=== 调试 Transfer 调用 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || toAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, TO_ADDRESS")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal("错误: 连接失败", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("错误: 解析私钥失败", err)
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	tokenAddress := common.HexToAddress(tokenAddressHex)
	toAddress := common.HexToAddress(toAddressHex)

	// 1. 查询发送方余额
	fmt.Println("\n1. 查询余额...")
	balance := queryBalance(client, tokenAddress, fromAddress)
	fmt.Printf("发送方余额: %s 代币\n", balance)

	bigFloat := new(big.Float)
	bigFloat.SetString(balance)
	balanceInt, _ := new(big.Int).SetString(balance, 10)
	fmt.Printf("发送方余额 (wei): %s\n", balanceInt.String())

	// 2. 构建 transfer 调用数据
	fmt.Println("\n2. 构建 transfer 调用...")
	amount := big.NewInt(1000) // 转账 1000 代币（最小单位）
	amountWei := new(big.Int).Mul(amount, big.NewInt(1e18))
	fmt.Printf("转账数量: %s wei\n", amountWei.String())

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte("transfer(address,uint256)"))
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amountWei.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	fmt.Printf("Method ID: %x\n", methodID)
	fmt.Printf("Data: %x\n", data)

	// 3. 尝试调用（不发送交易）
	fmt.Println("\n3. 模拟调用...")
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	}, nil)
	if err != nil {
		fmt.Printf("❌ 调用失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 调用成功，返回: %x\n", result)

	// 4. 估算 Gas
	fmt.Println("\n4. 估算 Gas...")
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		fmt.Printf("❌ 估算 Gas 失败: %v\n", err)
		return
	}
	fmt.Printf("✅ Gas Limit: %d\n", gasLimit)
}

func queryBalance(client *ethclient.Client, tokenContract, account common.Address) string {
	hash := crypto.Keccak256([]byte("balanceOf(address)"))
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(account.Bytes(), 32)
	data := append(methodID, paddedAddress...)

	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContract,
		Data: data,
	}, nil)
	if err != nil {
		return "0"
	}

	balance := new(big.Int).SetBytes(result)
	balanceFloat := new(big.Float).SetInt(balance)
	balanceFloat.Quo(balanceFloat, big.NewFloat(1e18))
	return balanceFloat.String()
}
