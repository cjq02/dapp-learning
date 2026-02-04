package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"}],"name":"getItem","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://sepolia.infura.io/v3/YOUR_API_KEY"
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 练习：从环境变量读取一笔已知的 setItem 交易哈希
	txHashHex := os.Getenv("TX_HASH")
	if txHashHex == "" {
		log.Fatal("请设置环境变量 TX_HASH（一笔调用 Store.setItem 的交易哈希）")
	}
	txHash := common.HexToHash(txHashHex)

	// 练习：使用 client.TransactionReceipt 获取该交易的收据
	// receipt, err := client.TransactionReceipt(context.Background(), txHash)
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	// ItemSet 事件签名：用于匹配 receipt.Logs 中哪条是 ItemSet
	eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	// 练习：遍历 receipt.Logs，若 log.Topics[0] == eventSig，则认为是 ItemSet，用 UnpackIntoInterface 解码并打印
	for i, vLog := range receipt.Logs {
		if len(vLog.Topics) == 0 {
			continue
		}
		if vLog.Topics[0] != eventSig {
			continue
		}
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Printf("log[%d] Unpack err: %v", i, err)
			continue
		}
		fmt.Println("ItemSet event from receipt:")
		fmt.Println("  Key(hex):  ", common.Bytes2Hex(event.Key[:]))
		fmt.Println("  Value(hex):", common.Bytes2Hex(event.Value[:]))
		fmt.Println("  BlockNumber:", receipt.BlockNumber)
		fmt.Println("  TxHash:", receipt.TxHash.Hex())
	}
}
