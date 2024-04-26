package network

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func SetupZetaGenesisState(t *testing.T, genesisState map[string]json.RawMessage, codec codec.Codec, observerList []string, setupChainNonces bool) {

	// Cross-chain genesis state
	var crossChainGenesis types.GenesisState
	require.NoError(t, codec.UnmarshalJSON(genesisState[types.ModuleName], &crossChainGenesis))
	nodeAccountList := make([]*observertypes.NodeAccount, len(observerList))
	for i, operator := range observerList {
		nodeAccountList[i] = &observertypes.NodeAccount{
			Operator:   operator,
			NodeStatus: observertypes.NodeStatus_Active,
		}
	}

	require.NoError(t, crossChainGenesis.Validate())
	crossChainGenesisBz, err := codec.MarshalJSON(&crossChainGenesis)
	require.NoError(t, err)

	// EVM genesis state
	var evmGenesisState evmtypes.GenesisState
	require.NoError(t, codec.UnmarshalJSON(genesisState[evmtypes.ModuleName], &evmGenesisState))
	evmGenesisState.Params.EvmDenom = cmdcfg.BaseDenom
	require.NoError(t, evmGenesisState.Validate())
	evmGenesisBz, err := codec.MarshalJSON(&evmGenesisState)
	require.NoError(t, err)

	// Staking genesis state
	var stakingGenesisState stakingtypes.GenesisState
	require.NoError(t, codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesisState))
	stakingGenesisState.Params.BondDenom = cmdcfg.BaseDenom
	stakingGenesisStateBz, err := codec.MarshalJSON(&stakingGenesisState)
	require.NoError(t, err)

	// Observer genesis state
	var observerGenesis observertypes.GenesisState
	require.NoError(t, codec.UnmarshalJSON(genesisState[observertypes.ModuleName], &observerGenesis))
	observerSet := observertypes.ObserverSet{
		ObserverList: observerList,
	}

	if setupChainNonces {
		privatenetChains := chains.ChainListByNetworkType(chains.NetworkType_privnet)
		chainNonceList := make([]observertypes.ChainNonces, len(privatenetChains))
		for i, chain := range privatenetChains {
			chainNonceList[i] = observertypes.ChainNonces{
				Index:   chain.ChainName.String(),
				ChainId: chain.ChainId,
				Nonce:   0,
			}
		}
		observerGenesis.ChainNonces = chainNonceList
	}

	observerGenesis.Observers = observerSet
	observerGenesis.NodeAccountList = nodeAccountList
	observerGenesis.Keygen = &observertypes.Keygen{
		Status:         observertypes.KeygenStatus_PendingKeygen,
		GranteePubkeys: observerList,
		BlockNumber:    5,
	}
	observerGenesis.CrosschainFlags = &observertypes.CrosschainFlags{
		IsInboundEnabled:  true,
		IsOutboundEnabled: true,
	}
	require.NoError(t, observerGenesis.Validate())
	observerGenesisBz, err := codec.MarshalJSON(&observerGenesis)
	require.NoError(t, err)

	// authority genesis state
	var authorityGenesis authoritytypes.GenesisState
	require.NoError(t, codec.UnmarshalJSON(genesisState[authoritytypes.ModuleName], &authorityGenesis))
	policies := authoritytypes.Policies{
		Items: []*authoritytypes.Policy{
			{
				PolicyType: authoritytypes.PolicyType_groupEmergency,
				Address:    "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
			},
			{
				PolicyType: authoritytypes.PolicyType_groupAdmin,
				Address:    "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
			},
			{
				PolicyType: authoritytypes.PolicyType_groupOperational,
				Address:    "zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
			},
		},
	}
	authorityGenesis.Policies = policies
	require.NoError(t, authorityGenesis.Validate())
	authorityGenesisBz, err := codec.MarshalJSON(&authorityGenesis)
	require.NoError(t, err)

	genesisState[types.ModuleName] = crossChainGenesisBz
	genesisState[stakingtypes.ModuleName] = stakingGenesisStateBz
	genesisState[observertypes.ModuleName] = observerGenesisBz
	genesisState[evmtypes.ModuleName] = evmGenesisBz
	genesisState[authoritytypes.ModuleName] = authorityGenesisBz
}

func AddObserverData(t *testing.T, n int, genesisState map[string]json.RawMessage, codec codec.Codec, ballots []*observertypes.Ballot) *observertypes.GenesisState {
	state := observertypes.GenesisState{}
	require.NoError(t, codec.UnmarshalJSON(genesisState[observertypes.ModuleName], &state))

	// set chain params with chains all enabled
	state.ChainParamsList = observertypes.GetDefaultChainParams()
	for i := range state.ChainParamsList.ChainParams {
		state.ChainParamsList.ChainParams[i].IsSupported = true
	}

	// set params
	if len(ballots) > 0 {
		state.Ballots = ballots
	}
	state.Params.BallotMaturityBlocks = 3
	state.Keygen = &observertypes.Keygen{BlockNumber: 10, GranteePubkeys: []string{}}

	// set tss
	tss := observertypes.TSS{
		TssPubkey:           "tssPubkey",
		TssParticipantList:  []string{"tssParticipantList"},
		OperatorAddressList: []string{"operatorAddressList"},
		FinalizedZetaHeight: 1,
		KeyGenZetaHeight:    1,
	}
	pendingNonces := make([]observertypes.PendingNonces, len(chains.DefaultChainsList()))
	for i, chain := range chains.DefaultChainsList() {
		pendingNonces[i] = observertypes.PendingNonces{
			ChainId:   chain.ChainId,
			NonceLow:  0,
			NonceHigh: 0,
			Tss:       tss.TssPubkey,
		}
	}
	state.Tss = &tss
	state.TssHistory = []observertypes.TSS{tss}
	state.PendingNonces = pendingNonces

	// set crosschain flags
	crosschainFlags := &observertypes.CrosschainFlags{
		IsInboundEnabled:             true,
		IsOutboundEnabled:            true,
		GasPriceIncreaseFlags:        &observertypes.DefaultGasPriceIncreaseFlags,
		BlockHeaderVerificationFlags: &observertypes.DefaultBlockHeaderVerificationFlags,
	}
	nullify.Fill(&crosschainFlags)
	state.CrosschainFlags = crosschainFlags

	for i := 0; i < n; i++ {
		state.ChainNonces = append(state.ChainNonces, observertypes.ChainNonces{Creator: "ANY", Index: strconv.Itoa(i), Signers: []string{}})
	}

	// check genesis state validity
	require.NoError(t, state.Validate())

	// marshal genesis state
	buf, err := codec.MarshalJSON(&state)
	require.NoError(t, err)
	genesisState[observertypes.ModuleName] = buf
	return &state
}
func AddCrosschainData(t *testing.T, n int, genesisState map[string]json.RawMessage, codec codec.Codec) *types.GenesisState {
	state := types.GenesisState{}
	require.NoError(t, codec.UnmarshalJSON(genesisState[types.ModuleName], &state))
	// TODO : Fix add EVM balance to deploy contracts
	for i := 0; i < n; i++ {
		state.CrossChainTxs = append(state.CrossChainTxs, &types.CrossChainTx{
			Creator: "ANY",
			Index:   strconv.Itoa(i),
			CctxStatus: &types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "",
				LastUpdateTimestamp: 0,
			},
			InboundTxParams:  &types.InboundTxParams{InboundTxObservedHash: fmt.Sprintf("Hash-%d", i), Amount: math.OneUint()},
			OutboundTxParams: []*types.OutboundTxParams{},
			ZetaFees:         math.OneUint()},
		)
	}

	for i := 0; i < n; i++ {
		state.GasPriceList = append(state.GasPriceList, &types.GasPrice{Creator: "ANY", ChainId: int64(i), Index: strconv.Itoa(i), Prices: []uint64{}, BlockNums: []uint64{}, Signers: []string{}})
	}
	for i := 0; i < n; i++ {
		state.LastBlockHeightList = append(state.LastBlockHeightList, &types.LastBlockHeight{Creator: "ANY", Index: strconv.Itoa(i)})
	}

	for i := 0; i < n; i++ {
		inTxTracker := types.InTxTracker{
			ChainId:  5,
			TxHash:   fmt.Sprintf("txHash-%d", i),
			CoinType: coin.CoinType_Gas,
		}
		nullify.Fill(&inTxTracker)
		state.InTxTrackerList = append(state.InTxTrackerList, inTxTracker)
	}

	for i := 0; i < n; i++ {
		inTxHashToCctx := types.InTxHashToCctx{
			InTxHash: strconv.Itoa(i),
		}
		nullify.Fill(&inTxHashToCctx)
		state.InTxHashToCctxList = append(state.InTxHashToCctxList, inTxHashToCctx)
	}
	for i := 0; i < n; i++ {
		outTxTracker := types.OutTxTracker{
			Index:   fmt.Sprintf("%d-%d", i, i),
			ChainId: int64(i),
			Nonce:   uint64(i),
		}
		nullify.Fill(&outTxTracker)
		state.OutTxTrackerList = append(state.OutTxTrackerList, outTxTracker)
	}

	require.NoError(t, state.Validate())

	buf, err := codec.MarshalJSON(&state)
	require.NoError(t, err)
	genesisState[types.ModuleName] = buf
	return &state
}
