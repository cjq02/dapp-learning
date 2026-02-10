// 使用 abigen 生成的 counter 包向 Sepolia 部署 Counter 合约，并输出合约地址。
//
// 环境变量：SEPOLIA_RPC_URL、PRIVATE_KEY
//
// 运行（在 task1 目录）：go run ./contract-codegen/deploy-counter
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	counter "task1/contract-codegen/contract"
)

func main() {
	rpcURL := strings.TrimSpace(os.Getenv("SEPOLIA_RPC_URL"))
	if rpcURL == "" {
		log.Fatal("请设置环境变量 SEPOLIA_RPC_URL")
	}
	privateKeyHex := strings.TrimSpace(os.Getenv("PRIVATE_KEY"))
	if privateKeyHex == "" {
		log.Fatal("请设置环境变量 PRIVATE_KEY")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("解析私钥失败:", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("获取 Chain ID 失败:", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("创建交易认证器失败:", err)
	}

	start := time.Now()
	address, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal("部署合约失败:", err)
	}
	fmt.Printf("部署交易已发送: %s\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("等待交易确认失败:", err)
	}
	if receipt.Status != 1 {
		log.Fatal("部署交易执行失败")
	}

	elapsed := time.Since(start)
	fmt.Printf("合约已部署，地址: %s\n", address.Hex())
	fmt.Printf("部署耗时: %s\n", elapsed.Round(time.Millisecond))
}
