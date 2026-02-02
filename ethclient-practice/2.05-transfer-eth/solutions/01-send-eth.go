// 01-send-eth.go - 发送 ETH 转账 - 答案

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== ETH 转账 ===")

	// 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}

	toAddressHex := os.Getenv("TO_ADDRESS")
	if toAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 TO_ADDRESS")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 加载私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// 获取发送方地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("发送方: %s\n", fromAddress.Hex())

	// 获取 Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Nonce: %d\n", nonce)

	// 设置转账金额（0.001 ETH）
	value := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e15)) // 0.001 ETH = 10^15 Wei

	// 设置 Gas 参数
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收地址
	toAddress := common.HexToAddress(toAddressHex)

	// 构建未签名交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 获取 Chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())

	// 等待交易确认
	fmt.Println("\n等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 简单轮询方式
	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			// 交易还未确认
			time.Sleep(10 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Printf("\n✅ 交易成功！\n")
			fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

			// 计算实际费用
			actualFee := new(big.Int).Mul(receipt.GasUsed, gasPrice)
			actualFeeEth := new(big.Float).Quo(
				new(big.Float).SetInt(actualFee),
				big.NewFloat(1e18),
			)
			fmt.Printf("实际费用: %.6f ETH\n", actualFeeEth)
		} else {
			fmt.Printf("\n❌ 交易失败！\n")
		}
		break
	}

	fmt.Println("=== 完成 ===")
}
