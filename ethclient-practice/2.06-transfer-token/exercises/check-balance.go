// check-balance.go - 查询 ERC20 代币余额
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	apiKey := os.Getenv("INFURA_API_KEY")
	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	toAddressHex := os.Getenv("TO_ADDRESS") // 可选：要查询的地址

	if apiKey == "" || tokenAddress == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, TOKEN_ADDRESS")
	}

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal("错误: 连接失败", err)
	}
	defer client.Close()

	tokenContract := common.HexToAddress(tokenAddress)

	// 确定要查询的地址
	var userAddress common.Address
	if toAddressHex != "" {
		// 如果指定了查询地址，使用该地址
		userAddress = common.HexToAddress(toAddressHex)
	} else if privateKeyHex != "" {
		// 否则使用私钥对应的地址
		privateKey, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			log.Fatal("错误: 解析私钥失败", err)
		}
		userAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
	} else {
		log.Fatal("错误: 请设置环境变量 PRIVATE_KEY 或 TO_ADDRESS")
	}

	// 调用 ERC20 合约的 balanceOf(address) 函数
	// 1. 计算 balanceOf(address) 的方法 ID
	hash := crypto.Keccak256([]byte("balanceOf(address)"))
	methodID := hash[:4]

	// 2. 将地址填充为 32 字节
	paddedAddress := common.LeftPadBytes(userAddress.Bytes(), 32)

	// 3. 组合数据
	data := append(methodID, paddedAddress...)

	// 4. 调用合约
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContract,
		Data: data,
	}, nil)
	if err != nil {
		log.Fatal("错误: 查询代币余额失败", err)
	}

	// 解析结果
	balance := new(big.Int).SetBytes(result)

	// 查询代币的 decimals
	decimalsHash := crypto.Keccak256([]byte("decimals()"))
	decimalsMethodID := decimalsHash[:4]
	decimalsResult, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContract,
		Data: decimalsMethodID,
	}, nil)
	if err != nil {
		log.Fatal("错误: 查询代币 decimals 失败", err)
	}
	if len(decimalsResult) == 0 {
		log.Fatal("错误: 合约未返回 decimals 值")
	}
	decimals := new(big.Int).SetBytes(decimalsResult).Uint64()

	// 查询合约地址本身的代币余额
	paddedContractAddress := common.LeftPadBytes(tokenContract.Bytes(), 32)
	contractData := append(methodID, paddedContractAddress...)
	contractResult, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContract,
		Data: contractData,
	}, nil)
	if err != nil {
		log.Fatal("错误: 查询合约地址代币余额失败", err)
	}
	contractBalance := new(big.Int).SetBytes(contractResult)

	// 转换为代币数量的辅助函数
	decimalsBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	formatBalance := func(b *big.Int) string {
		balanceFloat := new(big.Float).SetInt(b)
		balanceFloat.Quo(balanceFloat, new(big.Float).SetInt(decimalsBig))
		return balanceFloat.Text('f', 0)
	}

	fmt.Printf("代币地址: %s\n", tokenAddress)
	fmt.Printf("\n=== 接收地址余额 ===\n")
	fmt.Printf("查询地址: %s\n", userAddress.Hex())
	fmt.Printf("代币余额 (最小单位): %s\n", balance.String())
	fmt.Printf("代币余额: %s 代币 (decimals: %d)\n", formatBalance(balance), decimals)
	
	fmt.Printf("\n=== 合约地址余额 ===\n")
	fmt.Printf("合约地址: %s\n", tokenContract.Hex())
	fmt.Printf("代币余额 (最小单位): %s\n", contractBalance.String())
	fmt.Printf("代币余额: %s 代币 (decimals: %d)\n", formatBalance(contractBalance), decimals)
}
