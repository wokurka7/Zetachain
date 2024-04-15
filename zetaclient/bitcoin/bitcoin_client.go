package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DynamicDepositorFeeHeight = 834500    // The starting height (Bitcoin mainnet) from which dynamic depositor fee will take effect
	maxHeightDiff             = 10000     // in case the last block is too old when the observer starts
	btcBlocksPerDay           = 144       // for LRU block cache size
	bigValueSats              = 200000000 // 2 BTC
	bigValueConfirmationCount = 6         // 6 confirmations for value >= 2 BTC
)

var _ interfaces.ChainClient = &BTCChainClient{}

type BTCLog struct {
	Chain      zerolog.Logger // The parent logger for the chain
	InTx       zerolog.Logger // The logger for incoming transactions
	OutTx      zerolog.Logger // The logger for outgoing transactions
	UTXOS      zerolog.Logger // The logger for UTXOs management
	GasPrice   zerolog.Logger // The logger for gas price
	Compliance zerolog.Logger // The logger for compliance checks
}

// BTCChainClient represents a chain configuration for Bitcoin
// Filled with above constants depending on chain
type BTCChainClient struct {
	chain            chains.Chain
	netParams        *chaincfg.Params
	rpcClient        interfaces.BTCRPCClient
	zetaClient       interfaces.ZetaCoreBridger
	Tss              interfaces.TSSSigner
	lastBlock        int64
	lastBlockScanned int64
	BlockTime        uint64 // block time in seconds

	Mu                *sync.Mutex // lock for all the maps, utxos and core params
	pendingNonce      uint64
	includedTxHashes  map[string]bool                          // key: tx hash
	includedTxResults map[string]*btcjson.GetTransactionResult // key: chain-tss-nonce
	broadcastedTx     map[string]string                        // key: chain-tss-nonce, value: outTx hash
	utxos             []btcjson.ListUnspentResult
	params            observertypes.ChainParams
	coreContext       *corecontext.ZetaCoreContext

	db     *gorm.DB
	stop   chan struct{}
	logger BTCLog
	ts     *metrics.TelemetryServer

	BlockCache *lru.Cache
}

func (ob *BTCChainClient) WithZetaClient(bridge *zetabridge.ZetaCoreBridge) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.zetaClient = bridge
}

func (ob *BTCChainClient) WithLogger(logger zerolog.Logger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.logger = BTCLog{
		Chain:    logger,
		InTx:     logger.With().Str("module", "WatchInTx").Logger(),
		OutTx:    logger.With().Str("module", "WatchOutTx").Logger(),
		UTXOS:    logger.With().Str("module", "WatchUTXOS").Logger(),
		GasPrice: logger.With().Str("module", "WatchGasPrice").Logger(),
	}
}

func (ob *BTCChainClient) WithBtcClient(client interfaces.BTCRPCClient) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.rpcClient = client
}

func (ob *BTCChainClient) WithChain(chain chains.Chain) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chain = chain
}

func (ob *BTCChainClient) SetChainParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.params = params
}

func (ob *BTCChainClient) GetChainParams() observertypes.ChainParams {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.params
}

// NewBitcoinClient returns a new configuration based on supplied target chain
func NewBitcoinClient(
	appcontext *appcontext.AppContext,
	chain chains.Chain,
	bridge interfaces.ZetaCoreBridger,
	tss interfaces.TSSSigner,
	dbpath string,
	loggers clientcommon.ClientLogger,
	btcCfg config.BTCConfig,
	ts *metrics.TelemetryServer,
) (*BTCChainClient, error) {
	ob := BTCChainClient{
		ts: ts,
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
	netParams, err := chains.BitcoinNetParamsFromChainID(ob.chain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("error getting net params for chain %d: %s", ob.chain.ChainId, err)
	}
	ob.netParams = netParams
	ob.Mu = &sync.Mutex{}
	chainLogger := loggers.Std.With().Str("chain", chain.ChainName.String()).Logger()
	ob.logger = BTCLog{
		Chain:      chainLogger,
		InTx:       chainLogger.With().Str("module", "WatchInTx").Logger(),
		OutTx:      chainLogger.With().Str("module", "WatchOutTx").Logger(),
		UTXOS:      chainLogger.With().Str("module", "WatchUTXOS").Logger(),
		GasPrice:   chainLogger.With().Str("module", "WatchGasPrice").Logger(),
		Compliance: loggers.Compliance,
	}

	ob.zetaClient = bridge
	ob.Tss = tss
	ob.coreContext = appcontext.ZetaCoreContext()
	ob.includedTxHashes = make(map[string]bool)
	ob.includedTxResults = make(map[string]*btcjson.GetTransactionResult)
	ob.broadcastedTx = make(map[string]string)
	_, chainParams, found := appcontext.ZetaCoreContext().GetBTCChainParams()
	if !found {
		return nil, fmt.Errorf("btc chains params not initialized")
	}
	ob.params = *chainParams

	// initialize the Client
	ob.rpcClient, err = NewRPCClientFallback(btcCfg, ob.logger.Chain)
	if err != nil {
		return nil, err
	}

	ob.BlockCache, err = lru.New(btcBlocksPerDay)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("failed to create bitcoin block cache")
		return nil, err
	}

	//Load btc chain client DB
	err = ob.loadDB(dbpath)
	if err != nil {
		return nil, err
	}

	return &ob, nil
}

func (ob *BTCChainClient) Start() {
	ob.logger.Chain.Info().Msgf("BitcoinChainClient is starting")
	go ob.WatchInTx()        // watch bitcoin chain for incoming txs and post votes to zetacore
	go ob.WatchOutTx()       // watch bitcoin chain for outgoing txs status
	go ob.WatchUTXOS()       // watch bitcoin chain for UTXOs owned by the TSS address
	go ob.WatchGasPrice()    // watch bitcoin chain for gas rate and post to zetacore
	go ob.WatchIntxTracker() // watch zetacore for bitcoin intx trackers
	go ob.WatchRPCStatus()   // watch the RPC status of the bitcoin chain
}

// WatchRPCStatus watches the RPC status of the Bitcoin chain
func (ob *BTCChainClient) WatchRPCStatus() {
	ob.logger.Chain.Info().Msgf("RPCStatus is starting")
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}
			bn, err := ob.rpcClient.GetBlockCount()
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			hash, err := ob.rpcClient.GetBlockHash(bn)
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			header, err := ob.rpcClient.GetBlockHeader(hash)
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			blockTime := header.Timestamp
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 1200 {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: RPC down? ")
				continue
			}
			tssAddr := ob.Tss.BTCAddressWitnessPubkeyHash()
			res, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, 1000000, []btcutil.Address{tssAddr})
			if err != nil {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: can't list utxos of TSS address; wallet or loaded? TSS address is not imported? ")
				continue
			}
			if len(res) == 0 {
				ob.logger.Chain.Error().Err(err).Msg("RPC status check: TSS address has no utxos; TSS address is not imported? ")
				continue
			}
			ob.logger.Chain.Info().Msgf("[OK] RPC status check: latest block number %d, timestamp %s (%.fs ago), tss addr %s, #utxos: %d", bn, blockTime, elapsedSeconds, tssAddr, len(res))

		case <-ob.stop:
			return
		}
	}
}

func (ob *BTCChainClient) Stop() {
	ob.logger.Chain.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop
	ob.logger.Chain.Info().Msgf("%s observer stopped", ob.chain.String())
}

func (ob *BTCChainClient) SetLastBlockHeight(height int64) {
	if height < 0 {
		panic("lastBlock is negative")
	}
	atomic.StoreInt64(&ob.lastBlock, height)
}

func (ob *BTCChainClient) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	return height
}

func (ob *BTCChainClient) SetLastBlockHeightScanned(height int64) {
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	atomic.StoreInt64(&ob.lastBlockScanned, height)
	metrics.LastScannedBlockNumber.WithLabelValues(ob.chain.ChainName.String()).Set(float64(height))
}

func (ob *BTCChainClient) GetLastBlockHeightScanned() int64 {
	height := atomic.LoadInt64(&ob.lastBlockScanned)
	if height < 0 {
		panic("lastBlockScanned is negative")
	}
	return height
}

func (ob *BTCChainClient) GetPendingNonce() uint64 {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.pendingNonce
}

// GetBaseGasPrice ...
// TODO: implement
// https://github.com/zeta-chain/node/issues/868
func (ob *BTCChainClient) GetBaseGasPrice() *big.Int {
	return big.NewInt(0)
}

// WatchInTx watches Bitcoin chain for incoming txs and post votes to zetacore
func (ob *BTCChainClient) WatchInTx() {
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchInTx", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.InTx.Error().Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	ob.logger.InTx.Info().Msgf("WatchInTx started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.InTx.Sample(&zerolog.BasicSampler{N: 10})
	for {
		select {
		case <-ticker.C():
			if !corecontext.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchInTx: inbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			err := ob.ObserveInTx()
			if err != nil {
				ob.logger.InTx.Error().Err(err).Msg("WatchInTx error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.InTx)
		case <-ob.stop:
			ob.logger.InTx.Info().Msgf("WatchInTx stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *BTCChainClient) postBlockHeader(tip int64) error {
	ob.logger.InTx.Info().Msgf("postBlockHeader: tip %d", tip)
	bn := tip
	res, err := ob.zetaClient.GetBlockHeaderStateByChain(ob.chain.ChainId)
	if err == nil && res.BlockHeaderState != nil && res.BlockHeaderState.EarliestHeight > 0 {
		bn = res.BlockHeaderState.LatestHeight + 1
	}
	if bn > tip {
		return fmt.Errorf("postBlockHeader: must post block confirmed block header: %d > %d", bn, tip)
	}
	res2, err := ob.GetBlockByNumberCached(bn)
	if err != nil {
		return fmt.Errorf("error getting bitcoin block %d: %s", bn, err)
	}

	var headerBuf bytes.Buffer
	err = res2.Header.Serialize(&headerBuf)
	if err != nil { // should never happen
		ob.logger.InTx.Error().Err(err).Msgf("error serializing bitcoin block header: %d", bn)
		return err
	}
	blockHash := res2.Header.BlockHash()
	_, err = ob.zetaClient.PostAddBlockHeader(
		ob.chain.ChainId,
		blockHash[:],
		res2.Block.Height,
		proofs.NewBitcoinHeader(headerBuf.Bytes()),
	)
	ob.logger.InTx.Info().Msgf("posted block header %d: %s", bn, blockHash)
	if err != nil { // error shouldn't block the process
		ob.logger.InTx.Error().Err(err).Msgf("error posting bitcoin block header: %d", bn)
	}
	return err
}

func (ob *BTCChainClient) ObserveInTx() error {
	// get and update latest block height
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInTxBTC: error getting block number: %s", err)
	}
	if cnt < 0 {
		return fmt.Errorf("observeInTxBTC: block number is negative: %d", cnt)
	}
	if cnt < ob.GetLastBlockHeight() {
		return fmt.Errorf("observeInTxBTC: block number should not decrease: current %d last %d", cnt, ob.GetLastBlockHeight())
	}
	ob.SetLastBlockHeight(cnt)

	// skip if current height is too low
	// #nosec G701 always in range
	confirmedBlockNum := cnt - int64(ob.GetChainParams().ConfirmationCount)
	if confirmedBlockNum < 0 {
		return fmt.Errorf("observeInTxBTC: skipping observer, current block number %d is too low", cnt)
	}

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		return nil
	}

	// query incoming gas asset to TSS address
	{
		bn := lastScanned + 1
		res, err := ob.GetBlockByNumberCached(bn)
		if err != nil {
			ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error getting bitcoin block %d", bn)
			return err
		}
		ob.logger.InTx.Info().Msgf("observeInTxBTC: block %d has %d txs, current block %d, last block %d",
			bn, len(res.Block.Tx), cnt, lastScanned)

		// print some debug information
		if len(res.Block.Tx) > 1 {
			for idx, tx := range res.Block.Tx {
				ob.logger.InTx.Debug().Msgf("BTC InTX |  %d: %s\n", idx, tx.Txid)
				for vidx, vout := range tx.Vout {
					ob.logger.InTx.Debug().Msgf("vout %d \n value: %v\n scriptPubKey: %v\n", vidx, vout.Value, vout.ScriptPubKey.Hex)
				}
			}
		}

		// add block header to zetabridge
		// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
		// https://github.com/zeta-chain/node/issues/1847
		flags := ob.coreContext.GetCrossChainFlags()
		if flags.BlockHeaderVerificationFlags != nil && flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
			err = ob.postBlockHeader(bn)
			if err != nil {
				ob.logger.InTx.Warn().Err(err).Msgf("observeInTxBTC: error posting block header %d", bn)
			}
		}

		if len(res.Block.Tx) > 1 {
			// get depositor fee
			depositorFee := CalcDepositorFee(res.Block, ob.chain.ChainId, ob.netParams, ob.logger.InTx)

			// filter incoming txs to TSS address
			tssAddress := ob.Tss.BTCAddress()
			// #nosec G701 always positive
			inTxs, err := FilterAndParseIncomingTx(
				ob.rpcClient,
				res.Block.Tx,
				uint64(res.Block.Height),
				tssAddress,
				ob.logger.InTx,
				ob.netParams,
				depositorFee,
			)
			if err != nil {
				ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error filtering incoming txs for block %d", bn)
				return err // we have to re-scan this block next time
			}

			// post inbound vote message to zetacore
			for _, inTx := range inTxs {
				msg := ob.GetInboundVoteMessageFromBtcEvent(inTx)
				if msg != nil {
					zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
					if err != nil {
						ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error posting to zeta core for tx %s", inTx.TxHash)
						return err // we have to re-scan this block next time
					} else if zetaHash != "" {
						ob.logger.InTx.Info().Msgf("observeInTxBTC: PostVoteInbound zeta tx hash: %s inTx %s ballot %s fee %v",
							zetaHash, inTx.TxHash, ballot, depositorFee)
					}
				}
			}
		}

		// Save LastBlockHeight
		ob.SetLastBlockHeightScanned(bn)
		// #nosec G701 always positive
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(uint64(bn))).Error; err != nil {
			ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error writing last scanned block %d to db", bn)
		}
	}

	return nil
}

// ConfirmationsThreshold returns number of required Bitcoin confirmations depending on sent BTC amount.
func (ob *BTCChainClient) ConfirmationsThreshold(amount *big.Int) int64 {
	if amount.Cmp(big.NewInt(bigValueSats)) >= 0 {
		return bigValueConfirmationCount
	}
	if bigValueConfirmationCount < ob.GetChainParams().ConfirmationCount {
		return bigValueConfirmationCount
	}
	// #nosec G701 always in range
	return int64(ob.GetChainParams().ConfirmationCount)
}

// IsSendOutTxProcessed returns isIncluded(or inMempool), isConfirmed, Error
func (ob *BTCChainClient) IsSendOutTxProcessed(cctx *types.CrossChainTx, logger zerolog.Logger) (bool, bool, error) {
	params := *cctx.GetCurrentOutTxParam()
	sendHash := cctx.Index
	nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce

	// get broadcasted outtx and tx result
	outTxID := ob.GetTxID(nonce)
	logger.Info().Msgf("IsSendOutTxProcessed %s", outTxID)

	ob.Mu.Lock()
	txnHash, broadcasted := ob.broadcastedTx[outTxID]
	res, included := ob.includedTxResults[outTxID]
	ob.Mu.Unlock()

	if !included {
		if !broadcasted {
			return false, false, nil
		}
		// If the broadcasted outTx is nonce 0, just wait for inclusion and don't schedule more keysign
		// Schedule more than one keysign for nonce 0 can lead to duplicate payments.
		// One purpose of nonce mark UTXO is to avoid duplicate payment based on the fact that Bitcoin
		// prevents double spending of same UTXO. However, for nonce 0, we don't have a prior nonce (e.g., -1)
		// for the signer to check against when making the payment. Signer treats nonce 0 as a special case in downstream code.
		if nonce == 0 {
			return true, false, nil
		}

		// Try including this outTx broadcasted by myself
		txResult, inMempool := ob.checkIncludedTx(cctx, txnHash)
		if txResult == nil { // check failed, try again next time
			return false, false, nil
		} else if inMempool { // still in mempool (should avoid unnecessary Tss keysign)
			ob.logger.OutTx.Info().Msgf("IsSendOutTxProcessed: outTx %s is still in mempool", outTxID)
			return true, false, nil
		}
		// included
		ob.setIncludedTx(nonce, txResult)

		// Get tx result again in case it is just included
		res = ob.getIncludedTx(nonce)
		if res == nil {
			return false, false, nil
		}
		ob.logger.OutTx.Info().Msgf("IsSendOutTxProcessed: setIncludedTx succeeded for outTx %s", outTxID)
	}

	// It's safe to use cctx's amount to post confirmation because it has already been verified in observeOutTx()
	amountInSat := params.Amount.BigInt()
	if res.Confirmations < ob.ConfirmationsThreshold(amountInSat) {
		return true, false, nil
	}

	logger.Debug().Msgf("Bitcoin outTx confirmed: txid %s, amount %s\n", res.TxID, amountInSat.String())
	zetaHash, ballot, err := ob.zetaClient.PostVoteOutbound(
		sendHash,
		res.TxID,
		// #nosec G701 always positive
		uint64(res.BlockIndex),
		0,   // gas used not used with Bitcoin
		nil, // gas price not used with Bitcoin
		0,   // gas limit not used with Bitcoin
		amountInSat,
		chains.ReceiveStatus_Success,
		ob.chain,
		nonce,
		coin.CoinType_Gas,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("IsSendOutTxProcessed: error confirming bitcoin outTx %s, nonce %d ballot %s", res.TxID, nonce, ballot)
	} else if zetaHash != "" {
		logger.Info().Msgf("IsSendOutTxProcessed: confirmed Bitcoin outTx %s, zeta tx hash %s nonce %d ballot %s", res.TxID, zetaHash, nonce, ballot)
	}
	return true, true, nil
}

// WatchGasPrice watches Bitcoin chain for gas rate and post to zetacore
func (ob *BTCChainClient) WatchGasPrice() {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice()
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchGasPrice", ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.GasPrice.Error().Err(err).Msg("error creating ticker")
		return
	}
	ob.logger.GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.chain.ChainId, ob.GetChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.chain.ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.GasPrice)
		case <-ob.stop:
			ob.logger.GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *BTCChainClient) PostGasPrice() error {
	if ob.chain.ChainId == 18444 { //bitcoin regtest; hardcode here since this RPC is not available on regtest
		bn, err := ob.rpcClient.GetBlockCount()
		if err != nil {
			return err
		}
		// #nosec G701 always in range
		zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, 1, "100", uint64(bn))
		if err != nil {
			ob.logger.GasPrice.Err(err).Msg("PostGasPrice:")
			return err
		}
		_ = zetaHash
		//ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)
		return nil
	}
	// EstimateSmartFee returns the fees per kilobyte (BTC/kb) targeting given block confirmation
	feeResult, err := ob.rpcClient.EstimateSmartFee(1, &btcjson.EstimateModeEconomical)
	if err != nil {
		return err
	}
	if feeResult.Errors != nil || feeResult.FeeRate == nil {
		return fmt.Errorf("error getting gas price: %s", feeResult.Errors)
	}
	if *feeResult.FeeRate > math.MaxInt64 {
		return fmt.Errorf("gas price is too large: %f", *feeResult.FeeRate)
	}
	feeRatePerByte := FeeRateToSatPerByte(*feeResult.FeeRate)
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}
	// #nosec G701 always positive
	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, feeRatePerByte.Uint64(), "100", uint64(bn))
	if err != nil {
		ob.logger.GasPrice.Err(err).Msg("PostGasPrice:")
		return err
	}
	_ = zetaHash
	return nil
}

type BTCInTxEvnet struct {
	FromAddress string  // the first input address
	ToAddress   string  // some TSS address
	Value       float64 // in BTC, not satoshi
	MemoBytes   []byte
	BlockNumber uint64
	TxHash      string
}

// FilterAndParseIncomingTx given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(
	rpcClient interfaces.BTCRPCClient,
	txs []btcjson.TxRawResult,
	blockNumber uint64,
	tssAddress string,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) ([]*BTCInTxEvnet, error) {
	inTxs := make([]*BTCInTxEvnet, 0)
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}
		inTx, err := GetBtcEvent(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}
		if inTx != nil {
			inTxs = append(inTxs, inTx)
			logger.Info().Msgf("FilterAndParseIncomingTx: found btc event for tx %s in block %d", tx.Txid, blockNumber)
		}
	}
	return inTxs, nil
}

func (ob *BTCChainClient) GetInboundVoteMessageFromBtcEvent(inTx *BTCInTxEvnet) *types.MsgVoteOnObservedInboundTx {
	ob.logger.InTx.Debug().Msgf("Processing inTx: %s", inTx.TxHash)
	amount := big.NewFloat(inTx.Value)
	amount = amount.Mul(amount, big.NewFloat(1e8))
	amountInt, _ := amount.Int(nil)
	message := hex.EncodeToString(inTx.MemoBytes)

	// compliance check
	if ob.IsInTxRestricted(inTx) {
		return nil
	}

	return zetabridge.GetInBoundVoteMessage(
		inTx.FromAddress,
		ob.chain.ChainId,
		inTx.FromAddress,
		inTx.FromAddress,
		ob.zetaClient.ZetaChain().ChainId,
		cosmosmath.NewUintFromBigInt(amountInt),
		message,
		inTx.TxHash,
		inTx.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		0,
	)
}

// IsInTxRestricted returns true if the inTx contains restricted addresses
func (ob *BTCChainClient) IsInTxRestricted(inTx *BTCInTxEvnet) bool {
	receiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(inTx.MemoBytes))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(inTx.FromAddress, receiver) {
		compliance.PrintComplianceLog(ob.logger.InTx, ob.logger.Compliance,
			false, ob.chain.ChainId, inTx.TxHash, inTx.FromAddress, receiver, "BTC")
		return true
	}
	return false
}

// GetBtcEvent either returns a valid BTCInTxEvent or nil
// Note: the caller should retry the tx on error (e.g., GetSenderAddressByVin failed)
func GetBtcEvent(
	rpcClient interfaces.BTCRPCClient,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) (*BTCInTxEvnet, error) {
	found := false
	var value float64
	var memo []byte
	if len(tx.Vout) >= 2 {
		// 1st vout must have tss address as receiver with p2wpkh scriptPubKey
		vout0 := tx.Vout[0]
		script := vout0.ScriptPubKey.Hex
		if len(script) == 44 && script[:4] == "0014" { // P2WPKH output: 0x00 + 20 bytes of pubkey hash
			receiver, err := DecodeScriptP2WPKH(vout0.ScriptPubKey.Hex, netParams)
			if err != nil { // should never happen
				return nil, err
			}
			// skip irrelevant tx to us
			if receiver != tssAddress {
				return nil, nil
			}
			// deposit amount has to be no less than the minimum depositor fee
			if vout0.Value < depositorFee {
				logger.Info().Msgf("GetBtcEvent: btc deposit amount %v in txid %s is less than depositor fee %v", vout0.Value, tx.Txid, depositorFee)
				return nil, nil
			}
			value = vout0.Value - depositorFee

			// 2nd vout must be a valid OP_RETURN memo
			vout1 := tx.Vout[1]
			memo, found, err = DecodeOpReturnMemo(vout1.ScriptPubKey.Hex, tx.Txid)
			if err != nil {
				logger.Error().Err(err).Msgf("GetBtcEvent: error decoding OP_RETURN memo: %s", vout1.ScriptPubKey.Hex)
				return nil, nil
			}
		}
	}
	// event found, get sender address
	if found {
		if len(tx.Vin) == 0 { // should never happen
			return nil, fmt.Errorf("GetBtcEvent: no input found for intx: %s", tx.Txid)
		}
		fromAddress, err := GetSenderAddressByVin(rpcClient, tx.Vin[0], netParams)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting sender address for intx: %s", tx.Txid)
		}
		return &BTCInTxEvnet{
			FromAddress: fromAddress,
			ToAddress:   tssAddress,
			Value:       value,
			MemoBytes:   memo,
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}, nil
	}
	return nil, nil
}

// GetSenderAddressByVin get the sender address from the previous transaction
func GetSenderAddressByVin(rpcClient interfaces.BTCRPCClient, vin btcjson.Vin, net *chaincfg.Params) (string, error) {
	// query previous raw transaction by txid
	// GetTransaction requires reconfiguring the bitcoin node (txindex=1), so we use GetRawTransaction instead
	hash, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		return "", err
	}
	tx, err := rpcClient.GetRawTransaction(hash)
	if err != nil {
		return "", errors.Wrapf(err, "error getting raw transaction %s", vin.Txid)
	}
	// #nosec G701 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript
	scriptHex := hex.EncodeToString(pkScript)
	if IsPkScriptP2TR(pkScript) {
		return DecodeScriptP2TR(scriptHex, net)
	}
	if IsPkScriptP2WSH(pkScript) {
		return DecodeScriptP2WSH(scriptHex, net)
	}
	if IsPkScriptP2WPKH(pkScript) {
		return DecodeScriptP2WPKH(scriptHex, net)
	}
	if IsPkScriptP2SH(pkScript) {
		return DecodeScriptP2SH(scriptHex, net)
	}
	if IsPkScriptP2PKH(pkScript) {
		return DecodeScriptP2PKH(scriptHex, net)
	}
	// sender address not found, return nil and move on to the next tx
	return "", nil
}

// WatchUTXOS watches bitcoin chain for UTXOs owned by the TSS address
func (ob *BTCChainClient) WatchUTXOS() {
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchUTXOS", ob.GetChainParams().WatchUtxoTicker)
	if err != nil {
		ob.logger.UTXOS.Error().Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.FetchUTXOS()
			if err != nil {
				ob.logger.UTXOS.Error().Err(err).Msg("error fetching btc utxos")
			}
			ticker.UpdateInterval(ob.GetChainParams().WatchUtxoTicker, ob.logger.UTXOS)
		case <-ob.stop:
			ob.logger.UTXOS.Info().Msgf("WatchUTXOS stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *BTCChainClient) FetchUTXOS() error {
	defer func() {
		if err := recover(); err != nil {
			ob.logger.UTXOS.Error().Msgf("BTC fetchUTXOS: caught panic error: %v", err)
		}
	}()

	// This is useful when a zetaclient's pending nonce lagged behind for whatever reason.
	ob.refreshPendingNonce()

	// get the current block height.
	bh, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("btc: error getting block height : %v", err)
	}
	maxConfirmations := int(bh)

	// List all unspent UTXOs (160ms)
	tssAddr := ob.Tss.BTCAddress()
	address, err := chains.DecodeBtcAddress(tssAddr, ob.chain.ChainId)
	if err != nil {
		return fmt.Errorf("btc: error decoding wallet address (%s) : %s", tssAddr, err.Error())
	}
	utxos, err := ob.rpcClient.ListUnspentMinMaxAddresses(0, maxConfirmations, []btcutil.Address{address})
	if err != nil {
		return err
	}

	// rigid sort to make utxo list deterministic
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Amount == utxos[j].Amount {
			if utxos[i].TxID == utxos[j].TxID {
				return utxos[i].Vout < utxos[j].Vout
			}
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Amount < utxos[j].Amount
	})

	// filter UTXOs good to spend for next TSS transaction
	utxosFiltered := make([]btcjson.ListUnspentResult, 0)
	for _, utxo := range utxos {
		// UTXOs big enough to cover the cost of spending themselves
		if utxo.Amount < DefaultDepositorFee {
			continue
		}
		// we don't want to spend other people's unconfirmed UTXOs as they may not be safe to spend
		if utxo.Confirmations == 0 {
			if !ob.isTssTransaction(utxo.TxID) {
				continue
			}
		}
		utxosFiltered = append(utxosFiltered, utxo)
	}

	ob.Mu.Lock()
	metrics.NumberOfUTXO.Set(float64(len(utxosFiltered)))
	ob.utxos = utxosFiltered
	ob.Mu.Unlock()
	return nil
}

// isTssTransaction checks if a given transaction was sent by TSS itself.
// An unconfirmed transaction is safe to spend only if it was sent by TSS and verified by ourselves.
func (ob *BTCChainClient) isTssTransaction(txid string) bool {
	_, found := ob.includedTxHashes[txid]
	return found
}

// refreshPendingNonce tries increasing the artificial pending nonce of outTx (if lagged behind).
// There could be many (unpredictable) reasons for a pending nonce lagging behind, for example:
// 1. The zetaclient gets restarted.
// 2. The tracker is missing in zetabridge.
func (ob *BTCChainClient) refreshPendingNonce() {
	// get pending nonces from zetabridge
	p, err := ob.zetaClient.GetPendingNoncesByChain(ob.chain.ChainId)
	if err != nil {
		ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting pending nonces")
	}

	// increase pending nonce if lagged behind
	ob.Mu.Lock()
	pendingNonce := ob.pendingNonce
	ob.Mu.Unlock()

	// #nosec G701 always non-negative
	nonceLow := uint64(p.NonceLow)
	if nonceLow > pendingNonce {
		// get the last included outTx hash
		txid, err := ob.getOutTxidByNonce(nonceLow-1, false)
		if err != nil {
			ob.logger.Chain.Error().Err(err).Msg("refreshPendingNonce: error getting last outTx txid")
		}

		// set 'NonceLow' as the new pending nonce
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		ob.pendingNonce = nonceLow
		ob.logger.Chain.Info().Msgf("refreshPendingNonce: increase pending nonce to %d with txid %s", ob.pendingNonce, txid)
	}
}

func (ob *BTCChainClient) getOutTxidByNonce(nonce uint64, test bool) (string, error) {

	// There are 2 types of txids an observer can trust
	// 1. The ones had been verified and saved by observer self.
	// 2. The ones had been finalized in zetabridge based on majority vote.
	if res := ob.getIncludedTx(nonce); res != nil {
		return res.TxID, nil
	}
	if !test { // if not unit test, get cctx from zetabridge
		send, err := ob.zetaClient.GetCctxByNonce(ob.chain.ChainId, nonce)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting cctx for nonce %d", nonce)
		}
		txid := send.GetCurrentOutTxParam().OutboundTxHash
		if txid == "" {
			return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
		}
		// make sure it's a real Bitcoin txid
		_, getTxResult, err := GetTxResultByHash(ob.rpcClient, txid)
		if err != nil {
			return "", errors.Wrapf(err, "getOutTxidByNonce: error getting outTx result for nonce %d hash %s", nonce, txid)
		}
		if getTxResult.Confirmations <= 0 { // just a double check
			return "", fmt.Errorf("getOutTxidByNonce: outTx txid %s for nonce %d is not included", txid, nonce)
		}
		return txid, nil
	}
	return "", fmt.Errorf("getOutTxidByNonce: cannot find outTx txid for nonce %d", nonce)
}

func (ob *BTCChainClient) findNonceMarkUTXO(nonce uint64, txid string) (int, error) {
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	amount := chains.NonceMarkAmount(nonce)
	for i, utxo := range ob.utxos {
		sats, err := GetSatoshis(utxo.Amount)
		if err != nil {
			ob.logger.OutTx.Error().Err(err).Msgf("findNonceMarkUTXO: error getting satoshis for utxo %v", utxo)
		}
		if utxo.Address == tssAddress && sats == amount && utxo.TxID == txid && utxo.Vout == 0 {
			ob.logger.OutTx.Info().Msgf("findNonceMarkUTXO: found nonce-mark utxo with txid %s, amount %d satoshi", utxo.TxID, sats)
			return i, nil
		}
	}
	return -1, fmt.Errorf("findNonceMarkUTXO: cannot find nonce-mark utxo with nonce %d", nonce)
}

// SelectUTXOs selects a sublist of utxos to be used as inputs.
//
// Parameters:
//   - amount: The desired minimum total value of the selected UTXOs.
//   - utxos2Spend: The maximum number of UTXOs to spend.
//   - nonce: The nonce of the outbound transaction.
//   - consolidateRank: The rank below which UTXOs will be consolidated.
//   - test: true for unit test only.
//
// Returns:
//   - a sublist (includes previous nonce-mark) of UTXOs or an error if the qualifying sublist cannot be found.
//   - the total value of the selected UTXOs.
//   - the number of consolidated UTXOs.
//   - the total value of the consolidated UTXOs.
func (ob *BTCChainClient) SelectUTXOs(amount float64, utxosToSpend uint16, nonce uint64, consolidateRank uint16, test bool) ([]btcjson.ListUnspentResult, float64, uint16, float64, error) {
	idx := -1
	if nonce == 0 {
		// for nonce = 0; make exception; no need to include nonce-mark utxo
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
	} else {
		// for nonce > 0; we proceed only when we see the nonce-mark utxo
		preTxid, err := ob.getOutTxidByNonce(nonce-1, test)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		ob.Mu.Lock()
		defer ob.Mu.Unlock()
		idx, err = ob.findNonceMarkUTXO(nonce-1, preTxid)
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// select smallest possible UTXOs to make payment
	total := 0.0
	left, right := 0, 0
	for total < amount && right < len(ob.utxos) {
		if utxosToSpend > 0 { // expand sublist
			total += ob.utxos[right].Amount
			right++
			utxosToSpend--
		} else { // pop the smallest utxo and append the current one
			total -= ob.utxos[left].Amount
			total += ob.utxos[right].Amount
			left++
			right++
		}
	}
	results := make([]btcjson.ListUnspentResult, right-left)
	copy(results, ob.utxos[left:right])

	// include nonce-mark as the 1st input
	if idx >= 0 { // for nonce > 0
		if idx < left || idx >= right {
			total += ob.utxos[idx].Amount
			results = append([]btcjson.ListUnspentResult{ob.utxos[idx]}, results...)
		} else { // move nonce-mark to left
			for i := idx - left; i > 0; i-- {
				results[i], results[i-1] = results[i-1], results[i]
			}
		}
	}
	if total < amount {
		return nil, 0, 0, 0, fmt.Errorf("SelectUTXOs: not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}

	// consolidate biggest possible UTXOs to maximize consolidated value
	// consolidation happens only when there are more than (or equal to) consolidateRank (10) UTXOs
	utxoRank, consolidatedUtxo, consolidatedValue := uint16(0), uint16(0), 0.0
	for i := len(ob.utxos) - 1; i >= 0 && utxosToSpend > 0; i-- { // iterate over UTXOs big-to-small
		if i != idx && (i < left || i >= right) { // exclude nonce-mark and already selected UTXOs
			utxoRank++
			if utxoRank >= consolidateRank { // consolication starts from the 10-ranked UTXO based on value
				utxosToSpend--
				consolidatedUtxo++
				total += ob.utxos[i].Amount
				consolidatedValue += ob.utxos[i].Amount
				results = append(results, ob.utxos[i])
			}
		}
	}

	return results, total, consolidatedUtxo, consolidatedValue, nil
}

// SaveBroadcastedTx saves successfully broadcasted transaction
func (ob *BTCChainClient) SaveBroadcastedTx(txHash string, nonce uint64) {
	outTxID := ob.GetTxID(nonce)
	ob.Mu.Lock()
	ob.broadcastedTx[outTxID] = txHash
	ob.Mu.Unlock()

	broadcastEntry := clienttypes.ToOutTxHashSQLType(txHash, outTxID)
	if err := ob.db.Save(&broadcastEntry).Error; err != nil {
		ob.logger.OutTx.Error().Err(err).Msgf("SaveBroadcastedTx: error saving broadcasted txHash %s for outTx %s", txHash, outTxID)
	}
	ob.logger.OutTx.Info().Msgf("SaveBroadcastedTx: saved broadcasted txHash %s for outTx %s", txHash, outTxID)
}

// WatchOutTx watches Bitcoin chain for outgoing txs status
func (ob *BTCChainClient) WatchOutTx() {
	ticker, err := clienttypes.NewDynamicTicker("Bitcoin_WatchOutTx", ob.GetChainParams().OutTxTicker)
	if err != nil {
		ob.logger.OutTx.Error().Err(err).Msg("error creating ticker ")
		return
	}

	defer ticker.Stop()
	ob.logger.OutTx.Info().Msgf("WatchInTx started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.OutTx.Sample(&zerolog.BasicSampler{N: 10})
	for {
		select {
		case <-ticker.C():
			if !corecontext.IsOutboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutTx: outbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain.ChainId, interfaces.Ascending)
			if err != nil {
				ob.logger.OutTx.Error().Err(err).Msgf("WatchOutTx: error GetAllOutTxTrackerByChain for chain %d", ob.chain.ChainId)
				continue
			}
			for _, tracker := range trackers {
				// get original cctx parameters
				outTxID := ob.GetTxID(tracker.Nonce)
				cctx, err := ob.zetaClient.GetCctxByNonce(ob.chain.ChainId, tracker.Nonce)
				if err != nil {
					ob.logger.OutTx.Info().Err(err).Msgf("WatchOutTx: can't find cctx for chain %d nonce %d", ob.chain.ChainId, tracker.Nonce)
					break
				}
				nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
				if tracker.Nonce != nonce { // Tanmay: it doesn't hurt to check
					ob.logger.OutTx.Error().Msgf("WatchOutTx: tracker nonce %d not match cctx nonce %d", tracker.Nonce, nonce)
					break
				}
				if len(tracker.HashList) > 1 {
					ob.logger.OutTx.Warn().Msgf("WatchOutTx: oops, outTxID %s got multiple (%d) outTx hashes", outTxID, len(tracker.HashList))
				}
				// iterate over all txHashes to find the truly included one.
				// we do it this (inefficient) way because we don't rely on the first one as it may be a false positive (for unknown reason).
				txCount := 0
				var txResult *btcjson.GetTransactionResult
				for _, txHash := range tracker.HashList {
					result, inMempool := ob.checkIncludedTx(cctx, txHash.TxHash)
					if result != nil && !inMempool { // included
						txCount++
						txResult = result
						ob.logger.OutTx.Info().Msgf("WatchOutTx: included outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, tracker.Nonce)
						if txCount > 1 {
							ob.logger.OutTx.Error().Msgf(
								"WatchOutTx: checkIncludedTx passed, txCount %d chain %d nonce %d result %v", txCount, ob.chain.ChainId, tracker.Nonce, result)
						}
					}
				}
				if txCount == 1 { // should be only one txHash included for each nonce
					ob.setIncludedTx(tracker.Nonce, txResult)
				} else if txCount > 1 {
					ob.removeIncludedTx(tracker.Nonce) // we can't tell which txHash is true, so we remove all (if any) to be safe
					ob.logger.OutTx.Error().Msgf("WatchOutTx: included multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, tracker.Nonce)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutTxTicker, ob.logger.OutTx)
		case <-ob.stop:
			ob.logger.OutTx.Info().Msgf("WatchOutTx stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// checkIncludedTx checks if a txHash is included and returns (txResult, inMempool)
// Note: if txResult is nil, then inMempool flag should be ignored.
func (ob *BTCChainClient) checkIncludedTx(cctx *types.CrossChainTx, txHash string) (*btcjson.GetTransactionResult, bool) {
	outTxID := ob.GetTxID(cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
	hash, getTxResult, err := GetTxResultByHash(ob.rpcClient, txHash)
	if err != nil {
		ob.logger.OutTx.Error().Err(err).Msgf("checkIncludedTx: error GetTxResultByHash: %s", txHash)
		return nil, false
	}
	if txHash != getTxResult.TxID { // just in case, we'll use getTxResult.TxID later
		ob.logger.OutTx.Error().Msgf("checkIncludedTx: inconsistent txHash %s and getTxResult.TxID %s", txHash, getTxResult.TxID)
		return nil, false
	}
	if getTxResult.Confirmations >= 0 { // check included tx only
		err = ob.checkTssOutTxResult(cctx, hash, getTxResult)
		if err != nil {
			ob.logger.OutTx.Error().Err(err).Msgf("checkIncludedTx: error verify bitcoin outTx %s outTxID %s", txHash, outTxID)
			return nil, false
		}
		return getTxResult, false // included
	}
	return getTxResult, true // in mempool
}

// setIncludedTx saves included tx result in memory
func (ob *BTCChainClient) setIncludedTx(nonce uint64, getTxResult *btcjson.GetTransactionResult) {
	txHash := getTxResult.TxID
	outTxID := ob.GetTxID(nonce)

	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	res, found := ob.includedTxResults[outTxID]

	if !found { // not found.
		ob.includedTxHashes[txHash] = true
		ob.includedTxResults[outTxID] = getTxResult // include new outTx and enforce rigid 1-to-1 mapping: nonce <===> txHash
		if nonce >= ob.pendingNonce {               // try increasing pending nonce on every newly included outTx
			ob.pendingNonce = nonce + 1
		}
		ob.logger.OutTx.Info().Msgf("setIncludedTx: included new bitcoin outTx %s outTxID %s pending nonce %d", txHash, outTxID, ob.pendingNonce)
	} else if txHash == res.TxID { // found same hash.
		ob.includedTxResults[outTxID] = getTxResult // update tx result as confirmations may increase
		if getTxResult.Confirmations > res.Confirmations {
			ob.logger.OutTx.Info().Msgf("setIncludedTx: bitcoin outTx %s got confirmations %d", txHash, getTxResult.Confirmations)
		}
	} else { // found other hash.
		// be alert for duplicate payment!!! As we got a new hash paying same cctx (for whatever reason).
		delete(ob.includedTxResults, outTxID) // we can't tell which txHash is true, so we remove all to be safe
		ob.logger.OutTx.Error().Msgf("setIncludedTx: duplicate payment by bitcoin outTx %s outTxID %s, prior outTx %s", txHash, outTxID, res.TxID)
	}
}

// getIncludedTx gets the receipt and transaction from memory
func (ob *BTCChainClient) getIncludedTx(nonce uint64) *btcjson.GetTransactionResult {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.includedTxResults[ob.GetTxID(nonce)]
}

// removeIncludedTx removes included tx from memory
func (ob *BTCChainClient) removeIncludedTx(nonce uint64) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	txResult, found := ob.includedTxResults[ob.GetTxID(nonce)]
	if found {
		delete(ob.includedTxHashes, txResult.TxID)
		delete(ob.includedTxResults, ob.GetTxID(nonce))
	}
}

// Basic TSS outTX checks:
//   - should be able to query the raw tx
//   - check if all inputs are segwit && TSS inputs
//
// Returns: true if outTx passes basic checks.
func (ob *BTCChainClient) checkTssOutTxResult(cctx *types.CrossChainTx, hash *chainhash.Hash, res *btcjson.GetTransactionResult) error {
	params := cctx.GetCurrentOutTxParam()
	nonce := params.OutboundTxTssNonce
	rawResult, err := GetRawTxResult(ob.rpcClient, hash, res)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: error GetRawTxResultByHash %s", hash.String())
	}
	err = ob.checkTSSVin(rawResult.Vin, nonce)
	if err != nil {
		return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vin in outTx %s nonce %d", hash, nonce)
	}

	// differentiate between normal and restricted cctx
	if compliance.IsCctxRestricted(cctx) {
		err = ob.checkTSSVoutCancelled(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vout in cancelled outTx %s nonce %d", hash, nonce)
		}
	} else {
		err = ob.checkTSSVout(params, rawResult.Vout)
		if err != nil {
			return errors.Wrapf(err, "checkTssOutTxResult: invalid TSS Vout in outTx %s nonce %d", hash, nonce)
		}
	}
	return nil
}

// GetTxResultByHash gets the transaction result by hash
func GetTxResultByHash(rpcClient interfaces.BTCRPCClient, txID string) (*chainhash.Hash, *btcjson.GetTransactionResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTxResultByHash: error NewHashFromStr: %s", txID)
	}

	// The Bitcoin node has to be configured to watch TSS address
	txResult, err := rpcClient.GetTransaction(hash)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetOutTxByTxHash: error GetTransaction %s", hash.String())
	}
	return hash, txResult, nil
}

// GetRawTxResult gets the raw tx result
func GetRawTxResult(rpcClient interfaces.BTCRPCClient, hash *chainhash.Hash, res *btcjson.GetTransactionResult) (btcjson.TxRawResult, error) {
	if res.Confirmations == 0 { // for pending tx, we query the raw tx directly
		rawResult, err := rpcClient.GetRawTransactionVerbose(hash) // for pending tx, we query the raw tx
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetRawTransactionVerbose %s", res.TxID)
		}
		return *rawResult, nil
	} else if res.Confirmations > 0 { // for confirmed tx, we query the block
		blkHash, err := chainhash.NewHashFromStr(res.BlockHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error NewHashFromStr for block hash %s", res.BlockHash)
		}
		block, err := rpcClient.GetBlockVerboseTx(blkHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: error GetBlockVerboseTx %s", res.BlockHash)
		}
		if res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx)) {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "getRawTxResult: invalid outTx with invalid block index, TxID %s, BlockIndex %d", res.TxID, res.BlockIndex)
		}
		return block.Tx[res.BlockIndex], nil
	}

	// res.Confirmations < 0 (meaning not included)
	return btcjson.TxRawResult{}, fmt.Errorf("getRawTxResult: tx %s not included yet", hash)
}

// checkTSSVin checks vin is valid if:
//   - The first input is the nonce-mark
//   - All inputs are from TSS address
func (ob *BTCChainClient) checkTSSVin(vins []btcjson.Vin, nonce uint64) error {
	// vins: [nonce-mark, UTXO1, UTXO2, ...]
	if nonce > 0 && len(vins) <= 1 {
		return fmt.Errorf("checkTSSVin: len(vins) <= 1")
	}
	pubKeyTss := hex.EncodeToString(ob.Tss.PubKeyCompressedBytes())
	for i, vin := range vins {
		// The length of the Witness should be always 2 for SegWit inputs.
		if len(vin.Witness) != 2 {
			return fmt.Errorf("checkTSSVin: expected 2 witness items, got %d", len(vin.Witness))
		}
		if vin.Witness[1] != pubKeyTss {
			return fmt.Errorf("checkTSSVin: witness pubkey %s not match TSS pubkey %s", vin.Witness[1], pubKeyTss)
		}
		// 1st vin: nonce-mark MUST come from prior TSS outTx
		if nonce > 0 && i == 0 {
			preTxid, err := ob.getOutTxidByNonce(nonce-1, false)
			if err != nil {
				return fmt.Errorf("checkTSSVin: error findTxIDByNonce %d", nonce-1)
			}
			// nonce-mark MUST the 1st output that comes from prior TSS outTx
			if vin.Txid != preTxid || vin.Vout != 0 {
				return fmt.Errorf("checkTSSVin: invalid nonce-mark txid %s vout %d, expected txid %s vout 0", vin.Txid, vin.Vout, preTxid)
			}
		}
	}
	return nil
}

// checkTSSVout vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the correct payment to recipient
//   - The third output is the change to TSS (optional)
func (ob *BTCChainClient) checkTSSVout(params *types.OutboundTxParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, payment to recipient, change to TSS (optional)]
	if !(len(vouts) == 2 || len(vouts) == 3) {
		return fmt.Errorf("checkTSSVout: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.OutboundTxTssNonce
	tssAddress := ob.Tss.BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverExpected := tssAddress
		if vout.N == 1 {
			// the 2nd output is the payment to recipient
			receiverExpected = params.Receiver
		}
		receiverVout, amount, err := DecodeTSSVout(vout, receiverExpected, ob.chain)
		if err != nil {
			return err
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVout: nonce-mark address %s not match TSS address %s", receiverVout, tssAddress)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVout: nonce-mark amount %d not match nonce-mark amount %d", amount, chains.NonceMarkAmount(nonce))
			}
		case 1: // 2nd vout: payment to recipient
			if receiverVout != params.Receiver {
				return fmt.Errorf("checkTSSVout: output address %s not match params receiver %s", receiverVout, params.Receiver)
			}
			// #nosec G701 always positive
			if uint64(amount) != params.Amount.Uint64() {
				return fmt.Errorf("checkTSSVout: output amount %d not match params amount %d", amount, params.Amount)
			}
		case 2: // 3rd vout: change to TSS (optional)
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVout: change address %s not match TSS address %s", receiverVout, tssAddress)
			}
		}
	}
	return nil
}

// checkTSSVoutCancelled vout is valid if:
//   - The first output is the nonce-mark
//   - The second output is the change to TSS (optional)
func (ob *BTCChainClient) checkTSSVoutCancelled(params *types.OutboundTxParams, vouts []btcjson.Vout) error {
	// vouts: [nonce-mark, change to TSS (optional)]
	if !(len(vouts) == 1 || len(vouts) == 2) {
		return fmt.Errorf("checkTSSVoutCancelled: invalid number of vouts: %d", len(vouts))
	}

	nonce := params.OutboundTxTssNonce
	tssAddress := ob.Tss.BTCAddress()
	for _, vout := range vouts {
		// decode receiver and amount from vout
		receiverVout, amount, err := DecodeTSSVout(vout, tssAddress, ob.chain)
		if err != nil {
			return errors.Wrap(err, "checkTSSVoutCancelled: error decoding P2WPKH vout")
		}
		switch vout.N {
		case 0: // 1st vout: nonce-mark
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVoutCancelled: nonce-mark address %s not match TSS address %s", receiverVout, tssAddress)
			}
			if amount != chains.NonceMarkAmount(nonce) {
				return fmt.Errorf("checkTSSVoutCancelled: nonce-mark amount %d not match nonce-mark amount %d", amount, chains.NonceMarkAmount(nonce))
			}
		case 1: // 2nd vout: change to TSS (optional)
			if receiverVout != tssAddress {
				return fmt.Errorf("checkTSSVoutCancelled: change address %s not match TSS address %s", receiverVout, tssAddress)
			}
		}
	}
	return nil
}

func (ob *BTCChainClient) BuildBroadcastedTxMap() error {
	var broadcastedTransactions []clienttypes.OutTxHashSQLType
	if err := ob.db.Find(&broadcastedTransactions).Error; err != nil {
		ob.logger.Chain.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, entry := range broadcastedTransactions {
		ob.broadcastedTx[entry.Key] = entry.Hash
	}
	return nil
}

func (ob *BTCChainClient) LoadLastBlock() error {
	bn, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return err
	}

	//Load persisted block number
	var lastBlockNum clienttypes.LastBlockSQLType
	if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
		ob.logger.Chain.Info().Msg("LastBlockNum not found in DB, scan from latest")
		ob.SetLastBlockHeightScanned(bn)
	} else {
		// #nosec G701 always in range
		lastBN := int64(lastBlockNum.Num)
		ob.SetLastBlockHeightScanned(lastBN)

		//If persisted block number is too low, use the latest height
		if (bn - lastBN) > maxHeightDiff {
			ob.logger.Chain.Info().Msgf("LastBlockNum too low: %d, scan from latest", lastBlockNum.Num)
			ob.SetLastBlockHeightScanned(bn)
		}
	}

	if ob.chain.ChainId == 18444 { // bitcoin regtest: start from block 100
		ob.SetLastBlockHeightScanned(100)
	}
	ob.logger.Chain.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

	return nil
}

func (ob *BTCChainClient) loadDB(dbpath string) error {
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		err := os.MkdirAll(dbpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	path := fmt.Sprintf("%s/btc_chain_client", dbpath)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to connect database")
	}
	ob.db = db

	err = db.AutoMigrate(&clienttypes.TransactionResultSQLType{},
		&clienttypes.OutTxHashSQLType{},
		&clienttypes.LastBlockSQLType{})
	if err != nil {
		return err
	}

	//Load last block
	err = ob.LoadLastBlock()
	if err != nil {
		return err
	}

	//Load broadcasted transactions
	err = ob.BuildBroadcastedTxMap()

	return err
}

func (ob *BTCChainClient) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.BTCAddress()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

type BTCBlockNHeader struct {
	Header *wire.BlockHeader
	Block  *btcjson.GetBlockVerboseTxResult
}

func (ob *BTCChainClient) GetBlockByNumberCached(blockNumber int64) (*BTCBlockNHeader, error) {
	if result, ok := ob.BlockCache.Get(blockNumber); ok {
		return result.(*BTCBlockNHeader), nil
	}
	// Get the block hash
	hash, err := ob.rpcClient.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	// Get the block header
	header, err := ob.rpcClient.GetBlockHeader(hash)
	if err != nil {
		return nil, err
	}
	// Get the block with verbose transactions
	block, err := ob.rpcClient.GetBlockVerboseTx(hash)
	if err != nil {
		return nil, err
	}
	blockNheader := &BTCBlockNHeader{
		Header: header,
		Block:  block,
	}
	ob.BlockCache.Add(blockNumber, blockNheader)
	ob.BlockCache.Add(hash, blockNheader)
	return blockNheader, nil
}
