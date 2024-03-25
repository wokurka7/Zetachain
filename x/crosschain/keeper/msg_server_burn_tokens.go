package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) BurnTokens(goCtx context.Context, msg *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group2) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Update can only be executed by the correct policy account")
	}

	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrCannotCreateBurnCCTX, "cannot find current TSS")
	}
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, msg.ChainId)
	if !isFound {
		return nil, types.ErrUnableToGetGasPrice
	}

	multipliedGasPrice, err := common.MultiplyGasPrice(medianGasPrice, types.TssMigrationGasMultiplierEVM)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrCannotCreateBurnCCTX, fmt.Sprintf("cannot multiply gas price %s", err))
	}

	cctx := k.GetBurnCCTX(ctx, msg.ChainId, msg.Amount, tss, multipliedGasPrice.String())
	err = k.ProcessBurnCCTX(ctx, cctx)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrCannotCreateBurnCCTX, fmt.Sprintf("cannot process burn CCTX %s", err))
	}
	k.SaveBurnCCTX(ctx, cctx)
	return &types.MsgBurnTokensResponse{}, nil
}

func (k Keeper) GetBurnCCTX(ctx sdk.Context, chainID int64, amount sdkmath.Uint, tss observertypes.TSS, multipliedGasPrice string) types.CrossChainTx {

	indexString := GetIndexStingForBurnCCTX(chainID, ctx.BlockHeight())
	hash := crypto.Keccak256Hash([]byte(indexString))
	index := hash.Hex()

	cctx := types.CrossChainTx{
		Creator:        "",
		Index:          index,
		ZetaFees:       sdkmath.Uint{},
		RelayedMessage: fmt.Sprintf("%s:%s", common.CmdBurnTokens, "Funds Migrator Admin Cmd"),
		CctxStatus: &types.Status{
			Status:              types.CctxStatus_PendingOutbound,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
		},
		InboundTxParams: &types.InboundTxParams{
			Sender:                          "",
			SenderChainId:                   chainID,
			TxOrigin:                        "",
			CoinType:                        common.CoinType_Cmd,
			Asset:                           "",
			Amount:                          amount,
			InboundTxObservedHash:           tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash()).String(),
			InboundTxObservedExternalHeight: 0,
			InboundTxBallotIndex:            "",
			InboundTxFinalizedZetaHeight:    0,
		},
		OutboundTxParams: []*types.OutboundTxParams{{
			Receiver:                         types.BurnTokensEVMAddress,
			ReceiverChainId:                  chainID,
			CoinType:                         common.CoinType_Cmd,
			Amount:                           amount,
			OutboundTxTssNonce:               0,
			OutboundTxGasLimit:               common.EVMSend,
			OutboundTxGasPrice:               multipliedGasPrice,
			OutboundTxHash:                   "",
			OutboundTxBallotIndex:            "",
			OutboundTxObservedExternalHeight: 0,
			OutboundTxGasUsed:                0,
			OutboundTxEffectiveGasPrice:      sdkmath.Int{},
			OutboundTxEffectiveGasLimit:      0,
			TssPubkey:                        tss.TssPubkey,
		}}}

	return cctx
}

func (k Keeper) ProcessBurnCCTX(ctx sdk.Context, cctx types.CrossChainTx) error {
	err := k.UpdateNonce(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, &cctx)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SaveBurnCCTX(ctx sdk.Context, cctx types.CrossChainTx) {
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
}

func GetIndexStingForBurnCCTX(chainID int64, blockHeight int64) string {
	return fmt.Sprintf("%d-%d-%s-%s", chainID, blockHeight)
}
