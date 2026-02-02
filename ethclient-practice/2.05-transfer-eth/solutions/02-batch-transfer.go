// 02-batch-transfer.go - 批量转账 - 答案

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transfer struct {
	To     common.Address
	Amount *big.Int
}

type TransferResult struct {
	Index   int
	TxHash  string
	Success bool
	GasUsed uint64
	Error   error
}

func main() {
	fmt.Println("=== 批量 ETH 转账 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
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

	// 定义转账列表（示例地址）
	transfers := []Transfer{
		{
			To:     common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F"),
			Amount: big.NewInt(1000000000000000), // 0.001 ETH
		},
		{
			To:     common.HexToAddress("0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359"),
			Amount: big.NewInt(1000000000000000), // 0.001 ETH
		},
	}

	if len(transfers) == 0 {
		fmt.Println("错误: 请定义至少一个转账目标")
		return
	}

	// 获取起始 Nonce
	startNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("起始 Nonce: %d\n", startNonce)

	// 获取 Chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 获取 Gas Price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Gas Price: %s Wei\n", gasPrice.String())

	// 批量发送交易
	results := make([]TransferResult, len(transfers))
	var wg sync.WaitGroup

	for i, transfer := range transfers {
		wg.Add(1)
		go func(index int, t Transfer) {
			defer wg.Done()

			// 使用递增的 Nonce
			nonce := startNonce + uint64(index)

			// 构建交易
			tx := types.NewTransaction(nonce, t.To, t.Amount, 21000, gasPrice, nil)

			// 签名交易
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
			if err != nil {
				results[index] = TransferResult{
					Index:   index,
					Success: false,
					Error:   err,
				}
				return
			}

			// 发送交易
			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				results[index] = TransferResult{
					Index:   index,
					Success: false,
					Error:   err,
				}
				return
			}

			results[index] = TransferResult{
				Index:   index,
				TxHash:  signedTx.Hash().Hex(),
				Success: true,
				GasUsed: 21000, // ETH 转账固定值
			}

			fmt.Printf("[%d] 交易已发送: %s\n", index+1, signedTx.Hash().Hex())
		}(i, transfer)
	}

	wg.Wait()

	// 输出结果
	fmt.Println("\n转账结果:")
	fmt.Println("────────────────────────────────────────")

	var totalAmount *big.Int
	var totalGasUsed uint64
	successCount := 0

	totalAmount = big.NewInt(0)

	for i, r := range results {
		if r.Success {
			totalAmount.Add(totalAmount, transfers[i].Amount)
			totalGasUsed += r.GasUsed
			successCount++
			fmt.Printf("[%d] ✅ 成功 | Hash: %s | Gas: %d\n", r.Index+1, r.TxHash, r.GasUsed)
		} else {
			fmt.Printf("[%d] ❌ 失败 | Error: %v\n", r.Index+1, r.Error)
		}
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Printf("总转账金额: %s Wei\n", totalAmount.String())
	fmt.Printf("总 Gas 使用: %d\n", totalGasUsed)

	// 计算总费用
	totalFee := new(big.Int).Mul(big.NewInt(int64(totalGasUsed)), gasPrice)
	totalFeeEth := new(big.Float).Quo(new(big.Float).SetInt(totalFee), big.NewFloat(1e18))
	fmt.Printf("总 Gas 费用: %.6f ETH\n", totalFeeEth)
	fmt.Printf("成功/总数: %d/%d\n", successCount, len(transfers))

	fmt.Println("=== 完成 ===")
}
