package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgBurnTokens_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgBurnTokens
		err  require.ErrorAssertionFunc
	}{
		{
			name: "valid message",
			msg: &types.MsgBurnTokens{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: sample.EthAddress().String(),
			},
			err: require.NoError,
		},

		{
			name: "invalid creator address",
			msg: &types.MsgBurnTokens{
				Creator:     "invalid",
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: "",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid creator address")
			},
		},

		{
			name: "invalid chain ID",
			msg: &types.MsgBurnTokens{
				Creator: sample.AccAddress(),
				ChainId: 0,
				Amount:  sdkmath.NewUint(100),
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrInvalidChainID)
			},
		},

		{
			name: "invalid amount",
			msg: &types.MsgBurnTokens{
				Creator: sample.AccAddress(),
				ChainId: common.GoerliChain().ChainId,
				Amount:  sdkmath.NewUint(0),
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrInvalidAmount)
			},
		},

		{
			name: "invalid amount nil",
			msg: &types.MsgBurnTokens{
				Creator: sample.AccAddress(),
				ChainId: common.GoerliChain().ChainId,
				Amount:  sdkmath.Uint{},
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorIs(t, err, types.ErrInvalidAmount)
			},
		},

		{
			name: "invalid burn address",
			msg: &types.MsgBurnTokens{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: "invalid",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid burn address")
			},
		},

		{
			name: "valid burn address",
			msg: &types.MsgBurnTokens{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: "0x000000000000000000000000000000000000000",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid burn address")
			},
		},

		{
			name: "valid burn address",
			msg: &types.MsgBurnTokens{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: "0xbadbadbadbadbadbadbadbadbadbadbadbadbad",
			},
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid burn address")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}
