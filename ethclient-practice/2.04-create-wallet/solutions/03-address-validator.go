// 03-address-validator.go - 地址生成验证器练习 - 标准答案
//
// 运行：go run solutions/03-address-validator.go

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
	addressMethodA := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("\n步骤 2: 方法 A (内置函数)\n")
	fmt.Printf("地址: %s\n", addressMethodA)

	// 步骤 3: 方法 B - 手动使用 Keccak-256 哈希生成地址
	fmt.Printf("\n步骤 3: 方法 B (手动计算)\n")

	// 跳过第一个字节（0x04），然后计算 Keccak-256 哈希
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	hashBytes := hash.Sum(nil)

	fmt.Printf("完整哈希 (32字节): %s\n", hexutil.Encode(hashBytes))

	// 取哈希值的后 20 字节作为地址
	addressBytes := hashBytes[12:]
	addressMethodB := hexutil.Encode(addressBytes)
	fmt.Printf("地址 (后20字节): %s\n", addressMethodB)

	// 步骤 4: 验证两种方法结果一致
	fmt.Printf("\n步骤 4: 验证结果\n")
	if addressMethodA == addressMethodB {
		fmt.Printf("✓ 两种方法生成的地址完全一致！\n")
	} else {
		fmt.Printf("✗ 地址不一致！\n")
		fmt.Printf("  方法 A: %s\n", addressMethodA)
		fmt.Printf("  方法 B: %s\n", addressMethodB)
	}

	// 步骤 5: 详细解释
	fmt.Printf("\n步骤 5: 原理解释\n")
	fmt.Printf("1. 公钥格式: 0x04 (1字节) + X坐标 (32字节) + Y坐标 (32字节) = 65字节\n")
	fmt.Printf("2. 跳过第一个字节后: 64字节\n")
	fmt.Printf("3. Keccak-256 哈希: 32字节\n")
	fmt.Printf("4. 取后20字节: 以太坊地址\n")
	fmt.Printf("5. 为什么取后20字节? 以太坊设计选择，平衡安全性和效率\n")

	// 额外：展示完整的计算过程
	fmt.Printf("\n额外：完整计算过程\n")
	fmt.Printf("公钥 (完整): %s\n", hexutil.Encode(publicKeyBytes))
	fmt.Printf("公钥 (去掉0x04): %s\n", hexutil.Encode(publicKeyBytes[1:]))
	fmt.Printf("Keccak-256 哈希: %s\n", hexutil.Encode(hashBytes))
	fmt.Printf("哈希前12字节: %s (丢弃)\n", hexutil.Encode(hashBytes[:12]))
	fmt.Printf("哈希后20字节: %s (地址)\n", hexutil.Encode(hashBytes[12:]))

	fmt.Println("\n=== 完成 ===")
	fmt.Println("\n📝 关键要点:")
	fmt.Println("- 以太坊地址 = Keccak-256(公钥[1:]) 的后 20 字节")
	fmt.Println("- 公钥第一个字节 0x04 是 EC 前缀，不参与哈希计算")
	fmt.Println("- 32 字节哈希中，只保留后 20 字节作为地址")
	fmt.Println("- 这是一种安全性和效率的平衡设计")
}
