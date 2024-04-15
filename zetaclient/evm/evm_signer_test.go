package evm

import (
	"path"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

var (
	// Dummy addresses as they are just used as transaction data to be signed
	ConnectorAddress    = sample.EthAddress()
	ERC20CustodyAddress = sample.EthAddress()
)

func getNewEvmSigner() (*Signer, error) {
	mpiAddress := ConnectorAddress
	erc20CustodyAddress := ERC20CustodyAddress
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	cfg := config.NewConfig()
	evmCfg := config.EVMConfig{
		Chain:     chains.BscMainnetChain(),
		Endpoints: []string{stub.EVMRPCEnabled},
	}
	cfg.EVMChainConfigs[chains.BscMainnetChain().ChainId] = evmCfg
	return NewEVMSigner(
		evmCfg,
		stub.NewTSSMainnet(),
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		mpiAddress,
		erc20CustodyAddress,
		corecontext.NewZetaCoreContext(cfg),
		logger,
		ts)
}

func getNewEvmChainClient() (*ChainClient, error) {
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	cfg := config.NewConfig()
	tss := stub.NewTSSMainnet()

	evmcfg := config.EVMConfig{Chain: chains.BscMainnetChain(), Endpoints: []string{"http://localhost:8545"}}
	cfg.EVMChainConfigs[chains.BscMainnetChain().ChainId] = evmcfg
	coreCTX := corecontext.NewZetaCoreContext(cfg)
	appCTX := appcontext.NewAppContext(coreCTX, cfg)

	return NewEVMChainClient(appCTX, stub.NewMockZetaCoreBridge(), tss, "", logger, evmcfg, ts)
}

func getNewOutTxProcessor() *outtxprocessor.Processor {
	logger := zerolog.Logger{}
	return outtxprocessor.NewOutTxProcessorManager(logger)
}

func getCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	var cctx crosschaintypes.CrossChainTx
	testutils.LoadObjectFromJSONFile(t, &cctx, path.Join("../", testutils.TestDataPathCctx, "cctx_56_68270.json"))
	return &cctx
}

func getInvalidCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	var cctx crosschaintypes.CrossChainTx
	testutils.LoadObjectFromJSONFile(t, &cctx, path.Join("../", testutils.TestDataPathCctx, "cctx_56_68270_invalidChainID.json"))
	return &cctx
}

func TestSigner_SetGetConnectorAddress(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ConnectorAddress, evmSigner.GetZetaConnectorAddress())

	// Update and get again
	newConnector := sample.EthAddress()
	evmSigner.SetZetaConnectorAddress(newConnector)
	require.Equal(t, newConnector, evmSigner.GetZetaConnectorAddress())
}

func TestSigner_SetGetERC20CustodyAddress(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ERC20CustodyAddress, evmSigner.GetERC20CustodyAddress())

	// Update and get again
	newCustody := sample.EthAddress()
	evmSigner.SetERC20CustodyAddress(newCustody)
	require.Equal(t, newCustody, evmSigner.GetERC20CustodyAddress())
}

func TestSigner_TryProcessOutTx(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	cctx := getCCTX(t)
	processorManager := getNewOutTxProcessor()
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)

	evmSigner.TryProcessOutTx(cctx, processorManager, "123", mockChainClient, stub.NewMockZetaCoreBridge(), 123)

	//Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	require.Len(t, *list, 1)
}

func TestSigner_SignOutboundTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct

	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignOutboundTx - should successfully sign", func(t *testing.T) {
		// Call SignOutboundTx
		tx, err := evmSigner.SignOutboundTx(txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignRevertTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignRevertTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignWithdrawTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignCommandTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCommandTx CmdWhitelistERC20", func(t *testing.T) {
		cmd := constant.CmdWhitelistERC20
		params := ConnectorAddress.Hex()
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(txData, cmd, params)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})

	t.Run("SignCommandTx CmdMigrateTssFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateTssFunds
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(txData, cmd, "")
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignERC20WithdrawTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_BroadcastOutTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("BroadcastOutTx - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutTx(tx, cctx, zerolog.Logger{}, sdktypes.AccAddress{}, stub.NewMockZetaCoreBridge(), txData)

		//Check if cctx was signed and broadcasted
		list := evmSigner.GetReportedTxList()
		require.Len(t, *list, 1)
	})
}

func TestSigner_getEVMRPC(t *testing.T) {
	t.Run("getEVMRPC error dialing", func(t *testing.T) {
		evmCfg := config.EVMConfig{
			Chain:     chains.BscMainnetChain(),
			Endpoints: []string{"invalidEndpoint"},
		}
		clientLogger := common.ClientLogger{
			Std:        zerolog.Logger{},
			Compliance: zerolog.Logger{},
		}
		client, signer, err := getEVMRPC(evmCfg, clientLogger)
		require.Nil(t, client)
		require.Nil(t, signer)
		require.Error(t, err)
	})
}

func TestSigner_SignerErrorMsg(t *testing.T) {
	cctx := getCCTX(t)

	msg := SignerErrorMsg(cctx)
	require.Contains(t, msg, "nonce 68270 chain 56")
}

func TestSigner_SignWhitelistERC20Cmd(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	txData, skip, err := NewOutBoundTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	tx, err := evmSigner.SignWhitelistERC20Cmd(txData, "")
	require.Nil(t, tx)
	require.ErrorContains(t, err, "invalid erc20 address")
}
