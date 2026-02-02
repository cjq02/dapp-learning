package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 这里的代码假设已经使用 abigen 生成了 store.go
// 练习前需要先运行：
// cd contract && solcjs --bin Store.sol && solcjs --abi Store.sol
// abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=store.go

func main() {
	// 从环境变量读取配置
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY")
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

	// 加载私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("连接到网络，链 ID: %s\n", chainID.String())

	// 创建交易认证器
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 设置 Gas 参数
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	auth.Value = big.NewInt(0)

	fmt.Printf("Gas 价格: %s Wei\n", gasPrice.String())
	fmt.Printf("Gas 上限: %d\n", auth.GasLimit)

	// 注意：这里需要先编译生成 store.go
	// 在实际运行前，需要在 contract 目录执行：
	// solcjs --bin Store.sol && solcjs --abi Store.sol
	// abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=main --out=store.go

	// 部署合约（这里使用模拟代码，实际需要 abigen 生成的代码）
	// contractAddr, tx, instance, err := DeployStore(auth, client, "1.0")

	fmt.Println("\n开始部署合约...")

	// 模拟部署过程
	// contractAddr, tx, _, err := DeployStore(auth, client, "1.0")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("\n合约部署成功！")
	// fmt.Printf("合约地址: %s\n", contractAddr.Hex())
	// fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())

	// 等待交易确认
	// receipt, err := waitForReceipt(client, tx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if receipt.Status == 1 {
	// 	fmt.Println("✓ 合约已成功部署到链上")
	// 	fmt.Printf("✓ 合约地址: %s\n", receipt.ContractAddress.Hex())
	// 	fmt.Printf("✓ Gas 使用: %d\n", receipt.GasUsed)
	// } else {
	// 	log.Fatal("✗ 合约部署失败")
	// }

	fmt.Println("\n注意：此示例需要使用 abigen 生成 store.go 文件")
	fmt.Println("请先运行以下命令：")
	fmt.Println("  cd contract")
	fmt.Println("  solcjs --bin Store.sol")
	fmt.Println("  solcjs --abi Store.sol")
	fmt.Println("  abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=main --out=store.go")
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	fmt.Println("等待交易确认...")
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err.Error() != "not found" {
			return nil, err
		}
		time.Sleep(2 * time.Second)
	}
}
