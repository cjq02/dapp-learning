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
	// 从环境变量获取 RPC URL
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"
	}

	// 练习 1：定义多个合约地址（有效和无效的）
	// 提示：创建一个字符串数组，包含有效的合约地址和无效的地址
	contractAddresses := []string{
		"0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e", // 有效地址
		"0x0000000000000000000000000000000000000000", // 无效地址（零地址）
		"0x1234567890123456789012345678901234567890", // 可能无效的地址
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("开始批量加载合约...")

	// 练习 2：批量加载这些合约
	// 提示：遍历地址数组，使用 store.NewStore() 加载每个合约
	// 练习 3：捕获并统计成功/失败数量
	successCount := 0
	failCount := 0

	for _, addrStr := range contractAddresses {
		addr := common.HexToAddress(addrStr)

		// TODO: 检查地址是否有合约代码
		code, err := client.CodeAt(context.Background(), addr, nil)
		if err != nil {
			fmt.Printf("❌ 地址 %s: 检查代码失败 - %v\n", addr.Hex(), err)
			failCount++
			continue
		}

		if len(code) == 0 {
			fmt.Printf("❌ 地址 %s: 没有合约代码\n", addr.Hex())
			failCount++
			continue
		}

		// TODO: 尝试加载合约
		_, err = store.NewStore(addr, client)
		if err != nil {
			fmt.Printf("❌ 地址 %s: 加载失败 - %v\n", addr.Hex(), err)
			failCount++
			continue
		}

		fmt.Printf("✅ 地址 %s: 加载成功 (代码长度: %d 字节)\n", addr.Hex(), len(code))
		successCount++
	}

	// 练习 4：输出统计信息
	fmt.Println("\n=== 统计信息 ===")
	fmt.Printf("总数: %d\n", len(contractAddresses))
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", failCount)
	fmt.Printf("成功率: %.2f%%\n", float64(successCount)/float64(len(contractAddresses))*100)
}
