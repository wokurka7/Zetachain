package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "crosschain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_metacore"

	ProtocolFee = 2000000000000000000
)

func GetProtocolFee() sdk.Uint {
	return sdk.NewUint(ProtocolFee)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	TxinKey              = "Txin-value-"
	TxinVoterKey         = "TxinVoter-value-"
	NodeAccountKey       = "NodeAccount-value-"
	TxoutKey             = "Txout-value-"
	TxoutCountKey        = "Txout-count-"
	TxoutConfirmationKey = "TxoutConfirmation-value-"
	SendKey              = "Send-value-"
	VoteCounterKey       = "VoteCounter-value-"
	ReceiveKey           = "Receive-value-"
	LastBlockHeightKey   = "LastBlockHeight-value-"
	ChainNoncesKey       = "ChainNonces-value-"
	GasPriceKey          = "GasPrice-value-"

	GasBalanceKey         = "GasBalance-value-"
	TSSKey                = "TSS-value-"
	TSSVoterKey           = "TSSVoter-value-"
	KeygenKey             = "Keygen-value-"
	OutTxTrackerKeyPrefix = "OutTxTracker/value/"
)

// OutTxTrackerKey returns the store key to retrieve a OutTxTracker from the index fields
func OutTxTrackerKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// TODO: what's the purpose of this log identifier?
func (cctx CrossChainTx) LogIdentifierForCCTX() string {
	if len(cctx.OutboundTxParams) == 0 {
		return fmt.Sprintf("%s-%d", cctx.InboundTxParams.Sender, cctx.InboundTxParams.SenderChainId)
	} else {
		i := len(cctx.OutboundTxParams) - 1
		outTx := cctx.OutboundTxParams[i]
		return fmt.Sprintf("%s-%d-%d-%d", cctx.InboundTxParams.Sender, cctx.InboundTxParams.SenderChainId, outTx.ReceiverChainId, outTx.OutboundTxTssNonce)
	}
}

const (
	// this admin can do some cleanup services like adding/removing outTxTracker entries
	AdminKey = "zeta1rx9r8hff0adaqhr5tuadkzj4e7ns2ntg446vtt"
)

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
	//ModuleAddressEVM common.EVMAddress
	ModuleAddressEVM = common.BytesToAddress(ModuleAddress.Bytes())
	//0xB73C0Aac4C1E606C6E495d848196355e6CB30381
)

const (
	PermissionFlagsKey = "PermissionFlags-value-"
)
