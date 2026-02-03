package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 字节码从 contract 包引用：import store "github.com/.../contract"，使用 store.StoreContractBytecode

func main() {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}

	client, err := ethclient.Dial(rpcURL)
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
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 练习：获取当前 nonce
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)

	// 练习：计算预期合约地址
	// 提示：使用 crypto.CreateAddress(fromAddress, nonce)
	// predictedAddress := crypto.CreateAddress(fromAddress, nonce)

	fmt.Println("=== 合约地址预测 ===")
	fmt.Printf("发送者地址: %s\n", fromAddress.Hex())
	fmt.Printf("当前 Nonce: %d\n", 0)
	fmt.Printf("预期合约地址: %s\n", "predictedAddress.Hex()")

	// 练习：部署合约
	// gasPrice, _ := client.SuggestGasPrice(context.Background())
	// data, _ := hex.DecodeString(store.StoreContractBytecode)
	// tx := types.NewContractCreation(nonce, big.NewInt(0), 3000000, gasPrice, data)
	// chainID, _ := client.NetworkID(context.Background())
	// signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// client.SendTransaction(context.Background(), signedTx)

	fmt.Printf("\n交易已发送: %s\n", "signedTx.Hash().Hex()")

	// 练习：等待确认并验证地址
	// receipt, _ := waitForReceipt(client, signedTx.Hash())
	// actualAddress := receipt.ContractAddress

	fmt.Println("\n=== 验证结果 ===")
	fmt.Printf("预期地址: %s\n", "predictedAddress.Hex()")
	fmt.Printf("实际地址: %s\n", "actualAddress.Hex()")
	if "predictedAddress" == "actualAddress" {
		fmt.Println("✓ 地址匹配！")
	} else {
		fmt.Println("✗ 地址不匹配！")
	}
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	return nil, nil
}
