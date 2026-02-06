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

const storeABI = util.ReadABI("contract/Store_sol_Store.abi")

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
	parsedABI, err := abi.JSON(strings.NewReader(storeABI))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ABI 解析成功")

	// 练习 2：获取发送地址
	// 提示：从私钥获取公钥，然后转换为地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法转换公钥类型")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("发送地址: %s\n", fromAddress.Hex())

	// 练习 3：获取 nonce
	// 提示：使用 client.PendingNonceAt()
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Nonce: %d\n", nonce)

	// 练习 4：获取 Gas 价格
	// 提示：使用 client.SuggestGasPrice()
	gasPrice, err := client.SuggestGasPrice(context.Background())
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
	input, err := parsedABI.Pack("setItem", key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("调用数据打包成功")

	// 练习 6：创建交易
	// 提示：使用 types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	// 链 ID 可用 client.ChainID(context.Background())
	contractAddress := common.HexToAddress(contractAddressStr)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0), // 金额（0 ETH）
		300000,        // Gas 限制
		gasPrice,      // Gas 价格
		input,         // 调用数据
	)

	// 练习 7：签名交易
	// 提示：使用 types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("交易签名成功")

	// 练习 8：发送交易
	// 提示：使用 client.SendTransaction()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已发送: %s\n", signedTx.Hash().Hex())

	// 练习 9：等待交易确认
	// 提示：循环调用 client.TransactionReceipt() 直到成功
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("交易状态: %d\n", receipt.Status)

	// 练习 10：使用 eth_call 读取数据
	// 提示：parsedABI.Pack("getItem", key) 打包；client.CallContract(ctx, callMsg, nil) 调用；
	// parsedABI.UnpackIntoInterface(&unpacked, "getItem", result) 解析返回值
	// TODO: 打包 getItem 调用、构造 CallMsg、CallContract、Unpack 并打印
	// callInput, err := ???
	// callMsg := ethereum.CallMsg{ To: &contractAddress, Data: callInput }
	// result, err := client.CallContract(???)
	// var unpacked [32]byte
	// err = parsedABI.UnpackIntoInterface(&unpacked, "getItem", result)
	// fmt.Printf("读取的值: %x\n", unpacked)
}

// waitForReceipt 等待交易回执
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(2 * time.Second)
		fmt.Println("⏳ 等待交易确认...")
	}
}
