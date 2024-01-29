package local

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

// bitcoinTestRoutine runs Bitcoin related smoke tests
func bitcoinTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
	initBitcoinNetwork bool,
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("bitcoin panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for bitcoin test
		bitcoinRunner, err := initTestRunner(
			"bitcoin",
			conf,
			deployerRunner,
			UserBitcoinAddress,
			UserBitcoinPrivateKey,
			runner.NewLogger(verbose, color.FgYellow, "bitcoin"),
		)
		if err != nil {
			return err
		}

		bitcoinRunner.Logger.Print("🏃 starting Bitcoin tests")
		//startTime :=  time.Now()

		// funding the account
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserBitcoinAddress, 1000)
		bitcoinRunner.WaitForTxReceiptOnEvm(txUSDTSend)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := bitcoinRunner.DepositEther(false)
		txERC20Deposit := bitcoinRunner.DepositERC20()

		bitcoinRunner.WaitForMinedCCTX(txEtherDeposit)
		bitcoinRunner.WaitForMinedCCTX(txERC20Deposit)

		bitcoinRunner.SetupBitcoinAccount(initBitcoinNetwork)
		bitcoinRunner.DepositBTC(true)

		tx, err := bitcoinRunner.BTCZRC20.Approve(bitcoinRunner.ZevmAuth, bitcoinRunner.BTCZRC20Addr, big.NewInt(1e18))
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(bitcoinRunner.Ctx, bitcoinRunner.ZevmClient, tx, bitcoinRunner.Logger, bitcoinRunner.ReceiptTimeout)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}

		bitcoinRunner.MineBlocks()

		go WithdrawCCtx(bitcoinRunner)

		//// run bitcoin test
		//// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
		//// to make it faster to catch up with the latest block header
		//if err := bitcoinRunner.RunSmokeTestsFromNames(
		//	smoketests.AllSmokeTests,
		//	smoketests.TestBitcoinWithdrawName,
		//	//smoketests.TestSendZetaOutBTCRevertName,
		//	//smoketests.TestCrosschainSwapName,
		//); err != nil {
		//	return fmt.Errorf("bitcoin tests failed: %v", err)
		//}
		//
		//if err := bitcoinRunner.CheckBtcTSSBalance(); err != nil {
		//	return err
		//}
		//
		//bitcoinRunner.Logger.Print("🍾 Bitcoin tests completed in %s", time.Since(startTime).String())

		time.Sleep(100 * time.Second)

		return nil
	}
}

// TODO: Remove this part

var (
	zevmNonce = big.NewInt(1)
)

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func WithdrawCCtx(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			WithdrawBTCZRC20(sm)
		}
	}
}

func WithdrawBTCZRC20(sm *runner.SmokeTestRunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	sm.ZevmAuth.Nonce = zevmNonce

	sm.Logger.Print("nonce %d: starting withdraw", zevmNonce)
	tx, err := sm.BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), big.NewInt(100))
	if err != nil {
		panic(err)
	}

	go MonitorCCTXFromTxHash(sm, tx, zevmNonce.Int64())
}

func MonitorCCTXFromTxHash(sm *runner.SmokeTestRunner, tx *ethtypes.Transaction, nonce int64) {
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		sm.Logger.Print("nonce %d: withdraw evm tx failed", nonce)
		return
	}
	sm.Logger.Print("nonce %d: withdraw evm tx success, receipt: %+v", nonce, receipt)
	// mine 10 blocks to confirm the withdraw tx
	_, err := sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		sm.Logger.Print(
			"nonce %d: withdraw cctx failed with status %s, message %s",
			nonce,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
		)
		return
	}
	sm.Logger.Print("nonce %d: withdraw cctx success", nonce)
}
