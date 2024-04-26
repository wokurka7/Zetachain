package e2etests

import (
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// TODO : Add smoke test for abort refund
// https://github.com/zeta-chain/node/issues/1745
const (
	TestContextUpgradeName                = "context_upgrade"
	TestDepositAndCallRefundName          = "deposit_and_call_refund"
	TestMultipleERC20DepositName          = "erc20_multiple_deposit"
	TestMultipleWithdrawsName             = "erc20_multiple_withdraw"
	TestZetaWithdrawName                  = "zeta_withdraw"
	TestZetaWithdrawBTCRevertName         = "zeta_withdraw_btc_revert" // #nosec G101 - not a hardcoded password
	TestMessagePassingName                = "message_passing"
	TestZRC20SwapName                     = "zrc20_swap"
	TestBitcoinWithdrawSegWitName         = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName        = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawLegacyName         = "bitcoin_withdraw_legacy"
	TestBitcoinWithdrawP2WSHName          = "bitcoin_withdraw_p2wsh"
	TestBitcoinWithdrawP2SHName           = "bitcoin_withdraw_p2sh"
	TestBitcoinWithdrawInvalidAddressName = "bitcoin_withdraw_invalid"
	TestBitcoinWithdrawRestrictedName     = "bitcoin_withdraw_restricted"
	TestCrosschainSwapName                = "crosschain_swap"
	TestMessagePassingRevertFailName      = "message_passing_revert_fail"
	TestMessagePassingRevertSuccessName   = "message_passing_revert_success"
	TestMessagePassingEVMtoZEVMName       = "message_passing_zevm"
	TestMessagePassingZEVMToEVMName       = "message_passing_zevm_to_evm"

	TestMessagePassingZEVMtoEVMRevertName = "message_passing_zevm_to_evm_revert"
	TestMessagePassingEVMtoZEVMRevertName = "message_passing_revert_zevm"
	TestERC20DepositAndCallRefundName     = "erc20_deposit_and_call_refund"
	TestEtherDepositAndCallName           = "eth_deposit_and_call"
	TestDepositEtherLiquidityCapName      = "deposit_eth_liquidity_cap"
	TestMyTestName                        = "my_test"

	TestERC20WithdrawName = "erc20_withdraw"
	TestERC20DepositName  = "erc20_deposit"

	TestERC20DepositRestrictedName  = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestEtherDepositName            = "eth_deposit"
	TestEtherWithdrawName           = "eth_withdraw"
	TestEtherWithdrawRestrictedName = "eth_withdraw_restricted"
	TestBitcoinDepositName          = "bitcoin_deposit"
	TestZetaDepositName             = "zeta_deposit"
	TestZetaDepositRestrictedName   = "zeta_deposit_restricted"

	TestDonationEtherName = "donation_ether"

	TestStressEtherWithdrawName = "stress_eth_withdraw"
	TestStressBTCWithdrawName   = "stress_btc_withdraw"
	TestStressEtherDepositName  = "stress_eth_deposit"
	TestStressBTCDepositName    = "stress_btc_deposit"

	// Admin test
	TestMigrateChainSupportName     = "migrate_chain_support"
	TestPauseZRC20Name              = "pause_zrc20"
	TestUpdateBytecodeZRC20Name     = "update_bytecode_zrc20"
	TestUpdateBytecodeConnectorName = "update_bytecode_connector"
	TestRateLimiterName             = "rate_limiter"
)

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	runner.NewE2ETest(
		TestContextUpgradeName,
		"tests sending ETH on ZEVM and check context data using ContextApp",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "1000000000000000"},
		},
		TestContextUpgrade,
	),
	runner.NewE2ETest(
		TestDepositAndCallRefundName,
		"deposit ZRC20 into ZEVM and call a contract that reverts; should refund",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "10000000000000000000"},
		},
		TestDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestMultipleERC20DepositName,
		"deposit ERC20 into ZEVM in multiple deposits",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount", DefaultValue: "1000000000"},
			runner.ArgDefinition{Description: "count", DefaultValue: "3"},
		},
		TestMultipleERC20Deposit,
	),
	runner.NewE2ETest(
		TestERC20WithdrawName,
		"withdraw ERC20 from ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestMultipleWithdrawsName,
		"withdraw ERC20 from ZEVM in multiple withdrawals",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount", DefaultValue: "100"},
			runner.ArgDefinition{Description: "count", DefaultValue: "3"},
		},
		TestMultipleWithdraws,
	),
	runner.NewE2ETest(
		TestZetaWithdrawName,
		"withdraw ZETA from ZEVM to Ethereum",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestZetaWithdraw,
	),
	runner.NewE2ETest(
		TestZetaDepositName,
		"deposit ZETA from Ethereum to ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaDeposit,
	),
	runner.NewE2ETest(
		TestZetaWithdrawBTCRevertName,
		"sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaWithdrawBTCRevert,
	),
	runner.NewE2ETest(
		TestMessagePassingName,
		"evm->evm message passing (sending ZETA only)",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassing,
	),
	runner.NewE2ETest(
		TestZRC20SwapName,
		"swap ZRC20 ERC20 for ZRC20 ETH",
		[]runner.ArgDefinition{},
		TestZRC20Swap,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawSegWitName,
		"withdraw BTC from ZEVM to a SegWit address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "receiver address", DefaultValue: ""},
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawSegWit,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawTaprootName,
		"withdraw BTC from ZEVM to a Taproot address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "receiver address", DefaultValue: ""},
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawTaproot,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawLegacyName,
		"withdraw BTC from ZEVM to a legacy address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "receiver address", DefaultValue: ""},
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawLegacy,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2WSHName,
		"withdraw BTC from ZEVM to a P2WSH address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "receiver address", DefaultValue: ""},
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2WSH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2SHName,
		"withdraw BTC from ZEVM to a P2SH address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "receiver address", DefaultValue: ""},
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2SH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawInvalidAddressName,
		"withdraw BTC from ZEVM to an unsupported btc address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.00001"},
		},
		TestBitcoinWithdrawToInvalidAddress,
	),
	runner.NewE2ETest(
		TestCrosschainSwapName,
		"testing Bitcoin ERC20 cross-chain swap",
		[]runner.ArgDefinition{},
		TestCrosschainSwap,
	),
	runner.NewE2ETest(
		TestMessagePassingRevertFailName,
		"evm->evm message passing (revert fail)",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassingRevertFail,
	),
	runner.NewE2ETest(
		TestMessagePassingRevertSuccessName,
		"evm->evm message passing (revert success)",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassingRevertSuccess,
	),
	runner.NewE2ETest(
		TestPauseZRC20Name,
		"pausing ZRC20 on ZetaChain",
		[]runner.ArgDefinition{},
		TestPauseZRC20,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRefundName,
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{},
		TestERC20DepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeZRC20Name,
		"update ZRC20 bytecode swap",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeZRC20,
	),
	runner.NewE2ETest(
		TestEtherDepositAndCallName,
		"deposit ZRC20 into ZEVM and call a contract",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "1000000000000000000"},
		},
		TestEtherDepositAndCall,
	),
	runner.NewE2ETest(
		TestDepositEtherLiquidityCapName,
		"deposit Ethers into ZEVM with a liquidity cap",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000000000000"},
		},
		TestDepositEtherLiquidityCap,
	),
	runner.NewE2ETest(
		TestMyTestName,
		"performing custom test",
		[]runner.ArgDefinition{},
		TestMyTest,
	),
	runner.NewE2ETest(
		TestERC20DepositName,
		"deposit ERC20 into ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestEtherDepositName,
		"deposit Ether into ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestEtherDeposit,
	),
	runner.NewE2ETest(
		TestEtherWithdrawName,
		"withdraw Ether from ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdraw,
	),
	runner.NewE2ETest(
		TestBitcoinDepositName,
		"deposit Bitcoin into ZEVM",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDeposit,
	),
	runner.NewE2ETest(
		TestDonationEtherName,
		"donate Ether to the TSS",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000000000000000"},
		},
		TestDonationEther,
	),
	runner.NewE2ETest(
		TestStressEtherWithdrawName,
		"stress test Ether withdrawal",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000"},
			runner.ArgDefinition{Description: "count", DefaultValue: "100"},
		},
		TestStressEtherWithdraw,
	),
	runner.NewE2ETest(
		TestStressBTCWithdrawName,
		"stress test BTC withdrawal",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.01"},
			runner.ArgDefinition{Description: "count", DefaultValue: "100"},
		},
		TestStressBTCWithdraw,
	),
	runner.NewE2ETest(
		TestStressEtherDepositName,
		"stress test Ether deposit",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000"},
			runner.ArgDefinition{Description: "count", DefaultValue: "100"},
		},
		TestStressEtherDeposit,
	),
	runner.NewE2ETest(
		TestStressBTCDepositName,
		"stress test BTC deposit",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
			runner.ArgDefinition{Description: "count", DefaultValue: "100"},
		},
		TestStressBTCDeposit,
	),
	runner.NewE2ETest(
		TestZetaDepositRestrictedName,
		"deposit ZETA from Ethereum to ZEVM restricted address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaDepositRestricted,
	),
	runner.NewE2ETest(
		TestERC20DepositRestrictedName,
		"deposit ERC20 into ZEVM restricted address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositRestricted,
	),
	runner.NewE2ETest(
		TestEtherWithdrawRestrictedName,
		"withdraw Ether from ZEVM to restricted address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdrawRestricted,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawRestrictedName,
		"withdraw Bitcoin from ZEVM to restricted address",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawRestricted,
	),
	runner.NewE2ETest(
		TestMigrateChainSupportName,
		"migrate the evm chain from goerli to sepolia",
		[]runner.ArgDefinition{},
		TestMigrateChainSupport,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeConnectorName,
		"update zevm connector bytecode",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeConnector,
	),
	runner.NewE2ETest(
		TestRateLimiterName,
		"test sending cctxs with rate limiter enabled and show logs when processing cctxs",
		[]runner.ArgDefinition{},
		TestRateLimiter,
	),
	runner.NewE2ETest(
		TestMessagePassingZEVMToEVMName,
		"zevm -> evm message passing contract call",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000007"},
		},
		TestMessagePassingZEVMtoEVM,
	),
	runner.NewE2ETest(
		TestMessagePassingZEVMtoEVMRevertName,
		"zevm -> evm message passing contract call reverts",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000006"},
		},
		TestMessagePassingZEVMtoEVMRevert,
	),
	runner.NewE2ETest(
		TestMessagePassingEVMtoZEVMName,
		"evm -> zevm message passing contract call ",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000009"},
		},
		TestMessagePassingEVMtoZEVM,
	),
	runner.NewE2ETest(
		TestMessagePassingEVMtoZEVMRevertName,
		"evm -> zevm message passing and revert back to evm",
		[]runner.ArgDefinition{
			runner.ArgDefinition{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		TestMessagePassingEVMtoZEVMRevert,
	),
}
