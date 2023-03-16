package signature

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"evm-api/util"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateSignature(data []byte) (string, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return "", err
	}

	hash := crypto.Keccak256Hash(data) // keccak256 + hash
	fmt.Println("Hash:", hash.Hex())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	signature[64] += 27 // https://github.com/ethereum/go-ethereum/issues/19751#issuecomment-504900739

	fmt.Println("signature:", hexutil.Encode(signature))

	r, s, v := util.SigRSV(signature)
	fmt.Println("r", hexutil.Encode(r[:])[2:])
	fmt.Println("s", hexutil.Encode(s[:])[2:])
	fmt.Println("v", v)

	return hexutil.Encode(signature), nil
}

func GenerateSignatureWihEIP191(message string) (string, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return "", err
	}

	// hash, fullMessage := accounts.TextAndHash([]byte(message))
	// fmt.Println("Hash:", hash)
	// fmt.Println("FullMessage:", fullMessage)
	// or
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage)) // keccak256 + hash
	fmt.Println("Hash:", hash.Hex())
	fmt.Println("FullMessage:", fullMessage)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	signature[64] += 27 // https://github.com/ethereum/go-ethereum/issues/19751#issuecomment-504900739

	fmt.Println("signature:", hexutil.Encode(signature))

	r, s, v := util.SigRSV(signature)
	fmt.Println("r", hexutil.Encode(r[:])[2:])
	fmt.Println("s", hexutil.Encode(s[:])[2:])
	fmt.Println("v", v)

	return hexutil.Encode(signature), nil
}

func VerifySignature(data []byte, signature string) error {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	hash := crypto.Keccak256Hash(data)
	fmt.Println("Hash:", hash.Hex())

	signatureByte, err := hexutil.Decode(signature)
	if err != nil {
		fmt.Println("decode error")
		return err
	}

	signatureByte[64] -= 27

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signatureByte)
	if err != nil {
		fmt.Println("ecrecover error")
		return err
	}

	// pattern 1
	matches1 := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println("matches pattern1:", matches1)

	// pattern 2
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signatureByte)
	if err != nil {
		return err
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches2 := bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println("matches pattern2:", matches2)

	// pattern 3
	signatureNoRecoverID := signatureByte[:len(signatureByte)-1] // remove recovery ID
	matches3 := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println("matches pattern3:", matches3)

	return nil
}

func VerifySignatureWithEIP191(message string, signature string) error {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// hash, fullMessage := accounts.TextAndHash([]byte(message))
	// fmt.Println("Hash:", hash)
	// fmt.Println("FullMessage:", fullMessage)
	// or
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage)) // keccak256 + hash
	fmt.Println("Hash:", hash.Hex())

	signatureByte, err := hexutil.Decode(signature)
	if err != nil {
		fmt.Println("decode error")
		return err
	}

	signatureByte[64] -= 27

	// pattern 1
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signatureByte)
	if err != nil {
		fmt.Println("ecrecover error")
		return err
	}

	matches1 := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println("matches pattern1:", matches1)

	// pattern 2
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signatureByte)
	if err != nil {
		return err
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches2 := bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println("matches pattern2:", matches2)

	// pattern 3
	signatureNoRecoverID := signatureByte[:len(signatureByte)-1] // remove recovery ID
	matches3 := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println("matches pattern3:", matches3)

	return nil
}
