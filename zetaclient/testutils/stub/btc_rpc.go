package stub

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
)

// EvmClient interface
var _ interfaces.BTCRPCClient = &MockBTCRPCClient{}

// MockBTCRPCClient is a mock implementation of the BTCRPCClient interface
type MockBTCRPCClient struct {
	Txs []*btcutil.Tx
	err error
}

// NewMockBTCRPCClient creates a new mock BTC RPC client
func NewMockBTCRPCClient() *MockBTCRPCClient {
	client := &MockBTCRPCClient{}
	return client.Reset()
}

// Reset clears the mock data
func (c *MockBTCRPCClient) Reset() *MockBTCRPCClient {
	c.Txs = []*btcutil.Tx{}
	return c
}

func (c *MockBTCRPCClient) GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) CreateWallet(_ string, _ ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetNewAddress(_ string) (btcutil.Address, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GenerateToAddress(_ int64, _ btcutil.Address, _ *int64) ([]*chainhash.Hash, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetBalance(_ string) (btcutil.Amount, error) {
	return 0, c.err
}

func (c *MockBTCRPCClient) SendRawTransaction(_ *wire.MsgTx, _ bool) (*chainhash.Hash, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) ListUnspentMinMaxAddresses(_ int, _ int, _ []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) EstimateSmartFee(_ int64, _ *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetTransaction(_ *chainhash.Hash) (*btcjson.GetTransactionResult, error) {
	return nil, c.err
}

// GetRawTransaction returns a pre-loaded transaction or nil
func (c *MockBTCRPCClient) GetRawTransaction(_ *chainhash.Hash) (*btcutil.Tx, error) {
	// pop a transaction from the list
	if len(c.Txs) > 0 {
		tx := c.Txs[len(c.Txs)-1]
		c.Txs = c.Txs[:len(c.Txs)-1]
		return tx, nil
	}
	return nil, c.err
}

func (c *MockBTCRPCClient) GetRawTransactionVerbose(_ *chainhash.Hash) (*btcjson.TxRawResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetBlockCount() (int64, error) {
	return 0, c.err
}

func (c *MockBTCRPCClient) GetBlockHash(_ int64) (*chainhash.Hash, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetBlockVerbose(_ *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetBlockVerboseTx(_ *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	return nil, c.err
}

func (c *MockBTCRPCClient) GetBlockHeader(_ *chainhash.Hash) (*wire.BlockHeader, error) {
	return nil, c.err
}

// ----------------------------------------------------------------------------
// Feed data to the mock BTC RPC client for testing
// ----------------------------------------------------------------------------

func (c *MockBTCRPCClient) WithRawTransaction(tx *btcutil.Tx) *MockBTCRPCClient {
	c.Txs = append(c.Txs, tx)
	return c
}

func (c *MockBTCRPCClient) WithRawTransactions(txs []*btcutil.Tx) *MockBTCRPCClient {
	c.Txs = append(c.Txs, txs...)
	return c
}

func (c *MockBTCRPCClient) WithError(err error) {
	c.err = err
}
