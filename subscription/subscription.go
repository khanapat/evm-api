package subscription

import (
	"context"
	"evm-api/store"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ItemSetEvent struct {
	Key   [32]byte
	Value [32]byte
}

func SubscriptionBlock() {
	ws, err := ethclient.DialContext(context.Background(), os.Getenv("WS_RPC_NETWORK"))
	if err != nil {
		log.Fatal(err.Error())
	}
	headers := make(chan *types.Header)

	sub, err := ws.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := ws.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println("block hash:", block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println("block number:", block.Number().Uint64()) // 3477413
			fmt.Println("block time:", block.Time())              // 1529525947
			fmt.Println("block nonce:", block.Nonce())            // 130524141876765836
			fmt.Println("block txns:", len(block.Transactions())) // 7
		}
	}
}

func SubscriptionEvent() error {
	ws, err := ethclient.DialContext(context.Background(), os.Getenv("WS_RPC_NETWORK"))
	if err != nil {
		log.Fatal(err.Error())
	}
	contractAddress := common.HexToAddress("0x188992471F03D5bf1EaC66973fc3E7CA7ee5C0D3")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	logs := make(chan types.Log)
	sub, err := ws.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}

func ReadEvent(client *ethclient.Client) error {
	contractAddress := common.HexToAddress("0x188992471F03D5bf1EaC66973fc3E7CA7ee5C0D3")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(7920218),
		ToBlock:   big.NewInt(7920671),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		fmt.Println("Blockhash:", vLog.BlockHash.Hex())
		fmt.Println("BlockNumber:", vLog.BlockNumber)
		fmt.Println("TxHash:", vLog.TxHash.Hex())

		event, err := contractAbi.Unpack("ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("RawEvent:", event)

		itemSetEvent := ItemSetEvent{
			Key:   event[0].([32]byte),
			Value: event[1].([32]byte),
		}

		fmt.Println("Key:", string(itemSetEvent.Key[:]))
		fmt.Println("Value:", string(itemSetEvent.Value[:]))

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}
		fmt.Println("Topics:", topics)
		fmt.Println("Topic 0:", topics[0]) // topic ItemSet(bytes32,bytes32)

		fmt.Println("--------------")
	}

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature) // keccak256("ItemSet(bytes32,bytes32)")
	fmt.Println("Topic ItemSet(bytes32,bytes32):", hash.Hex())

	return nil
}
