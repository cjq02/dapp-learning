// Code generated - DO NOT EDIT.
// This file is a generated binding and any modifications will be lost.

package store

import (
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// StoreABI is the input ABI used to generate the binding from.
const StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

// Store is an auto generated Go binding around an Ethereum contract.
type Store struct {
	StoreCaller     // Read-only binding to the contract
	StoreTransactor // Write-only binding to the contract
	StoreFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreCaller struct {
	contract *bind.BoundContract
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreTransactor struct {
	contract *bind.BoundContract
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreFilterer struct {
	contract *bind.BoundContract
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreSession struct {
	Contract     *Store
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreCallerSession struct {
	Contract *StoreCaller
	CallOpts bind.CallOpts
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreTransactorSession struct {
	Contract     *StoreTransactor
	TransactOpts bind.TransactOpts
}

// StoreItemSet represents a ItemSet event raised by the Store contract.
type StoreItemSet struct {
	Key   [32]byte
	Value [32]byte
	Raw   types.Log
}

// StoreItemSetIterator is returned from FilterItemSet and is used to iterate over the raw logs.
type StoreItemSetIterator struct {
	event *event.Event
}

// Next advances the iterator to the subsequent event.
func (it *StoreItemSetIterator) Next() bool {
	return it.event.Next()
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StoreItemSetIterator) Error() error {
	return it.event.Err()
}

// Close terminates the iteration process.
func (it *StoreItemSetIterator) Close() error {
	it.event.Close()
	return nil
}

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func NewStore(address common.Address, backend bind.ContractBackend) (*Store, error) {
	contract, err := bindStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Store{StoreCaller: StoreCaller{contract: contract}, StoreTransactor: StoreTransactor{contract: contract}, StoreFilterer: StoreFilterer{contract: contract}}, nil
}

// bindStore binds a generic wrapper to an already deployed contract.
func bindStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Version is a free data retrieval call binding the contract method.
// Solidity: function version() view returns(string)
func (_Store *StoreCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "version")
	return out[0].(string), err
}

// Version is a free data retrieval call binding the contract method.
// Solidity: function version() view returns(string)
func (_Store *StoreSession) Version() (string, error) {
	return _Store.Contract.Version(&_Store.CallOpts)
}

// GetItem is a free data retrieval call binding the contract method.
// Solidity: function items(bytes32 ) view returns(bytes32)
func (_Store *StoreCaller) GetItem(opts *bind.CallOpts, key [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "items", key)
	return out[0].([32]byte), err
}

// GetItem is a free data retrieval call binding the contract method.
// Solidity: function items(bytes32 ) view returns(bytes32)
func (_Store *StoreSession) GetItem(key [32]byte) ([32]byte, error) {
	return _Store.Contract.GetItem(&_Store.CallOpts, key)
}

// SetItem is a paid mutator transaction binding the contract method.
// Solidity: function setItem(bytes32 key, bytes32 value) returns()
func (_Store *StoreTransactor) SetItem(opts *bind.TransactOpts, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "setItem", key, value)
}

// SetItem is a paid mutator transaction binding the contract method.
// Solidity: function setItem(bytes32 key, bytes32 value) returns()
func (_Store *StoreSession) SetItem(key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _Store.Contract.SetItem(&_Store.TransactOpts, key, value)
}

// FilterItemSet is a free log retrieval operation binding the contract event.
// Solidity: event ItemSet(bytes32 indexed key, bytes32 value)
func (_Store *StoreFilterer) FilterItemSet(opts *bind.FilterOpts, key [][32]byte) (*StoreItemSetIterator, error) {
	var keyRule []interface{}
	if len(key) == 0 {
		keyRule = []interface{}{new([32]byte)}
	} else {
		keyRule = make([]interface{}, len(key))
		for i, k := range key {
			keyRule[i] = k
		}
	}

	_, sub, err := _Store.contract.FilterLogs(opts, "ItemSet", keyRule)
	if err != nil {
		return nil, err
	}
	return &StoreItemSetIterator{event: sub}, nil
}

// WatchItemSet is a free log subscription operation binding the contract event.
// Solidity: event ItemSet(bytes32 indexed key, bytes32 value)
func (_Store *StoreFilterer) WatchItemSet(opts *bind.WatchOpts, sink chan<- *StoreItemSet, key [][32]byte) (event.Subscription, error) {
	var keyRule []interface{}
	if len(key) == 0 {
		keyRule = []interface{}{new([32]byte)}
	} else {
		keyRule = make([]interface{}, len(key))
		for i, k := range key {
			keyRule[i] = k
		}
	}

	logs, sub, err := _Store.contract.WatchLogs(opts, "ItemSet", keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				event := new(StoreItemSet)
				if err := _Store.contract.UnpackLog(event, "ItemSet", log); err != nil {
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

// ParseItemSet is a log parse operation binding the contract event.
// Solidity: event ItemSet(bytes32 indexed key, bytes32 value)
func (_Store *StoreFilterer) ParseItemSet(log types.Log) (*StoreItemSet, error) {
	event := new(StoreItemSet)
	if err := _Store.contract.UnpackLog(event, "ItemSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
