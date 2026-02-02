// 03-approve.go - ERC20 授权练习
//
// 任务：
// 1. 调用 approve(spender, amount) 授权指定地址使用你的代币
// 2. 等待授权交易确认
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TOKEN_ADDRESS=0x... && export SPENDER_ADDRESS=0x... && go run exercises/03-approve.go
//
// 说明：
// - PRIVATE_KEY: 代币持有者的私钥（你的私钥）
// - SPENDER_ADDRESS: 被授权的地址（你允许谁使用你的代币）
// - TOKEN_AMOUNT: 从控制台输入授权数量（人类可读单位，如 100 表示 100 个代币）

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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/dapp-learning/ethclient/util"
)

func main() {
	fmt.Println("=== ERC20 授权 ===")

	// 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	spenderAddressHex := os.Getenv("SPENDER_ADDRESS")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || spenderAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, SPENDER_ADDRESS")
	}

	// 从控制台输入授权数量
	var amountStr string
	fmt.Print("请输入授权数量（人类可读单位，如 100 表示 100 个代币）: ")
	fmt.Scanln(&amountStr)

	if amountStr == "" {
		log.Fatal("错误: 授权数量不能为空")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 加载私钥并获取地址
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tokenAddress := common.HexToAddress(tokenAddressHex)
	spenderAddress := common.HexToAddress(spenderAddressHex)

	// 查询代币小数位数（默认 18）
	var decimals uint64 = 18
	{
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: util.BuildCallData("decimals()"),
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询代币小数位数失败", err)
		}
		decimals = new(big.Int).SetBytes(result).Uint64()
	}

	// 将人类可读的数量转换为最小单位
	var amountFloat float64
	_, err = fmt.Sscanf(amountStr, "%f", &amountFloat)
	if err != nil {
		log.Fatalf("错误: 无法解析代币数量 %s: %v", amountStr, err)
	}
	amount := util.TokenAmountToWei(amountFloat, decimals)

	fmt.Printf("\n=== 授权信息 ===\n")
	fmt.Printf("代币持有者: %s\n", fromAddress.Hex())
	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("被授权地址: %s\n", spenderAddress.Hex())
	fmt.Printf("授权数量: %s (decimals: %d)\n", amountStr, decimals)
	fmt.Printf("转换为最小单位: %s\n", amount.String())

	// 构建 approve 数据
	approveData := util.BuildCallData("approve(address,uint256)", spenderAddress.Bytes(), amount.Bytes())

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("错误: 获取 Nonce 失败", err)
	}

	// 获取 gas price
	gasPrice, err := util.SuggestGasPrice(context.Background(), client, big.NewInt(10000000000)) // 最低 10 Gwei
	if err != nil {
		log.Fatal("错误: 获取 Gas Price 失败", err)
	}

	// 估算 gas（必须指定 From，否则模拟时 msg.sender 为零地址会导致 approve revert）
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: approveData,
	})
	if err != nil {
		log.Fatal("错误: 估算 Gas 失败", err)
	}

	fmt.Printf("\n=== 交易信息 ===\n")
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("Gas Price: %s wei\n", gasPrice.String())
	fmt.Printf("Gas Limit: %d\n", gasLimit)

	// 构建交易
	value := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, approveData)

	// 签名
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("错误: 获取 Chain ID 失败", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("错误: 签名交易失败", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("错误: 发送交易失败", err)
	}

	fmt.Printf("\n✅ 授权交易已发送\n")
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Printf("\n等待交易确认...\n")

	// 等待交易确认
	for {
		receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		if receipt.Status == 1 {
			fmt.Printf("\n✅ 授权已确认！\n")
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
			fmt.Printf("\n现在 %s 被授权可以使用 %s 的代币\n", spenderAddress.Hex(), fromAddress.Hex())
			break
		} else {
			log.Fatal("\n❌ 授权交易失败")
		}
	}

	fmt.Println("\n=== 完成 ===")
}
