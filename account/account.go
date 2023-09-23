package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func GenerateAccount() error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	// private key in wallet
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("private key:", hexutil.Encode(privateKeyBytes)[2:]) // [2:] is to cut 0x

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return err
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("public key:", hexutil.Encode(publicKeyBytes)[4:]) // [4:] is to cut 0x04

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address:", address)

	fmt.Println("or")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("address:", hexutil.Encode(hash.Sum(nil)[12:]))

	return nil
}

func CheckIsAddress(client *ethclient.Client, address string) error {
	account := common.HexToAddress(address)

	bytecode, err := client.CodeAt(context.Background(), account, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	isContract := len(bytecode) > 0 // if there is bytecode, it is a smartcontract not ethereum account

	fmt.Printf("is contract: %v\n", isContract) // is contract: true
	return nil
}

func Balance(client *ethclient.Client, address string) error {
	account := common.HexToAddress(address)
	balanceInWei, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return err
	}

	fbalance := new(big.Float)
	fbalance.SetString(balanceInWei.String())
	balanceInEth := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println("balance(wei):", balanceInWei)
	fmt.Println("balance(eth):", balanceInEth)

	return nil
}

// generate keystore (same as geth account new)
func GenerateKeystore(password string) error {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return err
	}

	fmt.Println("account:", account.Address.Hex())

	return nil
}

// import keystore (generate new file in tmp & remove old file in wallets)
func ImportKeystore(file string, password string) error {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		return err
	}

	fmt.Println("account:", account.Address.Hex())

	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}
