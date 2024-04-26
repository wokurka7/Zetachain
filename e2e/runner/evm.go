package runner

import (
	"log"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/pkg/proofs/ethereum"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
)

var blockHeaderETHTimeout = 5 * time.Minute

// WaitForTxReceiptOnEvm waits for a tx receipt on EVM
func (runner *E2ERunner) WaitForTxReceiptOnEvm(tx *ethtypes.Transaction) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

// MintERC20OnEvm mints ERC20 on EVM
// amount is a multiple of 1e18
func (runner *E2ERunner) MintERC20OnEvm(amountERC20 int64) {
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	tx, err := runner.ERC20.Mint(runner.EVMAuth, amount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("mint failed")
	}
	runner.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendERC20OnEvm sends ERC20 to an address on EVM
// this allows the ERC20 contract deployer to funds other accounts on EVM
// amountERC20 is a multiple of 1e18
func (runner *E2ERunner) SendERC20OnEvm(address ethcommon.Address, amountERC20 int64) *ethtypes.Transaction {
	// the deployer might be sending ERC20 in different goroutines
	defer func() {
		runner.Unlock()
	}()
	runner.Lock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	// transfer
	tx, err := runner.ERC20.Transfer(runner.EVMAuth, address, amount)
	if err != nil {
		panic(err)
	}
	return tx
}

func (runner *E2ERunner) DepositERC20() ethcommon.Hash {
	runner.Logger.Print("⏳ depositing ERC20 into ZEVM")

	return runner.DepositERC20WithAmountAndMessage(runner.DeployerAddress, big.NewInt(1e18), []byte{})
}

func (runner *E2ERunner) DepositERC20WithAmountAndMessage(to ethcommon.Address, amount *big.Int, msg []byte) ethcommon.Hash {
	// reset allowance, necessary for USDT
	tx, err := runner.ERC20.Approve(runner.EVMAuth, runner.ERC20CustodyAddr, big.NewInt(0))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	runner.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = runner.ERC20.Approve(runner.EVMAuth, runner.ERC20CustodyAddr, amount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	runner.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	tx, err = runner.ERC20Custody.Deposit(runner.EVMAuth, to.Bytes(), runner.ERC20Addr, amount, msg)
	runner.Logger.Info("TX: %v", tx)

	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, tx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	runner.Logger.Info("Deposit receipt tx hash: %s, status %d", receipt.TxHash.Hex(), receipt.Status)
	for _, log := range receipt.Logs {
		event, err := runner.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		runner.Logger.Info("Deposited event:")
		runner.Logger.Info("  Recipient address: %x", event.Recipient)
		runner.Logger.Info("  ERC20 address: %s", event.Asset.Hex())
		runner.Logger.Info("  Amount: %d", event.Amount)
		runner.Logger.Info("  Message: %x", event.Message)
	}
	return tx.Hash()
}

// DepositEther sends Ethers into ZEVM
func (runner *E2ERunner) DepositEther(testHeader bool) ethcommon.Hash {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 eth
	return runner.DepositEtherWithAmount(testHeader, amount)
}

// DepositEtherWithAmount sends Ethers into ZEVM
func (runner *E2ERunner) DepositEtherWithAmount(testHeader bool, amount *big.Int) ethcommon.Hash {
	runner.Logger.Print("⏳ depositing Ethers into ZEVM")

	signedTx, err := runner.SendEther(runner.TSSAddress, amount, nil)
	if err != nil {
		panic(err)
	}
	runner.Logger.EVMTransaction(*signedTx, "send to TSS")

	receipt := utils.MustWaitForTxReceipt(runner.Ctx, runner.EVMClient, signedTx, runner.Logger, runner.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("deposit failed")
	}
	runner.Logger.EVMReceipt(*receipt, "send to TSS")

	// due to the high block throughput in localnet, ZetaClient might catch up slowly with the blocks
	// to optimize block header proof test, this test is directly executed here on the first deposit instead of having a separate test
	if testHeader {
		runner.ProveEthTransaction(receipt)
	}

	return signedTx.Hash()
}

// SendEther sends ethers to the TSS on EVM
func (runner *E2ERunner) SendEther(_ ethcommon.Address, value *big.Int, data []byte) (*ethtypes.Transaction, error) {
	evmClient := runner.EVMClient

	nonce, err := evmClient.PendingNonceAt(runner.Ctx, runner.DeployerAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(30000) // in units
	gasPrice, err := evmClient.SuggestGasPrice(runner.Ctx)
	if err != nil {
		return nil, err
	}

	tx := ethtypes.NewTransaction(nonce, runner.TSSAddress, value, gasLimit, gasPrice, data)
	chainID, err := evmClient.ChainID(runner.Ctx)
	if err != nil {
		return nil, err
	}

	deployerPrivkey, err := crypto.HexToECDSA(runner.DeployerPrivateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), deployerPrivkey)
	if err != nil {
		return nil, err
	}
	err = evmClient.SendTransaction(runner.Ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// ProveEthTransaction proves an ETH transaction on ZetaChain
func (runner *E2ERunner) ProveEthTransaction(receipt *ethtypes.Receipt) {
	startTime := time.Now()

	txHash := receipt.TxHash
	blockHash := receipt.BlockHash

	// #nosec G701 test - always in range
	txIndex := int(receipt.TransactionIndex)

	block, err := runner.EVMClient.BlockByHash(runner.Ctx, blockHash)
	if err != nil {
		panic(err)
	}
	for {
		// check timeout
		if time.Since(startTime) > blockHeaderETHTimeout {
			panic("timeout waiting for block header")
		}

		_, err := runner.LightclientClient.BlockHeader(runner.Ctx, &lightclienttypes.QueryGetBlockHeaderRequest{
			BlockHash: blockHash.Bytes(),
		})
		if err != nil {
			runner.Logger.Info("WARN: block header not found; retrying... error: %s", err.Error())
		} else {
			runner.Logger.Info("OK: block header found")
			break
		}

		time.Sleep(2 * time.Second)
	}

	trie := ethereum.NewTrie(block.Transactions())
	if trie.Hash() != block.Header().TxHash {
		panic("tx root hash & block tx root mismatch")
	}
	txProof, err := trie.GenerateProof(txIndex)
	if err != nil {
		panic("error generating txProof")
	}
	val, err := txProof.Verify(block.TxHash(), txIndex)
	if err != nil {
		panic("error verifying txProof")
	}
	var txx ethtypes.Transaction
	err = txx.UnmarshalBinary(val)
	if err != nil {
		panic("error unmarshalling txProof'd tx")
	}
	res, err := runner.LightclientClient.Prove(runner.Ctx, &lightclienttypes.QueryProveRequest{
		BlockHash: blockHash.Hex(),
		TxIndex:   int64(txIndex),
		TxHash:    txHash.Hex(),
		Proof:     proofs.NewEthereumProof(txProof),
		ChainId:   chains.GoerliLocalnetChain().ChainId,
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof invalid") // FIXME: don't do this in production
	}
	runner.Logger.Info("OK: txProof verified")
}

// AnvilMineBlocks mines blocks on Anvil localnet
// the block time is provided in seconds
// the method returns a function to stop the mining
func (runner *E2ERunner) AnvilMineBlocks(url string, blockTime int) (func(), error) {
	stop := make(chan struct{})

	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				time.Sleep(time.Duration(blockTime) * time.Second)

				var result interface{}
				err = client.CallContext(runner.Ctx, &result, "evm_mine")
				if err != nil {
					log.Fatalf("Failed to mine a new block: %v", err)
				}
			}
		}
	}()
	return func() {
		close(stop)
	}, nil
}
