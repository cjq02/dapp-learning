// cancel-tx.go - 取消或加速 pending 交易
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 取消/加速 Pending 交易 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	if apiKey == "" || privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal("错误: 连接到以太坊节点失败", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("错误: 解析私钥失败", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 获取当前 pending nonce（这会包含 pending 交易的 nonce）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("错误: 获取 Nonce 失败", err)
	}

	// 旧交易的 nonce（需要手动设置，通常比当前 nonce 小 1）
	// 如果您的旧交易是该账户的第 0 笔交易，那么 oldNonce = 0
	var oldNonce uint64
	fmt.Print("请输入要取消的交易的 nonce (通常是当前 nonce - 1): ")
	fmt.Scanf("%d", &oldNonce)

	fmt.Printf("当前 pending nonce: %d\n", nonce)
	fmt.Printf("要取消的交易 nonce: %d\n", oldNonce)

	// 设置更高的 Gas Price (至少比原交易高 10%)
	// 原交易是 0.923 Gwei，我们用 20 Gwei 来加速
	gasPrice := big.NewInt(20000000000) // 20 Gwei
	fmt.Printf("新 Gas Price: %.2f Gwei\n", new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9)))

	// 创建一个取消交易：发送 0 ETH 到自己
	// 这会使用相同的 nonce，但 gas price 更高，从而替换旧交易
	tx := types.NewTransaction(oldNonce, fromAddress, big.NewInt(0), 21000, gasPrice, nil)

	// 签名
	chainID := big.NewInt(11155111) // Sepolia
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("错误: 签名失败", err)
	}

	// 发送取消交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("错误: 发送取消交易失败", err)
	}

	fmt.Printf("\n✅ 取消交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
	fmt.Println("\n这个交易会替换掉旧的 pending 交易。")
	fmt.Println("一旦确认，旧交易就会失效，您就可以使用新的 nonce 发送交易了。")
}
