// 使用 abigen 生成的 counter 包与 Sepolia 上的 Counter 合约交互：调用 increment，再查询 getCount。
//
// 环境变量：SEPOLIA_RPC_URL、PRIVATE_KEY、COUNTER_CONTRACT_ADDRESS（已部署的合约地址）
//
// 运行（在 task1 目录）：go run ./contract-codegen/call-counter
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	addrHex := strings.TrimSpace(os.Getenv("COUNTER_CONTRACT_ADDRESS"))
	if addrHex == "" {
		log.Fatal("请设置环境变量 COUNTER_CONTRACT_ADDRESS")
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

	contractAddr := common.HexToAddress(addrHex)
	instance, err := counter.NewCounter(contractAddr, client)
	if err != nil {
		log.Fatal("加载合约失败:", err)
	}

	n, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("GetCount 失败:", err)
	}
	fmt.Printf("当前 count: %s\n", n.String())

	start := time.Now()
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal("Increment 失败:", err)
	}
	fmt.Printf("Increment 交易已发送: %s\n", tx.Hash().Hex())

	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal("等待交易确认失败:", err)
	}
	if receipt.Status != 1 {
		log.Fatal("交易执行失败")
	}

	elapsed := time.Since(start)
	fmt.Printf("交易耗时: %s\n", elapsed.Round(time.Millisecond))

	n, err = instance.GetCount(nil)
	if err != nil {
		log.Fatal("GetCount 失败:", err)
	}
	fmt.Printf("调用后 count: %s\n", n.String())
}
