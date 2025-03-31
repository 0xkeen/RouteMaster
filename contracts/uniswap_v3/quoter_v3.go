// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswap_v3

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

// IQuoterQuoteExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterQuoteExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	AmountIn          *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterQuoteExactInputSingleWithPoolParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterQuoteExactInputSingleWithPoolParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	AmountIn          *big.Int
	Pool              common.Address
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterQuoteExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterQuoteExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Amount            *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterQuoteExactOutputSingleWithPoolParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterQuoteExactOutputSingleWithPoolParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Amount            *big.Int
	Fee               *big.Int
	Pool              common.Address
	SqrtPriceLimitX96 *big.Int
}

// QuoterV3MetaData contains all meta data concerning the QuoterV3 contract.
var QuoterV3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoter.QuoteExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReceived\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoter.QuoteExactInputSingleWithPoolParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactInputSingleWithPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountReceived\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoter.QuoteExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoter.QuoteExactOutputSingleWithPoolParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactOutputSingleWithPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// QuoterV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use QuoterV3MetaData.ABI instead.
var QuoterV3ABI = QuoterV3MetaData.ABI

// QuoterV3 is an auto generated Go binding around an Ethereum contract.
type QuoterV3 struct {
	QuoterV3Caller     // Read-only binding to the contract
	QuoterV3Transactor // Write-only binding to the contract
	QuoterV3Filterer   // Log filterer for contract events
}

// QuoterV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type QuoterV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type QuoterV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuoterV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuoterV3Session struct {
	Contract     *QuoterV3         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuoterV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuoterV3CallerSession struct {
	Contract *QuoterV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// QuoterV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuoterV3TransactorSession struct {
	Contract     *QuoterV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// QuoterV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type QuoterV3Raw struct {
	Contract *QuoterV3 // Generic contract binding to access the raw methods on
}

// QuoterV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuoterV3CallerRaw struct {
	Contract *QuoterV3Caller // Generic read-only contract binding to access the raw methods on
}

// QuoterV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuoterV3TransactorRaw struct {
	Contract *QuoterV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewQuoterV3 creates a new instance of QuoterV3, bound to a specific deployed contract.
func NewQuoterV3(address common.Address, backend bind.ContractBackend) (*QuoterV3, error) {
	contract, err := bindQuoterV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QuoterV3{QuoterV3Caller: QuoterV3Caller{contract: contract}, QuoterV3Transactor: QuoterV3Transactor{contract: contract}, QuoterV3Filterer: QuoterV3Filterer{contract: contract}}, nil
}

// NewQuoterV3Caller creates a new read-only instance of QuoterV3, bound to a specific deployed contract.
func NewQuoterV3Caller(address common.Address, caller bind.ContractCaller) (*QuoterV3Caller, error) {
	contract, err := bindQuoterV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterV3Caller{contract: contract}, nil
}

// NewQuoterV3Transactor creates a new write-only instance of QuoterV3, bound to a specific deployed contract.
func NewQuoterV3Transactor(address common.Address, transactor bind.ContractTransactor) (*QuoterV3Transactor, error) {
	contract, err := bindQuoterV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterV3Transactor{contract: contract}, nil
}

// NewQuoterV3Filterer creates a new log filterer instance of QuoterV3, bound to a specific deployed contract.
func NewQuoterV3Filterer(address common.Address, filterer bind.ContractFilterer) (*QuoterV3Filterer, error) {
	contract, err := bindQuoterV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuoterV3Filterer{contract: contract}, nil
}

// bindQuoterV3 binds a generic wrapper to an already deployed contract.
func bindQuoterV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QuoterV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuoterV3 *QuoterV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuoterV3.Contract.QuoterV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuoterV3 *QuoterV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuoterV3.Contract.QuoterV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuoterV3 *QuoterV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuoterV3.Contract.QuoterV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuoterV3 *QuoterV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuoterV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuoterV3 *QuoterV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuoterV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuoterV3 *QuoterV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuoterV3.Contract.contract.Transact(opts, method, params...)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV3 *QuoterV3Caller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV3 *QuoterV3Session) Factory() (common.Address, error) {
	return _QuoterV3.Contract.Factory(&_QuoterV3.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV3 *QuoterV3CallerSession) Factory() (common.Address, error) {
	return _QuoterV3.Contract.Factory(&_QuoterV3.CallOpts)
}

// QuoteExactInput is a free data retrieval call binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) view returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactInput(opts *bind.CallOpts, path []byte, amountIn *big.Int) (struct {
	AmountOut                   *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactInput", path, amountIn)

	outstruct := new(struct {
		AmountOut                   *big.Int
		SqrtPriceX96AfterList       []*big.Int
		InitializedTicksCrossedList []uint32
		GasEstimate                 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountOut = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96AfterList = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.InitializedTicksCrossedList = *abi.ConvertType(out[2], new([]uint32)).(*[]uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactInput is a free data retrieval call binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) view returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactInput(path []byte, amountIn *big.Int) (struct {
	AmountOut                   *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInput(&_QuoterV3.CallOpts, path, amountIn)
}

// QuoteExactInput is a free data retrieval call binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) view returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactInput(path []byte, amountIn *big.Int) (struct {
	AmountOut                   *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInput(&_QuoterV3.CallOpts, path, amountIn)
}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactInputSingle(opts *bind.CallOpts, params IQuoterQuoteExactInputSingleParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactInputSingle", params)

	outstruct := new(struct {
		AmountReceived          *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountReceived = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96After = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactInputSingle(params IQuoterQuoteExactInputSingleParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInputSingle(&_QuoterV3.CallOpts, params)
}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactInputSingle(params IQuoterQuoteExactInputSingleParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInputSingle(&_QuoterV3.CallOpts, params)
}

// QuoteExactInputSingleWithPool is a free data retrieval call binding the contract method 0xd85c3d63.
//
// Solidity: function quoteExactInputSingleWithPool((address,address,uint256,address,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactInputSingleWithPool(opts *bind.CallOpts, params IQuoterQuoteExactInputSingleWithPoolParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactInputSingleWithPool", params)

	outstruct := new(struct {
		AmountReceived          *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountReceived = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96After = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactInputSingleWithPool is a free data retrieval call binding the contract method 0xd85c3d63.
//
// Solidity: function quoteExactInputSingleWithPool((address,address,uint256,address,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactInputSingleWithPool(params IQuoterQuoteExactInputSingleWithPoolParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInputSingleWithPool(&_QuoterV3.CallOpts, params)
}

// QuoteExactInputSingleWithPool is a free data retrieval call binding the contract method 0xd85c3d63.
//
// Solidity: function quoteExactInputSingleWithPool((address,address,uint256,address,uint24,uint160) params) view returns(uint256 amountReceived, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactInputSingleWithPool(params IQuoterQuoteExactInputSingleWithPoolParams) (struct {
	AmountReceived          *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactInputSingleWithPool(&_QuoterV3.CallOpts, params)
}

// QuoteExactOutput is a free data retrieval call binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) view returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactOutput(opts *bind.CallOpts, path []byte, amountOut *big.Int) (struct {
	AmountIn                    *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactOutput", path, amountOut)

	outstruct := new(struct {
		AmountIn                    *big.Int
		SqrtPriceX96AfterList       []*big.Int
		InitializedTicksCrossedList []uint32
		GasEstimate                 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountIn = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96AfterList = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.InitializedTicksCrossedList = *abi.ConvertType(out[2], new([]uint32)).(*[]uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactOutput is a free data retrieval call binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) view returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactOutput(path []byte, amountOut *big.Int) (struct {
	AmountIn                    *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutput(&_QuoterV3.CallOpts, path, amountOut)
}

// QuoteExactOutput is a free data retrieval call binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) view returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactOutput(path []byte, amountOut *big.Int) (struct {
	AmountIn                    *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutput(&_QuoterV3.CallOpts, path, amountOut)
}

// QuoteExactOutputSingle is a free data retrieval call binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactOutputSingle(opts *bind.CallOpts, params IQuoterQuoteExactOutputSingleParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactOutputSingle", params)

	outstruct := new(struct {
		AmountIn                *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountIn = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96After = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactOutputSingle is a free data retrieval call binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactOutputSingle(params IQuoterQuoteExactOutputSingleParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutputSingle(&_QuoterV3.CallOpts, params)
}

// QuoteExactOutputSingle is a free data retrieval call binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactOutputSingle(params IQuoterQuoteExactOutputSingleParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutputSingle(&_QuoterV3.CallOpts, params)
}

// QuoteExactOutputSingleWithPool is a free data retrieval call binding the contract method 0x86e3a7cf.
//
// Solidity: function quoteExactOutputSingleWithPool((address,address,uint256,uint24,address,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Caller) QuoteExactOutputSingleWithPool(opts *bind.CallOpts, params IQuoterQuoteExactOutputSingleWithPoolParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	var out []interface{}
	err := _QuoterV3.contract.Call(opts, &out, "quoteExactOutputSingleWithPool", params)

	outstruct := new(struct {
		AmountIn                *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountIn = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SqrtPriceX96After = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactOutputSingleWithPool is a free data retrieval call binding the contract method 0x86e3a7cf.
//
// Solidity: function quoteExactOutputSingleWithPool((address,address,uint256,uint24,address,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3Session) QuoteExactOutputSingleWithPool(params IQuoterQuoteExactOutputSingleWithPoolParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutputSingleWithPool(&_QuoterV3.CallOpts, params)
}

// QuoteExactOutputSingleWithPool is a free data retrieval call binding the contract method 0x86e3a7cf.
//
// Solidity: function quoteExactOutputSingleWithPool((address,address,uint256,uint24,address,uint160) params) view returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV3 *QuoterV3CallerSession) QuoteExactOutputSingleWithPool(params IQuoterQuoteExactOutputSingleWithPoolParams) (struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}, error) {
	return _QuoterV3.Contract.QuoteExactOutputSingleWithPool(&_QuoterV3.CallOpts, params)
}
