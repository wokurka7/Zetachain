// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	common "github.com/zeta-chain/zetacore/common"

	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleObserverKeeper is an autogenerated mock type for the FungibleObserverKeeper type
type FungibleObserverKeeper struct {
	mock.Mock
}

// GetAllBallots provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetAllBallots(ctx types.Context) []*observertypes.Ballot {
	ret := _m.Called(ctx)

	var r0 []*observertypes.Ballot
	if rf, ok := ret.Get(0).(func(types.Context) []*observertypes.Ballot); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*observertypes.Ballot)
		}
	}

	return r0
}

// GetAllObserverMappers provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetAllObserverMappers(ctx types.Context) []*observertypes.ObserverMapper {
	ret := _m.Called(ctx)

	var r0 []*observertypes.ObserverMapper
	if rf, ok := ret.Get(0).(func(types.Context) []*observertypes.ObserverMapper); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*observertypes.ObserverMapper)
		}
	}

	return r0
}

// GetBallot provides a mock function with given fields: ctx, index
func (_m *FungibleObserverKeeper) GetBallot(ctx types.Context, index string) (observertypes.Ballot, bool) {
	ret := _m.Called(ctx, index)

	var r0 observertypes.Ballot
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (observertypes.Ballot, bool)); ok {
		return rf(ctx, index)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) observertypes.Ballot); ok {
		r0 = rf(ctx, index)
	} else {
		r0 = ret.Get(0).(observertypes.Ballot)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, index)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetCoreParamsByChainID provides a mock function with given fields: ctx, chainID
func (_m *FungibleObserverKeeper) GetCoreParamsByChainID(ctx types.Context, chainID int64) (*common.CoreParams, bool) {
	ret := _m.Called(ctx, chainID)

	var r0 *common.CoreParams
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, int64) (*common.CoreParams, bool)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(types.Context, int64) *common.CoreParams); ok {
		r0 = rf(ctx, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.CoreParams)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, int64) bool); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetObserverMapper provides a mock function with given fields: ctx, chain
func (_m *FungibleObserverKeeper) GetObserverMapper(ctx types.Context, chain *common.Chain) (observertypes.ObserverMapper, bool) {
	ret := _m.Called(ctx, chain)

	var r0 observertypes.ObserverMapper
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, *common.Chain) (observertypes.ObserverMapper, bool)); ok {
		return rf(ctx, chain)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *common.Chain) observertypes.ObserverMapper); ok {
		r0 = rf(ctx, chain)
	} else {
		r0 = ret.Get(0).(observertypes.ObserverMapper)
	}

	if rf, ok := ret.Get(1).(func(types.Context, *common.Chain) bool); ok {
		r1 = rf(ctx, chain)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetParams provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetParams(ctx types.Context) observertypes.Params {
	ret := _m.Called(ctx)

	var r0 observertypes.Params
	if rf, ok := ret.Get(0).(func(types.Context) observertypes.Params); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(observertypes.Params)
	}

	return r0
}

// SetBallot provides a mock function with given fields: ctx, ballot
func (_m *FungibleObserverKeeper) SetBallot(ctx types.Context, ballot *observertypes.Ballot) {
	_m.Called(ctx, ballot)
}

// SetObserverMapper provides a mock function with given fields: ctx, om
func (_m *FungibleObserverKeeper) SetObserverMapper(ctx types.Context, om *observertypes.ObserverMapper) {
	_m.Called(ctx, om)
}

// NewFungibleObserverKeeper creates a new instance of FungibleObserverKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleObserverKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleObserverKeeper {
	mock := &FungibleObserverKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
