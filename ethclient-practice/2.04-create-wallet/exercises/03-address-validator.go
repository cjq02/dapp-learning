// 03-address-validator.go - 地址生成验证器练习
//
// 任务：
// 1. 生成新钱包
// 2. 使用 crypto.PubkeyToAddress() 生成地址（方法 A）
// 3. 手动使用 Keccak-256 哈希生成地址（方法 B）
// 4. 验证两种方法生成的地址一致
// 5. 输出详细步骤，帮助理解地址生成过程
//
// 运行：go run exercises/03-address-validator.go

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func main() {
	fmt.Println("=== 地址生成验证器 ===")

	// 步骤 1: 生成新钱包
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	fmt.Printf("\n步骤 1: 公钥信息\n")
	fmt.Printf("原始公钥长度: %d 字节\n", len(publicKeyBytes))
	fmt.Printf("原始公钥 (Hex): %s\n", hexutil.Encode(publicKeyBytes))

	// 步骤 2: 方法 A - 使用内置函数生成地址
	// TODO 1: 使用 crypto.PubkeyToAddress() 生成地址
	var addressMethodA string
	{
		// 在这里填写代码
	}
	fmt.Printf("\n步骤 2: 方法 A (内置函数)\n")
	fmt.Printf("地址: %s\n", addressMethodA)

	// 步骤 3: 方法 B - 手动使用 Keccak-256 哈希生成地址
	fmt.Printf("\n步骤 3: 方法 B (手动计算)\n")

	// TODO 2: 跳过公钥的第一个字节（0x04），然后计算 Keccak-256 哈希
	// 提示：使用 publicKeyBytes[1:] 和 sha3.NewLegacyKeccak256()
	var hash []byte
	{
		// 在这里填写代码
	}
	fmt.Printf("完整哈希 (32字节): %s\n", hexutil.Encode(hash))

	// TODO 3: 取哈希值的后 20 字节作为地址
	var addressMethodB string
	{
		// 在这里填写代码
	}
	fmt.Printf("地址 (后20字节): %s\n", addressMethodB)

	// 步骤 4: 验证两种方法结果一致
	fmt.Printf("\n步骤 4: 验证结果\n")
	// TODO 4: 比较两种方法生成的地址是否一致
	{
		// 在这里填写代码
	}

	// 步骤 5: 详细解释
	fmt.Printf("\n步骤 5: 原理解释\n")
	fmt.Printf("1. 公钥格式: 0x04 (1字节) + X坐标 (32字节) + Y坐标 (32字节) = 65字节\n")
	fmt.Printf("2. 跳过第一个字节后: 64字节\n")
	fmt.Printf("3. Keccak-256 哈希: 32字节\n")
	fmt.Printf("4. 取后20字节: 以太坊地址\n")
	fmt.Printf("5. 为什么取后20字节? 以太坊设计选择，平衡安全性和效率\n")

	fmt.Println("=== 完成 ===")
}
