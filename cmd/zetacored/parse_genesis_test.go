package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/app"
	zetacored "github.com/zeta-chain/zetacore/cmd/zetacored"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setCosmosConfig() {
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()
}
func Test_ModifyCrossChainState(t *testing.T) {
	setCosmosConfig()
	t.Run("successfully modify cross chain state to reduce data", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState()
		importData := GetImportData(t, cdc, 100)
		err := zetacored.ModifyCrossChainState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedCrosschainAppState := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedCrosschainAppState.CrossChainTxs, 10)
		require.Len(t, modifiedCrosschainAppState.InTxHashToCctxList, 10)
		require.Len(t, modifiedCrosschainAppState.FinalizedInbounds, 10)
	})

	t.Run("successfully modify cross chain state without changing data when not needed", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState()
		importData := GetImportData(t, cdc, 8)
		err := zetacored.ModifyCrossChainState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedCrosschainAppState := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedCrosschainAppState.CrossChainTxs, 8)
		require.Len(t, modifiedCrosschainAppState.InTxHashToCctxList, 8)
		require.Len(t, modifiedCrosschainAppState.FinalizedInbounds, 8)
	})
}

func GetImportData(t *testing.T, cdc *codec.ProtoCodec, n int) map[string]json.RawMessage {
	importData := sample.AppState()

	importedCrossChainGenState := crosschaintypes.GetGenesisStateFromAppState(cdc, importData)
	cctxList := make([]*crosschaintypes.CrossChainTx, n)
	intxHashToCctxList := make([]crosschaintypes.InTxHashToCctx, n)
	finalLizedInbounds := make([]string, n)
	for i := 0; i < n; i++ {
		cctxList[i] = sample.CrossChainTx(t, fmt.Sprintf("crosschain-%d", i))
		intxHashToCctxList[i] = sample.InTxHashToCctx(t, fmt.Sprintf("intxHashToCctxList-%d", i))
		finalLizedInbounds[i] = fmt.Sprintf("finalLizedInbounds-%d", i)
	}
	importedCrossChainGenState.CrossChainTxs = cctxList
	importedCrossChainGenState.InTxHashToCctxList = intxHashToCctxList
	importedCrossChainGenState.FinalizedInbounds = finalLizedInbounds
	importedCrossChainStateBz, err := json.Marshal(importedCrossChainGenState)
	require.NoError(t, err)
	importData[crosschaintypes.ModuleName] = importedCrossChainStateBz

	importedObserverGenState := observertypes.GetGenesisStateFromAppState(cdc, importData)
	ballots := make([]*observertypes.Ballot, n)
	nonceToCctx := make([]observertypes.NonceToCctx, n)
	for i := 0; i < n; i++ {
		ballots[i] = sample.Ballot(t, fmt.Sprintf("ballots-%d", i))
		nonceToCctx[i] = sample.NonceToCCTX(t, fmt.Sprintf("nonceToCctx-%d", i))
	}
	importedObserverGenState.Ballots = ballots
	importedObserverGenState.NonceToCctx = nonceToCctx
	importedObserverStateBz, err := cdc.MarshalJSON(&importedObserverGenState)
	require.NoError(t, err)
	importData[observertypes.ModuleName] = importedObserverStateBz

	return importData
}
