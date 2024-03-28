package smoketests

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestAdjustZetaSupply(sm *runner.SmokeTestRunner) {
	sm.Logger.Print("🏃 starting burn zeta tokens test")

	chainID, err := sm.GoerliClient.NetworkID(sm.Ctx)
	if err != nil {
		panic(err)
	}
	amount := sdkmath.NewUint(1000000000)
	// get token balances before reducing supply
	connectorBalance, burnAddressBalance := GetTokenBalances(sm)

	msg := crosschaintypes.NewMsgReduceZetaSupply(
		sm.ZetaTxServer.GetAccountAddress(0),
		chainID.Int64(),
		amount,
		"",
	)

	res, err := sm.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	cctxIndex := crosschainkeeper.GetIndexForReduceZetaSupplyCMD(chainID.Int64(), res.Height)

	// wait for the cctx to be mined on the zeta chain through the outbound votes from obsevers
	sm.WaitForMinedCCTXfromIndex(cctxIndex)

	// get the cctx from the zeta chain and assert values are correct
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
			sm.Logger.Print("  Zeta Value: %d", event.ZetaValue)
			sm.Logger.Info("  src ChainId: %d", event.SourceChainId)
			if event.DestinationAddress.String() != crosschaintypes.BurnTokensEVMAddress {
				panic("Destination address mismatch")
			}
			if event.ZetaValue.Cmp(cctxRes.CrossChainTx.GetCurrentOutTxParam().Amount.BigInt()) != 0 {
				panic("Zeta value mismatch")
			}
		}
	}
	// get token balances after reducing supply and compare balances
	connectorBalancesAfter, burnAddressBalanceAfter := GetTokenBalances(sm)
	connectorBalanceDiff := new(big.Int).Sub(connectorBalance, connectorBalancesAfter)
	burnAddressBalanceDiff := new(big.Int).Sub(burnAddressBalanceAfter, burnAddressBalance)
	if connectorBalanceDiff.Cmp(amount.BigInt()) != 0 {
		panic("Connector balance mismatch")
	}
	if burnAddressBalanceDiff.Cmp(amount.BigInt()) != 0 {
		panic("Burn address balance mismatch")
	}
}

func GetTokenBalances(sm *runner.SmokeTestRunner) (*big.Int, *big.Int) {
	tokenBalanceAfter, err := sm.ZetaEth.BalanceOf(&bind.CallOpts{}, sm.ConnectorEthAddr)
	if err != nil {
		panic(err)
	}
	burnAddress := ethcommon.HexToAddress(crosschaintypes.BurnTokensEVMAddress)
	burnAddressBalanceAfter, err := sm.ZetaEth.BalanceOf(&bind.CallOpts{}, burnAddress)
	if err != nil {
		panic(err)
	}
	return tokenBalanceAfter, burnAddressBalanceAfter
}
