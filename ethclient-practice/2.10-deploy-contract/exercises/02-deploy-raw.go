package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dapp-learning/ethclient-practice/util"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 纯 ethclient 部署：用 util.ReadBin / ReadABI 从 contract 目录读 .bin 和 .abi，不依赖 store.go。

func main() {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}

	// 练习：连接到以太坊节点
	// var client *ethclient.Client
	// client, err = ???
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("已连接 RPC:", rpcURL)

	// 练习：加载私钥并获取发送者地址
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
	fmt.Printf("发送者地址: %s\n", fromAddress.Hex())

	// 练习：获取 nonce
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("当前 Nonce: %d\n", nonce)

	// 练习：获取 Gas 价格
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	minGasPrice := new(big.Int).Mul(big.NewInt(10), big.NewInt(1e9)) // 10 gwei
	gasPrice, err := util.SuggestGasPrice(context.Background(), client, minGasPrice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GasPrice: %s wei\n", gasPrice.String())

	// 从 contract 目录读取 .bin / .abi（支持从 2.10-deploy-contract 或 exercises 运行）
	binPath, abiPath := "contract/Store_sol_Store.bin", "contract/Store_sol_Store.abi"
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		binPath, abiPath = "../contract/Store_sol_Store.bin", "../contract/Store_sol_Store.abi"
	}
	bytecode, err := util.ReadBin(binPath)
	if err != nil {
		log.Fatal(err)
	}
	parsedABI, err := util.ReadABI(abiPath)
	if err != nil {
		log.Fatal(err)
	}
	constructorArgs, err := parsedABI.Constructor.Inputs.Pack("1.0")
	if err != nil {
		log.Fatal("打包构造参数失败:", err)
	}
	data := append(bytecode, constructorArgs...)

	// 练习：创建合约部署交易
	// 提示：使用 types.NewContractCreation，to 地址为 nil
	// tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)
	gasLimit := uint64(3000000)
	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data)
	fmt.Printf("GasLimit: %d\n", gasLimit)
	fmt.Printf("字节码长度: %d bytes\n", len(data))

	// 练习：获取链 ID 并签名交易
	// chainID, _ := client.NetworkID(context.Background())
	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("链 ID: %s，交易已签名\n", chainID.String())

	// 练习：发送交易
	// err = client.SendTransaction(context.Background(), signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	txHash := signedTx.Hash()
	fmt.Printf("交易已发送，交易哈希: %s\n", txHash.Hex())

	// 练习：等待交易确认并获取合约地址
	// receipt, err := waitForReceipt(client, signedTx.Hash())
	fmt.Printf("等待交易回执: %s\n", txHash.Hex())
	receipt, err := util.WaitForReceipt(context.Background(), client, txHash, 2*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("已获取交易回执")
	fmt.Printf("状态 Status: %d (1=成功,0=失败)\n", receipt.Status)
	fmt.Printf("区块号 BlockNumber: %s\n", receipt.BlockNumber.String())
	fmt.Printf("GasUsed: %d\n", receipt.GasUsed)
	fmt.Printf("合约地址 ContractAddress: %s\n", receipt.ContractAddress.Hex())
	fmt.Printf("交易索引 TxIndex: %d\n", receipt.TransactionIndex)
	if receipt.Status == 1 {
		fmt.Println("合约部署成功")
	} else {
		fmt.Println("合约部署失败（可到 Etherscan 查看该交易的错误详情）")
	}
}
