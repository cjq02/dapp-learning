// 03-approve-transfer.go - 授权并代理转账 - 答案

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
	fmt.Println("=== ERC20 授权并代理转账 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	spenderAddressHex := os.Getenv("SPENDER_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")
	amountStr := os.Getenv("AMOUNT")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || spenderAddressHex == "" || toAddressHex == "" || amountStr == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, SPENDER_ADDRESS, TO_ADDRESS, AMOUNT")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/"+apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

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
	toAddress := common.HexToAddress(toAddressHex)

	amount := new(big.Int)
	amount.SetString(amountStr, 10)

	// 步骤 1: 构建 approve 函数调用数据
	approveSignature := []byte("approve(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(approveSignature)
	methodID := hash.Sum(nil)[:4]

	paddedSpender := common.LeftPadBytes(spenderAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var approveData []byte
	approveData = append(approveData, methodID...)
	approveData = append(approveData, paddedSpender...)
	approveData = append(approveData, paddedAmount...)

	fmt.Printf("授权地址: %s\n", spenderAddress.Hex())
	fmt.Printf("授权金额: %s\n", amount.String())

	// 步骤 2: 发送 approve 交易
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, _ := client.SuggestGasPrice(context.Background())

	value := big.NewInt(0)
	gasLimit, _ := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: approveData,
	})

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, approveData)

	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n授权交易已发送: %s\n", signedTx.Hash().Hex())

	// 步骤 3: 等待授权交易确认
	fmt.Println("等待授权交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Println("\n✅ 授权已确认！")
		} else {
			log.Fatal("\n❌ 授权交易失败")
		}
		break
	}

	// 步骤 4: 构建 transferFrom 函数调用数据
	transferFromSignature := []byte("transferFrom(address,address,uint256)")
	hash = sha3.NewLegacyKeccak256()
	hash.Write(transferFromSignature)
	methodID = hash.Sum(nil)[:4]

	paddedFrom := common.LeftPadBytes(fromAddress.Bytes(), 32)
	paddedTo := common.LeftPadBytes(toAddress.Bytes(), 32)

	var transferFromData []byte
	transferFromData = append(transferFromData, methodID...)
	transferFromData = append(transferFromData, paddedFrom...)
	transferFromData = append(transferFromData, paddedTo...)
	transferFromData = append(transferFromData, paddedAmount...)

	// 步骤 5: 发送 transferFrom 交易
	nonce, _ = client.PendingNonceAt(context.Background(), fromAddress)

	gasLimit, _ = client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: transferFromData,
	})

	tx = types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, transferFromData)
	signedTx, _ = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n代理转账交易已发送: %s\n", signedTx.Hash().Hex())

	// 步骤 6: 等待转账交易确认
	fmt.Println("等待转账交易确认...")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel2()

	for {
		receipt, err := client.TransactionReceipt(ctx2, signedTx.Hash())
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Println("\n✅ 代理转账完成！")
			fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		} else {
			fmt.Println("\n❌ 转账交易失败")
		}
		break
	}

	fmt.Println("=== 完成 ===")
}
