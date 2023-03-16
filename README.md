# Go Ethereum

network: Goerli test network

store contract address: 0x188992471F03D5bf1EaC66973fc3E7CA7ee5C0D3

erc20 contract address: 0x9846a60C6ab6733b599C84062Dc92254A836Ab41

docs:

-   https://ethereum.org/en/developers/docs/programming-languages/golang/
-   https://goethereumbook.org/en/

## get base fee

-   https://ethereum.stackexchange.com/questions/107814/getting-current-base-fee-from-json-rpc

## cumulative gas used & gas used by txn

-   https://ethereum.stackexchange.com/questions/3346/what-is-and-how-to-calculate-cumulative-gas-used
-   https://ethereum.stackexchange.com/questions/40172/whats-difference-between-gas-used-by-txn-and-cumulative-gas-used-fields

## compile solidity

-   https://github.com/crytic/solc-select

```bash
# check cli
solc --version

# switch solc version
solc-select versions # list version
solc-select install 0.8.13 # install version
solc-select use 0.8.13 # use version

# generate abi
solc --optimize --abi --bin ./artifacts/Store.sol -o ./artifacts/build

# convert abi to go file
export PATH=$(go env GOPATH)/bin:$PATH # check abigen binary in $GOPATH/bin

abigen --bin=./artifacts/build/Store.bin --abi=./artifacts/build/Store.abi --pkg=store --out=./contract/store.go
```

## sign message

-   https://github.com/etaaa/Golang-Ethereum-Personal-Sign
-   https://github.com/storyicon/sigverify

## send txn VS send raw txn

-   https://stackoverflow.com/questions/50985798/difference-between-sendtransaction-and-sendrawtransaction-in-web3-py
