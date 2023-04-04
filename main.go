package main

import (
	"context"
	"evm-api/account"
	"evm-api/signature"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	https, err := ethclient.DialContext(ctx, os.Getenv("HTTPS_RPC_NETWORK"))
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := account.Balance(https, "0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"); err != nil {
		log.Fatal(err.Error())
	}

	// util.Encode()

	// encode, err := util.AbiEncodeUint()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// if err := util.AbiDecodeUint(encode); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// encode, err := util.AbiEncodeTest()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// fmt.Println(encode)

	// if err := account.GenerateAccount(); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := account.GenerateKeystore("secret"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := account.ImportKeystore("./wallets/UTC--2022-11-04T10-36-03.153579000Z--8f4c26dee2be4d391712a3971f440b11e20fe27a", "secret"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// 0xD1F4FEFc66BECC13cEdAcF61B32607BE8a239827 is contract
	// if err := account.CheckIsAddress(https, "0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := transaction.QueryBlocks(https, 7908442); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// example
	// 0x22fc5e3c57b451dd61621ee2d07e5ae2e25afe2360726edfab0e97c98de0e1e0
	// 0x93ed284eaf3e0ae96b83c4dcb3a2a2b45acd8f6c638a1a354f272ff8d6f90fe4
	// 0x5d3b42b79c37d7abee044cd47ae0e51e7e28d5d753c6dd07458f812894201b23
	// 0x04fa6f2064697037717cfe7983d525ca6b13884492ed4f05c122b7a3d98ed747
	// if err := transaction.QueryTxn(https, "0x22fc5e3c57b451dd61621ee2d07e5ae2e25afe2360726edfab0e97c98de0e1e0"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := transaction.SendETH(https); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := transaction.SendERC20(https); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// subscription.SubscriptionBlock()

	// if err := smartcontract.Deploy(https); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// instance, err := smartcontract.Load(https, "0x188992471F03D5bf1EaC66973fc3E7CA7ee5C0D3")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := smartcontract.SetStore(https, instance); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := smartcontract.GetBytecode(https, "0x188992471F03D5bf1EaC66973fc3E7CA7ee5C0D3"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := subscription.SubscriptionEvent(); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := subscription.ReadEvent(https); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// generate & verify signature
	// 	sign, err := signature.GenerateSignature([]byte("bobo"))
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	if err := signature.VerifySignature([]byte("bobo"), sign); err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	signWithEIP191, err := signature.GenerateSignatureWihEIP191("bobo")
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	if err := signature.VerifySignatureWithEIP191("bobo", signWithEIP191); err != nil {
	// 		log.Fatal(err.Error())
	// 	}

	// 	// use third-party lib to verify
	// 	// standard
	// 	isValid, err := sigverify.VerifyEllipticCurveHexSignatureEx(common.HexToAddress("0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"), []byte("bobo"), signWithEIP191) // add prefix EIP-191 inside function
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// 	fmt.Println("standard:", isValid)

	//	// eip-712
	//	const ExampleTypedData = `
	//
	//	{
	//	    "types": {
	//	        "EIP712Domain": [
	//	            {
	//	                "name": "name",
	//	                "type": "string"
	//	            },
	//				{
	//					"name": "version",
	//					"type": "string"
	//				},
	//	            {
	//	                "name": "chainId",
	//	                "type": "uint256"
	//	            },
	//				{
	//					"name": "verifyingContract",
	//					"type": "address"
	//				}
	//	        ],
	//	        "Person": [
	//	            {
	//	                "name": "name",
	//	                "type": "string"
	//	            },
	//	            {
	//	                "name": "wallet",
	//	                "type": "address"
	//	            }
	//	        ]
	//	    },
	//	    "domain": {
	//	        "name": "Test EIP712",
	//			"version": "1",
	//	        "chainId": "80001",
	//			"verifyingContract": "0xd8b934580fcE35a11B58C6D73aDeE468a2833fa8"
	//	    },
	//	    "primaryType": "Person",
	//	    "message": {
	//	        "name": "bobo",
	//	        "wallet": "0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"
	//	    }
	//	}
	//
	// `
	//
	//	var typedData apitypes.TypedData
	//	if err := json.Unmarshal([]byte(ExampleTypedData), &typedData); err != nil {
	//		log.Fatal(err.Error())
	//	}
	//	valid, err := sigverify.VerifyTypedDataHexSignatureEx(
	//		common.HexToAddress("0xc083EB69aa7215f4AFa7a22dcbfCC1a33999371C"),
	//		typedData,
	//		"0x60db464f22d6677d5bfaede105cab6eaaecd33a63197f9cdedc86052e537774055b2525992821b5f69aca000eaabf90f622cd0d7079168ce9b60034fff36491e1b",
	//	)
	//	if err != nil {
	//		log.Fatal(err.Error())
	//	}
	//	fmt.Println("eip-712:", valid)

	hash, sign, err := signature.GenerateEIP712()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := signature.VerifyEIP712(hash, sign); err != nil {
		log.Fatal(err.Error())
	}

	// if err := signature.Test(); err != nil {
	// 	log.Fatal(err.Error())
	// }
}
