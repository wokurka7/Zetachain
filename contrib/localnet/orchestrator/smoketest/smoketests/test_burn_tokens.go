package smoketests

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestBurnZetaTokens(sm *runner.SmokeTestRunner) {
	sm.Logger.Print("🏃 starting burn zeta tokens test")

	connectorBalance, _ := sm.GoerliClient.BalanceAt(sm.Ctx, sm.ConnectorEthAddr, nil)
	sm.Logger.Print("Connector balance: %s", connectorBalance.String())

	chainID, err := sm.GoerliClient.NetworkID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	sm.Logger.Print("Chain ID: %d", chainID.Int64())

	msg := crosschaintypes.NewMsgBurnTokens(
		sm.ZetaTxServer.GetAccountAddress(0),
		chainID.Int64(),
		sdkmath.NewUint(1000000000),
	)
	res, err := sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	cctxIndex := crosschainkeeper.GetIndexForBurnCCTX(chainID.Int64(), res.Height)
	sm.Logger.Print(fmt.Sprintf("Burn tx hash: %s, cctx index: %s ,height %d", res.TxHash, cctxIndex, res.Height))
}
