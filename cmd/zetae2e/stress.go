package main

import (
	"context"
	"errors"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/cmd/zetae2e/local"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const (
	StatInterval      = 5
	StressTestTimeout = 100 * time.Minute
)

var (
	zevmNonce = big.NewInt(1)
)

type stressArguments struct {
	deployerAddress    string
	deployerPrivateKey string
	network            string
	txnInterval        int64
	contractsDeployed  bool
	config             string
}

var stressTestArgs = stressArguments{}

func NewStressTestCmd() *cobra.Command {
	var StressCmd = &cobra.Command{
		Use:   "stress",
		Short: "Run Stress Test",
		Run:   StressTest,
	}

	StressCmd.Flags().StringVar(&stressTestArgs.deployerAddress, "addr", "0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC", "--addr <eth address>")
	StressCmd.Flags().StringVar(&stressTestArgs.deployerPrivateKey, "privKey", "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263", "--privKey <eth private key>")
	StressCmd.Flags().StringVar(&stressTestArgs.network, "network", "LOCAL", "--network TESTNET")
	StressCmd.Flags().Int64Var(&stressTestArgs.txnInterval, "tx-interval", 500, "--tx-interval [TIME_INTERVAL_MILLISECONDS]")
	StressCmd.Flags().BoolVar(&stressTestArgs.contractsDeployed, "contracts-deployed", false, "--contracts-deployed=false")
	StressCmd.Flags().StringVar(&stressTestArgs.config, local.FlagConfigFile, "", "config file to use for the smoketest")
	StressCmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

	local.DeployerAddress = ethcommon.HexToAddress(stressTestArgs.deployerAddress)

	return StressCmd
}

func StressTest(cmd *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(StressTestTimeout)
		fmt.Println("Smoke test timed out after", StressTestTimeout)
		os.Exit(1)
	}()

	// set account prefix to zeta
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()

	// initialize smoke tests config
	conf, err := local.GetConfig(cmd)
	if err != nil {
		panic(err)
	}

	// -----------------------------------------------------------------------------------

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	logger := runner.NewLogger(verbose, color.FgHiBlue, "setup")

	if stressTestArgs.network == "LOCAL" {
		logger.Print("⏳ wait 70s for genesis")
		time.Sleep(70 * time.Second)
	}

	// initialize smoke test runner
	smokeTest, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		local.DeployerAddress,
		local.DeployerPrivateKey,
		utils.FungibleAdminName,
		FungibleAdminMnemonic,
		logger,
	)
	if err != nil {
		panic(err)
	}

	// wait for keygen
	local.WaitKeygenHeight(ctx, smokeTest.CctxClient, logger)

	// setup TSS addresses
	smokeTest.SetTSSAddresses()

	smokeTest.SetupEVM(stressTestArgs.contractsDeployed)

	// If stress test is running on local docker environment
	if stressTestArgs.network == "LOCAL" {
		// deploy and set zevm contract
		smokeTest.SetZEVMContracts()

		// deposit on ZetaChain
		smokeTest.DepositEther(false)
		smokeTest.DepositZeta()
		smokeTest.SetupBitcoinAccount(true)
		smokeTest.DepositBTC(false)

	} else if stressTestArgs.network == "TESTNET" {
		ethZRC20Addr, err := smokeTest.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(5))
		if err != nil {
			panic(err)
		}
		smokeTest.ETHZRC20Addr = ethZRC20Addr
		smokeTest.ETHZRC20, err = zrc20.NewZRC20(smokeTest.ETHZRC20Addr, smokeTest.ZevmClient)
		if err != nil {
			panic(err)
		}
	} else {
		err := errors.New("invalid network argument: " + stressTestArgs.network)
		panic(err)
	}

	// Check zrc20 balance of Deployer address
	ethZRC20Balance, err := smokeTest.ETHZRC20.BalanceOf(nil, local.DeployerAddress)
	if err != nil {
		panic(err)
	}
	fmt.Printf("eth zrc20 balance: %s Wei \n", ethZRC20Balance.String())

	//Pre-approve ETH withdraw on ZEVM
	fmt.Printf("approving ZRC20...\n")
	tx, err := smokeTest.ETHZRC20.Approve(smokeTest.ZevmAuth, smokeTest.ETHZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(ctx, smokeTest.ZevmClient, tx, logger, smokeTest.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	tx, err = smokeTest.BTCZRC20.Approve(smokeTest.ZevmAuth, smokeTest.BTCZRC20Addr, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(ctx, smokeTest.ZevmClient, tx, logger, smokeTest.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}

	// Get current nonce on zevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum, err := smokeTest.ZevmClient.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	// #nosec G701 smoketest - always in range
	nonce, err := smokeTest.ZevmClient.NonceAt(context.Background(), local.DeployerAddress, big.NewInt(int64(blockNum)))
	if err != nil {
		panic(err)
	}

	// #nosec G701 smoketest - always in range
	zevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from ZEVM to EVM - goerli")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	smokeTest.WG.Add(2)
	go WithdrawCCtx(smokeTest)       // Withdraw USDT from ZEVM to EVM - goerli
	go EchoNetworkMetrics(smokeTest) // Display Network metrics periodically to monitor performance

	smokeTest.WG.Wait()
}

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func WithdrawCCtx(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(stressTestArgs.txnInterval))
	for {
		select {
		case <-ticker.C:
			WithdrawBTCZRC20(sm)
		}
	}
}

func EchoNetworkMetrics(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Second * StatInterval)
	var queue = make([]uint64, 0)
	var numTicks = 0
	var totalMinedTxns = uint64(0)
	var previousMinedTxns = uint64(0)
	chainID, err := getChainID(sm.GoerliClient)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			cctxResp, err := sm.CctxClient.CctxListPending(context.Background(), &crosschaintypes.QueryListCctxPendingRequest{
				ChainId: chainID.Int64(),
			})
			if err != nil {
				continue
			}
			sends := cctxResp.CrossChainTx
			sort.Slice(sends, func(i, j int) bool {
				return sends[i].GetCurrentOutTxParam().OutboundTxTssNonce < sends[j].GetCurrentOutTxParam().OutboundTxTssNonce
			})
			if len(sends) > 0 {
				fmt.Printf("pending nonces %d to %d\n", sends[0].GetCurrentOutTxParam().OutboundTxTssNonce, sends[len(sends)-1].GetCurrentOutTxParam().OutboundTxTssNonce)
			} else {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := sm.CctxClient.OutTxTrackerAll(context.Background(), &crosschaintypes.QueryAllOutTxTrackerRequest{})
			if err != nil {
				continue
			}

			currentMinedTxns := sends[0].GetCurrentOutTxParam().OutboundTxTssNonce
			newMinedTxCnt := currentMinedTxns - previousMinedTxns
			previousMinedTxns = currentMinedTxns

			// Add new mined txn count to queue and remove the oldest entry
			queue = append(queue, newMinedTxCnt)
			if numTicks > 60/StatInterval {
				totalMinedTxns -= queue[0]
				queue = queue[1:]
				numTicks = 60/StatInterval + 1 //prevent overflow
			}

			//Calculate rate -> tx/min
			totalMinedTxns += queue[len(queue)-1]
			rate := totalMinedTxns

			numPending := len(cctxResp.CrossChainTx)
			numTrackers := len(trackerResp.OutTxTracker)

			fmt.Println("Network Stat => Num of Pending cctx: ", numPending, "Num active trackers: ", numTrackers, "Tx Rate: ", rate, " tx/min")
		}
	}
}

func WithdrawETHZRC20(sm *runner.SmokeTestRunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	sm.ZevmAuth.Nonce = zevmNonce
	_, err := sm.ETHZRC20.Withdraw(sm.ZevmAuth, local.DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
}

func WithdrawBTCZRC20(sm *runner.SmokeTestRunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	sm.ZevmAuth.Nonce = zevmNonce
	tx, err := sm.BTCZRC20.Withdraw(sm.ZevmAuth, local.DeployerAddress.Bytes(), big.NewInt(100))
	if err != nil {
		panic(err)
	}
	MonitorCCTXFromTxHash(sm, tx, zevmNonce.Int64())
}

func MonitorCCTXFromTxHash(sm *runner.SmokeTestRunner, tx *ethtypes.Transaction, nonce int64) {
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		sm.Logger.Info("nonce %d: withdraw evm tx failed", nonce)
		return
	}
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		sm.Logger.Info(
			"nonce %d: withdraw cctx failed with status %s, message %s",
			nonce,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
		)
		return
	}
	sm.Logger.Info("nonce %d: withdraw cctx success", nonce)
}

// Get ETH based chain ID
func getChainID(client *ethclient.Client) (*big.Int, error) {
	return client.ChainID(context.Background())
}
