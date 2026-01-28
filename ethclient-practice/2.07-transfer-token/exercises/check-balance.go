// 保存为 check-balance.go，然后运行 go run check-balance.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	client, _ := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	defer client.Close()

	tokenAddressParsed := common.HexToAddress(tokenAddress)

	// 获取您的地址
	privateKey, _ := crypto.HexToECDSA(privateKeyHex)
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 查询余额
	balance, err := client.BalanceAt(context.Background(), tokenAddressParsed, nil)
	if err != nil {
		log.Fatal("查询余额失败", err)
	}

	fmt.Printf("代币地址: %s\n", tokenAddress)
	fmt.Printf("您的地址: %s\n", fromAddress.Hex())
	fmt.Printf("您的余额: %s wei\n", balance.String())
	fmt.Printf("您的余额: %s 代币\n", new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(1e18),
	).String())
}
