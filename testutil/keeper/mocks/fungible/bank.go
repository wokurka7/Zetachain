// Code generated by mockery v2.35.2. DO NOT EDIT.

package mocks

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleBankKeeper is an autogenerated mock type for the FungibleBankKeeper type
type FungibleBankKeeper struct {
	mock.Mock
}

// BlockedAddr provides a mock function with given fields: addr
func (_m *FungibleBankKeeper) BlockedAddr(addr types.AccAddress) bool {
	ret := _m.Called(addr)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.AccAddress) bool); ok {
		r0 = rf(addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// BurnCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *FungibleBankKeeper) BurnCoins(ctx types.Context, moduleName string, amt types.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: ctx, addr, denom
func (_m *FungibleBankKeeper) GetBalance(ctx types.Context, addr types.AccAddress, denom string) types.Coin {
	ret := _m.Called(ctx, addr, denom)

	var r0 types.Coin
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, string) types.Coin); ok {
		r0 = rf(ctx, addr, denom)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	return r0
}

// GetDenomMetaData provides a mock function with given fields: ctx, denom
func (_m *FungibleBankKeeper) GetDenomMetaData(ctx types.Context, denom string) (banktypes.Metadata, bool) {
	ret := _m.Called(ctx, denom)

	var r0 banktypes.Metadata
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (banktypes.Metadata, bool)); ok {
		return rf(ctx, denom)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) banktypes.Metadata); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(banktypes.Metadata)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, denom)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// HasSupply provides a mock function with given fields: ctx, denom
func (_m *FungibleBankKeeper) HasSupply(ctx types.Context, denom string) bool {
	ret := _m.Called(ctx, denom)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) bool); ok {
		r0 = rf(ctx, denom)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsSendEnabledCoin provides a mock function with given fields: ctx, coin
func (_m *FungibleBankKeeper) IsSendEnabledCoin(ctx types.Context, coin types.Coin) bool {
	ret := _m.Called(ctx, coin)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context, types.Coin) bool); ok {
		r0 = rf(ctx, coin)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MintCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *FungibleBankKeeper) MintCoins(ctx types.Context, moduleName string, amt types.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromAccountToModule provides a mock function with given fields: ctx, senderAddr, recipientModule, amt
func (_m *FungibleBankKeeper) SendCoinsFromAccountToModule(ctx types.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	ret := _m.Called(ctx, senderAddr, recipientModule, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, string, types.Coins) error); ok {
		r0 = rf(ctx, senderAddr, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *FungibleBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.AccAddress, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDenomMetaData provides a mock function with given fields: ctx, denomMetaData
func (_m *FungibleBankKeeper) SetDenomMetaData(ctx types.Context, denomMetaData banktypes.Metadata) {
	_m.Called(ctx, denomMetaData)
}

// NewFungibleBankKeeper creates a new instance of FungibleBankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleBankKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleBankKeeper {
	mock := &FungibleBankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
