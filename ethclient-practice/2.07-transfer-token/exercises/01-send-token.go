// 01-send-token.go - 发送 ERC20 代币转账练习
//
// 任务：
// 1. 从环境变量读取配置
// 2. 构建 transfer 函数调用数据
// 3. 估算 Gas 并发送交易
// 4. 等待交易确认
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TOKEN_ADDRESS=0x... && export TO_ADDRESS=0x... && export AMOUNT=1000... && go run exercises/01-send-token.go

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {
	fmt.Println("=== ERC20 代币转账 ===")

	// TODO 1: 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")
	amountStr := os.Getenv("AMOUNT")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || toAddressHex == "" || amountStr == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, TO_ADDRESS, AMOUNT")
	}

	// TODO 2: 连接到以太坊节点
	var client *ethclient.Client
	{
		// 在这里填写代码
	}
	defer client.Close()

	// TODO 3: 加载私钥并获取发送方地址
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	{
		// 在这里填写代码
	}

	// TODO 4: 获取 Nonce 和 Gas Price
	var nonce uint64
	var gasPrice *big.Int
	{
		// 在这里填写代码
	}

	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("Nonce: %d\n", nonce)

	// TODO 5: 设置地址
	toAddress := common.HexToAddress(toAddressHex)
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// TODO 6: 构建 transfer 函数调用数据
	// 函数签名: transfer(address,uint256)
	var data []byte
	{
		// 6.1 生成 Method ID
		var methodID []byte
		{
			// 在这里填写代码
			// 提示：使用 sha3.NewLegacyKeccak256() 计算 "transfer(address,uint256)" 的哈希
		}
		fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

		// 6.2 填充地址到 32 字节
		var paddedAddress []byte
		{
			// 在这里填写代码
			// 提示：使用 common.LeftPadBytes()
		}
		fmt.Printf("Padded Address: %s\n", hexutil.Encode(paddedAddress))

		// 6.3 设置金额并填充
		amount := new(big.Int)
		amount.SetString(amountStr, 10)
		var paddedAmount []byte
		{
			// 在这里填写代码
		}
		fmt.Printf("Amount: %s\n", amount.String())
		fmt.Printf("Padded Amount: %s\n", hexutil.Encode(paddedAmount))

		// 6.4 组合数据
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)
	}

	// TODO 7: 估算 Gas
	var gasLimit uint64
	{
		// 在这里填写代码
		// 提示：使用 client.EstimateGas()
		// 注意：To 字段应该是代币合约地址
	}
	fmt.Printf("Estimated Gas: %d\n", gasLimit)

	// TODO 8: 构建交易
	// 注意：value = 0（ERC20 转账不发送 ETH）
	// 注意：to 是代币合约地址，不是接收代币的地址
	value := big.NewInt(0)
	var tx *types.Transaction
	{
		// 在这里填写代码
	}

	// TODO 9: 获取 Chain ID
	var chainID *big.Int
	{
		// 在这里填写代码
	}

	// TODO 10: 签名并发送交易
	var signedTx *types.Transaction
	{
		// 在这里填写代码
	}

	fmt.Printf("\n交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n\n", signedTx.Hash().Hex())

	// TODO 11: 等待交易确认
	fmt.Println("等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Printf("\n✅ 交易成功！\n")
			fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		} else {
			fmt.Printf("\n❌ 交易失败\n")
		}
		break
	}

	fmt.Println("=== 完成 ===")
}
