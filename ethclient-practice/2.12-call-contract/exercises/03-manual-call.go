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

	"ethclient/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	// 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	fmt.Println("已连接到以太坊节点")

	// 从私钥创建私钥实例
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	// 练习 1：解析 ABI 字符串
	// 提示：使用 abi.JSON(strings.NewReader(storeABI))
	var parsedABI abi.ABI
	// TODO: parsedABI, err = abi.JSON(strings.NewReader(storeABI))
	parsedABI, err = util.ReadABI("../store/Store_sol_Store.abi")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ABI 解析成功")

	// 练习 2：获取发送地址
	// 提示：从私钥获取公钥，然后转换为地址（privateKey.Public() -> 类型断言 *ecdsa.PublicKey -> crypto.PubkeyToAddress）
	var fromAddress common.Address
	// TODO: 在此实现并赋值 fromAddress（需 publicKey、类型断言为 *ecdsa.PublicKey、crypto.PubkeyToAddress）
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("发送地址: %s\n", fromAddress.Hex())

	// 练习 3：获取 nonce
	// 提示：使用 client.PendingNonceAt(context.Background(), fromAddress)
	var nonce uint64
	// TODO: nonce, err = client.PendingNonceAt(...)
	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Nonce: %d\n", nonce)

	// 练习 4：获取 Gas 价格
	// 提示：使用 client.SuggestGasPrice(context.Background())
	var gasPrice *big.Int
	// TODO: gasPrice, err = client.SuggestGasPrice(...)
	gasPrice, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Gas 价格: %s Wei\n", gasPrice.String())

	// 准备数据
	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("manual_key"))
	copy(value[:], []byte("manual_value"))

	// 练习 5：使用 ABI 打包函数调用数据
	// 提示：使用 parsedABI.Pack("setItem", key, value)
	var input []byte
	// TODO: input, err = parsedABI.Pack("setItem", key, value)
	input, err = parsedABI.Pack("setItem", key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("调用数据打包成功")

	// 练习 6：创建交易
	// 提示：先获取 chainID（client.ChainID）和 contractAddress，再 types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	contractAddress := common.HexToAddress(contractAddressStr)
	var chainID *big.Int
	var tx *types.Transaction
	// TODO: 获取 chainID（client.ChainID），检查 err；再 tx = types.NewTransaction(nonce, contractAddress, big.NewInt(0), 300000, gasPrice, input)
	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx = types.NewTransaction(nonce, contractAddress, big.NewInt(0), 300000, gasPrice, input)
	_ = chainID
	_ = tx

	// 练习 7：签名交易
	// 提示：使用 types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	var signedTx *types.Transaction
	// TODO: signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("交易签名成功")

	// 练习 8：发送交易
	// 提示：使用 client.SendTransaction(context.Background(), signedTx)
	// TODO: err = client.SendTransaction(context.Background(), signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易已发送: %s\n", signedTx.Hash().Hex())

	// 练习 9：等待交易确认
	// 提示：可调用 waitForReceipt(client, signedTx.Hash())，或自己写循环 client.TransactionReceipt 直到成功
	var receipt *types.Receipt
	// TODO: receipt, err = waitForReceipt(client, signedTx.Hash())
	receipt, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("交易状态: %d\n", receipt.Status)

	// 练习 10：使用 eth_call 读取刚写入的数据
	// 提示：parsedABI.Pack("getItem", key) 打包；ethereum.CallMsg{To: &contractAddress, Data: callInput} 构造消息；
	// client.CallContract(ctx, callMsg, nil) 调用；parsedABI.UnpackIntoInterface(&unpacked, "getItem", result) 解析返回值
	// TODO: 打包 getItem、CallContract、Unpack，并打印读取到的值（如 strings.TrimRight(string(unpacked[:]), "\x00")）
	callInput, err := parsedABI.Pack("getItem", key)
	if err != nil {
		log.Fatal(err)
	}
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callInput,
	}
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}
	var unpacked [32]byte
	err = parsedABI.UnpackIntoInterface(&unpacked, "getItem", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("读取的值: %s\n", strings.TrimRight(string(unpacked[:]), "\x00"))
	if unpacked == value {
		fmt.Println("验证成功")
	} else {
		fmt.Println("验证失败")
	}
}

// waitForReceipt 等待交易回执（练习 9 可直接调用）
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		time.Sleep(2 * time.Second)
		fmt.Println("⏳ 等待交易确认...")
	}
}
