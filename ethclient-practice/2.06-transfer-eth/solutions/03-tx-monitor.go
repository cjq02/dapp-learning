// 03-tx-monitor.go - 交易监控器 - 答案

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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 交易监控器 ===")

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

	// 连接并加载私钥
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

	// 获取 Nonce 和 Gas Price
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000)

	// 构建并发送交易
	value := big.NewInt(1000000000000000) // 0.001 ETH
	toAddress := common.HexToAddress(toAddressHex)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	txHash := signedTx.Hash()
	fmt.Printf("\n交易已发送: %s\n", txHash.Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n\n", txHash.Hex())

	// 开始监控
	fmt.Println("开始监控交易状态...")
	fmt.Println("────────────────────────────────────────")

	// 状态标志
	inMempool := false
	isPending := true
	isConfirmed := false
	var receipt *types.Receipt

	// 轮询监控
	for {
		// 检查交易是否在 mempool
		// 注意：有些 RPC 不支持 TransactionInPool，用查询交易代替
		isInMempool := false
		_, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			// 收据不存在，可能还在 mempool
			// 检查是否是 "not found" 错误
			if err == ethereum.NotFound {
				isInMempool = true
			}
		}

		if isInMempool && !inMempool {
			inMempool = true
			fmt.Printf("[%s] ✅ 交易在 Mempool 中\n", formatTime())
		}

		// 检查交易是否已确认
		if isPending {
			receipt, err = client.TransactionReceipt(context.Background(), txHash)
			if err == nil && receipt != nil {
				isPending = false
				isConfirmed = true

				fmt.Printf("\n[%s] 🎉 交易已打包！\n", formatTime())
				fmt.Printf("  区块号: %d\n", receipt.BlockNumber.Uint64())
				fmt.Printf("  区块哈希: %s\n", receipt.BlockHash.Hex())
				fmt.Printf("  交易索引: %d\n", receipt.TransactionIndex)

				// 检查状态
				if receipt.Status == 1 {
					fmt.Printf("\n[%s] ✅ 交易成功！\n", formatTime())
					fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
					fmt.Printf("  Gas Limit: %d\n", gasLimit)

					// 计算实际费用
					actualFee := new(big.Int).Mul(receipt.GasUsed, gasPrice)
					actualFeeEth := weiToEth(actualFee)
					fmt.Printf("  实际费用: %.6f ETH\n", actualFeeEth)

					// Gas 使用率
					gasEfficiency := float64(receipt.GasUsed) / float64(gasLimit) * 100
					fmt.Printf("  Gas 使用率: %.2f%%\n", gasEfficiency)

				} else {
					fmt.Printf("\n[%s] ❌ 交易失败\n", formatTime())
					fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
				}

				break
			}
		}

		// 等待一段时间再检查
		time.Sleep(5 * time.Second)
		fmt.Printf("[%s] ⏳ 等待确认...\n", formatTime())
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Println("=== 监控结束 ===")
}

// 辅助函数：格式化时间
func formatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 辅助函数：Wei 转 ETH
func weiToEth(wei *big.Int) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e18),
	)
}
