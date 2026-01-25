// 01-send-eth.go - 发送 ETH 转账练习
//
// 任务：
// 1. 从环境变量读取私钥
// 2. 连接到 Sepolia 测试网
// 3. 获取发送方地址和 Nonce
// 4. 构建并发送 ETH 转账交易
// 5. 等待交易确认
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TO_ADDRESS=0x... && go run exercises/01-send-eth.go

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ethToWei 将 ETH 数量转为 Wei。amount 为 ETH 数量，如 0.1 表示 0.1 ETH。
// 实现：amount * 1 ETH（1 ETH = 10^18 Wei）
func ethToWei(ethAmount float64) *big.Int {
	oneEth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1 ETH in Wei
	wei := new(big.Float).Mul(big.NewFloat(ethAmount), new(big.Float).SetInt(oneEth))
	result, _ := wei.Int(nil)
	return result
}

func main() {
	fmt.Println("=== ETH 转账 ===")

	// TODO 1: 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}

	toAddressHex := os.Getenv("TO_ADDRESS")
	if toAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 TO_ADDRESS")
	}

	// TODO 2: 连接到以太坊节点
	var client *ethclient.Client
	{
		// 在这里填写代码
		// 提示：使用 ethclient.Dial()
		var err error
		client, err = ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
		if err != nil {
			log.Fatal("错误: 连接到以太坊节点失败", err)
		}
	}
	defer func() {
		// 关闭连接
		if client != nil {
			client.Close()
		}
	}()

	// TODO 3: 加载私钥
	var privateKey *ecdsa.PrivateKey
	var err error
	{
		// 在这里填写代码
		// 提示：使用 crypto.HexToECDSA()
		privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			log.Fatal("错误: 解析私钥失败", err)
		}
	}

	// TODO 4: 获取发送方地址
	var fromAddress common.Address
	{
		// 在这里填写代码
		// 提示：从私钥派生公钥，然后获取地址
		fromAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	fmt.Printf("发送方: %s\n", fromAddress.Hex())

	// TODO 5: 获取 Nonce
	var nonce uint64
	{
		// 在这里填写代码
		// 提示：使用 client.PendingNonceAt()
		nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal("错误: 获取 Nonce 失败", err)
		}
	}

	fmt.Printf("Nonce: %d\n", nonce)

	// TODO 6: 设置转账金额（0.1 ETH）
	value := ethToWei(0.1)
	fmt.Printf("转账金额: %s Wei\n", value.String())

	// TODO 7: 设置 Gas 参数
	var gasLimit uint64 = 21000
	var gasPrice *big.Int
	{
		// 在这里填写代码
		// 提示：使用 client.SuggestGasPrice()
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("错误: 获取 Gas Price 失败", err)
		}
	}

	// TODO 8: 设置接收地址
	toAddress := common.HexToAddress(toAddressHex)

	// TODO 9: 构建未签名交易
	var tx *types.Transaction
	{
		// 在这里填写代码
		// 提示：使用 types.NewTransaction()
		// ETH 转账的 data 参数为 nil
		tx = types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	}

	// TODO 10: 获取 Chain ID
	var chainID *big.Int
	{
		// 在这里填写代码
		// 提示：使用 client.NetworkID()
		chainID, err = client.NetworkID(context.Background())
		if err != nil {
			log.Fatal("错误: 获取 Chain ID 失败", err)
		}
	}

	// TODO 11: 签名交易
	var signedTx *types.Transaction
	{
		// 在这里填写代码
		// 提示：使用 types.SignTx()
		// 需要使用 types.NewEIP155Signer(chainID)
		signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal("错误: 签名交易失败", err)
		}
	}

	// TODO 12: 发送交易
	{
		// 在这里填写代码
		// 提示：使用 client.SendTransaction()
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatal("错误: 发送交易失败", err)
		}
	}

	fmt.Printf("\n交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())

	// TODO 13: 等待交易确认
	fmt.Println("\n等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 简单轮询方式
	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			// 交易还未确认
			time.Sleep(10 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Printf("\n交易成功！\n")
			fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

			// 计算实际费用
			gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
			actualFee := new(big.Int).Mul(gasUsed, gasPrice)
			actualFeeEth := new(big.Float).Quo(
				new(big.Float).SetInt(actualFee),
				big.NewFloat(1e18),
			)
			fmt.Printf("实际费用: %.6f ETH\n", actualFeeEth)
		} else {
			fmt.Printf("\n交易失败！\n")
		}
		break
	}

	fmt.Println("=== 完成 ===")
}
