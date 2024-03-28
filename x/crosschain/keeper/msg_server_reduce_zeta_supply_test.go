package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	crosschainmocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/crosschain"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGetIndexForBurnCCTX(t *testing.T) {
	t.Run("get valid index", func(t *testing.T) {
		require.Equal(t, keeper.GetIndexForReduceZetaSupplyCMD(1, 1), "0xdd9b5c76369780e1d367068c34d05aa94055beb50c57475a38bbbd6656baddb2")
	})
}

func Test_GetCmdCCTX(t *testing.T) {
	t.Run("get valid cmd cctx", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		tss := sample.Tss()
		chainID := common.GoerliChain().ChainId
		amount := sample.UintInRange(1, 100)
		gasPrices := sample.UintInRange(1, 100).String()
		burnAddress := sample.EthAddress().String()
		cctx := keeper.GetCmdCCTX(ctx, chainID, amount, tss, gasPrices, burnAddress)
		require.NotNil(t, cctx)
		require.Equal(t, cctx.Creator, "")
		require.Equal(t, cctx.GetCurrentOutTxParam().ReceiverChainId, chainID)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount, amount)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasPrice, gasPrices)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxGasLimit, uint64(1_000_000))
		require.Equal(t, cctx.GetCurrentOutTxParam().Receiver, burnAddress)
		require.Equal(t, cctx.GetCurrentOutTxParam().TssPubkey, tss.TssPubkey)
	})
}

func Test_ProcessCmdCCTX(t *testing.T) {
	t.Run("fail process cmd cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		cctx := sample.CrossChainTx(t, "sample")
		nonceBefore := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
		chain := common.GoerliChain()
		cctx.GetCurrentOutTxParam().ReceiverChainId = chain.ChainId
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, chain.ChainId).
			Return(nil)

		err := k.ProcessCmdCCTX(ctx, *cctx)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, nonceBefore)
	})

	t.Run("process cmd cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		cctx := sample.CrossChainTx(t, "sample")
		chain := common.GoerliChain()
		cctx.GetCurrentOutTxParam().ReceiverChainId = chain.ChainId
		nonce := MockUpdateNonce(observerMock, chain, sample.Tss())

		err := k.ProcessCmdCCTX(ctx, *cctx)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().OutboundTxTssNonce, nonce)
	})
}

func TestKeeper_SaveCmdCCTX(t *testing.T) {
	t.Run("successfully save cmd cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "sample")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		k.SaveCmdCCTX(ctx, *cctx)

		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		intxToCctx, found := k.GetInTxHashToCctx(ctx, cctx.InboundTxParams.InboundTxObservedHash)
		require.True(t, found)
		require.Equal(t, c.GetIndex(), intxToCctx.CctxIndex[0])
		nonceToCcctx, found := zk.ObserverKeeper.GetNonceToCctx(ctx, tss.TssPubkey,
			cctx.GetCurrentOutTxParam().ReceiverChainId,
			int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce))
		require.True(t, found)
		require.Equal(t, c.GetIndex(), nonceToCcctx.CctxIndex)
	})

	t.Run("unable to save cmd cctx if tss not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "sample")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		tss := sample.Tss()

		k.SaveCmdCCTX(ctx, *cctx)

		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		intxToCctx, found := k.GetInTxHashToCctx(ctx, cctx.InboundTxParams.InboundTxObservedHash)
		require.True(t, found)
		require.Equal(t, c.GetIndex(), intxToCctx.CctxIndex[0])
		_, found = zk.ObserverKeeper.GetNonceToCctx(ctx, tss.TssPubkey,
			cctx.GetCurrentOutTxParam().ReceiverChainId,
			int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce))
		require.False(t, found)
	})
}

func TestMsgServer_ReduceZetaSupply(t *testing.T) {
	t.Run("successfully create reduce zeta supply cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		admin := observertypes.GroupID1Address
		tss := sample.Tss()

		chain := common.GoerliChain()
		amount := sample.UintInRange(1, 100)
		gasPrices := sample.UintInRange(1, 100)
		k.SetGasPrice(ctx, crosschaintypes.GasPrice{
			ChainId:     chain.ChainId,
			Prices:      []uint64{gasPrices.Uint64()},
			MedianIndex: 0,
		})

		observerMock.On("GetParams", mock.Anything).Return(observertypes.Params{
			AdminPolicy: observertypes.DefaultAdminPolicy(),
		})
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		MockUpdateNonce(observerMock, chain, tss)
		observerMock.On("SetNonceToCctx", mock.Anything, mock.Anything).Return()

		msg := crosschaintypes.MsgReduceZetaSupply{
			Creator:     admin,
			ChainId:     chain.ChainId,
			Amount:      amount,
			BurnAddress: "",
		}
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.ReduceZetaSupply(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)
		_, found := k.GetCrossChainTx(ctx, keeper.GetIndexForReduceZetaSupplyCMD(chain.ChainId, ctx.BlockHeight()))
		require.True(t, found)
	})

	t.Run("fail to  create reduce zeta supply cctx if process cmd fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		admin := observertypes.GroupID1Address
		tss := sample.Tss()

		chain := common.GoerliChain()
		amount := sample.UintInRange(1, 100)
		gasPrices := sample.UintInRange(1, 100)
		k.SetGasPrice(ctx, crosschaintypes.GasPrice{
			ChainId:     chain.ChainId,
			Prices:      []uint64{gasPrices.Uint64()},
			MedianIndex: 0,
		})

		observerMock.On("GetParams", mock.Anything).Return(observertypes.Params{
			AdminPolicy: observertypes.DefaultAdminPolicy(),
		})
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, chain.ChainId).Return(nil)

		msg := crosschaintypes.MsgReduceZetaSupply{
			Creator:     admin,
			ChainId:     chain.ChainId,
			Amount:      amount,
			BurnAddress: "",
		}
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.ReduceZetaSupply(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorContains(t, err, observertypes.ErrSupportedChains.Error())
	})

	t.Run("fail to create reduce zeta supply cctx if sender is not admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		chain := common.GoerliChain()
		amount := sample.UintInRange(1, 100)
		gasPrices := sample.UintInRange(1, 100)
		k.SetGasPrice(ctx, crosschaintypes.GasPrice{
			ChainId:     chain.ChainId,
			Prices:      []uint64{gasPrices.Uint64()},
			MedianIndex: 0,
		})

		observerMock.On("GetParams", mock.Anything).Return(observertypes.Params{
			AdminPolicy: observertypes.DefaultAdminPolicy(),
		})

		msg := crosschaintypes.MsgReduceZetaSupply{
			Creator:     sample.AccAddress(),
			ChainId:     chain.ChainId,
			Amount:      amount,
			BurnAddress: "",
		}
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.ReduceZetaSupply(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("fail to create reduce zeta supply cctx if tss is not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		admin := observertypes.GroupID1Address
		tss := sample.Tss()

		chain := common.GoerliChain()
		amount := sample.UintInRange(1, 100)
		gasPrices := sample.UintInRange(1, 100)
		k.SetGasPrice(ctx, crosschaintypes.GasPrice{
			ChainId:     chain.ChainId,
			Prices:      []uint64{gasPrices.Uint64()},
			MedianIndex: 0,
		})

		observerMock.On("GetParams", mock.Anything).Return(observertypes.Params{
			AdminPolicy: observertypes.DefaultAdminPolicy(),
		})
		observerMock.On("GetTSS", mock.Anything).Return(tss, false)

		msg := crosschaintypes.MsgReduceZetaSupply{
			Creator:     admin,
			ChainId:     chain.ChainId,
			Amount:      amount,
			BurnAddress: "",
		}
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.ReduceZetaSupply(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorContains(t, err, "cannot find current TSS")
	})

	t.Run("fail to create reduce zeta supply cctx if gasprice is not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		admin := observertypes.GroupID1Address
		tss := sample.Tss()

		chain := common.GoerliChain()
		amount := sample.UintInRange(1, 100)

		observerMock.On("GetParams", mock.Anything).Return(observertypes.Params{
			AdminPolicy: observertypes.DefaultAdminPolicy(),
		})
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)

		msg := crosschaintypes.MsgReduceZetaSupply{
			Creator:     admin,
			ChainId:     chain.ChainId,
			Amount:      amount,
			BurnAddress: "",
		}
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.ReduceZetaSupply(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, crosschaintypes.ErrUnableToGetGasPrice)
	})
}

func MockUpdateNonce(m *crosschainmocks.CrosschainObserverKeeper, senderChain common.Chain, tss observertypes.TSS) (nonce uint64) {
	nonce = uint64(1)

	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(&senderChain)
	m.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
		Return(observertypes.ChainNonces{Nonce: nonce}, true)
	m.On("GetTSS", mock.Anything).
		Return(tss, true)
	m.On("GetPendingNonces", mock.Anything, tss.TssPubkey, mock.Anything).
		Return(observertypes.PendingNonces{NonceHigh: int64(nonce)}, true)
	m.On("SetChainNonces", mock.Anything, mock.Anything)
	m.On("SetPendingNonces", mock.Anything, mock.Anything)
	return
}
