// 03-transfer-from.go - ERC20 代理转账练习
//
// 任务：
// 1. 调用 transferFrom(from, to, amount) 代理转账
// 2. 等待转账交易确认
//
// 运行：export INFURA_API_KEY=your-key && export SPENDER_PRIVATE_KEY=your-key && export TOKEN_ADDRESS=0x... && export FROM_ADDRESS=0x... && export TO_ADDRESS=0x... && go run exercises/03-transfer-from.go
//
// 说明：
// - SPENDER_PRIVATE_KEY: 被授权者（spender）的私钥
// - FROM_ADDRESS: 代币持有者的地址
// - TO_ADDRESS: 接收代币的地址
// - TOKEN_AMOUNT: 从控制台输入转账数量（人类可读单位，如 100 表示 100 个代币）
//
// 前提条件：
// - FROM_ADDRESS 必须已经通过 approve 授权给 SPENDER（即 PRIVATE_KEY 对应的地址）

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
	fmt.Println("=== ERC20 代理转账 ===")

	// 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	spenderPrivateKeyHex := os.Getenv("SPENDER_PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	fromAddressHex := os.Getenv("FROM_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")

	if apiKey == "" || spenderPrivateKeyHex == "" || tokenAddressHex == "" || fromAddressHex == "" || toAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, SPENDER_PRIVATE_KEY, TOKEN_ADDRESS, FROM_ADDRESS, TO_ADDRESS")
	}

	// 从控制台输入转账数量
	var amountStr string
	fmt.Print("请输入转账数量（人类可读单位，如 100 表示 100 个代币）: ")
	fmt.Scanln(&amountStr)

	if amountStr == "" {
		log.Fatal("错误: 转账数量不能为空")
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 加载私钥并获取地址（这是被授权者/spender的地址）
	privateKey, err := crypto.HexToECDSA(spenderPrivateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	spenderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tokenAddress := common.HexToAddress(tokenAddressHex)
	fromAddress := common.HexToAddress(fromAddressHex)
	toAddress := common.HexToAddress(toAddressHex)

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

	fmt.Printf("\n=== 代理转账信息 ===\n")
	fmt.Printf("代理（被授权者）: %s\n", spenderAddress.Hex())
	fmt.Printf("代币持有者: %s\n", fromAddress.Hex())
	fmt.Printf("接收地址: %s\n", toAddress.Hex())
	fmt.Printf("代币合约: %s\n", tokenAddress.Hex())
	fmt.Printf("转账数量: %s (decimals: %d)\n", amountStr, decimals)
	fmt.Printf("转换为最小单位: %s\n", amount.String())

	// 预检查：授权额度与持有者余额，避免 "execution reverted" 时难以排查
	{
		allowanceData := util.BuildCallData("allowance(address,address)", fromAddress.Bytes(), spenderAddress.Bytes())
		allowanceResult, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: allowanceData,
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询授权额度失败", err)
		}
		allowance := new(big.Int).SetBytes(allowanceResult)
		if allowance.Cmp(amount) < 0 {
			log.Fatalf("错误: 授权额度不足。当前授权: %s，需要: %s。请让代币持有者 %s 先用 03-approve 对该 spender 授权足够数量后再试。",
				allowance.String(), amount.String(), fromAddress.Hex())
		}
		balanceData := util.BuildCallData("balanceOf(address)", fromAddress.Bytes())
		balanceResult, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: balanceData,
		}, nil)
		if err != nil {
			log.Fatal("错误: 查询代币持有者余额失败", err)
		}
		balance := new(big.Int).SetBytes(balanceResult)
		if balance.Cmp(amount) < 0 {
			log.Fatalf("错误: 代币持有者余额不足。当前余额: %s，需要: %s。", balance.String(), amount.String())
		}
	}

	// 构建 transferFrom 数据
	// transferFrom(address from, address to, uint256 amount)
	transferFromData := util.BuildCallData("transferFrom(address,address,uint256)",
		fromAddress.Bytes(),
		toAddress.Bytes(),
		amount.Bytes(),
	)

	// 获取 nonce（使用 spender 的 nonce）
	nonce, err := client.PendingNonceAt(context.Background(), spenderAddress)
	if err != nil {
		log.Fatal("错误: 获取 Nonce 失败", err)
	}

	// 获取 gas price
	gasPrice, err := util.SuggestGasPrice(context.Background(), client, big.NewInt(10000000000)) // 最低 10 Gwei
	if err != nil {
		log.Fatal("错误: 获取 Gas Price 失败", err)
	}

	// 估算 gas
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: spenderAddress,
		To:   &tokenAddress,
		Data: transferFromData,
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
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, transferFromData)

	// 签名（使用 spender 的私钥）
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

	fmt.Printf("\n✅ 代理转账交易已发送\n")
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
			fmt.Printf("\n✅ 代理转账完成！\n")
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
			fmt.Printf("\n成功从 %s 转账 %s 代币到 %s\n", fromAddress.Hex(), amountStr, toAddress.Hex())
			break
		} else {
			log.Fatal("\n❌ 转账交易失败")
		}
	}

	fmt.Println("\n=== 完成 ===")
}
