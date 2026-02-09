package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dapp-learning/ethclient-practice/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}

	txHashHex := os.Getenv("TX_HASH")
	if txHashHex == "" {
		log.Fatal("错误: 请设置环境变量 TX_HASH（例如 0xb230...）")
	}
	txHashHex = strings.TrimSpace(txHashHex)
	if err := validateTxHash(txHashHex); err != nil {
		log.Fatal(err)
	}
	txHash := common.HexToHash(txHashHex)

	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rpcClient.Close()
	client := ethclient.NewClient(rpcClient)

	ctx := context.Background()
	fmt.Printf("查询交易回执: %s\n", txHash.Hex())
	receipt, err := util.WaitForReceipt(ctx, client, txHash, 2*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("已获取交易回执")
	fmt.Printf("状态 Status: %d (1=成功,0=失败)\n", receipt.Status)
	fmt.Printf("区块号 BlockNumber: %s\n", receipt.BlockNumber.String())
	fmt.Printf("GasUsed: %d\n", receipt.GasUsed)
	fmt.Printf("合约地址 ContractAddress: %s\n", receipt.ContractAddress.Hex())
	fmt.Printf("交易索引 TxIndex: %d\n", receipt.TransactionIndex)

	if receipt.Status == 0 {
		msg := getRevertReason(ctx, rpcClient, txHash)
		if msg != "" {
			fmt.Printf("失败原因: %s\n", msg)
		} else {
			fmt.Println("失败原因: （当前节点可能未开放 debug 接口，无法获取；可到 Etherscan 查看该交易的错误详情）")
		}
	}
}

// getRevertReason 通过 debug_traceTransaction 获取交易失败原因（需节点支持 debug 接口）
func getRevertReason(ctx context.Context, rpcClient *rpc.Client, txHash common.Hash) string {
	// 1) 默认 tracer，部分 Geth 节点在顶层返回 error
	var result map[string]interface{}
	err := rpcClient.CallContext(ctx, &result, "debug_traceTransaction", txHash.Hex(), map[string]interface{}{})
	if err != nil {
		return ""
	}
	if e := extractErrorString(result); e != "" {
		return e
	}
	// 2) 使用 callTracer，部分节点在 callTracer 结果中返回 error
	var callResult map[string]interface{}
	err = rpcClient.CallContext(ctx, &callResult, "debug_traceTransaction", txHash.Hex(), map[string]interface{}{
		"tracer": "callTracer",
	})
	if err != nil {
		return ""
	}
	return extractErrorString(callResult)
}

// extractErrorString 从 trace 结果中提取 error 字符串（支持顶层或嵌套在 calls 中）
func extractErrorString(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	if e, ok := m["error"].(string); ok && e != "" {
		return strings.TrimSpace(e)
	}
	// callTracer 可能把 error 放在 calls[0].error
	if calls, ok := m["calls"].([]interface{}); ok && len(calls) > 0 {
		if first, ok := calls[0].(map[string]interface{}); ok {
			if e, ok := first["error"].(string); ok && e != "" {
				return strings.TrimSpace(e)
			}
		}
	}
	return ""
}

func validateTxHash(txHashHex string) error {
	s := strings.TrimSpace(txHashHex)
	if strings.HasPrefix(s, "0X") {
		s = "0x" + s[2:]
	}
	if !strings.HasPrefix(s, "0x") {
		return fmt.Errorf("错误: TX_HASH 必须以 0x 开头: %q", txHashHex)
	}
	if len(s) != 66 {
		return fmt.Errorf("错误: TX_HASH 长度应为 66（0x + 64 hex），当前=%d: %q", len(s), txHashHex)
	}
	_, err := hex.DecodeString(s[2:])
	if err != nil {
		return fmt.Errorf("错误: TX_HASH 不是合法十六进制: %q, err=%v", txHashHex, err)
	}
	return nil
}
