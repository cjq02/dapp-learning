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

	blockNumber := big.NewInt(5671744)

	// 练习：获取区块和所有交易
	// block, err := ???

	// 练习：获取链 ID，用于恢复发送者地址
	// chainID, err := ???

	// 练习：创建 EIP155 签名器
	// signer := types.NewEIP155Signer(???)

	// 练习：遍历所有交易并显示信息
	fmt.Println("=== 交易列表 ===")
	// for _, tx := range block.Transactions() {
	//     // 恢复发送者地址
	//     sender, err := types.Sender(signer, tx)
	//     if err != nil {
	//         log.Fatal(err)
	//     }
	//
	//     // TODO: 转换金额为 Ether (1 Ether = 10^18 Wei)
	//     // valueInEther := new(big.Float).Quo(...)
	//
	//     // TODO: 转换 Gas Price 为 Gwei (1 Gwei = 10^9 Wei)
	//     // gasPriceInGwei := new(big.Float).Quo(...)
	//
	//     fmt.Printf("Hash: %s\n", tx.Hash().Hex())
	//     fmt.Printf("From: %s\n", sender.Hex())
	//     fmt.Printf("To: %s\n", tx.To().Hex())
	//     fmt.Printf("Value: %s Ether\n", valueInEther)
	//     fmt.Printf("Gas Price: %s Gwei\n", gasPriceInGwei)
	//     fmt.Println("---")
	// }
}
