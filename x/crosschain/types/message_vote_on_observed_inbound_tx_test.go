package types_test

import (
	"testing"

	"math/rand"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgVoteOnObservedInboundTx_ValidateBasic(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	tests := []struct {
		name string
		msg  *types.MsgVoteOnObservedInboundTx
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgVoteOnObservedInboundTx(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				pkg.CoinType_Zeta,
				sample.String(),
				42,
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgVoteOnObservedInboundTx(
				"invalid_address",
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				pkg.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid sender chain ID",
			msg: types.NewMsgVoteOnObservedInboundTx(
				sample.AccAddress(),
				sample.AccAddress(),
				-1,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				pkg.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid receiver chain ID",
			msg: types.NewMsgVoteOnObservedInboundTx(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				-1,
				math.NewUint(42),
				sample.String(),
				sample.String(),
				42,
				42,
				pkg.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: types.ErrInvalidChainID,
		},
		{
			name: "invalid message length",
			msg: types.NewMsgVoteOnObservedInboundTx(
				sample.AccAddress(),
				sample.AccAddress(),
				42,
				sample.String(),
				sample.String(),
				42,
				math.NewUint(42),
				sample.StringRandom(r, types.MaxMessageLength+1),
				sample.String(),
				42,
				42,
				pkg.CoinType_Zeta,
				sample.String(),
				42,
			),
			err: sdkerrors.ErrInvalidRequest,
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

func TestMsgVoteOnObservedInboundTx_Digest(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	msg := types.MsgVoteOnObservedInboundTx{
		Creator:       sample.AccAddress(),
		Sender:        sample.AccAddress(),
		SenderChainId: 42,
		TxOrigin:      sample.String(),
		Receiver:      sample.String(),
		ReceiverChain: 42,
		Amount:        math.NewUint(42),
		Message:       sample.String(),
		InTxHash:      sample.String(),
		InBlockHeight: 42,
		GasLimit:      42,
		CoinType:      pkg.CoinType_Zeta,
		Asset:         sample.String(),
		EventIndex:    42,
	}
	hash := msg.Digest()
	require.NotEmpty(t, hash, "hash should not be empty")

	// creator not used
	msg2 := msg
	msg2.Creator = sample.AccAddress()
	hash2 := msg2.Digest()
	require.Equal(t, hash, hash2, "creator should not change hash")

	// in block height not used
	msg2 = msg
	msg2.InBlockHeight = 43
	hash2 = msg2.Digest()
	require.Equal(t, hash, hash2, "in block height should not change hash")

	// sender used
	msg2 = msg
	msg2.Sender = sample.AccAddress()
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "sender should change hash")

	// sender chain ID used
	msg2 = msg
	msg2.SenderChainId = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "sender chain ID should change hash")

	// tx origin used
	msg2 = msg
	msg2.TxOrigin = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "tx origin should change hash")

	// receiver used
	msg2 = msg
	msg2.Receiver = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "receiver should change hash")

	// receiver chain ID used
	msg2 = msg
	msg2.ReceiverChain = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "receiver chain ID should change hash")

	// amount used
	msg2 = msg
	msg2.Amount = math.NewUint(43)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "amount should change hash")

	// message used
	msg2 = msg
	msg2.Message = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "message should change hash")

	// in tx hash used
	msg2 = msg
	msg2.InTxHash = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "in tx hash should change hash")

	// gas limit used
	msg2 = msg
	msg2.GasLimit = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "gas limit should change hash")

	// coin type used
	msg2 = msg
	msg2.CoinType = pkg.CoinType_ERC20
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "coin type should change hash")

	// asset used
	msg2 = msg
	msg2.Asset = sample.StringRandom(r, 32)
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "asset should change hash")

	// event index used
	msg2 = msg
	msg2.EventIndex = 43
	hash2 = msg2.Digest()
	require.NotEqual(t, hash, hash2, "event index should change hash")
}

func TestMsgVoteOnObservedInboundTx_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgVoteOnObservedInboundTx
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgVoteOnObservedInboundTx{
				Creator: "invalid",
			},
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

func TestMsgVoteOnObservedInboundTx_Type(t *testing.T) {
	msg := types.MsgVoteOnObservedInboundTx{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, pkg.InboundVoter.String(), msg.Type())
}

func TestMsgVoteOnObservedInboundTx_Route(t *testing.T) {
	msg := types.MsgVoteOnObservedInboundTx{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteOnObservedInboundTx_GetSignBytes(t *testing.T) {
	msg := types.MsgVoteOnObservedInboundTx{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
