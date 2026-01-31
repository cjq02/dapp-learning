package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/dapp-learning/ethclient/call-contract/store"
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

	fmt.Println("✅ 已连接到以太坊节点")

	// 加载合约实例
	contractAddress := common.HexToAddress(contractAddressStr)
	storeContract, err := store.NewStore(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 合约实例加载成功: %s\n", contractAddress.Hex())

	// 调用 Version() 函数获取合约版本
	version, err := storeContract.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 合约版本: %s\n", version)

	// 调用 GetItem() 函数获取存储的值
	var key [32]byte
	copy(key[:], []byte("demo_key"))

	value, err := storeContract.GetItem(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 存储值: %x\n", value)
}
