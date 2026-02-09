package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/dapp-learning/ethclient-practice/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 从环境变量读取配置（与 2.12 一致）
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		log.Fatal("错误: 请设置环境变量 CONTRACT_ADDRESS（可填 2.10 部署得到的合约地址）")
	}
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		log.Fatal("错误: 请设置环境变量 SEPOLIA_RPC_URL")
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	contractAddress := common.HexToAddress(contractAddressStr)

	// -------------------------------------------------------------------------
	// 练习 1：构造区块范围
	// 提示：client.HeaderByNumber(ctx, nil) 获取最新区块头，header.Number 为区块号
	//       设置 toBlock = 最新，fromBlock = 最新 - 100（用 new(big.Int).Sub）
	// -------------------------------------------------------------------------
	var fromBlock, toBlock *big.Int
	var header *types.Header
	var query ethereum.FilterQuery
	// TODO: 请在此实现 fromBlock、toBlock 的赋值
	header, err = client.HeaderByNumber(context.Background(), big.NewInt(10204764))
	if err != nil {
		log.Fatal(err)
	}
	toBlock = header.Number
	fromBlock = new(big.Int).Sub(header.Number, big.NewInt(100))
	query = ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{contractAddress},
	}

	// -------------------------------------------------------------------------
	// 练习 2：调用 FilterLogs 获取该区块范围内、该合约的日志
	// 提示：logs, err := client.FilterLogs(context.Background(), query)
	// -------------------------------------------------------------------------
	var logs []types.Log // 练习 2 完成后此处会由 FilterLogs 填充
	var filterErr error
	// TODO: 请在此调用 client.FilterLogs(context.Background(), query)，将结果赋给 logs 和 filterErr
	logs, filterErr = client.FilterLogs(context.Background(), query)
	if filterErr != nil {
		log.Fatal(filterErr)
	}

	// -------------------------------------------------------------------------
	// 练习 3：解析 ABI，遍历 logs，解码 ItemSet 事件并打印
	// 提示：abi.JSON(strings.NewReader(StoreABI)) 得到 contractAbi
	//       event 结构体：struct{ Key [32]byte; Value [32]byte }
	//       contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data) 只解 Data
	//       打印：BlockHash、BlockNumber、TxHash、Key(hex)、Value(hex)、Topics
	// -------------------------------------------------------------------------
	parsedABI, err := util.ReadABI("../contract/Store_sol_Store.abi")
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		// TODO: 在此解码 vLog 为 ItemSet 事件并打印（Key、Value、Topics 等）
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err = parsedABI.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Printf("Unpack err: %v", err)
			continue
		}
		fmt.Println("BlockHash:", vLog.BlockHash.Hex())
		fmt.Println("BlockNumber:", vLog.BlockNumber)
		fmt.Println("TxHash:", vLog.TxHash.Hex())
		fmt.Println("Key(hex):", common.Bytes2Hex(event.Key[:]))
		fmt.Println("Value(hex):", common.Bytes2Hex(event.Value[:]))
		for i, t := range vLog.Topics {
			fmt.Printf("Topic[%d]=%s\n", i, t.Hex())
		}
		fmt.Println("---")
	}
	fmt.Printf("共 %d 条日志\n", len(logs))
}
