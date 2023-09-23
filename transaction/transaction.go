package transaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"evm-api/util"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"golang.org/x/crypto/sha3"
)

func QueryBlocks(client *ethclient.Client, number int64) error {
	// pattern 1
	header, err := client.HeaderByNumber(context.Background(), nil) // nil for latest
	if err != nil {
		return err
	}

	fmt.Println("block numer:", header.Number.String())

	// pattern 2
	blockNumber := big.NewInt(number)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return err
	}

	fmt.Println("block number:", block.Number().Uint64())         // 5671744
	fmt.Println("block timestamp:", block.Time())                 // 1527211625
	fmt.Println("block difficulty:", block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println("block hash:", block.Hash().Hex())                // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9	fmt.Println("number of txn in block:", len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		return nil
	}

	fmt.Println("number of txn in block:", count)

	return nil
}

func QueryTxn(client *ethclient.Client, hash string) error {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Chain ID:", chainID)

	txHash := common.HexToHash(hash)
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return err
	}

	if !isPending {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return err
		}
		fmt.Println("Gas Used:", receipt.GasUsed)
		fmt.Println("Cumulative Gas Used:", receipt.CumulativeGasUsed)
		fmt.Println("Status:", receipt.Status) // 1 (success) | 0 (fail)
		fmt.Println("BlockNumber:", receipt.BlockNumber)
		fmt.Println("BlockHash: ", receipt.BlockHash)
		fmt.Println("ContractAddress: ", receipt.ContractAddress)
		fmt.Println("TxHash: ", receipt.TxHash)
		fmt.Println("Type: ", receipt.Type) // type-0 (legacy) | type-2 (EIP-1559)
		fmt.Println("TxnIndex: ", receipt.TransactionIndex)
		fmt.Println("PostState: ", receipt.PostState)
		fmt.Println("Logs: ", receipt.Logs)

		fmt.Println("--------------")

		fmt.Println("Max Gas Limit:", tx.Gas())                               // max gas limit
		fmt.Println("Max Gas Fee in wei:", tx.GasPrice().Uint64())            // max gas fee
		fmt.Println("Max Gas Fee in gwei:", util.ToDecimal(tx.GasPrice(), 9)) // in eth use 18
		fmt.Println("Max Txn Fee:", tx.Cost())                                // max transaction fee
		fmt.Println("Nonce:", tx.Nonce())                                     // nonce
		fmt.Println("To:", tx.To().Hex())                                     // to address
		v, r, s := tx.RawSignatureValues()                                    // v, r, s
		fmt.Println("RawSignatureValue: ", v.String(), r.String(), s.String())
		a, _ := tx.MarshalJSON()
		fmt.Println("Tx marshal: ", string(a))

		block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
		if err != nil {
			return err
		}
		fmt.Println("Usage by Txn:", receipt.GasUsed)                        // gas limit
		fmt.Println("Base fee in wei:", block.BaseFee())                     // base fee in each block
		fmt.Println("Base fee in gwei:", util.ToDecimal(block.BaseFee(), 9)) // in eth use 18
		// fmt.Println("base fee:", misc.CalcBaseFee(params.GoerliChainConfig, block.Header())) // read at readme.md (previous block number)
		fmt.Println("Gas Tip in wei:", tx.GasTipCap())                     // max priority
		fmt.Println("Gas Tip in gwei:", util.ToDecimal(tx.GasTipCap(), 9)) // in eth use 18

		gasPrice := big.NewInt(0).Add(block.BaseFee(), tx.GasTipCap())
		if gasPrice.Cmp(tx.GasPrice()) == 1 { // if base fee + tip > max gas fee, use max gas fee
			gasPrice = tx.GasPrice()
		}
		fmt.Println("Gas Price in wei:", gasPrice)                     // base fee + tip or max gas fee
		fmt.Println("Gas Price in gwei:", util.ToDecimal(gasPrice, 9)) // in eth use 18

		fmt.Println("Transaction Fee in wei: ", util.CalcGasCost(receipt.GasUsed, gasPrice))                     // usage by txn * gas price
		fmt.Println("Transaction Fee in eth: ", util.ToDecimal(util.CalcGasCost(receipt.GasUsed, gasPrice), 18)) // in gwei use 9

		fmt.Println("Burnt in wei:", util.CalcGasCost(receipt.GasUsed, block.BaseFee()))                     // usage by txn * base fee
		fmt.Println("Burnt in eth:", util.ToDecimal(util.CalcGasCost(receipt.GasUsed, block.BaseFee()), 18)) // in eth use 18

		saving := big.NewInt(0).Sub(tx.GasPrice(), gasPrice)
		fmt.Println("Txn Savings in wei:", util.CalcGasCost(receipt.GasUsed, saving))                     // usage by txn * (max fee - (base fee + tip))
		fmt.Println("Txn Savings in eth:", util.ToDecimal(util.CalcGasCost(receipt.GasUsed, saving), 18)) // in eth use 9

		fmt.Println("--------------")

		// in go-ethereum 1.11 remove this function
		// msg, err := tx.AsMessage(types.LatestSignerForChainID(chainID), nil)
		// if err != nil {
		// 	return err
		// }
		// fmt.Println("From: ", msg.From().Hex())     // get from address
		// fmt.Println("GasFeeCap: ", msg.GasFeeCap()) // max gas fee
		// fmt.Println("GasPrice: ", msg.GasPrice())   // max gas fee
		// fmt.Println("GasTipCap: ", msg.GasTipCap()) // gas tip
		// fmt.Println("Gas: ", msg.Gas())
		// fmt.Println("Nonce: ", msg.Nonce())
		// fmt.Println("AccessList: ", msg.AccessList())
		// fmt.Println("Value:", msg.Value())
	}

	return nil
}

func SendETH(client *ethclient.Client) error {
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

	value := big.NewInt(1000000000) // in wei (1 geth)
	gasLimit := uint64(21000)       // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("GasPrice:", gasPrice)
	toAddress := common.HexToAddress("0xa9B6D99bA92D7d691c6EF4f49A1DC909822Cee46")

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	return nil
}

func SendERC20(client *ethclient.Client) error {
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

	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("GasPrice:", gasPrice)
	toAddress := common.HexToAddress("0xa9B6D99bA92D7d691c6EF4f49A1DC909822Cee46")
	tokenAddress := common.HexToAddress("0x9846a60C6ab6733b599C84062Dc92254A836Ab41")

	// signature transfer(address,uint256)
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println("Method - transfer(address,uint256)", hexutil.Encode(methodID))

	// padding address
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println("ToAddress:", hexutil.Encode(paddedAddress))

	// padding amount
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println("Amount:", hexutil.Encode(paddedAmount))

	// data
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	fmt.Println("Data:", hexutil.Encode(data))

	// estimate gas (if have problem with gas fee, add this)
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		return err
	}
	fmt.Println("GasLimit:", gasLimit+10000)

	tx := types.NewTransaction(nonce, tokenAddress, value, 60000, gasPrice.Add(gasPrice, big.NewInt(10000000000)), data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	ts := types.Transactions{signedTx}
	b := new(bytes.Buffer)
	ts.EncodeIndex(0, b)
	rawTxBytes := b.Bytes()
	rawTxHex := hex.EncodeToString(rawTxBytes)
	fmt.Println("Raw Txn:", rawTxHex)

	// if err := client.SendTransaction(context.Background(), signedTx); err != nil {
	// 	return err
	// }
	// fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	return nil
}

func RawTransaction(client *ethclient.Client) error {
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

	value := big.NewInt(1000000000) // in wei (1 geth)
	gasLimit := uint64(21000)       // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("GasPrice:", gasPrice)
	toAddress := common.HexToAddress("0xa9B6D99bA92D7d691c6EF4f49A1DC909822Cee46")

	// legacy transaction
	legacyTx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	legacySignedTx, err := types.SignTx(legacyTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	fmt.Printf("legacy tx sent: %s\n", legacySignedTx.Hash().Hex())
	c, _ := legacySignedTx.MarshalJSON()
	fmt.Println("legacy c0:", string(c))

	d, _ := legacySignedTx.MarshalBinary()
	fmt.Println("legacy d2:", hexutil.Encode(d))
	// or
	// ts := types.Transactions{legacySignedTx}
	// b := new(bytes.Buffer)
	// ts.EncodeIndex(0, b)
	// rawTxBytes := b.Bytes()
	// rawTxHex := hex.EncodeToString(rawTxBytes)
	// fmt.Println("legacy d2:", rawTxHex)

	// dynamic fee transaction
	tip := big.NewInt(2000000000) // 2 gwei

	dynamicTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: tip,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      nil,
	})
	dynamicSignedTx, err := types.SignTx(dynamicTx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return err
	}

	fmt.Printf("dynamic tx sent: %s\n", dynamicSignedTx.Hash().Hex())
	e, _ := dynamicSignedTx.MarshalJSON()
	fmt.Println("dynamic e0:", string(e))

	f, _ := dynamicSignedTx.MarshalBinary()
	fmt.Println("dynamic f2:", hexutil.Encode(f))

	return nil
}
