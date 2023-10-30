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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Load(client *ethclient.Client, contractAddress string) (*store.Store, error) {
	address := common.HexToAddress(contractAddress)

	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Byte Code:", hexutil.Encode(bytecode))

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

func SetStore(client *ethclient.Client, contract string, instance *store.Store) error {
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

	realNonce, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	fmt.Println("Real Nonce:", realNonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	storeABI, err := abi.JSON(strings.NewReader(store.StoreABI))
	if err != nil {
		return err
	}

	// key := [32]byte{}
	// value := [32]byte{}
	// copy(key[:], []byte("foo"))
	// copy(value[:], []byte("bobo"))

	data, err := storeABI.Pack("setItem", "earing", big.NewInt(1))
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress(contract)
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
	if gasPrice.Cmp(big.NewInt(0)) == 0 {
		gasPrice = big.NewInt(10000000000)
	}
	fmt.Println("GasPrice:", gasPrice)

	// arise don't have this method
	// gasTipCap, err := client.SuggestGasTipCap(context.Background())
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("GasTipCap:", gasTipCap)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	// send txn
	key1 := "trust"
	// value1 := big.NewInt(1)
	// tx1, err := instance.SetItem(auth, key1, value1)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx1 sent: %s\n", tx1.Hash().Hex())

	// auth.Nonce = big.NewInt(int64(nonce + 1))

	key2 := "note"
	// value2 := big.NewInt(6)
	// tx2, err := instance.SetItem(auth, key2, value2)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx2 sent: %s\n", tx2.Hash().Hex())

	// auth.Nonce = big.NewInt(int64(nonce + 2))

	key3 := "earing"
	// value3 := big.NewInt(7)
	// tx3, err := instance.SetItem(auth, key3, value3)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx3 sent: %s\n", tx3.Hash().Hex())

	// auth.Nonce = big.NewInt(int64(nonce + 3))

	key4 := "chawin"
	// value4 := big.NewInt(8)
	// tx4, err := instance.SetItem(auth, key4, value4)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("tx4 sent: %s\n", tx4.Hash().Hex())

	nonce2, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	fmt.Println("Nonce:", nonce2)

	realNonce2, err := client.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return err
	}
	fmt.Println("Real Nonce:", realNonce2)

	// receipt, err := bind.WaitMined(context.Background(), client, tx)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("usage by txn:", receipt.GasUsed) // usage by txn
	// fmt.Println("cumulative gas used:", receipt.CumulativeGasUsed)
	// fmt.Println("effective gas price:", receipt.EffectiveGasPrice) // gas price

	// call
	result1, err := instance.Items(nil, key1) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items1:", result1)

	result2, err := instance.Items(nil, key2) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items2:", result2)

	result3, err := instance.Items(nil, key3) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items3:", result3)

	result4, err := instance.Items(nil, key4) // nil or &bind.CallOpts{} is the same
	if err != nil {
		return err
	}

	fmt.Println("items4:", result4)

	// fmt.Println("items:", string(result[:]))

	return nil
}
