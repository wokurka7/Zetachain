package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	mocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/crosschain"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func getEthereumChainID() int64 {
	return 5 // Goerli
}

// TODO: Add a test case with proof and Bitcoin chain
// https://github.com/zeta-chain/node/issues/1994
func TestMsgServer_AddToOutTxTracker(t *testing.T) {

	// mockCctxByNonce mocks the methods called by CctxByNonce to directly return the given cctx or error
	mockCctxByNonce := func(
		t *testing.T,
		ctx sdk.Context,
		k keeper.Keeper,
		observerKeeper *mocks.CrosschainObserverKeeper,
		cctxStatus types.CctxStatus,
		isErr bool,
	) {
		if isErr {
			// return error on GetTSS to make CctxByNonce return error
			observerKeeper.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false).Once()
			return
		}

		cctx := sample.CrossChainTx(t, sample.StringRandom(sample.Rand(), 10))
		cctx.CctxStatus = &types.Status{
			Status: cctxStatus,
		}
		k.SetCrossChainTx(ctx, *cctx)

		observerKeeper.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true).Once()
		observerKeeper.On("GetNonceToCctx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(observertypes.NonceToCctx{
			CctxIndex: cctx.Index,
		}, true).Once()
	}

	t.Run("admin can add tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		chainID := getEthereumChainID()
		hash := sample.Hash().Hex()

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    hash,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.True(t, found)
		require.Equal(t, hash, tracker.HashList[0].TxHash)
	})

	t.Run("observer can add tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(true)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		chainID := getEthereumChainID()
		hash := sample.Hash().Hex()

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    hash,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.True(t, found)
		require.Equal(t, hash, tracker.HashList[0].TxHash)
	})

	t.Run("can add hash to existing tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		chainID := getEthereumChainID()
		existinghHash := sample.Hash().Hex()
		newHash := sample.Hash().Hex()

		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   42,
			HashList: []*types.TxHashList{
				{
					TxHash: existinghHash,
				},
			},
		})

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    newHash,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 42)
		require.True(t, found)
		require.Len(t, tracker.HashList, 2)
		require.EqualValues(t, existinghHash, tracker.HashList[0].TxHash)
		require.EqualValues(t, newHash, tracker.HashList[1].TxHash)
	})

	t.Run("should return early if cctx not pending", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})

		// set cctx status to outbound mined
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_OutboundMined, false)

		chainID := getEthereumChainID()

		res, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		require.Equal(t, &types.MsgAddToOutTxTrackerResponse{IsRemoved: true}, res)

		// check if tracker is removed
		_, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.False(t, found)
	})

	t.Run("should error for unsupported chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(nil)

		// set cctx status to outbound mined
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_OutboundMined, false)

		chainID := getEthereumChainID()

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})

	t.Run("should error if no CctxByNonce", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, true)

		chainID := getEthereumChainID()

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.ErrorIs(t, err, types.ErrCannotFindCctx)
	})

	t.Run("should fail if max tracker hashes reached", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		hashes := make([]*types.TxHashList, keeper.MaxOutTxTrackerHashes)
		for i := 0; i < keeper.MaxOutTxTrackerHashes; i++ {
			hashes[i] = &types.TxHashList{
				TxHash: sample.Hash().Hex(),
			}
		}

		chainID := getEthereumChainID()
		newHash := sample.Hash().Hex()

		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId:  chainID,
			Nonce:    42,
			HashList: hashes,
		})

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    newHash,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.ErrorIs(t, err, types.ErrMaxTxOutTrackerHashesReached)
	})

	t.Run("no hash added if already exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)

		chainID := getEthereumChainID()
		existinghHash := sample.Hash().Hex()

		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   42,
			HashList: []*types.TxHashList{
				{
					TxHash: existinghHash,
				},
			},
		})

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    existinghHash,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 42)
		require.True(t, found)
		require.Len(t, tracker.HashList, 1)
		require.EqualValues(t, existinghHash, tracker.HashList[0].TxHash)
	})

	t.Run("can add tracker with proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseObserverMock:    true,
			UseLightclientMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		ethTx, ethTxBytes, tssAddress := sample.EthTxSigned(t, chainID, sample.EthAddress(), 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, nil)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ethTxBytes, nil)

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 42)
		require.True(t, found)
		require.EqualValues(t, txHash, tracker.HashList[0].TxHash)
		require.True(t, tracker.HashList[0].Proved)
	})

	t.Run("adding existing hash with proof make it proven", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseObserverMock:    true,
			UseLightclientMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		ethTx, ethTxBytes, tssAddress := sample.EthTxSigned(t, chainID, sample.EthAddress(), 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, nil)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ethTxBytes, nil)

		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   42,
			HashList: []*types.TxHashList{
				{
					TxHash: sample.Hash().Hex(),
					Proved: false,
				},
				{
					TxHash: txHash,
					Proved: false,
				},
			},
		})

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 42)
		require.True(t, found)
		require.Len(t, tracker.HashList, 2)
		require.EqualValues(t, txHash, tracker.HashList[1].TxHash)
		require.True(t, tracker.HashList[1].Proved)
	})

	t.Run("should fail if verify proof fail", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseObserverMock:    true,
			UseLightclientMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		ethTx, ethTxBytes, _ := sample.EthTxSigned(t, chainID, sample.EthAddress(), 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ethTxBytes, errors.New("error"))

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.ErrorIs(t, err, types.ErrProofVerificationFail)
	})

	t.Run("should fail if no tss when adding hash with proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseObserverMock:    true,
			UseLightclientMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		ethTx, ethTxBytes, tssAddress := sample.EthTxSigned(t, chainID, sample.EthAddress(), 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(ethTxBytes, nil)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, errors.New("error"))

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.ErrorIs(t, err, observertypes.ErrTssNotFound)
	})

	t.Run("should fail if body verification fail with proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock:   true,
			UseObserverMock:    true,
			UseLightclientMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		chainID := getEthereumChainID()
		ethTx, _, tssAddress := sample.EthTxSigned(t, chainID, sample.EthAddress(), 42)
		txHash := ethTx.Hash().Hex()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		lightclientMock := keepertest.GetCrosschainLightclientMock(t, k)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{})
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)
		mockCctxByNonce(t, ctx, *k, observerMock, types.CctxStatus_PendingOutbound, false)
		observerMock.On("GetTssAddress", mock.Anything, mock.Anything).Return(&observertypes.QueryGetTssAddressResponse{
			Eth: tssAddress.Hex(),
		}, nil)

		// makes VerifyProof returning an invalid hash
		lightclientMock.On("VerifyProof", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sample.Bytes(), nil)

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    txHash,
			Proof:     &proofs.Proof{},
			BlockHash: "",
			TxIndex:   0,
			Nonce:     42,
		})
		require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
	})
}
