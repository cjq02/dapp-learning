package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dapp-learning/ethclient/call-contract/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取配置（与 2.11 一致：必填项缺失则直接报错）
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		log.Fatal("错误: 请设置环境变量 CONTRACT_ADDRESS（可填 2.10 部署得到的合约地址）")
	}
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}

	// 练习 1：连接到以太坊节点
	// 提示：使用 ethclient.Dial(rpcURL)
	// client, err := ???
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("已连接到以太坊节点")

	// 练习 2：加载合约实例
	// 提示：使用 store.NewStore(address, client)
	contractAddress := common.HexToAddress(contractAddressStr)
	// storeContract, err := ???
	storeContract, err := store.NewStore(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合约实例加载成功: %s\n", contractAddress.Hex())

	// 练习 3：调用 Version() 函数获取合约版本
	// 提示：调用 storeContract.Version(nil)
	// version, err := ???
	version, err := storeContract.Version(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合约版本: %s\n", version)

	// 练习 4：调用 GetItem() 函数获取存储的值
	// 提示：传入一个 32 字节的 key，使用 storeContract.GetItem(nil, key)
	fmt.Print("请输入 key: ")
	reader := bufio.NewReader(os.Stdin)
	keyStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	keyStr = strings.TrimSpace(keyStr)
	var key [32]byte
	copy(key[:], []byte(keyStr))

	// TODO: 调用 GetItem 获取存储值并打印
	// value, err := ???
	// if err != nil { ... }
	// fmt.Printf("存储值: %x\n", value)
	value, err := storeContract.GetItem(nil, key)
	if err != nil {
		log.Fatal(err)
	}
	// bytes32 存的是字符串时：去掉尾部 0 再转成字符串
	valueStr := strings.TrimRight(string(value[:]), "\x00")
	fmt.Printf("存储值(hex): %x\n", value)
	fmt.Printf("存储值: %s\n", valueStr)

}
