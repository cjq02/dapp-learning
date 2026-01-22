// 02-batch-transfer.go - 批量转账练习
//
// 任务：
// 1. 定义多个接收地址和转账金额
// 2. 按顺序依次发送多笔转账
// 3. 每笔交易使用递增的 Nonce
// 4. 统计总转账金额和总 Gas 费用
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && go run exercises/02-batch-transfer.go

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
	Index     int
	TxHash    string
	Success   bool
	GasUsed   uint64
	Error     error
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

	// TODO 1: 连接并加载私钥
	var client *ethclient.Client
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	{
		// 在这里填写代码
	}

	defer client.Close()

	// TODO 2: 定义转账列表
	transfers := []Transfer{
		// 示例：需要替换为实际地址
		// {
		// 	To:     common.HexToAddress("0x..."),
		// 	Amount: big.NewInt(1000000000000000), // 0.001 ETH
		// },
	}

	if len(transfers) == 0 {
		fmt.Println("错误: 请定义至少一个转账目标")
		return
	}

	// TODO 3: 获取起始 Nonce
	var startNonce uint64
	{
		// 在这里填写代码
	}

	fmt.Printf("起始 Nonce: %d\n", startNonce)

	// TODO 4: 获取 Gas Price
	var gasPrice *big.Int
	{
		// 在这里填写代码
	}

	// TODO 5: 批量发送交易
	results := make([]TransferResult, len(transfers))
	var wg sync.WaitGroup

	for i, transfer := range transfers {
		wg.Add(1)
		go func(index int, t Transfer) {
			defer wg.Done()

			// TODO: 在这里填写发送交易的代码
			// 提示：使用 startNonce + uint64(index) 作为 Nonce
			// 提示：构建交易、签名、发送
			// 提示：将结果存入 results[index]

		}(i, transfer)
	}

	wg.Wait()

	// TODO 6: 输出结果
	fmt.Println("\n转账结果:")
	fmt.Println("────────────────────────────────────────")

	var totalAmount *big.Int
	var totalGasUsed uint64
	successCount := 0

	for _, r := range results {
		// TODO: 统计并输出每个转账的结果
	}

	fmt.Println("────────────────────────────────────────")
	fmt.Printf("总转账金额: %s Wei\n", totalAmount.String())
	fmt.Printf("总 Gas 使用: %d\n", totalGasUsed)
	fmt.Printf("成功/总数: %d/%d\n", successCount, len(transfers))

	fmt.Println("=== 完成 ===")
}
