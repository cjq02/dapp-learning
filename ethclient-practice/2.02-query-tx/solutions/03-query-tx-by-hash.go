package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	// 连接 Sepolia 测试网络
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 交易哈希（可以通过命令行参数传入，这里硬编码示例）
	txHashStr := "0xb7cedb112cb9b246faed99ffdfd8bfcdca9215e32dc2b4d006d7afc529a0c625"
	if len(os.Args) > 1 {
		txHashStr = os.Args[1]
	}

	txHash := common.HexToHash(txHashStr)

	// 通过交易哈希查询交易
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal("查询交易失败: ", err)
	}

	if isPending {
		fmt.Println("⚠️  交易还在待处理队列中，尚未被打包")
		return
	}

	// 获取交易收据（包含区块号等信息）
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal("获取交易收据失败: ", err)
	}

	// 使用 MakeSigner 自动适配所有交易类型（Legacy, EIP-1559, EIP-4844 Blob）
	signer := types.MakeSigner(params.SepoliaChainConfig, receipt.BlockNumber, uint64(params.SepoliaChainConfig.ChainID.Uint64()))

	// 恢复发送者地址
	sender, err := types.Sender(signer, tx)
	if err != nil {
		log.Fatal(err)
	}

	// 转换金额为 Ether
	valueInEther := new(big.Float).Quo(
		new(big.Float).SetInt(tx.Value()),
		big.NewFloat(1e18),
	)

	// 格式化 Gas Price
	gasPriceStr := formatGasPrice(tx)

	// 显示交易详细信息
	fmt.Println("=== 交易详情 ===")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("区块位置: %d\n", receipt.TransactionIndex)
	fmt.Printf("发送者: %s\n", sender.Hex())

	to := "<合约创建>"
	if tx.To() != nil {
		to = tx.To().Hex()
	}
	fmt.Printf("接收者: %s\n", to)

	fmt.Printf("金额: %.6f ETH\n", valueInEther)
	fmt.Printf("Gas 限制: %d\n", tx.Gas())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("Gas 价格: %s\n", gasPriceStr)
	fmt.Printf("交易状态: %s\n", getStatus(receipt.Status))
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Printf("交易类型: %d", tx.Type())
	if tx.Type() == 2 {
		fmt.Print(" (EIP-1559)")
	} else if tx.Type() == 3 {
		fmt.Print(" (EIP-4844 Blob)")
	}
	fmt.Println()

	// 显示输入数据（如果是合约调用）
	if len(tx.Data()) > 0 {
		fmt.Printf("\n输入数据长度: %d 字节\n", len(tx.Data()))
		if len(tx.Data()) <= 100 {
			fmt.Printf("输入数据: %x\n", tx.Data())
		} else {
			fmt.Printf("输入数据 (前100字节): %x...\n", tx.Data()[:100])
		}
	}

	// 显示日志数量
	if len(receipt.Logs) > 0 {
		fmt.Printf("\n事件日志数量: %d\n", len(receipt.Logs))
	}
}

// 格式化 Gas Price 显示
func formatGasPrice(tx *types.Transaction) string {
	weiToGwei := big.NewFloat(1e9)

	if tx.Type() >= 2 {
		maxFee := new(big.Float).Quo(new(big.Float).SetInt(tx.GasFeeCap()), weiToGwei)
		priority := new(big.Float).Quo(new(big.Float).SetInt(tx.GasTipCap()), weiToGwei)
		result := fmt.Sprintf("MaxFee: %.2f, Priority: %.2f Gwei", maxFee, priority)
		if tx.Type() == 3 {
			result += " (Blob Tx)"
		}
		return result
	}

	gasPrice := new(big.Float).Quo(new(big.Float).SetInt(tx.GasPrice()), weiToGwei)
	return fmt.Sprintf("%.2f Gwei", gasPrice)
}

// 获取交易状态
func getStatus(status uint64) string {
	if status == 1 {
		return "✅ 成功"
	}
	return "❌ 失败"
}
