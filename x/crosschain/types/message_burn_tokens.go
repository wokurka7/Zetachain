package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
)

const BurnTokens = "BurnTokens"

var _ sdk.Msg = &MsgBurnTokens{}

func NewMsgBurnTokens(creator string, chainID int64, amount sdkmath.Uint) *MsgBurnTokens {
	return &MsgBurnTokens{
		Creator: creator,
		ChainId: chainID,
		Amount:  amount,
	}
}

func (msg *MsgBurnTokens) Route() string {
	return RouterKey
}

func (msg *MsgBurnTokens) Type() string {
	return BurnTokens
}

func (msg *MsgBurnTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !common.IsEVMChain(msg.ChainId) {
		return ErrInvalidChainID
	}
	if msg.Amount.IsZero() || msg.Amount.IsNil() {
		return ErrInvalidAmount

	}
	return nil
}
