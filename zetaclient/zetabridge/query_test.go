package zetabridge

import (
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
	"net"
	"testing"
)

func setupMockServer(t *testing.T, serviceFunc any, method string, input any, expectedOutput any) *grpcmock.Server {
	l, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(serviceFunc),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(l),
		func(s *grpcmock.Server) {
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(input).
				Return(expectedOutput)
		},
	)(t)

	return server
}

func closeMockServer(t *testing.T, server *grpcmock.Server) {
	err := server.Close()
	require.NoError(t, err)
}

func setupCorBridge() (*ZetaCoreBridge, error) {
	return NewZetaCoreBridge(
		&keys.Keys{},
		"127.0.0.1",
		"",
		"zetachain_7000-1",
		false,
		&metrics.TelemetryServer{})
}

func TestZetaCoreBridge_GetBallot(t *testing.T) {
	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "456",
		Voters:           nil,
		ObservationType:  0,
		BallotStatus:     0,
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "123"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetBallotByID("123")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestZetaCoreBridge_GetCrosschainFlags(t *testing.T) {
	expectedOutput := observertypes.QueryGetCrosschainFlagsResponse{CrosschainFlags: observertypes.CrosschainFlags{
		IsInboundEnabled:             true,
		IsOutboundEnabled:            false,
		GasPriceIncreaseFlags:        nil,
		BlockHeaderVerificationFlags: nil,
	}}
	input := observertypes.QueryGetCrosschainFlagsRequest{}
	method := "/zetachain.zetacore.observer.Query/CrosschainFlags"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetCrosschainFlags()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrosschainFlags, resp)
}

func TestZetaCoreBridge_GetChainParamsForChainID(t *testing.T) {
	expectedOutput := observertypes.QueryGetChainParamsForChainResponse{ChainParams: &observertypes.ChainParams{
		ChainId: 123,
	}}
	input := observertypes.QueryGetChainParamsForChainRequest{ChainId: 123}
	method := "/zetachain.zetacore.observer.Query/GetChainParamsForChain"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetChainParamsForChainID(123)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams, resp)
}

func TestZetaCoreBridge_GetChainParams(t *testing.T) {
	expectedOutput := observertypes.QueryGetChainParamsResponse{ChainParams: &observertypes.ChainParamsList{
		ChainParams: []*observertypes.ChainParams{
			{
				ChainId: 123,
			},
		},
	}}
	input := observertypes.QueryGetChainParamsRequest{}
	method := "/zetachain.zetacore.observer.Query/GetChainParams"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetChainParams()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams.ChainParams, resp)
}

func TestZetaCoreBridge_GetUpgradePlan(t *testing.T) {
	expectedOutput := upgradetypes.QueryCurrentPlanResponse{
		Plan: &upgradetypes.Plan{
			Name:   "big upgrade",
			Height: 100,
		},
	}
	input := upgradetypes.QueryCurrentPlanRequest{}
	method := "/cosmos.upgrade.v1beta1.Query/CurrentPlan"
	server := setupMockServer(t, upgradetypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetUpgradePlan()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Plan, resp)
}

func TestZetaCoreBridge_GetAllCctx(t *testing.T) {
	expectedOutput := crosschainTypes.QueryAllCctxResponse{
		CrossChainTx: []*crosschainTypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		Pagination: nil,
	}
	input := crosschainTypes.QueryAllCctxRequest{}
	method := "/zetachain.zetacore.crosschain.Query/CctxAll"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetAllCctx()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetaCoreBridge_GetCctxByHash(t *testing.T) {
	expectedOutput := crosschainTypes.QueryGetCctxResponse{CrossChainTx: &crosschainTypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschainTypes.QueryGetCctxRequest{Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3"}
	method := "/zetachain.zetacore.crosschain.Query/Cctx"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetCctxByHash("9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetaCoreBridge_GetCctxByNonce(t *testing.T) {
	expectedOutput := crosschainTypes.QueryGetCctxResponse{CrossChainTx: &crosschainTypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschainTypes.QueryGetCctxByNonceRequest{
		ChainID: 7000,
		Nonce:   55,
	}
	method := "/zetachain.zetacore.crosschain.Query/CctxByNonce"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetCctxByNonce(7000, 55)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestZetaCoreBridge_GetObserverList(t *testing.T) {
	expectedOutput := observertypes.QueryObserverSetResponse{
		Observers: []string{
			"zeta19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
			"zeta1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj",
			"zeta1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
		},
	}
	input := observertypes.QueryObserverSet{}
	method := "/zetachain.zetacore.observer.Query/ObserverSet"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetObserverList()
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Observers, resp)
}

func TestZetaCoreBridge_ListPendingCctx(t *testing.T) {
	expectedOutput := crosschainTypes.QueryListCctxPendingResponse{
		CrossChainTx: []*crosschainTypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		TotalPending: 1,
	}
	input := crosschainTypes.QueryListCctxPendingRequest{ChainId: 7000}
	method := "/zetachain.zetacore.crosschain.Query/CctxListPending"
	server := setupMockServer(t, crosschainTypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, totalPending, err := zetabridge.ListPendingCctx(7000)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
	require.Equal(t, expectedOutput.TotalPending, totalPending)
}

func TestZetaCoreBridge_GetAbortedZetaAmount(t *testing.T) {
}

func TestZetaCoreBridge_GetGenesisSupply(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaTokenSupplyOnNode(t *testing.T) {
}

func TestZetaCoreBridge_GetLastBlockHeight(t *testing.T) {
}

func TestZetaCoreBridge_GetLatestZetaBlock(t *testing.T) {
}

func TestZetaCoreBridge_GetNodeInfo(t *testing.T) {
}

func TestZetaCoreBridge_GetLastBlockHeightByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaBlockHeight(t *testing.T) {
}

func TestZetaCoreBridge_GetBaseGasPrice(t *testing.T) {
}

func TestZetaCoreBridge_GetNonceByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetAllNodeAccounts(t *testing.T) {
}

func TestZetaCoreBridge_GetKeyGen(t *testing.T) {
}

func TestZetaCoreBridge_GetBallotByID(t *testing.T) {
}

func TestZetaCoreBridge_GetInboundTrackersForChain(t *testing.T) {
}

func TestZetaCoreBridge_GetCurrentTss(t *testing.T) {
}

func TestZetaCoreBridge_GetEthTssAddress(t *testing.T) {
}

func TestZetaCoreBridge_GetBtcTssAddress(t *testing.T) {
}

func TestZetaCoreBridge_GetTssHistory(t *testing.T) {
}

func TestZetaCoreBridge_GetOutTxTracker(t *testing.T) {
}

func TestZetaCoreBridge_GetAllOutTxTrackerByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetPendingNoncesByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetBlockHeaderStateByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetSupportedChains(t *testing.T) {
}

func TestZetaCoreBridge_GetPendingNonces(t *testing.T) {
}

func TestZetaCoreBridge_Prove(t *testing.T) {
}

func TestZetaCoreBridge_HasVoted(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaHotKeyBalance(t *testing.T) {
}
