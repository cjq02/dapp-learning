// 02-query-token-balance.go - 查询代币余额 - 答案

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

	// 查询代币余额 - balanceOf(address)
	balanceOfSignature := []byte("balanceOf(address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(balanceOfSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(targetAddress.Bytes(), 32)

	var balanceOfData []byte
	balanceOfData = append(balanceOfData, methodID...)
	balanceOfData = append(balanceOfData, paddedAddress...)

	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: balanceOfData,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	balance := new(big.Int).SetBytes(result)

	// 查询代币小数位数 - decimals()
	decimalsSignature := []byte("decimals()")
	hash = sha3.NewLegacyKeccak256()
	hash.Write(decimalsSignature)
	methodID = hash.Sum(nil)[:4]

	var decimalsData []byte
	decimalsData = append(decimalsData, methodID...)

	result, err = client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: decimalsData,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	decimals := new(big.Int).SetBytes(result).Uint64()

	// 转换余额为人类可读格式
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	readableBalance := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		divisor,
	)

	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("查询地址: %s\n", targetAddress.Hex())
	fmt.Printf("\n原始余额: %s\n", balance.String())
	fmt.Printf("小数位数: %d\n", decimals)
	fmt.Printf("可读余额: %s\n", readableBalance.String())

	fmt.Println("=== 完成 ===")
}
