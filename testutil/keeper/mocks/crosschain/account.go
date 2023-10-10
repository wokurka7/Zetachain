// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CrosschainAccountKeeper is an autogenerated mock type for the CrosschainAccountKeeper type
type CrosschainAccountKeeper struct {
	mock.Mock
}

// GetAccount provides a mock function with given fields: ctx, addr
func (_m *CrosschainAccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) authtypes.AccountI {
	ret := _m.Called(ctx, addr)

	var r0 authtypes.AccountI
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) authtypes.AccountI); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(authtypes.AccountI)
		}
	}

	return r0
}

// GetModuleAccount provides a mock function with given fields: ctx, name
func (_m *CrosschainAccountKeeper) GetModuleAccount(ctx types.Context, name string) authtypes.ModuleAccountI {
	ret := _m.Called(ctx, name)

	var r0 authtypes.ModuleAccountI
	if rf, ok := ret.Get(0).(func(types.Context, string) authtypes.ModuleAccountI); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(authtypes.ModuleAccountI)
		}
	}

	return r0
}

// GetModuleAddress provides a mock function with given fields: name
func (_m *CrosschainAccountKeeper) GetModuleAddress(name string) types.AccAddress {
	ret := _m.Called(name)

	var r0 types.AccAddress
	if rf, ok := ret.Get(0).(func(string) types.AccAddress); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AccAddress)
		}
	}

	return r0
}

// SetModuleAccount provides a mock function with given fields: _a0, _a1
func (_m *CrosschainAccountKeeper) SetModuleAccount(_a0 types.Context, _a1 authtypes.ModuleAccountI) {
	_m.Called(_a0, _a1)
}

// NewCrosschainAccountKeeper creates a new instance of CrosschainAccountKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainAccountKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainAccountKeeper {
	mock := &CrosschainAccountKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
