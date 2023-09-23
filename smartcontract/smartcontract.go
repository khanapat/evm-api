package smartcontract

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"evm-api/contract/store"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Load(client *ethclient.Client, contractAddress string) (*store.Store, error) {
	address := common.HexToAddress(contractAddress)
	instance, err := store.NewStore(address, client)
	if err != nil {
		return nil, err
	}

	version, err := instance.Version(nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Version:", version)

	return instance, nil
}

func GetBytecode(client *ethclient.Client, contractAddress string) error {
	address := common.HexToAddress(contractAddress)
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		return err
	}

	fmt.Println("Bytecode:", hex.EncodeToString(bytecode))

	return nil
}

func Deploy(client *ethclient.Client) error {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("Address:", fromAddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	fmt.Println("Nonce:", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("GasPrice:", gasPrice)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	version := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, version)
	if err != nil {
		return err
	}

	fmt.Println("Contract Address:", address.Hex())
	fmt.Println("Txn Hash:", tx.Hash().Hex())

	_ = instance

	return nil
}

func SetStore(client *ethclient.Client, instance *store.Store) error {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("Address:", fromAddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	fmt.Println("Nonce:", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	storeABI, err := abi.JSON(strings.NewReader(store.StoreABI))
	if err != nil {
		return err
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	data, err := storeABI.Pack("setItem", key, value)
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress("0xC8A0FE1489cCF266c3011b49c38769f7ba7624C2")
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: big.NewInt(0),
		Data:  data,
	})
	if err != nil {
		return err
	}
	fmt.Println("GasLimit:", gasLimit)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("GasPrice:", gasPrice)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	// send txn
	// tx, err := instance.SetItem(auth, key, value)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	// receipt, err := bind.WaitMined(context.Background(), client, tx)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("usage by txn:", receipt.GasUsed) // usage by txn
	// fmt.Println("cumulative gas used:", receipt.CumulativeGasUsed)
	// fmt.Println("effective gas price:", receipt.EffectiveGasPrice) // gas price

	// call
	result, err := instance.Items(nil, key) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items:", string(result[:]))

	return nil
}
