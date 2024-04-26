package keeper

import (
	"encoding/hex"
	"fmt"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// HandleEVMDeposit handles a deposit from an inbound tx
// returns (isContractReverted, err)
// (true, non-nil) means CallEVM() reverted
func (k Keeper) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) (bool, error) {
	to := ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver)
	var ethTxHash ethcommon.Hash
	inboundAmount := cctx.GetInboundTxParams().Amount.BigInt()
	inboundSender := cctx.GetInboundTxParams().Sender
	inboundSenderChainID := cctx.GetInboundTxParams().SenderChainId
	inboundCoinType := cctx.InboundTxParams.CoinType
	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash = ethcommon.BytesToHash(hash)
		cctx.GetCurrentOutTxParam().OutboundTxHash = ethTxHash.String()
		// #nosec G701 always positive
		cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = uint64(ctx.BlockHeight())
	}

	if inboundCoinType == coin.CoinType_Zeta {
		// if coin type is Zeta, this is a deposit ZETA to zEVM cctx.
		err := k.fungibleKeeper.DepositCoinZeta(ctx, to, inboundAmount)
		if err != nil {
			return false, err
		}
	} else {
		// cointype is Gas or ERC20; then it could be a ZRC20 deposit/depositAndCall cctx.
		parsedAddress, data, err := chains.ParseAddressAndData(cctx.RelayedMessage)
		if err != nil {
			return false, errors.Wrap(types.ErrUnableToParseAddress, err.Error())
		}
		if parsedAddress != (ethcommon.Address{}) {
			to = parsedAddress
		}

		from, err := chains.DecodeAddressFromChainID(inboundSenderChainID, inboundSender)
		if err != nil {
			return false, fmt.Errorf("HandleEVMDeposit: unable to decode address: %s", err.Error())
		}

		evmTxResponse, contractCall, err := k.fungibleKeeper.ZRC20DepositAndCallContract(
			ctx,
			from,
			to,
			inboundAmount,
			inboundSenderChainID,
			data,
			inboundCoinType,
			cctx.InboundTxParams.Asset,
		)
		if fungibletypes.IsContractReverted(evmTxResponse, err) || errShouldRevertCctx(err) {
			return true, err
		} else if err != nil {
			return false, err
		}

		// non-empty msg.Message means this is a contract call; therefore the logs should be processed.
		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if !evmTxResponse.Failed() && contractCall {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			if len(logs) > 0 {
				ctx = ctx.WithValue("inCctxIndex", cctx.Index)
				txOrigin := cctx.InboundTxParams.TxOrigin
				if txOrigin == "" {
					txOrigin = inboundSender
				}

				err = k.ProcessLogs(ctx, logs, to, txOrigin)
				if err != nil {
					// ProcessLogs should not error; error indicates exception, should abort
					return false, errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
				}
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(sdk.EventTypeMessage,
						sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
						sdk.NewAttribute("action", "DepositZRC20AndCallContract"),
						sdk.NewAttribute("contract", to.String()),
						sdk.NewAttribute("data", hex.EncodeToString(data)),
						sdk.NewAttribute("cctxIndex", cctx.Index),
					),
				)
			}
		}
	}
	return false, nil
}

// errShouldRevertCctx returns true if the cctx should revert from the error of the deposit
// we revert the cctx if a non-contract is tried to be called, if the liquidity cap is reached, or if the zrc20 is paused
func errShouldRevertCctx(err error) bool {
	return errors.Is(err, fungibletypes.ErrForeignCoinCapReached) ||
		errors.Is(err, fungibletypes.ErrCallNonContract) ||
		errors.Is(err, fungibletypes.ErrPausedZRC20)
}
