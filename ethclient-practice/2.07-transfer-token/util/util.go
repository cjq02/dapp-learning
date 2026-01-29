// Package util provides common utilities for ERC20 token operations
package util

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
