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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Store 合约的 ABI（简化版，仅包含 SetItem 和 GetItem 函数）
	storeABI = `[{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"getItem","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"}]`
)

func main() {
	// 从环境变量获取配置
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal("请设置环境变量 PRIVATE_KEY")
	}

	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		contractAddressStr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"
	}

	// 连接到以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("✅ 已连接到以太坊节点")

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

	fmt.Println("✅ ABI 解析成功")

	// 练习 2：获取发送地址
	// 提示：从私钥获取公钥，然后转换为地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法转换公钥类型")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("✅ 发送地址: %s\n", fromAddress.Hex())

	// 练习 3：获取 nonce
	// 提示：使用 client.PendingNonceAt()
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Nonce: %d\n", nonce)

	// 练习 4：获取 Gas 价格
	// 提示：使用 client.SuggestGasPrice()
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Gas 价格: %s Wei\n", gasPrice.String())

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

	fmt.Println("✅ 调用数据打包成功")

	// 练习 6：创建交易
	// 提示：使用 types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	contractAddress := common.HexToAddress(contractAddressStr)
	chainID := big.NewInt(11155111) // Sepolia

	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),      // 金额（0 ETH）
		300000,             // Gas 限制
		gasPrice,           // Gas 价格
		input,              // 调用数据
	)

	// 练习 7：签名交易
	// 提示：使用 types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ 交易签名成功")

	// 练习 8：发送交易
	// 提示：使用 client.SendTransaction()
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 交易已发送: %s\n", signedTx.Hash().Hex())

	// 练习 9：等待交易确认
	// 提示：循环调用 client.TransactionReceipt() 直到成功
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ 交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("✅ Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("✅ 交易状态: %d\n", receipt.Status)

	// 练习 10：使用 eth_call 读取数据
	// 提示：使用 parsedABI.Pack("getItem", key) 打包查询数据
	// 然后使用 client.CallContract() 发送调用
	// 最后使用 parsedABI.UnpackIntoInterface() 解析返回值
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

	fmt.Printf("✅ 读取的值: %s\n", string(unpacked[:]))

	// 验证
	if unpacked == value {
		fmt.Println("✅ 验证成功：存储的值与原始值一致")
	} else {
		fmt.Println("❌ 验证失败：存储的值与原始值不一致")
	}
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
