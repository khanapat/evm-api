package util

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Encode() {
	key := [32]byte{} // bytes32
	copy(key[:], []byte("key"))

	key2 := []byte("key") // bytes

	fmt.Println("key bytes32:", string(key[:])) // key
	fmt.Println("key bytes:", string(key2))     // key

	// hexutil do like hex with prefix 0x
	fmt.Println("key bytes32 encode:", hexutil.Encode(key[:])) // 0x6b65790000000000000000000000000000000000000000000000000000000000
	fmt.Println("key bytes encode:", hexutil.Encode(key2))     // 0x6b6579

	// hexutil convert
	u := uint64(123)
	fmt.Println("uint64 to hex:", u, hexutil.Uint64(u)) // uint64 to hex
	i := big.NewInt(100000)
	fmt.Println("big.Int to hex:", i, (*hexutil.Big)(i)) // big.Int to hex
	b := []byte{}
	fmt.Println("bytes to hex:", b, hexutil.Bytes(b)) // bytes to hex
}
