package zetabridge

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"google.golang.org/grpc"
)

func (b *ZetaCoreBridge) GetCrosschainFlags() (observertypes.CrosschainFlags, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.CrosschainFlags(context.Background(), &observertypes.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return observertypes.CrosschainFlags{}, err
	}
	return resp.CrosschainFlags, nil
}

func (b *ZetaCoreBridge) GetVerificationFlags() (lightclienttypes.VerificationFlags, error) {
	client := lightclienttypes.NewQueryClient(b.grpcConn)
	resp, err := client.VerificationFlags(context.Background(), &lightclienttypes.QueryVerificationFlagsRequest{})
	if err != nil {
		return lightclienttypes.VerificationFlags{}, err
	}
	return resp.VerificationFlags, nil
}

func (b *ZetaCoreBridge) GetChainParamsForChainID(externalChainID int64) (*observertypes.ChainParams, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.GetChainParamsForChain(context.Background(), &observertypes.QueryGetChainParamsForChainRequest{ChainId: externalChainID})
	if err != nil {
		return &observertypes.ChainParams{}, err
	}
	return resp.ChainParams, nil
}

func (b *ZetaCoreBridge) GetChainParams() ([]*observertypes.ChainParams, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	var err error

	resp := &observertypes.QueryGetChainParamsResponse{}
	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err = client.GetChainParams(context.Background(), &observertypes.QueryGetChainParamsRequest{})
		if err == nil {
			return resp.ChainParams.ChainParams, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get chain params | err %s", err.Error())
}

func (b *ZetaCoreBridge) GetUpgradePlan() (*upgradetypes.Plan, error) {
	client := upgradetypes.NewQueryClient(b.grpcConn)

	resp, err := client.CurrentPlan(context.Background(), &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Plan, nil
}

func (b *ZetaCoreBridge) GetAllCctx() ([]*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.CctxAll(context.Background(), &types.QueryAllCctxRequest{})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetCctxByHash(sendHash string) (*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.Cctx(context.Background(), &types.QueryGetCctxRequest{Index: sendHash})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetCctxByNonce(chainID int64, nonce uint64) (*types.CrossChainTx, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.CctxByNonce(context.Background(), &types.QueryGetCctxByNonceRequest{
		ChainID: chainID,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

func (b *ZetaCoreBridge) GetObserverList() ([]string, error) {
	var err error
	client := observertypes.NewQueryClient(b.grpcConn)

	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err := client.ObserverSet(context.Background(), &observertypes.QueryObserverSet{})
		if err == nil {
			return resp.Observers, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

// ListPendingCctx returns a list of pending cctxs for a given chainID
// the returned list has a limited size of crosschainkeeper.MaxPendingCctxs
// the total number of pending cctxs is returned
func (b *ZetaCoreBridge) ListPendingCctx(chainID int64) ([]*types.CrossChainTx, uint64, error) {
	client := types.NewQueryClient(b.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.CctxListPending(context.Background(), &types.QueryListCctxPendingRequest{ChainId: chainID}, maxSizeOption)
	if err != nil {
		return nil, 0, err
	}
	return resp.CrossChainTx, resp.TotalPending, nil
}

func (b *ZetaCoreBridge) GetAbortedZetaAmount() (string, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.ZetaAccounting(context.Background(), &types.QueryZetaAccountingRequest{})
	if err != nil {
		return "", err
	}
	return resp.AbortedZetaAmount, nil
}

func (b *ZetaCoreBridge) GetGenesisSupply() (sdkmath.Int, error) {
	tmURL := fmt.Sprintf("http://%s", b.cfg.ChainRPC)
	s, err := tmhttp.New(tmURL, "/websocket")
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	res, err := s.Genesis(context.Background())
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	appState, err := genutiltypes.GenesisStateFromGenDoc(*res.Genesis)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	bankstate := banktypes.GetGenesisStateFromAppState(b.encodingCfg.Codec, appState)
	return bankstate.Supply.AmountOf(config.BaseDenom), nil
}

func (b *ZetaCoreBridge) GetZetaTokenSupplyOnNode() (sdkmath.Int, error) {
	client := banktypes.NewQueryClient(b.grpcConn)
	resp, err := client.SupplyOf(context.Background(), &banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	return resp.GetAmount().Amount, nil
}

func (b *ZetaCoreBridge) GetLastBlockHeight() ([]*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &types.QueryAllLastBlockHeightRequest{})
	if err != nil {
		b.logger.Error().Err(err).Msg("query GetBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetLatestZetaBlock() (*tmservice.Block, error) {
	client := tmservice.NewServiceClient(b.grpcConn)
	res, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return res.SdkBlock, nil
}

func (b *ZetaCoreBridge) GetNodeInfo() (*tmservice.GetNodeInfoResponse, error) {
	var err error

	client := tmservice.NewServiceClient(b.grpcConn)
	for i := 0; i <= DefaultRetryCount; i++ {
		res, err := client.GetNodeInfo(context.Background(), &tmservice.GetNodeInfoRequest{})
		if err == nil {
			return res, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

func (b *ZetaCoreBridge) GetLastBlockHeightByChain(chain chains.Chain) (*types.LastBlockHeight, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastBlockHeight(context.Background(), &types.QueryGetLastBlockHeightRequest{Index: chain.ChainName.String()})
	if err != nil {
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

func (b *ZetaCoreBridge) GetZetaBlockHeight() (int64, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.LastZetaHeight(context.Background(), &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Height, nil
}

func (b *ZetaCoreBridge) GetBaseGasPrice() (int64, error) {
	client := feemarkettypes.NewQueryClient(b.grpcConn)
	resp, err := client.Params(context.Background(), &feemarkettypes.QueryParamsRequest{})
	if err != nil {
		return 0, err
	}
	if resp.Params.BaseFee.IsNil() {
		return 0, fmt.Errorf("base fee is nil")
	}
	return resp.Params.BaseFee.Int64(), nil
}

func (b *ZetaCoreBridge) GetBallotByID(id string) (*observertypes.QueryBallotByIdentifierResponse, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	return client.BallotByIdentifier(context.Background(), &observertypes.QueryBallotByIdentifierRequest{
		BallotIdentifier: id,
	})
}

func (b *ZetaCoreBridge) GetNonceByChain(chain chains.Chain) (observertypes.ChainNonces, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.ChainNonces(context.Background(), &observertypes.QueryGetChainNoncesRequest{Index: chain.ChainName.String()})
	if err != nil {
		return observertypes.ChainNonces{}, err
	}
	return resp.ChainNonces, nil
}

func (b *ZetaCoreBridge) GetAllNodeAccounts() ([]*observertypes.NodeAccount, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.NodeAccountAll(context.Background(), &observertypes.QueryAllNodeAccountRequest{})
	if err != nil {
		return nil, err
	}
	b.logger.Debug().Msgf("GetAllNodeAccounts: %d", len(resp.NodeAccount))
	return resp.NodeAccount, nil
}

func (b *ZetaCoreBridge) GetKeyGen() (*observertypes.Keygen, error) {
	var err error
	client := observertypes.NewQueryClient(b.grpcConn)

	for i := 0; i <= ExtendedRetryCount; i++ {
		resp, err := client.Keygen(context.Background(), &observertypes.QueryGetKeygenRequest{})
		if err == nil {
			return resp.Keygen, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get keygen | err %s", err.Error())

}

func (b *ZetaCoreBridge) GetBallot(ballotIdentifier string) (*observertypes.QueryBallotByIdentifierResponse, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.BallotByIdentifier(context.Background(), &observertypes.QueryBallotByIdentifierRequest{
		BallotIdentifier: ballotIdentifier,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *ZetaCoreBridge) GetInboundTrackersForChain(chainID int64) ([]types.InTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.InTxTrackerAllByChain(context.Background(), &types.QueryAllInTxTrackerByChainRequest{ChainId: chainID})
	if err != nil {
		return nil, err
	}
	return resp.InTxTracker, nil
}

func (b *ZetaCoreBridge) GetCurrentTss() (observertypes.TSS, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.TSS(context.Background(), &observertypes.QueryGetTSSRequest{})
	if err != nil {
		return observertypes.TSS{}, err
	}
	return resp.TSS, nil
}

func (b *ZetaCoreBridge) GetEthTssAddress() (string, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &observertypes.QueryGetTssAddressRequest{})
	if err != nil {
		return "", err
	}
	return resp.Eth, nil
}

func (b *ZetaCoreBridge) GetBtcTssAddress(chainID int64) (string, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &observertypes.QueryGetTssAddressRequest{
		BitcoinChainId: chainID,
	})
	if err != nil {
		return "", err
	}
	return resp.Btc, nil
}

func (b *ZetaCoreBridge) GetTssHistory() ([]observertypes.TSS, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.TssHistory(context.Background(), &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return nil, err
	}
	return resp.TssList, nil
}

func (b *ZetaCoreBridge) GetOutTxTracker(chain chains.Chain, nonce uint64) (*types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTracker(context.Background(), &types.QueryGetOutTxTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return &resp.OutTxTracker, nil
}

func (b *ZetaCoreBridge) GetAllOutTxTrackerByChain(chainID int64, order interfaces.Order) ([]types.OutTxTracker, error) {
	client := types.NewQueryClient(b.grpcConn)
	resp, err := client.OutTxTrackerAllByChain(context.Background(), &types.QueryAllOutTxTrackerByChainRequest{
		Chain: chainID,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	})
	if err != nil {
		return nil, err
	}
	if order == interfaces.Ascending {
		sort.SliceStable(resp.OutTxTracker, func(i, j int) bool {
			return resp.OutTxTracker[i].Nonce < resp.OutTxTracker[j].Nonce
		})
	}
	if order == interfaces.Descending {
		sort.SliceStable(resp.OutTxTracker, func(i, j int) bool {
			return resp.OutTxTracker[i].Nonce > resp.OutTxTracker[j].Nonce
		})
	}
	return resp.OutTxTracker, nil
}

func (b *ZetaCoreBridge) GetPendingNoncesByChain(chainID int64) (observertypes.PendingNonces, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.PendingNoncesByChain(context.Background(), &observertypes.QueryPendingNoncesByChainRequest{ChainId: chainID})
	if err != nil {
		return observertypes.PendingNonces{}, err
	}
	return resp.PendingNonces, nil
}

func (b *ZetaCoreBridge) GetBlockHeaderChainState(chainID int64) (lightclienttypes.QueryGetChainStateResponse, error) {
	client := lightclienttypes.NewQueryClient(b.grpcConn)
	resp, err := client.ChainState(context.Background(), &lightclienttypes.QueryGetChainStateRequest{ChainId: chainID})
	if err != nil {
		return lightclienttypes.QueryGetChainStateResponse{}, err
	}
	return *resp, nil
}

func (b *ZetaCoreBridge) GetSupportedChains() ([]*chains.Chain, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.SupportedChains(context.Background(), &observertypes.QuerySupportedChains{})
	if err != nil {
		return nil, err
	}
	return resp.GetChains(), nil
}

func (b *ZetaCoreBridge) GetPendingNonces() (*observertypes.QueryAllPendingNoncesResponse, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.PendingNoncesAll(context.Background(), &observertypes.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *ZetaCoreBridge) Prove(blockHash string, txHash string, txIndex int64, proof *proofs.Proof, chainID int64) (bool, error) {
	client := lightclienttypes.NewQueryClient(b.grpcConn)
	resp, err := client.Prove(context.Background(), &lightclienttypes.QueryProveRequest{
		BlockHash: blockHash,
		TxIndex:   txIndex,
		Proof:     proof,
		ChainId:   chainID,
		TxHash:    txHash,
	})
	if err != nil {
		return false, err
	}
	return resp.Valid, nil
}

func (b *ZetaCoreBridge) HasVoted(ballotIndex string, voterAddress string) (bool, error) {
	client := observertypes.NewQueryClient(b.grpcConn)
	resp, err := client.HasVoted(context.Background(), &observertypes.QueryHasVotedRequest{
		BallotIdentifier: ballotIndex,
		VoterAddress:     voterAddress,
	})
	if err != nil {
		return false, err
	}
	return resp.HasVoted, nil
}

func (b *ZetaCoreBridge) GetZetaHotKeyBalance() (sdkmath.Int, error) {
	client := banktypes.NewQueryClient(b.grpcConn)
	resp, err := client.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: b.keys.GetAddress().String(),
		Denom:   config.BaseDenom,
	})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	return resp.Balance.Amount, nil
}
