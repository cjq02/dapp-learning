// Package util provides common utilities for ERC20 token operations and contract deployment.
package util

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// TokenAmountToWei 将人类可读的代币数量转换为最小单位
func TokenAmountToWei(amount float64, decimals uint64) *big.Int {
	decimalsBig := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountFloat := big.NewFloat(amount)
	wei := new(big.Float).Mul(amountFloat, new(big.Float).SetInt(decimalsBig))
	result, _ := wei.Int(nil)
	return result
}

// BuildCallData 构建以太坊合约调用数据
// signature: 函数签名，如 "balanceOf(address)" 或 "transfer(address,uint256)"
// args: 函数参数的字节形式
func BuildCallData(signature string, args ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(signature))
	methodID := hash.Sum(nil)[:4]

	var data []byte
	data = append(data, methodID...)

	for _, arg := range args {
		padded := common.LeftPadBytes(arg, 32)
		data = append(data, padded...)
	}

	return data
}

// WeiToTokenAmount 将最小单位转换为人了可读的代币数量
func WeiToTokenAmount(balance *big.Int, decimals uint64) *big.Float {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(divisor))
}

// SuggestGasPrice 获取建议的 Gas Price，并设置最低值
// minGasPrice: 最低 Gas Price（单位：wei），例如 10000000000 表示 10 Gwei
func SuggestGasPrice(ctx context.Context, client *ethclient.Client, minGasPrice *big.Int) (*big.Int, error) {
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	// 确保不低于最低值
	if gasPrice.Cmp(minGasPrice) < 0 {
		gasPrice = minGasPrice
	}
	return gasPrice, nil
}

// WaitForReceipt 轮询等待交易回执，直到超时或拿到回执。
// interval: 轮询间隔；timeout: 总超时时间。
func WaitForReceipt(
	ctx context.Context,
	client *ethclient.Client,
	txHash common.Hash,
	interval time.Duration,
	timeout time.Duration,
) (*types.Receipt, error) {
	deadline := time.Now().Add(timeout)
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("等待回执超时（%s），最后错误: %w", timeout, err)
		}
		time.Sleep(interval)
	}
}

// formatTokenBalance 格式化代币余额
func FormatTokenBalance(balance *big.Int, decimals uint8) *big.Float {
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	return new(big.Float).Quo(fbal, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))
}

// ReadBin 从 solc 生成的 .bin 文件读取部署字节码。
// 文件内容为十六进制字符串，可有可无 0x 前缀；返回解码后的 []byte。
func ReadBin(path string) ([]byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 bin 文件 %s: %w", path, err)
	}
	hexStr := strings.TrimSpace(string(raw))
	hexStr = strings.TrimPrefix(hexStr, "0x")
	return hex.DecodeString(hexStr)
}

// ReadABI 从 .abi 文件（JSON）读取并解析为 go-ethereum 的 abi.ABI。
// 用于纯 ethclient 部署时打包构造函数参数等，无需依赖 abigen 生成的 binding。
func ReadABI(path string) (abi.ABI, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("读取 abi 文件 %s: %w", path, err)
	}
	parsed, err := abi.JSON(strings.NewReader(string(raw)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("解析 ABI %s: %w", path, err)
	}
	return parsed, nil
}
