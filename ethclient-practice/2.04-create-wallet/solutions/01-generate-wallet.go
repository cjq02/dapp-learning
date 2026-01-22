// 01-generate-wallet.go - 生成新钱包练习 - 标准答案
//
// 运行：go run solutions/01-generate-wallet.go

package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("=== 生成新钱包 ===")

	// 步骤 1: 生成新的随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// 步骤 2: 将私钥转换为字节，然后转为十六进制字符串（去掉 '0x' 前缀）
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Printf("私钥: %s\n", privateKeyHex)

	// 步骤 3: 从私钥获取公钥，并进行类型断言
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 步骤 4: 将公钥转换为字节，然后转为十六进制字符串（去掉 '0x' 和 '0x04' 前缀）
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Printf("公钥: %s\n", publicKeyHex)

	// 步骤 5: 从公钥生成地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("地址: %s\n", address)

	fmt.Println("=== 完成 ===")
	fmt.Println("\n📝 提示:")
	fmt.Println("- 私钥是 64 个十六进制字符（32 字节）")
	fmt.Println("- 公钥是 128 个十六进制字符（64 字节）")
	fmt.Println("- 地址是 42 个字符（0x + 40 个十六进制字符 = 20 字节）")
	fmt.Println("- ⚠️  永远不要分享私钥！")
}
