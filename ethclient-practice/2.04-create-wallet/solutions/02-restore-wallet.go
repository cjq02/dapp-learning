// 02-restore-wallet.go - 从私钥恢复钱包练习 - 标准答案
//
// 运行：go run solutions/02-restore-wallet.go

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
	privateKeyHex := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"

	// 步骤 1: 从十六进制字符串恢复私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// 步骤 2: 从私钥获取公钥，并进行类型断言
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 步骤 3: 从公钥生成地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("恢复的地址: %s\n", address)

	// 步骤 4: 验证：从同一私钥再次恢复，验证地址是否一致
	privateKey2, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey2 := privateKey2.Public()
	publicKeyECDSA2, ok := publicKey2.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey2 is not of type *ecdsa.PublicKey")
	}

	address2 := crypto.PubkeyToAddress(*publicKeyECDSA2).Hex()

	if address == address2 {
		fmt.Printf("验证：从同一私钥恢复的地址一致 ✓\n")
	} else {
		fmt.Printf("验证失败！地址不一致\n")
	}

	// 额外验证：比较私钥本身
	if comparePrivateKeys(privateKey, privateKey2) {
		fmt.Printf("额外验证：私钥完全相同 ✓\n")
	}

	fmt.Println("=== 完成 ===")
	fmt.Println("\n📝 提示:")
	fmt.Println("- 私钥到地址的映射是确定性的")
	fmt.Println("- 相同的私钥总是生成相同的地址")
	fmt.Println("- 这就是为什么可以通过私钥恢复钱包")
}

// 辅助函数：比较两个私钥是否相同
func comparePrivateKeys(k1, k2 *ecdsa.PrivateKey) bool {
	// 比较公钥的 X 和 Y 坐标
	return k1.X.Cmp(k2.X) == 0 && k1.Y.Cmp(k2.Y) == 0
}
