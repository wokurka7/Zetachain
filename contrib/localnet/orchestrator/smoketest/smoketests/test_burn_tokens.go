package smoketests

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestBurnZetaTokens(sm *runner.SmokeTestRunner) {
	sm.Logger.Print("🏃 starting burn zeta tokens test")

	chainID, err := sm.GoerliClient.NetworkID(sm.Ctx)
	if err != nil {
		panic(err)
	}

	msg := crosschaintypes.NewMsgReduceZetaSupply(
		sm.ZetaTxServer.GetAccountAddress(0),
		chainID.Int64(),
		sdkmath.NewUint(1000000000),
		"",
	)
	res, err := sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	cctxIndex := crosschainkeeper.GetIndexForReduceZetaSupplyCMD(chainID.Int64(), res.Height)

	sm.Logger.Info(fmt.Sprintf("Burn tx hash: %s, cctx index: %s ,height %d", res.TxHash, cctxIndex, res.Height))
	sm.WaitForMinedCCTXfromIndex(cctxIndex)

	cctxRes, err := sm.CctxClient.Cctx(sm.Ctx, &crosschaintypes.QueryGetCctxRequest{
		Index: cctxIndex,
	})
	if err != nil {
		panic(err)
	}

	receipt, err := sm.GoerliClient.TransactionReceipt(sm.Ctx, ethcommon.HexToHash(cctxRes.CrossChainTx.GetCurrentOutTxParam().OutboundTxHash))
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
	for _, log := range receipt.Logs {
		event, err := sm.ConnectorEth.ParseZetaReceived(*log)
		if err == nil {
			sm.Logger.Info("Received ZetaSent event:")
			sm.Logger.Info("  Dest Addr: %s", event.DestinationAddress)
			sm.Logger.Info("  Zeta Value: %d", event.ZetaValue)
			sm.Logger.Info("  src ChainId: %d", event.SourceChainId)
			if event.DestinationAddress.String() != crosschaintypes.BurnTokensEVMAddress {
				panic("Destination address mismatch")
			}
			if event.ZetaValue.Cmp(cctxRes.CrossChainTx.GetCurrentOutTxParam().Amount.BigInt()) != 0 {
				panic("Zeta value mismatch")
			}
		}
	}

}
