package signature

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func GenerateEIP712() ([]byte, string, error) {
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{ // order is significated
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Person": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
		},
		PrimaryType: "Person",
		Domain: apitypes.TypedDataDomain{
			Name:              "Test EIP712",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(80001),
			VerifyingContract: "0x962D47852d0a92B9fec46ce3Afea22eA615068a3",
		},
		Message: apitypes.TypedDataMessage{
			"name":   "bobo",
			"wallet": "0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C",
		},
	}

	hash2, rawData2, _ := apitypes.TypedDataAndHash(typedData)
	fmt.Println("Raw data:", hexutil.Encode([]byte(rawData2)))
	fmt.Println("Hash:", common.BytesToHash(hash2))
	// or
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Domain separator:", domainSeparator.String())

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Typed data hash:", typedDataHash.String())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	fmt.Println("Raw data:", hexutil.Encode(rawData))

	hash := common.BytesToHash(crypto.Keccak256(rawData))
	fmt.Println("Hash:", hash)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, "", err
	}

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, "", err
	}

	fmt.Println("signature formal V 0/1:", hexutil.Encode(signature))

	signature[64] += 27

	fmt.Println("signature:", hexutil.Encode(signature))

	return hash.Bytes(), hexutil.Encode(signature), nil
}

func VerifyEIP712(hash []byte, signature string) error {
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

	signatureByte, err := hexutil.Decode(signature)
	if err != nil {
		fmt.Println("decode error")
		return err
	}

	sig := make([]byte, len(signatureByte))
	copy(sig, signatureByte)
	if len(sig) != 65 {
		return errors.New(fmt.Sprintf("invalid length of signature: %d", len(sig)))
	}
	if sig[64] != 27 && sig[64] != 28 && sig[64] != 1 && sig[64] != 0 {
		return errors.New("invalid signature type")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// pattern 1
	sigPublicKey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return err
	}
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		return err
	}

	fmt.Println("address:", crypto.PubkeyToAddress(*pubKey))

	matches1 := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println("matches pattern1:", matches1)

	// pattern 2
	sigPublicKeyECDSA, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return err
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches2 := bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println("matches pattern2:", matches2)

	// pattern 3
	signatureNoRecoverID := sig[:len(sig)-1] // remove recovery ID
	matches3 := crypto.VerifySignature(publicKeyBytes, hash, signatureNoRecoverID)
	fmt.Println("matches pattern3:", matches3)

	return nil
}

func Test() error {

	// data := apitypes.TypedData{
	// 	Types: apitypes.Types{
	// 		"NFT": []apitypes.Type{
	// 			{Name: "cId", Type: "uint256"},
	// 			{Name: "nft", Type: "address"},
	// 			{Name: "signId", Type: "bytes32"},
	// 			{Name: "mintTo", Type: "address"},
	// 		},
	// 	},
	// 	PrimaryType: "NFT",
	// 	Domain: apitypes.TypedDataDomain{
	// 		Name:              "Aster Station",
	// 		Version:           "1.0.0",
	// 		ChainId:           math.NewHexOrDecimal256(4833),
	// 		VerifyingContract: "0xEA0BD2c194944C719b255C48F327d177702b4A34",
	// 	},
	// 	Message: apitypes.TypedDataMessage{
	// 		"cId":    1,
	// 		"nft":    "0xA67FCD5Aef83732D93D7b44ce8564de7340c72ec",
	// 		"signId": "1",
	// 		"mintTo": "0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C",
	// 	},
	// }

	return nil
}
