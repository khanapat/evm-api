package util

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// https://stackoverflow.com/questions/50772811/how-can-i-get-the-same-return-value-as-solidity-abi-encodepacked-in-golang

func AbiEncodeUint() (string, error) {
	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return "", err
	}
	fmt.Println("type:", uint256Ty.GetType(), uint256Ty.String())

	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}
	fmt.Println("arguments:", arguments)

	bytes, err := arguments.Pack(big.NewInt(42))
	if err != nil {
		return "", err
	}
	data := hexutil.Encode(bytes)
	fmt.Println("encoded data:", data)
	// ethers.utils.defaultAbiCoder.encode + keccak256
	fmt.Println("keccak256 encoded data:", crypto.Keccak256Hash(bytes)) // hexutil.Encode(crypto.Keccak256(bytes))

	return data, nil
}

func AbiDecodeUint(encode string) error {
	data, err := hexutil.Decode(encode)
	if err != nil {
		return err
	}
	fmt.Println("data:", data)

	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return err
	}
	fmt.Println("type:", uint256Ty.GetType(), uint256Ty.String())

	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}
	fmt.Println("arguments:", arguments)

	decode, err := arguments.Unpack(data)
	if err != nil {
		return err
	}
	fmt.Println("decoded data:", decode)

	return nil
}

func AbiEncodeString() (string, error) {
	stringTy, err := abi.NewType("string", "string", nil)
	if err != nil {
		return "", err
	}

	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}

	bytes, err := arguments.Pack("key")
	if err != nil {
		return "", err
	}
	data := hexutil.Encode(bytes)
	fmt.Println("encoded data:", data)
	fmt.Println("keccak256 encoded data:", crypto.Keccak256Hash(bytes)) // hexutil.Encode(crypto.Keccak256(bytes))

	return data, nil
}

func AbiDecodeString(encode string) error {
	data, err := hexutil.Decode(encode)
	if err != nil {
		return err
	}
	stringTy, err := abi.NewType("string", "string", nil)
	if err != nil {
		return err
	}
	arguments := abi.Arguments{
		{
			Type: stringTy,
		},
	}
	decode, err := arguments.Unpack(data)
	if err != nil {
		return err
	}
	fmt.Println("decoded data:", decode)
	return nil
}

func AbiEncodeTest() (string, error) {
	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return "", err
	}
	addressTy, err := abi.NewType("address", "address", nil)
	if err != nil {
		return "", err
	}

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
	}
	fmt.Println("arguments:", arguments)

	bytes, err := arguments.Pack(
		common.HexToAddress("0x7b2b2C058722F8c60eB9287D51C6716f04E442AE"),
		common.HexToAddress("0xA67FCD5Aef83732D93D7b44ce8564de7340c72ec"),
		big.NewInt(1),
	)
	if err != nil {
		return "", err
	}
	data := hexutil.Encode(bytes)
	fmt.Println("encoded data:", data)
	// ethers.utils.defaultAbiCoder.encode + keccak256
	fmt.Println("keccak256 encoded data:", crypto.Keccak256Hash(bytes)) // hexutil.Encode(crypto.Keccak256(bytes))

	return data, nil
}
