// mint-token.go - 铸造代币（payable mint 函数）
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
)

func main() {
	fmt.Println("=== 铸造 ERC20 代币 ===")

	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal("错误: 连接到以太坊节点失败", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("错误: 解析私钥失败", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// 检查余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("错误: 获取 ETH 余额失败", err)
	}
	fmt.Printf("您的 ETH 余额: %s wei\n", balance.String())
	balanceETH := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("您的 ETH 余额: %s ETH\n\n", balanceETH.String())

	// 最小 0.001 ETH
	minETH := big.NewInt(1000000000000000) // 0.001 ETH = 10^15 wei
	if balance.Cmp(minETH) < 0 {
		log.Fatal("错误: 您的 ETH 余额不足 0.001 ETH，无法铸造代币")
	}

	// 使用 0.001 ETH 铸造
	ethToSend := minETH
	fmt.Printf("准备发送: %s wei (0.001 ETH)\n", ethToSend.String())
	fmt.Printf("预计获得: 100,000,000 代币 (1 ETH = 100M 代币)\n\n")

	// 获取 Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("错误: 获取 Nonce 失败", err)
	}

	// 获取 Gas Price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("错误: 获取 Gas Price 失败", err)
	}

	// 调用 mint() 函数
	fmt.Println("调用 mint() 函数...")
	mintToken(client, privateKey, fromAddress, tokenAddress, ethToSend, nonce, gasPrice)

	// 等待交易确认
	time.Sleep(3 * time.Second)

	// 查询余额
	fmt.Println("\n=== 检查余额 ===")
	checkBalance(client, tokenAddress, fromAddress)
}

func mintToken(client *ethclient.Client, privateKey *ecdsa.PrivateKey, fromAddress, tokenAddress common.Address, value *big.Int, nonce uint64, gasPrice *big.Int) {
	// 计算方法 ID: mint()
	hash := crypto.Keccak256([]byte("mint()"))
	methodID := hash[:4]

	// 组合数据（mint() 不需要参数）
	data := methodID

	// 估算 Gas
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: &fromAddress,
		To:   &tokenAddress,
		Data: data,
		Value: value,
	})
	if err != nil {
		log.Fatal("错误: 估算 Gas 失败", err)
	}

	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Price: %s wei\n", gasPrice.String())
	fmt.Printf("Gas 费用: %s wei\n", new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice).String())

	// 构建交易
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// 签名
	chainID := big.NewInt(11155111) // Sepolia
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("错误: 签名失败", err)
	}

	// 发送交易
	fmt.Printf("\n发送交易...\n")
	fmt.Printf("数据: %s\n", hexutil.Encode(data))
	fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("错误: 发送交易失败", err)
	}

	fmt.Printf("\n✅ 交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
	fmt.Println("\n等待交易确认...")
}

func checkBalance(client *ethclient.Client, tokenAddress, account common.Address) {
	// 计算方法 ID: balanceOf(address)
	hash := crypto.Keccak256([]byte("balanceOf(address)"))
	methodID := hash[:4]

	// 填充地址
	paddedAddress := common.LeftPadBytes(account.Bytes(), 32)

	// 组合数据
	data := append(methodID, paddedAddress...)

	// 调用合约
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}, nil)

	if err != nil {
		fmt.Printf("❌ 查询余额失败: %v\n", err)
		return
	}

	balance := new(big.Int).SetBytes(result)
	fmt.Printf("当前代币余额: %s wei\n", balance.String())

	// 转换为代币数量
	balanceFloat := new(big.Float).SetInt(balance)
	balanceFloat.Quo(balanceFloat, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
	fmt.Printf("当前代币余额: %s 代币\n", balanceFloat.String())
}
