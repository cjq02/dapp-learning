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
	// 从环境变量读取配置（与 2.12 一致）
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	txHashHex := os.Getenv("TX_HASH")
	if txHashHex == "" {
		log.Fatal("错误: 请设置环境变量 TX_HASH（一笔调用 Store.setItem 的交易哈希）")
	}
	txHash := common.HexToHash(txHashHex)

	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

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
