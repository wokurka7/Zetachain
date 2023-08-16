package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestMsgDeployFungibleCoinZRC4_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeployFungibleCoinZRC20
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeployFungibleCoinZRC20{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid gas limit",
			msg: MsgDeployFungibleCoinZRC20{
				Creator:  sample.AccAddress(),
				GasLimit: -1,
			},
			err: sdkerrors.ErrInvalidGasLimit,
		},
		{
			name: "valid message",
			msg: MsgDeployFungibleCoinZRC20{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
