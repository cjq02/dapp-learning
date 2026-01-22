// 02-restore-wallet.go - 从私钥恢复钱包练习
//
// 任务：
// 1. 从已有的私钥（十六进制字符串）恢复钱包
// 2. 输出恢复的地址
// 3. 验证：从同一私钥多次恢复，地址应该一致
//
// 运行：go run exercises/02-restore-wallet.go

package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("=== 从私钥恢复钱包 ===")

	// 已有的私钥（十六进制格式，不带 '0x' 前缀）
	privateKeyHex := "fbf3eb52157e3e91e76d042b06eb77d516d95eed87646ef5e792302850dd2bc4"

	// TODO 1: 从十六进制字符串恢复私钥
	// 提示：使用 crypto.HexToECDSA()
	var privateKey *ecdsa.PrivateKey
	{
		// 在这里填写代码
		var err error
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO 2: 从私钥获取公钥，并进行类型断言
	var publicKeyECDSA *ecdsa.PublicKey
	{
		// 在这里填写代码
		publicKeyECDSA = privateKey.Public().(*ecdsa.PublicKey)
	}

	// TODO 3: 从公钥生成地址
	var address string
	{
		// 在这里填写代码
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	}
	fmt.Printf("恢复的地址: %s\n", address)

	// TODO 4: 验证：从同一私钥再次恢复，验证地址是否一致
	// 提示：再次调用 crypto.HexToECDSA() 并比较地址
	{
		// 在这里填写代码
		privateKey2, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			log.Fatal(err)
		}
		publicKeyECDSA2 := privateKey2.Public().(*ecdsa.PublicKey)
		address2 := crypto.PubkeyToAddress(*publicKeyECDSA2).Hex()
		if address == address2 {
			fmt.Println("验证：从同一私钥恢复的地址一致")
		} else {
			fmt.Println("验证：从同一私钥恢复的地址不一致")
		}
	}

	fmt.Println("=== 完成 ===")
}

// 辅助函数：比较两个私钥是否相同
func comparePrivateKeys(k1, k2 *ecdsa.PrivateKey) bool {
	// 比较公钥的 X 和 Y 坐标
	return k1.X.Cmp(k2.X) == 0 && k1.Y.Cmp(k2.Y) == 0
}
