package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/dapp-learning/ethclient/call-contract/store"
)

func main() {
	// 从环境变量获取配置
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal("请设置环境变量 PRIVATE_KEY")
	}

	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		contractAddressStr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
	}

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

	// 从私钥字符串创建私钥实例
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	// 创建交易认证器
	chainID := big.NewInt(11155111) // Sepolia
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ 交易认证器创建成功")

	// 加载合约实例
	contractAddress := common.HexToAddress(contractAddressStr)
	storeContract, err := store.NewStore(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 合约实例加载成功: %s\n", contractAddress.Hex())

	// 准备数据
	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("exercise_key"))
	copy(value[:], []byte("exercise_value"))

	// 调用 SetItem() 函数发送交易
	tx, err := storeContract.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 交易已发送: %s\n", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("✅ Gas 使用: %d\n", receipt.GasUsed)

	// 验证存储的值
	storedValue, err := storeContract.GetItem(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	if storedValue == value {
		fmt.Println("✅ 验证成功：存储的值与原始值一致")
	} else {
		fmt.Println("⚠️  注意：存储的值可能与原始值不同")
	}
}
