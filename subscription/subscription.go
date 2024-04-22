package subscription

import (
	"context"
	"evm-api/contract/getset"
	"evm-api/contract/store"
	"evm-api/contract/token"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robfig/cron/v3"
)

type ItemSetEvent struct {
	Key   [32]byte
	Value [32]byte
}

func SubscriptionBlock() {
	ws, err := ethclient.DialContext(context.Background(), os.Getenv("WSS_RPC_MUMBAI_NETWORK"))
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
	ws, err := ethclient.DialContext(context.Background(), os.Getenv("WSS_RPC_MUMBAI_NETWORK"))
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
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}

func SubscriptionEvents() error {
	wss, err := ethclient.DialContext(context.Background(), os.Getenv("WSS_RPC_MUMBAI_NETWORK"))
	if err != nil {
		return err
	}

	tokenAbi, err := abi.JSON(strings.NewReader(token.TokenABI))
	if err != nil {
		return err
	}

	tokenContractAddress := common.HexToAddress("0x65Dd393D0Fdd2866e37bDBC2F4ff46CCD5DfD82A")
	tokenQuery := ethereum.FilterQuery{
		Addresses: []common.Address{tokenContractAddress},
	}

	getsetAbi, err := abi.JSON(strings.NewReader(getset.GetsetABI))
	if err != nil {
		return err
	}

	getsetContractAddress := common.HexToAddress("0xB0fE721b5258Ef48AF16af3d435279100174B4AC")
	getsetQuery := ethereum.FilterQuery{
		Addresses: []common.Address{getsetContractAddress},
	}

	nftTransferAddresses := make([]common.Address, 0)
	for _, v := range []string{"0xD7AFDF92c9db81414628A21fE59085C596F85a9D", "0xA6BE298A6f6363AdB579958F93ED078B961BD8b1", "0x1af0C121Ed05626ecF8c2d36bF12Fc435d20bE9f"} {
		nftTransferAddresses = append(nftTransferAddresses, common.HexToAddress(v))
	}
	nftTransferSignature := []byte("Transfer(address,address,uint256)")
	nftTransferTopics := [][]common.Hash{{crypto.Keccak256Hash(nftTransferSignature)}}
	nftTransferQuery := ethereum.FilterQuery{
		Addresses: nftTransferAddresses,
		Topics:    nftTransferTopics,
	}

	// header
	headerLogs := make(chan *types.Header)
	headerSub, err := wss.SubscribeNewHead(context.Background(), headerLogs)
	if err != nil {
		return err
	}
	defer headerSub.Unsubscribe()

	// token
	tokenLogs := make(chan types.Log)
	tokenSub, err := wss.SubscribeFilterLogs(context.Background(), tokenQuery, tokenLogs)
	if err != nil {
		return err
	}
	defer tokenSub.Unsubscribe()

	// getset
	getsetSub, err := wss.SubscribeFilterLogs(context.Background(), getsetQuery, tokenLogs)
	if err != nil {
		return err
	}
	defer getsetSub.Unsubscribe()

	// nft
	transferNftLogs := make(chan types.Log)
	transferNftSub, err := wss.SubscribeFilterLogs(context.Background(), nftTransferQuery, transferNftLogs)
	if err != nil {
		return err
	}
	defer transferNftSub.Unsubscribe()

	for {
		select {
		case err := <-headerSub.Err():
			log.Fatal(err)
		case err := <-tokenSub.Err():
			log.Fatal(err)
		case err := <-getsetSub.Err():
			log.Fatal(err)
		case err := <-transferNftSub.Err():
			log.Fatal(err)
		case hLog := <-headerLogs:
			fmt.Println(hLog.Number)
		case tLog := <-tokenLogs:
			switch tLog.Topics[0].Hex() {
			case crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex():
				fmt.Println("topic:", tLog.Topics)
				fmt.Println("data:", hexutil.Encode(tLog.Data))
				fmt.Println("index:", tLog.Index)
				fmt.Println("tx index:", tLog.TxIndex)
				fmt.Println("removed", tLog.Removed)
				fmt.Println(tLog)

				transferSignature := []byte("Transfer(address,address,uint256)")
				transferTopics := crypto.Keccak256Hash(transferSignature)
				fmt.Println(transferTopics)

				fmt.Println(common.TrimLeftZeroes(tLog.Topics[1].Bytes()))
				fmt.Println(hexutil.Encode(common.TrimLeftZeroes(tLog.Topics[1].Bytes())))
				fmt.Println(tLog.Topics[1].Hex())
				fmt.Println(common.HexToAddress(tLog.Topics[2].Hex()))
				fmt.Println(tLog.Topics[2].Hex())

				var result token.TokenTransfer
				if err := tokenAbi.UnpackIntoInterface(&result, "Transfer", tLog.Data); err != nil {
					log.Fatal(err)
				}
				fmt.Println(result)
			case crypto.Keccak256Hash([]byte("SetA(uint256)")).Hex():
				fmt.Println("topic:", tLog.Topics)
				fmt.Println("data:", hexutil.Encode(tLog.Data))
				fmt.Println("index:", tLog.Index)
				fmt.Println("tx index:", tLog.TxIndex)
				fmt.Println("removed", tLog.Removed)
				fmt.Println(tLog)

				setASignature := []byte("SetA(uint256)")
				setATopics := crypto.Keccak256Hash(setASignature)
				fmt.Println(setATopics)

				var result getset.GetsetSetA
				if err := getsetAbi.UnpackIntoInterface(&result, "SetA", tLog.Data); err != nil {
					log.Fatal(err)
				}
				fmt.Println(result)
			}
		case nLog := <-transferNftLogs:
			fmt.Println(nLog)
		}
	}
}

func SubscriptionEventWABI() error {
	wss, err := ethclient.DialContext(context.Background(), os.Getenv("WSS_RPC_MUMBAI_NETWORK"))
	if err != nil {
		return err
	}
	getsetContractAddress := common.HexToAddress("0xB0fE721b5258Ef48AF16af3d435279100174B4AC")

	instance, err := getset.NewGetset(getsetContractAddress, wss)
	if err != nil {
		return err
	}
	logs := make(chan *getset.GetsetSetA)
	watchOpts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	sub, err := instance.WatchSetA(watchOpts, logs)
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case log := <-logs:
			fmt.Println(log)
			fmt.Println(log.A)
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

func ReadEvents() error {
	https, err := ethclient.DialContext(context.Background(), os.Getenv("HTTPS_RPC_MUMBAI_NETWORK"))
	if err != nil {
		return err
	}

	getsetAbi, err := abi.JSON(strings.NewReader(getset.GetsetABI))
	if err != nil {
		return err
	}

	getsetContractAddress := common.HexToAddress("0xB0fE721b5258Ef48AF16af3d435279100174B4AC")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(34517182),
		ToBlock:   big.NewInt(34519009),
		Addresses: []common.Address{getsetContractAddress},
	}

	logs, err := https.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)

		fmt.Println(vLog)
		var result getset.GetsetSetA
		if err := getsetAbi.UnpackIntoInterface(&result, "SetA", vLog.Data); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.A)
	}

	return nil
}

func ReadEventWABI() error {
	https, err := ethclient.DialContext(context.Background(), os.Getenv("HTTPS_RPC_MUMBAI_NETWORK"))
	if err != nil {
		return err
	}
	getsetContractAddress := common.HexToAddress("0xB0fE721b5258Ef48AF16af3d435279100174B4AC")

	instance, err := getset.NewGetset(getsetContractAddress, https)
	if err != nil {
		return err
	}
	end := uint64(34519009)
	filterOpts := &bind.FilterOpts{Context: context.Background(), Start: 34517182, End: &end}
	rows, err := instance.FilterSetA(filterOpts)
	if err != nil {
		return err
	}

	for rows.Next() {
		fmt.Printf("Log Block Number: %d\n", rows.Event.Raw.BlockNumber)
		fmt.Println(rows.Event)
		fmt.Println(rows.Event.A)
	}
	defer rows.Close()

	return nil
}

func PullingInterval() {
	c := cron.New(cron.WithLocation(time.Local))
	_, _ = c.AddFunc("* * * * *", func() {
		fmt.Println("[Job 1]Every minute job")
	})

	_, _ = c.AddFunc("*/2 * * * *", func() {
		fmt.Println("[Job 2]Every two minutes job")
	})

	_, _ = c.AddFunc("@hourly", func() {
		fmt.Println("Every hour")
	})

	fmt.Printf("Cron Info: %+v\n", c.Entries())
	c.Start()

	time.Sleep(2 * time.Hour)
}
