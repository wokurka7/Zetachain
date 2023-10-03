//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	etherminttypes "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ZetaTxServer is a ZetaChain tx server for smoke test
type ZetaTxServer struct {
	clientCtx      client.Context
	txFactory      tx.Factory
	Lock           sync.Mutex
	blockHeight    int64
	accountNumber  uint64
	sequenceNumber uint64
}

// NewZetaTxServer returns a new TxServer with provided account
func NewZetaTxServer(rpcAddr string, names []string, mnemonics []string) (*ZetaTxServer, error) {
	ctx := context.Background()

	if len(names) != len(mnemonics) {
		return nil, errors.New("invalid names and mnemonics")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return nil, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc)
	clientCtx := newContext(rpc, cdc, reg, kr)
	txf := newFactory(clientCtx)

	var accountNum, seqNum uint64
	// create accounts
	for i := range names {
		r, err := kr.NewAccount(names[i], mnemonics[i], "", sdktypes.FullFundraiserPath, hd.Secp256k1)
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %s", err.Error())
		}
		addr, err := r.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get account address: %s", err.Error())
		}
		//fmt.Printf(
		//	"Added account for Zeta tx server\nname: %s\nmnemonic: %s\naddress: %s\n",
		//	names[i],
		//	mnemonics[i],
		//	addr.String(),
		//)
		accountNum, seqNum, err = clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to get account number and sequence: %s", err.Error())
		}
	}

	return &ZetaTxServer{
		clientCtx:      clientCtx,
		txFactory:      txf,
		accountNumber:  accountNum,
		sequenceNumber: seqNum,
	}, nil
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg from the account
func (zts *ZetaTxServer) BroadcastTx(account string, msg sdktypes.Msg) (*sdktypes.TxResponse, error) {
	// Find number and sequence and set it
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return nil, err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return nil, err
	}
	accountNumber, accountSeq, err := zts.clientCtx.AccountRetriever.GetAccountNumberSequence(zts.clientCtx, addr)
	if err != nil {
		return nil, err
	}
	zts.txFactory = zts.txFactory.WithAccountNumber(accountNumber).WithSequence(accountSeq)

	txBuilder, err := zts.txFactory.BuildUnsignedTx(msg)
	if err != nil {
		return nil, err
	}

	// Sign tx
	err = tx.Sign(zts.txFactory, account, txBuilder, true)
	if err != nil {
		return nil, err
	}
	txBytes, err := zts.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// Broadcast tx
	return zts.clientCtx.BroadcastTx(txBytes)
}

// mimics the zetaclient broadcast to investigate the sequence mismatch issue
func (b *ZetaTxServer) BroadcastLikeZetaclient(account string, gaslimit uint64, msg sdktypes.Msg) (*sdktypes.TxResponse, error) {
	gaslimit = gaslimit * 3
	b.Lock.Lock()
	defer b.Lock.Unlock()
	var err error
	acc, err := b.clientCtx.Keyring.Key(account)
	if err != nil {
		return nil, err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return nil, err
	}

	//blockHeight, err := b.GetZetaBlockHeight()
	//if err != nil {
	//	return nil, err
	//}
	a, blockHeight, err := b.clientCtx.AccountRetriever.GetAccountWithHeight(b.clientCtx, addr)
	if err != nil {
		return nil, err
	}

	//if blockHeight > b.blockHeight {
	b.blockHeight = blockHeight
	accountNumber := a.GetAccountNumber()
	seqNumber := a.GetSequence()
	b.accountNumber = accountNumber
	//if b.sequenceNumber < seqNumber {
	fmt.Printf("[WARN] block %d Reset seq num %d => %d\n", blockHeight, b.sequenceNumber, seqNumber)
	b.sequenceNumber = seqNumber
	//}
	//}
	//b.logger.Info().Uint64("account_number", b.accountNumber).Uint64("sequence_number", b.seqNumber).Msg("account info")
	fmt.Printf("[INFO] account_number: %d, sequence_number: %d\n", b.accountNumber, b.sequenceNumber)
	b.txFactory = b.txFactory.WithAccountNumber(b.accountNumber).WithSequence(b.sequenceNumber)

	builder, err := b.txFactory.BuildUnsignedTx(msg)
	if err != nil {
		return nil, err
	}
	builder.SetGasLimit(gaslimit)
	fee := sdktypes.NewCoins(sdktypes.NewCoin("azeta", sdktypes.NewInt(40000)))
	builder.SetFeeAmount(fee)
	//fmt.Printf("signing from name: %s\n", ctx.GetFromName())
	err = tx.Sign(b.txFactory, account, builder, true)
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(builder.GetTx())

	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	commit, err := b.clientCtx.BroadcastTxSync(txBytes)
	if err != nil {
		fmt.Printf("[ERROR] BroadcastTxSync tx %s", err.Error())
		return nil, err
	}
	// Code will be the tendermint ABICode , it start at 1 , so if it is an error , code will not be zero
	if commit.Code > 0 {
		if commit.Code == 32 {
			errMsg := commit.RawLog
			re := regexp.MustCompile(`account sequence mismatch, expected ([0-9]*), got ([0-9]*)`)
			matches := re.FindStringSubmatch(errMsg)
			if len(matches) != 3 {
				return nil, err
			}
			expectedSeq, err := strconv.Atoi(matches[1])
			if err != nil {
				fmt.Printf("cannot parse expected seq %s", matches[1])
				return nil, err
			}
			gotSeq, err := strconv.Atoi(matches[2])
			if err != nil {
				fmt.Printf("cannot parse got seq %s", matches[2])
				return nil, err
			}
			b.sequenceNumber = uint64(expectedSeq)
			fmt.Printf("[WARN] Reset seq number to %d (from err msg) from %d", expectedSeq, gotSeq)
		}
		return commit, fmt.Errorf("fail to broadcast code:%d, log:%s", commit.Code, commit.RawLog)
	}

	b.sequenceNumber++
	fmt.Printf("[OK] seq number increments to %d\n", b.sequenceNumber)

	return commit, nil
}

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	etherminttypes.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(rpc *rpchttp.HTTP, cdc *codec.ProtoCodec, reg codectypes.InterfaceRegistry, kr keyring.Keyring) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(ZetaChainID).
		WithInterfaceRegistry(reg).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithBroadcastMode(flags.BroadcastSync).
		WithClient(rpc).
		WithSkipConfirmation(true).
		WithFromName("creator").
		WithFromAddress(sdktypes.AccAddress{}).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(300000).
		WithGasAdjustment(1.0).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithFees("50azeta")
}

func (zts *ZetaTxServer) GetZetaBlockHeight() (int64, error) {
	c := crosschaintypes.NewQueryClient(zts.clientCtx)
	resp, err := c.LastZetaHeight(context.Background(), &crosschaintypes.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Height, nil
}
