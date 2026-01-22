// 01-generate-wallet.go - 生成新钱包练习
//
// 任务：
// 1. 使用 crypto.GenerateKey() 生成一个新的随机私钥
// 2. 从私钥派生公钥
// 3. 从公钥生成地址
// 4. 输出所有信息
//
// 运行：go run exercises/01-generate-wallet.go

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

	// TODO 1: 生成新的随机私钥
	// 提示：使用 crypto.GenerateKey()
	var privateKey *ecdsa.PrivateKey
	{
		// 在这里填写代码
		var err error
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO 2: 将私钥转换为字节，然后转为十六进制字符串（去掉 '0x' 前缀）
	// 提示：使用 crypto.FromECDSA() 和 hexutil.Encode()
	var privateKeyHex string
	{
		// 在这里填写代码
		privateKeyHex = hexutil.Encode(crypto.FromECDSA(privateKey))
	}
	fmt.Printf("私钥: %s\n", privateKeyHex)

	// TODO 3: 从私钥获取公钥，并进行类型断言
	// 提示：使用 privateKey.Public() 和类型断言到 *ecdsa.PublicKey
	var publicKeyECDSA *ecdsa.PublicKey
	{
		// 在这里填写代码
		publicKeyECDSA = privateKey.Public().(*ecdsa.PublicKey)
	}

	// TODO 4: 将公钥转换为字节，然后转为十六进制字符串（去掉 '0x' 和 '0x04' 前缀）
	// 提示：使用 crypto.FromECDSAPub() 和 hexutil.Encode()
	var publicKeyHex string
	{
		// 在这里填写代码
		publicKeyHex = hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))
	}
	fmt.Printf("公钥: %s\n", publicKeyHex)

	// TODO 5: 从公钥生成地址
	// 提示：使用 crypto.PubkeyToAddress()
	var address string
	{
		// 在这里填写代码
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	}
	fmt.Printf("地址: %s\n", address)

	fmt.Println("=== 完成 ===")
}
