// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package getset

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// GetsetMetaData contains all meta data concerning the Getset contract.
var GetsetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"SetA\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"a\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_a\",\"type\":\"uint256\"}],\"name\":\"setA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060eb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630dbe671f146037578063ee919d50146051575b600080fd5b603f60005481565b60405190815260200160405180910390f35b6060605c366004609d565b6062565b005b60008190556040518181527f577ef53da0f5052c44a6d7f830d9245d67a4d07bcb4967c6c4a6e10d69ecc0429060200160405180910390a150565b60006020828403121560ae57600080fd5b503591905056fea2646970667358221220b664282008bb59c94ba723b755ac5829d75307b64de26e78343eba212c9af52e64736f6c634300080d0033",
}

// GetsetABI is the input ABI used to generate the binding from.
// Deprecated: Use GetsetMetaData.ABI instead.
var GetsetABI = GetsetMetaData.ABI

// GetsetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GetsetMetaData.Bin instead.
var GetsetBin = GetsetMetaData.Bin

// DeployGetset deploys a new Ethereum contract, binding an instance of Getset to it.
func DeployGetset(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Getset, error) {
	parsed, err := GetsetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GetsetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Getset{GetsetCaller: GetsetCaller{contract: contract}, GetsetTransactor: GetsetTransactor{contract: contract}, GetsetFilterer: GetsetFilterer{contract: contract}}, nil
}

// Getset is an auto generated Go binding around an Ethereum contract.
type Getset struct {
	GetsetCaller     // Read-only binding to the contract
	GetsetTransactor // Write-only binding to the contract
	GetsetFilterer   // Log filterer for contract events
}

// GetsetCaller is an auto generated read-only Go binding around an Ethereum contract.
type GetsetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GetsetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GetsetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GetsetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GetsetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GetsetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GetsetSession struct {
	Contract     *Getset           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GetsetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GetsetCallerSession struct {
	Contract *GetsetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GetsetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GetsetTransactorSession struct {
	Contract     *GetsetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GetsetRaw is an auto generated low-level Go binding around an Ethereum contract.
type GetsetRaw struct {
	Contract *Getset // Generic contract binding to access the raw methods on
}

// GetsetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GetsetCallerRaw struct {
	Contract *GetsetCaller // Generic read-only contract binding to access the raw methods on
}

// GetsetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GetsetTransactorRaw struct {
	Contract *GetsetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGetset creates a new instance of Getset, bound to a specific deployed contract.
func NewGetset(address common.Address, backend bind.ContractBackend) (*Getset, error) {
	contract, err := bindGetset(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Getset{GetsetCaller: GetsetCaller{contract: contract}, GetsetTransactor: GetsetTransactor{contract: contract}, GetsetFilterer: GetsetFilterer{contract: contract}}, nil
}

// NewGetsetCaller creates a new read-only instance of Getset, bound to a specific deployed contract.
func NewGetsetCaller(address common.Address, caller bind.ContractCaller) (*GetsetCaller, error) {
	contract, err := bindGetset(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GetsetCaller{contract: contract}, nil
}

// NewGetsetTransactor creates a new write-only instance of Getset, bound to a specific deployed contract.
func NewGetsetTransactor(address common.Address, transactor bind.ContractTransactor) (*GetsetTransactor, error) {
	contract, err := bindGetset(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GetsetTransactor{contract: contract}, nil
}

// NewGetsetFilterer creates a new log filterer instance of Getset, bound to a specific deployed contract.
func NewGetsetFilterer(address common.Address, filterer bind.ContractFilterer) (*GetsetFilterer, error) {
	contract, err := bindGetset(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GetsetFilterer{contract: contract}, nil
}

// bindGetset binds a generic wrapper to an already deployed contract.
func bindGetset(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GetsetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Getset *GetsetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Getset.Contract.GetsetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Getset *GetsetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Getset.Contract.GetsetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Getset *GetsetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Getset.Contract.GetsetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Getset *GetsetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Getset.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Getset *GetsetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Getset.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Getset *GetsetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Getset.Contract.contract.Transact(opts, method, params...)
}

// A is a free data retrieval call binding the contract method 0x0dbe671f.
//
// Solidity: function a() view returns(uint256)
func (_Getset *GetsetCaller) A(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Getset.contract.Call(opts, &out, "a")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// A is a free data retrieval call binding the contract method 0x0dbe671f.
//
// Solidity: function a() view returns(uint256)
func (_Getset *GetsetSession) A() (*big.Int, error) {
	return _Getset.Contract.A(&_Getset.CallOpts)
}

// A is a free data retrieval call binding the contract method 0x0dbe671f.
//
// Solidity: function a() view returns(uint256)
func (_Getset *GetsetCallerSession) A() (*big.Int, error) {
	return _Getset.Contract.A(&_Getset.CallOpts)
}

// SetA is a paid mutator transaction binding the contract method 0xee919d50.
//
// Solidity: function setA(uint256 _a) returns()
func (_Getset *GetsetTransactor) SetA(opts *bind.TransactOpts, _a *big.Int) (*types.Transaction, error) {
	return _Getset.contract.Transact(opts, "setA", _a)
}

// SetA is a paid mutator transaction binding the contract method 0xee919d50.
//
// Solidity: function setA(uint256 _a) returns()
func (_Getset *GetsetSession) SetA(_a *big.Int) (*types.Transaction, error) {
	return _Getset.Contract.SetA(&_Getset.TransactOpts, _a)
}

// SetA is a paid mutator transaction binding the contract method 0xee919d50.
//
// Solidity: function setA(uint256 _a) returns()
func (_Getset *GetsetTransactorSession) SetA(_a *big.Int) (*types.Transaction, error) {
	return _Getset.Contract.SetA(&_Getset.TransactOpts, _a)
}

// GetsetSetAIterator is returned from FilterSetA and is used to iterate over the raw logs and unpacked data for SetA events raised by the Getset contract.
type GetsetSetAIterator struct {
	Event *GetsetSetA // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GetsetSetAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GetsetSetA)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GetsetSetA)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GetsetSetAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GetsetSetAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GetsetSetA represents a SetA event raised by the Getset contract.
type GetsetSetA struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSetA is a free log retrieval operation binding the contract event 0x577ef53da0f5052c44a6d7f830d9245d67a4d07bcb4967c6c4a6e10d69ecc042.
//
// Solidity: event SetA(uint256 a)
func (_Getset *GetsetFilterer) FilterSetA(opts *bind.FilterOpts) (*GetsetSetAIterator, error) {

	logs, sub, err := _Getset.contract.FilterLogs(opts, "SetA")
	if err != nil {
		return nil, err
	}
	return &GetsetSetAIterator{contract: _Getset.contract, event: "SetA", logs: logs, sub: sub}, nil
}

// WatchSetA is a free log subscription operation binding the contract event 0x577ef53da0f5052c44a6d7f830d9245d67a4d07bcb4967c6c4a6e10d69ecc042.
//
// Solidity: event SetA(uint256 a)
func (_Getset *GetsetFilterer) WatchSetA(opts *bind.WatchOpts, sink chan<- *GetsetSetA) (event.Subscription, error) {

	logs, sub, err := _Getset.contract.WatchLogs(opts, "SetA")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GetsetSetA)
				if err := _Getset.contract.UnpackLog(event, "SetA", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetA is a log parse operation binding the contract event 0x577ef53da0f5052c44a6d7f830d9245d67a4d07bcb4967c6c4a6e10d69ecc042.
//
// Solidity: event SetA(uint256 a)
func (_Getset *GetsetFilterer) ParseSetA(log types.Log) (*GetsetSetA, error) {
	event := new(GetsetSetA)
	if err := _Getset.contract.UnpackLog(event, "SetA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
