// 发送交易：连接 Sepolia，构造并发送 ETH 转账，输出交易哈希。
//
// 环境变量：SEPOLIA_RPC_URL、PRIVATE_KEY、TO_ADDRESS（必填）
// 可选：VALUE_ETH  默认 0.001
//
// 运行：在 task1 目录下执行 go run ./blockchain-read-write/send-eth
package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getRPCURL() string {
	u := strings.TrimSpace(os.Getenv("SEPOLIA_RPC_URL"))
	if u == "" {
		log.Fatal("请设置环境变量 SEPOLIA_RPC_URL")
	}
	return u
}

func parseEthToWei(ethStr string) *big.Int {
	eth, _ := strconv.ParseFloat(ethStr, 64)
	if eth <= 0 {
		eth = 0.001
	}
	oneEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	wei := new(big.Float).Mul(big.NewFloat(eth), new(big.Float).SetInt(oneEth))
	out, _ := wei.Int(nil)
	return out
}

func main() {
	fmt.Println("=== 发送 ETH 转账 (Sepolia) ===")

	privateKeyHex := strings.TrimSpace(os.Getenv("PRIVATE_KEY"))
	if privateKeyHex == "" {
		log.Fatal("请设置环境变量 PRIVATE_KEY")
	}
	toAddressHex := strings.TrimSpace(os.Getenv("TO_ADDRESS"))
	if toAddressHex == "" {
		log.Fatal("请设置环境变量 TO_ADDRESS")
	}

	client, err := ethclient.Dial(getRPCURL())
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("解析私钥失败:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型错误")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(toAddressHex)

	valueEth := os.Getenv("VALUE_ETH")
	if valueEth == "" {
		valueEth = "0.001"
	}
	value := parseEthToWei(valueEth)

	ctx := context.Background()
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal("获取 Nonce 失败:", err)
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("获取 Gas 价格失败:", err)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal("获取 Chain ID 失败:", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("签名失败:", err)
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatal("发送交易失败:", err)
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("接收方: %s\n", toAddress.Hex())
	fmt.Printf("金额: %s ETH\n", valueEth)
	fmt.Printf("交易哈希: %s\n", txHash)
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n", txHash)

	fmt.Println("\n等待交易确认...")
	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	for {
		receipt, err := client.TransactionReceipt(waitCtx, signedTx.Hash())
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if receipt.Status == 1 {
			fmt.Println("✅ 交易成功")
		} else {
			fmt.Println("❌ 交易失败")
		}
		break
	}
	fmt.Println("=== 完成 ===")
}
