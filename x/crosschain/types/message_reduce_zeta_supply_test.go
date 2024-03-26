package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgBurnTokens_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgReduceZetaSupply
		err  require.ErrorAssertionFunc
	}{
		{
			name: "valid message",
			msg: &types.MsgReduceZetaSupply{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: sample.EthAddress().String(),
			},
			err: require.NoError,
		},

		{
			name: "invalid creator address",
			msg: &types.MsgReduceZetaSupply{
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
			msg: &types.MsgReduceZetaSupply{
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
			msg: &types.MsgReduceZetaSupply{
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
			msg: &types.MsgReduceZetaSupply{
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
			msg: &types.MsgReduceZetaSupply{
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
			msg: &types.MsgReduceZetaSupply{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: "0x000000000000000000000000000000000000000",
			},
			err: require.NoError,
		},

		{
			name: "valid burn address",
			msg: &types.MsgReduceZetaSupply{
				Creator:     sample.AccAddress(),
				ChainId:     common.GoerliChain().ChainId,
				Amount:      sdkmath.NewUint(100),
				BurnAddress: types.BurnTokensEVMAddress,
			},
			err: require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgReduceZetaSupply_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgReduceZetaSupply
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgReduceZetaSupply(signer, 1, sdkmath.OneUint(), ""),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgReduceZetaSupply("invalid", 1, sdkmath.OneUint(), ""),
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgReduceZetaSupply_Type(t *testing.T) {
	msg := types.NewMsgReduceZetaSupply(sample.AccAddress(), 1, sdkmath.OneUint(), "")
	require.Equal(t, types.ReduceZetaSupply, msg.Type())
}

func TestMsgReduceZetaSupply_Route(t *testing.T) {
	msg := types.NewMsgReduceZetaSupply(sample.AccAddress(), 1, sdkmath.OneUint(), "")
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgReduceZetaSupply_GetSignBytes(t *testing.T) {
	msg := types.NewMsgReduceZetaSupply(sample.AccAddress(), 1, sdkmath.OneUint(), "")
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
