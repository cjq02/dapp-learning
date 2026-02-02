package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	// 练习：连接到以太坊节点
	// 提示：使用 ethclient.Dial
	// var client *ethclient.Client
	// client, err = ???

	// 练习：加载私钥
	// 提示：使用 crypto.HexToECDSA，注意去掉 0x 前缀
	// var privateKey *ecdsa.PrivateKey
	// privateKey, err = ???

	// 练习：获取链 ID
	// 提示：使用 client.NetworkID
	// var chainID *big.Int
	// chainID, err = ???

	// 练习：创建交易认证器
	// 提示：使用 bind.NewKeyedTransactorWithChainID
	// var auth *bind.TransactOpts
	// auth, err = ???

	// 练习：设置 Gas 参数
	// 提示：
	// gasPrice, _ := client.SuggestGasPrice(context.Background())
	// auth.GasLimit = uint64(300000)
	// auth.GasPrice = gasPrice
	// auth.Value = big.NewInt(0)

	// 练习：部署合约
	// 提示：使用 DeployStore（需要先编译生成 store.go）
	// contractAddr, tx, instance, err := DeployStore(auth, client, "1.0")

	fmt.Println("合约部署成功！")
	fmt.Printf("合约地址: %s\n", "contractAddr.Hex()")
	fmt.Printf("交易哈希: %s\n", "tx.Hash().Hex()")

	// 练习：等待交易确认
	// 提示：使用 client.TransactionReceipt 轮询查询
	// receipt, err := waitForReceipt(client, tx.Hash())
	// if receipt.Status == 1 {
	//     fmt.Println("合约已部署到链上")
	// }
}
