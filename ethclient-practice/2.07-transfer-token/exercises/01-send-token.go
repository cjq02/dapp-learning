// 01-send-token.go - 发送 ERC20 代币转账练习
//
// 任务：
// 1. 从环境变量读取配置
// 2. 构建 transfer 函数调用数据
// 3. 估算 Gas 并发送交易
// 4. 等待交易确认
//
// 运行：export INFURA_API_KEY=your-key && export PRIVATE_KEY=your-key && export TOKEN_ADDRESS=0x... && export TO_ADDRESS=0x... && export TOKEN_AMOUNT=1000... && go run exercises/01-send-token.go

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// tokenAmountToWei 将人类可读的代币数量转换为最小单位（Wei）
// amount: 代币数量，如 1000 表示 1000 个代币
// decimals: 代币小数位数，如 18 表示 18 位小数（大多数 ERC20 代币）
// 返回: 转换后的最小单位数量
func tokenAmountToWei(amount float64, decimals uint64) *big.Int {
	decimalsBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountFloat := big.NewFloat(amount)
	wei := new(big.Float).Mul(amountFloat, new(big.Float).SetInt(decimalsBig))
	result, _ := wei.Int(nil)
	return result
}

func main() {
	fmt.Println("=== ERC20 代币转账 ===")

	// TODO 1: 从环境变量读取配置
	apiKey := os.Getenv("INFURA_API_KEY")
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	tokenAddressHex := os.Getenv("TOKEN_ADDRESS")
	toAddressHex := os.Getenv("TO_ADDRESS")
	tokenAmount := os.Getenv("TOKEN_AMOUNT")

	if apiKey == "" || privateKeyHex == "" || tokenAddressHex == "" || toAddressHex == "" || tokenAmount == "" {
		log.Fatal("错误: 请设置环境变量 INFURA_API_KEY, PRIVATE_KEY, TOKEN_ADDRESS, TO_ADDRESS, TOKEN_AMOUNT")
	}

	// TODO 2: 连接到以太坊节点
	var client *ethclient.Client
	var err error
	{
		// 在这里填写代码
		client, err = ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
		if err != nil {
			log.Fatal("错误: 连接到以太坊节点失败", err)
		}
	}
	defer client.Close()

	// TODO 3: 加载私钥并获取发送方地址
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	{
		// 在这里填写代码
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			log.Fatal("错误: 解析私钥失败", err)
		}
		fromAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	// TODO 4: 获取 Nonce 和 Gas Price
	var nonce uint64
	var gasPrice *big.Int
	{
		// 在这里填写代码
		nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal("错误: 获取 Nonce 失败", err)
		}
		// 获取建议 Gas Price，但设置最低值为 10 Gwei
		suggestedGasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("错误: 获取 Gas Price 失败", err)
		}
		// 确保最低 Gas Price 为 10 Gwei (Sepolia 测试网的推荐值)
		minGasPrice := big.NewInt(10000000000) // 10 Gwei = 10^10 wei
		if suggestedGasPrice.Cmp(minGasPrice) < 0 {
			gasPrice = minGasPrice
		} else {
			gasPrice = suggestedGasPrice
		}
	}

	fmt.Printf("发送方: %s\n", fromAddress.Hex())
	fmt.Printf("Nonce: %d\n", nonce)
	fmt.Printf("Gas Price: %s wei (%.2f Gwei)\n", gasPrice.String(), new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9)))

	// TODO 5: 设置地址
	toAddress := common.HexToAddress(toAddressHex)
	tokenAddress := common.HexToAddress(tokenAddressHex)

	// TODO 5.1: 查询代币小数位数（可选，默认使用 18）
	// 大多数 ERC20 代币使用 18 位小数，如果查询失败则使用默认值
	var decimals uint64 = 18 // 默认值
	{
		// 可选：查询代币的 decimals()
		// 如果查询失败，使用默认值 18
		// 提示：构建 "decimals()" 函数调用数据，使用 client.CallContract()
		hash := crypto.Keccak256([]byte("decimals()"))
		methodID := hash[:4]
		result, err := client.CallContract(context.Background(), ethereum.CallMsg{
			To:   &tokenAddress,
			Data: methodID,
		}, nil)
		decimals = new(big.Int).SetBytes(result).Uint64()
		if err != nil {
			log.Fatal("错误: 查询代币小数位数失败", err)
		}
	}

	// TODO 5.2: 将人类可读的数量转换为最小单位
	var amount *big.Int
	{
		// 在这里填写代码
		// 提示：将 tokenAmount 解析为 float64，然后使用 tokenAmountToWei() 转换
		// 例如：tokenAmount = "1000" 表示 1000 个代币
		// 如果代币有 18 位小数，则转换为 1000 * 10^18
		tokenAmountFloat, err := strconv.ParseFloat(tokenAmount, 64)
		if err != nil {
			log.Fatal("错误: 转换代币数量失败", err)
		}
		amount = tokenAmountToWei(tokenAmountFloat, decimals)
		if err != nil {
			log.Fatal("错误: 转换代币数量失败", err)
		}
	}
	fmt.Printf("转账数量: %s 代币 (decimals: %d)\n", tokenAmount, decimals)
	fmt.Printf("转换为最小单位: %s\n", amount.String())

	// TODO 6: 构建 transfer 函数调用数据
	// 函数签名: transfer(address,uint256)
	var data []byte
	{
		// 6.1 生成 Method ID
		var methodID []byte
		{
			// 在这里填写代码
			// 提示：使用 sha3.NewLegacyKeccak256() 计算 "transfer(address,uint256)" 的哈希
			hash := sha3.NewLegacyKeccak256()
			hash.Write([]byte("transfer(address,uint256)"))
			methodID = hash.Sum(nil)[:4]
		}
		fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

		// 6.2 填充地址到 32 字节
		var paddedAddress []byte
		{
			// 在这里填写代码
			// 提示：使用 common.LeftPadBytes()
			paddedAddress = common.LeftPadBytes(toAddress.Bytes(), 32)
			if err != nil {
				log.Fatal("错误: 填充地址失败", err)

			}
		}
		fmt.Printf("Padded Address: %s\n", hexutil.Encode(paddedAddress))

		// 6.3 填充金额到 32 字节
		var paddedAmount []byte
		{
			// 在这里填写代码
			// 提示：使用 common.LeftPadBytes(amount.Bytes(), 32)
			paddedAmount = common.LeftPadBytes(amount.Bytes(), 32)
			if err != nil {
				log.Fatal("错误: 填充金额失败", err)
			}
		}
		fmt.Printf("Padded Amount: %s\n", hexutil.Encode(paddedAmount))

		// 6.4 组合数据
		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)
	}

	// TODO 7: 估算 Gas
	var gasLimit uint64
	{
		// 在这里填写代码
		// 提示：使用 client.EstimateGas()
		// 注意：To 字段应该是代币合约地址
		// 注意：需要指定 From 字段，某些合约依赖发送者地址
		gasLimit, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			From: fromAddress,
			To:   &tokenAddress,
			Data: data,
		})
		if err != nil {
			log.Fatal("错误: 估算 Gas 失败", err)
		}
	}
	fmt.Printf("Estimated Gas: %d\n", gasLimit)

	// TODO 8: 构建交易
	// 注意：value = 0（ERC20 转账不发送 ETH）
	// 注意：to 是代币合约地址，不是接收代币的地址
	value := big.NewInt(0)
	var tx *types.Transaction
	{
		// 在这里填写代码
		// ERC20 转账：交易的 to 是代币合约地址，接收代币的地址在 data 中编码
		tx = types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
		if err != nil {
			log.Fatal("错误: 构建交易失败", err)
		}
	}

	// TODO 9: 获取 Chain ID
	var chainID *big.Int
	{
		// 在这里填写代码
		chainID, err = client.NetworkID(context.Background())
		if err != nil {
			log.Fatal("错误: 获取 Chain ID 失败", err)
		}
	}

	// TODO 10: 签名并发送交易
	var signedTx *types.Transaction
	{
		// 在这里填写代码
		signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal("错误: 签名交易失败", err)
		}
	}

	fmt.Printf("\n交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("查看: https://sepolia.etherscan.io/tx/%s\n\n", signedTx.Hash().Hex())

	// TODO 11: 等待交易确认
	fmt.Println("等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		if receipt.Status == 1 {
			fmt.Printf("\n✅ 交易成功！\n")
			fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		} else {
			fmt.Printf("\n❌ 交易失败\n")
		}
		break
	}

	fmt.Println("=== 完成 ===")
}
