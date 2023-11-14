// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"

	zetacorecommon "github.com/zeta-chain/zetacore/common"
)

// CrosschainFungibleKeeper is an autogenerated mock type for the CrosschainFungibleKeeper type
type CrosschainFungibleKeeper struct {
	mock.Mock
}

// CallUniswapV2RouterSwapExactETHForToken provides a mock function with given fields: ctx, sender, to, amountIn, outZRC4, noEthereumTxEvent
func (_m *CrosschainFungibleKeeper) CallUniswapV2RouterSwapExactETHForToken(ctx types.Context, sender common.Address, to common.Address, amountIn *big.Int, outZRC4 common.Address, noEthereumTxEvent bool) ([]*big.Int, error) {
	ret := _m.Called(ctx, sender, to, amountIn, outZRC4, noEthereumTxEvent)

	var r0 []*big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, bool) ([]*big.Int, error)); ok {
		return rf(ctx, sender, to, amountIn, outZRC4, noEthereumTxEvent)
	}
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, bool) []*big.Int); ok {
		r0 = rf(ctx, sender, to, amountIn, outZRC4, noEthereumTxEvent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, bool) error); ok {
		r1 = rf(ctx, sender, to, amountIn, outZRC4, noEthereumTxEvent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallUniswapV2RouterSwapExactTokensForTokens provides a mock function with given fields: ctx, sender, to, amountIn, inZRC4, outZRC4, noEthereumTxEvent
func (_m *CrosschainFungibleKeeper) CallUniswapV2RouterSwapExactTokensForTokens(ctx types.Context, sender common.Address, to common.Address, amountIn *big.Int, inZRC4 common.Address, outZRC4 common.Address, noEthereumTxEvent bool) ([]*big.Int, error) {
	ret := _m.Called(ctx, sender, to, amountIn, inZRC4, outZRC4, noEthereumTxEvent)

	var r0 []*big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, common.Address, bool) ([]*big.Int, error)); ok {
		return rf(ctx, sender, to, amountIn, inZRC4, outZRC4, noEthereumTxEvent)
	}
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, common.Address, bool) []*big.Int); ok {
		r0 = rf(ctx, sender, to, amountIn, inZRC4, outZRC4, noEthereumTxEvent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, common.Address, common.Address, *big.Int, common.Address, common.Address, bool) error); ok {
		r1 = rf(ctx, sender, to, amountIn, inZRC4, outZRC4, noEthereumTxEvent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CallZRC20Approve provides a mock function with given fields: ctx, owner, zrc20address, spender, amount, noEthereumTxEvent
func (_m *CrosschainFungibleKeeper) CallZRC20Approve(ctx types.Context, owner common.Address, zrc20address common.Address, spender common.Address, amount *big.Int, noEthereumTxEvent bool) error {
	ret := _m.Called(ctx, owner, zrc20address, spender, amount, noEthereumTxEvent)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, common.Address, *big.Int, bool) error); ok {
		r0 = rf(ctx, owner, zrc20address, spender, amount, noEthereumTxEvent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CallZRC20Burn provides a mock function with given fields: ctx, sender, zrc20address, amount, noEthereumTxEvent
func (_m *CrosschainFungibleKeeper) CallZRC20Burn(ctx types.Context, sender common.Address, zrc20address common.Address, amount *big.Int, noEthereumTxEvent bool) error {
	ret := _m.Called(ctx, sender, zrc20address, amount, noEthereumTxEvent)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int, bool) error); ok {
		r0 = rf(ctx, sender, zrc20address, amount, noEthereumTxEvent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeployZRC20Contract provides a mock function with given fields: ctx, name, symbol, decimals, chainID, coinType, erc20Contract, gasLimit
func (_m *CrosschainFungibleKeeper) DeployZRC20Contract(ctx types.Context, name string, symbol string, decimals uint8, chainID int64, coinType zetacorecommon.CoinType, erc20Contract string, gasLimit *big.Int) (common.Address, error) {
	ret := _m.Called(ctx, name, symbol, decimals, chainID, coinType, erc20Contract, gasLimit)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, string, string, uint8, int64, zetacorecommon.CoinType, string, *big.Int) (common.Address, error)); ok {
		return rf(ctx, name, symbol, decimals, chainID, coinType, erc20Contract, gasLimit)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string, string, uint8, int64, zetacorecommon.CoinType, string, *big.Int) common.Address); ok {
		r0 = rf(ctx, name, symbol, decimals, chainID, coinType, erc20Contract, gasLimit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, string, string, uint8, int64, zetacorecommon.CoinType, string, *big.Int) error); ok {
		r1 = rf(ctx, name, symbol, decimals, chainID, coinType, erc20Contract, gasLimit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DepositCoinZeta provides a mock function with given fields: ctx, to, amount
func (_m *CrosschainFungibleKeeper) DepositCoinZeta(ctx types.Context, to common.Address, amount *big.Int) error {
	ret := _m.Called(ctx, to, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, *big.Int) error); ok {
		r0 = rf(ctx, to, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DepositZRC20 provides a mock function with given fields: ctx, contract, to, amount
func (_m *CrosschainFungibleKeeper) DepositZRC20(ctx types.Context, contract common.Address, to common.Address, amount *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
	ret := _m.Called(ctx, contract, to, amount)

	var r0 *evmtypes.MsgEthereumTxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int) (*evmtypes.MsgEthereumTxResponse, error)); ok {
		return rf(ctx, contract, to, amount)
	}
	if rf, ok := ret.Get(0).(func(types.Context, common.Address, common.Address, *big.Int) *evmtypes.MsgEthereumTxResponse); ok {
		r0 = rf(ctx, contract, to, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*evmtypes.MsgEthereumTxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, common.Address, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, contract, to, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FundGasStabilityPool provides a mock function with given fields: ctx, chainID, amount
func (_m *CrosschainFungibleKeeper) FundGasStabilityPool(ctx types.Context, chainID int64, amount *big.Int) error {
	ret := _m.Called(ctx, chainID, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, int64, *big.Int) error); ok {
		r0 = rf(ctx, chainID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllForeignCoins provides a mock function with given fields: ctx
func (_m *CrosschainFungibleKeeper) GetAllForeignCoins(ctx types.Context) []fungibletypes.ForeignCoin {
	ret := _m.Called(ctx)

	var r0 []fungibletypes.ForeignCoin
	if rf, ok := ret.Get(0).(func(types.Context) []fungibletypes.ForeignCoin); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fungibletypes.ForeignCoin)
		}
	}

	return r0
}

// GetAllForeignCoinsForChain provides a mock function with given fields: ctx, foreignChainID
func (_m *CrosschainFungibleKeeper) GetAllForeignCoinsForChain(ctx types.Context, foreignChainID int64) []fungibletypes.ForeignCoin {
	ret := _m.Called(ctx, foreignChainID)

	var r0 []fungibletypes.ForeignCoin
	if rf, ok := ret.Get(0).(func(types.Context, int64) []fungibletypes.ForeignCoin); ok {
		r0 = rf(ctx, foreignChainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fungibletypes.ForeignCoin)
		}
	}

	return r0
}

// GetForeignCoinFromAsset provides a mock function with given fields: ctx, asset, chainID
func (_m *CrosschainFungibleKeeper) GetForeignCoinFromAsset(ctx types.Context, asset string, chainID int64) (fungibletypes.ForeignCoin, bool) {
	ret := _m.Called(ctx, asset, chainID)

	var r0 fungibletypes.ForeignCoin
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string, int64) (fungibletypes.ForeignCoin, bool)); ok {
		return rf(ctx, asset, chainID)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string, int64) fungibletypes.ForeignCoin); ok {
		r0 = rf(ctx, asset, chainID)
	} else {
		r0 = ret.Get(0).(fungibletypes.ForeignCoin)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string, int64) bool); ok {
		r1 = rf(ctx, asset, chainID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetForeignCoins provides a mock function with given fields: ctx, zrc20Addr
func (_m *CrosschainFungibleKeeper) GetForeignCoins(ctx types.Context, zrc20Addr string) (fungibletypes.ForeignCoin, bool) {
	ret := _m.Called(ctx, zrc20Addr)

	var r0 fungibletypes.ForeignCoin
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (fungibletypes.ForeignCoin, bool)); ok {
		return rf(ctx, zrc20Addr)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) fungibletypes.ForeignCoin); ok {
		r0 = rf(ctx, zrc20Addr)
	} else {
		r0 = ret.Get(0).(fungibletypes.ForeignCoin)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, zrc20Addr)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetGasCoinForForeignCoin provides a mock function with given fields: ctx, chainID
func (_m *CrosschainFungibleKeeper) GetGasCoinForForeignCoin(ctx types.Context, chainID int64) (fungibletypes.ForeignCoin, bool) {
	ret := _m.Called(ctx, chainID)

	var r0 fungibletypes.ForeignCoin
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, int64) (fungibletypes.ForeignCoin, bool)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(types.Context, int64) fungibletypes.ForeignCoin); ok {
		r0 = rf(ctx, chainID)
	} else {
		r0 = ret.Get(0).(fungibletypes.ForeignCoin)
	}

	if rf, ok := ret.Get(1).(func(types.Context, int64) bool); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetSystemContract provides a mock function with given fields: ctx
func (_m *CrosschainFungibleKeeper) GetSystemContract(ctx types.Context) (fungibletypes.SystemContract, bool) {
	ret := _m.Called(ctx)

	var r0 fungibletypes.SystemContract
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context) (fungibletypes.SystemContract, bool)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(types.Context) fungibletypes.SystemContract); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(fungibletypes.SystemContract)
	}

	if rf, ok := ret.Get(1).(func(types.Context) bool); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetUniswapV2Router02Address provides a mock function with given fields: ctx
func (_m *CrosschainFungibleKeeper) GetUniswapV2Router02Address(ctx types.Context) (common.Address, error) {
	ret := _m.Called(ctx)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context) (common.Address, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(types.Context) common.Address); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryGasLimit provides a mock function with given fields: ctx, contract
func (_m *CrosschainFungibleKeeper) QueryGasLimit(ctx types.Context, contract common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, contract)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) (*big.Int, error)); ok {
		return rf(ctx, contract)
	}
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) *big.Int); ok {
		r0 = rf(ctx, contract)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, common.Address) error); ok {
		r1 = rf(ctx, contract)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryProtocolFlatFee provides a mock function with given fields: ctx, contract
func (_m *CrosschainFungibleKeeper) QueryProtocolFlatFee(ctx types.Context, contract common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, contract)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) (*big.Int, error)); ok {
		return rf(ctx, contract)
	}
	if rf, ok := ret.Get(0).(func(types.Context, common.Address) *big.Int); ok {
		r0 = rf(ctx, contract)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, common.Address) error); ok {
		r1 = rf(ctx, contract)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QuerySystemContractGasCoinZRC20 provides a mock function with given fields: ctx, chainID
func (_m *CrosschainFungibleKeeper) QuerySystemContractGasCoinZRC20(ctx types.Context, chainID *big.Int) (common.Address, error) {
	ret := _m.Called(ctx, chainID)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int) (common.Address, error)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int) common.Address); ok {
		r0 = rf(ctx, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *big.Int) error); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryUniswapV2RouterGetZRC4AmountsIn provides a mock function with given fields: ctx, amountOut, inZRC4
func (_m *CrosschainFungibleKeeper) QueryUniswapV2RouterGetZRC4AmountsIn(ctx types.Context, amountOut *big.Int, inZRC4 common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, amountOut, inZRC4)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address) (*big.Int, error)); ok {
		return rf(ctx, amountOut, inZRC4)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address) *big.Int); ok {
		r0 = rf(ctx, amountOut, inZRC4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *big.Int, common.Address) error); ok {
		r1 = rf(ctx, amountOut, inZRC4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn provides a mock function with given fields: ctx, amountOut, inZRC4, outZRC4
func (_m *CrosschainFungibleKeeper) QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(ctx types.Context, amountOut *big.Int, inZRC4 common.Address, outZRC4 common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, amountOut, inZRC4, outZRC4)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address, common.Address) (*big.Int, error)); ok {
		return rf(ctx, amountOut, inZRC4, outZRC4)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address, common.Address) *big.Int); ok {
		r0 = rf(ctx, amountOut, inZRC4, outZRC4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *big.Int, common.Address, common.Address) error); ok {
		r1 = rf(ctx, amountOut, inZRC4, outZRC4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryUniswapV2RouterGetZetaAmountsIn provides a mock function with given fields: ctx, amountOut, outZRC4
func (_m *CrosschainFungibleKeeper) QueryUniswapV2RouterGetZetaAmountsIn(ctx types.Context, amountOut *big.Int, outZRC4 common.Address) (*big.Int, error) {
	ret := _m.Called(ctx, amountOut, outZRC4)

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address) (*big.Int, error)); ok {
		return rf(ctx, amountOut, outZRC4)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, common.Address) *big.Int); ok {
		r0 = rf(ctx, amountOut, outZRC4)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *big.Int, common.Address) error); ok {
		r1 = rf(ctx, amountOut, outZRC4)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetForeignCoins provides a mock function with given fields: ctx, foreignCoins
func (_m *CrosschainFungibleKeeper) SetForeignCoins(ctx types.Context, foreignCoins fungibletypes.ForeignCoin) {
	_m.Called(ctx, foreignCoins)
}

// SetGasPrice provides a mock function with given fields: ctx, chainID, gasPrice
func (_m *CrosschainFungibleKeeper) SetGasPrice(ctx types.Context, chainID *big.Int, gasPrice *big.Int) (uint64, error) {
	ret := _m.Called(ctx, chainID, gasPrice)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, *big.Int) (uint64, error)); ok {
		return rf(ctx, chainID, gasPrice)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *big.Int, *big.Int) uint64); ok {
		r0 = rf(ctx, chainID, gasPrice)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(types.Context, *big.Int, *big.Int) error); ok {
		r1 = rf(ctx, chainID, gasPrice)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithdrawFromGasStabilityPool provides a mock function with given fields: ctx, chainID, amount
func (_m *CrosschainFungibleKeeper) WithdrawFromGasStabilityPool(ctx types.Context, chainID int64, amount *big.Int) error {
	ret := _m.Called(ctx, chainID, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, int64, *big.Int) error); ok {
		r0 = rf(ctx, chainID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ZRC20DepositAndCallContract provides a mock function with given fields: ctx, from, to, amount, senderChain, data, coinType, asset
func (_m *CrosschainFungibleKeeper) ZRC20DepositAndCallContract(ctx types.Context, from []byte, to common.Address, amount *big.Int, senderChain *zetacorecommon.Chain, data []byte, coinType zetacorecommon.CoinType, asset string) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	ret := _m.Called(ctx, from, to, amount, senderChain, data, coinType, asset)

	var r0 *evmtypes.MsgEthereumTxResponse
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(types.Context, []byte, common.Address, *big.Int, *zetacorecommon.Chain, []byte, zetacorecommon.CoinType, string) (*evmtypes.MsgEthereumTxResponse, bool, error)); ok {
		return rf(ctx, from, to, amount, senderChain, data, coinType, asset)
	}
	if rf, ok := ret.Get(0).(func(types.Context, []byte, common.Address, *big.Int, *zetacorecommon.Chain, []byte, zetacorecommon.CoinType, string) *evmtypes.MsgEthereumTxResponse); ok {
		r0 = rf(ctx, from, to, amount, senderChain, data, coinType, asset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*evmtypes.MsgEthereumTxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, []byte, common.Address, *big.Int, *zetacorecommon.Chain, []byte, zetacorecommon.CoinType, string) bool); ok {
		r1 = rf(ctx, from, to, amount, senderChain, data, coinType, asset)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(types.Context, []byte, common.Address, *big.Int, *zetacorecommon.Chain, []byte, zetacorecommon.CoinType, string) error); ok {
		r2 = rf(ctx, from, to, amount, senderChain, data, coinType, asset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewCrosschainFungibleKeeper creates a new instance of CrosschainFungibleKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainFungibleKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainFungibleKeeper {
	mock := &CrosschainFungibleKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
