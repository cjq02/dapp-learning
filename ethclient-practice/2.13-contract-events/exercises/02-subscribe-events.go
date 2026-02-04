package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"}],"name":"getItem","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {
	// 订阅事件必须使用 WebSocket
	wsURL := os.Getenv("SEPOLIA_WS_URL")
	if wsURL == "" {
		wsURL = "wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"
	}
	// Alchemy WS: wss://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
	// Infura WS:  wss://sepolia.infura.io/ws/v3/YOUR_API_KEY
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	contractAddrHex := os.Getenv("STORE_CONTRACT_ADDRESS")
	if contractAddrHex == "" {
		log.Fatal("请设置环境变量 STORE_CONTRACT_ADDRESS")
	}
	contractAddress := common.HexToAddress(contractAddrHex)

	// 练习：构造 FilterQuery（订阅时可不设 FromBlock/ToBlock），仅填 Addresses
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// 练习：创建 types.Log 的 channel，并调用 client.SubscribeFilterLogs
	// 提示：logsCh := make(chan types.Log)；sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	contractAbi, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		log.Fatal(err)
	}

	// 练习：for + select 循环，从 sub.Err() 和 logsCh 读取，收到日志时解码 ItemSet 并打印
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logsCh:
			fmt.Println("BlockNumber:", vLog.BlockNumber, "TxHash:", vLog.TxHash.Hex())
			event := struct {
				Key   [32]byte
				Value [32]byte
			}{}
			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Printf("Unpack err: %v", err)
				continue
			}
			fmt.Println("Key:", common.Bytes2Hex(event.Key[:]), "Value:", common.Bytes2Hex(event.Value[:]))
			for i, t := range vLog.Topics {
				fmt.Printf("Topic[%d]=%s\n", i, t.Hex())
			}
		}
	}
}
