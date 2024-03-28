package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ReduceZetaSupply reduces the Zeta supply on external chains
// This is done by burning Zeta tokens on the external chain by sending it to the burn address.
// The message uses the provided chainID, amount, and burn address to create a new CCTX. If a burn address is not provided, the default burn address is used.
// This cctx can only be created by the policy account of group2
func (k msgServer) ReduceZetaSupply(goCtx context.Context, msg *types.MsgReduceZetaSupply) (*types.MsgReduceZetaSupplyResponse, error) {
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

	burnAddress := types.BurnTokensEVMAddress
	if msg.BurnAddress != "" {
		burnAddress = ethcommon.HexToAddress(msg.BurnAddress).Hex()
	}

	cctx := GetCmdCCTX(ctx, msg.ChainId, msg.Amount, tss, multipliedGasPrice.String(), burnAddress)
	err = k.ProcessCmdCCTX(ctx, cctx)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrCannotCreateBurnCCTX, fmt.Sprintf("cannot process burn CCTX %s", err))
	}
	k.SaveCmdCCTX(ctx, cctx)
	return &types.MsgReduceZetaSupplyResponse{}, nil
}

// GetCmdCCTX returns a new CrossChainTx for of Cointype Cmd which is used to reduce the Zeta supply on external chains
func GetCmdCCTX(ctx sdk.Context, chainID int64, amount sdkmath.Uint, tss observertypes.TSS, multipliedGasPrice string, burnAddress string) types.CrossChainTx {
	height := ctx.BlockHeight()
	// use chainID and height to create a unique index for the CCTX.
	index := GetIndexForReduceZetaSupplyCMD(chainID, height)

	cctx := types.CrossChainTx{
		Creator:        "",
		Index:          index,
		ZetaFees:       sdkmath.Uint{},
		RelayedMessage: fmt.Sprintf("%s:%s", common.CmdReduceZetaSupply, fmt.Sprintf("Reduce Zeta supply on chain %d by %s", chainID, amount.String())),
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
			InboundTxFinalizedZetaHeight:    uint64(height),
		},
		OutboundTxParams: []*types.OutboundTxParams{{
			Receiver:                         burnAddress,
			ReceiverChainId:                  chainID,
			CoinType:                         common.CoinType_Cmd,
			Amount:                           amount,
			OutboundTxTssNonce:               0,
			OutboundTxGasLimit:               1_000_000,
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

// ProcessCmdCCTX updates the nonce for the given chainID and CCTX
func (k Keeper) ProcessCmdCCTX(ctx sdk.Context, cctx types.CrossChainTx) error {
	err := k.UpdateNonce(ctx, cctx.GetCurrentOutTxParam().ReceiverChainId, &cctx)
	if err != nil {
		return err
	}
	return nil
}

// SaveCmdCCTX saves the given CCTX to the store
func (k Keeper) SaveCmdCCTX(ctx sdk.Context, cctx types.CrossChainTx) {
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
}

// GetIndexForReduceZetaSupplyCMD returns the index for the given chainID and blockHeight this is used to identify the CCTX for reducing the Zeta supply
func GetIndexForReduceZetaSupplyCMD(chainID int64, blockHeight int64) string {
	indexString := fmt.Sprintf("%d-%d", chainID, blockHeight)
	hash := crypto.Keccak256Hash([]byte(indexString))
	return hash.Hex()
}
