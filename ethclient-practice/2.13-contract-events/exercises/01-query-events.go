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
	var toHeader *types.Header
	var fromHeader *types.Header
	var query ethereum.FilterQuery
	// TODO: 请在此实现 fromBlock、toBlock 的赋值
	fromHeader, err = client.HeaderByNumber(context.Background(), big.NewInt(10202504))
	if err != nil {
		log.Fatal(err)
	}
	toHeader, err = client.HeaderByNumber(context.Background(), big.NewInt(10222940))
	if err != nil {
		log.Fatal(err)
	}
	toBlock = toHeader.Number
	fromBlock = fromHeader.Number
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
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		// Data 里只有非 indexed 参数，即 value
		err = parsedABI.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Printf("Unpack err: %v", err)
			continue
		}
		// indexed 的 key 在 topics[1]，32 字节即 bytes32
		if len(vLog.Topics) > 1 {
			copy(event.Key[:], vLog.Topics[1].Bytes())
		}
		fmt.Println("BlockHash:", vLog.BlockHash.Hex())
		fmt.Println("BlockNumber:", vLog.BlockNumber)
		fmt.Println("TxHash:", vLog.TxHash.Hex())

		var block *types.Block
		block, err = client.BlockByNumber(context.Background(), big.NewInt(int64(vLog.BlockNumber)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BlockTime:", block.Time(), "->", time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"))
		// bytes32 若存的是 ASCII，去掉尾部 0 可读
		// keyStr := strings.TrimRight(string(event.Key[:]), "\x00")
		keyStr := string(event.Key[:])
		// valueStr := strings.TrimRight(string(event.Value[:]), "\x00")
		valueStr := string(event.Value[:])
		fmt.Println("Key(hex):", common.Bytes2Hex(event.Key[:]), "->", keyStr)
		fmt.Println("Value(hex):", common.Bytes2Hex(event.Value[:]), "->", valueStr)
		for i, t := range vLog.Topics {
			desc := ""
			switch i {
			case 0:
				desc = " (事件签名哈希 ItemSet(bytes32,bytes32))"
			case 1:
				// desc = " (indexed key) -> " + strings.TrimRight(string(t.Bytes()), "\x00")
				desc = " (indexed key) -> " + string(t.Bytes())
			}
			fmt.Printf("Topic[%d]=%s%s\n", i, t.Hex(), desc)
		}
		fmt.Println("Data:", common.Bytes2Hex(vLog.Data), "->", string(vLog.Data))
		fmt.Println("---")
	}
	fmt.Printf("共 %d 条日志\n", len(logs))
}
