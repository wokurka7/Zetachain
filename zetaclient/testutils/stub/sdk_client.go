package stub

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/rpc/client/mock"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
)

type SDKClient struct {
	mock.Client
	err  error
	code uint32
}

func (c SDKClient) BroadcastTxCommit(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTxCommit, error) {
	return nil, c.err
}

func (c SDKClient) BroadcastTxAsync(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTx, error) {
	return nil, c.err
}

func (c SDKClient) BroadcastTxSync(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTx, error) {
	log := ""
	if c.err != nil {
		log = c.err.Error()
	}
	return &coretypes.ResultBroadcastTx{
		Code:      c.code,
		Data:      bytes.HexBytes{},
		Log:       log,
		Codespace: "",
		Hash:      bytes.HexBytes{},
	}, c.err
}

func (c SDKClient) Tx(_ context.Context, _ []byte, _ bool) (*coretypes.ResultTx, error) {
	return &coretypes.ResultTx{
		Hash:   bytes.HexBytes{},
		Height: 0,
		Index:  0,
		TxResult: abci.ResponseDeliverTx{
			Log: "",
		},
		Tx:    []byte{},
		Proof: tmtypes.TxProof{},
	}, c.err
}

func (c SDKClient) Block(_ context.Context, _ *int64) (*coretypes.ResultBlock, error) {
	return &coretypes.ResultBlock{Block: &tmtypes.Block{
		Header:   tmtypes.Header{},
		Data:     tmtypes.Data{},
		Evidence: tmtypes.EvidenceData{},
	}}, c.err
}

func NewSDKClientWithErr(err error, code uint32) *SDKClient {
	return &SDKClient{
		Client: mock.Client{},
		err:    err,
		code:   code,
	}
}
