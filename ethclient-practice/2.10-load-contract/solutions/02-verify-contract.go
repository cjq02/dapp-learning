package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/dapp-learning/ethclient/load-contract/store"
)

func main() {
	// 从环境变量获取合约地址
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		// 默认测试地址
		contractAddressStr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
	}

	// 从环境变量获取 RPC URL
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(contractAddressStr)

	// 检查地址是否有合约代码
	code, err := client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 验证合约代码长度
	if len(code) == 0 {
		log.Fatal("❌ 该地址没有合约代码")
	}

	fmt.Printf("✅ 合约代码长度: %d 字节\n", len(code))

	// 尝试加载合约并捕获错误
	storeContract, err := store.NewStore(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 合约加载成功: %s\n", contractAddress.Hex())

	// 验证：尝试调用合约函数
	version, err := storeContract.Version(nil)
	if err != nil {
		log.Fatal("❌ 无法调用合约函数")
	}

	fmt.Printf("✅ 合约版本: %s\n", version)
	fmt.Println("✅ 合约验证通过！")
}
