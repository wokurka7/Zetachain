package evm

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
)

var _ interfaces.EthClientFallback = &EthClientFallback{}

// EthClientFallback is a decorator combining client interfaces used by evm chains. Also encapsulates list of clients
// defined by endpoints from the config.
type EthClientFallback struct {
	evmCfg         config.EVMConfig
	ethClients     *common.ClientQueue
	jsonRPCClients *common.ClientQueue
	logger         zerolog.Logger
}

// NewEthClientFallback creates new instance of eth client used by evm chain client.
func NewEthClientFallback(evmCfg config.EVMConfig, logger zerolog.Logger) (*EthClientFallback, error) {
	if len(evmCfg.Endpoints) == 0 {
		return nil, errors.New("invalid endpoints")
	}
	ethClientFallback := EthClientFallback{}
	ethClientFallback.ethClients = common.NewClientQueue()
	ethClientFallback.jsonRPCClients = common.NewClientQueue()

	// Initialize clients
	for _, endpoint := range evmCfg.Endpoints {
		//Initialize go-ethereum clients
		client, err := ethclient.Dial(endpoint)
		if err != nil {
			logger.Error().Err(err).Msg("eth Client Dial")
			return nil, err
		}
		ethClientFallback.ethClients.Append(client)

		//Initialize jsonRPC clients from https://github.com/onrik/ethrpc
		jsonRPCClient := ethrpc.NewEthRPC(endpoint)
		ethClientFallback.jsonRPCClients.Append(jsonRPCClient)
	}

	ethClientFallback.evmCfg = evmCfg
	ethClientFallback.logger = logger

	return &ethClientFallback, nil
}

// The following functions are wrappers for EVMRPCClient interface. The logic is similar for all functions, the first client
// in the queue is used to attempt the rpc call. If this fails then it will attempt to call the next client in the list
// until it is successful or returns an error when the list of clients have been exhausted.

func (e *EthClientFallback) CodeAt(ctx context.Context, contract ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	var res []byte
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).CodeAt(ctx, contract, blockNumber)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var res []byte
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).CallContract(ctx, call, blockNumber)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) HeaderByNumber(ctx context.Context, number *big.Int) (*ethtypes.Header, error) {
	var res *ethtypes.Header
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).HeaderByNumber(ctx, number)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) PendingCodeAt(ctx context.Context, account ethcommon.Address) ([]byte, error) {
	var res []byte
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).PendingCodeAt(ctx, account)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) PendingNonceAt(ctx context.Context, account ethcommon.Address) (uint64, error) {
	var res uint64
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).PendingNonceAt(ctx, account)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).SuggestGasPrice(ctx)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).SuggestGasTipCap(ctx)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			gas, err = client.(interfaces.EVMRPCClient).EstimateGas(ctx, call)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return
}

func (e *EthClientFallback) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]ethtypes.Log, error) {
	var res []ethtypes.Log
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).FilterLogs(ctx, query)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- ethtypes.Log) (ethereum.Subscription, error) {
	var res ethereum.Subscription
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).SubscribeFilterLogs(ctx, query, ch)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) SendTransaction(ctx context.Context, tx *ethtypes.Transaction) error {
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			err = client.(interfaces.EVMRPCClient).SendTransaction(ctx, tx)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return err
}

func (e *EthClientFallback) BlockNumber(ctx context.Context) (uint64, error) {
	var res uint64
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).BlockNumber(ctx)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) BlockByNumber(ctx context.Context, number *big.Int) (*ethtypes.Block, error) {
	var res *ethtypes.Block
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).BlockByNumber(ctx, number)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) TransactionByHash(ctx context.Context, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			tx, isPending, err = client.(interfaces.EVMRPCClient).TransactionByHash(ctx, hash)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return
}

func (e *EthClientFallback) TransactionReceipt(ctx context.Context, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	var res *ethtypes.Receipt
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).TransactionReceipt(ctx, txHash)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) TransactionSender(ctx context.Context, tx *ethtypes.Transaction, block ethcommon.Hash, index uint) (ethcommon.Address, error) {
	var res ethcommon.Address
	var err error

	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			res, err = client.(interfaces.EVMRPCClient).TransactionSender(ctx, tx, block, index)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.ethClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (e *EthClientFallback) ChainID(ctx context.Context) (res *big.Int, err error) {
	for i := 0; i < e.ethClients.Length(); i++ {
		if client := e.ethClients.First(); client != nil {
			rpcClient := client.(*ethclient.Client)
			res, err = rpcClient.ChainID(ctx)
		}
		if err != nil {
			e.ethClients.Next()
			continue
		}
		break
	}
	return
}

// Implementation of interface for jsonRPC eth client - https://github.com/onrik/ethrpc

// EthGetBlockByNumber implementation of interfaces.EVMJSONRPCClient
func (e *EthClientFallback) EthGetBlockByNumber(number int, withTransactions bool) (*ethrpc.Block, error) {
	var res *ethrpc.Block
	var err error

	for i := 0; i < e.jsonRPCClients.Length(); i++ {
		if client := e.jsonRPCClients.First(); client != nil {
			res, err = client.(interfaces.EVMJSONRPCClient).EthGetBlockByNumber(number, withTransactions)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.jsonRPCClients.Next()
			continue
		}
		break
	}
	return res, err
}

// EthGetTransactionByHash implementation of interfaces.EVMJSONRPCClient
func (e *EthClientFallback) EthGetTransactionByHash(hash string) (*ethrpc.Transaction, error) {
	var res *ethrpc.Transaction
	var err error

	for i := 0; i < e.jsonRPCClients.Length(); i++ {
		if client := e.jsonRPCClients.First(); client != nil {
			res, err = client.(interfaces.EVMJSONRPCClient).EthGetTransactionByHash(hash)
		}
		if err != nil {
			e.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			e.jsonRPCClients.Next()
			continue
		}
		break
	}
	return res, err
}
