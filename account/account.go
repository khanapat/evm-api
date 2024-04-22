package account

import (
	"context"
	"crypto/ecdsa"
	"evm-api/util"
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
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

// https://idhww.medium.com/making-your-own-safety-cold-ethereum-hd-wallet-using-golang-b6f34b359c8f
// https://github.com/topics/mnemonic?l=go
func GenerateMnemonic() error {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return err
	}

	fmt.Println("entropy:", hexutil.Encode(entropy))

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return err
	}

	fmt.Println("mnemonic:", mnemonic)

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	fmt.Println("Master private key:", masterKey)
	fmt.Println("Master public key:", publicKey)

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	// Ethereum default path is m/44'/60'/0'/0/0
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return err
	}

	fmt.Println("Account1:", account.Address.Hex()) // account 1

	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	account, err = wallet.Derive(path, false)
	if err != nil {
		return err
	}

	fmt.Println("Account2:", account.Address.Hex()) // account 2

	return nil
}

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

	// https://github.com/ethereum/go-ethereum/issues/21221
	amount, _ := new(big.Int).SetString("2775176240359883548376", 10)
	fmt.Println("amount(wei):", amount)
	fmt.Println("amount(eth):", util.WeiToEther(amount))

	fmt.Println("Discount 30%")
	dAmountInWei := new(big.Int)
	dAmountInWei.Mul(amount, big.NewInt(70))
	dAmountInWei.Div(dAmountInWei, big.NewInt(100))
	fmt.Println("discount amount(wei):", dAmountInWei)
	fmt.Println("discount amount(eth):", util.WeiToEther(dAmountInWei))

	dString := new(big.Float)
	dString.SetString(dAmountInWei.String())
	fmt.Println("discount set string", dString)
	dInt := new(big.Float)
	dInt.SetInt(dAmountInWei)
	fmt.Println("discount set int (more precise)", dInt)

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
