package zetabridge

import (
	"encoding/hex"
	"errors"
	"net"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
)

func TestHandleBroadcastError(t *testing.T) {
	type response struct {
		retry  bool
		report bool
	}
	testCases := map[error]response{
		errors.New("nonce too low"):                       {retry: false, report: false},
		errors.New("replacement transaction underpriced"): {retry: false, report: false},
		errors.New("already known"):                       {retry: false, report: true},
		errors.New(""):                                    {retry: true, report: false},
	}
	for input, output := range testCases {
		retry, report := HandleBroadcastError(input, "", "", "")
		require.Equal(t, output.report, report)
		require.Equal(t, output.retry, retry)
	}
}

func TestBroadcast(t *testing.T) {
	address := types.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())

	//Setup server for multiple grpc calls
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)
	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(feemarkettypes.RegisterQueryServer),
		grpcmock.RegisterService(authtypes.RegisterQueryServer),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
		func(s *grpcmock.Server) {
			method := "/zetachain.zetacore.crosschain.Query/LastZetaHeight"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(crosschaintypes.QueryLastZetaHeightRequest{}).
				Return(crosschaintypes.QueryLastZetaHeightResponse{Height: 0})

			method = "/ethermint.feemarket.v1.Query/Params"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(feemarkettypes.QueryParamsRequest{}).
				Return(feemarkettypes.QueryParamsResponse{
					Params: feemarkettypes.Params{
						BaseFee: types.NewInt(23455),
					},
				})
		},
	)(t)

	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	t.Run("broadcast success", func(t *testing.T) {
		zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(nil, 0))
		blockHash, err := hex.DecodeString(ethBlockHash)
		require.NoError(t, err)
		msg := observerTypes.NewMsgVoteBlockHeader(address.String(), chains.EthChain.ChainId, blockHash, 18495266, getHeaderData(t))
		authzMsg, authzSigner, err := zetabridge.WrapMessageWithAuthz(msg)
		require.NoError(t, err)
		_, err = BroadcastToZetaCore(zetabridge, 10000, authzMsg, authzSigner)
		require.NoError(t, err)
	})

	t.Run("broadcast failed", func(t *testing.T) {
		zetabridge.EnableMockSDKClient(stub.NewSDKClientWithErr(errors.New("account sequence mismatch, expected 5 got 4"), 32))
		blockHash, err := hex.DecodeString(ethBlockHash)
		require.NoError(t, err)
		msg := observerTypes.NewMsgVoteBlockHeader(address.String(), chains.EthChain.ChainId, blockHash, 18495266, getHeaderData(t))
		authzMsg, authzSigner, err := zetabridge.WrapMessageWithAuthz(msg)
		require.NoError(t, err)
		_, err = BroadcastToZetaCore(zetabridge, 10000, authzMsg, authzSigner)
		require.Error(t, err)
	})

}

func TestZetaCoreBridge_GetContext(t *testing.T) {
	address := types.AccAddress(stub.TestKeyringPair.PubKey().Address().Bytes())
	zetabridge, err := setupCoreBridge()
	require.NoError(t, err)
	zetabridge.keys = keys.NewKeysWithKeybase(stub.NewKeyring(), address, "", "")

	_, err = zetabridge.GetContext()
	require.NoError(t, err)

}
