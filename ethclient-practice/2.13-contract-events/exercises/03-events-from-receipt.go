package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dapp-learning/ethclient-practice/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
	// 练习 1：从环境变量读取交易哈希，并获取该交易的收据
	// TODO: 调用 client.TransactionReceipt(context.Background(), txHash)，将结果赋给 receipt 和 receiptErr
	// -------------------------------------------------------------------------
	txHashHex := os.Getenv("TX_HASH")
	if txHashHex == "" {
		log.Fatal("错误: 请设置环境变量 TX_HASH（一笔调用 Store.setItem 的交易哈希）")
	}
	txHash := common.HexToHash(txHashHex)

	var receipt *types.Receipt
	var receiptErr error
	// TODO: receipt, receiptErr = client.TransactionReceipt(context.Background(), txHash)
	receipt, receiptErr = client.TransactionReceipt(context.Background(), txHash)
	if receiptErr != nil {
		log.Fatal(receiptErr)
	}

	// ItemSet 事件签名：keccak256("ItemSet(bytes32,bytes32)")，用于判断 log 是否为 ItemSet
	eventSig := crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)"))

	contractAbi, err := util.ReadABI("../contract/Store_sol_Store.abi")
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// 练习 2：遍历 receipt.Logs，识别 ItemSet 事件并解码打印
	// TODO: for range receipt.Logs {
	//   1. 若 len(vLog.Topics)==0 或 vLog.Topics[0]!=eventSig，则 continue（不是 ItemSet）
	//   2. 定义 event 结构体 struct{ Key, Value [32]byte }
	//   3. contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data) 解出 value（Data 里只有非 indexed）
	//   4. 从 vLog.Topics[1] 拷贝到 event.Key（indexed 的 key 在 topics[1]）
	//   5. 打印 Key、Value、BlockNumber、TxHash 等
	// }
	// -------------------------------------------------------------------------
	for i, vLog := range receipt.Logs {
		// TODO: 在此实现上述步骤
		_, _ = i, vLog
		if len(vLog.Topics) == 0 || vLog.Topics[0] != eventSig {
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
		if len(vLog.Topics) > 1 {
			copy(event.Key[:], vLog.Topics[1].Bytes())
		}
		fmt.Println("BlockNumber:", vLog.BlockNumber, "TxHash:", vLog.TxHash.Hex())
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(vLog.BlockNumber)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BlockTime:", block.Time(), "->", time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"))
		keyStr := string(event.Key[:])
		valueStr := string(event.Value[:])
		fmt.Println("key(hex):", common.Bytes2Hex(event.Key[:]), "->", keyStr)
		fmt.Println("value(hex):", common.Bytes2Hex(event.Value[:]), "->", valueStr)
	}
	fmt.Println("解析完成")
}
