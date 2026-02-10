package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dapp-learning/ethclient/call-contract/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取配置（与 2.11 一致：必填项缺失则直接报错）
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		log.Fatal("错误: 请设置环境变量 CONTRACT_ADDRESS")
	}
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}

	// 练习 1：连接到以太坊节点并获取链 ID
	// 提示：ethclient.Dial(rpcURL)，再用 client.ChainID(context.Background())
	// client, err := ???
	// chainID, err := ???
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("已连接到以太坊节点")

	// 练习 2：从私钥创建交易认证器
	// 提示：crypto.HexToECDSA(privateKeyStr)，再 bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	// privateKey, err := ???
	// auth, err := ???
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 练习 3：加载合约实例
	contractAddress := common.HexToAddress(contractAddressStr)
	// storeContract, err := ???
	storeContract, err := store.NewStore(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合约实例加载成功: %s\n", contractAddress.Hex())

	// 从控制台读取 key 和 value
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入 key: ")
	keyInput, _ := reader.ReadString('\n')
	keyInput = strings.TrimSpace(keyInput)
	fmt.Print("请输入 value: ")
	valueInput, _ := reader.ReadString('\n')
	valueInput = strings.TrimSpace(valueInput)

	var key [32]byte
	var value [32]byte
	copy(key[:], []byte(keyInput))
	copy(value[:], []byte(valueInput))

	// 练习 4：调用 SetItem() 函数发送交易
	// 提示：使用 storeContract.SetItem(auth, key, value)
	// TODO: tx, err := ???
	// if err != nil { ... }
	// fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())
	start := time.Now()
	tx, err := storeContract.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易已发送: %s\n", tx.Hash().Hex())

	// 练习 5：等待交易确认
	// 提示：使用 bind.WaitMined(context.Background(), client, tx)
	// TODO: receipt, err := ???
	// if err != nil { ... }
	// fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	// fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("发送到确认耗时: %v\n", elapsed)
}
