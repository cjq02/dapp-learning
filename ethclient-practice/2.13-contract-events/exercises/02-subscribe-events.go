package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/dapp-learning/ethclient-practice/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取配置（与 2.12 一致；订阅需 WebSocket）
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		log.Fatal("错误: 请设置环境变量 CONTRACT_ADDRESS（可填 2.10 部署得到的合约地址）")
	}
	wsURL := os.Getenv("SEPOLIA_WS_URL")
	if wsURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_WS_URL（订阅需用 wss://，见 contract-events.md 测试网资源）")
	}
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contractAddress := common.HexToAddress(contractAddressStr)

	// -------------------------------------------------------------------------
	// 练习 1：构造 FilterQuery（订阅时可不设 FromBlock/ToBlock，仅填 Addresses）
	// TODO: 请将 Addresses 设为 []common.Address{contractAddress}，使订阅只收到该合约的日志
	// -------------------------------------------------------------------------
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// -------------------------------------------------------------------------
	// 练习 2：创建 types.Log 的 channel，调用 SubscribeFilterLogs，并 defer Unsubscribe
	// TODO: 请在此补全：logsCh := make(chan types.Log)；sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)；defer sub.Unsubscribe()
	// -------------------------------------------------------------------------
	// logsCh := make(chan types.Log)
	// sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer sub.Unsubscribe()
	var logsCh chan types.Log
	logsCh = make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	contractAbi, err := util.ReadABI("../contract/Store_sol_Store.abi")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("等待 ItemSet 事件（可另开终端调用 setItem 触发）...")

	// -------------------------------------------------------------------------
	// 练习 3：for + select 循环，处理 sub.Err() 与 logsCh
	// 收到 vLog 时用 contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data) 解码并打印
	// 提示：event 结构体 struct{ Key, Value [32]byte }；可打印 BlockNumber、TxHash、Key、Value、Topics
	// -------------------------------------------------------------------------
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logsCh:
			// TODO: 在此解码 vLog 为 ItemSet 事件并打印
			event := struct {
				Key   [32]byte
				Value [32]byte
			}{}
			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Printf("Unpack err: %v", err)
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
	}
}
