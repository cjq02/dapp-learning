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
	"github.com/ethereum/go-ethereum/core/types"
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

	// -------------------------------------------------------------------------
	// 练习 1：从环境变量读取一笔 setItem 交易哈希，并获取该交易的收据
	// 提示：TX_HASH 环境变量；client.TransactionReceipt(context.Background(), txHash)
	// -------------------------------------------------------------------------
	txHashHex := os.Getenv("TX_HASH")
	if txHashHex == "" {
		log.Fatal("错误: 请设置环境变量 TX_HASH（一笔调用 Store.setItem 的交易哈希）")
	}
	txHash := common.HexToHash(txHashHex)

	var receipt *types.Receipt
	var receiptErr error
	// TODO: 请在此调用 client.TransactionReceipt(context.Background(), txHash)，将结果赋给 receipt 和 receiptErr
	_, _ = context.Background(), txHash
	if receiptErr != nil {
		log.Fatal(receiptErr)
	}

	// ItemSet 事件签名，用于匹配 log.Topics[0]
	eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// 练习 2：遍历 receipt.Logs，若 log.Topics[0] == eventSig 则认为是 ItemSet 事件
	// 使用 contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data) 解码并打印 Key、Value 等
	// -------------------------------------------------------------------------
	for i, vLog := range receipt.Logs {
		// TODO: 在此判断是否为 ItemSet（len(vLog.Topics)>0 且 vLog.Topics[0]==eventSig），若是则解码并打印
		_, _ = i, vLog
		_ = eventSig
		_ = contractAbi
	}
	fmt.Println("解析完成")
}
