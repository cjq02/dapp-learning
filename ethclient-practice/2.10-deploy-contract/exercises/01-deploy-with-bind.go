package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dapp-learning/ethclient-practice/util"

	store "github.com/dapp-learning/ethclient/deploy-contract/contract"
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
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 练习：加载私钥
	// 提示：使用 crypto.HexToECDSA，注意去掉 0x 前缀
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// 练习：获取链 ID
	// 提示：使用 client.NetworkID
	// var chainID *big.Int
	// chainID, err = ???
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 练习：创建交易认证器
	// 提示：使用 bind.NewKeyedTransactorWithChainID
	// var auth *bind.TransactOpts
	// auth, err = ???
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 练习：设置 Gas 参数
	// 提示：
	// gasPrice, _ := client.SuggestGasPrice(context.Background())
	// auth.GasLimit = uint64(300000)
	// auth.GasPrice = gasPrice
	// auth.Value = big.NewInt(0)
	// 使用 util.SuggestGasPrice：获取建议 gasPrice，并设置最低值（避免过低导致长时间 pending）
	minGasPrice := new(big.Int).Mul(big.NewInt(10), big.NewInt(1e9)) // 10 gwei
	gasPrice, err := util.SuggestGasPrice(context.Background(), client, minGasPrice)
	if err != nil {
		log.Fatal(err)
	}
	auth.GasPrice = gasPrice
	// 部署合约时用较大的 GasLimit 上限，避免 out of gas。实际只扣 gasPrice×gasUsed，不会多付。
	auth.GasLimit = 2_000_000
	auth.Value = big.NewInt(0)

	// 练习：部署合约
	// 提示：使用 DeployStore（需要先编译生成 store.go）
	// contractAddr, tx, instance, err := DeployStore(auth, client, "1.0")
	contractAddr, tx, instance, err := store.DeployStore(auth, client, "1.0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("合约地址: %s\n", contractAddr.Hex())
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("实例: %v\n", instance)

	// 练习：等待交易确认（与 02-deploy-raw 一致，用 util.WaitForReceipt 轮询）
	fmt.Printf("等待交易回执: %s\n", tx.Hash().Hex())
	receipt, err := util.WaitForReceipt(context.Background(), client, tx.Hash(), 2*time.Second, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已获取交易回执")
	fmt.Printf("状态 Status: %d (1=成功,0=失败)\n", receipt.Status)
	fmt.Printf("区块号 BlockNumber: %s\n", receipt.BlockNumber.String())
	fmt.Printf("GasUsed: %d\n", receipt.GasUsed)
	fmt.Printf("合约地址 ContractAddress: %s\n", receipt.ContractAddress.Hex())
	if receipt.Status == 1 {
		fmt.Println("合约已部署到链上")
	} else {
		fmt.Println("合约部署失败（可到 Etherscan 查看该交易的错误详情）")
	}
}
