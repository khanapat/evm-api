package smartcontract

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"evm-api/store"
	"fmt"
	"math/big"
	"os"

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

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	// send txn
	// tx, err := instance.SetItem(auth, key, value)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx sent: %s", tx.Hash().Hex())

	// call
	result, err := instance.Items(nil, key) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items:", string(result[:]))

	return nil
}
