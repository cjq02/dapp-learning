package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BidPlaced 事件签名
// event BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount, bool isETH)
var (
	// 事件签名 = keccak256("BidPlaced(uint256,address,uint256,bool)")
	BidPlacedSignature = crypto.Keccak256Hash([]byte("BidPlaced(uint256,address,uint256,bool)"))
)

func main() {
	// 从环境变量读取 API Key
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY\n例如: export INFURA_API_KEY=your-key-here")
	}
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 替换为实际的交易哈希（包含 BidPlaced 事件的交易）
	txHash := common.HexToHash("0xb7cedb112cb9b246faed99ffdfd8bfcdca9215e32dc2b4d006d7afc529a0c625")

	// 查询交易收据
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== BidPlaced 事件解析 ===\n")

	// 遍历日志，查找 BidPlaced 事件
	for i, log := range receipt.Logs {
		// Topic[0] 是事件签名
		if len(log.Topics) == 0 {
			continue
		}

		eventSignature := log.Topics[0]
		fmt.Printf("[日志 %d] 事件签名: %s\n", i+1, eventSignature.Hex())

		// 检查是否是 BidPlaced 事件
		if eventSignature != BidPlacedSignature {
			fmt.Println("  → 不是 BidPlaced 事件，跳过\n")
			continue
		}

		fmt.Println("  → 找到 BidPlaced 事件！")

		// 解析 indexed 参数
		// Topic[1] = auctionId (uint256 indexed)
		// Topic[2] = bidder (address indexed)
		if len(log.Topics) >= 3 {
			auctionId := new(big.Int).SetBytes(log.Topics[1].Bytes())
			bidder := common.BytesToAddress(log.Topics[2].Bytes())

			fmt.Printf("\n  [Indexed 参数]\n")
			fmt.Printf("    auctionId: %s\n", auctionId.String())
			fmt.Printf("    bidder: %s\n", bidder.Hex())
		}

		// 解析非 indexed 参数（在 Data 中）
		// Data = amount (uint256, 32 bytes) + isETH (bool, 1 byte)
		if len(log.Data) >= 32 {
			amount := new(big.Int).SetBytes(log.Data[0:32])

			var isETH bool
			if len(log.Data) > 32 {
				isETH = log.Data[32] == 1
			}

			// 转换为 Ether
			amountInEther := new(big.Float).Quo(
				new(big.Float).SetInt(amount),
				big.NewFloat(1e18),
			)

			fmt.Printf("\n  [非 Indexed 参数]\n")
			fmt.Printf("    amount: %s Wei (%.6f ETH)\n", amount.String(), amountInEther)
			fmt.Printf("    isETH: %v\n", isETH)
		}

		// 显示原始数据（调试用）
		fmt.Printf("\n  [原始数据]\n")
		fmt.Printf("    Topics: %d 个\n", len(log.Topics))
		for j, topic := range log.Topics {
			fmt.Printf("      Topic[%d]: %s\n", j, topic.Hex())
		}
		fmt.Printf("    Data (hex): %s\n", hex.EncodeToString(log.Data))
		fmt.Println()
	}

	fmt.Println("=== 解析规则说明 ===")
	fmt.Println("事件: BidPlaced(uint256 indexed auctionId, address indexed bidder, uint256 amount, bool isETH)")
	fmt.Println("")
	fmt.Println("参数位置:")
	fmt.Println("  Topic[0] = 事件签名 (keccak256)")
	fmt.Println("  Topic[1] = auctionId (indexed)")
	fmt.Println("  Topic[2] = bidder (indexed)")
	fmt.Println("  Data[0:32] = amount (uint256)")
	fmt.Println("  Data[32] = isETH (bool, 0x01=true, 0x00=false)")
}
